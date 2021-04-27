package core

import (
	"fmt"
	"github.com/dterei/gotsc"
	"time"
)

const CalibrationPeriod = uint64(1E7)
const WantedTicks = uint64(3E9)
const RealtimeOverheadRounds = 1000

type Clock interface {
	NextTick()
}

type HighPrecisionClock struct {
	freq uint64
	period uint64
	totalTicks uint64
	lastTsc uint64
	overruns uint64
	backwards uint64
	tickDeficit uint64
	lastCalibration time.Time
}

func (h *HighPrecisionClock) calibrate() {
	overhead := gotsc.TSCOverhead()
	toSum := uint64(0)
	for i := 0; i < RealtimeOverheadRounds; i++ {
		timeOverhead := gotsc.BenchStart()
		_ = time.Now()
		timeOverhead = gotsc.BenchEnd() - timeOverhead
		toSum += timeOverhead
	}
	overhead += toSum / RealtimeOverheadRounds

	start := time.Now()
	var realTime time.Time
	t0 := gotsc.BenchStart()
	var actualTicks uint64
	for {
		t1 := gotsc.BenchEnd()
		if t1 - t0 > WantedTicks {
			actualTicks = t1 - t0
			realTime = time.Now()
			break
		}
	}
	actualTicks -= overhead
	elapsed := realTime.Sub(start)
	h.period = (actualTicks / uint64(elapsed.Microseconds()) * h.freq) / 1e6
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
	h.period = uint64(float64(h.period) * (1.0 + factor)/2.0)
	h.lastCalibration = now
	h.totalTicks = 0
	fmt.Println(timeElapsed, factor, h.period, h.overruns)
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
	if h.tickDeficit > 0 {
		// We're running behind. Execute ticks without delay until we're caught up!
		h.tickDeficit--
		return
	}
	if t0 - h.lastTsc > h.period {
		h.overruns++
		h.tickDeficit = (t0 - h.lastTsc) / h.period
		h.lastTsc = gotsc.BenchStart()
		return
	}
	nextTick := t0 + h.period - t0 % h.period
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
	}
	c.calibrate()
	return c
}
