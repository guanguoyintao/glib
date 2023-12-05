package uwrapper

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ucontext"
	"git.umu.work/be/goframework/logger"
	"golang.org/x/sync/errgroup"
	"runtime"
	"sync"
	"time"
)

type wrapperFunc func(ctx context.Context) error

func GoroutineRecover(ctx context.Context) {
	r := recover()
	if r != nil {
		pc, file, line, ok := runtime.Caller(2)
		if !ok {
			logger.GetLogger(ctx).Error(fmt.Sprintf("goroutine err is: %+v", r))
			return
		}
		funcName := runtime.FuncForPC(pc).Name()

		logger.GetLogger(ctx).Error(fmt.Sprintf("goroutine err is: %+v", r),
			"file", file,
			"line", line,
			"funcName", funcName)
	}
}

func RetryWithBackoff(ctx context.Context, fn wrapperFunc, maxRetries int, maxDelay time.Duration) error {
	var retryDelay = time.Second
	var retries int
	for {
		err := fn(ctx)
		if err == nil {
			return nil
		}
		retries++
		if retries > maxRetries {
			return err
		}
		// Calculate the next retry delay using exponential backoff.
		retryDelay *= 2
		if retryDelay > maxDelay {
			retryDelay = maxDelay
		}

		time.Sleep(retryDelay)
	}
}

func GoExecWithSemaphore(ctx context.Context, f func(ctx context.Context) error, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer GoroutineRecover(ctx)
		err := f(ctx)
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
		}
	}()
}

func GoExecWithErrSemaphore(ctx context.Context, f func(ctx context.Context) error, eg *errgroup.Group) {
	eg.Go(func() error {
		defer GoroutineRecover(ctx)
		err := f(ctx)
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
			return err
		}

		return nil
	})
}

func LogWithStack(ctx context.Context, f func(ctx context.Context) error) error {
	err := f(ctx)
	if err != nil {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)

		callers := make([]uintptr, 10)
		n = runtime.Callers(0, callers)
		callers = callers[:n]

		frames := runtime.CallersFrames(callers)
		for {
			frame, more := frames.Next()
			logger.GetLogger(ctx).Error(fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))
			if !more {
				break
			}
		}

		return err
	}

	return nil
}

func GoExecHandlerWithoutTimeout(ctx context.Context, f wrapperFunc) {
	ctx = ucontext.NewUValueContext(ctx)
	go func() {
		defer GoroutineRecover(ctx)
		err := f(ctx)
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
		}
	}()

	return
}

func GoExecHandlerWithTimeout(ctx context.Context, f wrapperFunc, timeout time.Duration) {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(ctx, timeout)
	go func() {
		defer cancel() // 确保在函数退出时取消上下文
		defer GoroutineRecover(ctx)
		// 开始执行 f(ctx)
		result := make(chan error)
		exit := false
		go func(innerCtx context.Context) {
			err := f(innerCtx)
			if !exit {
				result <- err
			}
		}(ctx)
		// 监听 ctx.Done() 和结果通道
		select {
		case <-ctx.Done():
			exit = true
			return
		case err := <-result:
			if err != nil {
				logger.GetLogger(ctx).Error(err.Error())
			}
			exit = true
			return
		}
	}()
}
