package core

import (
	"fmt"
	"testing"
	"time"
)

func TestHighPrecisionClock_NextTick(t *testing.T) {
	freq := uint64(1000000)
	c := &HighPrecisionClock{
		freq: freq,
		period: 1000000000 / freq,
	}
	c.calibrate()
	start := time.Now()
	for i := 0; i < 2E6; i++ {
		c.NextTick()
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Elapsed time: %s", elapsed)
}