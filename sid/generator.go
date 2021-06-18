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

type Source interface {
	ReadOutput() soundSample
}

type Generator struct {
	test     bool // Test mode (voice disabled, pulse output perpetual high)
	ringMod  bool // Ring modulator enabled
	waveform reg8 // Waveform register

	syncTrigger bool // Triggers voice synch

	accumulator   reg24 // Counts up to 2^12 once every period
	shiftRegister reg24 // Shift register for

	freq       reg16 // Fout  = (Fn*Fclk/16777216)Hz
	pulseWidth reg12

	syncSource *Generator
	syncTarget *Generator

	// Wave combination tables
	wavePS  []reg12
	wavePT  []reg12
	waveST  []reg12
	wavePST []reg12
}

func NewGenerator() *Generator {
	g := &Generator{
		wavePS:  wave8580PS,
		wavePT:  wave8580PT,
		waveST:  wave8580ST,
		wavePST: wave8580PST,
	}
	g.syncSource = g
	g.reset()
	return g
}

func (g *Generator) Clock() {
	// No operation if test bit is set.
	if g.test {
		return
	}

	oldAcc := g.accumulator

	// Calculate new accumulator value
	g.accumulator += reg24(g.freq)
	g.accumulator &= 0xffffff

	// Did we flip bit 19?
	g.syncTrigger = oldAcc&0x800000 == 0 && (g.accumulator&0x800000) != 0

	// Shift noise register once for each time accumulator bit 19 is set high.
	if g.syncTrigger {
		bit0 := ((g.shiftRegister >> 22) ^ (g.shiftRegister >> 17)) & 0x1
		g.shiftRegister <<= 1
		g.shiftRegister &= 0x7fffff
		g.shiftRegister |= bit0
	}
}

func (g *Generator) genSawtooth() reg12 {
	return reg12(g.accumulator >> 12)
}

func (g *Generator) genTriangle() reg12 {
	msb := g.accumulator
	if g.ringMod {
		msb ^= g.syncSource.accumulator
	}
	msb &= 0x800000
	if msb != 0 {
		return reg12(^g.accumulator>>11) & 0xfff
	} else {
		return reg12(g.accumulator>>11) & 0xfff
	}
}

func (g *Generator) genPulse() reg12 {
	// If test bit is set, we lock the output high. Otherwise, output high
	// if accumulator is greater than pulseWidth.
	if g.test || reg12(g.accumulator>>12) >= g.pulseWidth {
		return 0xfff
	} else {
		return 0x000
	}
}

// Generate noise using a simple LFSR. This is apparently how the real SID does it.
func (g *Generator) genNoise() reg12 {
	return reg12((g.shiftRegister&0x400000)>>11 |
		(g.shiftRegister&0x100000)>>10 |
		(g.shiftRegister&0x010000)>>7 |
		(g.shiftRegister&0x002000)>>5 |
		(g.shiftRegister&0x000800)>>4 |
		(g.shiftRegister&0x000080)>>1 |
		(g.shiftRegister&0x000010)<<1 |
		(g.shiftRegister&0x000004)<<2)
}

func (g *Generator) ReadOutput() reg12 {
	switch g.waveform {
	default:
	case 0x0:
		return 0
	case 0x1:
		return g.genTriangle()
	case 0x2:
		return g.genSawtooth()
	case 0x3:
		return g.waveST[g.genSawtooth()] << 4
	case 0x4:
		return g.genPulse()
	case 0x5:
		return (g.wavePT[g.genTriangle()>>1] << 4) & g.genPulse()
	case 0x6:
		return (g.wavePS[g.genSawtooth()] << 4) & g.genPulse()
	case 0x7:
		return (g.wavePST[g.genSawtooth()] << 4) & g.genPulse()
	case 0x8:
		return g.genNoise()
	case 0x9:
		return 0
	case 0xa:
		return 0
	case 0xb:
		return 0
	case 0xc:
		return 0
	case 0xd:
		return 0
	case 0xe:
		return 0
	case 0xf:
		return 0
	}
	return 0 // Shouldn't happen
}

func (g *Generator) reset() {
	g.accumulator = 0
	g.shiftRegister = 0x7ffff8
	g.freq = 0
	g.pulseWidth = 0

	g.test = false
	g.ringMod = false
	//g.sync = false

	g.syncTrigger = false
}
