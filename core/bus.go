package core

type AddressSpace interface {
	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)
}

type Clockable interface {
	Clock()
}

type TriState struct {
	pullers int
}

type ROM struct {
	Bytes []uint8
}

type RAM struct {
	Bytes []uint8
}

type Device struct {
	start uint16
	end uint16
	device AddressSpace
}

type Bus struct {
	devices []Device
	phase1[] Clockable
	phase2[] Clockable

	// Event pins
	RDY    TriState
	NotIRQ TriState
	NotNMI TriState
}

func (b *Bus) Connect(device AddressSpace, start, end uint16) {
	b.devices = append(b.devices, Device{ start, end, device })
}

func(b *Bus) ConnectClockablePh1(device Clockable) {
	b.phase1 = append(b.phase1, device)
}

func(b *Bus) ConnectClockablePh2(device Clockable) {
	b.phase2 = append(b.phase2, device)
}

func (b* Bus) ClockPh1() {
	for _, c := range b.phase1 {
		c.Clock()
	}
}

func (b* Bus) ClockPh2() {
	for _, c := range b.phase1 {
		c.Clock()
	}
}

func (b *Bus) ReadByte(addr uint16) uint8 {
	for _, d := range b.devices {
		if addr >= d.start && addr <= d.end {
			return d.device.ReadByte(addr - d.start)
		}
	}
	return 0
}

func (b *Bus) WriteByte(addr uint16, data uint8) {
	for _, d := range b.devices {
		if addr >= d.start && addr <= d.end {
			d.device.WriteByte(addr - d.start, data)
		}
	}
}

func (r *ROM) ReadByte(addr uint16) uint8 {
	if int(addr) < len(r.Bytes) {
		return r.Bytes[int(addr)]
	} else {
		return 0
	}
}

func (r *ROM) WriteByte(addr uint16, data uint8) {
}

func (r *RAM) ReadByte(addr uint16) uint8 {
	if int(addr) < len(r.Bytes) {
		return r.Bytes[int(addr)]
	} else {
		return 0
	}
}

func (r *RAM) WriteByte(addr uint16, data uint8) {
	if int(addr) < len(r.Bytes) {
		r.Bytes[int(addr)] = data
	}
}

func (t *TriState) PullDown() {
	t.pullers++
}

func (t *TriState) Release() {
	t.pullers--
}

func (t *TriState) Get() bool {
	return t.pullers == 0
}





