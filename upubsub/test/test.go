package test

import (
	"context"
	"git.umu.work/AI/uglib/ujson"
)

type Event struct {
	EventTopic   string            `json:"topic"`
	EventHeader  map[string]string `json:"header"`
	EventMessage []byte            `json:"message"`
	EventType    int32             `json:"type"`
	EventID      string            `json:"id"`
}

func (e *Event) String() string {
	s, err := ujson.Marshal(e)
	if err != nil {
		return ""
	}

	return string(s)
}

func (e *Event) Topic(ctx context.Context) string {
	return e.EventTopic
}

func (e *Event) Header(ctx context.Context) map[string]string {
	return e.EventHeader
}

func (e *Event) ID(ctx context.Context) string {
	return e.EventID
}

func (e *Event) Message(ctx context.Context) []byte {
	return e.EventMessage
}

func (e *Event) Ack(ctx context.Context) error {
	return nil
}

func (e *Event) Nack(ctx context.Context) error {
	return nil
}

func (e *Event) Type(ctx context.Context) int32 {
	return e.EventType
}
