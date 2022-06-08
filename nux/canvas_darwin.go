// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
)

type canvas struct {
	ctx darwin.CGContext
}

func newCanvas(ctx darwin.CGContext) *canvas {
	return &canvas{
		ctx: ctx,
	}
}

func (me *canvas) native() *canvas {
	return me
}

func (me *canvas) ResetClip() {
	darwin.CGContextResetClip(me.ctx)
}

func (me *canvas) Save() {
	darwin.CGContextSaveGState(me.ctx)
}

func (me *canvas) Restore() {
	darwin.CGContextRestoreGState(me.ctx)
}

func (me *canvas) Translate(x, y float32) {
	darwin.CGContextTranslateCTM(me.ctx, x, y)
}

func (me *canvas) Scale(x, y float32) {
	darwin.CGContextScaleCTM(me.ctx, x, y)
}

func (me *canvas) Rotate(angle float32) {
	darwin.CGContextRotateCTM(me.ctx, angle)
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
	darwin.CGContextClipToRect(me.ctx, darwin.CGRectMake(x, y, width, height))
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	rect := darwin.CGRectMake(x, y, width, height)
	path := darwin.CGPathCreateWithRoundedRect(rect, cornerX, cornerY, nil)
	darwin.CGContextAddPath(me.ctx, path)
	darwin.CGContextClip(me.ctx)
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
	darwin.CGContextSetAlpha(me.ctx, alpha)
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

	rect := darwin.CGRectMake(x, y, width, height)

	switch paint.Style() {
	case PaintStyle_Fill:
		darwin.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		darwin.CGContextFillRect(me.ctx, rect)
	case PaintStyle_Stroke:
		darwin.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		darwin.CGContextStrokeRectWithWidth(me.ctx, rect, paint.Width())
	case PaintStyle_Both:
		darwin.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		darwin.CGContextFillRect(me.ctx, rect)
		darwin.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		darwin.CGContextStrokeRectWithWidth(me.ctx, rect, paint.Width())
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

	path := darwin.CGPathCreateMutable()
	darwin.CGPathAddRoundRectPath(path, x, y, width, height, rLT, rRT, rRB, rLB)
	darwin.CGPathCloseSubpath(path)
	darwin.CGContextAddPath(me.ctx, darwin.CGPath(path))
	darwin.CGContextSetLineWidth(me.ctx, paint.Width())

	hasShadow := false
	if sc, sx, sy, sb := paint.Shadow(); sc != 0 && sb > 0 {
		hasShadow = true
		me.Save()
		r0, g0, b0, a0 := sc.RGBAf()
		darwin.CGContextSetShadowWithColor(me.ctx, darwin.CGSizeMake(sx, -sy), sb, darwin.CGColorMake(r0, g0, b0, a0))
	}

	switch paint.Style() {
	case PaintStyle_Fill:
		darwin.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		darwin.CGContextFillPath(me.ctx)
	case PaintStyle_Stroke:
		darwin.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		darwin.CGContextStrokePath(me.ctx)
	case PaintStyle_Both:
		darwin.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
		darwin.CGContextFillPath(me.ctx)
		darwin.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)
		darwin.CGContextStrokePath(me.ctx)
	}
	darwin.CGPathRelease(darwin.CGPath(path))

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
	darwin.CGContextSetRGBFillColor(me.ctx, r, g, b, a)
	darwin.CGContextSetRGBStrokeColor(me.ctx, r, g, b, a)

	var clockwise int = 0
	if useCenter {
		clockwise = 1
	}
	darwin.CGContextAddArc(me.ctx, x, y, radius, startAngle, endAngle, clockwise)
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
	darwin.CGContextFillEllipseInRect(me.ctx, darwin.CGRectMake(x, y, width, height))
}

func (me *canvas) DrawPath(path Path) {
	// TODO::
}

func (me *canvas) DrawImage(img Image) {
	w, h := img.Size()
	darwin.CGContextDrawImage(me.ctx, darwin.CGRectMake(0, 0, float32(w), float32(h)), img.(*nativeImage).ref)
}

func (me *canvas) Flush() {
	// darwin.CGContextFlush(me.ctx)
}

func (me *canvas) Destroy() {
}
