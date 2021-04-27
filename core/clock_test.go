package core

import (
	"fmt"
	"github.com/dterei/gotsc"
	"testing"
	"time"
)

func TestTSC(t *testing.T) {

	high := uint64(0)
	low := uint64(1e15)
	for i :=0; i < 10; i++ {
		t1 := gotsc.BenchStart()
		time.Sleep(10 * time.Second)
		t2 := gotsc.BenchEnd()
		elapsed := t2 - t1
		if elapsed > high {
			high = elapsed
		}
		if elapsed < low {
			low = elapsed
		}
		fmt.Println(float64(t2 - t1))
	}
	fmt.Println(high - low)
}

func TestHighPrecisionClock_NextTick(t *testing.T) {
	freq := uint64(1000000)
	c := &HighPrecisionClock{
		freq: freq,
		period: 1000000000 / freq,
	}
	c.calibrate()
	start := time.Now()
	for i := 0; i < 60E6; i++ {
		c.NextTick()
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Elapsed time: %s", elapsed)
}