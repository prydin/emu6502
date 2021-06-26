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

package core

import "fmt"

type AddressSpace interface {
	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)

	IsWriteable() bool
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
	pages      [256]Device
	phase1     []Clockable
	phase2     []Clockable
	dmaAllowed bool
	switcher   *BankSwitcher

	// Event pins
	RDY    TriState
	NotIRQ TriState
	NotNMI TriState
}

type PagedSpace struct {
	pages     []AddressSpace
	writeable bool
}

func MakeRAM(size uint16) *RAM {
	return &RAM{Bytes: make([]uint8, size)}
}

func (b *Bus) Connect(device AddressSpace, start, end uint16) {
	d := Device{start, end, device}
	for page := start >> 8; page <= end>>8; page++ {
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

	d := b.pages[addr>>8]
	if d.device != nil {
		if addr == 0xdd00 {
			fmt.Printf("CIA2 read: %02x\n", d.device.ReadByte(addr-d.start))
		}
		return d.device.ReadByte(addr - d.start)
	} else {
		return 0
	}
}

func (b *Bus) WriteByte(addr uint16, data uint8) {
	// Special case: Write to the "CPU port" configures the bank switcher
	/*if addr == 0xdd00 {
		fmt.Printf("CIA2 write: %02x\n", data)
	}*/
	if addr == 0x0001 && b.switcher != nil {
//		fmt.Printf("CPU port write: %02x\n", data)
		b.switcher.Switch(int(data & 0x07))
	}
	d := b.pages[addr>>8]
	d.device.WriteByte(addr-d.start, data)
}

func (b *Bus) IsWriteable() bool {
	return true // Uhm... Yeah sure. At least partially.
}

func (b *Bus) SetSwitcher(switcher *BankSwitcher) {
	b.switcher = switcher
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

func (r *ROM) IsWriteable() bool {
	return false
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
		fmt.Printf("WARNING: Attempt to write outside RAM: %04x\n", addr)
	}
}

func (r *RAM) IsWriteable() bool {
	return true
}

func (r *RAM) Page(page int) *RAM {
	return &RAM{
		Bytes: r.Bytes[page<<8 : (page+1)<<8],
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
		return page.ReadByte(addr & 0xff)
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
		page.WriteByte(addr&0xff, data)
	}
}

func (p *PagedSpace) IsWriteable() bool {
	return p.writeable
}

func NewPagedSpace(pages []AddressSpace, writeable bool) *PagedSpace {
	return &PagedSpace{pages: pages, writeable: writeable}
}
