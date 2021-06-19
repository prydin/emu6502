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

package cia

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPort_ReadOutputs(t *testing.T) {
	p := Port{}
	for ddr := uint8(0); ; {
		p.ddr = ddr
		for data := uint8(0); ; {
			p.internalWrite(data)
			d := p.ReadOutputs()
			require.Equal(t, data&ddr, d, "Data mismatch. ddr=%02x, data=%02x", ddr, data)
			if data == 255 {
				break
			}
			data++
		}
		if ddr == 255 {
			break
		}
		ddr++
	}
}

func TestPort_SetInputs(t *testing.T) {
	p := Port{}

	for ddr := uint8(0); ; {
		p.ddr = ddr
		for data := uint8(0); ; {
			p.SetInputs(data)
			d := p.internalRead()
			require.Equal(t, data, d, "Data mismatch. ddr=%02x, data=%02x", ddr, data)
			if data == 255 {
				break
			}
			data++
		}
		if ddr == 255 {
			break
		}
		ddr++
	}
}
