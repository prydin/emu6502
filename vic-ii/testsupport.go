package vic_ii

import (
	"image"
	"image/color"
)

type ImageRaster struct {
	Img *image.RGBA
}

func (i *ImageRaster) setPixel(x, y uint16, color color.Color) {
	i.Img.Set(int(x), int(y),color)
}
