// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/ios"
	"runtime"
)

type path struct {
	ptr  ios.CGMutablePathRef
	curX float32
	curY float32
}

func newPath() *path {
	me := &path{}
	me.ptr = ios.CGPathCreateMutable()
	runtime.SetFinalizer(me, freePath)
	return me
}

func freePath(me *path) {
	ios.CGPathRelease(ios.CGPathRef(me.ptr))
}

func (me *path) native() *path {
	return me
}

func (me *path) Rect(x, y, width, height float32) {
	ios.CGPathAddRect(me.ptr, nil, ios.CGRectMake(x, y, width, height))
}

func (me *path) RoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32) {
}

func (me *path) Ellipse(cx, cy, rx, ry float32) {
	ios.CGPathAddEllipseInRect(me.ptr, nil, ios.CGRectMake(cx-rx, cy-ry, 2*rx, 2*ry))
}

func (me *path) MoveTo(x, y float32) {
	ios.CGPathMoveToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) LineTo(x, y float32) {
	ios.CGPathAddLineToPoint(me.ptr, nil, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	ios.CGPathAddCurveToPoint(me.ptr, nil, x1, y1, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) CurveToV(x2, y2, x3, y3 float32) {
	ios.CGPathAddQuadCurveToPoint(me.ptr, nil, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) Close() {
	ios.CGPathCloseSubpath(me.ptr)
}
