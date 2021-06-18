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

const (
	ATTACK = iota
	DECAY_SUSTAIN
	RELEASE
)

var rateCounterPeriod = []reg16{
	9,     //   2ms*1.0MHz/256 =     7.81
	32,    //   8ms*1.0MHz/256 =    31.25
	63,    //  16ms*1.0MHz/256 =    62.50
	95,    //  24ms*1.0MHz/256 =    93.75
	149,   //  38ms*1.0MHz/256 =   148.44
	220,   //  56ms*1.0MHz/256 =   218.75
	267,   //  68ms*1.0MHz/256 =   265.63
	313,   //  80ms*1.0MHz/256 =   312.50
	392,   // 100ms*1.0MHz/256 =   390.63
	977,   // 250ms*1.0MHz/256 =   976.56
	1954,  // 500ms*1.0MHz/256 =  1953.13
	3126,  // 800ms*1.0MHz/256 =  3125.00
	3907,  //   1 s*1.0MHz/256 =  3906.25
	11720, //   3 s*1.0MHz/256 = 11718.75
	19532, //   5 s*1.0MHz/256 = 19531.25
	31251, //   8 s*1.0MHz/256 = 31250.00
}

var sustainLevel = []reg8{
	0x00,
	0x11,
	0x22,
	0x33,
	0x44,
	0x55,
	0x66,
	0x77,
	0x88,
	0x99,
	0xaa,
	0xbb,
	0xcc,
	0xdd,
	0xee,
	0xff,
}

type EnvelopeGenerator struct {
	rateCounter              reg16
	ratePeriod               reg16
	exponentialCounter       reg8
	exponentialCounterPeriod reg8
	envelopeCounter          reg8
	holdZero                 bool

	attack  reg4
	decay   reg4
	sustain reg4
	release reg4
	state   int

	gate bool
}

func (e *EnvelopeGenerator) ReadOutput() reg8 {
	return e.envelopeCounter
}

func (e *EnvelopeGenerator) setControl(control reg8) {
	gateNext := control&0x01 != 0

	// The rate counter is never reset, thus there will be a delay before the
	// envelope counter starts counting up (attack) or down (release).

	// Gate bit on: Start attack, decay, sustain.
	if !e.gate && gateNext {
		e.state = ATTACK
		e.ratePeriod = rateCounterPeriod[e.attack]

		// Switching to attack state unlocks the zero freeze.
		e.holdZero = false
	} else if e.gate && !gateNext {
		e.state = RELEASE
		e.ratePeriod = rateCounterPeriod[e.release]
	}
	e.gate = gateNext
}

func (e *EnvelopeGenerator) setAttackDecay(attackDecay uint8) {
	e.attack = reg4((attackDecay >> 4) & 0x0f)
	e.decay = reg4(attackDecay & 0x0f)
	if e.state == ATTACK {
		e.ratePeriod = rateCounterPeriod[e.attack]
	} else if e.state == DECAY_SUSTAIN {
		e.ratePeriod = rateCounterPeriod[e.decay]
	}
}

func (e *EnvelopeGenerator) setSustainRelease(sustainRelease uint8) {
	e.sustain = reg4((sustainRelease >> 4) & 0x0f)
	e.release = reg4(sustainRelease & 0x0f)
	if e.state == RELEASE {
		e.ratePeriod = rateCounterPeriod[e.release]
	}
}

func (e *EnvelopeGenerator) Clock() {
	// Check for ADSR delay bug.
	// If the rate counter comparison value is set below the current value of the
	// rate counter, the counter will continue counting up until it wraps around
	// to zero at 2^15 = 0x8000, and then count rate_period - 1 before the
	// envelope can finally be stepped.
	// This has been verified by sampling ENV3.
	//
	e.rateCounter++
	if +e.rateCounter&0x8000 != 0 {
		e.rateCounter++
		e.rateCounter &= 0x7fff
	}

	if e.rateCounter != e.ratePeriod {
		return
	}

	e.rateCounter = 0

	// The first envelope step in the attack state also resets the exponential
	// counter. This has been verified by sampling ENV3.
	//
	e.exponentialCounter++
	if e.state == ATTACK || e.exponentialCounter == e.exponentialCounterPeriod {
		e.exponentialCounter = 0

		// Check whether the envelope counter is frozen at zero.
		if e.holdZero {
			return
		}

		switch e.state {
		case ATTACK:
			// The envelope counter can flip from 0xff to 0x00 by changing state to
			// release, then to attack. The envelope counter is then frozen at
			// zero; to unlock this situation the state must be changed to release,
			// then to attack. This has been verified by sampling ENV3.
			//
			e.envelopeCounter++
			e.envelopeCounter &= 0xff
			if e.envelopeCounter == 0xff {
				e.state = DECAY_SUSTAIN
				e.ratePeriod = rateCounterPeriod[e.decay]
			}
		case DECAY_SUSTAIN:
			if e.envelopeCounter != sustainLevel[e.sustain] {
				e.envelopeCounter--
			}
		case RELEASE:
			// The envelope counter can flip from 0x00 to 0xff by changing state to
			// attack, then to release. The envelope counter will then continue
			// counting down in the release state.
			// This has been verified by sampling ENV3.
			// NB! The operation below requires two's complement integer.
			//
			e.envelopeCounter--
			e.envelopeCounter &= 0xff
		}

		// Check for change of exponential counter period.
		switch e.envelopeCounter {
		case 0xff:
			e.exponentialCounterPeriod = 1
		case 0x5d:
			e.exponentialCounterPeriod = 2
		case 0x36:
			e.exponentialCounterPeriod = 4
		case 0x1a:
			e.exponentialCounterPeriod = 8
		case 0x0e:
			e.exponentialCounterPeriod = 16
		case 0x06:
			e.exponentialCounterPeriod = 30
		case 0x00:
			e.exponentialCounterPeriod = 1

			// When the envelope counter is changed to zero, it is frozen at zero.
			// This has been verified by sampling ENV3.
			e.holdZero = true
		}
	}
}

func (e *EnvelopeGenerator) reset() {
	e.envelopeCounter = 0
	e.attack = 0
	e.decay = 0
	e.sustain = 0
	e.release = 0
	e.gate = false
	e.rateCounter = 0
	e.exponentialCounter = 0
	e.exponentialCounterPeriod = 1
	e.state = RELEASE
	e.ratePeriod = rateCounterPeriod[e.release]
	e.holdZero = true
}
