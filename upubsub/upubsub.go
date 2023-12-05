package upubsub

import (
	"git.umu.work/AI/uglib/uconcurrent"
	"time"
)

type subscriberWorker struct {
	group      string
	topic      string
	h          SubscriberFunc
	executor   *uconcurrent.ConcurrentControlExecutor
	maxRetries int
	maxDelay   time.Duration
}
