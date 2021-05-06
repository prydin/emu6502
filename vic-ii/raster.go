package vic_ii

import "image/color"

type Raster interface {
	setPixel(x, y uint16, color color.Color)
}
