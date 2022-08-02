// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/ios"
)

type canvas struct {
	ctx ios.CGContextRef
}

func newCanvas(ctx ios.CGContextRef) *canvas {
	return &canvas{
		ctx: ctx,
	}
}

func (me *canvas) native() *canvas {
	return me
}

func (me *canvas) ResetClip() {
	ios.CGContextResetClip(me.ctx)
}

func (me *canvas) Save() {
	ios.CGContextSaveGState(me.ctx)
}

func (me *canvas) Restore() {
	ios.CGContextRestoreGState(me.ctx)
}

func (me *canvas) Translate(x, y float32) {
	ios.CGContextTranslateCTM(me.ctx, x, y)
}

func (me *canvas) Scale(x, y float32) {
	ios.CGContextScaleCTM(me.ctx, x, y)
}

func (me *canvas) Rotate(angle float32) {
	ios.CGContextRotateCTM(me.ctx, angle)
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
	ios.CGContextClipToRect(me.ctx, ios.CGRectMake(x, y, width, height))
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	rect := ios.CGRectMake(x, y, width, height)
	path := ios.CGPathCreateWithRoundedRect(rect, cornerX, cornerY, nil)
	ios.CGContextAddPath(me.ctx, path)
	ios.CGContextClip(me.ctx)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
	ios.CGContextSetAlpha(me.ctx, alpha)
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	r, g, b, a := paint.Color().RGBAf()
	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		x += 1
		y += 1
	}

	if fix {
		me.Save()
		me.Translate(-0.5, -0.5)
	}

	rect := ios.CGRectMake(x, y, width, height)

	switch paint.Style() {
	case PaintStyle_Fill:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillRect(me.ctx, rect)
	case PaintStyle_Stroke:
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokeRectWithWidth(me.ctx, rect, paint.Width())
	case PaintStyle_Both:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillRect(me.ctx, rect)
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokeRectWithWidth(me.ctx, rect, paint.Width())
	}

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawRoundRect(x, y, width, height, rLT, rRT, rRB, rLB float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}

	r, g, b, a := paint.Color().RGBAf()
	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		x += 1
		y += 1
	}

	if fix {
		me.Save()
		me.Translate(-0.5, -0.5)
	}

	path := ios.CGPathCreateMutable()
	ios.CGPathAddRoundRectPath(path, x, y, width, height, rLT, rRT, rRB, rLB)
	ios.CGPathCloseSubpath(path)
	ios.CGContextAddPath(me.ctx, ios.CGPathRef(path))
	ios.CGContextSetLineWidth(me.ctx, paint.Width())

	hasShadow := false
	if sc, sx, sy, sb := paint.Shadow(); sc != 0 && sb > 0 {
		hasShadow = true
		me.Save()
		r0, g0, b0, a0 := sc.RGBAf()
		ios.CGContextSetShadowWithColor(me.ctx, ios.CGSizeMake(sx, sy), sb, ios.CGColorMake(r0, g0, b0, a0))
	}

	ios.CGContextSetLineDash(me.ctx, 0, paint.Dash(), len(paint.Dash()))

	switch paint.Style() {
	case PaintStyle_Fill:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
	case PaintStyle_Stroke:
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	case PaintStyle_Both:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	}
	ios.CGPathRelease(ios.CGPathRef(path))

	if hasShadow {
		me.Restore()
	}

	if fix {
		me.Restore()
	}
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {
	// TODO:: useCenter
	r, g, b, a := paint.Color().RGBAf()
	ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
	ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)

	var clockwise int = 0
	if useCenter {
		clockwise = 1
	}
	ios.CGContextAddArc(me.ctx, x, y, radius, startAngle, endAngle, clockwise)

	switch paint.Style() {
	case PaintStyle_Fill:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
	case PaintStyle_Stroke:
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	case PaintStyle_Both:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	}
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
	ios.CGContextFillEllipseInRect(me.ctx, ios.CGRectMake(x, y, width, height))
}

func (me *canvas) DrawPath(path Path, paint Paint) {
	ios.CGContextAddPath(me.ctx, ios.CGPathRef(path.native().ptr))
	r, g, b, a := paint.Color().RGBAf()
	switch paint.Style() {
	case PaintStyle_Fill:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
	case PaintStyle_Stroke:
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	case PaintStyle_Both:
		ios.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		ios.CGContextFillPath(me.ctx)
		ios.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		ios.CGContextStrokePath(me.ctx)
	}
}

func (me *canvas) DrawImage(img Image) {
	img.Draw(me)
}

func (me *canvas) Flush() {
	// ios.CGContextFlush(me.ctx)
}

func (me *canvas) Destroy() {
}
