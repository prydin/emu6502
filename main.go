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
	"fmt"
	"github.com/beevik/go6502/asm"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/prydin/emu6502/computer"
	"github.com/prydin/emu6502/screen"
	vic_ii "github.com/prydin/emu6502/vic-ii"
	"image"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var loadasm = flag.String("loadasm", "", "load assembly language file")
var loadprg = flag.String("loadprg", "", "load binary program file")

var PalFPS = 50.125
var PalFrameTime = time.Duration((1 / PalFPS) * 1e9)

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

	var code *asm.Assembly
	var sourceMap *asm.SourceMap
	if *loadasm != "" {
		in, err := os.Open(*loadasm)
		if err != nil {
			panic(err)
		}
		code, sourceMap, err = asm.Assemble(in, *loadasm, os.Stderr, 0)
		for _, parseErr := range code.Errors {
			println(parseErr)
		}
		if err != nil {
			panic(err)
		}
	}
	if *loadprg != "" {
		prg, err := os.Open(*loadprg)
		if err != nil {
			panic(err)
		}
		defer prg.Close()
		bytes, err := ioutil.ReadAll(prg)
		if err != nil {
			panic(err)
		}
		code = &asm.Assembly{
			Code: bytes[2:],
		}
		sourceMap = asm.NewSourceMap()
		sourceMap.Size = uint32(len(bytes) - 2)
		sourceMap.Origin = uint16(bytes[0]) + uint16(bytes[1])<<8
		fmt.Printf("Loaded %d bytes starting at %d (%04x)\n", sourceMap.Size, sourceMap.Origin, sourceMap.Origin)
	}

	pixelgl.Run(func() {
		c64 := computer.Commodore64{}
		cfg := pixelgl.WindowConfig{
			Title:     "Gommodore64",
			Bounds:    pixel.R(0, 0, 1024, 768),
			VSync:     true,
			Resizable: true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		win.SetSmooth(true) // Gives a nice blurry retro look!
		scr := screen.New(win, image.Rectangle{
			Min: image.Point{},
			Max: image.Point{vic_ii.PalVisibleWidth, vic_ii.PalVisibleHeight},
		})

		c64.Cpu.CrashOnInvalidInst = true // TODO: Make configurable
		c64.Init(scr, vic_ii.PALDimensions)
		c64.Keyboard.SetProvider(win)
		//c64.cpu.Trace = true
		c64.Cpu.Reset()

		var lastVSynch time.Time
		n := 0
		for {
			if c64.Vic.IsVSynch() {
				now := time.Now()
				frameTime := now.Sub(lastVSynch)

				// Sleeping precision is too low, so we spin instead
				if frameTime < PalFrameTime {
					target := now.Add(PalFrameTime - frameTime)
					for time.Now().Before(target) {
						// Do nothing
					}
				}
				lastVSynch = time.Now()
			}
			c64.Clock()
			if code != nil && n > 10000000 {
				for i := uint16(0); i < uint16(sourceMap.Size); i++ {
					c64.Bus.WriteByte(i+sourceMap.Origin, code.Code[i])
				}
				// Set pointer to end of BASIC program. // TODO: We should probably have a flag for this
				end := uint16(sourceMap.Size + 1) + sourceMap.Origin
				c64.Bus.WriteByte(0x002d, uint8(end & 0xff))
				c64.Bus.WriteByte(0x002e, uint8(end >> 8))
				code = nil
			}
			if n%1000000 == 0 {
				if win.Closed() {
					break
				}
			}
			n++
		}
	})
}
