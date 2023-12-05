package usort

import (
	"time"
)

type DurationList []time.Duration

func (s DurationList) Len() int {
	return len(s)
}
func (s DurationList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s DurationList) Less(i, j int) bool {
	return s[i] < s[j]
}
