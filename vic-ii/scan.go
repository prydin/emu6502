package vic_ii

func (v *VicII) Clock() {
	if v.clockPhase2 {
		v.renderCycle()
		v.bus.ClockPh2()
	} else {
		// Stun the CPU if we hit a bad line
		if v.GetRasterLine()&0x07==v.scrollY {
			v.badLine = true
			v.bus.RDY.PullDown()
		} else {
			if v.badLine {
				// Getting out of a bad line. Let the RDY pin float.
				v.badLine = false
				v.bus.RDY.Release()
			}
		}
		v.bus.ClockPh1()
	}
	v.clockPhase2 = !v.clockPhase2
}

func (v *VicII) renderCycle() {
	line := v.cycle / (v.dimensions.ScreenWidth >> 3)
	localPixel := (v.cycle % (v.dimensions.ScreenWidth >> 3)) << 3

	// In the visible area?
	if line >= v.dimensions.TopBorder &&
		line < v.dimensions.BottomBlank &&
		localPixel >= v.dimensions.LeftBorder &&
		localPixel < v.dimensions.RightBlank {

		rasterLine := v.GetRasterLine()

		// Hit raster interrupt line?
		if localPixel == 0 && v.irqRasterEnabled && rasterLine == v.rasterLineTrigger {
			v.irqRaster = true
			v.bus.NotIRQ.PullDown()
		}
		// Border
		if localPixel < v.dimensions.LeftContent ||
			line < v.dimensions.TopContent ||
			localPixel >= v.dimensions.RightBorder ||
			line >= v.dimensions.BottomBorder {
			for i := 0; i < 8; i++ {
				v.screen.setPixel(localPixel-v.dimensions.LeftBorder+uint16(i), line, C64Colors[v.borderCol])
			}
		} else if v.bitmapMode {
			/// TODO
		} else {
			v.renderText(localPixel-v.dimensions.LeftContent, line)
		}
	}

	// Move to next cycle
	v.cycle++
	if v.cycle >= v.dimensions.Cycles {
		v.cycle = 0
	}
}

func (v *VicII) renderText(x, y uint16) {
	// TODO: Handle smooth scroll
	leftBorderOffset := v.dimensions.LeftContent - v.dimensions.LeftBorder
	line := y - v.dimensions.TopContent
	memOffset := x >> 3 + (line >> 3) * 40 // TODO: Handle 38 cols
	bgIndex := 0
	ch := v.bus.ReadByte(v.screenMemPtr + memOffset)
	if v.extendedClr {
		bgIndex = int(ch >> 6)
		ch &= 0x3f
	}
	pattern := uint8(0)
	if ch != 0 {
		pattern = v.bus.ReadByte(uint16(ch-1)<<3 + v.charSetPtr + line&0x07) // TODO: What about shifted characters?
	}
	fgColor := v.bus.ReadByte(memOffset + 0xd800) // TODO: Prefetch during bad lines
	for i := uint16(0); i < 8; i++ {
		color := uint8(0)
		if pattern & 0x80 != 0 {
			color = fgColor
		} else {
			color = v.backgroundColors[bgIndex]
		}
		pattern <<= 1
		v.screen.setPixel(x + i + leftBorderOffset, y, C64Colors[color])
	}
}
