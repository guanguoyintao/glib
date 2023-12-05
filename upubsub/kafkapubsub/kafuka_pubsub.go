package kafkapubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/AI/uglib/umetadata"
	"git.umu.work/AI/uglib/upubsub"
	goframeworkMQ "git.umu.work/be/goframework/async/mq"
	"git.umu.work/be/goframework/config"
	"git.umu.work/be/goframework/logger"
	"git.umu.work/be/goframework/metadata"
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

var once = sync.Once{}

type kafkaEvent struct {
	goframeworkMQ.Event
	Type   int32  `json:"type,omitempty"`
	ID     string `json:"ID"`
	Offset int64  `json:"offset"`
}

type KafkaEvent struct {
	Event        *kafkaEvent
	EventGroup   string                      `json:"group"`
	EventTopic   string                      `json:"topic"`
	kafkaMessage *sarama.ConsumerMessage     `json:"-"`
	session      sarama.ConsumerGroupSession `json:"-"`
	opts         *upubsub.Options            `json:"-"`
}

func kafkaHeaders(ctx context.Context) []sarama.RecordHeader {
	var headers []sarama.RecordHeader
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return headers
	}
	for k, v := range md {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}
	return headers
}

func NewKafukaEvent(ctx context.Context, id string, body []byte) upubsub.Event {
	header := umetadata.ToMap(ctx)
	return &KafkaEvent{
		Event: &kafkaEvent{
			Event: goframeworkMQ.Event{
				Headers:     header,
				Body:        string(body),
				CreatedTime: time.Now().UnixNano(),
			},
			ID: id,
		},
	}
}

func NewKafukaBusEvent(ctx context.Context, id string, eventType int32, body []byte) upubsub.Event {
	header := umetadata.ToMap(ctx)
	return &KafkaEvent{
		Event: &kafkaEvent{
			Event: goframeworkMQ.Event{
				Headers:     header,
				Body:        string(body),
				CreatedTime: time.Now().UnixNano(),
			},
			Type: eventType,
			ID:   id,
		},
	}
}

func (m *KafkaEvent) Header(ctx context.Context) map[string]string {
	if m.Event != nil && m.Event.Headers != nil {
		return m.Event.Headers
	}

	return nil
}

func (m *KafkaEvent) String() string {
	s, err := ujson.Marshal(m)
	if err != nil {
		return ""
	}

	return string(s)
}

func (m *KafkaEvent) ID(ctx context.Context) string {
	if m.Event != nil {
		return m.Event.ID
	}

	return ""
}

func (m *KafkaEvent) Message(ctx context.Context) []byte {
	if m.Event != nil {
		return []byte(m.Event.Body)
	}

	return nil
}

func (m *KafkaEvent) Topic(ctx context.Context) string {
	if m.Event != nil {
		return m.Event.Topic
	}

	return ""
}

func (m *KafkaEvent) Type(ctx context.Context) int32 {
	if m.Event != nil {
		return m.Event.Type
	}

	return 0
}

func (m *KafkaEvent) Ack(ctx context.Context) error {
	var err error
	var header []byte
	if m.Event != nil && m.Event.Headers != nil {
		// header, err = ujson.Marshal(m.event.Headers)
		if err != nil {
			if m.opts != nil && m.opts.MonitorMetrics != nil {
				name := fmt.Sprintf("%s-%s", m.EventTopic, m.EventGroup)
				m.opts.MonitorMetrics.ConsumerMonitorRecorder(ctx, name, m.EventTopic, m.EventGroup, "event_ack_fail")
			}
			return err
		}
	}
	m.session.MarkMessage(m.kafkaMessage, string(header))
	m.session.Commit()
	logger.GetLogger(ctx).Info(fmt.Sprintf("kfk %+v event offset %d ack", m.session.MemberID(), m.kafkaMessage.Offset))

	return nil
}

