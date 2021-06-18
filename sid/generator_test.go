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

package sid

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerator_Sawtooth(t *testing.T) {
	gen := Generator{
		freq:     1,
		waveform: 2,
	}
	for i := reg24(0); i < (1<<24)-1; i++ {
		gen.Clock()
		gen.ReadOutput()
		sample := gen.ReadOutput()
		if reg12((i+1)>>12) != sample {
			require.Equalf(t, reg12((i+1)>>12), sample, "Sample mismatch at %d", i)
		}
	}
}

func TestGenerator_Triangle(t *testing.T) {
	gen := Generator{
		freq:     1,
		waveform: 1,
	}
	for i := reg24(0); i < (1<<24)-1; i++ {
		gen.Clock()
		gen.ReadOutput()
		sample := gen.ReadOutput()
		if i < 0x007fffff {
			if reg12((i+1)>>11) != sample {
				require.Equalf(t, reg12((i+1)>>11), sample, "Sample mismatch at %08x", i)
			}
		} else {
			actual := (^reg12((i + 1) >> 11)) & 0xfff
			if actual != sample {
				require.Equalf(t, actual, sample, "Sample mismatch at %08x", i)
			}
		}
	}
}
