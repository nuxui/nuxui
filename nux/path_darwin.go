// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
	"runtime"
)

type path struct {
	ptr  darwin.CGMutablePathRef
	curX float32
	curY float32
}

func newPath() *path {
	me := &path{}
	me.ptr = darwin.CGPathCreateMutable()
	runtime.SetFinalizer(me, freePath)
	return me
}

func freePath(me *path) {
	darwin.CGPathRelease(darwin.CGPathRef(me.ptr))
}

func (me *path) native() *path {
	return me
}

func (me *path) Rect(x, y, width, height float32) {
	darwin.CGPathAddRect(me.ptr, nil, darwin.CGRectMake(x, y, width, height))
}

func (me *path) RoundRect(x, y, width, height, rx, ry float32) {
	darwin.CGPathAddRoundedRect(me.ptr, nil, darwin.CGRectMake(x, y, width, height), rx, ry)
}

func (me *path) Ellipse(cx, cy, rx, ry float32) {
	darwin.CGPathAddEllipseInRect(me.ptr, nil, darwin.CGRectMake(cx-rx, cy-ry, 2*rx, 2*ry))
}

func (me *path) MoveTo(x, y float32) {
	darwin.CGPathMoveToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) LineTo(x, y float32) {
	darwin.CGPathAddLineToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	darwin.CGPathAddCurveToPoint(me.ptr, nil, x1, y1, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) CurveToV(x2, y2, x3, y3 float32) {
	darwin.CGPathAddQuadCurveToPoint(me.ptr, nil, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) Close() {
	darwin.CGPathCloseSubpath(me.ptr)
}
