package usubtitle

import "time"

type Item struct {
	No      int
	StartAt time.Duration
	EndAt   time.Duration
	Texts   []string
}
