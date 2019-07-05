package util

import (
	"time"
)

func Next(t time.Time, interval time.Duration) time.Time {
	ts := time.Duration(t.Hour()*3600+t.Minute()*60+t.Second()) * time.Second
	div := ts / interval
	next := time.Duration((div + 1) * interval)
	return t.Add(next - ts)
}

func Interval(next, end time.Time) time.Duration {
	if interval := next.Sub(end); interval > 0 {
		return time.Duration(interval.Nanoseconds())
	}
	return 0
}
