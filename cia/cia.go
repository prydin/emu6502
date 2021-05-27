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

package cia

import "github.com/prydin/emu6502/core"

const (
	PRA     = 0x00 // Port A
	PRB     = 0x01 // Port B
	DDRA    = 0x02 // Data direction A
	DDRB    = 0x03 // Data direction B
	TALO    = 0x04 // Timer A Low
	TAHO    = 0x05 // Timer A High
	TBLO    = 0x06 // Timer B Low
	TBHI    = 0x07 // Timer B High
	TOD10TH = 0x08 // Time of day 10th seconds
	TODSEC  = 0x09 // Time of day seconds
	TODMIN  = 0x0a // Time of day minutes
	TODHR   = 0x0b // Time of day hours
	SDR     = 0x0c // Serial shift register
	ICR     = 0x0d // Interrupt control
	CRA     = 0x0e // Timer A control
	CRB     = 0x0f // Timer B control
)

type CIA struct {
	bus    *core.Bus
	PortA  Port
	PortB  Port
	TimerA Timer
	TimerB Timer
}

type Port struct {
	data    uint8 // Input and output bits
	ddr     uint8 // Corresponding it is 0 for input, 1 for output
	PullUps uint8 // Corresponding bit set 1 simulates pullup-resistor
}

func (p *Port) internalRead() uint8 {
	return p.data & ^p.ddr | p.PullUps&p.ddr
}

func (p *Port) internalWrite(data uint8) {
	p.data = p.data & ^p.ddr | data&p.ddr
}

func (p *Port) ReadOutputs() uint8 {
	return p.data&p.ddr | p.PullUps & ^p.ddr
}

func (p *Port) SetInputs(data uint8) {
	p.data = p.data&p.ddr | data & ^p.ddr
}

type Timer struct {
	counter uint16
	latch   uint16
}
