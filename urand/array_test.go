package urand

import (
	"math"
	"testing"
)

func TestRandItem(t *testing.T) {
	const numTests = 100000
	const probTolerance = 0.01

	tests := []struct {
		name  string
		items []*Item
	}{
		{
			name: "test1",
			items: []*Item{
				{Probability: 0.5, Data: "A"},
				{Probability: 0.5, Data: "B"},
			},
		},
		{
			name: "test2",
			items: []*Item{
				{Probability: 0.25, Data: "A"},
				{Probability: 0.25, Data: "B"},
				{Probability: 0.25, Data: "C"},
				{Probability: 0.25, Data: "D"},
			},
		},
		{
			name: "test3",
			items: []*Item{
				{Probability: 0.1, Data: "A"},
				{Probability: 0.3, Data: "B"},
				{Probability: 0.2, Data: "C"},
				{Probability: 0.15, Data: "D"},
				{Probability: 0.25, Data: "E"},
			},
		},
		{
			name: "test4",
			items: []*Item{
				{Probability: 0, Data: "A"},
				{Probability: 0, Data: "B"},
				{Probability: 0, Data: "C"},
				{Probability: 0, Data: "D"},
				{Probability: 1, Data: "E"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itemCounts := make(map[interface{}]int)

			for i := 0; i < numTests; i++ {
				item, err := RandItem(tt.items)
				if err != nil {
					t.Errorf(err.Error())
				}
				itemCounts[item.Data]++
			}

			// Check that the probability of each item is within the tolerance range
			for _, item := range tt.items {
				actualProb := float64(itemCounts[item.Data]) / float64(numTests)
				if math.Abs(actualProb-item.Probability) > probTolerance {
					t.Errorf("test %q: item %v has probability %v, expected %v", tt.name, item.Data, actualProb, item.Probability)
				}
			}
		})
	}

	// Test for edge case when all probabilities are zero
	// Test for edge case when all probabilities are zero
	t.Run("test5", func(t *testing.T) {
		items := []*Item{
			{Probability: 0, Data: "A"},
			{Probability: 0, Data: "B"},
			{Probability: 0, Data: "C"},
		}
		_, err := RandItem(items)
		if err == nil {
			t.Errorf("expected an error when all probabilities are zero")
		}
	})
}
