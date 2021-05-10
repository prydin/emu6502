package vic_ii

import (
	"fmt"
	"github.com/prydin/emu6502/charset"
	"github.com/prydin/emu6502/core"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"
)

type ImageRaster struct {
	img *image.RGBA
}

func (i *ImageRaster) setPixel(x, y uint16, color color.Color) {
	i.img.Set(int(x), int(y),color)
}

func initVicII() (*VicII, *image.RGBA) {
	img := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{403, 312}})
	vicii := VicII{}
	vicii.Init()
	vicii.dimensions = PALDimensions
	vicii.bus = &core.Bus{}
	vicii.screen = &ImageRaster{img}
	vicii.borderCol = 14
	return &vicii, img
}

func Test_BlankScreen(t *testing.T) {
	vicii, img := initVicII()
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
	vicii, img := initVicII()
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xd7ff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i % 10 + 0x30)
	}
	vicii.bus.Connect(&core.RAM{ Bytes: screenMem[:]}, 0x0400, 0x07ff)
	colorMem := make([]uint8, 1024)
	for i := range colorMem {
		colorMem[i] = 14
	}
	vicii.bus.Connect(&core.RAM{Bytes: colorMem[:]}, 0xd800, 0xdbff)
	vicii.screenMemPtr = 0x0400
	vicii.charSetPtr = 0xd000
	vicii.backgroundColors[0] = 6
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
	vicii, img := initVicII()
	vicii.bus.Connect(&charset.CharacterROM, 0xd000, 0xd7ff)
	screenMem := make([]uint8, 1024)
	for i := range screenMem {
		screenMem[i] = uint8(i)
	}
	vicii.bus.Connect(&core.RAM{ Bytes: screenMem[:]}, 0x0400, 0x07ff)
	colorMem := make([]uint8, 1024)
	for i := range colorMem {
		colorMem[i] = 14
	}
	vicii.bus.Connect(&core.RAM{Bytes: colorMem[:]}, 0xd800, 0xdbff)
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
