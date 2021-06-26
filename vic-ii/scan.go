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
/*	if v.sprites[1].dma {
		println(v.cycle%63, v.sprites[1].mcBase, v.sprites[1].mc)
	} */
	// IMPORTANT NOTICE: Some VIC-II literature, including Christian Bauer's famous
	// text file use one-based numbering of cycles within a line. In this code, we
	// use zero-based values. Thus, everything will be off by one! Pay attention!
	if v.clockPhase2 {
		v.cpuBus.ClockPh2() // TODO: Not sure about the order here? In theory it's simulateneous
		v.phi2()

	} else {
		v.cpuBus.ClockPh1() // TODO: Not sure about the order here? In theory it's simulateneous
		v.phi1()

	}
	v.clockPhase2 = !v.clockPhase2
}

func (v *VicII) phi1() {
	if v.rasterLine == 48 && v.enable {
		v.skipFrame = false
	}
	localCycle := v.cycle % v.dimensions.CyclesPerLine
	switch localCycle { // TODO: This is hardcoded for PAL
	case 0:
		v.initLine()
		v.pAccess(3)
	case 1:
		v.sAccess(3)
	case 2:
		v.pAccess(4)
	case 3:
		v.sAccess(4)
	case 4:
		v.pAccess(5)
	case 5:
		v.sAccess(5)
	case 6:
		v.pAccess(6)
	case 7:
		v.sAccess(6)
	case 8:
		v.pAccess(7)
	case 9:
		v.sAccess(7)
	case 10:
		// Nothing
	case 11:
		v.checkForBadLine()
	case 12:
		v.checkForBadLine()
	case 13:
		// In the first phase of cycle 14 of each line, VC is loaded from VCBASE
		// (VCBASE->VC) and VMLI is cleared. If there is a Bad Line Condition in
		// this phase, RC is also reset to zero.
		v.checkForBadLine()
		v.initVc()
	case 14:
		v.checkForBadLine()
		v.handleSpriteYExpStep()
	case 15:
		v.checkForBadLine()
		if v.badLine {
			v.displayState = true
		}
		v.initMc()
		v.gAccess()
	case 54:
		// Bad lines always end here regardless of how they started
		v.resetBadLine()
		v.checkSpriteAndInit(true)
		v.gAccess()
	case 55:
		v.checkSpriteAndInit(false)
		fallthrough
	case 56:
		v.idleGap()
	case 57:
		v.prepareVc()
		v.resetSpriteCounters()
		v.pAccess(0)
	case 58:
		v.sAccess(0)
	case 59:
		v.pAccess(1)
	case 60:
		v.sAccess(1)
	case 61:
		v.pAccess(2)
	case 62:
		v.sAccess(2)
	default: // 16-53
		v.gAccess()
	}
}

func (v *VicII) phi2() {
	if v.cycle == 0 {
		v.screen.Flip()
	}
	localCycle := v.cycle % v.dimensions.CyclesPerLine

	// The negative edge of IRQ on a raster interrupt has been used to define the
	// beginning of a line (this is also the moment in which the RASTER register
	// is incremented). Raster line 0 is, however, an exception: In this line, IRQ
	// and incrementing (and resetting) of RASTER are performed one cycle later
	// than in the other lines. But for simplicity we assume equal line lengths
	// and define the beginning of raster line 0 to be one cycle before the
	// occurrence of the IRQ.
	if localCycle == 0 && v.cycle != 0 {
		v.rasterLine++
		v.rasterInterrupt()
	} else if localCycle == 1 && v.cycle == 1 {
		// First line is one clock cycle delayed in triggering interrupt
		v.rasterLine = 0
		v.rasterInterrupt()
	}

	if localCycle >= PalFirstVisibleCycle && localCycle <= PalLastVisibleCycle &&
		v.rasterLine >= PalFirstVisibleLine && v.rasterLine <= PalLastVisibleLine {
		v.renderCycle()
	}
	switch localCycle {
	case 0:
		fallthrough
	case 1:
		v.sAccess(3)
	case 2:
		fallthrough
	case 3:
		v.sAccess(4)
	case 4:
		fallthrough
	case 5:
		v.sAccess(5)
	case 6:
		fallthrough
	case 7:
		v.sAccess(6)
	case 8:
		fallthrough
	case 9:
		v.sAccess(7)
	case 10:
		// Nothing
	case 11:
		// Nothing
	case 12:
		// Nothing
	case 13:
		// Nothing
	case 54:
		// Nothing
	case 55:
		// Nothing
	case 56:
		// Nothing
	case 57:
		fallthrough
	case 58:
		v.sAccess(0)
	case 59:
		fallthrough
	case 60:
		v.sAccess(1)
	case 61:
		fallthrough
	case 62:
		v.sAccess(2)
	default: // 14-54
		if v.badLine {
			v.cAccess()
		}
	}

	// Let the CPU do it's thing!

	// Move to next cycle
	v.cycle++
	if v.cycle >= v.dimensions.Cycles {
		v.cycle = 0
		v.skipFrame = true // Skip unless DEN is set at line 48
	}
}

