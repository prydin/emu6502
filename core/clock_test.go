/*
 * Copyright (c) 2021 Pontus Rydin
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the Software
 * without restriction, including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
 * to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

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
	for i :=0; i < 2; i++ {
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