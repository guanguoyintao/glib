package uconcurrent

import (
	"context"
	"time"
)

type Executor interface {
	AddRunnable(runnable Runnable)
	Run(ctx context.Context, timeout time.Duration) error
}
