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

package keyboard

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/prydin/emu6502/cia"
	"github.com/prydin/emu6502/core"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestKeyProvider struct {
	keystrokes []pixelgl.Button
	pos int
}

func (t* TestKeyProvider) Pressed(b pixelgl.Button) bool {
	pressed := b == t.keystrokes[t.pos]
	return pressed
}

func TestKeyboard_ScanAll(t *testing.T) {
	colBits := []uint8{8, 2, 32, 32, 16}
	//rowBits := []uint8 { 16, 64, 2, 2, 4}
	bus := core.Bus{}
	c := cia.CIA{}
	c.Init(&bus)
	k := Keyboard{}
	k.Init(&c)
	bus.Connect(&c, 0xdc00, 0xdcff)
	p := TestKeyProvider{
		keystrokes: []pixelgl.Button{
			pixelgl.KeyH,
			pixelgl.KeyE,
			pixelgl.KeyL,
			pixelgl.KeyL,
			pixelgl.KeyO,
		},
	}
	k.SetProvider(&p)
	bus.WriteByte(0xdc02, 0xff)
	bus.WriteByte(0xdc03, 0x00)
	i := 0

	// Scan entire matrix
	for p.pos = 0; p.pos < len(p.keystrokes); p.pos++ {
		bus.WriteByte(0xdc00, 0x00)
		k.Clock()
		key := ^bus.ReadByte(0xdc01)
		require.Equal(t, colBits[i], key, "Keystroke column mismatch")
		i++
	}
}

func TestKeyboard_ScanRow(t *testing.T) {
	colBits := []uint8{8, 2, 32, 32, 16}
	//rowBits := []uint8 { 16, 64, 2, 2, 4}
	bus := core.Bus{}
	c := cia.CIA{}
	c.Init(&bus)
	k := Keyboard{}
	k.Init(&c)
	bus.Connect(&c, 0xdc00, 0xdcff)
	p := TestKeyProvider{
		keystrokes: []pixelgl.Button{
			pixelgl.KeyH,
			pixelgl.KeyE,
			pixelgl.KeyL,
			pixelgl.KeyL,
			pixelgl.KeyO,
		},
	}
	k.SetProvider(&p)
	bus.WriteByte(0xdc02, 0xff)
	bus.WriteByte(0xdc03, 0x00)
	i := 0


	// Scan row by row
	var key uint8
	i = 0
	for p.pos = 0; p.pos < len(p.keystrokes); p.pos++ {
		mask := uint8(0x80)
		for mask != 0 {
			bus.WriteByte(0xdc00, ^mask)
			k.Clock()
			key = ^bus.ReadByte(0xdc01)
			if key != 0x00 {
				break
			}
			mask >>= 1
		}
		require.Equal(t, colBits[i], key, "Keystroke column mismatch")
		i++
	}
}