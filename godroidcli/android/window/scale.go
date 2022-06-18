// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"image"
	"image/color"
	"math"
)

// Interp interpolates an image's color at fractional co-ordinates.
type Interp interface {
	// Interp interpolates (x, y).
	Interp(src image.Image, x, y float64) color.Color
}

// RGBA is a fast-path interpolation implementation for image.RGBA.
// It is common for an Interp to also implement RGBA.
type RGBA interface {
	// RGBA interpolates (x, y).
	RGBA(src *image.RGBA, x, y float64) color.RGBA
}

var Bilinear Interp = bilinear{}

type bilinear struct{}

func (i bilinear) Interp(src image.Image, x, y float64) color.Color {
	if src, ok := src.(*image.RGBA); ok {
		return i.RGBA(src, x, y)
	}
	return bilinearGeneral(src, x, y)
}

func bilinearGeneral(src image.Image, x, y float64) color.Color {
	p := findLinearSrc(src.Bounds(), x, y)
	var fr, fg, fb, fa float64
	var r, g, b, a uint32

	r, g, b, a = src.At(p.low.X, p.low.Y).RGBA()
	fr += float64(r) * p.frac00
	fg += float64(g) * p.frac00
	fb += float64(b) * p.frac00
	fa += float64(a) * p.frac00

	r, g, b, a = src.At(p.high.X, p.low.Y).RGBA()
	fr += float64(r) * p.frac01
	fg += float64(g) * p.frac01
	fb += float64(b) * p.frac01
	fa += float64(a) * p.frac01

	r, g, b, a = src.At(p.low.X, p.high.Y).RGBA()
	fr += float64(r) * p.frac10
	fg += float64(g) * p.frac10
	fb += float64(b) * p.frac10
	fa += float64(a) * p.frac10

	r, g, b, a = src.At(p.high.X, p.high.Y).RGBA()
	fr += float64(r) * p.frac11
	fg += float64(g) * p.frac11
	fb += float64(b) * p.frac11
	fa += float64(a) * p.frac11

	var c color.RGBA64
	c.R = uint16(fr + 0.5)
	c.G = uint16(fg + 0.5)
	c.B = uint16(fb + 0.5)
	c.A = uint16(fa + 0.5)
	return c
}

func (bilinear) RGBA(src *image.RGBA, x, y float64) color.RGBA {
	p := findLinearSrc(src.Bounds(), x, y)

	// Array offsets for the surrounding pixels.
	off00 := offRGBA(src, p.low.X, p.low.Y)
	off01 := offRGBA(src, p.high.X, p.low.Y)
	off10 := offRGBA(src, p.low.X, p.high.Y)
	off11 := offRGBA(src, p.high.X, p.high.Y)

	var fr, fg, fb, fa float64

	fr += float64(src.Pix[off00+0]) * p.frac00
	fg += float64(src.Pix[off00+1]) * p.frac00
	fb += float64(src.Pix[off00+2]) * p.frac00
	fa += float64(src.Pix[off00+3]) * p.frac00

	fr += float64(src.Pix[off01+0]) * p.frac01
	fg += float64(src.Pix[off01+1]) * p.frac01
	fb += float64(src.Pix[off01+2]) * p.frac01
	fa += float64(src.Pix[off01+3]) * p.frac01

	fr += float64(src.Pix[off10+0]) * p.frac10
	fg += float64(src.Pix[off10+1]) * p.frac10
	fb += float64(src.Pix[off10+2]) * p.frac10
	fa += float64(src.Pix[off10+3]) * p.frac10

	fr += float64(src.Pix[off11+0]) * p.frac11
	fg += float64(src.Pix[off11+1]) * p.frac11
	fb += float64(src.Pix[off11+2]) * p.frac11
	fa += float64(src.Pix[off11+3]) * p.frac11

	var c color.RGBA
	c.R = uint8(fr + 0.5)
	c.G = uint8(fg + 0.5)
	c.B = uint8(fb + 0.5)
	c.A = uint8(fa + 0.5)
	return c
}

type bilinearSrc struct {
	// Top-left and bottom-right interpolation sources
	low, high image.Point
	// Fraction of each pixel to take. The 0 suffix indicates
	// top/left, and the 1 suffix indicates bottom/right.
	frac00, frac01, frac10, frac11 float64
}

