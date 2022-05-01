// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"nuxui.org/nuxui/nux/internal/win32"
	"syscall"
)

func newCanvas(hdcBuffer uintptr) *canvas {
	me := &canvas{
		ptr:    &win32.GpGraphics{},
		pen:    &win32.GpPen{},
		brush:  &win32.GpBrush{},
		states: []win32.GpState{},
		clip:   &win32.RectF{0, 0, 99999, 99999},
	}
	win32.GdipCreateFromHDC(hdcBuffer, &me.ptr)
	win32.GdipCreatePen1(0, 1, win32.UnitWorld, &me.pen)
	win32.GdipCreateSolidFill(0, &me.brush)
	return me
}

type canvas struct {
	// hdc    uintptr
	ptr    *win32.GpGraphics
	pen    *win32.GpPen
	brush  *win32.GpBrush
	states []win32.GpState
	clip   *win32.RectF
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
	var s win32.GpState
	win32.GdipSaveGraphics(me.ptr, &s)
	me.states = append(me.states, s)
}

func (me *canvas) Restore() {
	l := len(me.states)
	s := me.states[l-1]
	win32.GdipRestoreGraphics(me.ptr, s)
	me.states = me.states[0 : l-1]
}

func (me *canvas) Translate(x, y float32) {
	win32.GdipTranslateWorldTransform(me.ptr, x, y, win32.MatrixOrderAppend)
}

func (me *canvas) Scale(x, y float32) {
	// win32.GdipScaleWorldTransform(me.ptr, x, y, win32.MatrixOrderAppend)
}

func (me *canvas) Rotate(angle float32) {
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
	return Matrix{}
}

func (me *canvas) ClipRect(x, y, width, height float32) {
	win32.GdipSetClipRect(me.ptr, x, y, width, height, win32.CombineModeReplace)
	me.clip.X = x
	me.clip.Y = y
	me.clip.Width = width
	me.clip.Height = height
}

func (me *canvas) ClipRoundRect(x, y, width, height, radius float32) {
	// TODO::

	win32.GdipSetClipRect(me.ptr, x, y, width, height, win32.CombineModeReplace)
	me.clip.X = x
	me.clip.Y = y
	me.clip.Width = width
	me.clip.Height = height
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}
func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	switch paint.Style() {
	case PaintStyle_Stroke:
		{
			win32.GdipSetPenColor(me.pen, win32.ARGB(paint.Color()))
			win32.GdipSetPenWidth(me.pen, paint.Width())
			win32.GdipDrawRectangle(me.ptr, me.pen, x, y, width, height)
		}
	case PaintStyle_Fill:
		{
			win32.GdipSetSolidFillColor(me.brush, win32.ARGB(paint.Color()))
			win32.GdipFillRectangle(me.ptr, me.brush, x, y, width, height)
		}
	case PaintStyle_Both:
		{
			win32.GdipSetPenColor(me.pen, win32.ARGB(paint.Color()))
			win32.GdipSetPenWidth(me.pen, paint.Width())
			win32.GdipDrawRectangle(me.ptr, me.pen, x, y, width, height)

			win32.GdipSetSolidFillColor(me.brush, win32.ARGB(paint.Color()))
			win32.GdipFillRectangle(me.ptr, me.brush, x, y, width, height)
		}
	}
}

func (me *canvas) DrawRoundRect(x, y, width, height float32, radius float32, paint Paint) {
	// TODO::
	me.DrawRect(x, y, width, height, paint)
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {
	// TODO:: useCenter
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
}

func (me *canvas) DrawPath(path Path) {
	// TODO::
}

func (me *canvas) DrawImage(img Image) {
	w, h := img.Size()
	win32.GdipDrawImageRect(me.ptr, img.(*nativeImage).ptr, 0, 0, float32(w), float32(h))
}

func (me *canvas) DrawText(text string, width, height float32, paint Paint) {
	font := &win32.GpFont{}
	family := &win32.GpFontFamily{}
	win32.GdipGetGenericFontFamilyMonospace(&family)
	win32.GdipCreateFont(family, paint.TextSize(), 0, 0, &font)

	str, _ := syscall.UTF16FromString(text)
	layout := &win32.RectF{0, 0, width, height}
	brush := &win32.GpBrush{}
	win32.GdipCreateSolidFill(win32.ARGB(paint.Color()), &brush)

	var mode win32.GpSmoothingMode
	win32.GdipGetSmoothingMode(me.ptr, &mode)
	win32.GdipSetSmoothingMode(me.ptr, win32.SmoothingModeAntiAlias)
	win32.GdipDrawString(me.ptr, &str[0], int32(len(str)), font, layout, nil, brush)
	win32.GdipSetSmoothingMode(me.ptr, mode)
}

func (me *canvas) Flush() {
	win32.GdipFlush(me.ptr, win32.FlushIntentionFlush)
}

func (me *canvas) Destroy() {
	win32.GdipDeletePen(me.pen)
	win32.GdipDeleteBrush(me.brush)
	win32.GdipDeleteGraphics(me.ptr)
}

func (me *paint) MeasureText(text string, width, height float32) (outWidth float32, outHeight float32) {
	if text == "" {
		return 0, 0
	}

	// TODO:: use hdc as args for newCanvas
	hwnd := theApp.window.hwnd
	font := &win32.GpFont{}
	family := &win32.GpFontFamily{}
	win32.GdipGetGenericFontFamilyMonospace(&family)
	win32.GdipCreateFont(family, me.textSize, 0, 0, &font)

	str, _ := syscall.UTF16FromString(text)
	layout := &win32.RectF{0, 0, width, height}
	size := &win32.RectF{}
	g := &win32.GpGraphics{}
	win32.GdipCreateFromHWND(hwnd, &g)
	win32.GdipMeasureString(g, &str[0], int32(len(str)), font, layout, nil, size, nil, nil)
	win32.GdipDeleteGraphics(g)
	return size.Width, size.Height
}

func (me *paint) CharacterIndexForPoint(text string, width, height float32, x, y float32) uint32 {
	return 0
}
