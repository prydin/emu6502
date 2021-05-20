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
	selector int
	devices []AddressSpace
}

type BankSwitcher struct {
	banks []Bank
}

func (b *Bank) ReadByte(addr uint16) uint8 {
	return b.devices[b.selector].ReadByte(addr)
}

func (b *Bank) WriteByte(addr uint16, data uint8) {
	b.devices[b.selector].WriteByte(addr, data)
}

func (bs* BankSwitcher) Switch(selector int) {
	for i := range bs.banks {
		bs.banks[i].selector = selector
	}
}

func (bs *BankSwitcher) GetBank(index int) AddressSpace{
	return &bs.banks[index]
}

func NewBankSwitcher(devices [][]AddressSpace) *BankSwitcher {
	bs := BankSwitcher{ banks: make([]Bank, len(devices)) }
	for i, d := range devices {
		bs.banks[i] = Bank{ devices: d}
	}
	return &bs
}




