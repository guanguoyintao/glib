package ucurry

import (
	"context"
)

type CurryFunctor struct {
	handler func(ctx context.Context, args ...interface{}) interface{}
	result  interface{}
}

func Curry(f func(ctx context.Context, args ...interface{}) interface{}) *CurryFunctor {
	return &CurryFunctor{
		handler: f,
	}
}

func (c *CurryFunctor) Call(ctx context.Context, args ...interface{}) interface{} {
	if c.result != nil {
		args = append([]interface{}{c.result}, args)
	}
	res := c.handler(ctx, args)
	c.result = res

	return res
}
