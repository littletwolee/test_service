package util

import (
	"fmt"
	"testing"
	"time"
)

func Test_Next(t *testing.T) {
	tt := Next(time.Date(2019, 01, 01, 01, 02, 10, 0, time.UTC), time.Duration(1800)*time.Second)
	fmt.Println(tt)
}

func Test_Interval(t *testing.T) {
	next := time.Date(2019, 01, 01, 01, 02, 10, 0, time.UTC)
	end := time.Date(2019, 01, 01, 01, 01, 10, 0, time.UTC)
	fmt.Println(Interval(next, end))
}
