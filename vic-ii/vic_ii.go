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

package vic_ii

import "github.com/prydin/emu6502/core"

const (
	REG_M0X             = 0x00
	REG_M0Y             = 0x01
	REG_M1X             = 0x02
	REG_M1Y             = 0x03
	REG_M2X             = 0x04
	REG_M2Y             = 0x05
	REG_M3X             = 0x06
	REG_M3Y             = 0x07
	REG_M4X             = 0x08
	REG_M4Y             = 0x09
	REG_M5X             = 0x0a
	REG_M5Y             = 0x0b
	REG_M6X             = 0x0c
	REG_M6Y             = 0x0d
	REG_M7X             = 0x0e
	REG_M7Y             = 0x0f
	REG_MX_HIGH         = 0x10
	REG_CTRL1           = 0x11
	REG_RASTER_CNT      = 0x12
	REG_LPX             = 0x13
	REG_LPY             = 0x14
	REG_SPRITE_ENABLE   = 0x15
	REG_CTRL2           = 0x16
	REG_MY_EXP          = 0x17
	REG_MEMPTR          = 0x18
	REG_IRQ             = 0x19
	REG_IRQ_ENABLE      = 0x1a
	REG_SPRITE_PRIO     = 0x1b
	REG_SPRITE_MULTICLR = 0x1c
	REG_MX_EXP          = 0x1d
	REG_SPRITE_COLL     = 0x1e
	REG_DATA_COLL       = 0x1f
	REG_BORDER          = 0x20
	REG_BG0             = 0x21
	REG_BG1             = 0x22
	REG_BG2             = 0x23
	REG_BG3             = 0x24
	REG_SPRITE_MC0      = 0x25
	REG_SPRITE_MC1      = 0x26
	REG_SPRITE_COL0     = 0x27
	REG_SPRITE_COL1     = 0x28
	REG_SPRITE_COL2     = 0x29
	REG_SPRITE_COL3     = 0x2a
	REG_SPRITE_COL4     = 0x2b
	REG_SPRITE_COL5     = 0x2c
	REG_SPRITE_COL6     = 0x2d
	REG_SPRITE_COL7     = 0x2e
)

// Control register 1
const (
	CR1_RASTERHI = 0x80
	CR1_EXTCLR   = 0x40
	CR1_BITMAP   = 0x20
	CR1_ENABLE   = 0x10
	CR1_25LINES  = 0x08
)

// Control register 2
const (
	CR2_MULTICOLOR = 0x10
	CR2_40COLS     = 0x08
)

// Interrupt control bits
const (
	IRQ_RASTER             = 0x01
	IRQ_SPRITE_BG_COLL     = 0x02
	IRQ_SPRITE_SPRITE_COLL = 0x04
	IRQ_LP                 = 0x08
	IRQ_HAPPENED           = 0x80
	IRQ_ENABLED            = 0x80
)

// Screen refresh states
const (
	STATE_VBLANKT = iota
	STATE_VBLANKB
	STATE_HBLANKL
	STATE_HBLANKR
	STATE_BORDER
	STATE_CONTENT
)

type Sprite struct {
	enabled     bool
	x           uint16
	y           uint8
	color       uint8
	expandedX   bool
	expandedY   bool
	hasPriority bool
	multicolor  bool
}

type VicII struct {
	bus                    *core.Bus
	clockSink              *core.Bus
	colorRam               core.AddressSpace
	dimensions             ScreenDimensions
	clockPhase2            bool
	borderCol              uint8 // Border color
	sprites                [8]Sprite
	rasterLineTrigger      uint16
	scrollX                uint16
	scrollY                uint16
	line25                 bool
	col40                  bool
	enable                 bool
	bitmapMode             bool
	extendedClr            bool
	multiColor             bool
	irqRaster              bool
	irqSpriteBg            bool
	irqSpriteSprite        bool
	irqLp                  bool
	irqEnabled             bool
	irqRasterEnabled       bool
	irqSpriteBgEnabled     bool
	irqSpriteSpriteEnabled bool
	irqLpEnabled           bool
	backgroundColors       [4]uint8
	spriteMultiClr0        uint8
	spriteMultiClr1        uint8
	charSetPtr             uint16
	screenMemPtr           uint16
	rasterLine             uint16

	// Internal registers
	vc        uint16
	rc        uint16
	vcBase    uint16
	vmli      uint16
	skipFrame bool

	// Bad line flags
	badLine bool

	// Internal color and character buffers
	cBuf       [40]uint16
	gBuf       [40]uint8
	charBufPtr uint16

	// Screen refresh state
	displayState bool
	cycle        uint16
	screen       Raster

	// Border flip-flops
	vBorderFf bool
	hBorderFf bool
}

