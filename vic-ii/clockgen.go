package vic_ii

func (v *VicII) Clock() {
	if v.clockPhase2 {
		v.renderCycle()
		v.bus.ClockPh2()
	} else {
		// Let the CPU tick unless it's stunned due to a bad line
		if v.rasterLine < v.dimensions.TopContent ||
			v.rasterLine >= v.dimensions.BottomBorder ||
			v.rasterLine&0x07 != v.scrollY {
			v.bus.ClockPh1()
		}
	}
	v.clockPhase2 = !v.clockPhase2
}
