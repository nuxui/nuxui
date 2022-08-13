// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
)

type path struct {
	ref  android.Path
	curX float32
	curY float32
}

func newPath() *path {
	me := &path{}
	runtime.SetFinalizer(me, freePath)
	return me
}

func freePath(me *path) {
}

func (me *path) native() *path {
	return me
}

func (me *path) Rect(x, y, width, height float32) {
}

func (me *path) RoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32) {
}

func (me *path) Ellipse(cx, cy, rx, ry float32) {
}

func (me *path) MoveTo(x, y float32) {
	// darwin.CGPathMoveToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) LineTo(x, y float32) {
	// darwin.CGPathAddLineToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	// darwin.CGPathAddCurveToPoint(me.ptr, nil, x1, y1, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) CurveToV(x2, y2, x3, y3 float32) {
	// darwin.CGPathAddQuadCurveToPoint(me.ptr, nil, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) Close() {
	// darwin.CGPathCloseSubpath(me.ptr)
}
