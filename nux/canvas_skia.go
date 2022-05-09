// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build skia

package nux

/*
#cgo LDFLAGS: -L/Users/mustodo/Documents/skia/skia/out/Static
#cgo LDFLAGS: -lstdc++
#cgo LDFLAGS: -lskia
#cgo LDFLAGS: -framework Cocoa

#cgo CFLAGS: -I/Users/mustodo/Documents/skia/skia
#include "include/c/sk_types.h"
#include "include/c/sk_canvas.h"
#include "include/c/sk_data.h"
#include "include/c/sk_image.h"
#include "include/c/sk_imageinfo.h"
#include "include/c/sk_paint.h"
#include "include/c/sk_path.h"
#include "include/c/sk_surface.h"
*/
import "C"

type canvas struct {
	ptr *C.sk_canvas_t
}

func (me *canvas) Save() {
	C.sk_canvas_save(me.ptr)
}

func (me *canvas) Restore() {
	C.sk_canvas_restore(me.ptr)
}

func (me *canvas) Translate(x, y float32) {
	C.sk_canvas_translate(me.ptr, C.float(x), C.float(y))
}

func (me *canvas) Scale(x, y float32) {
	C.sk_canvas_scale(me.ptr, C.float(x), C.float(y))
}

func (me *canvas) Rotate(degrees float32) {
	C.sk_canvas_rotate_degrees(me.ptr, C.float(degrees))
}

func (me *canvas) Skew(x, y float32) {
	// TODO::
}

func (me *canvas) Transform(a, b, c, d, e, f float32) {
	// TODO::
}

func (me *canvas) SetMatrix(matrix Matrix) {
	// TODO::
}

func (me *canvas) GetMatrix() Matrix {
	// TODO::
	return nil
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	var rect C.sk_rect_t
	rect.left = C.float(x)
	rect.top = C.float(y)
	rect.right = C.float(x + width)
	rect.bottom = C.float(y + height)
	C.sk_canvas_clip_rect(me.ptr, &rect)
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// rect := C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
	// path := C.CGPathCreateWithRoundedRect(rect, C.CGFloat(cornerX), C.CGFloat(cornerY), nil)
	// C.CGContextAddPath(me.ptr, path)
	// C.CGContextClip(me.ptr)
	me.ClipRect(x, y, width, height)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
	// TODO::
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	var rect C.sk_rect_t
	rect.left = C.float(x)
	rect.top = C.float(y)
	rect.right = C.float(x + width)
	rect.bottom = C.float(y + height)
	C.sk_canvas_draw_rect(me.ptr, &rect, paint.toC())
}
