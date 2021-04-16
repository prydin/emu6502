package core

import (
	"fmt"
	"github.com/dterei/gotsc"
	"time"
)

const CalibrationPeriod = uint64(1E7)

type Clock interface {
	NextTick()
}

type HighPrecisionClock struct {
	freq uint64
	period uint64
	totalTicks uint64
	lastTsc uint64
	lastCalibration time.Time
}

func (h *HighPrecisionClock) calibrate() {
	oh := gotsc.TSCOverhead()
	t0 := gotsc.BenchStart()
	time.Sleep(1000 * time.Millisecond)
	t1 := gotsc.BenchEnd()
	scale := t1 - t0 - oh
	h.period = (scale * h.period) /1000000000
}

func (h *HighPrecisionClock) recalibrate() {
	if h.lastCalibration.IsZero() {
		h.lastCalibration = time.Now()
		h.totalTicks = 0
		return
	}
	now := time.Now()
	timeElapsed := now.Sub(h.lastCalibration)
	factor := float64(1E9*CalibrationPeriod/h.freq) / float64(timeElapsed)
	h.period = uint64(float64(h.period) * (2.0 + factor)/3)
	h.lastCalibration = now
	h.totalTicks = 0
	fmt.Println(timeElapsed, factor, h.period)
}

func (h *HighPrecisionClock) NextTick() {
	h.totalTicks++
	if h.totalTicks >= CalibrationPeriod {
		h.recalibrate()
	}
	t0 := gotsc.BenchStart()
	if h.lastTsc == 0 {
		h.lastTsc = t0
	}
	nextTick := h.lastTsc + h.period
	t1 := gotsc.BenchStart()
	for t1 < nextTick {
		t1 = gotsc.BenchStart()
	}
	h.lastTsc = t1
}

func NewClock(freq uint64) Clock {
	// TODO: Create lower precision clock on non-x86 architectures
	c := &HighPrecisionClock{
		freq: freq,
		period: 1000000000 / freq,
	}
	c.calibrate()
	return c
}
