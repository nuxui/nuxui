// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && cairo

package cairo

/*
#cgo LDFLAGS: -limm32
#cgo pkg-config: cairo

#include <cairo/cairo.h>
#include <cairo/cairo-win32.h>

#include <windows.h>
#include <windowsx.h>

HDC toHDC(uintptr_t hdc){
	return (HDC)hdc;
}
*/
import "C"

func Win32SurfaceCreate(hdc uintptr) *Surface {
	return (*Surface)(C.cairo_win32_surface_create(C.toHDC(C.uintptr_t(hdc))))
}
