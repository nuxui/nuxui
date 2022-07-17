// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || (windows && cairo)

package nux

import (
	"nuxui.org/nuxui/nux/internal/cairo"
	"runtime"
)

type path struct {
	cairo   *cairo.Cairo
	surface *cairo.Surface
	curX    float32
	curY    float32
}

func newPath() *path {
	me := &path{}
	me.surface = cairo.ImageSurfaceCreate(cairo.CAIRO_FORMAT_ARGB32, 1, 1)
	me.cairo = cairo.Create(me.surface)
	runtime.SetFinalizer(me, freePath)
	return me
}

func freePath(me *path) {
	me.cairo.Destroy()
	me.surface.Finish()
	me.surface.Destroy()
}

func (me *path) native() *path {
	return me
}

func (me *path) Rect(x, y, width, height float32) {
}

func (me *path) RoundRect(x, y, width, height, rx, ry float32) {
}

func (me *path) Ellipse(cx, cy, rx, ry float32) {
}

func (me *path) MoveTo(x, y float32) {
	me.cairo.MoveTo(x, y)
	me.curX = x
	me.curY = y
}

func (me *path) LineTo(x, y float32) {
	me.cairo.LineTo(x, y)
	me.curX = x
	me.curY = y
}

func (me *path) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	me.cairo.CurveTo(x1, y1, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) CurveToV(x2, y2, x3, y3 float32) {
	cx1 := me.curX + 2.0/3.0*(x2-me.curX)
	cy1 := me.curY + 2.0/3.0*(y2-me.curY)
	cx2 := x3 + 2.0/3.0*(x2-x3)
	cy2 := y3 + 2.0/3.0*(y2-y3)

	me.cairo.CurveTo(cx1, cy1, cx2, cy2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) Close() {
	me.cairo.ClosePath()
}
