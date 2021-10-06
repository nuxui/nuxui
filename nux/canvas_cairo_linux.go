// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package nux

/*
#cgo pkg-config: cairo
#include <cairo/cairo.h>
#include <cairo/cairo-pdf.h>
#include <cairo/cairo-ps.h>
#include <cairo/cairo-svg.h>
#include <cairo/cairo-xlib.h>

#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
#include <pango/pangocairo.h>

#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>

#include <X11/Xlib.h>

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

func newSurfaceXlib(display *C.Display, drawable C.Drawable, visual *C.Visual, width, height int32) *Surface {
	s := C.cairo_xlib_surface_create(display, drawable, visual, C.int(width), C.int(height))
	cairo := C.cairo_create(s)
	return &Surface{surface: s, canvas: &canvas{ptr: cairo}}
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

func (me *Surface) Flush() {
	C.cairo_surface_flush(me.surface)
}

func (me *Surface) Destroy() {
	C.cairo_surface_finish(me.surface)
	C.cairo_surface_destroy(me.surface)
	me.surface = nil
	me.canvas = nil
}
