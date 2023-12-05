package utopk

var (
	defaultTopKFloat64 TopKFloat64 = &bfprtFloat64{}
)

func MaxTopKFloat64(array []float64, k int) []float64 {
	return defaultTopKFloat64.MaxTopK(array, k)
}

func MinTopKFloat64(array []float64, k int) []float64 {
	return defaultTopKFloat64.MinTopK(array, k)
}
