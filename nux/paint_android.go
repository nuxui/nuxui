// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"runtime"

	"nuxui.org/nuxui/nux/internal/android"
)

func NewPaint() Paint {
	me := &paint{
		ref: android.NewPaint(), // JObject ref
	}
	runtime.SetFinalizer(me, freePaint)
	return me
}

func freePaint(me *paint) {
	android.DeleteGlobalRef(android.JObject(me.ref))
}

type paint struct {
	ref android.Paint
}

func (me *paint) Color() Color {
	return FromARGB(me.ref.GetColor())
}

func (me *paint) SetColor(color Color) {
	me.ref.SetColor(color.ARGB())
}

func (me *paint) AntiAlias() bool {
	return false
}

func (me *paint) SetAntiAlias(antialias bool) {
}

func (me *paint) Width() float32 {
	return 0
}

func (me *paint) SetWidth(width float32) {
}

func (me *paint) Style() PaintStyle {
	style := me.ref.GetStyle()
	switch style {
	case 0:
		return PaintStyle_Fill
	case 1:
		return PaintStyle_Stroke
	case 2:
		return PaintStyle_Both
	}
	return PaintStyle_Fill
}

func (me *paint) SetStyle(style PaintStyle) {
	s := 0
	switch style {
	case PaintStyle_Fill:
		s = 0
	case PaintStyle_Stroke:
		s = 1
	case PaintStyle_Both:
		s = 2
	}
	me.ref.SetStyle(s)
}

func (me *paint) TextSize() float32 {
	return 0
}

func (me *paint) SetTextSize(size float32) {
}

func (me *paint) SetDash(dash []float32) {
}

func (me *paint) Dash() []float32 {
	return []float32{}
}

func (me *paint) SetShadow(color Color, offsetX, offsetY, blur float32) {

}

func (me *paint) Shadow() (color Color, offsetX, offsetY, blur float32) {
	return
}

func (me *paint) HasShadow() bool {
	return false
}

func (me *paint) native() *paint {
	return me
}
