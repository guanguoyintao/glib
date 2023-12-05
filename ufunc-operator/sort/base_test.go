package usort

import (
	"github.com/bmizerany/assert"
	"sort"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	tests := []struct {
		name string
		list DurationList
		want DurationList
	}{
		{
			name: "1",
			list: DurationList{
				1 * time.Hour,
				1 * time.Second,
				1 * time.Minute,
				2 * time.Second,
			},
			want: DurationList{
				1 * time.Second,
				2 * time.Second,
				1 * time.Minute,
				1 * time.Hour,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.list)
			assert.Equal(t, tt.want, tt.list)
		})
	}
}
