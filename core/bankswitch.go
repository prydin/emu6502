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

type Bank struct {
	selector  int
	devices   []AddressSpace
	writeable []bool
}

type BankSwitcher struct {
	banks []Bank
}

func (b *Bank) init() {
	b.writeable = make([]bool, len(b.devices))
	for i, d := range b.devices {
		b.writeable[i] = d.IsWriteable()
	}
}

func (b *Bank) ReadByte(addr uint16) uint8 {
	return b.devices[b.selector].ReadByte(addr)
}

func (b *Bank) WriteByte(addr uint16, data uint8) {
	// If we're trying to write to a read-only area (like a ROM), we should
	// instead write to the RAM that's mapped to the same address. We do
	// this by temporarily switching to bank 0 (all RAM)
	if !b.writeable[b.selector] {
		b.devices[0].WriteByte(addr, data)
	} else {
		b.devices[b.selector].WriteByte(addr, data)
	}

}

func (b *Bank) IsWriteable() bool {
	return b.devices[b.selector].IsWriteable()
}

func (bs *BankSwitcher) Switch(selector int) {
	/*if selector != bs.banks[0].selector {
		println(selector)
	}*/
	for i := range bs.banks {
		bs.banks[i].selector = selector
	}
}

func (bs *BankSwitcher) GetBank(index int) AddressSpace {
	return &bs.banks[index]
}

func NewBankSwitcher(devices [][]AddressSpace) *BankSwitcher {
	bs := BankSwitcher{banks: make([]Bank, len(devices))}
	for i, d := range devices {
		bs.banks[i] = Bank{devices: d}
		bs.banks[i].init()
	}
	return &bs
}
