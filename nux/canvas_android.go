// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && android

package nux

/*
#include <jni.h>
void deleteGlobalRef(jobject globalRef);
void deleteLocalRef(jobject localRef);

jint canvas_save(jobject canvas);
void canvas_restore(jobject canvas);
void canvas_translate(jobject canvas, jfloat x, jfloat y);
void canvas_scale(jobject canvas, jfloat x, jfloat y);
void canvas_rotate(jobject canvas, jfloat degrees);
void canvas_clipRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom);
// void canvas_drawColor(jobject canvas, uint32_t color);
void canvas_drawRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint);
void canvas_drawRoundRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jfloat rx, jfloat ry, jobject paint);
void canvas_drawText(jobject canvas, char* text, jint width, jobject paint);
void canvas_drawBitmap(jobject canvas, jobject bitmap, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint);

jobject new_Paint();
void paint_setColor(jobject paint, uint32_t color);
void paint_setTextSize(jobject paint, jfloat textSize);
void paint_setStyle(jobject paint, jint style);
void paint_setAntiAlias(jobject paint, jboolean aa);
void paint_measureText(jobject paint, char* text, jint width, jint *outWidth, jint* outHeight);
*/
import "C"
import (
// "nuxui.org/nuxui/log"
// "unicode/utf16"
)

type canvas struct {
	ptr C.jobject
}

func newCanvas(canvasPtr C.jobject) *canvas {
	me := &canvas{
		ptr: canvasPtr,
	}
	return me
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
	C.canvas_save(me.ptr)
}

func (me *canvas) Restore() {
	C.canvas_restore(me.ptr)
}

func (me *canvas) Translate(x, y float32) {
	C.canvas_translate(me.ptr, C.jfloat(x), C.jfloat(y))
}

func (me *canvas) Scale(x, y float32) {
	C.canvas_scale(me.ptr, C.jfloat(x), C.jfloat(y))
}

func (me *canvas) Rotate(degrees float32) {
	C.canvas_rotate(me.ptr, C.jfloat(degrees))
}

func (me *canvas) Skew(x, y float32) {
}

func (me *canvas) Transform(a, b, c, d, e, f float32) {
}

func (me *canvas) SetMatrix(matrix Matrix) {
}

func (me *canvas) GetMatrix() Matrix {
	return Matrix{}
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	C.canvas_clipRect(me.ptr, C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height))
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// TODO::
	me.ClipRect(x, y, width, height)
}

func (me *canvas) ClipPath(path Path) {
}

func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	p := C.new_Paint()
	C.paint_setColor(p, C.uint32_t(paint.Color().ARGB()))
	C.paint_setTextSize(p, C.jfloat(paint.TextSize()))
	C.paint_setStyle(p, C.jint(paint.Style()))
	C.paint_setAntiAlias(p, C.jboolean(1))
	C.canvas_drawRect(me.ptr, C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height), p)
	C.deleteLocalRef(p)
}

func (me *canvas) DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint) {
	p := C.new_Paint()
	C.paint_setColor(p, C.uint32_t(paint.Color().ARGB()))
	C.paint_setTextSize(p, C.jfloat(paint.TextSize()))
	C.paint_setStyle(p, C.jint(paint.Style())) // TODO:: the style is not same with android
	C.paint_setAntiAlias(p, C.jboolean(1))
	// TODO:: use arc
	C.canvas_drawRoundRect(me.ptr, C.jfloat(x), C.jfloat(y), C.jfloat(x+width), C.jfloat(y+height), C.jfloat(rLT), C.jfloat(rLT), p)
	C.deleteLocalRef(p)
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {

}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
}

func (me *canvas) DrawPath(path Path) {
}

// func (me *canvas) DrawColor(color Color) {
// 	C.canvas_drawColor(me.ptr, C.uint32_t(color))
// }

func (me *canvas) DrawImage(img Image) {
	w, h := img.Size()
	p := C.new_Paint()
	C.canvas_drawBitmap(me.ptr, img.(*nativeImage).ptr, 0, 0, C.jfloat(w), C.jfloat(h), p)
	C.deleteLocalRef(p)
}

func (me *canvas) DrawText(text string, width, height float32, paint Paint) {
	chars := ([]byte)(text)
	p := C.new_Paint()
	C.paint_setColor(p, C.uint32_t(paint.Color().ARGB()))
	C.paint_setTextSize(p, C.jfloat(paint.TextSize()))
	C.paint_setStyle(p, C.jint(paint.Style()))
	C.paint_setAntiAlias(p, C.jboolean(1))
	C.canvas_drawText(me.ptr, (*C.char)(&chars[0]), C.jint(width), C.jobject(p))
	C.deleteLocalRef(p)
}

func (me *canvas) Flush() {
}

func (me *canvas) Destroy() {

}

func (me *paint) MeasureText(text string, width, height float32) (outWidth float32, outHeight float32) {
	if text == "" {
		return 0, 0
	}

	chars := ([]byte)(text)
	p := C.new_Paint()
	C.paint_setColor(p, C.uint32_t(me.Color().ARGB()))
	C.paint_setTextSize(p, C.jfloat(me.TextSize()))
	C.paint_setStyle(p, C.jint(me.Style()))
	C.paint_setAntiAlias(p, C.jboolean(1))
	var w, h C.jint = 0, 0
	C.paint_measureText(p, (*C.char)(&chars[0]), C.jint(int32(width)), &w, &h)
	C.deleteLocalRef(p)
	return float32(w), float32(h)
}

func (me *paint) CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32 {
	return 0
}
