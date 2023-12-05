package uexp

import (
	"git.umu.work/AI/uglib/ufunc-operator/sort"
	"sort"
	"time"
)

func SumFloat64(data []float64) float64 {
	var result float64
	for _, d := range data {
		result += d
	}

	return result
}

func SumDuration(data []time.Duration) time.Duration {
	var result time.Duration
	for _, d := range data {
		result += d
	}

	return result
}

func MaxDuration(data []time.Duration) time.Duration {
	var result time.Duration
	for _, d := range data {
		if d > result {
			result = d
		}
	}

	return result
}

func MinDuration(data []time.Duration) time.Duration {
	if len(data) == 0 {
		return 0 * time.Second
	}
	result := data[0]
	for _, d := range data {
		if d < result {
			result = d
		}
	}

	return result
}

func AvgDuration(data []time.Duration) time.Duration {
	var sum time.Duration
	for _, d := range data {
		sum += d
	}
	result := time.Duration(float64(sum.Microseconds())/float64(len(data))) * time.Microsecond

	return result
}

func PercentileDuration(data []time.Duration, th float64) []time.Duration {
	sort.Sort(usort.DurationList(data))
	percentileIndex := int(float64(len(data)) * th)
	if percentileIndex == 0 {
		return data
	}

	return data[:percentileIndex]
}
