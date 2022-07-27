// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || (windows && cairo)

package nux

import (
	"nuxui.org/nuxui/log"
)

func (me *canvas) native() *canvas {
	return me
}

func (me *canvas) ResetClip() {
}

func (me *canvas) Save() {
	me.cairo.Save()
}

func (me *canvas) Restore() {
	me.cairo.Restore()
}

func (me *canvas) Translate(x, y float32) {
	me.cairo.Translate(x, y)
}

func (me *canvas) Scale(x, y float32) {
	me.cairo.Scale(x, y)
}

func (me *canvas) Rotate(angle float32) {
	me.cairo.Rotate(_RADIAN * angle)
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
	if width < 0 || height < 0 {
		log.Fatal("nuxui", "invalid rect for clip")
	}
	me.cairo.Rectangle(x, y, width, height)
	me.cairo.Clip()
}

func (me *canvas) ClipRoundRect(x, y, width, height, cornerX, cornerY float32) {
	// TODO::
}

func (me *canvas) ClipPath(path Path) {
	// TODO::
}

func (me *canvas) SetAlpha(alpha float32) {
}

func (me *canvas) DrawRect(x, y, width, height float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}

	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		// C.cairo_identity_matrix(me.ptr)
		me.cairo.Save()
		me.cairo.Translate(0.5, 0.5)

	}

	me.cairo.Rectangle(x, y, width, height)
	me.drawPaint(paint)

	if fix {
		me.cairo.Restore()
	}
}

func (me *canvas) DrawRoundRect(x, y, width, height float32, rLT, rRT, rRB, rLB float32, paint Paint) {
	if width < 0 || height < 0 {
		return
	}

	fix := paint.Style() == PaintStyle_Stroke && int32(paint.Width())%2 != 0
	if fix {
		me.Save()
		me.Translate(0.5, 0.5)
	}

	path := newPath()
	path.RoundRect(x, y, width, height, rLT, rRT, rRB, rLB)

	// me.cairo.NewSubPath()
	// me.cairo.Arc(x+width-rRT, y+rRT, rRT, -90*_RADIAN, 0)
	// me.cairo.Arc(x+width-rRB, y+height-rRB, rRB, 0, 90*_RADIAN)
	// me.cairo.Arc(x+rLB, y+height-rLB, rLB, 90*_RADIAN, 180*_RADIAN)
	// me.cairo.Arc(x+rLT, y+rLT, rLT, 180*_RADIAN, 270*_RADIAN)
	// me.cairo.ClosePath()

	// TODO:: shadow
	if sc, sx, sy, sb := paint.Shadow(); sc != 0 && sb > 0 {
		me.Save()
		me.Translate(sx, sy)

		me.cairo.NewPath()
		copy := path.native().cairo.CopyPath()
		me.cairo.AppendPath(copy)
		copy.Destroy()
		r, g, b, a := sc.RGBAf()
		me.cairo.SetSourceRGBA(r, g, b, a)
		me.cairo.Fill()
		me.Restore()
	}

	me.DrawPath(path, paint)

	if fix {
		me.cairo.Restore()
	}
}

func (me *canvas) DrawArc(x, y, radius, startAngle, endAngle float32, useCenter bool, paint Paint) {
	// TODO:: useCenter
	if useCenter {
		me.cairo.Arc(x, y, radius, startAngle*_RADIAN, endAngle*_RADIAN)
		me.drawPaint(paint)
	} else {
		me.cairo.Arc(x, y, radius, startAngle*_RADIAN, endAngle*_RADIAN)
		me.drawPaint(paint)
	}
}

func (me *canvas) DrawOval(x, y, width, height float32, paint Paint) {
	me.cairo.Save()
	var centerX, centerY, scaleX, scaleY, radius float32
	if width > height {
		centerX = x + width/2.0
		centerY = y + width/2.0
		scaleX = 1.0
		scaleY = height / width
		radius = width / 2.0
	} else {
		centerX = x + height/2.0
		centerY = y + height/2.0
		scaleX = width / height
		scaleY = 1.0
		radius = height / 2.0
	}

	me.cairo.Scale(scaleX, scaleY)
	me.cairo.Arc(centerX, centerY, radius, 0, _PI2)
	me.drawPaint(paint)
	me.cairo.Restore()
}

func (me *canvas) drawPaint(paint Paint) {
	r, g, b, a := paint.Color().RGBAf()
	// C.cairo_fill_preserve(me.ptr)
	me.cairo.SetSourceRGBA(r, g, b, a)
	me.cairo.SetLineWidth(paint.Width())
	switch paint.Style() {
	case PaintStyle_Stroke:
		me.cairo.Stroke()
	case PaintStyle_Fill:
		me.cairo.Fill()
	case PaintStyle_Both:
		me.cairo.Stroke()
		me.cairo.Fill()
	}
}

func (me *canvas) DrawColor(color Color) {
	r, g, b, a := color.RGBAf()
	me.cairo.SetSourceRGBA(r, g, b, a)
	me.cairo.Paint()
}

func (me *canvas) DrawPath(path Path, paint Paint) {
	me.cairo.NewPath()
	copy := path.native().cairo.CopyPath()
	me.cairo.AppendPath(copy)
	copy.Destroy()
	me.drawPaint(paint)
}

func (me *canvas) DrawImage(img Image) {
	img.Draw(me)
}

func (me *canvas) DrawText(text string, width, height float32, paint Paint) {
	// cfamily := C.CString("")
	// ctext := C.CString(text)
	// r, g, b, a := paint.Color().RGBAf()
	// me.cairo.SetSourceRGBA(r,g,b,a)
	// C.drawText(me.ptr, cfamily, C.int(1), C.int(paint.TextSize()), ctext, C.int(width), C.int(height))
	// me.drawPaint(paint)
	// C.free(unsafe.Pointer(cfamily))
	// C.free(unsafe.Pointer(ctext))
}

func (me *canvas) Flush() {
	me.surface.Flush()
}

func (me *canvas) Destroy() {
}
