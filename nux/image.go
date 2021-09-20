// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

/*
#cgo pkg-config: cairo
#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <cairo/cairo-quartz.h>

#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
#include <pango/pangocairo.h>

#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include <stdio.h>
#include <AppKit/NSGraphicsContext.h>
#import <Cocoa/Cocoa.h>

#cgo pkg-config: libjpeg
#include "image_cairo_jpg.h"

typedef struct png{
	cairo_surface_t* image;
	int32_t width;
	int32_t height;
}png_t;

png_t * create_png(char* src){
	png_t* png = malloc (sizeof (png_t));
	memset(png, 0, sizeof(png_t));
	png->image = cairo_image_surface_create_from_png (src);
	png->width = cairo_image_surface_get_width (png->image);
	png->height = cairo_image_surface_get_height (png->image);
	return png;
}

void printData(unsigned char*data, int width, int height){
	int i,j=0;
	for(i=0; i!=width; i++){
		for(j=0; j!= height; j++){
			printf("%d", data[i*j]);
		}
		printf("\n");
	}
}

*/
import "C"
import (
	"bufio"
	"image"
	"image/color"
	"image/draw"

	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"
	"math/bits"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

type Image interface {
	Width() int32
	Height() int32
	Buffer() *C.cairo_surface_t
}

type myimage struct {
	width  int32
	height int32
	stride int32
	buffer *C.cairo_surface_t
}

func (me *myimage) Width() int32 {
	return me.width
}

func (me *myimage) Height() int32 {
	return me.height
}
func (me *myimage) Buffer() *C.cairo_surface_t {
	return me.buffer
}

func CreateImage(src string) Image {
	ext := strings.ToLowerSpecial(unicode.TurkishCase, filepath.Ext(src))
	switch ext {
	case ".png":
		return createPNGImage(src)
	case ".jpg", ".jpeg":
		return createJPGIMagec(src)
	}
	return nil
}

func createPNGImage(src string) Image {
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	png := (*pngImage)(unsafe.Pointer(C.create_png(csrc)))
	return png
}

func createJPGIMagec(src string) Image {
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	img := &jpgImage{
		image: C.cairo_image_surface_create_from_jpeg(csrc),
	}
	return img
}

func createJPGIMage(src string) Image {
	f, err := os.Open(src)
	if err != nil {
		log.E("nux", "open file error: %s", err.Error())
		return nil
	}

	img, str, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		log.E("nux", "image decode error: %s", err.Error())
		return nil
	}

	// rgba := image.NewRGBA(img.Bounds())
	rgba := newABGR(img.Bounds())
	log.V("nux", "decode str=%s, width=%d, height=%d Stride=%d", str, img.Bounds().Dx(), img.Bounds().Dy(), rgba.Stride)
	draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)
	// for i := 0; i != len(rgba.Pix); i++ {
	// 	fmt.Print(rgba.Pix[i])
	// }

	data := (*C.uchar)(unsafe.Pointer(&rgba.Pix[0]))
	// surface := C.cairo_image_surface_create(C.CAIRO_FORMAT_RGB24, C.int(img.Bounds().Dx()), C.int(img.Bounds().Dy()))
	surface := C.cairo_image_surface_create_for_data(data, C.CAIRO_FORMAT_ARGB32, C.int(img.Bounds().Dx()), C.int(img.Bounds().Dy()), C.int(rgba.Stride))
	log.V("nux", "decode str=%s, width=%d, height=%d Stride=%d", str, img.Bounds().Dx(), img.Bounds().Dy(), rgba.Stride)
	// mime := C.CString("image/jpeg")
	// defer C.free(unsafe.Pointer(mime))
	// C.cairo_surface_set_mime_data(surface, mime, data, C.ulong(len(rgba.Pix)), nil, nil)

	// C.printData(data, C.int(img.Bounds().Dx()), C.int(img.Bounds().Dy()))

	return &myimage{
		width:  int32(img.Bounds().Dx()),
		height: int32(img.Bounds().Dy()),
		stride: int32(rgba.Stride),
		buffer: surface,
	}
}

type jpgImage struct {
	image *C.cairo_surface_t
}

func (me *jpgImage) Width() int32 {
	return int32(C.cairo_image_surface_get_width(me.image))
}

func (me *jpgImage) Height() int32 {
	return int32(C.cairo_image_surface_get_height(me.image))
}

func (me *jpgImage) Buffer() *C.cairo_surface_t {
	return me.image
}

type pngImage struct {
	image  *C.cairo_surface_t
	width  int32
	height int32
}

func (me *pngImage) Width() int32 {
	return me.width
}

func (me *pngImage) Height() int32 {
	return me.height
}
func (me *pngImage) Buffer() *C.cairo_surface_t {
	return me.image
}

///////////////////////////////////////////////////////////

type abgr struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func abgrModel(c color.Color) color.Color {
	// if _, ok := c.(color.RGBA); ok {
	// 	return c
	// }
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(b >> 8), uint8(g >> 8), uint8(r >> 8), uint8(a >> 8)}
}
func (p *abgr) ColorModel() color.Model { return color.ModelFunc(abgrModel) }

func (p *abgr) Bounds() image.Rectangle { return p.Rect }

func (p *abgr) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *abgr) RGBA64At(x, y int) color.RGBA64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	r := uint16(s[0])
	g := uint16(s[1])
	b := uint16(s[2])
	a := uint16(s[3])
	return color.RGBA64{
		(r << 8) | r,
		(g << 8) | g,
		(b << 8) | b,
		(a << 8) | a,
	}
}

func (p *abgr) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return color.RGBA{s[0], s[1], s[2], s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *abgr) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *abgr) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.R
	s[1] = c1.G
	s[2] = c1.B
	s[3] = c1.A
}

func (p *abgr) SetRGBA64(x, y int, c color.RGBA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c.R >> 8)
	s[1] = uint8(c.G >> 8)
	s[2] = uint8(c.B >> 8)
	s[3] = uint8(c.A >> 8)
}

func (p *abgr) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *abgr) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &abgr{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &abgr{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *abgr) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// newABGR returns a new abgr image with the given bounds.
func newABGR(r image.Rectangle) *abgr {
	return &abgr{
		Pix:    make([]uint8, pixelBufferLength(4, r, "abgr")),
		Stride: 4 * r.Dx(),
		Rect:   r,
	}
}

func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}
