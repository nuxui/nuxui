// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

package nux

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#cgo pkg-config: cairo
#cgo pkg-config: pango
#cgo pkg-config: pangocairo
#cgo pkg-config: gobject-2.0
#cgo pkg-config: libjpeg

#include <cairo/cairo.h>
#include <cairo/cairo-quartz.h>

#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>

#include <AppKit/NSGraphicsContext.h>
#import <Cocoa/Cocoa.h>

cairo_surface_t * nux_cairo_quartz_surface_create_for_cg_context(uintptr_t cgContext, unsigned int width, unsigned int height){
    return cairo_quartz_surface_create_for_cg_context(((NSGraphicsContext *)cgContext).CGContext, width, height);
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

func newSurfaceQuartzWithCGContext(context uintptr, width, height int32) *Surface {
	s := C.nux_cairo_quartz_surface_create_for_cg_context(C.uintptr_t(context), C.uint(width), C.uint(height))
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
