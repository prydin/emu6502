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

import (
	"github.com/prydin/emu6502/core"
	"sync/atomic"
)

const (
	PRA     = 0x00 // Port A
	PRB     = 0x01 // Port B
	DDRA    = 0x02 // Data direction A
	DDRB    = 0x03 // Data direction B
	TALO    = 0x04 // Timer A Low
	TAHI    = 0x05 // Timer A High
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

const (
	internalClock     = 0x00
	externalClock     = 0x01
	chainedClock      = 0x02
	chainedGatedClock = 0x03
)

type CIA struct {
	bus       *core.Bus
	PortA     Port
	PortB     Port
	TimerA    Timer
	TimerB    Timer
	irqActive bool
}

type Port struct {
	data    uint8 // Input and output bits
	ddr     uint8 // Corresponding it is 0 for input, 1 for output
	PullUps uint8 // Corresponding bit set 1 simulates pullup-resistor
}

type Timer struct {
	pendingTicks int64
	counter      uint16
	latch        uint16
	running      bool
	source       int
	continuous   bool
	irqEnabled   bool
	irqOccurred  bool
	secondary    bool
	linkedTimer  *Timer
}

func (c *CIA) Init(bus *core.Bus) {
	c.bus = bus
	c.TimerA.linkedTimer = &c.TimerB
	c.TimerB.secondary = true
}

func (c *CIA) WriteByte(addr uint16, data uint8) {
	addr &= 0x0f
	switch addr {
	case PRA:
		c.PortA.internalWrite(data)
	case PRB:
		c.PortB.internalWrite(data)
	case DDRA:
		c.PortA.ddr = data
	case DDRB:
		c.PortB.ddr = data
	case TALO:
		c.TimerA.latch = c.TimerA.latch&0xff00 | uint16(data)
	case TAHI:
		c.TimerA.latch = c.TimerA.latch&0x00ff | uint16(data)<<8
	case TBLO:
		c.TimerB.latch = c.TimerB.latch&0xff00 | uint16(data)
	case TBHI:
		c.TimerB.latch = c.TimerB.latch&0x00ff | uint16(data)<<8
	case ICR:
		if data&0x80 != 0 {
			// Set bits
			if data&0x01 != 0 {
				c.TimerA.irqEnabled = true
			}
			if data&0x02 != 0 {
				c.TimerB.irqEnabled = true
			}
			// TODO: More flags
		} else {
			// Clear bits
			if data&0x01 != 0 {
				c.TimerA.irqEnabled = false
			}
			if data&0x02 != 0 {
				c.TimerB.irqEnabled = false
			}
			// TODO: More flags
		}
	case CRA:
		c.TimerA.setControlFlags(data)
	case CRB:
		c.TimerB.setControlFlags(data)
	}
}

func (c *CIA) ReadByte(addr uint16) uint8 {
	addr &= 0x0f
	switch addr {
	case PRA:
		return c.PortA.internalRead()
	case PRB:
		return c.PortB.internalRead()
	case DDRA:
		return c.PortA.ddr
	case DDRB:
		return c.PortB.ddr
	case TALO:
		return uint8(c.TimerA.counter & 0xff)
	case TAHI:
		return uint8(c.TimerA.counter >> 8)
	case TBLO:
		return uint8(c.TimerB.counter & 0xff)
	case TBHI:
		return uint8(c.TimerB.counter >> 8)
	case ICR:
		irqFlags := uint8(0)
		if c.TimerA.irqOccurred {
			irqFlags |= 0x01
			c.TimerA.irqOccurred = false
		}
		if c.TimerB.irqOccurred {
			irqFlags |= 0x01
			c.TimerB.irqOccurred = false
		}
		if irqFlags != 0 {
			irqFlags |= 0x80
		}
		// TODO: More interrupt sources
		if c.irqActive {
			c.bus.NotIRQ.Release()
		}
		return irqFlags
	}
	return 0xff // TODO: Is this correct?
}

func (c *CIA) Clock() {
	irqA := c.TimerA.irqOccurred
	irqB := c.TimerB.irqOccurred
	c.TimerA.Clock() // Preserve the order A -> B since B can tick A if it reaches zero
	c.TimerB.Clock()
	doIrq := (c.TimerA.irqOccurred && !irqA) || (c.TimerB.irqOccurred && !irqB) // Trigger on positive edge

	// TODO: Other stuff that might cause an interrupt

	if doIrq {
		c.irqActive = true
		c.bus.NotIRQ.PullDown()
	}
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

func (t *Timer) PulseCNT() {
	// TODO: Implement pin state instead
	if t.running && t.source == externalClock {
		atomic.AddInt64(&t.pendingTicks, 1)
	}
}

func (t *Timer) Clock() {
	if !t.running {
		return
	}
	if t.source == internalClock {
		t.tick()
	}
	n := atomic.SwapInt64(&t.pendingTicks, 0)
	for i := int64(0); i < n; i++ {
		t.tick()
	}
}

func (t *Timer) tick() {
	if t.counter == 0 {
		if t.continuous {
			t.counter = t.latch
			if t.linkedTimer != nil && t.linkedTimer.source == chainedClock && t.linkedTimer.running {
				t.linkedTimer.tick()
			}
		} else {
			t.running = false
		}
		if t.irqEnabled {
			t.irqOccurred = true
		}
	} else {
		t.counter--
	}
}

func (t *Timer) setControlFlags(flags uint8) {
	t.running = flags&0x01 != 0
	// TODO: Handle output modes
	t.continuous = flags&0x08 == 0
	if flags&0x10 != 0 {
		t.counter = t.latch
	}
	if t.secondary {
		t.source = int(flags & 0x60 >> 6)
	} else {
		t.source = int(flags & 0x20 >> 6)
	}

	// TODO: Handle serial port mode
	// TODO: Handle TOD frequency
}

func (t *Timer) getControlFlags() uint8 {
	flags := uint8(0)
	if t.running {
		flags |= 0x01
	}
	if t.continuous {
		flags |= 0x08
	}
	// TODO: Handle output modes
	// TODO: Handle force latch?
	if t.source == internalClock {
		flags |= 0x20
	}
	return flags
}
