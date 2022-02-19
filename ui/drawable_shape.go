// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
)

type ShapeDrawable interface {
	nux.Drawable
	Solid() nux.Color
	SetSolid(nux.Color)
}

func NewShapeDrawable(owner nux.Widget, attrs ...nux.Attr) ShapeDrawable {
	attr := nux.MergeAttrs(attrs...)
	me := &shapeDrawable{
		shape:        attr.GetString("shape", "rect"),
		solid:        attr.GetColor("solid", 0),
		stroke:       attr.GetColor("stroke", 0),
		strokeWidth:  attr.GetDimen("strokeWidth", "1px").Value(),
		cornerRadius: attr.GetDimen("cornerRadius", "0").Value(),
		paint:        nux.NewPaint(),
	}

	if shadow := attr.GetAttr("shadow", nil); shadow != nil {
		me.paint.SetShadow(shadow.GetColor("color", 0),
			shadow.GetDimen("x", "0").Value(),
			shadow.GetDimen("y", "0").Value(),
			shadow.GetDimen("blur", "0").Value())
	}

	return me
}

type shapeDrawable struct {
	x            int32
	y            int32
	width        int32
	height       int32
	shape        string
	solid        nux.Color
	stroke       nux.Color
	strokeWidth  float32
	cornerRadius float32
	shadowColor  nux.Color
	shadowX      float32
	shadowY      float32
	shadowBlur   float32
	paint        nux.Paint
}

func (me *shapeDrawable) Size() (width, height int32) {
	return me.width, me.height
}

func (me *shapeDrawable) SetBounds(x, y, width, height int32) {
	me.x = x
	me.y = y
	me.width = width
	me.height = height
}
func (me *shapeDrawable) Solid() nux.Color {
	return me.solid
}

func (me *shapeDrawable) SetSolid(color nux.Color) {
	me.solid = color
}

func (me *shapeDrawable) Stroke() nux.Color {
	return me.stroke
}

func (me *shapeDrawable) SetStroke(color nux.Color) {
	me.stroke = color
}

func (me *shapeDrawable) Draw(canvas nux.Canvas) {
	if me.solid != 0 {
		me.paint.SetColor(me.solid)
		me.paint.SetStyle(nux.PaintStyle_Fill)
		if me.cornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.cornerRadius, me.paint)
		} else {
			canvas.DrawRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.paint)
		}
	}

	if me.stroke != 0 {
		me.paint.SetColor(me.stroke)
		me.paint.SetStyle(nux.PaintStyle_Stroke)
		me.paint.SetWidth(me.strokeWidth)
		if me.cornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.cornerRadius, me.paint)
		} else {
			canvas.DrawRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.paint)
		}
	}

}

func (me *shapeDrawable) Equal(drawable nux.Drawable) bool {
	if c, ok := drawable.(*shapeDrawable); ok {
		if me.paint.Color().Equal(c.paint.Color()) &&
			me.x == c.x &&
			me.y == c.y &&
			me.width == c.width &&
			me.height == c.height {
			return true
		}
	}
	return false
}
