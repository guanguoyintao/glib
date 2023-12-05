package upubsub

import (
	"context"
)

type Event interface {
	Topic(ctx context.Context) string
	Header(ctx context.Context) map[string]string
	ID(ctx context.Context) string
	Message(ctx context.Context) []byte
	Ack(ctx context.Context) error
	Nack(ctx context.Context) error
	Type(ctx context.Context) int32
	String() string
}

type Consumer interface {
	Options(ctx context.Context) *Options
	Consume(ctx context.Context) (Event, error)
	Close(ctx context.Context) error
}

type Producer interface {
	Options(ctx context.Context) *Options
	Publish(ctx context.Context, event Event) error
	Close(ctx context.Context) error
}

type MQ interface {
	NewConsumer(ctx context.Context, group, topic string) (Consumer, error)
	NewProducer(ctx context.Context, topics ...string) (Producer, error)
}

type SubscriberFunc func(ctx context.Context, event Event) error

type MessageFunc func(ctx context.Context, id string, body []byte) Event

// SubscriberWrapper wraps the SubscriberFunc and returns the equivalent
type SubscriberWrapper func(SubscriberFunc) SubscriberFunc
