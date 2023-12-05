package umonitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Timer struct {
	begin    time.Time
	observer prometheus.Observer
}

func NewTimerWithBegin(o prometheus.Observer, begin time.Time) *Timer {
	return &Timer{
		begin:    begin,
		observer: o,
	}
}

func (t *Timer) ObserveDuration() time.Duration {
	d := time.Since(t.begin)
	if t.observer != nil {
		t.observer.Observe(d.Seconds())
	}
	return d
}
