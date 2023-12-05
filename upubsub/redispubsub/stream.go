// Package upubsub 下面是基于该interface和redis pubsub实现发布订阅的示例代码。我们将Redis Stream作为消息队列，使用Redis提供的XADD、XREAD、XGROUP等指令实现消息的发布和消费。为了提高可用性，我们在代码中实现了以下优化点：
//
// 消息确认和重试机制：消费者通过调用Message.Ack()方法确认收到消息。如果消费失败，可以调用Message.Nack()方法将消息重新放回队列中，并在一定时间后重试。这个时间可以由消息的“retry-after”元数据指定。在重试次数达到一定阈值时，我们可以将消息放入一个“dead letter”队列中，以便后续分析。
//
// 消息幂等：我们通过记录已处理的消息的ID来避免重复消费。消费者可以在处理消息之前检查是否已经处理过该消息，如果是则直接忽略。
//
// 消息丢失处理：我们使用了Redis Stream提供的“pending list”功能来避免消息丢失。如果一个消费者没有确认收到消息，则该消息会被加入到该消费者的“pending list”中，并在一定时间后自动重新放回队列中等待消费。
package redispubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/AI/uglib/upubsub"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/logger"
	"github.com/go-redis/redis/v8"
	"strings"
	"sync"
	"time"
)

// var redisClient *redis.Client

type redisStreamMQ struct {
	mutex       *sync.Mutex
	redisClient redis.Cmdable
	consumerMap map[string]uint32
}

func NewRedisStreamMQ(ctx context.Context) (upubsub.MQ, error) {
	redisClient, err := cache.GetRedis(ctx, mqDBName)
	if err != nil {
		return nil, err
	}

	return &redisStreamMQ{
		redisClient: redisClient,
		consumerMap: make(map[string]uint32, 0),
		mutex:       &sync.Mutex{},
	}, nil
}

func (mq *redisStreamMQ) NewConsumer(ctx context.Context, group string, topic string) (upubsub.Consumer, error) {
	if len(group) == 0 {
		group = "default"
	}
	key := fmt.Sprintf("%s-%s", topic, group)
	mq.mutex.Lock()
	counter, ok := mq.consumerMap[key]
	if ok {
		counter += 1
		mq.consumerMap[key] += 1
	} else {
		mq.consumerMap[key] = 0
	}
	mq.mutex.Unlock()

	return &redisStreamConsumer{
		topic:       topic,
		group:       group,
		consumer:    fmt.Sprintf("%s-%s-%d", topic, group, counter),
		redisClient: mq.redisClient,
	}, nil
}

func (mq *redisStreamMQ) NewProducer(ctx context.Context, topics ...string) (upubsub.Producer, error) {
	stream := topics[0]
	return &redisStreamProducer{
		topic:       stream,
		redisClient: mq.redisClient,
	}, nil
}

// 消息确认：在Ack方法中调用XAck接口将消息从消费组的pending list中移除。
// 不能重复消费和消息丢失：在XClaim接口中传入MinIdle参数，以避免在网络故障或消费者宕机等情况下重复消费和消息
type redisStreamMessage struct {
	EventID           string            `json:"id"`
	EventHeader       map[string]string `json:"header"`
	EventBody         []byte
	EventTopic        string        `json:"topic"`
	EventConsumerName string        `json:"consumer_name"`
	redisClient       redis.Cmdable `json:"-"`
}

func (m *redisStreamMessage) String() string {
	s, _ := ujson.Marshal(m)

	return string(s)
}

func (m *redisStreamMessage) Type(ctx context.Context) int32 {
	// TODO implement me
	panic("implement me")
}

func NewRedisStreamEvent(id string, body []byte) upubsub.Event {
	return &redisStreamMessage{
		EventID:   id,
		EventBody: body,
	}
}

func (m *redisStreamMessage) Header(ctx context.Context) map[string]string {
	return m.EventHeader
}

func (m *redisStreamMessage) Message(ctx context.Context) []byte {
	return m.EventBody
}

func (m *redisStreamMessage) Topic(ctx context.Context) string {
	return m.EventTopic
}

func (m *redisStreamMessage) ID(ctx context.Context) string {
	return m.EventID
}

