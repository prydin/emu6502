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

package sid

import (
	"github.com/cnkei/gospline"
)

// Maximum cutoff frequency is specified as
// FCmax = 2.6e-5/C = 2.6e-5/2200e-12 = 11818.
//
// Measurements indicate a cutoff frequency range of approximately
// 220Hz - 18kHz on a MOS6581 fitted with 470pF capacitors. The function
// mapping FC to cutoff frequency has the shape of the tanh function, with
// a discontinuity at FCHI = 0x80.
// In contrast, the MOS8580 almost perfectly corresponds with the
// specification of a linear mapping from 30Hz to 12kHz.
//
// The mappings have been measured by feeding the SID with an external
// signal since the chip itself is incapable of generating waveforms of
// higher fundamental frequency than 4kHz. It is best to use the bandpass
// output at full resonance to pick out the cutoff frequency at any given
// FC setting.
//
// The mapping function is specified with spline interpolation points and
// the function values are retrieved via table lookup.
//
// NB! Cutoff frequency characteristics may vary, we have modeled two
// particular Commodore 64s.

var f0Points6581X = []float64{
	0, 0, 128, 256, 384, 512, 640, 768, 832, 896, 960, 992, 1008, 1016, 1023, 1023, 1024,
	1024, 1032, 1056, 1088, 1120, 1152, 1280, 1408, 1536, 1664, 1792, 1920, 2047, 2047}

var f0Points6581Y = []float64{
	220, 220, 230, 250, 300, 420, 780, 1600, 2300, 3200, 4300, 5000, 5400, 5700, 6000, 6000,
	4600, 4600, 4800, 5300, 6000, 6600, 7200, 9500, 12000, 14500, 16000, 17100, 17700, 18000, 18000}

var f0Points8580X = []float64{
	0, 0, 128, 256, 384, 512, 640, 768, 896, 1024, 1152, 1280, 1408, 1536, 1664, 1792, 1920, 2047, 2047}

var f0Points8580Y = []float64{0, 0, 800, 1600, 2500, 3300, 4100, 4800, 5600, 6500, 7500, 8400, 9200, 9800,
	10500, 11000, 11700, 12500, 12500}

var curve6581 gospline.Spline
var curve8580 gospline.Spline

type Filter struct {
	// User accessible filters
	enabled   bool  // Filter enabled.
	fc        reg12 // Filter cutoff frequency.
	res       reg8  // Filter resonance.
	filt      reg8  // Selects which inputs to route through filter.
	voice3off bool  // Switch voice 3 off.
	mode      reg8  // Highpass, bandpass, and lowpass filter modes.
	vol       reg4  // Master volume

	mixerDC soundSample // Mixer DC offset.

	// State of filter.
	vhp soundSample // highpass
	vbp soundSample // bandpass
	vlp soundSample // lowpass
	vnf soundSample // not filtered

	// Cutoff frequency, resonance.
	w0, w0Ceil1, w0CeilDt soundSample
	Q1024div              soundSample

	// Cutoff frequency tables.
	// FC is an 11 bit register.
	/*sound_sample f0_6581[2048]
	sound_sample f0_8580[2048]
	sound_sample* f0
	fc_point* f0_points
	int f0_count */
}

func (f *Filter) reset() {
	// Reset user registers
	f.fc = 0
	f.res = 0
	f.filt = 0
	f.voice3off = false
	f.mode = 0
	f.vol = 0

	// State of filter.
	f.vhp = 0
	f.vbp = 0
	f.vlp = 0
	f.vnf = 0

	// Initialize filter parameters
	//f.setW0()
	//f.setQ()
}

/*

func (f *Filter) setW0() {
	// Multiply with 1.048576 to facilitate division by 1 000 000 by right-
	// shifting 20 times (2 ^ 20 = 1048576).
	w0 := soundSample(2 * math.Pi * f0[f.fc] * 1.048576)

	// Limit f0 to 16kHz to keep 1 cycle filter stable.
	const sound_sample w0_max_1 = static_cast < sound_sample > (2 * pi * 16000 * 1.048576)
	w0_ceil_1 = w0 <= w0_max_1 ? w0:
	w0_max_1

	// Limit f0 to 4kHz to keep delta_t cycle filter stable.
	const sound_sample w0_max_dt = static_cast < sound_sample > (2 * pi * 4000 * 1.048576)
	w0_ceil_dt = w0 <= w0_max_dt ? w0:
	w0_max_dt
}

func init() {
	// Pre-calculate smoothed filter curves
	curve6581 = gospline.NewCubicSpline(f0Points6581X, f0Points6581Y)
	curve8580 = gospline.NewCubicSpline(f0Points8580X, f0Points8580Y)
}
*/
