package redispubsub

import (
	"context"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/AI/uglib/umetadata"
	"git.umu.work/AI/uglib/upubsub"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/logger"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"io"
	"sync"
	"time"
)

type RedisListEvent struct {
	Group       string            `json:"group"`
	EventTopic  string            `json:"topic"`
	EventID     string            `json:"id"`
	EventHeader map[string]string `json:"header"`
	CreatedTime time.Time         `json:"created_time"`
	body        []byte            `json:"-"`
}

func (r *RedisListEvent) Topic(ctx context.Context) string {
	return r.EventTopic
}

func (r *RedisListEvent) Header(ctx context.Context) map[string]string {
	return r.EventHeader
}

func (r *RedisListEvent) ID(ctx context.Context) string {
	return r.EventID
}

func (r *RedisListEvent) Message(ctx context.Context) []byte {
	return r.body
}

func (r *RedisListEvent) Ack(ctx context.Context) error {
	return nil
}

func (r *RedisListEvent) Nack(ctx context.Context) error {
	return nil
}

func (r *RedisListEvent) Type(ctx context.Context) int32 {
	// TODO implement me
	panic("implement me")
}

func (r *RedisListEvent) String() string {
	es, err := ujson.Marshal(r)
	if err != nil {
		return ""
	}
	var m map[string]interface{}
	err = ujson.Unmarshal(r.body, &m)
	if err != nil {
		return ""
	}
	err = ujson.Unmarshal(es, &m)
	if err != nil {
		return ""
	}
	res, err := ujson.Marshal(m)
	if err != nil {
		return ""
	}

	return string(res)
}

func NewRedisListEvent(ctx context.Context, id string, body []byte) upubsub.Event {
	header := umetadata.ToMap(ctx)
	return &RedisListEvent{
		EventID:     id,
		EventHeader: header,
		CreatedTime: time.Now(),
		body:        body,
	}
}

type redisListMQ struct {
	mutex        *sync.Mutex
	redisClient  redis.Cmdable
	onceConsumer map[string]<-chan upubsub.Event
	opts         *upubsub.Options
}

func NewRedisListMQ(ctx context.Context, opts ...upubsub.Option) (upubsub.MQ, error) {
	redisClient, err := cache.GetRedis(ctx, mqDBName)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}
	opt := &upubsub.Options{}
	for _, o := range opts {
		o(opt)
	}

	return &redisListMQ{
		redisClient:  redisClient,
		mutex:        &sync.Mutex{},
		onceConsumer: make(map[string]<-chan upubsub.Event, 0),
		opts:         opt,
	}, nil
}

type redisListConsumer struct {
	topic       string
	redisClient redis.Cmdable
	opts        *upubsub.Options
	eventChan   <-chan upubsub.Event
}

func (r *redisListConsumer) Options(ctx context.Context) *upubsub.Options {
	return r.opts
}

func (r *redisListConsumer) Consume(ctx context.Context) (upubsub.Event, error) {
	select {
	case event, ok := <-r.eventChan:
		if !ok {
			return nil, io.EOF
		}
		return event, nil
	case <-ctx.Done():
		return nil, io.EOF
	}
}

func (r *redisListConsumer) Close(ctx context.Context) error {
	logger.GetLogger(ctx).Info("redis list mq is closed")
	return nil
}

type redisListProducer struct {
	topics      []string
	redisClient redis.Cmdable
	opts        *upubsub.Options
}

func (r *redisListProducer) Options(ctx context.Context) *upubsub.Options {
	return r.opts
}

func (r *redisListProducer) Publish(ctx context.Context, event upubsub.Event) error {
	for _, topic := range r.topics {
		err := r.redisClient.RPush(ctx, topic, event.String()).Err()
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return err
		}
	}

	return nil
}

func (r *redisListProducer) Close(ctx context.Context) error {
	return nil
}

func (r *redisListMQ) newConsumer(ctx context.Context, group, topic string, eventChan chan upubsub.Event) error {
	var timeout time.Duration
	if r.opts.Timeout > 0 {
		timeout = r.opts.Timeout
	} else {
		timeout = 1 * time.Minute
	}
	msgs, err := r.redisClient.BLPop(ctx, timeout, topic).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return errors.Wrapf(err, "lpop failed topic=%s", topic)
		}
	}
	for _, msg := range msgs {
		var event RedisListEvent
		err = ujson.Unmarshal([]byte(msg), &event)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return err
		}
		// 发送event
		event.EventTopic = topic
		event.Group = group
		eventChan <- &event
	}

	return nil
}

func (r *redisListMQ) NewConsumer(ctx context.Context, group, topic string) (upubsub.Consumer, error) {
	consumer := &redisListConsumer{
		topic:       topic,
		redisClient: r.redisClient,
		opts:        r.opts,
	}
	r.mutex.Lock()
	_, ok := r.onceConsumer[topic]
	if !ok {
		eventChan := make(chan upubsub.Event, 0)
		r.onceConsumer[topic] = eventChan
		go func() {
			for {
				err := r.newConsumer(ctx, group, topic, eventChan)
				if err != nil {
					logger.GetLogger(ctx).Warn(err.Error())
					break
				}
			}
		}()

		consumer.eventChan = eventChan
	} else {
		eventChan := r.onceConsumer[topic]
		consumer.eventChan = eventChan
	}
	r.mutex.Unlock()

	return consumer, nil
}

func (r *redisListMQ) NewProducer(ctx context.Context, topics ...string) (upubsub.Producer, error) {
	return &redisListProducer{
		topics:      topics,
		redisClient: r.redisClient,
		opts:        r.opts,
	}, nil
}
