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

func (v *VicII) Clock() {
	localCycle := v.cycle % v.dimensions.CyclesPerLine

	if v.clockPhase2 {
		if v.cycle == 0 {
			v.screen.Flip()
		}
		// ******** CLOCK PHASE 2 ********

		if localCycle >= PalFirstVisibleCycle && localCycle <= PalLastVisibleCycle &&
			v.rasterLine >= PalFirstVisibleLine && v.rasterLine <= PalLastVisibleLine {
			v.renderCycle()
		}

		// Perform c-access on bad lines
		if v.badLine && localCycle >= 15 && localCycle <= 54 {
			v.cAccess()
		}
		v.cpuBus.ClockPh2()

		// Perform s-access if sprites are enabled
		sprite, _ := getSpriteForCycle(localCycle)
		if sprite != 0x0ff {
			v.sAccess(sprite)
		}

		// Move to next cycle
		v.cycle++
		if v.cycle >= v.dimensions.Cycles {
			v.cycle = 0
			v.skipFrame = true // Skip unless DEN is set at line 48
		}
	} else {
		// ******** CLOCK PHASE 1 ********
		if v.cycle == 0 {
			v.vcBase = 0 // vic-ii.txt: 3.7.2.2
			v.badLine = false
			v.rasterLine = 0
		}

		switch localCycle {
		case 0:
			if v.cycle > 0 {
				v.rasterLine++
				v.rasterInterrupt()
			}
		case 1:
			if v.cycle == 0 {
				v.rasterLine = 0
				v.rasterInterrupt()
			}
		}

		// Handle sprite access if we're in the right part of the line.
		sprite, pCycle := getSpriteForCycle(localCycle)
		if sprite != 0xff {
			if pCycle {
				v.pAccess(sprite)
			} else {
				v.sAccess(sprite)
			}
		}

		// Are we three cycles away from a sprite with DMA enabled?
		// Then we need to lower the RDY line to stun the CPU while
		// we do phase 2 accesses to get the sprite data. The
		// RDY line will be raised at the end of the last S-access.
		if v.spriteWillNeedBus(localCycle) {
			v.bus.RDY.PullDown()
		}

		// Handle visible stuff
		contentBottom := v.dimensions.ContentBottom24Lines
		if v.line25 {
			contentBottom = v.dimensions.ContentBottom25Lines
		}
		if v.rasterLine >= 48 && v.rasterLine <= contentBottom {
			if v.rasterLine == 48 && v.enable {
				v.skipFrame = false
			}
			switch localCycle {
			case 14:
				v.vc = v.vcBase // vic-ii.txt: 3.7.2.2
				v.vmli = 0
				if v.badLine {
					v.rc = 0
				}
			case 15:
				// In the first phase of cycle 15, it is checked if the expansion flip flop
				// is set. If so, MCBASE is incremented by 2.
				for i := range v.sprites {
					if v.sprites[i].expand {
						v.sprites[i].mcBase += 2
					}
				}
			case 16:
				v.displayState = true

				// In the first phase of cycle 16, it is checked if the expansion flip flop
				// is set. If so, MCBASE is incremented by 1. After that, the VIC checks if
				// MCBASE is equal to 63 and turns of the DMA and the display of the sprite
				// if it is.
				for i := range v.sprites {
					s := &v.sprites[i]
					if s.expand {
						s.mcBase = (s.mcBase + 1) & 0x3f
						if s.mcBase == 0x3f {
							s.dma = false
						}
					}
				}
			case 55:
				// Bad lines always end here regardless of how they started
				if v.badLine {
					v.badLine = false
					v.bus.RDY.Release()
				}
				v.checkSpriteAndInit(true)
			case 56:
				v.checkSpriteAndInit(false)
			case 58:
				if v.rc == 0x07 {
					v.rc = 0 // vic-ii.txt: 3.7.2.5
					v.vcBase = v.vc
					v.displayState = false
				}

				/* 4. In the first phase of cycle 58, the MC of every sprite is loaded from
				   its belonging MCBASE (MCBASE->MC) and it is checked if the DMA for the
				   sprite is turned on and the Y coordinate of the sprite matches the lower
				   8 bits of RASTER. If this is the case, the display of the sprite is
				   turned on.
				*/
				for i := range v.sprites {
					s := &v.sprites[i]
					s.mc = s.mcBase
					if s.dma {
						s.displayEnabled = true
					}
				}
			case 59:
				if v.displayState {
					v.rc = (v.rc + 1) & 0x07
				}
			}

			// Misc stuff that's better handled outside the switch
			if localCycle >= 16 && localCycle < 56 {
				v.gAccess()
			}
			if localCycle >= 12 && localCycle <= 54 {
				if !v.skipFrame && v.rasterLine&0x07 == v.scrollY {
					if !v.badLine {
						// Flipping from normal to bad
						v.badLine = true
						v.bus.RDY.PullDown()
					}
				} else {
					// Flipping from bad to normal
					if v.badLine {
						v.badLine = false
						v.bus.RDY.Release()
					}
				}
			}
		}
		v.cpuBus.ClockPh1()
	}
	v.clockPhase2 = !v.clockPhase2
}