func (v *VicII) initLine() {
	if v.cycle == 0 {
		v.vcBase = 0 // vic-ii.txt: 3.7.2.2
		v.badLine = false
		v.rasterLine = 0
	}
}

func (v *VicII) initVc() {
	v.vc = v.vcBase
	v.vmli = 0
	if v.badLine {
		v.rc = 0
	}
}

func (v *VicII) prepareVc() {
	// In the first phase of cycle 58, the VIC checks if RC=7. If so, the video
	// logic goes to idle state and VCBASE is loaded from VC (VC->VCBASE). If
	// the video logic is in display state afterwards (this is always the case
	// if there is a Bad Line Condition), RC is incremented.
	if v.rc == 0x07 {
		v.rc = 0
		v.vcBase = v.vc
		v.displayState = false
		//fmt.Printf("vc=%d, vcBase=%d, cycle=%d, vc/40=%d, cycle/63=%d\n", v.vc, v.vcBase, v.cycle, v.vc/40, (v.cycle-3711)/63/8)
	} else {
		if v.displayState {
			v.rc = (v.rc + 1) & 0x07
		}
	}
}

func (v *VicII) resetBadLine() {
	if v.badLine {
		v.badLine = false
		v.bus.RDY.Release()
	}
}

func (v *VicII) initMc() {
	// In the first phase of cycle 16, it is checked if the expansion flip flop
	// is set. If so, MCBASE is incremented by 1. After that, the VIC checks if
	// MCBASE is equal to 63 and turns of the DMA and the display of the sprite
	// if it is.
	for i := range v.sprites {
		s := &v.sprites[i]
		if s.expandYFF {
			s.mcBase = (s.mcBase + 1) & 0x3f
			if s.mcBase == 63 {
				s.dma = false
			}
		}
	}
}

