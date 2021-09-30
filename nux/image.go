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

#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
#include <pango/pangocairo.h>

#include <stdlib.h>

#cgo pkg-config: libjpeg
#include "image_cairo_jpg.h"

*/
import "C"
import (

	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"

	"path/filepath"
	"strings"
	"unicode"
	"unsafe"
)

type Image interface {
	Width() int32
	Height() int32
	Buffer() *C.cairo_surface_t
}

func CreateImage(src string) Image {
	ext := strings.ToLowerSpecial(unicode.TurkishCase, filepath.Ext(src))
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))

	switch ext {
	case ".png":
		return &cimage{image: C.cairo_image_surface_create_from_png(csrc)}
	case ".svg":
		img := &cimage{image: C.cairo_svg_surface_create(csrc, 400, 400)}
		C.cairo_svg_surface_set_document_unit(img.image, C.CAIRO_SVG_UNIT_PX)
		return img
	case ".jpg", ".jpeg":
		return &cimage{image: C.cairo_image_surface_create_from_jpeg(csrc)}
	}
	return nil
}

type cimage struct {
	image *C.cairo_surface_t
}

func (me *cimage) Width() int32 {
	return int32(C.cairo_image_surface_get_width(me.image))
}

func (me *cimage) Height() int32 {
	return int32(C.cairo_image_surface_get_height(me.image))
}

func (me *cimage) Buffer() *C.cairo_surface_t {
	return me.image
}
