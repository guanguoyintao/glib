package unboundedqueue

import (
	"context"
	disruptorqueue "git.umu.work/AI/uglib/udatastructure/unbounded_queue/disruptor"
)

type UnboundedQueue interface {
	Enqueue(ctx context.Context, item interface{}) (ok bool, quantity uint32)
	Dequeue(ctx context.Context) (item interface{}, ok bool, quantity uint32)
}

func NewUnboundedQueue(capacity uint32) UnboundedQueue {
	return disruptorqueue.NewQueue(capacity)
}
