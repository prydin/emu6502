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

package main

import (
	"flag"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/prydin/emu6502/screen"
	vic_ii "github.com/prydin/emu6502/vic-ii"
	"image"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	pixelgl.Run(func() {
		c64 := Commodore64{}
		cfg := pixelgl.WindowConfig{
			Title:  "Gommodore64",
			Bounds: pixel.R(0, 0, 1024, 768),
			VSync:  true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		scr := screen.New(win, image.Rectangle{
			Min: image.Point{},
			Max: image.Point{vic_ii.PalVisibleWidth, vic_ii.PalVisibleHeight},
		})
		c64.Init(scr, vic_ii.PALDimensions)
		//c64.cpu.Trace = true
		c64.cpu.Reset()

		n := 0
		for {
			c64.Clock()
			if n % 1000000 == 0 {
				if win.Closed() {
					break
				}
			}
			n++
		}
	})
}
