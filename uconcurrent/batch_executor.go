package uconcurrent

import (
	"context"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/logger"
	"sync"
	"sync/atomic"
	"time"
)

type innerResult struct {
	Index  int
	Result interface{}
}

type BatchExecutor struct {
	// 是否超时
	HasTimeout bool
	// 结果集，对应AddRunnable时的顺序
	Results []interface{}
	err     chan error
	// 操作成功个数
	SuccessCount uint32
	// 操作返回error个数
	ErrorCount uint32

	wg           sync.WaitGroup
	runnables    []Runnable
	resultsChan  chan *innerResult
	success      chan int
	runnableExit chan struct{}
}

func NewBatchExecutor() *BatchExecutor {
	return &BatchExecutor{
		wg: sync.WaitGroup{},
	}
}

// AddRunnable 添加需并行执行的可执行单元
func (b *BatchExecutor) AddRunnable(runnable Runnable) {
	b.runnables = append(b.runnables, runnable)
}

// Run 并行执行所有已添加的可执行单元，等待最长timeout时间后结束
func (b *BatchExecutor) Run(ctx context.Context, timeout time.Duration) error {
	// 初始化
	b.Results = make([]interface{}, len(b.runnables))
	b.resultsChan = make(chan *innerResult, len(b.runnables))
	b.success = make(chan int)
	b.err = make(chan error)
	// 设置超时
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()
	// 批量执行
	go b.executeRunnables(ctx)
	// 等待全部完成或超时
	for {
		select {
		case result := <-b.resultsChan:
			b.Results[result.Index] = result.Result
		case <-b.success:
			b.HasTimeout = false
			return nil
		case <-timeoutCtx.Done():
			b.HasTimeout = true
			return nil
		case err := <-b.err:
			return err
		}
	}
}

func (b *BatchExecutor) executeRunnables(ctx context.Context) {
	b.wg.Add(len(b.runnables))
	for index, runnable := range b.runnables {
		f := runnable
		i := index
		go func() {
			defer uwrapper.GoroutineRecover(ctx)
			result, err := f(ctx)
			if err != nil {
				logger.GetLogger(ctx).Error(err.Error())
				atomic.AddUint32(&b.ErrorCount, 1)
				b.err <- err
			} else {
				atomic.AddUint32(&b.SuccessCount, 1)
				b.resultsChan <- &innerResult{Index: i, Result: result}
			}
			b.wg.Done()
		}()
	}

	b.wg.Wait()
	if !b.HasTimeout {
		b.success <- 1
	}
}
