// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jint nux_canvas_save(jobject canvas);
void nux_canvas_restore(jobject canvas);
void nux_canvas_translate(jobject canvas, jfloat x, jfloat y);
void nux_canvas_scale(jobject canvas, jfloat x, jfloat y);
void nux_canvas_rotate(jobject canvas, jfloat degrees);
void nux_canvas_clipRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom);
void nux_canvas_drawRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint);
void nux_canvas_drawRoundRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jfloat rx, jfloat ry, jobject paint);
void nux_canvas_drawBitmap(jobject canvas, jobject bitmap, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint);
*/
import "C"

func (me Canvas) Save() {
	C.nux_canvas_save(C.jobject(me))
}

func (me Canvas) Restore() {
	C.nux_canvas_restore(C.jobject(me))
}

func (me Canvas) Translate(x, y float32) {
	C.nux_canvas_translate(C.jobject(me), C.jfloat(x), C.jfloat(y))
}

func (me Canvas) Scale(x, y float32) {
	C.nux_canvas_scale(C.jobject(me), C.jfloat(x), C.jfloat(y))
}

func (me Canvas) Rotate(degrees float32) {
	C.nux_canvas_rotate(C.jobject(me), C.jfloat(degrees))
}

func (me Canvas) ClipRect(x, y, width, height float32) {
	C.nux_canvas_clipRect(C.jobject(me), C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height))
}

func (me Canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// TODO::
}

func (me Canvas) DrawRect(x, y, width, height float32, paint Paint) {
	C.nux_canvas_drawRect(C.jobject(me), C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height), C.jobject(paint))
}

func (me Canvas) DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint) {
	// TODO:: use arc
	C.nux_canvas_drawRoundRect(C.jobject(me), C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height), C.jfloat(rLT), C.jfloat(rLT), C.jobject(paint))
}

func (me Canvas) DrawBitmap(bitmap Bitmap, x, y, width, height float32, paint Paint) {
	C.nux_canvas_drawBitmap(C.jobject(me), C.jobject(bitmap), C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height), C.jobject(paint))
}
