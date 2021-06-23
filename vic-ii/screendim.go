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

// IMPORTANT NOTICE: Some VIC-II literature, including Christian Bauer's famous
// text file use one-based numbering of cycles within a line. In this code, we
// use zero-based values. Thus, everything will be off by one! Pay attention!
//
// PAL screen constants
const (
	PalFirstVisibleCycle = 10
	PalFirstVisibleLine  = 16
	PalLastVisibleLine   = 287
	PalCyclesPerLine     = 63
	PalLastVisibleCycle  = 57

	PalLeftBorderWidth40Cols  = 32
	PalRightBorderWidth40Cols = 32
	PalLeftBorderWidth38Cols  = 46
	PalRightBorderWidth38Cols = 36

	PalTopBorderHeight40Cols = 36
	PalTopBorderHeight38Cols = 40

	PalContentTop25Lines    = PalFirstVisibleLine + PalTopBorderHeight40Cols
	PalContentTop24Lines    = PalFirstVisibleLine + PalTopBorderHeight38Cols
	PalContentBottom25Lines = PalContentTop25Lines + 25*8
	PalContentBottom24Lines = PalContentTop24Lines + 24*8

	PalContentWidth40Cols = 320
	PalContentWidth38Cols = PalContentWidth40Cols - 16

	PalFirstContentCycle = 14
	PalLastContentCycle  = 55
	PalScreenWidth       = PalCyclesPerLine * 8
	PalScreenHeight      = 312

	PalOptimalYScroll25Lines = 3 // Top and bottom lines fully visible
	PalOptimalYScroll24Lines = 7 // Top and bottom lines fully visible

	PalCycles = PalCyclesPerLine * 312

	PalVisibleWidth  = PalContentWidth40Cols + PalRightBorderWidth40Cols + PalLeftBorderWidth40Cols
	PalVisibleHeight = PalContentBottom25Lines + PalTopBorderHeight40Cols

	DMAStart = 0x30
	DMAEnd   = 0xf7
)

type ScreenDimensions struct {
	ScreenHeight      uint16
	ScreenWidth       uint16
	FirstVisibleCycle uint16

	ContentTop25Lines      uint16
	ContentTop24Lines      uint16
	ContentBottom25Lines   uint16
	ContentBottom24Lines   uint16
	LeftBorderWidth40Cols  uint16
	LeftBorderWidth38Cols  uint16
	RightBorderWidth40Cols uint16
	RightBorderWidth38Cols uint16
	ContentWidth40Cols     uint16
	ContentWidth38Cols     uint16
	FirstContentCycle      uint16
	LastContentCycle       uint16
	OptimalYScroll25Lines  uint16
	OptimalYScroll24Lines  uint16

	VisibleHeight uint16
	VisibleWidth  uint16

	FirstVisibleLine uint16
	LastVisibleLine  uint16
	CyclesPerLine    uint16
	Cycles           uint16
}

var PALDimensions = ScreenDimensions{
	ScreenHeight: PalScreenHeight,
	ScreenWidth:  PalScreenWidth,

	ContentTop25Lines:    PalContentTop25Lines,
	ContentTop24Lines:    PalContentTop24Lines,
	ContentBottom25Lines: PalContentBottom25Lines,
	ContentBottom24Lines: PalContentBottom24Lines,

	LeftBorderWidth38Cols:  PalLeftBorderWidth38Cols,
	LeftBorderWidth40Cols:  PalLeftBorderWidth40Cols,
	RightBorderWidth38Cols: PalRightBorderWidth38Cols,
	RightBorderWidth40Cols: PalLeftBorderWidth40Cols,


	ContentWidth40Cols:    PalContentWidth40Cols,
	ContentWidth38Cols:    PalContentWidth38Cols,
	FirstContentCycle:     PalFirstContentCycle,
	LastContentCycle:      PalLastContentCycle,
	OptimalYScroll25Lines: PalOptimalYScroll25Lines,
	OptimalYScroll24Lines: PalOptimalYScroll24Lines,

	VisibleHeight: PalVisibleHeight,
	VisibleWidth:  PalVisibleWidth,

	FirstVisibleLine:  PalFirstVisibleLine,
	LastVisibleLine:   PalLastVisibleLine,
	FirstVisibleCycle: PalFirstVisibleCycle,
	CyclesPerLine:     PalCyclesPerLine,
	Cycles:            PalCycles,
}