func (v *VicII) checkForBadLine() {
	// If there is a Bad Line Condition in cycles 12-54, BA is set low and the
	// c-accesses are started. Once started, one c-access is done in the second
	// phase of every clock cycle in the range 15-54. The read data is stored
	// in the video matrix/color line at the position specified by VMLI. These
	// data is internally read from the position specified by VMLI as well on
	// each g-access in display state.
	if !v.skipFrame && v.rasterLine&0x07 == v.scrollY && v.rasterLine >= DMAStart && v.rasterLine <= DMAEnd {
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

func (v *VicII) idleGap() {
	// TODO: Do idle accesses here.
}

func (v *VicII) handleSpriteYExpStep() {
	// In the first phase of cycle 15, it is checked if the expansion flip flop
	// is set. If so, MCBASE is incremented by 2.
	for i := range v.sprites {
		if v.sprites[i].expandYFF {
			v.sprites[i].mcBase += 2
		}
	}
}

func (v *VicII) cAccess() {
	// If the CPU hasn't given us the bus yet, we read $FF
	if !v.cpuBus.IsDMAAllowed() {
		v.cBuf[v.vmli] = 0xff
	}
	ch := v.bus.ReadByte(v.screenMemPtr | v.vc | v.getBankBase())
	col := v.colorRam.ReadByte(v.vc)
	if v.bitmapMode {
		v.cBuf[v.vmli] = uint16(ch)
	} else {
		v.cBuf[v.vmli] = uint16(col)<<8 | uint16(ch)
	}
}

func (v *VicII) gAccess() {
	if !v.displayState {
		// In idle state, only g-accesses occur. The access is always to address
		// $3fff ($39ff when the ECM bit in register $d016 is set). The graphics
		// are displayed by the sequencer exactly as in display state, but with
		// the video matrix data treated as "0" bits.
		addr := uint16(0x3fff)
		if v.extendedClr {
			addr = 0x39ff
		}
		v.gBuf[v.vmli] = v.bus.ReadByte(addr + v.getBankBase())
		v.vmli = (v.vmli + 1) & 0x003f
		return
	}
	basePtr := v.charSetPtr | v.getBankBase()
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
	addr := v.screenMemPtr | 0x03f8 | spriteIndex | v.getBankBase()
	v.sprites[spriteIndex].pointer = uint16(v.bus.ReadByte(addr)) << 6

	// Check if the next sprite will need bus. In that case, grab the bus now.
	next := (spriteIndex + 1) & 0x07
	if v.sprites[next].dma {
		v.bus.RDY.PullDown()
	}
}

func (v *VicII) sAccess(spriteIndex uint16) {
	s := &v.sprites[spriteIndex]
	if !s.dma {
		s.gData[0] = 0 // TODO: Should this really be needed?
		s.gData[1] = 0
		s.gData[2] = 0
		return
	}
	s.gData[s.sIndex] = v.bus.ReadByte(s.pointer | uint16(s.mc) | v.getBankBase())
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
	} else {
		contentTop = v.dimensions.ContentTop24Lines
		contentBottom = v.dimensions.ContentBottom24Lines
	}

	localCycle := v.cycle % v.dimensions.CyclesPerLine
	startPixel := (localCycle - v.dimensions.FirstVisibleCycle) << 3

	lastFgColor := uint8(0)
	cycleSpriteColl := uint8(0) // Sprite collisions during this cycle
	cycleBgColl := uint8(0)     // Sprite-bg collisions this cycle

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
		if v.rasterLine == contentTop && v.enable {
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
				c, lpSprites := v.generateSpriteColor(pixel, false) // Low priority sprites
				if lpSprites != 0 {
					color = c
				}
				graphics := false
				if v.bitmapMode {
					c, graphics = v.generateBitmapColor()
				} else {
					c, graphics = v.generateTextColor()
				}
				// Draw graphics background if an only if no low priority sprite was drawn.
				// In other words, screen background will have the lowest priority.
				if graphics || !(lpSprites != 0) {
					color = c
				}
				c, hpSprites := v.generateSpriteColor(pixel, true) // High priority sprites
				if hpSprites != 0 {
					color = c
				}
				lastFgColor = color

				// Update collision accumulators
				allSprites := lpSprites | hpSprites
				if graphics {
					cycleBgColl |= allSprites
				}
				// If zero or one bits are set, there were no collisions. Otherwise, all sprites that were
				// drawn collided. Use a single round of Kernighan's algorithm to check.
				if allSprites != 0 && allSprites&(allSprites-1) != 0 {
					cycleSpriteColl |= allSprites
				}
			}
		}
		v.screen.SetPixel(pixel, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[color&0x0f])
		v.sequencer <<= 1
	}
	// New collisions during this cycle? Fire interrupt if enabled!
	if v.spriteSpriteColl == 0 && cycleSpriteColl != 0 {
		v.irqSpriteSprite = true
		if v.irqSpriteSpriteEnabled {
			v.bus.NotIRQ.PullDown()
		}
	}
	if v.spriteDataColl == 0 && cycleBgColl != 0 {
		v.irqSpriteBg = true
		if v.irqSpriteBgEnabled {
			v.bus.NotIRQ.PullDown()
		}
	}

	v.spriteSpriteColl |= cycleSpriteColl
	v.spriteDataColl |= cycleBgColl
}

