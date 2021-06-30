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

import (
	"fmt"
	"github.com/prydin/emu6502/cia"
	"github.com/prydin/emu6502/core"
)

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

// Unused bits must be set to 1. This table provides a value we can or the result with
// when reading any register.
var unusedBitsMask = []uint8 {
	0x00 /* $D000 */, 0x00 /* $D001 */, 0x00 /* $D002 */, 0x00 /* $D003 */,
	0x00 /* $D004 */, 0x00 /* $D005 */, 0x00 /* $D006 */, 0x00 /* $D007 */,
	0x00 /* $D008 */, 0x00 /* $D009 */, 0x00 /* $D00A */, 0x00 /* $D00B */,
	0x00 /* $D00C */, 0x00 /* $D00D */, 0x00 /* $D00E */, 0x00 /* $D00F */,
	0x00 /* $D010 */, 0x00 /* $D011 */, 0x00 /* $D012 */, 0x00 /* $D013 */,
	0x00 /* $D014 */, 0x00 /* $D015 */, 0xc0 /* $D016 */, 0x00 /* $D017 */,
	0x01 /* $D018 */, 0x70 /* $D019 */, 0xf0 /* $D01A */, 0x00 /* $D01B */,
	0x00 /* $D01C */, 0x00 /* $D01D */, 0x00 /* $D01E */, 0x00 /* $D01F */,
	0xf0 /* $D020 */, 0xf0 /* $D021 */, 0xf0 /* $D022 */, 0xf0 /* $D023 */,
	0xf0 /* $D024 */, 0xf0 /* $D025 */, 0xf0 /* $D026 */, 0xf0 /* $D027 */,
	0xf0 /* $D028 */, 0xf0 /* $D029 */, 0xf0 /* $D02A */, 0xf0 /* $D02B */,
	0xf0 /* $D02C */, 0xf0 /* $D02D */, 0xf0 /* $D02E */, 0xff /* $D02F */,
	0xff /* $D030 */, 0xff /* $D031 */, 0xff /* $D032 */, 0xff /* $D033 */,
	0xff /* $D034 */, 0xff /* $D035 */, 0xff /* $D036 */, 0xff /* $D037 */,
	0xff /* $D038 */, 0xff /* $D039 */, 0xff /* $D03A */, 0xff /* $D03B */,
	0xff /* $D03C */, 0xff /* $D03D */, 0xff /* $D03E */, 0xff /* $D03F */,
}

type Sprite struct {
	enabled     bool
	x           uint16
	y           uint8
	color       uint8
	expandedX   bool
	expandedY   bool
	hasPriority bool
	multicolor  bool

	// Internal registers
	pointer        uint16
	gData          [3]uint8 // Current graphics data
	sequencer      uint8    // Shifted graphics data
	mc             uint8
	sIndex         int
	mcBase         uint8
	dma            bool  // Will do DMA before/after visible area
	expandYFF      bool  // Internal Y expansion flip flop
	repeatCnt      uint8 // Number of times to repeat a pixel due to multicolor and expansion
	mcFF           bool  // Internal multicolor flip flop
	displayEnabled bool
	lastColor      uint8 // Temporary color storage for multicolor sprites
	lastWasDrawn   bool  // Was the last pixel drawn or transparent?
}

type VicII struct {
	bus                    *core.Bus
	cpuBus                 *core.Bus
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
	spriteSpriteColl       uint8
	spriteDataColl         uint8
	cia2                   *cia.CIA // Pointer to CIA managing bank switch

	// Internal registers
	vc        uint16
	rc        uint16
	vcBase    uint16
	vmli      uint16
	skipFrame bool

	// Bad line flags
	badLine bool

	// Internal color and character buffers
	cBuf [40]uint16
	gBuf [40]uint8

	// Screen refresh state
	displayState bool
	cycle        uint16
	screen       Raster

	// Border flip flops
	vBorderFF bool
	hBorderFF bool

	// Sequencer internal registers
	sequencer uint8  // Sequencer shift register
	cData     uint16 // Current c-data
}

func (v *VicII) Init(bus *core.Bus, cpuBus *core.Bus, cia2 *cia.CIA, colorRam core.AddressSpace, screen Raster, dimensions ScreenDimensions) {
	v.screenMemPtr = 0x0400
	v.scrollY = 3
	v.scrollX = 0
	v.line25 = true
	v.col40 = true
	v.enable = true
	v.bitmapMode = false
	v.extendedClr = false
	v.multiColor = false
	v.borderCol = 4
	v.bus = bus
	v.screen = screen
	v.dimensions = dimensions
	v.colorRam = colorRam
	v.cpuBus = cpuBus
	v.hBorderFF = true
	v.vBorderFF = true
	v.cia2 = cia2
	for i := range v.sprites {
		v.sprites[i].expandYFF = true
	}
}

