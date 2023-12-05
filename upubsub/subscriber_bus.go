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

type subscriberEventBusWorker struct {
	group string
	topic string
	loop  map[int32]*subscriberWorker
}

type SubscriberEventBusWorker struct {
	mutex  sync.RWMutex
	mq     MQ
	router map[string]*subscriberEventBusWorker
	ms     []SubscriberWrapper
}

func NewSubscriberEventBusWorker(mq MQ) *SubscriberEventBusWorker {
	return &SubscriberEventBusWorker{
		mq:     mq,
		mutex:  sync.RWMutex{},
		router: make(map[string]*subscriberEventBusWorker),
		ms:     make([]SubscriberWrapper, 0),
	}
}

func (w *SubscriberEventBusWorker) RegisterSubscriber(ctx context.Context, messageType int32, topic, group string,
	concurrency, maxRetries int, maxDelay time.Duration, handler SubscriberFunc) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	name := fmt.Sprintf("%s-%s", topic, group)
	_, ok := w.router[name]
	if !ok {
		w.router[name] = &subscriberEventBusWorker{
			group: group,
			topic: topic,
			loop:  make(map[int32]*subscriberWorker),
		}
	}
	w.router[name].loop[messageType] = &subscriberWorker{
		group:      group,
		topic:      topic,
		h:          handler,
		executor:   uconcurrent.NewConcurrentControlExecutor(uint32(concurrency)),
		maxRetries: maxRetries,
		maxDelay:   maxDelay,
	}

	return nil
}

func (w *SubscriberEventBusWorker) Middleware(ctx context.Context, middleware SubscriberWrapper) {
	w.ms = append(w.ms, middleware)
}

// SubscriberChain returns a Middleware that specifies the chained handler for endpoint.
func (w *SubscriberEventBusWorker) subscriberChain(m ...SubscriberWrapper) SubscriberWrapper {
	return func(next SubscriberFunc) SubscriberFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// Run TODO: goroutine 飙升, event的内存逃逸
func (w *SubscriberEventBusWorker) Run(ctx context.Context) error {
	w.mutex.RLock()
	router := w.router
	w.mutex.RUnlock()
	for _, r := range router {
		logger.GetLogger(ctx).Info(fmt.Sprintf("subscriber event bus worker handler %+v", r.loop))
	}
	eg, ctx := errgroup.WithContext(ctx)
	for _, bus := range router {
		eg.Go(func() error {
			defer logger.GetLogger(ctx).Info(fmt.Sprintf("subscriber event bus worker quit"))
			consumer, err := w.mq.NewConsumer(ctx, bus.group, bus.topic)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				return err
			}
			defer consumer.Close(ctx)
			opts := consumer.Options(ctx)
			name := fmt.Sprintf("%s-%s", bus.topic, bus.group)
			topic := bus.topic
			group := bus.group
			for {
				select {
				case <-ctx.Done():
					logger.GetLogger(ctx).Info(fmt.Sprintf("group %s, topic %s consumer stop", bus.group, bus.topic))
					return nil
				default:
					logger.GetLogger(ctx).Debug("event bus subscriber worker start queue listening")
					et, err := consumer.Consume(ctx)
					if err != nil {
						logger.GetLogger(ctx).Error(err.Error())
					}
					ctx = umetadata.Merge(ctx, et.Header(ctx))
					logger.GetLogger(ctx).Debug(fmt.Sprintf("mq type %+v, event %+v, loop %+v",
						et.Type(ctx), et.String(), bus.loop))
					// 只读不写不会有资源竞争
					ws, ok := bus.loop[et.Type(ctx)]
					if !ok {
						logger.GetLogger(ctx).Debug(fmt.Sprintf("message type %+v not found register", et.Type(ctx)))
						continue
					}
					// 包装中间件
					handler := w.subscriberChain(w.ms...)(ws.h)
					// 异步执行函数，并且不附带超时
					uwrapper.GoExecHandlerWithoutTimeout(ctx, func(ctx context.Context) error {
						// handler 重试
						err = uwrapper.RetryWithBackoff(ctx, func(ctx context.Context) error {
							// 并发处理 event
							logger.GetLogger(ctx).Debug("concurrent control executor")
							_, err = ws.executor.Run(ctx, func(ctx context.Context) (interface{}, error) {
								logger.GetLogger(ctx).Debug("handler exec")
								err := handler(ctx, et)
								if err != nil {
									logger.GetLogger(ctx).Warn(err.Error())
									return nil, err
								}
								return nil, nil
							})
							if err != nil {
								logger.GetLogger(ctx).Warn(err.Error())
								return err
							}
							return nil
						}, ws.maxRetries, ws.maxDelay)
						if err != nil {
							logger.GetLogger(ctx).Debug("kafka consumer process failed", "error", err)
							// 死信队列
							err = w.deadLetterQueuePublish(ctx, topic, et)
							if err != nil {
								logger.GetLogger(ctx).Error(fmt.Sprintf("dead letter queue topic: %s group: %s event: (id %+v, message %+v, heders %+v) publish failed",
									topic, group, et.ID(ctx), et.Message(ctx), et.Header(ctx)), "error", err)
							}
							if opts != nil && opts.MonitorMetrics != nil {
								opts.MonitorMetrics.ConsumerMonitorRecorder(ctx, name, topic, group, "event_processor_fail")
							}
						}
						logger.GetLogger(ctx).Info("handler success")
						if opts != nil && opts.MonitorMetrics != nil {
							opts.MonitorMetrics.ConsumerMonitorRecorder(ctx, name, topic, group, "OK")
						}
						err = et.Ack(ctx)
						if err != nil {
							logger.GetLogger(ctx).Error(fmt.Sprintf("topic: %s group: %s event: (id %+v, message %+v, heders %+v) ack failed",
								topic, group, et.ID(ctx), et.Message(ctx), et.Header(ctx)), "error", err)
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

func (w *SubscriberEventBusWorker) deadLetterQueuePublish(ctx context.Context, topic string, event Event) error {
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
