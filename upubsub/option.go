package upubsub

import (
	"context"
	goframeworkMQ "git.umu.work/be/goframework/async/mq"
	"git.umu.work/be/goframework/metrics"
	"time"
)

type Metrics struct {
	consumerTimer    *metrics.Timer
	producerTimer    *metrics.Timer
	consumerRecorder metrics.Recorder
	producerRecorder metrics.Recorder
}

func (m *Metrics) ConsumerMonitorStart(ctx context.Context) {
	m.consumerRecorder = m.consumerTimer.Timer()
}

func (m *Metrics) ConsumerMonitorRecorder(ctx context.Context, name, topic, group, key string) {
	m.consumerRecorder(name, topic, group, key)
}

func (m *Metrics) ProducerMonitorStart(ctx context.Context) {
	m.producerRecorder = m.producerTimer.Timer()
}

func (m *Metrics) ProducerMonitorRecorder(ctx context.Context, name, topic, key string) {
	m.producerRecorder(name, topic, key)
}

type Options struct {
	Timeout        time.Duration
	PartitionKey   string
	MonitorMetrics *Metrics
}

type Option func(*Options)

func WithMonitor() Option {
	return func(o *Options) {
		o.MonitorMetrics = &Metrics{
			consumerTimer: goframeworkMQ.GetConsumerTimerMetric(),
			producerTimer: goframeworkMQ.GetProducerTimerMetric(),
		}
	}
}

func WithPartitionKey(key string) Option {
	return func(o *Options) {
		o.PartitionKey = key
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}
