package utopk

type TopKInt interface {
	MaxTopK(array []int, k int) []int
	MinTopK(array []int, k int) []int
}

type TopKFloat64 interface {
	MaxTopK(array []float64, k int) []float64
	MinTopK(array []float64, k int) []float64
}
