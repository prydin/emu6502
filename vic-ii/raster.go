package vic_ii

import "image/color"

type Raster interface {
	setPixel(color color.Color)
}
