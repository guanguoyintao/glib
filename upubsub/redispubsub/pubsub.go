package redispubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/AI/uglib/upubsub"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/logger"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

const (
	mqDBName                 = "ai-mq"
	maxRetries int           = 10
	maxDelay   time.Duration = time.Hour
)

type redisPubSubMQConsumerGroup struct {
	topic       string
	messageChan *chan upubsub.Event
}

type redisPubSubMQ struct {
	once        sync.Once
	groupMap    map[string]*chan redisPubSubMQConsumerGroup
	consumerMap map[string]*chan upubsub.Event
	redisClient *redis.Client
}

func newRedisPubSubMQ(ctx context.Context) (*redisPubSubMQ, error) {
	redisClient, err := cache.GetRedis(ctx, mqDBName)
	if err != nil {
		return nil, err
	}

	return &redisPubSubMQ{
		once:        sync.Once{},
		groupMap:    make(map[string]*chan redisPubSubMQConsumerGroup),
		consumerMap: make(map[string]*chan upubsub.Event),
		redisClient: redisClient.(*redis.Client),
	}, nil
}

func (mq *redisPubSubMQ) consumerMaster(ctx context.Context, group string, c *chan redisPubSubMQConsumerGroup) error {
	topics := make([]string, 0)
	consumerMap := make(map[string]*chan upubsub.Event)
	pubsub := mq.redisClient.Subscribe(ctx, topics...)
	pubsub.Channel()
	for {
		select {
		case consumerGroup := <-*c:
			_, ok := consumerMap[consumerGroup.topic]
			if !ok {
				topics = append(topics, consumerGroup.topic)
				consumerMap[consumerGroup.topic] = consumerGroup.messageChan
				pubsub = mq.redisClient.Subscribe(ctx, topics...)
			}
		case msg := <-pubsub.Channel():
			messageChan := consumerMap[msg.Channel]
			message := &redisPubSubEvent{
				EventHeader: map[string]string{
					"topic": msg.Channel,
					"group": group,
				},
				EventTopic: msg.Channel,
				Body:       []byte(msg.Payload),
			}
			go func() {
				*messageChan <- message
			}()
		}
	}
}

func (mq *redisPubSubMQ) NewConsumer(ctx context.Context, group string, topic string) (upubsub.Consumer, error) {
	if len(group) == 0 {
		group = "default"
	}
	messageChan, ok := mq.consumerMap[fmt.Sprintf("%v-%v", group, topic)]
	if !ok {
		messageChanTmp := make(chan upubsub.Event)
		messageChan = &messageChanTmp
		mq.consumerMap[fmt.Sprintf("%v-%v", group, topic)] = messageChan
		consumerGroupChan, groupIsExit := mq.groupMap[group]
		if !groupIsExit {
			consumerGroupChanTmp := make(chan redisPubSubMQConsumerGroup)
			consumerGroupChan = &consumerGroupChanTmp
			mq.groupMap[group] = consumerGroupChan
			go func() {
				err := mq.consumerMaster(ctx, group, consumerGroupChan)
				if err != nil {
					logger.GetLogger(ctx).Warn(err.Error())
				}
			}()
		}
		*consumerGroupChan <- redisPubSubMQConsumerGroup{
			topic:       topic,
			messageChan: messageChan,
		}
	}

	return &redisPubSubConsumer{
		messageChan: messageChan,
		redisClient: mq.redisClient,
		group:       group,
		topic:       topic,
	}, nil
}

func (mq *redisPubSubMQ) NewProducer(ctx context.Context, topics ...string) (upubsub.Producer, error) {
	producer := &redisPubSubProducer{
		redisClient: mq.redisClient,
		topics:      topics,
	}

	return producer, nil
}

type redisPubSubConsumer struct {
	messageChan *chan upubsub.Event
	redisClient *redis.Client
	group       string
	topic       string
}

func (c *redisPubSubConsumer) Options(ctx context.Context) *upubsub.Options {
	// TODO implement me
	panic("implement me")
}

func (c *redisPubSubConsumer) Consume(ctx context.Context) (upubsub.Event, error) {
	msg, _ := <-*c.messageChan

	return msg, nil
}

func (c *redisPubSubConsumer) Close(ctx context.Context) error {
	// result, err := c.redisClient.PubSubNumSub(ctx, c.topics...).Result()
	// if err != nil {
	//	return err
	// }

	return nil
}

type redisPubSubProducer struct {
	redisClient *redis.Client
	topics      []string
}

func (p *redisPubSubProducer) Options(ctx context.Context) *upubsub.Options {
	// TODO implement me
	panic("implement me")
}

func (p *redisPubSubProducer) Publish(ctx context.Context, msg upubsub.Event) error {
	// Publish message to each topic
	for _, topic := range p.topics {
		err := uwrapper.RetryWithBackoff(ctx, func(ctx context.Context) error {
			status, err := p.redisClient.Publish(ctx, topic, msg.Message(ctx)).Result()
			if err != nil {
				return err
			}
			if status != 1 {
				logger.GetLogger(ctx).Warn(fmt.Sprintf("redis pub/sub pulish status is %+v\n", status))
				return uerrors.UErrorSystemError
			}

			return nil
		}, maxRetries, maxDelay)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *redisPubSubProducer) Close(ctx context.Context) error {
	// nothing to do
	return nil
}

type redisPubSubEvent struct {
	EventHeader map[string]string `json:"header"`
	EventID     string            `json:"id"`
	EventTopic  string            `json:"topic"`
	Body        []byte            `json:"body"`
}

func (m *redisPubSubEvent) String() string {
	s, err := ujson.Marshal(m)
	if err != nil {
		return ""
	}

	return string(s)
}

func (m *redisPubSubEvent) Type(ctx context.Context) int32 {
	// TODO implement me
	panic("implement me")
}

func newRedisPubSubEvent(id string, body []byte) upubsub.Event {
	return &redisPubSubEvent{
		Body:    body,
		EventID: id,
	}
}

func (m *redisPubSubEvent) Topic(ctx context.Context) string {
	return m.EventTopic
}

func (m *redisPubSubEvent) Message(ctx context.Context) []byte {
	return m.Body
}

func (m *redisPubSubEvent) Header(ctx context.Context) map[string]string {
	return m.EventHeader
}

func (m *redisPubSubEvent) ID(ctx context.Context) string {
	return m.EventID
}

func (m *redisPubSubEvent) Ack(ctx context.Context) error {
	return nil
}
func (m *redisPubSubEvent) Nack(ctx context.Context) error {
	return nil
}