func findLinearSrc(b image.Rectangle, sx, sy float64) bilinearSrc {
	maxX := float64(b.Max.X)
	maxY := float64(b.Max.Y)
	minX := float64(b.Min.X)
	minY := float64(b.Min.Y)
	lowX := math.Floor(sx - 0.5)
	lowY := math.Floor(sy - 0.5)
	if lowX < minX {
		lowX = minX
	}
	if lowY < minY {
		lowY = minY
	}

	highX := math.Ceil(sx - 0.5)
	highY := math.Ceil(sy - 0.5)
	if highX >= maxX {
		highX = maxX - 1
	}
	if highY >= maxY {
		highY = maxY - 1
	}

	// In the variables below, the 0 suffix indicates top/left, and the
	// 1 suffix indicates bottom/right.

	// Center of each surrounding pixel.
	x00 := lowX + 0.5
	y00 := lowY + 0.5
	x01 := highX + 0.5
	y01 := lowY + 0.5
	x10 := lowX + 0.5
	y10 := highY + 0.5
	x11 := highX + 0.5
	y11 := highY + 0.5

	p := bilinearSrc{
		low:  image.Pt(int(lowX), int(lowY)),
		high: image.Pt(int(highX), int(highY)),
	}

	// Literally, edge cases. If we are close enough to the edge of
	// the image, curtail the interpolation sources.
	if lowX == highX && lowY == highY {
		p.frac00 = 1.0
	} else if sy-minY <= 0.5 && sx-minX <= 0.5 {
		p.frac00 = 1.0
	} else if maxY-sy <= 0.5 && maxX-sx <= 0.5 {
		p.frac11 = 1.0
	} else if sy-minY <= 0.5 || lowY == highY {
		p.frac00 = x01 - sx
		p.frac01 = sx - x00
	} else if sx-minX <= 0.5 || lowX == highX {
		p.frac00 = y10 - sy
		p.frac10 = sy - y00
	} else if maxY-sy <= 0.5 {
		p.frac10 = x11 - sx
		p.frac11 = sx - x10
	} else if maxX-sx <= 0.5 {
		p.frac01 = y11 - sy
		p.frac11 = sy - y01
	} else {
		p.frac00 = (x01 - sx) * (y10 - sy)
		p.frac01 = (sx - x00) * (y11 - sy)
		p.frac10 = (x11 - sx) * (sy - y00)
		p.frac11 = (sx - x10) * (sy - y01)
	}

	return p
}

// TODO(crawshaw): When we have inlining, consider func (p *RGBA) Off(x, y) int
func offRGBA(src *image.RGBA, x, y int) int {
	return (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4
}

// I is the identity Affine transform matrix.
var I = Affine{
	1, 0, 0,
	0, 1, 0,
	0, 0, 1,
}

// Affine is a 3x3 2D affine transform matrix.
// M(i,j) is Affine[i*3+j].
type Affine [9]float64

// Mul returns the multiplication of two affine transform matrices.
func (a Affine) Mul(b Affine) Affine {
	return Affine{
		a[0]*b[0] + a[1]*b[3] + a[2]*b[6],
		a[0]*b[1] + a[1]*b[4] + a[2]*b[7],
		a[0]*b[2] + a[1]*b[5] + a[2]*b[8],
		a[3]*b[0] + a[4]*b[3] + a[5]*b[6],
		a[3]*b[1] + a[4]*b[4] + a[5]*b[7],
		a[3]*b[2] + a[4]*b[5] + a[5]*b[8],
		a[6]*b[0] + a[7]*b[3] + a[8]*b[6],
		a[6]*b[1] + a[7]*b[4] + a[8]*b[7],
		a[6]*b[2] + a[7]*b[5] + a[8]*b[8],
	}
}

func (a Affine) Transform(width, height int, src image.Image, i Interp,
	f func(int, int, color.Color)) {
	srcb := src.Bounds()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sx, sy := a.pt(x, y)
			if inBounds(srcb, sx, sy) {
				f(x, y, i.Interp(src, sx, sy))
			}
		}
	}
}

func inBounds(b image.Rectangle, x, y float64) bool {
	if x < float64(b.Min.X) || x >= float64(b.Max.X) {
		return false
	}
	if y < float64(b.Min.Y) || y >= float64(b.Max.Y) {
		return false
	}
	return true
}

func (a Affine) pt(x0, y0 int) (x1, y1 float64) {
	fx := float64(x0) + 0.5
	fy := float64(y0) + 0.5
	x1 = fx*a[0] + fy*a[1] + a[2]
	y1 = fx*a[3] + fy*a[4] + a[5]
	return x1, y1
}

// Scale produces a scaling transform of factors x and y.
func (a Affine) Scale(x, y float64) Affine {
	return a.Mul(Affine{
		1 / x, 0, 0,
		0, 1 / y, 0,
		0, 0, 1,
	})
}