func (v *VicII) Init(bus *core.Bus, clockSink *core.Bus, colorRam core.AddressSpace, screen Raster, dimensions ScreenDimensions) {
	v.screenMemPtr = 0x0400
	v.scrollY = 3
	v.scrollX = 0
	v.line25 = true
	v.col40 = true
	v.enable = true
	v.bitmapMode = false
	v.extendedClr = false
	v.multiColor = false
	v.hBorderFf = true
	v.vBorderFf = true
	v.borderCol = 4
	v.bus = bus
	v.screen = screen
	v.dimensions = dimensions
	v.colorRam = colorRam
	v.clockSink = clockSink
}

func (v *VicII) ReadByte(addr uint16) uint8 {
	addr &= 0x003f
	if addr >= 0x002f {
		return 0xff // Read from unassigned address
	}
	switch {
	case addr <= REG_M7X && addr&0x01 == 0:
		return uint8(v.sprites[addr>>1].x & 0xff)
	case addr <= REG_M7Y && addr&0x01 == 1:
		return v.sprites[addr>>1].y & 0xff
	case addr == REG_MX_HIGH:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			r |= uint8((v.sprites[i].x & 0x100) >> 1)
		}
		return r
	case addr >= REG_BG0 && addr <= REG_BG3:
		return v.backgroundColors[addr-REG_BG0] | 0xf0
	case addr >= REG_SPRITE_COL0 && addr <= REG_SPRITE_COL7:
		return v.sprites[addr-REG_SPRITE_COL0].color | 0xf0
	}
	switch addr {
	case REG_CTRL1:
		r := (v.rasterLine & 0x100) >> 1
		if v.extendedClr {
			r |= CR1_EXTCLR
		}
		if v.bitmapMode {
			r |= CR1_BITMAP
		}
		if v.enable {
			r |= CR1_ENABLE
		}
		if v.line25 {
			r |= CR1_25LINES
		}
		r |= v.scrollY & 0x07
		return uint8(r)
	case REG_RASTER_CNT:
		return uint8(v.rasterLine & 0xff)
	case REG_LPX:
		return 0
	case REG_LPY:
		return 0
	case REG_SPRITE_ENABLE:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			if v.sprites[i].enabled {
				r |= 0x80
			}
		}
		return r
	case REG_CTRL2:
		r := uint8(v.scrollX)
		if v.col40 {
			r |= CR2_40COLS
		}
		if v.multiColor {
			r |= CR2_MULTICOLOR
		}
		return r | 0xc0
	case REG_MX_EXP:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			if v.sprites[i].expandedX {
				r |= 0x80
			}
		}
		return r
	case REG_MEMPTR:
		r := uint8(v.charSetPtr >> 10)
		r |= uint8(v.screenMemPtr >> 6)
		return r | 0x01
	case REG_IRQ:
		r := uint8(0)
		if v.irqRaster {
			r |= IRQ_RASTER
		}
		if v.irqSpriteBg {
			r |= IRQ_SPRITE_BG_COLL
		}
		if v.irqSpriteSprite {
			r |= IRQ_SPRITE_SPRITE_COLL
		}
		if v.irqLp {
			r |= IRQ_LP
		}
		if r != 0 {
			r |= IRQ_HAPPENED
		}
		return r | 0x70
	case REG_IRQ_ENABLE:
		r := uint8(0)
		if v.irqRasterEnabled {
			r |= IRQ_RASTER
		}
		if v.irqSpriteBgEnabled {
			r |= IRQ_SPRITE_BG_COLL
		}
		if v.irqSpriteSpriteEnabled {
			r |= IRQ_SPRITE_SPRITE_COLL
		}
		if v.irqLpEnabled {
			r |= IRQ_LP
		}
		if v.irqEnabled {
			r |= IRQ_ENABLED
		}
		return r | 0xf0
	case REG_SPRITE_PRIO:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			if v.sprites[i].hasPriority {
				r |= 0x80
			}
		}
		return r
	case REG_SPRITE_MULTICLR:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			if v.sprites[i].multicolor {
				r |= 0x80
			}
		}
		return r
	case REG_MY_EXP:
		r := uint8(0)
		for i := range v.sprites {
			r >>= 1
			if v.sprites[i].expandedY {
				r |= 0x80
			}
		}
		return r
		// TODO: Collisions
	case REG_BORDER:
		return v.borderCol
	case REG_SPRITE_MC0:
		return v.spriteMultiClr0
	case REG_SPRITE_MC1:
		return v.spriteMultiClr1
	}
	return 0xff
}