func (v *VicII) cAccess() {
	ch := v.bus.ReadByte(v.screenMemPtr | v.vc)
	col := v.colorRam.ReadByte(v.vc)
	if v.bitmapMode {
		v.cBuf[v.vmli] = uint16(ch)
	} else {
		v.cBuf[v.vmli] = uint16(col)<<8 | uint16(ch)
	}
}

func (v *VicII) gAccess() {
	basePtr := v.charSetPtr
	var addr uint16
	if v.bitmapMode {
		basePtr &= 0x2000
		addr = basePtr + v.rc + v.vc<<3
	} else {
		mask := uint16(0xff)
		if v.extendedClr {
			mask = 0x3f
		}
		addr = basePtr + ((v.cBuf[v.vmli]&mask)<<3 | v.rc) // TODO: Wild poking to the VIC can make vmli out of range
	}
	v.gBuf[v.vmli] = v.bus.ReadByte(addr)

	// Increment counters and make sure they stay within 6 and 10 bits boundary respectively
	v.vmli = (v.vmli + 1) & 0x003f
	v.vc = (v.vc + 1) & 0x03ff
}

func (v *VicII) pAccess(spriteIndex uint16) {
	addr := v.screenMemPtr | 0x03f8 | spriteIndex
	v.sprites[spriteIndex].pointer = uint16(v.bus.ReadByte(addr)) << 6
}

func (v *VicII) sAccess(spriteIndex uint16) {
	s := &v.sprites[spriteIndex]
	if !s.dma {
		s.gData[0] = 0
		s.gData[1] = 0
		s.gData[2] = 0
		return
	}
	s.gData[s.sIndex] = v.bus.ReadByte(s.pointer + uint16(s.mc))
	s.sIndex++
	s.mc++
	if s.sIndex > 2 {
		s.sIndex = 0
		v.bus.RDY.Release()
	}
}

func (v *VicII) rasterInterrupt() {
	if !v.irqRaster && v.irqRasterEnabled && v.rasterLine == v.rasterLineTrigger {
		v.irqRaster = true
		v.cpuBus.NotIRQ.PullDown()
	}
}

