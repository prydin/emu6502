package core

import "fmt"

type AddressSpace interface {
	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)
}

type Clockable interface {
	Clock()
}

type TriState struct {
	pullers int
	edge    int
}

type ROM struct {
	Bytes []uint8
}

type RAM struct {
	Bytes []uint8
}

type Device struct {
	start  uint16
	end    uint16
	device AddressSpace
}

type Bus struct {
	pages 	   [256]Device
	phase1     []Clockable
	phase2     []Clockable
	dmaAllowed bool

	// Event pins
	RDY    TriState
	NotIRQ TriState
	NotNMI TriState
}

type PagedSpace struct {
	pages []AddressSpace
}

func MakeRAM(size uint16) *RAM {
	return &RAM{Bytes: make([]uint8, size)}
}

func (b *Bus) Connect(device AddressSpace, start, end uint16) {
	d := Device{start, end, device}
	for page := start >> 8; page <= end >> 8; page++ {
		b.pages[page] = d
	}
}

func (b *Bus) ConnectClockablePh1(device Clockable) {
	b.phase1 = append(b.phase1, device)
}

func (b *Bus) ConnectClockablePh2(device Clockable) {
	b.phase2 = append(b.phase2, device)
}

func (b *Bus) ClockPh1() {
	for _, c := range b.phase1 {
		c.Clock()
	}
}

func (b *Bus) ClockPh2() {
	for _, c := range b.phase2 {
		c.Clock()
	}
}

func (b *Bus) ReadByte(addr uint16) uint8 {
	d := b.pages[addr >> 8]
	if d.device != nil {
		return d.device.ReadByte(addr - d.start)
	} else {
		return 0
	}
}

func (b *Bus) WriteByte(addr uint16, data uint8) {
	d := b.pages[addr >> 8]
	d.device.WriteByte(addr-d.start, data)
}

func (b *Bus) CPUClaimBus() {
	b.dmaAllowed = false
}

func (b *Bus) CPUReleaseBus() {
	b.dmaAllowed = true
}

func (b *Bus) IsDMAAllowed() bool {
	return b.dmaAllowed
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
	} else {
		fmt.Printf("WARNING: Attempt to write outside RAM: %04x", addr)
	}
}

func (t *TriState) PullDown() {
	if t.pullers == 0 {
		t.edge = -1
	}
	t.pullers++
}

func (t *TriState) Release() {
	if t.pullers == 0 {
		return
	}
	t.pullers--
	if t.pullers == 0 {
		t.edge = 1
	}
}

func (t *TriState) GetEdge() int {
	e := t.edge
	t.edge = 0
	return e
}

func (t *TriState) Get() bool {
	return t.pullers == 0
}

func (p *PagedSpace) ReadByte(addr uint16) uint8 {
	n := int(addr >> 8)
	if n >= len(p.pages) {
		return 0
	}
	page := p.pages[n]
	if page != nil {
		return page.ReadByte(addr)
	}
	return 0
}

func (p *PagedSpace) WriteByte(addr uint16, data uint8) {
	n := int(addr >> 8)
	if n >= len(p.pages) {
		return
	}
	page := p.pages[n]
	if page != nil {
		page.WriteByte(addr, data)
	}
}

func NewPagedSpace(pages []AddressSpace) *PagedSpace {
	return &PagedSpace{ pages: pages }
}