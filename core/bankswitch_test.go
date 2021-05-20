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

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBankSwitcher(t *testing.T) {
	ram0_1 := MakeRAM(1024)
	ram0_2 := MakeRAM(1024)
	ram1_1 := MakeRAM(1024)
	ram1_2 := MakeRAM(1024)
	ram2_1 := MakeRAM(1024)
	ram2_2 := MakeRAM(1024)
	bs := NewBankSwitcher([][]AddressSpace{
		{ ram0_1, ram0_2 },
		{ ram1_1, ram1_2 },
		{ ram2_1, ram2_2 },
	})

	bs.Switch(0)
	bus := Bus{}
	bus.Connect(bs.GetBank(0), 0x0000, 0x00ff)
	bus.Connect(bs.GetBank(1), 0x1000, 0x10ff)
	bus.WriteByte(0x0000, 0)
	bus.WriteByte(0x1000, 1)
	bs.Switch(1)
	bus.WriteByte(0x0000, 2)
	bus.WriteByte(0x1000, 3)
	bs.Switch(0)
	require.Equal(t, uint8(0), bus.ReadByte(0x0000))
	require.Equal(t, uint8(1), bus.ReadByte(0x1000))
	bs.Switch(1)
	require.Equal(t, uint8(2), bus.ReadByte(0x0000))
	require.Equal(t, uint8(3), bus.ReadByte(0x1000))
}
