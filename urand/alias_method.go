package urand

// @Title   Alias Method实现
// @Description  Alias Method, 算法的时间复杂度为 O(1)

import "math/rand"

type AliasMethod struct {
	prob  []float64
	alias []int
	n     int
}

func NewAliasMethod(p []float64) *AliasMethod {
	n := len(p)
	prob := make([]float64, n)
	alias := make([]int, n)

	sum := 0.0
	for i := 0; i < n; i++ {
		sum += p[i]
	}
	for i := 0; i < n; i++ {
		prob[i] = p[i] * float64(n) / sum
	}

	small := make([]int, 0)
	large := make([]int, 0)
	for i := 0; i < n; i++ {
		if prob[i] < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	for len(small) > 0 && len(large) > 0 {
		l := small[len(small)-1]
		small = small[:len(small)-1]
		g := large[len(large)-1]
		large = large[:len(large)-1]
		alias[l] = g
		prob[g] = prob[g] - (1.0 - prob[l])
		if prob[g] < 1.0 {
			small = append(small, g)
		} else {
			large = append(large, g)
		}
	}

	for len(large) > 0 {
		g := large[len(large)-1]
		large = large[:len(large)-1]
		prob[g] = 1.0
	}

	return &AliasMethod{prob, alias, n}
}

func (am *AliasMethod) Rand() int {
	i := rand.Intn(am.n)
	if rand.Float64() < am.prob[i] {
		return i
	} else {
		return am.alias[i]
	}
}
