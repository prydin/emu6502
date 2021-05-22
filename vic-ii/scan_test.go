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

package vic_ii

import (
	"fmt"
	"github.com/prydin/emu6502/charset"
	"github.com/prydin/emu6502/core"
	"image"
	"image/png"
	"os"
	"testing"
	"time"
)

func initVicII(colorRam core.AddressSpace) (*VicII, *image.RGBA) {
	img := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{403, 312}})
	vicii := VicII{}
	bus := &core.Bus{}
	vicii.Init(bus, bus, colorRam, &ImageRaster{img}, PALDimensions)
	vicii.borderCol = 14
	return &vicii, img
}

func Test_BlankScreen(t *testing.T) {
	vicii, img := initVicII(core.MakeRAM(1024))
	start := time.Now()
	for i := 0; i < int(PalScreenWidth) * int(PalScreenHeight) / 4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("blankscreen.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_CharacterMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xd7ff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		if i > 1000 {
			screenMem[i] = 0xff
		} else if i % 40 == 0 {
			screenMem[i] = uint8((i / 40) % 10) + 0x30
		} else {
			screenMem[i] = uint8(i%10 + 0x30)
		}
	}
	vicii.bus.Connect(&core.RAM{ Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 3
	start := time.Now()
	for i := 0; i < int(PalScreenWidth) * int(PalScreenHeight) / 4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("characters.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_ExtendedCharacterMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xd7ff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i + 64)
	}
	vicii.bus.Connect(&core.RAM{ Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 1
	vicii.backgroundColors[1] = 2
	vicii.backgroundColors[2] = 3
	vicii.backgroundColors[3] = 4
	vicii.extendedClr = true

	start := time.Now()
	for i := 0; i < int(PalScreenWidth) * int(PalScreenHeight) / 4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("characters_ext.png")
	png.Encode(f, img)
	// TODO: Check image
}