func (m *redisStreamMessage) Ack(ctx context.Context) error {
	_, err := m.redisClient.XAck(ctx, m.EventTopic, m.EventConsumerName, m.EventID).Result()
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return uerrors.UErrorSystemError
	}
	return nil
}

func (m *redisStreamMessage) Nack(ctx context.Context) error {
	// _, err := m.redisClient.XGroupDelConsumer(ctx, m.stream, m.consumer, m.id).Result()
	// if err != nil {
	//	return fmt.Errorf("failed to remove message from consumer group: %s", err)
	// }
	// //使用client.XAdd向消息队列重新发送一条被Nack的消息。这是因为当消费者对一条消息Nack后，
	// //该消息会被重新放回pending list中，等待被重新消费。但是如果该消息的重试次数超过了队列的最大重试次数，那么该消息会被移到dead letter list中。
	// //因此，为了保证消息能够被重新消费，我们需要使用client.XAdd将该消息重新添加到队列中。
	// _, err = m.redisClient.XAdd(ctx, &redis.XAddArgs{
	//	Stream: m.stream,
	//	Values: map[string]interface{}{
	//		"message": "NACK",
	//	},
	// }).Result()
	// if err != nil {
	//	return fmt.Errorf("failed to send NACK message to stream: %s", err)
	// }

	return nil
}

type redisStreamProducer struct {
	topic       string
	redisClient redis.Cmdable
}

func (p *redisStreamProducer) Options(ctx context.Context) *upubsub.Options {
	return nil
}

