package vic_ii

func (v *VicII) Clock() {
	localCycle := v.cycle % v.dimensions.CyclesPerLine

	if v.clockPhase2 {
		// ******** CLOCK PHASE 2 ********
		if v.badLine && localCycle >= 15 && localCycle <= 54 {
			v.cAccess()
		}
		if localCycle >= PalFirstVisibleCycle && localCycle <= PalLastVisibleCycle &&
			v.rasterLine >= PalFirstVisibleLine && v.rasterLine <= PalLastVisibleLine {
			v.renderCycle()
		}
		v.clockSink.ClockPh2()

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
			case 55:
				// Bad lines always end here regardless of how they started
				if v.badLine {
					v.badLine = false
					v.bus.RDY.Release()
				}
			case 58:
				if v.rc == 0x07 {
					v.rc = 0 // vic-ii.txt: 3.7.2.5
					v.vcBase = v.vc
					v.displayState = false
				}
			case 59:
				if v.displayState {
					v.rc = (v.rc + 1) & 0x07
				}
			default:
				if localCycle >= 16 && localCycle <= 56 {
					v.gAccess()
					if localCycle == 16 { // TODO: Documentation says display state is only enabled on bad lines, but that just can't be true. Or can it?
						v.displayState = true
					}
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
		}
		v.clockSink.ClockPh1()
	}
	v.clockPhase2 = !v.clockPhase2
}

func (v *VicII) cAccess() {
	if v.bitmapMode {
		// TODO
	} else {
		if v.extendedClr {
			// TODO
		} else {
			ch := v.bus.ReadByte(v.screenMemPtr | v.vc)
			col := v.colorRam.ReadByte(v.vc)
			v.cBuf[v.vmli] = uint16(ch) | uint16(col)<<8
		}
	}
}

func (v *VicII) gAccess() {
	if v.bitmapMode {
		// TODO
	} else {
		if v.extendedClr {
			// TODO
		} else {
			addr := v.charSetPtr + ((v.cBuf[v.vmli]&0xff)<<3 | v.rc)
			v.gBuf[v.vmli] = v.bus.ReadByte(addr)
		}
	}
	// Increment counters and make sure they stay within 6 and 10 bits boundary respectively
	v.vmli = (v.vmli + 1) & 0x003f
	v.vc = (v.vc + 1) & 0x03ff
}

func (v *VicII) rasterInterrupt() {
	if v.irqRasterEnabled && v.rasterLine == v.rasterLineTrigger {
		v.irqRaster = true
		v.bus.NotIRQ.PullDown()
	}
}

func (v *VicII) renderCycle() {
	var contentLeft uint16
	var contentRight uint16
	var contentTop uint16
	var contentBottom uint16
	var scrollOffset uint16
	if v.col40 {
		contentLeft = v.dimensions.LeftBorderWidth40Cols
		contentRight = contentLeft + v.dimensions.ContentWidth40Cols
	} else {
		contentLeft = v.dimensions.LeftBorderWidth40Cols
		contentRight = contentLeft + v.dimensions.ContentWidth38Cols
	}
	if v.line25 {
		contentTop = v.dimensions.ContentTop25Lines
		contentBottom = v.dimensions.ContentBottom25Lines
		scrollOffset = v.dimensions.OptimalYScroll25Lines
	} else {
		contentTop = v.dimensions.ContentTop24Lines
		contentBottom = v.dimensions.ContentBottom24Lines
		scrollOffset = v.dimensions.OptimalYScroll24Lines
	}

	localCycle := (v.cycle  % v.dimensions.CyclesPerLine) - v.dimensions.FirstVisibleCycle
	startPixel := localCycle << 3

	switch {
	// Vertical border?
	case v.rasterLine < contentTop || v.rasterLine >= contentBottom:
		v.drawBorder(startPixel, 8)
	// Left border?
	case startPixel < contentLeft:
		v.drawBorder(startPixel, min(8, contentLeft-1))
	// Right border?
	case startPixel > contentRight-8:
		v.drawBorder(max(startPixel, contentRight+1), 8)
	default:
		// Visible area? Maybe so, but check yscroll first!
		if v.rasterLine+v.scrollY-scrollOffset <= contentBottom - 8 {
			if v.bitmapMode {
				/// TODO
			} else {
				v.renderText(startPixel)
			}
		} else {
			v.drawBackground(startPixel, 8)
		}
	}
}

func (v *VicII) drawBorder(x, n uint16) {
	for i := uint16(0); i < n; i++ {
		v.screen.setPixel(i+x, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[v.borderCol])
	}
}

func (v *VicII) drawBackground(x, n uint16) {
	for i := uint16(0); i < n; i++ {
		v.screen.setPixel(i+x, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[v.backgroundColors[0]]) // TODO: Support multicolor bg?
	}
}

func (v *VicII) renderText(x uint16) {
	// TODO: Handle smooth scroll
	leftBorderOffset := uint16(0)
	index := v.cycle%v.dimensions.CyclesPerLine - v.dimensions.FirstContentCycle
	data := v.cBuf[index]
	fgColor := uint8(data >> 8)
	bgIndex := 0
	/*if v.extendedClr {
		bgIndex = int(ch >> 6)
		ch &= 0x3f
	}*/
	pattern := v.gBuf[index]
	for i := uint16(0); i < 8; i++ {
		color := uint8(0)
		if pattern&0x80 != 0 {
			color = fgColor
		} else {
			color = v.backgroundColors[bgIndex]
		}
		pattern <<= 1
		v.screen.setPixel(x+i+leftBorderOffset, v.rasterLine-v.dimensions.FirstVisibleLine, C64Colors[color & 0x0f])
	}
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