func (m *KafkaEvent) Nack(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

type kfkClaim struct {
	topic     string
	group     string
	ctx       context.Context
	opts      *upubsub.Options
	eventChan chan upubsub.Event
}

func (c *kfkClaim) claimHandler(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.GetLogger(c.ctx).Info("claim handler start")
	for message := range claim.Messages() {
		if c.opts.MonitorMetrics != nil {
			c.opts.MonitorMetrics.ConsumerMonitorStart(c.ctx)
		}
		logger.GetLogger(c.ctx).Info(fmt.Sprintf("kafuka consumer offset %d", message.Offset))
		logger.GetLogger(c.ctx).Info("kafuka consumer begin", "message", message)
		// todo: header
		var event kafkaEvent
		err := ujson.Unmarshal(message.Value, &event)
		if err != nil {
			panic(err)
		}
		event.Topic = c.topic
		event.Offset = message.Offset
		c.eventChan <- &KafkaEvent{
			EventTopic:   c.topic,
			EventGroup:   c.group,
			opts:         c.opts,
			session:      session,
			kafkaMessage: message,
			Event:        &event,
		}
	}

	return nil
}

type KafkaConsumer struct {
	group string
	topic string
	opts  *upubsub.Options
	goframeworkMQ.Consumer
	claim *kfkClaim
}

func (c *KafkaConsumer) Consume(ctx context.Context) (upubsub.Event, error) {
	event := <-c.claim.eventChan

	return event, nil
}

func (c *KafkaConsumer) Options(ctx context.Context) *upubsub.Options {
	return c.opts
}

func (c *KafkaConsumer) Close(ctx context.Context) error {
	err := c.Consumer.Stop()
	if err != nil {
		return err
	}
	logger.GetLogger(ctx).Info(fmt.Sprintf("kfk topic %s group %s consumer closed", c.topic, c.group))

	return nil
}

type KafkaProducer struct {
	opts     *upubsub.Options
	producer map[string]goframeworkMQ.Producer
}

func (p *KafkaProducer) Publish(ctx context.Context, event upubsub.Event) error {
	e := event
	if p.opts.MonitorMetrics != nil {
		p.opts.MonitorMetrics.ProducerMonitorStart(ctx)
	}
	for topic, producer := range p.producer {
		logger.GetLogger(ctx).Info(fmt.Sprintf("kfk topic %s producer publish event %+v", topic, event))
		var key string
		if len(p.opts.PartitionKey) > 0 {
			key = p.opts.PartitionKey
		}
		publishEvent, ok := e.(*KafkaEvent)
		if !ok {
			publishEvent = &KafkaEvent{
				Event: &kafkaEvent{
					Event: goframeworkMQ.Event{
						Topic:   topic,
						Headers: event.Header(ctx),
						Body:    string(event.Message(ctx)),
					},
					ID: event.ID(ctx),
				},
			}
		}
		publishEvent.Event.CreatedTime = time.Now().UnixNano()
		eb, err := ujson.Marshal(publishEvent.Event)
		if err != nil {
			if p.opts.MonitorMetrics != nil {
				p.opts.MonitorMetrics.ProducerMonitorRecorder(ctx, topic, topic, "event_producer_msg_err")
			}
			return err
		}
		err = producer.PubCustom(ctx, topic, key, string(eb))
		if err != nil {
			if p.opts.MonitorMetrics != nil {
				p.opts.MonitorMetrics.ProducerMonitorRecorder(ctx, topic, topic, "event_send_message_fail")
			}
			return err
		}
		if p.opts.MonitorMetrics != nil {
			p.opts.MonitorMetrics.ProducerMonitorRecorder(ctx, topic, topic, "OK")
		}
	}

	return nil
}

func (p *KafkaProducer) Options(ctx context.Context) *upubsub.Options {
	return p.opts
}

func (p *KafkaProducer) Close(ctx context.Context) error {
	for topic, producer := range p.producer {
		err := producer.Close()
		if err != nil {
			return err
		}
		logger.GetLogger(ctx).Info(fmt.Sprintf("kfk topic %s producer closed", topic))
	}
	return nil
}

type KafkaMQ struct {
	cfg      config.Config
	opts     *upubsub.Options
	claimCtx context.Context
}

func NewKafkaMQ(cfg config.Config, opts ...upubsub.Option) upubsub.MQ {
	once.Do(func() {
		goframeworkMQ.InitConsumer(cfg)
		goframeworkMQ.InitProducer(cfg)
	})
	opt := &upubsub.Options{}
	for _, o := range opts {
		o(opt)
	}
	return &KafkaMQ{
		cfg:  cfg,
		opts: opt,
	}
}

func (mq *KafkaMQ) NewConsumer(ctx context.Context, group string, topic string) (upubsub.Consumer, error) {
	name := fmt.Sprintf("%s-%s", topic, group)
	consumer, err := goframeworkMQ.GetConsumer(name)
	if err != nil {
		return nil, err
	}
	c := &kfkClaim{
		topic:     topic,
		group:     group,
		ctx:       ctx,
		opts:      mq.opts,
		eventChan: make(chan upubsub.Event),
	}
	consumer.SetCustomConsumer(c.claimHandler)
	err = consumer.Start(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error("consumer err is %+v", err)
	}

	return &KafkaConsumer{
		group:    group,
		topic:    topic,
		opts:     mq.opts,
		Consumer: consumer,
		claim:    c,
	}, nil
}

func (mq *KafkaMQ) NewProducer(ctx context.Context, topics ...string) (upubsub.Producer, error) {
	producer := make(map[string]goframeworkMQ.Producer, len(topics))
	for _, topic := range topics {
		p, err := goframeworkMQ.GetProducer(topic)
		if err != nil {
			return nil, err
		}
		producer[topic] = p
	}

	return &KafkaProducer{
		opts:     mq.opts,
		producer: producer,
	}, nil
}
