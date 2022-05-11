// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && cairo

package nux

/*
#cgo pkg-config: libjpeg

#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <stdlib.h>

#include "image_jpg_cairo.h"

*/
import "C"

import (
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unsafe"
)

func CreateImage(src string) Image {
	path, _ = filepath.Abs(path)
	ext := strings.ToLowerSpecial(unicode.TurkishCase, filepath.Ext(path))
	csrc := C.CString(path)
	defer C.free(unsafe.Pointer(csrc))

	var img *nativeImage

	switch ext {
	case ".png":
		img = &nativeImage{ptr: C.cairo_image_surface_create_from_png(csrc)}
	case ".svg":
		img = &nativeImage{ptr: C.cairo_svg_surface_create(csrc, 400, 400)}
		C.cairo_svg_surface_set_document_unit(img.ptr, C.CAIRO_SVG_UNIT_PX)
		return img
	case ".jpg", ".jpeg":
		img = &nativeImage{ptr: C.cairo_image_surface_create_from_jpeg(csrc)}
	}

	runtime.SetFinalizer(img, freeImage)
	return img
}

func freeImage(img *nativeImage) {
	C.cairo_surface_destroy(img.ptr)
}

type nativeImage struct {
	ptr *C.cairo_surface_t
}

func (me *nativeImage) Size() (width, height int32) {
	if me.ptr == nil {
		return 0, 0
	}
	return int32(C.cairo_image_surface_get_width(me.ptr)), int32(C.cairo_image_surface_get_height(me.ptr))
}
