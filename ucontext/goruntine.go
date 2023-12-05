package ucontext

import (
	"context"
	"time"
)

type UValueContext struct {
	context.Context
}

func (c *UValueContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *UValueContext) Done() <-chan struct{} {

	return nil
}

func (c *UValueContext) Err() error {
	return nil
}

func NewUValueContext(ctx context.Context) context.Context {
	return &UValueContext{ctx}
}

func NewUTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	valueCxt := &UValueContext{ctx}
	newCtx, cancel := context.WithTimeout(valueCxt, timeout)

	return newCtx, cancel
}
