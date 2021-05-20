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

import "image/color"

var C64Colors = []color.Color{
	color.NRGBA{0, 0, 0, 255},          // Black
	color.NRGBA{0xff, 0xff, 0xff, 255}, // White
	color.NRGBA{0x68, 0x37, 0x2b, 255}, // Red
	color.NRGBA{0x70, 0xa4, 0xb2, 255}, // Cyan
	color.NRGBA{0x6f, 0x3d, 0x86, 255}, // Purple
	color.NRGBA{0x58, 0x8d, 0x43, 255}, // Green
	color.NRGBA{0x35, 0x28, 0x79, 255}, // Blue
	color.NRGBA{0xb8, 0xc7, 0x6f, 255}, // Yellow
	color.NRGBA{0x6f, 0x4f, 0x25, 255}, // Orange
	color.NRGBA{0x43, 0x39, 0x00, 255}, // Brown
	color.NRGBA{0x9a, 0x67, 0x59, 255}, // Light red
	color.NRGBA{0x44, 0x44, 0x44, 255}, // Dark grey
	color.NRGBA{0x6c, 0x6c, 0x6c, 255}, // Grey 2
	color.NRGBA{0x9a, 0xd2, 0x84, 255}, // Light green
	color.NRGBA{0x6c, 0x5e, 0xb5, 255}, // Light blue
	color.NRGBA{0x95, 0x95, 0x95, 255}, // Light grey
}
