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

package screen

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"time"
)

type Screen struct {
	window  *pixelgl.Window
	front   *pixel.PictureData
	back    *pixel.PictureData
	lastFps time.Time
	frames  int
	fps     int64
	jobs    chan *pixel.PictureData
	atlas   *text.Atlas
}

func New(win *pixelgl.Window, bounds image.Rectangle) *Screen {
	s := &Screen{
		window:  win,
		front:   pixel.MakePictureData(toRect(bounds)),
		back:    pixel.MakePictureData(toRect(bounds)),
		lastFps: time.Now(),
		jobs:    make(chan *pixel.PictureData),
		atlas:   text.NewAtlas(basicfont.Face7x13, text.ASCII),
	}
	go s.runScreenUpdate()
	return s
}

func (s *Screen) runScreenUpdate() {
	for {
		pd := <-s.jobs
		// draw.NearestNeighbor.Scale(screenImage, screenImage.Bounds(), image, image.Bounds(), draw.Src, nil)
		s.window.Canvas().Clear(color.Black)
		sprite := pixel.NewSprite(pd, pd.Bounds())
		sprite.Draw(s.window.Canvas(), pixel.IM.Moved(s.window.Bounds().Center()).Scaled(s.window.Bounds().Center(), 2))

		//s.window.Canvas().SetPixels(screenImage.Pix)
		txt := text.New(pixel.V(10, 10), s.atlas)
		fmt.Fprintf(txt, "FPS: %d", s.fps)
		txt.Draw(s.window.Canvas(), pixel.IM)
		s.window.Update()
	}
}

func (s *Screen) Flip() {
	s.frames++
	now := time.Now()
	if now.Sub(s.lastFps) > 1*time.Second {
		s.fps = int64(s.frames*1e9) / now.Sub(s.lastFps).Nanoseconds()
		s.lastFps = now
		s.frames = 0
	}
	tmp := s.front

	// Swap internal buffers so background task can work on rendering while we're building the next frame
	s.front = s.back
	s.back = tmp
	select {
	case s.jobs <- s.front:
	// All is well
	default:
		fmt.Println("Discarded frame")
	}
}

func (s *Screen) SetPixel(x, y uint16, color color.RGBA) {
	//s.back.SetRGBA(int(x), int(y), color)
	s.back.Pix[int(x)+int(y)*s.back.Stride] = color
}

func toRectangle(rect pixel.Rect) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(rect.Min.X), Y: int(rect.Min.Y)},
		Max: image.Point{X: int(rect.Max.X), Y: int(rect.Max.Y)},
	}
}

func toRect(r image.Rectangle) pixel.Rect {
	return pixel.R(float64(r.Min.X), float64(r.Min.Y), float64(r.Max.X), float64(r.Max.Y))
}
