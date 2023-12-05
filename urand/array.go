package urand

import (
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type Item struct {
	Probability float64
	Data        interface{}
}

func RandDiscreteProbabilityDistributionArray(items []*Item) (*Item, error) {
	n := len(items)
	p := make([]float64, n)
	allZero := true // flag to check if all probabilities are zero
	for i := 0; i < n; i++ {
		p[i] = items[i].Probability
		if items[i].Probability != 0 {
			allZero = false
		}
	}
	if allZero {
		return nil, errors.New("all probabilities are zero")
	}

	am := NewAliasMethod(p)
	i := am.Rand()

	return items[i], nil
}

func RandUniformDistributionArray(array []interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(array))

	return array[index]
}
