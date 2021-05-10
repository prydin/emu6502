package vic_ii

// Hardware independent constants
const (
	ContentTop25Lines    = 51
	ContentTop24Lines    = 55
	ContentBottom25Lines = 251
	ContentBottom24Lines = 247
	ContentLeft40Cols    = 128 // Zero based
	ContentLeft38Cols    = 132 // Zero based
	ContentRight40Cols   = 448
	ContentRight38Cols   = 439
	FirstContentCycle    = 16
	LastContentCycle     = 56
)

// PAL screen constants
const (
	PalScreenHeight  = uint16(312)
	PalVisibleHeight = uint16(284)
	PalContentHeight = uint16(200)
	PalScreenWidth   = uint16(504)
	PalVisibleWidth  = uint16(403)
	PalContentWidth  = uint16(320)
	PalLeftBorder    = PalScreenWidth/2 - PalVisibleWidth/2
	PalLeftContent   = PalScreenWidth/2 - PalContentWidth/2
	PalRightBorder   = PalScreenWidth - PalLeftContent
	PalRightBlank    = PalScreenWidth - PalLeftBorder
	PalTopBorder     = PalScreenHeight/2 - PalVisibleHeight/2
	PalTopContent    = PalScreenHeight/2 - PalContentHeight/2
	PalBottomBorder  = PalScreenHeight - PalTopContent
	PalBottomBlank   = PalScreenHeight - PalTopBorder

	PalFirstVisibleLine      = 16
	PalLastVisibleLine       = 287
	PalCyclesPerLine         = 63
	PalLeftmostVisiblePixel  = 60  // TODO: Check!
	PalRightmostVisiblePixel = 462 // TODO: Check
	PalCycles                = uint16((uint32(PalScreenWidth) * uint32(PalScreenHeight)) / 8)
)

type ScreenDimensions struct {
	ScreenHeight  uint16
	VisibleHeight uint16
	ContentHeight uint16
	ScreenWidth   uint16
	VisibleWidth  uint16
	ContentWidth  uint16
	LeftBorder    uint16
	LeftContent   uint16
	RightBorder   uint16
	RightBlank    uint16
	TopBorder     uint16
	TopContent    uint16
	BottomBorder  uint16
	BottomBlank   uint16

	FirstVisibleLine      uint16
	LastVisibleLine       uint16
	LeftmostVisiblePixel  uint16
	RightmostVisiblePixel uint16
	CyclesPerLine         uint16
	Cycles                uint16
}

var PALDimensions = ScreenDimensions{
	ScreenHeight:          PalScreenHeight,
	VisibleHeight:         PalVisibleHeight,
	ContentHeight:         PalContentHeight,
	ScreenWidth:           PalScreenWidth,
	VisibleWidth:          PalScreenWidth,
	ContentWidth:          PalContentWidth,
	LeftBorder:            PalLeftBorder,
	LeftContent:           PalLeftContent,
	RightBorder:           PalRightBorder,
	RightBlank:            PalRightBlank,
	TopBorder:             PalTopBorder,
	TopContent:            PalTopContent,
	BottomBorder:          PalBottomBorder,
	BottomBlank:           PalBottomBlank,
	FirstVisibleLine:      PalFirstVisibleLine,
	LastVisibleLine:       PalLastVisibleLine,
	LeftmostVisiblePixel:  PalLeftmostVisiblePixel,
	RightmostVisiblePixel: PalRightmostVisiblePixel,
	CyclesPerLine:         PalCyclesPerLine,
	Cycles:                PalCycles,
}
