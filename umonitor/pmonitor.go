package umonitor

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uwrapper"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	hostname               string
	reqTotalMetricsName    = "umu_req_total"
	reqDurationMetricsName = "umu_req_duration"
	counterVecMap          map[string]*prometheus.CounterVec
	histogramVecMap        map[string]*prometheus.HistogramVec
)

type RequestLabels struct {
	ServiceName   string
	Client        string
	InterfaceName string
	ErrCode       int64
	ErrMsg        string
	Begin         time.Time
}

func init() {
	name, err := os.Hostname()
	if err != nil {
		hostname = "unknownhost"
	} else {
		hostname = name
	}
	counterVecMap = make(map[string]*prometheus.CounterVec)
	histogramVecMap = make(map[string]*prometheus.HistogramVec)
	initRequestMetrics()
}

func initRequestMetrics() {
	counterVecMap[reqTotalMetricsName] = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: reqTotalMetricsName,
			Help: reqTotalMetricsName,
		},
		[]string{"client", "server", "host", "interface", "err_code", "err_msg"})
	histogramVecMap[reqDurationMetricsName] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    reqDurationMetricsName,
		Help:    "Histogram for the runtime of a simple example function.",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	},
		[]string{"client", "server", "host", "interface", "err_code"})
	fmt.Println("register histogram:", reqDurationMetricsName)
	prometheus.MustRegister(histogramVecMap[reqDurationMetricsName])

}

// IncRequest 记录请求信息
func IncRequest(ctx context.Context, req *RequestLabels) {
	defer uwrapper.GoroutineRecover(ctx)
	counterVecMap[reqTotalMetricsName].With(prometheus.Labels{
		"client":    req.Client,
		"server":    req.ServiceName,
		"host":      hostname,
		"interface": req.InterfaceName,
		"err_code":  fmt.Sprintf("%d", req.ErrCode),
		"err_msg":   req.ErrMsg}).Inc()
	// 记录请求时长
	addReqDur(ctx, req)
}

// 记录请求耗时
func addReqDur(ctx context.Context, req *RequestLabels) {
	defer uwrapper.GoroutineRecover(ctx)
	timer := NewTimerWithBegin(histogramVecMap[reqDurationMetricsName].With(prometheus.Labels{
		"client":    req.Client,
		"server":    req.ServiceName,
		"host":      hostname,
		"interface": req.InterfaceName,
		"err_code":  fmt.Sprintf("%d", req.ErrCode)}), req.Begin)
	timer.ObserveDuration()
}
