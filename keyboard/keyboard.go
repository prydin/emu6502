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
	"fmt"
	p "github.com/faiface/pixel/pixelgl"
	"github.com/prydin/emu6502/cia"
)

const Shift = 1 << 32

var keyMatrix = [][]p.Button{
	{p.KeyEscape, p.KeyQ, p.KeyLeftSuper, p.KeySpace, p.Key2, p.KeySpace, p.KeyBackspace, p.Key1},
	{p.KeySlash, p.Key6 | Shift, p.KeyEqual, p.KeyRightShift, p.KeyHome, p.KeySemicolon, p.Key8 | Shift, p.KeyBackslash},
	{p.KeyComma, p.Key2 | Shift, p.KeySemicolon | Shift, p.KeyPeriod, p.KeyMinus, p.KeyL, p.KeyP, p.KeyEqual | Shift},
	{p.KeyN, p.KeyO, p.KeyK, p.KeyM, p.Key0, p.KeyJ, p.KeyI, p.Key9},
	{p.KeyV, p.KeyU, p.KeyH, p.KeyB, p.Key8, p.KeyG, p.KeyY, p.Key7},
	{p.KeyX, p.KeyT, p.KeyF, p.KeyC, p.Key6, p.KeyD, p.KeyR, p.Key5},
	{p.KeyLeftShift, p.KeyE, p.KeyS, p.KeyZ, p.Key4, p.KeyA, p.KeyW, p.Key3},
	{p.KeyDown, p.KeyF5, p.KeyF3, p.KeyF1, p.KeyF7, p.KeyRight, p.KeyEnter, p.KeyDelete},
}

type KeyProvider interface {
	Pressed(b p.Button) bool
}

type Keyboard struct {
	cia      *cia.CIA
	provider KeyProvider
}

func (k *Keyboard) Init(cia *cia.CIA) {
	k.cia = cia
}

func (k *Keyboard) SetProvider(provider KeyProvider) {
	k.provider = provider
}

func (k *Keyboard) Clock() {
	k.cia.PortB.SetInputs(k.scan())
}

func (k *Keyboard) scan() uint8 {
	mask := k.cia.PortA.ReadOutputs()
	if mask == 0xff {
		return 0xff
	}

	// Count cleared bits and find position of the last one
	m := ^mask
	n := 0
	p := 0
	i := 0
	for m != 0 {
		if m&0x01 != 0 {
			p = i
			n++
		}
		m >>= 1
		i++
	}

	// Fast path: Only one bit cleared
	if n == 1 {
		return k.scanRow(7-p)
	}
	// Slow path: Scan entire matrix
	result := uint8(0xff)
	for row := 0; row < 8; row++ {
		if mask&0x80 == 0 {
			cr := k.scanRow(row)
			result &= cr
		}
		mask <<= 1
	}
	return result
}

func (k *Keyboard) scanRow(row int) uint8 {
	result := uint8(0)
	for col := 0; col < 8; col++ {
		result <<= 1
		// TODO: Handle shift
		if !k.provider.Pressed(keyMatrix[row][col] &^ Shift) {
			result |= 0x01
		}
	}
	if result != 0xff {
		fmt.Println("Key pressed:", ^result)
	}
	return result
}
