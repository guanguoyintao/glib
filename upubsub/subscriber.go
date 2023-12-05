package upubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uconcurrent"
	"git.umu.work/AI/uglib/umetadata"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/logger"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type SubscriberWorker struct {
	mutex  sync.RWMutex
	once   sync.Once
	mq     MQ
	router map[string]*subscriberWorker
	ms     []SubscriberWrapper
}

func NewSubscriberWorker(mq MQ) *SubscriberWorker {
	return &SubscriberWorker{
		mq:     mq,
		mutex:  sync.RWMutex{},
		once:   sync.Once{},
		router: make(map[string]*subscriberWorker),
		ms:     make([]SubscriberWrapper, 0),
	}
}

func (w *SubscriberWorker) RegisterSubscriber(ctx context.Context, topic, group string, concurrency int,
	maxRetries int, maxDelay time.Duration, handler SubscriberFunc) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.router[fmt.Sprintf("%s-%s", topic, group)] = &subscriberWorker{
		group:      group,
		topic:      topic,
		h:          handler,
		executor:   uconcurrent.NewConcurrentControlExecutor(uint32(concurrency)),
		maxRetries: maxRetries,
		maxDelay:   maxDelay,
	}

	return nil
}

func (w *SubscriberWorker) Middleware(ctx context.Context, middleware SubscriberWrapper) {
	w.ms = append(w.ms, middleware)
}

// SubscriberChain returns a Middleware that specifies the chained handler for endpoint.
func (w *SubscriberWorker) subscriberChain(m ...SubscriberWrapper) SubscriberWrapper {
	return func(next SubscriberFunc) SubscriberFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// Run TODO: goroutine 飙升, event的内存逃逸
func (w *SubscriberWorker) Run(ctx context.Context) error {
	w.mutex.RLock()
	router := w.router
	w.mutex.RUnlock()
	eg, ctx := errgroup.WithContext(ctx)
	for _, worker := range router {
		eg.Go(func() error {
			consumer, err := w.mq.NewConsumer(ctx, worker.group, worker.topic)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				return err
			}
			defer consumer.Close(ctx)
			opts := consumer.Options(ctx)
			name := fmt.Sprintf("%s-%s", worker.topic, worker.group)
			for {
				select {
				case <-ctx.Done():
					logger.GetLogger(ctx).Info(fmt.Sprintf("group %s, topic %s consumer %d stop",
						worker.group, worker.topic))
					return nil
				default:
					et, err := consumer.Consume(ctx)
					if err != nil {
						logger.GetLogger(ctx).Error(err.Error())
					}
					ctx = umetadata.Merge(ctx, et.Header(ctx))
					// 包装中间件
					handler := w.subscriberChain(w.ms...)(worker.h)
					// 异步执行函数，并且不附带超时
					uwrapper.GoExecHandlerWithoutTimeout(ctx, func(ctx context.Context) error {
						err = uwrapper.RetryWithBackoff(ctx, func(ctx context.Context) error {
							// 并发处理 event
							_, err = worker.executor.Run(ctx, func(ctx context.Context) (interface{}, error) {
								err := handler(ctx, et)
								if err != nil {
									logger.GetLogger(ctx).Warn(err.Error())
									return nil, err
								}
								return nil, nil
							})
							if err != nil {
								logger.GetLogger(ctx).Error(err.Error())
								return err
							}
							logger.GetLogger(ctx).Info("handler success")
							return nil
						}, worker.maxRetries, worker.maxDelay)
						if err != nil {
							logger.GetLogger(ctx).Debug("kafka consumer process failed", "error", err)
							// 死信队列
							err = w.deadLetterQueuePublish(ctx, worker.topic, et)
							if err != nil {
								logger.GetLogger(ctx).Error(fmt.Sprintf("dead letter queue topic: %s group: %s event: (id %+v, message %+v, heders %+v) publish failed",
									worker.topic, worker.group, et.ID(ctx), et.Message(ctx), et.Header(ctx)), "error", err)
							}
							if opts != nil && opts.MonitorMetrics != nil {
								opts.MonitorMetrics.ConsumerMonitorRecorder(ctx, name, worker.topic, worker.group, "event_processor_fail")
							}
						}
						if opts != nil && opts.MonitorMetrics != nil {
							opts.MonitorMetrics.ConsumerMonitorRecorder(ctx, name, worker.topic, worker.group, "OK")
						}
						err = et.Ack(ctx)
						if err != nil {
							logger.GetLogger(ctx).Error(fmt.Sprintf("topic: %s group: %s event: (id %+v, message %+v, heders %+v) ack failed",
								worker.topic, worker.group, et.ID(ctx), et.Message(ctx), et.Header(ctx)), "error", err)
						}

						return nil
					})
				}
			}
		})

	}
	err := eg.Wait()
	if err != nil {
		panic(err)
		return err
	}
	logger.GetLogger(ctx).Info("all consumer stop")

	return nil
}

func (w *SubscriberWorker) deadLetterQueuePublish(ctx context.Context, topic string, event Event) error {
	logger.GetLogger(ctx).Info(fmt.Sprintf("dlq message %+v", event.String()))
	dlqTopic := fmt.Sprintf("%s_%s", topic, "dlq")
	producer, err := w.mq.NewProducer(ctx, dlqTopic)
	if err != nil {
		return err
	}
	err = producer.Publish(ctx, event)
	if err != nil {
		return err
	}
	err = event.Ack(ctx)
	if err != nil {
		return err
	}

	return nil
}
