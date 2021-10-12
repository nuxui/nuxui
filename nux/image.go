// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

/*
#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <stdlib.h>


#include "image_cairo_jpg.h"

*/
import "C"
import (

	// _ "image/gif"
	// _ "image/jpeg"
	// _ "image/png"

	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

// TODO:: release memory
type Image interface {
	Width() int32
	Height() int32
	Buffer() *C.cairo_surface_t
}

func CreateImage(src string) Image {
	ext := strings.ToLowerSpecial(unicode.TurkishCase, filepath.Ext(src))
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))

	var img *cimage

	switch ext {
	case ".png":
		img = &cimage{image: C.cairo_image_surface_create_from_png(csrc)}
	case ".svg":
		img = &cimage{image: C.cairo_svg_surface_create(csrc, 400, 400)}
		C.cairo_svg_surface_set_document_unit(img.image, C.CAIRO_SVG_UNIT_PX)
		return img
	case ".jpg", ".jpeg":
		img = &cimage{image: C.cairo_image_surface_create_from_jpeg(csrc)}
	}

	runtime.SetFinalizer(img, freeCairoImage)
	return img
}

func freeCairoImage(img *cimage) {
	log.V("nux", "runtime finalizer => freeCairoImage")
	C.cairo_surface_destroy(img.image)
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