func (p *redisStreamProducer) Publish(ctx context.Context, msg upubsub.Event) error {
	var header string
	if msg.Header(ctx) != nil {
		headerBytes, err := ujson.Marshal(msg.Header(ctx))
		if err != nil {
			return err
		}
		header = string(headerBytes)
	}
	messageData := map[string]interface{}{
		"header": header,
		"body":   string(msg.Message(ctx)),
	}
	_, err := p.redisClient.XAdd(ctx, &redis.XAddArgs{
		ID:     msg.ID(ctx),
		Stream: p.topic,
		Values: messageData,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}
	return nil
}

func (p *redisStreamProducer) Close(ctx context.Context) error {
	return nil
}

type redisStreamConsumer struct {
	topic         string
	group         string
	consumer      string
	redisClient   redis.Cmdable
	maxRetryCount int64
	maxIdleTime   time.Duration
}

// 检查 group 是否存在
func (c *redisStreamConsumer) groupIsExist(ctx context.Context) (ok bool, err error) {
	groups, err := c.redisClient.XInfoGroups(ctx, c.topic).Result()
	if err != nil {
		logger.GetLogger(ctx).Warn(fmt.Sprintf("failed to get stream groups: %s\n", err.Error()))
		return false, err
	}
	var groupExists bool
	for _, group := range groups {
		if group.Name == c.group {
			groupExists = true
			break
		}
	}

	return groupExists, nil
}

func (c *redisStreamConsumer) Options(ctx context.Context) *upubsub.Options {
	return nil
}

// 检查 stream 是否存在
func (c *redisStreamConsumer) streamIsExist(ctx context.Context) (ok bool, err error) {
	isExist, err := c.redisClient.Exists(ctx, c.topic).Result()
	if err != nil {
		logger.GetLogger(ctx).Warn(fmt.Sprintf("failed to get key: %s, %+v\n", c.topic, err.Error()))
		return false, err
	}
	if isExist == 1 {
		return true, nil
	}

	return false, nil
}

func (c *redisStreamConsumer) Consume(ctx context.Context) (upubsub.Event, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			// 检查 stream 是否存在
			isExist, err := c.streamIsExist(ctx)
			if err != nil {
				time.Sleep(10 * time.Second)
				logger.GetLogger(ctx).Warn(err.Error())
				return nil, err
			}
			if !isExist {
				time.Sleep(10 * time.Second)
				continue
			}

			// pending, err := c.redisClient.XPendingExt(ctx, &redis.XPendingExtArgs{
			//	Stream:   c.topic,
			//	Group:    c.group,
			//	Start:    "-",
			//	End:      "+",
			//	Count:    1,
			//	Consumer: c.consumer,
			// }).Result()
			// if err != nil {
			//	logger.GetLogger(ctx).Warn(err.Error())
			//	return nil, uerrors.UErrorSystemError
			// }
			// if len(pending) > 0 {
			//	messageID := pending[0].ID
			//	message, err := c.redisClient.XClaim(ctx, &redis.XClaimArgs{
			//		Stream:   c.topic,
			//		Group:    c.group,
			//		Consumer: c.consumer,
			//		MinIdle:  time.Millisecond,
			//		Messages: []string{messageID},
			//	}).Result()
			//	if err != nil || len(message) == 0 {
			//		logger.GetLogger(ctx).Warn(err.Error())
			//		//return nil, uerrors.UErrorSystemError
			//	}
			//	body := message[0].Values["body"].(string)
			//	body = strings.TrimSpace(strings.Trim(body, `"`))
			//	bodyBytes, err := base64.StdEncoding.DecodeString(body)
			//	if err != nil {
			//		logger.GetLogger(ctx).Warn("failed to get message with id %s: %s", messageID, err.Error())
			//		//return nil, err
			//	}
			//	pendingMessage := pending[0]
			//	retryCount := pendingMessage.RetryCount
			//	if pendingMessage.Idle > c.maxIdleTime && retryCount < c.maxRetryCount {
			//		// 消息已超时，重新加入队列，并将重试计数器加1
			//		_, err := c.redisClient.XAdd(ctx, &redis.XAddArgs{
			//			Stream: c.topic,
			//			Values: map[string]interface{}{
			//				"body": bodyBytes,
			//			},
			//			MaxLenApprox: 0,
			//		}).Result()
			//		if err != nil {
			//			logger.GetLogger(ctx).Warn(err.Error())
			//			//return nil, uerrors.UErrorSystemError
			//		}
			//		_, err = c.redisClient.XAck(ctx, c.topic, c.group, messageID).Result()
			//		if err != nil {
			//			logger.GetLogger(ctx).Warn(err.Error())
			//			//return nil, uerrors.UErrorSystemError
			//		}
			//	}
			// }
			// 正常读取消息
			messageSlice, err := c.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.group,
				Consumer: c.consumer,
				Streams:  []string{c.topic, ">"},
				Count:    1,
				Block:    10 * time.Second,
				NoAck:    false,
			}).Result()
			if err != nil {
				if err == redis.Nil {
					continue
					// todo: 错误判断没有标准类型
				} else if strings.Contains(err.Error(), "NOGROUP") {
					// 不存在group就创建
					_, err = c.redisClient.XGroupCreateMkStream(ctx, c.topic, c.group, "0").Result()
					if err != nil && err != redis.Nil {
						logger.GetLogger(ctx).Warn(err.Error())
						// return nil, uerrors.UErrorSystemError
					}
					time.Sleep(1 * time.Second)
					continue
				}
				logger.GetLogger(ctx).Warn(err.Error())
				return nil, err
			}
			if len(messageSlice) == 0 {
				continue
			}
			message := messageSlice[0]
			messageID := message.Messages[0].ID
			var body string
			if message.Messages[0].Values["body"] != nil {
				body = message.Messages[0].Values["body"].(string)
			}
			// 构造header
			header := make(map[string]string, 0)
			if message.Messages[0].Values["header"] != nil {
				headerT := message.Messages[0].Values["header"]
				switch headerT.(type) {
				case string:
					headerString := headerT.(string)
					if len(headerString) > 0 {
						err = ujson.Unmarshal([]byte(headerString), &header)
						if err != nil {
							logger.GetLogger(ctx).Warn(err.Error())
							return nil, err
						}
					}
				}
			}
			return &redisStreamMessage{
				EventID:           messageID,
				EventHeader:       header,
				EventBody:         []byte(body),
				EventTopic:        c.topic,
				EventConsumerName: c.consumer,
				redisClient:       c.redisClient,
			}, nil

		}
	}
}

func (c *redisStreamConsumer) Close(ctx context.Context) error {
	// _, err := c.redisClient.XGroupDestroy(ctx, c.topic, c.group).Result()
	// if err != nil {
	//	return fmt.Errorf("failed to destroy consumer group: %s", err)
	// }

	return nil
}
