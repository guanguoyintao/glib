package utopk

import "sort"

const cutoff = 15

var _ TopKFloat64 = (*bfprtFloat64)(nil)

type bfprtFloat64 struct {
}

func (b *bfprtFloat64) MinTopK(arr []float64, k int) []float64 {
	n := len(arr)

	// 递归结束条件
	// 优化1：使用插入排序对小数组进行排序
	if n <= cutoff {
		sort.Float64s(arr)
		return arr[:k]
	}

	// 优化2：将数组划分为若干个大小为5的子数组，对每个子数组进行排序，取出中位数组成新的数组
	subArrCount := (n + 4) / 5
	subArrs := make([][]float64, subArrCount)
	for i := 0; i < subArrCount; i++ {
		start := i * 5
		end := start + 5
		if end > n {
			end = n
		}
		subArrs[i] = arr[start:end]
		sort.Float64s(subArrs[i])
	}

	medians := make([]float64, subArrCount)
	for i := 0; i < subArrCount; i++ {
		medians[i] = subArrs[i][len(subArrs[i])/2]
	}

	// 优化3：递归调用bfprt算法，找到中位数的中位数作为pivot
	pivot := b.MinTopK(medians, len(medians)/2)[0]

	// 优化4：根据pivot将数组分为3部分
	var (
		lessThan    []float64
		equalTo     []float64
		greaterThan []float64
	)
	for _, num := range arr {
		switch {
		case num < pivot:
			lessThan = append(lessThan, num)
		case num == pivot:
			equalTo = append(equalTo, num)
		case num > pivot:
			greaterThan = append(greaterThan, num)
		}
	}

	// 优化5：根据划分后的数组长度关系，确定下一步递归的数组
	switch {
	case k < len(lessThan):
		return b.MinTopK(lessThan, k)
	case k < len(lessThan)+len(equalTo):
		return equalTo
	default:
		return append(equalTo, b.MinTopK(greaterThan, k-len(lessThan)-len(equalTo))...)
	}
}

func (b *bfprtFloat64) MaxTopK(a []float64, k int) []float64 {
	n := len(a)
	arr := make([]float64, len(a), len(a))
	copy(arr, a)

	// 递归结束条件
	// 优化1：使用插入排序对小数组进行排序
	if n <= cutoff {
		sort.Sort(sort.Reverse(sort.Float64Slice(arr)))
		return arr[:k]
	}

	// 优化2：将数组划分为若干个大小为5的子数组，对每个子数组进行排序，取出中位数组成新的数组
	subArrCount := (n + 4) / 5
	subArrs := make([][]float64, subArrCount)
	for i := 0; i < subArrCount; i++ {
		start := i * 5
		end := start + 5
		if end > n {
			end = n
		}
		subArrs[i] = arr[start:end]
		sort.Sort(sort.Reverse(sort.Float64Slice(subArrs[i])))
	}

	medians := make([]float64, subArrCount)
	for i := 0; i < subArrCount; i++ {
		medians[i] = subArrs[i][len(subArrs[i])/2]
	}

	// 优化3：递归调用bfprt算法，找到中位数的中位数作为pivot
	pivot := b.MaxTopK(medians, len(medians)/2)[0]

	// 优化4：根据pivot将数组分为3部分
	var (
		lessThan    []float64
		equalTo     []float64
		greaterThan []float64
	)
	for _, num := range arr {
		switch {
		case num < pivot:
			greaterThan = append(greaterThan, num)
		case num == pivot:
			equalTo = append(equalTo, num)
		case num > pivot:
			lessThan = append(lessThan, num)
		}
	}

	// 优化5：根据划分后的数组长度关系，确定下一步递归的数组
	switch {
	case k < len(lessThan):
		return b.MaxTopK(lessThan, k)
	case k < len(lessThan)+len(equalTo):
		return equalTo
	default:
		return append(b.MaxTopK(greaterThan, k-len(lessThan)-len(equalTo)), equalTo...)
	}
}
