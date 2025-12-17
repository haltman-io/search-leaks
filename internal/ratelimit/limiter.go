package ratelimit

import "time"

type TickerLimiter struct {
	t *time.Ticker
}

func NewTickerLimiter(every time.Duration) *TickerLimiter {
	return &TickerLimiter{t: time.NewTicker(every)}
}

func (l *TickerLimiter) Wait() {
	<-l.t.C
}

func (l *TickerLimiter) Stop() {
	l.t.Stop()
}