func (v *VicII) renderCycle() {
	var contentLeft uint16
	var contentRight uint16
	var contentTop uint16
	var contentBottom uint16
	//var scrollOffset uint16
	if v.col40 {
		contentLeft = v.dimensions.LeftBorderWidth40Cols
		contentRight = contentLeft + v.dimensions.ContentWidth40Cols
	} else {
		contentLeft = v.dimensions.LeftBorderWidth38Cols
		contentRight = contentLeft + v.dimensions.ContentWidth38Cols
	}
	if v.line25 {
		contentTop = v.dimensions.ContentTop25Lines
		contentBottom = v.dimensions.ContentBottom25Lines
		//scrollOffset = 3
	} else {
		contentTop = v.dimensions.ContentTop24Lines
		contentBottom = v.dimensions.ContentBottom24Lines
		//scrollOffset = 7
	}

	localCycle := (v.cycle % v.dimensions.CyclesPerLine) - v.dimensions.FirstVisibleCycle
	startPixel := localCycle << 3

	lastFgColor := uint8(0)

	for pixel := startPixel; pixel < startPixel+8; pixel++ {
		color := v.borderCol

		// 1. If the X coordinate reaches the right comparison value, the main border
		//    flip flop is set.
		// 2. If the Y coordinate reaches the bottom comparison value in cycle 63, the
		//    vertical border flip flop is set.
		// 3. If the Y coordinate reaches the top comparison value in cycle 63 and the
		//    DEN bit in register $d011 is set, the vertical border flip flop is
		//    reset.
		// 4. If the X coordinate reaches the left comparison value and the Y
		//    coordinate reaches the bottom one, the vertical border flip flop is set.
		// 5. If the X coordinate reaches the left comparison value and the Y
		//    coordinate reaches the top one and the DEN bit in register $d011 is set,
		//	  the vertical border flip flop is reset.
		// 6. If the X coordinate reaches the left comparison value and the vertical
		//    border flip flop is not set, the main flip flop is reset.
		if pixel == contentRight {
			v.hBorderFF = true
		}
		if v.rasterLine == contentBottom {
			v.vBorderFF = true
		}
		if v.rasterLine == contentTop {
			v.vBorderFF = false
		}
		if pixel == contentLeft && !v.vBorderFF && v.enable {
			v.hBorderFF = false
		}
		if !v.hBorderFF {
			// We're in the content area!

			// If we're at an odd pixel and multicolor mode is on, we just need to repeat the last pixel color
			if v.multiColor && pixel&1 == 1 {
				color = lastFgColor
			} else {
				localPixel := pixel - startPixel
				if localPixel&0x07 == v.scrollX {
					dataIdx := v.cycle%v.dimensions.CyclesPerLine - v.dimensions.FirstContentCycle
					v.cData = v.cBuf[dataIdx]
					v.sequencer = v.gBuf[dataIdx]
				}
				c, valid := v.generateSpriteColor(pixel, false) // Low priority sprites
				if valid {
					color = c
				} else {
					// TODO: Handle collisions
					if v.bitmapMode {
						color = v.generateBitmapColor()
					} else {
						color = v.generateTextColor()
					}
				}
				c, valid = v.generateSpriteColor(pixel, true) // High priority sprites
				if valid {
					// TODO: Handle collisions
					color = c
				}
				lastFgColor = color
			}
		}
		v.screen.SetPixel(pixel, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[color])
		v.sequencer <<= 1
	}
}

func (v *VicII) drawBorder(n uint16, segment []uint8) {
	for i := uint16(0); i < n; i++ {
		segment[i] = v.borderCol
		// v.screen.SetPixel(i+x, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[v.borderCol&0x0f])
	}
}

func (v *VicII) drawBackground(n uint16, segment []uint8) {
	for i := uint16(0); i < n; i++ {
		segment[i] = v.backgroundColors[0] & 0x0f
		// v.screen.SetPixel(i+x, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[v.backgroundColors[0]&0x0f])
	}
}

func (v *VicII) generateTextColor() uint8 {
	fgColor := uint8(v.cData>>8) & 0x0f
	bgIndex := 0
	if v.extendedClr {
		bgIndex = int(v.cData>>6) & 0x03
	}

	// Multicolor mode
	if v.multiColor {
		// N.Bus. This part must only be execute when x is even. Caller is responsible for that.
		// TODO: Handle illegal modes

		cIndex := v.sequencer & 0xc0 >> 6
		if cIndex == 0x03 {
			return fgColor & 0x0f
		} else {
			return v.backgroundColors[cIndex] & 0x0f
		}
	} else {
		if v.sequencer&0x80 != 0 {
			return fgColor
		} else {
			return v.backgroundColors[bgIndex]
		}
	}
}

