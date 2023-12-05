package utest

import (
	"context"
	"time"
)

type B struct {
	preheat   interface{}
	startTime time.Time
}

func (b *B) ResetStartTime() {
	b.startTime = time.Now()
}

func (b *B) getSpendTime() time.Duration {
	spendTime := time.Since(b.startTime)

	return spendTime
}

func (b *B) setPreheat(preheat interface{}) {
	b.preheat = preheat
}

func (b *B) GetPreheat() interface{} {
	return b.preheat
}

type Resource interface {
	Close(ctx context.Context) error
	Load(ctx context.Context, index int) interface{}
	Len(ctx context.Context) int
}

type TestFunc func(ctx context.Context) error

type TestBenchmark struct {
	preheat func(ctx context.Context) (Resource, error)
	f       func(ctx context.Context, b *B) error
}

func NewBenchmark(f func(ctx context.Context, b *B) error) *TestBenchmark {
	return &TestBenchmark{
		preheat: nil,
		f:       f,
	}
}
