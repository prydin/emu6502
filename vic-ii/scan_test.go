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
	"github.com/beevik/go6502/asm"
	"github.com/prydin/emu6502/charset"
	"github.com/prydin/emu6502/core"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"os"
	"strings"
	"testing"
	"time"
)

func assemble(program string) ([]byte, error) {
	assy, _, err := asm.Assemble(strings.NewReader(program), "test.asm", os.Stderr, 0)
	for _, e := range assy.Errors {
		println(e)
	}
	return assy.Code, err
}

func initVicII(mainBus *core.Bus, colorRam core.AddressSpace) (*VicII, *image.RGBA) {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{403, 312}})
	vicii := VicII{}
	bus := &core.Bus{}
	if mainBus == nil {
		mainBus = bus
	}
	vicii.Init(bus, mainBus, colorRam, &ImageRaster{img}, PALDimensions)
	vicii.borderCol = 14
	return &vicii, img
}

func Test_BlankScreen(t *testing.T) {
	vicii, img := initVicII(nil, core.MakeRAM(1024))
	start := time.Now()
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("blankscreen.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_CharacterMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		if i > 1000 {
			screenMem[i] = 0xff
		} else if i%40 == 0 {
			screenMem[i] = uint8((i/40)%10) + 0x30
		} else {
			screenMem[i] = uint8(i%10 + 0x30)
		}
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 0
	start := time.Now()
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("characters.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_ExtendedCharacterMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i + 64)
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
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
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("characters_ext.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_MultiColorCharacterMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i % 64)
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 0
	vicii.backgroundColors[1] = 2
	vicii.backgroundColors[2] = 3
	vicii.multiColor = true

	start := time.Now()
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("characters_multi.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_BitmapMode(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
	vicii.bus.Connect(core.MakeRAM(0x2000), 0x2000, 0x3fff)
	b := uint8(1)
	for i := 0x2000; i < 0x3fff; i++ {
		vicii.bus.WriteByte(uint16(i), b)
		b <<= 1
		if b == 0 {
			b = 1
		}
	}
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i / 40)
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0x2000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 0
	vicii.bitmapMode = true
	start := time.Now()
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("bitmap.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_BitmapModeMultiColor(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
	vicii.bus.Connect(core.MakeRAM(0x2000), 0x2000, 0x3fff)
	for i := 0x2000; i < 0x3fff; i++ {
		vicii.bus.WriteByte(uint16(i), uint8(i))
	}
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i / 40)
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0x2000
	vicii.backgroundColors[0] = 6
	vicii.multiColor = true
	vicii.scrollY = 3
	vicii.scrollX = 0
	vicii.bitmapMode = true
	start := time.Now()
	for i := 0; i < int(PalScreenWidth)*int(PalScreenHeight)/4; i++ {
		vicii.Clock()
	}
	fmt.Printf("Rendering time: %s", time.Now().Sub(start))
	f, _ := os.Create("bitmap_multi.png")
	png.Encode(f, img)
	// TODO: Check image
}

func Test_RasterlinePolling(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	mainBus := core.Bus{}
	vicii, img := initVicII(&mainBus, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		if i > 1000 {
			screenMem[i] = 0xff
		} else if i%40 == 0 {
			screenMem[i] = uint8((i/40)%10) + 0x30
		} else {
			screenMem[i] = uint8(i%10 + 0x30)
		}
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 3

	mainBus.Connect(vicii, 0xd000, 0xd0ff)
	mainBus.Connect(core.MakeRAM(100), 0x0000, 0x00ff)
	cpu := core.CPU{}
	cpu.Init(&mainBus)
	mainBus.ConnectClockablePh1(&cpu)

	code, err := assemble(`
		LDA #$FF
		STA $02
LOOP	LDA	$D012
		CMP $02
		BEQ LOOP
		STA $02
		AND #$07
		BNE NOTZERO
		INC $D021
		JMP LOOP
NOTZERO	CMP #$04
		BNE	LOOP
		INC $D021
		JMP LOOP
`)
	if err != nil {
		panic(err)
	}
	mainBus.Connect(&core.RAM{Bytes: code}, 0x1000, 0x1000+uint16(len(code)))
	cpu.SetPC(0x1000)
	for i := 0; i < 100000; i++ {
		vicii.Clock()
	}
	f, _ := os.Create("raster_poll.png")
	png.Encode(f, img)
}

func Test_RasterlineInterrupt(t *testing.T) {
	colorRam := core.MakeRAM(1024)
	mainBus := core.Bus{}
	vicii, img := initVicII(&mainBus, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		if i > 1000 {
			screenMem[i] = 0xff
		} else if i%40 == 0 {
			screenMem[i] = uint8((i/40)%10) + 0x30
		} else {
			screenMem[i] = uint8(i%10 + 0x30)
		}
	}
	vicii.bus.Connect(&core.RAM{Bytes: screenMem[:]}, 0x0400, 0x07ff)
	for i := uint16(0); i < 1024; i++ {
		colorRam.WriteByte(i, 14)
	}
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 3

	mainBus.Connect(vicii, 0xd000, 0xd0ff)
	mainBus.Connect(core.MakeRAM(0x300), 0x0000, 0x02ff)
	mainBus.Connect(core.MakeRAM(0x1000), 0xf000, 0xffff)

	cpu := core.CPU{}
	cpu.Init(&mainBus)
	mainBus.ConnectClockablePh1(&cpu)

	code, err := assemble(`
		SEI
		LDA #<IRQ	; Set IRQ vector
		STA $FFFE
		LDA #>IRQ
		STA $FFFF
		LDA #$00
		STA $D021
		LDA #$30
		STA $D012	; Trigger IRQ on raster line $30
		LDA $D011	; Make sure high bit is zero
		AND	#$7f
		STA $D011
		LDA $D01A
		ORA	#$01	; Enable raster IRQ
		STA $D01A
		CLI
HANG	JMP HANG	; Wait forever
IRQ		INC $D021	; Change background color
		DEC $D019	; Acknowledge IRQ
		LDA $D012	
		ADC	#$04	; Set new interrupt to be triggered 4 lines from here
		CMP #$FC	; Reached bottom?
		BCC INCR
		LDA #$30	; Reset raster line IRQ to top of screen
		STA $D012	; Trigger IRQ on raster line $30
		LDA $D011	; Make sure high bit is zero
		LDA #$01
		STA $D021	; Start with a black band again
		AND	#$7f
		STA $D011
		LDA $D01A
		RTI
INCR	STA $D012	; Not at bottom. Keep going
		RTI
`)
	if err != nil {
		panic(err)
	}
	mainBus.Connect(&core.RAM{Bytes: code}, 0x1000, 0x1000+uint16(len(code)))
	cpu.SetPC(0x1000)
	for i := 0; i < 10000000; i++ {
		vicii.Clock()
	}
	f, _ := os.Create("raster_irq.png")
	png.Encode(f, img)
}

func Test_getSpriteFromCycle(t *testing.T) {
	sprites := map[uint16]uint16{
		0:  2,
		1:  3,
		2:  3,
		3:  4,
		4:  4,
		5:  5,
		6:  5,
		7:  6,
		8:  6,
		9:  7,
		10: 7,
		58: 0,
		59: 0,
		60: 1,
		61: 1,
		62: 2,
	}
	for i := uint16(0); i < 63; i++ {
		s, pAccess := getSpriteForCycle(i)
		index, present := sprites[i]
		if !present {
			index = 0xff
		}
		require.Equalf(t, index, s, "Sprite index mismatch at %d", i)
		if present && i < 11 && i&1 == 1 || i > 57 && i&1 == 0 {
			require.Truef(t, pAccess, "pAccess mismatch at %d", i)
		} else {
			require.Falsef(t, pAccess, "pAccess mismatch at %d", i)
		}
	}
}

func TestPAccess(t *testing.T) {
	cycles := []uint16{58, 60, 62, 1, 3, 5, 7, 9}
	for i := 0; i < 8; i++ {
		colorRam := core.MakeRAM(1024)
		vicii, _ := initVicII(nil, colorRam)
		vicii.bus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
		vicii.bus.Connect(core.MakeRAM(0x2000), 0x2000, 0x3fff)
		screenRAM := core.MakeRAM(1024)
		vicii.bus.Connect(screenRAM, 0x0400, 0x07ff)
		screenRAM.WriteByte(0x03f8+uint16(i), uint8(0x2000>>6))
		vicii.sprites[i].enabled = true
		vicii.sprites[i].x = 100
		vicii.sprites[i].y = 200
		vicii.cycle = cycles[i]
		vicii.Clock()
		require.Equalf(t, uint16(0x2000), vicii.sprites[i].pointer, "Pointer mismatch at sprite %d", i)
	}
}

func TestSAccess(t *testing.T) {
	cycles := []uint16{58, 60, 62, 1, 3, 5, 7, 9}
	for i := 0; i < 8; i++ {
		colorRam := core.MakeRAM(1024)
		vicii, _ := initVicII(nil, colorRam)
		vicii.bus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
		vicii.bus.Connect(core.MakeRAM(0x2000), 0x2000, 0x3fff)
		screenRAM := core.MakeRAM(1024)
		vicii.bus.Connect(screenRAM, 0x0400, 0x07ff)
		screenRAM.WriteByte(0x03f8+uint16(i), uint8(0x2000>>6))
		for c := uint16(0); c < 3; c++ {
			vicii.bus.WriteByte(0x2000+uint16(c), uint8(c))
		}
		vicii.sprites[i].enabled = true
		vicii.sprites[i].x = 100
		vicii.sprites[i].y = 200
		vicii.sprites[i].dma = true
		vicii.cycle = cycles[i]
		for c := 0; c < 4; c++ {
			vicii.Clock()
		}
		require.Equalf(t, uint16(0x2000), vicii.sprites[i].pointer, "Pointer mismatch at sprite %d", i)
		require.Equalf(t, uint32(0x00000102), vicii.sprites[i].sequencer, "Data mismatch on sprite %d", i)
	}
}

func TestDrawSprite(t *testing.T) {
	cycles := []uint16{58, 60, 62, 1, 3, 5, 7, 9}
	colorRam := core.MakeRAM(1024)
	vicii, img := initVicII(nil, colorRam)
	vicii.bus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
	vicii.bus.Connect(core.MakeRAM(0x2000), 0x2000, 0x3fff)
	screenRAM := core.MakeRAM(1024)
	vicii.bus.Connect(screenRAM, 0x0400, 0x07ff)
	vicii.bus.Connect(colorRam, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
	vicii.scrollY = 3
	vicii.scrollX = 3
	for i := 0; i < 8; i++ {
		screenRAM.WriteByte(0x03f8+uint16(i), uint8(0x2000>>6))
		for c := uint16(0); c < 63; c++ {
			data := uint8(1 << ((c / 3) % 8))
			if c%3 == 1 {
				data = ^data
			}
			vicii.bus.WriteByte(0x2000+uint16(c)+uint16(i*63), data)
		}
		vicii.sprites[i].enabled = true
		vicii.sprites[i].x = uint16(26 + i*48)
		vicii.sprites[i].y = uint8(50 + i*21)
		vicii.sprites[i].color = 1
		vicii.sprites[i].expandedY = false
		vicii.cycle = cycles[i]
	}
	for c := 0; c < PalScreenWidth*PalScreenHeight/4; c++ {
		vicii.Clock()
	}
	f, _ := os.Create("sprite.png")
	png.Encode(f, img)
}