func (v *VicII) generateBitmapColor() uint8 {
	bgColor := uint8(v.cData & 0x0f)
	fgColor := uint8(v.cData>>4) & 0x0f
	if v.multiColor {
		// N.Bus. This part must only be execute when x is even. Caller is responsible for that.
		cIndex := v.sequencer & 0xc0 >> 6
		switch cIndex {
		case 0:
			return v.backgroundColors[0]
		case 1:
			return uint8(v.cData>>4) & 0x0f
		case 2:
			return uint8(v.cData) & 0x0f
		case 3:
			return uint8(v.cData>>8) & 0x0f
		}
	} else {
		if v.sequencer&0x80 != 0 {
			return fgColor
		} else {
			return bgColor
		}
	}
	return 0 // This should never happen
}

func (v *VicII) generateSpriteColor(x uint16, highPriority bool) (uint8, bool) {
	color := uint8(0)
	valid := false
	for i := range v.sprites {
		s := &v.sprites[7-i] // Draw from lowest to highest to respect sprite priority
		if s.hasPriority != highPriority {
			continue
		}
		idx := x - s.x
		if idx <= 16 && idx&0x07 == 0 {
			s.sequencer = s.gData[idx>>3]
		}
		if s.enabled && s.sequencer&0x80 != 0 {
			color = s.color & 0x0f
			valid = true
		}
		s.sequencer <<= 1
	}
	return color, valid
}

// 1. The expansion flip flop is set as long as the bit in MxYE in register
//    $d017 corresponding to the sprite is cleared.
// 2. If the MxYE bit is set in the first phase of cycle 55, the expansion
//    flip flop is inverted.
// 3. In the first phases of cycle 55 and 56, the VIC checks for every sprite
//    if the corresponding MxE bit in register $d015 is set and the Y
//    coordinate of the sprite (odd registers $d001-$d00f) match the lower 8
//    bits of RASTER. If this is the case and the DMA for the sprite is still
//    off, the DMA is switched on, MCBASE is cleared, and if the MxYE bit is
//    set the expansion flip flip is reset.
func (v *VicII) checkSpriteAndInit(first bool) {
	// TODO: Handle expansion
	var start int
	var end int
	if first { // TODO: This may be backwards. Needs more research
		start = 3
		end = 7
	} else {
		start = 0
		end = 2
	}
	for i := start; i <= end; i++ {
		s := &v.sprites[i]
		if s.expandedY {
			s.expand = !s.expand
		}
		if s.enabled && s.y == uint8(v.rasterLine&0xff) && !s.dma {
			s.mcBase = 0
			s.dma = true
			s.expand = !s.expandedY // If not expanded, flip-flop will remain set forever
		}
	}
}

// Predict whether a sprite will need DMA in three cycles to allow
// the VIC to pull RDY low.
func (v *VicII) spriteWillNeedBus(localCycle uint16) bool {
	sIndex, pCycle := getSpriteForCycle((localCycle - 3) % 63)
	return sIndex != 0x0ff && pCycle && v.sprites[sIndex].dma
}

// Returns the sprite for this cycle (or 0xff if no sprite) and 'true' if this is
// a p-access cycle.
func getSpriteForCycle(localCycle uint16) (uint16, bool) {
	if localCycle == 0 {
		return 2, false // Wrap around from cycle 63 on previous line
	}
	if localCycle <= 10 {
		return (localCycle-1)/2 + 3, localCycle&1 == 1
	}
	if localCycle >= 58 {
		return (localCycle - 58) / 2, localCycle&1 == 0
	}
	return 0xff, false
}

func min(a uint16, b uint16) uint16 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a uint16, b uint16) uint16 {
	if a > b {
		return a
	} else {
		return b
	}
}
