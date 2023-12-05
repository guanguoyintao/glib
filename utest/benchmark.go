package utest

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uexp"
	"github.com/liushuochen/gotable"
	"sync"
	"sync/atomic"
	"time"
)

type Rt struct {
	Rt50  int64
	Rt75  int64
	Rt90  int64
	MinRt int64
	MaxRt int64
	AvgRt int64
}

type Metric struct {
	CompleteNumber uint32
	Concurrent     uint32
	Qps90          float64
	QPS            float64
	Rt             Rt
	Error          float64
}

func (tb *TestBenchmark) PreHeat(f func(ctx context.Context) (Resource, error)) {
	tb.preheat = f
}

func (tb *TestBenchmark) Testing(ctx context.Context, total int32, concurrence int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	var counter int32 = 0
	var errorCounter uint32 = 0
	timeConsumeList := make([]time.Duration, 0, concurrence)
	for {
		var resource Resource
		var err error
		if tb.preheat != nil {
			resource, err = tb.preheat(ctx)
			if err != nil {
				panic(err)
			}
		}
		for i := 0; i < concurrence; i++ {
			timeoutCtx, cancelFunc := context.WithTimeout(ctx, timeout)
			wg.Add(1)
			index := i
			go func() {
				b := &B{}
				defer wg.Done()
				defer cancelFunc()
				b.ResetStartTime()
				if resource != nil {
					b.setPreheat(resource.Load(ctx, index))
				}
				err = tb.f(timeoutCtx, b)
				spendTime := b.getSpendTime()
				if err != nil {
					atomic.AddUint32(&errorCounter, 1)
				} else {
					timeConsumeList = append(timeConsumeList, spendTime)
				}
				atomic.AddInt32(&counter, 1)
			}()
			if atomic.LoadInt32(&counter) > total {
				cancelFunc()
				break
			}
		}
		wg.Wait()
		if resource != nil {
			resource.Close(ctx)
		}
		if atomic.LoadInt32(&counter) > total {
			fmt.Printf("Completed %d\n", total)
			break
		} else {
			fmt.Printf("Completed %d\n", atomic.LoadInt32(&counter))
		}
		time.Sleep(2 * time.Second)
	}
	timeConsumeList50 := uexp.PercentileDuration(timeConsumeList, 0.5)
	timeConsumeList75 := uexp.PercentileDuration(timeConsumeList, 0.75)
	timeConsumeList90 := uexp.PercentileDuration(timeConsumeList, 0.9)
	ms := []Metric{
		{
			CompleteNumber: uint32(total),
			Concurrent:     uint32(concurrence),
			QPS:            float64(len(timeConsumeList)) / uexp.SumDuration(timeConsumeList).Seconds(),
			Qps90:          float64(len(timeConsumeList90)) / uexp.SumDuration(timeConsumeList90).Seconds(),
			Rt: Rt{
				Rt50:  uexp.MaxDuration(timeConsumeList50).Milliseconds(),
				Rt75:  uexp.MaxDuration(timeConsumeList75).Milliseconds(),
				Rt90:  uexp.MaxDuration(timeConsumeList90).Milliseconds(),
				MinRt: uexp.MinDuration(timeConsumeList).Milliseconds(),
				MaxRt: uexp.MaxDuration(timeConsumeList).Milliseconds(),
				AvgRt: uexp.AvgDuration(timeConsumeList).Milliseconds(),
			},
			Error: float64(errorCounter) / float64(atomic.LoadInt32(&counter)),
		},
	}
	err := tb.print(ms)
	if err != nil {
		panic(err)
	}
}

func (tb *TestBenchmark) print(ms []Metric) error {
	completeNumberColumn := "completed num"
	concurrentNumColumn := "concurrent num"
	qpsColumn90 := "90 qps"
	qpsColumn := "qps"
	rt50Column := "50 rt(ms)"
	rt75Column := "75 rt(ms)"
	rt90Column := "90 rt(ms)"
	minRtColumn := "min rt(ms)"
	maxRtColumn := "max rt(ms)"
	avgRtColumn := "avg rt(ms)"
	errorRateColumn := "error rate"
	table, err := gotable.Create(completeNumberColumn, concurrentNumColumn, qpsColumn, qpsColumn90, minRtColumn,
		maxRtColumn, avgRtColumn, errorRateColumn, rt50Column, rt75Column, rt90Column)
	if err != nil {
		return err
	}

	for _, m := range ms {
		row := make(map[string]string)
		row[completeNumberColumn] = fmt.Sprintf("%d", m.CompleteNumber)
		row[concurrentNumColumn] = fmt.Sprintf("%d", m.Concurrent)
		row[qpsColumn] = fmt.Sprintf("%tb", m.QPS)
		row[qpsColumn90] = fmt.Sprintf("%tb", m.Qps90)
		row[rt50Column] = fmt.Sprintf("%d", m.Rt.Rt50)
		row[rt75Column] = fmt.Sprintf("%d", m.Rt.Rt75)
		row[rt90Column] = fmt.Sprintf("%d", m.Rt.Rt90)
		row[minRtColumn] = fmt.Sprintf("%d", m.Rt.MinRt)
		row[maxRtColumn] = fmt.Sprintf("%d", m.Rt.MaxRt)
		row[avgRtColumn] = fmt.Sprintf("%d", m.Rt.AvgRt)
		row[errorRateColumn] = fmt.Sprintf("%tb", m.Error)
		err := table.AddRow(row)
		if err != nil {
			return err
		}
	}

	table.CloseBorder()
	fmt.Println(table.String())

	return nil
}
