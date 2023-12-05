package uconcurrent

import (
	"context"
	"git.umu.work/AI/uglib/ucontext"
	"git.umu.work/AI/uglib/uwrapper"
	"sync"
	"time"
)

type ErrorExecutor struct {
	err error
	wg  sync.WaitGroup
	rw  sync.RWMutex
}

func NewErrorExecutor() *ErrorExecutor {
	return &ErrorExecutor{
		wg: sync.WaitGroup{},
		rw: sync.RWMutex{},
	}
}

// AddAsyncRunnable 添加异步并发可执行单元
func (e *ErrorExecutor) AddAsyncRunnable(ctx context.Context, runnable RunnableWithoutResult, timeout time.Duration) {
	e.wg.Add(1)
	go func(innerCtx context.Context) {
		innerCtx, _ = context.WithTimeout(innerCtx, timeout)
		defer uwrapper.GoroutineRecover(ctx)
		defer e.wg.Done()
		innerErr := runnable.WithTimeout(innerCtx)
		if innerErr != nil {
			e.rw.Lock()
			e.err = innerErr
			e.rw.Unlock()
		}
	}(ucontext.NewUValueContext(ctx))
}

// Wait 等待异步任务全部结束，堵塞操作
func (e *ErrorExecutor) Wait(ctx context.Context) error {
	e.wg.Wait()
	e.rw.RLock()
	err := e.err
	e.rw.RUnlock()
	if err != nil {
		return err
	}

	return nil
}
