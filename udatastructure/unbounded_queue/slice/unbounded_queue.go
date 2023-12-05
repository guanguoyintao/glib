package slicequeue

import (
	"container/list"
	"context"
	"sync"
)

type Queue struct {
	queue *list.List
	mu    sync.Mutex
	wg    sync.WaitGroup
}

func NewUnboundedQueue(capacity uint32) *Queue {
	return &Queue{
		queue: list.New(),
	}
}

func (q *Queue) Enqueue(ctx context.Context, item interface{}) (bool, uint32) {
	q.mu.Lock()
	q.queue.PushBack(item)
	quantity := q.queue.Len()
	q.mu.Unlock()

	return true, uint32(quantity)
}

func (q *Queue) Dequeue(ctx context.Context) (interface{}, bool, uint32) {
	for {
		q.mu.Lock()
		if q.queue.Len() == 0 {
			q.mu.Unlock()
			q.wg.Wait() // 队列为空，等待生产者入队
		} else {
			item := q.queue.Front()
			quantity := q.queue.Len()
			q.queue.Remove(item)
			q.mu.Unlock()

			return item.Value, true, uint32(quantity)
		}
	}

}
