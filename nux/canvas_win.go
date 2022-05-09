// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && !cairo

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
	win32.GdipTranslateWorldTransform(me.ptr, x, y, win32.MatrixOrderPrepend)
}

func (me *canvas) Scale(x, y float32) {
	win32.GdipScaleWorldTransform(me.ptr, x, y, win32.MatrixOrderPrepend)
}

func (me *canvas) Rotate(angle float32) {
	win32.GdipRotateWorldTransform(me.ptr, angle, win32.MatrixOrderPrepend)
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
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// TODO::

	win32.GdipSetClipRect(me.ptr, x, y, width, height, win32.CombineModeReplace)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}
func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	if paint.Style()&PaintStyle_Stroke == PaintStyle_Stroke {
		win32.GdipSetPenColor(me.pen, win32.ARGB(paint.Color()))
		win32.GdipSetPenWidth(me.pen, paint.Width())
		win32.GdipDrawRectangle(me.ptr, me.pen, x, y, width, height)
	}
	if paint.Style()&PaintStyle_Fill == PaintStyle_Fill {
		win32.GdipSetSolidFillColor(me.brush, win32.ARGB(paint.Color()))
		win32.GdipFillRectangle(me.ptr, me.brush, x, y, width, height)
	}
}

// typedef enum SmoothingMode {
// 	SmoothingModeInvalid,
// 	SmoothingModeDefault,
// 	SmoothingModeHighSpeed,
// 	SmoothingModeHighQuality,
// 	SmoothingModeNone,
// 	SmoothingModeAntiAlias,
// 	SmoothingModeAntiAlias8x4,
// 	SmoothingModeAntiAlias8x8
//   } ;

func (me *canvas) DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint) {
	// if zero, path can not close
	if rLT <= 0 {
		rLT = 1
	}
	if rRT <= 0 {
		rRT = 1
	}
	if rRB <= 0 {
		rRB = 1
	}
	if rLB <= 0 {
		rLB = 1
	}

	var path *win32.GpPath
	var lastSmoothingMode win32.GpSmoothingMode
	win32.GdipGetSmoothingMode(me.ptr, &lastSmoothingMode)
	win32.GdipSetSmoothingMode(me.ptr, win32.SmoothingModeAntiAlias)

	win32.GdipCreatePath(win32.FillModeAlternate, &path)
	win32.GdipAddPathArc(path, x+width-rRT-rRT, y, rRT+rRT, rRT+rRT, -90, 90)
	win32.GdipAddPathArc(path, x+width-rRB-rRB, y+height-rRB-rRB, rRB+rRB, rRB+rRB, 0, 90)
	win32.GdipAddPathArc(path, x, y+height-rLB-rLB, rLB+rLB, rLB+rLB, 90, 90)
	win32.GdipAddPathArc(path, x, y, rLT+rLT, rLT+rLT, 180, 90)
	win32.GdipClosePathFigure(path)

	if paint.Style()&PaintStyle_Stroke == PaintStyle_Stroke {
		win32.GdipSetPenColor(me.pen, win32.ARGB(paint.Color()))
		win32.GdipSetPenWidth(me.pen, paint.Width())
		win32.GdipDrawPath(me.ptr, me.pen, path)
	}
	if paint.Style()&PaintStyle_Fill == PaintStyle_Fill {
		win32.GdipSetSolidFillColor(me.brush, win32.ARGB(paint.Color()))
		win32.GdipFillPath(me.ptr, me.brush, path)
	}

	win32.GdipDeletePath(path)
	win32.GdipSetSmoothingMode(me.ptr, lastSmoothingMode)
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

	// var mode win32.GpSmoothingMode
	// win32.GdipGetSmoothingMode(me.ptr, &mode)
	// win32.GdipSetSmoothingMode(me.ptr, win32.SmoothingModeAntiAlias)
	win32.GdipDrawString(me.ptr, &str[0], int32(len(str)), font, layout, nil, brush)
	// win32.GdipSetSmoothingMode(me.ptr, mode)
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
