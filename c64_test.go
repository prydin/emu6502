package main

import (
	vic_ii "github.com/prydin/emu6502/vic-ii"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestCommodore64_Boot(t *testing.T) {
	c64 := Commodore64{}
	img := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{403, 312}})
	c64.Init(&vic_ii.ImageRaster{img}, vic_ii.PALDimensions)
	//c64.cpu.Trace = true
	c64.cpu.Reset()

	// Run 10,000,000 clock cycles
	for i := 0; i < 10000000; i++ {
		c64.Clock()
	}
	f, _ := os.Create("basic.png")
	png.Encode(f, img)
}
