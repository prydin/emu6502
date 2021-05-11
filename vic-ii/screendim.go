package vic_ii

// PAL screen constants
const (
	PalFirstVisibleCycle = 11
	PalFirstVisibleLine  = 16
	PalLastVisibleLine   = 287
	PalCyclesPerLine     = 63
	PalLastVisibleCycle = 58

	PalLeftBorderWidth40Cols  = 32
	PalRightBorderWidth40Cols = 32
	PalLeftBorderWidth38Cols  = 46
	PalRightBorderWidth38Cols = 36

	PalContentTop25Lines    = 51
	PalContentTop24Lines    = 55
	PalContentBottom25Lines = 251
	PalContentBottom24Lines = 247

	PalContentWidth40Cols = 320
	PalContentWidth38Cols = PalContentWidth40Cols - 16

	PalFirstContentCycle = 15
	PalLastContentCycle  = 56
	PalScreenWidth       = PalCyclesPerLine * 8
	PalScreenHeight      = 312

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


	ContentWidth40Cols: PalContentWidth40Cols,
	ContentWidth38Cols: PalContentWidth38Cols,
	FirstContentCycle:  PalFirstContentCycle,
	LastContentCycle:   PalLastContentCycle,

	FirstVisibleLine: PalFirstVisibleLine,
	LastVisibleLine:  PalLastVisibleLine,
	FirstVisibleCycle: PalFirstVisibleCycle,
	CyclesPerLine:    PalCyclesPerLine,
	Cycles:           PalCycles,
}
