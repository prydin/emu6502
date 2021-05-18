package vic_ii

// PAL screen constants
const (
	PalFirstVisibleCycle = 11
	PalFirstVisibleLine  = 16
	PalLastVisibleLine   = 287
	PalCyclesPerLine     = 63
	PalLastVisibleCycle  = 58

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

	PalFirstContentCycle = 15
	PalLastContentCycle  = 56
	PalScreenWidth       = PalCyclesPerLine * 8
	PalScreenHeight      = 312

	PalOptimalYScroll25Lines = 3 // Top and bottom lines fully visible
	PalOptimalYScroll24Lines = 7 // Top and bottom lines fully visible

	PalCycles = PalCyclesPerLine * 312
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

	FirstVisibleLine:  PalFirstVisibleLine,
	LastVisibleLine:   PalLastVisibleLine,
	FirstVisibleCycle: PalFirstVisibleCycle,
	CyclesPerLine:     PalCyclesPerLine,
	Cycles:            PalCycles,
}
