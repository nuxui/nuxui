// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !skia
// +build android

package nux

/*
// #cgo pkg-config: cairo

#cgo LDFLAGS: -L~/Documents/skia/engine/out/arm64 -lcairo -lm -lz
#cgo CFLAGS: -I/usr/local/Cellar/cairo/1.16.0_3/include
#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>


#include <stdlib.h>
#include <string.h>
#include <stdint.h>

#include <stdio.h>

void drawImage2(cairo_t* cr, cairo_surface_t *image){
	cairo_set_source_surface (cr, image, 0, 0);
	cairo_paint (cr);
}
*/
import "C"

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                        Surface                          ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////

// Surface c
type Surface struct {
	surface *C.cairo_surface_t
	canvas  *canvas
}

func (me *Surface) WriteToPng(fileName string) {
	name := C.CString(fileName)
	defer C.free(unsafe.Pointer(name))
	C.cairo_surface_write_to_png(me.surface, name)
}

func (me *Surface) GetCanvas() Canvas {
	if me.canvas == nil {
		log.Fatal("nuxui", "create surface by call NewSurface...")
	}
	return me.canvas
}
