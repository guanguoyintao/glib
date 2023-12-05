package uconcurrent

import (
	"context"
	"fmt"
	"git.umu.work/be/goframework/logger"
)

type ConcurrentControlExecutor struct {
	concurrencyControl chan struct{}
}

func NewConcurrentControlExecutor(maxConcurrency uint32) *ConcurrentControlExecutor {
	return &ConcurrentControlExecutor{
		concurrencyControl: make(chan struct{}, maxConcurrency),
	}
}

func (q *ConcurrentControlExecutor) Run(ctx context.Context, runnable Runnable) (interface{}, error) {
	logger.GetLogger(ctx).Debug(fmt.Sprintf("concurrent control executor enter, concurrent control queue len is %+v", cap(q.concurrencyControl)))
	q.concurrencyControl <- struct{}{}
	logger.GetLogger(ctx).Debug("concurrent control executor run")
	defer func() {
		<-q.concurrencyControl
	}()
	res, err := runnable(ctx)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}
	logger.GetLogger(ctx).Debug("runnable run end")

	return res, nil
}