func (v *VicII) WriteByte(addr uint16, data uint8) {
	addr &= 0x003f
	if addr >= 0x002f {
		return
	}
	switch {
	case addr <= REG_M7X && addr&0x01 == 0:
		v.sprites[addr>>1].x = (v.sprites[addr>>1].x & 0xff00) | uint16(data)
		return
	case addr <= REG_M7Y && addr&0x01 == 1:
		v.sprites[addr>>1].y = data
	case addr == REG_MX_HIGH:
		for i := range v.sprites {
			v.sprites[i].x = (v.sprites[i].x & 0x00ff) | uint16(data)<<8
			data >>= 1
		}
	case addr >= REG_BG0 && addr <= REG_BG3:
		v.backgroundColors[addr-REG_BG0] = data
	case addr >= REG_SPRITE_COL0 && addr <= REG_SPRITE_COL7:
		v.sprites[addr-REG_SPRITE_COL0].color = data
	}
	switch addr {
	case REG_CTRL1:
		v.rasterLineTrigger = (v.rasterLineTrigger & 0x00ff) | uint16(data&CR1_RASTERHI)<<1
		v.extendedClr = data&CR1_EXTCLR != 0
		v.bitmapMode = data&CR1_BITMAP != 0
		v.enable = data&CR1_ENABLE != 0
		v.line25 = data&CR1_25LINES != 0
		v.scrollY = uint16(data & 0x07)
	case REG_RASTER_CNT:
		v.rasterLineTrigger = (v.rasterLineTrigger & 0xff00) | uint16(data)
	case REG_LPX:
		return
	case REG_LPY:
		return
	case REG_SPRITE_ENABLE:
		for i := range v.sprites {
			v.sprites[i].enabled = data&0x01 != 0
			data >>= 1
		}
	case REG_CTRL2:
		v.scrollX = uint16(data & 0x07)
		v.col40 = data&CR2_40COLS != 0
		v.multiColor = data&CR2_MULTICOLOR != 0
	case REG_MX_EXP:
		for i := range v.sprites {
			v.sprites[i].expandedX = data&0x01 != 0
			data >>= 1
		}
	case REG_MEMPTR:
		if data&0x0e>>1 == 2 || data&0x0e>>1 == 3 {
			v.charSetPtr = 0x1000 // TODO: Temporary hack. Fix when banking is implemented
		} else {
			v.charSetPtr = uint16(data&0x0e) << 10
		}
		v.screenMemPtr = uint16(data&0xf0) << 6
	case REG_IRQ:
		v.irqRaster = data&IRQ_RASTER != 0
		v.irqSpriteBg = data&IRQ_SPRITE_BG_COLL != 0
		v.irqSpriteSprite = data&IRQ_SPRITE_SPRITE_COLL != 0
		v.irqLp = data&IRQ_LP != 0
	case REG_IRQ_ENABLE:
		v.irqRasterEnabled = data&IRQ_RASTER != 0
		v.irqSpriteBgEnabled = data&IRQ_SPRITE_BG_COLL != 0
		v.irqSpriteSpriteEnabled = data&IRQ_SPRITE_SPRITE_COLL != 0
		v.irqLpEnabled = data&IRQ_LP != 0
		v.irqEnabled = data&IRQ_ENABLED != 0
	case REG_BORDER:
		v.borderCol = data
	case REG_SPRITE_PRIO:
		for i := range v.sprites {
			v.sprites[i].hasPriority = data&0x01 != 0
			data >>= 1
		}
	case REG_SPRITE_MULTICLR:
		for i := range v.sprites {
			v.sprites[i].multicolor = data&0x01 != 0
			data >>= 1
		}
	case REG_MY_EXP:
		for i := range v.sprites {
			v.sprites[i].expandedY = data&0x01 != 0
			data >>= 1
		}
	case REG_SPRITE_MC0:
		v.spriteMultiClr0 = data
	case REG_SPRITE_MC1:
		v.spriteMultiClr1 = data
	}
}
