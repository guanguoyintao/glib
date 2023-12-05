package utakenumberqueue

import (
	"context"
	"fmt"
	"git.umu.work/be/goframework/logger"
	"time"
)

type EndNumberValue struct{}

type Result struct {
	Value  interface{}
	Err    error
	Number int32
}

// TakeNumberQueue 取号队列
type TakeNumberQueue struct {
	offset  int32
	queue   map[int32]interface{}
	out     chan *Result
	exit    chan struct{}
	timer   *time.Timer
	timeout time.Duration
}

func NewTakeNumberQueue(ctx context.Context, size int, timeout time.Duration) *TakeNumberQueue {
	timer := time.NewTimer(timeout)
	tq := &TakeNumberQueue{
		offset:  0,
		queue:   make(map[int32]interface{}, 100),
		out:     make(chan *Result, size),
		exit:    make(chan struct{}),
		timer:   timer,
		timeout: timeout,
	}
	go func() {
		for {
			select {
			case <-timer.C:
				close(tq.exit)
				err := tq.Close(ctx)
				if err != nil {
					logger.GetLogger(ctx).Warn(err.Error())
				}
			}
		}
	}()

	return tq
}

func (t *TakeNumberQueue) checkIsEnd(ctx context.Context, value interface{}) bool {
	select {
	case <-t.exit:
		return true
	default:
		_, ok := value.(EndNumberValue)
		if ok {
			return true
		}
	}

	return false
}

// TakeNumber 取号
func (t *TakeNumberQueue) TakeNumber(ctx context.Context, number int32, value interface{}) error {
	t.timer.Reset(t.timeout)
	if number == t.offset {
		isEnd := t.checkIsEnd(ctx, value)
		if isEnd {
			t.timer.Stop()
			close(t.out)
			return nil
		}
		t.out <- &Result{
			Value:  value,
			Number: number,
		}
		t.offset += 1
		ok := true
		// 自旋
		for ok {
			var v interface{}
			v, ok = t.queue[t.offset]
			if ok {
				isEnd = t.checkIsEnd(ctx, v)
				if isEnd {
					logger.GetLogger(ctx).Info("task number queue is closed")
					t.timer.Stop()
					close(t.out)
					return nil
				}
				t.out <- &Result{
					Value:  v,
					Number: number,
				}
				t.offset += 1
			}
		}
	} else if number < t.offset {
		logger.GetLogger(ctx).Debug(fmt.Sprintf("number %+v is less than offset %+v", number, t.offset))
		isEnd := t.checkIsEnd(ctx, value)
		if isEnd {
			logger.GetLogger(ctx).Info("task number queue is closed")
			t.timer.Stop()
			close(t.out)
			return nil
		}
		t.out <- &Result{
			Value:  value,
			Number: number,
		}
	} else {
		t.queue[number] = value
	}

	return nil
}

// Call 叫号
func (t *TakeNumberQueue) Call(ctx context.Context) (<-chan *Result, error) {
	return t.out, nil
}

func (t *TakeNumberQueue) TakeEndNumber(ctx context.Context, number int32) error {
	logger.GetLogger(ctx).Info(fmt.Sprintf("task number queue end number is %+v", number))
	var end EndNumberValue = struct{}{}
	err := t.TakeNumber(ctx, number, end)
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭，退出queue
func (t *TakeNumberQueue) Close(ctx context.Context) error {
	select {
	case _, ok := <-t.out:
		if ok {
			close(t.out)
		}
	case <-t.exit:
	default:
		close(t.exit)
	}

	return nil
}
