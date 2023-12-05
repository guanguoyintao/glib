package uconcurrent

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
)

// Runnable 带返回值和错误的标准可执行单元
type Runnable func(ctx context.Context) (interface{}, error)

// RunnableWithoutResult 不带返回值的只有错误的可执行单元
type RunnableWithoutResult func(ctx context.Context) error

func (r *RunnableWithoutResult) WithTimeout(ctx context.Context) error {
	done := make(chan struct{})
	errChan := make(chan error)
	go func() {
		f := *r
		err := f(ctx)
		if err != nil {
			errChan <- err
		} else {
			done <- struct{}{}
		}
	}()
	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return uerrors.UErrorTimeout
	}
}