func (v *VicII) generateTextColor() (uint8, bool) {
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
			return fgColor & 0x0f, true
		} else {
			// Color 1 is treated as background too.
			return v.backgroundColors[cIndex] & 0x0f, cIndex != 1
		}
	} else {
		if v.sequencer&0x80 != 0 {
			return fgColor, true
		} else {
			return v.backgroundColors[bgIndex], false
		}
	}
}

func (v *VicII) generateBitmapColor() (uint8, bool) {
	bgColor := uint8(v.cData & 0x0f)
	fgColor := uint8(v.cData>>4) & 0x0f
	if v.multiColor {
		// N.B. This part must only be execute when x is even. Caller is responsible for that.
		cIndex := v.sequencer & 0xc0 >> 6
		switch cIndex {
		case 0:
			return v.backgroundColors[0], false
		case 1:
			return uint8(v.cData>>4) & 0x0f, false // 01 is treated as background
		case 2:
			return uint8(v.cData) & 0x0f, true
		case 3:
			return uint8(v.cData>>8) & 0x0f, true
		}
	} else {
		if v.sequencer&0x80 != 0 {
			return fgColor, true
		} else {
			return bgColor, false
		}
	}
	return 0, false // This should never happen
}

func (v *VicII) generateSpriteColor(x uint16, highPriority bool) (uint8, uint8) {
	color := uint8(0)
	spritesDrawn := uint8(0)
	for i := range v.sprites {
		spritesDrawn <<= 1
		s := &v.sprites[7-i] // Draw from lowest to highest to respect sprite priority
		if s.hasPriority != highPriority {
			continue
		}
		// Are we repeating pixels? Just return the previous one!
		if s.repeatCnt > 0 {
			s.repeatCnt--
			color = s.lastColor
			if s.lastWasDrawn {
				spritesDrawn |= 0x01
			}
			continue
		}

		// Determine for how many pixels (if any) we should repeat
		idx := x - s.x
		r := uint8(1)
		if s.multicolor {
			r <<= 1
		}
		if s.expandedX {
			r <<= 1
			idx >>= 1
		}
		s.repeatCnt = r - 1
		if idx <= 16 && idx&0x07 == 0 {
			s.sequencer = s.gData[idx>>3]
		}
		drawn := false
		if s.enabled {
			if s.multicolor {
				switch (s.sequencer & 0xc0) >> 6 {
				case 1:
					color = v.spriteMultiClr0
					drawn = true
				case 2:
					color = s.color
					drawn = true
				case 3:
					color = v.spriteMultiClr1
					drawn = true
				}
				s.sequencer <<= 1 // Advance an extra step since we used up two bits
			} else {
				if s.sequencer&0x80 != 0 {
					color = s.color & 0x0f
					drawn = true
				}
			}
			if drawn {
				spritesDrawn |= 0x01
			}
			s.lastColor = color
			s.lastWasDrawn = drawn
		}
		s.sequencer <<= 1
	}
	return color, spritesDrawn
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
			s.expandYFF = !s.expandYFF
		}
		if s.enabled && s.y == uint8(v.rasterLine&0xff) && !s.dma {
			s.mcBase = 0
			s.dma = true
			s.expandYFF = !s.expandedY // If not expanded, flip-flop will remain set forever
		}
	}
}

 // In the first phase of cycle 58, the MC of every sprite is loaded from
 // its belonging MCBASE (MCBASE->MC) and it is checked if the DMA for the
 // sprite is turned on and the Y coordinate of the sprite matches the lower
 // 8 bits of RASTER. If this is the case, the display of the sprite is
 // turned on.
func (v *VicII) resetSpriteCounters() {
	for i := range v.sprites {
		s := &v.sprites[i]
		s.mc = s.mcBase
		if s.dma && uint8(v.rasterLine & 0xff) == s.y {
			s.displayEnabled = true
		}
	}
}

func (v *VicII) getBankBase() uint16 {
	return uint16(^v.cia2.PortA.ReadOutputs()&0x03) << 14
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
