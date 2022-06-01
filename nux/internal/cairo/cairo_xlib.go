// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package cairo

/*
#cgo pkg-config: cairo
#cgo pkg-config: x11

#include <X11/Xlib.h>
#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>

*/
import "C"

import (
	"nuxui.org/nuxui/nux/internal/linux/xlib"
	"unsafe"
)

func XlibSurfaceCreate(display *xlib.Display, drawable xlib.Drawable, visual *xlib.Visual, width, height int32) *Surface {
	return (*Surface)(C.cairo_xlib_surface_create((*C.Display)(display), C.Drawable(drawable), (*C.Visual)(unsafe.Pointer(visual)), C.int(width), C.int(height)))
}
