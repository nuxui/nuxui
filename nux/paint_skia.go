// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build skia

package nux

type paint struct {
	ref *C.sk_paint_t
}

func newPaint() *paint {
	me := paint{
		ref: C.sk_paint_new(),
	}
	runtime.SetFinalizer(me, deletePaint)
	return me
}

func deletePaint(p *paint) {
	C.sk_paint_delete(p.ref)
}

func (me *paint) Color() Color {
	return Color(C.sk_paint_get_color(me.ref))
}

func (me *paint) SetColor(color Color) {
	C.sk_paint_set_color(me.ref, C.sk_color_t(color))
}

func (me *paint) AntiAlias() bool {
	if C.sk_paint_is_antialias(me.ref) > 0 {
		return true
	}
	return false
}

func (me *paint) SetAntiAlias(antialias bool) {
	C.sk_paint_set_antialias(me.ref, C.bool(antialias))
}

func (me *paint) Width() float32 {
	return float32(C.sk_paint_get_stroke_width(me.ref))
}

func (me *paint) SetWidth(width float32) {
	C.sk_paint_set_stroke_width(me.ref, C.float(width))
}