func (v *VicII) ReadByte(addr uint16) uint8 {
	fmt.Printf("READ $D0%02X at cycle %d of current_line $%04X\n",
		addr, v.cycle % 63, v.rasterLine)
	addr &= 0x003f
	data := uint8(0xff)
	switch {
	case addr <= REG_M7X && addr&0x01 == 0:
		data =  uint8(v.sprites[addr>>1].x & 0xff)
	case addr <= REG_M7Y && addr&0x01 == 1:
		data =  v.sprites[addr>>1].y & 0xff
	case addr == REG_MX_HIGH:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			data |= uint8((v.sprites[i].x & 0x100) >> 1)
		}
	case addr >= REG_BG0 && addr <= REG_BG3:
		data =  v.backgroundColors[addr-REG_BG0] | 0xf0
	case addr >= REG_SPRITE_COL0 && addr <= REG_SPRITE_COL7:
		data =  v.sprites[addr-REG_SPRITE_COL0].color | 0xf0
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
		data =  uint8(r)
	case REG_RASTER_CNT:
		data =  uint8(v.rasterLine & 0xff)
	case REG_LPX:
		data =  0
	case REG_LPY:
		data =  0
	case REG_SPRITE_ENABLE:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			if v.sprites[i].enabled {
				data |= 0x80
			}
		}
		data =  data
	case REG_CTRL2:
		r := uint8(v.scrollX)
		if v.col40 {
			r |= CR2_40COLS
		}
		if v.multiColor {
			r |= CR2_MULTICOLOR
		}
		data =  r | 0xc0
	case REG_MX_EXP:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			if v.sprites[i].expandedX {
				data |= 0x80
			}
		}
		data = data
	case REG_MEMPTR:
		data = uint8(v.charSetPtr >> 10)
		data |= uint8(v.screenMemPtr >> 6)
	case REG_IRQ:
		data = uint8(0)
		if v.irqRaster {
			data |= IRQ_RASTER
		}
		if v.irqSpriteBg {
			data |= IRQ_SPRITE_BG_COLL
		}
		if v.irqSpriteSprite {
			data |= IRQ_SPRITE_SPRITE_COLL
		}
		if v.irqLp {
			data |= IRQ_LP
		}
		if data != 0 {
			data |= IRQ_HAPPENED
		}
	case REG_IRQ_ENABLE:
		data = uint8(0)
		if v.irqRasterEnabled {
			data |= IRQ_RASTER
		}
		if v.irqSpriteBgEnabled {
			data |= IRQ_SPRITE_BG_COLL
		}
		if v.irqSpriteSpriteEnabled {
			data |= IRQ_SPRITE_SPRITE_COLL
		}
		if v.irqLpEnabled {
			data |= IRQ_LP
		}
		if v.irqEnabled {
			data |= IRQ_ENABLED
		}
	case REG_SPRITE_PRIO:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			if v.sprites[i].hasPriority {
				data |= 0x80
			}
		}
	case REG_SPRITE_MULTICLR:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			if v.sprites[i].multicolor {
				data |= 0x80
			}
		}
	case REG_MY_EXP:
		data = uint8(0)
		for i := range v.sprites {
			data >>= 1
			if v.sprites[i].expandedY {
				data |= 0x80
			}
		}
	case REG_SPRITE_COLL:
		data = v.spriteSpriteColl
		v.spriteSpriteColl = 0
	case REG_DATA_COLL:
		data = v.spriteDataColl
		v.spriteDataColl = 0
	case REG_BORDER:
		data =  v.borderCol
	case REG_SPRITE_MC0:
		data =  v.spriteMultiClr0
	case REG_SPRITE_MC1:
		data =  v.spriteMultiClr1
	}
	return data | unusedBitsMask[addr]
}

func (v *VicII) WriteByte(addr uint16, data uint8) {
	fmt.Printf("WRITE $D0%02X at cycle %d of current_line $%04X\n",
		addr, v.cycle % 63, v.rasterLine)
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
			v.sprites[i].x = (v.sprites[i].x & 0x00ff) | uint16(data&0x01)<<8
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
			v.charSetPtr = 0x1000 // TODO: Temporary hack. Fix when banking is implemented. Not valid for bitmap mode!
		} else {
			v.charSetPtr = uint16(data&0x0e) << 10
		}
		v.screenMemPtr = uint16(data&0xf0) << 6
	case REG_IRQ:
		newRaster := data&IRQ_RASTER != 0
		if newRaster && v.irqRaster {
			v.cpuBus.NotIRQ.Release()
			v.irqRaster = false
		}
		newSpritqBg := v.irqSpriteBg && data&IRQ_SPRITE_BG_COLL != 0
		if newSpritqBg && v.irqSpriteBg {
			v.cpuBus.NotIRQ.Release()
			v.irqSpriteBg = false
		}
		newSpriteSprite := data&IRQ_SPRITE_SPRITE_COLL != 0
		if newSpriteSprite && v.irqSpriteSprite {
			v.cpuBus.NotIRQ.Release()
			v.irqSpriteBg = false
		}
		newLp := data&IRQ_LP != 0
		if newLp && v.irqLp {
			v.cpuBus.NotIRQ.Release()
			v.irqLp = false
		}
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
			v.sprites[i].hasPriority = data&0x01 == 0
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
	// We don't handle writes to the collision register since it's assumed that they're read-only.
	// TODO: Check this!
}

func (v *VicII) GetCycle() uint16 {
	return v.cycle
}

func (v *VicII) IsVSynch() bool {
	return v.cycle == 0 && !v.clockPhase2
}

func (v *VicII) IsWriteable() bool {
	return true
}
