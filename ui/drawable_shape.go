// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type Shape struct {
	Name         string
	Solid        nux.Color
	Stroke       nux.Color
	StrokeWidth  float32
	CornerRadius float32
	ShadowColor  nux.Color
	ShadowX      float32
	ShadowY      float32
	ShadowBlur   float32
}

type ShapeDrawable interface {
	nux.Drawable
	SetShape(shape Shape)
	Shape() Shape
}

func NewShapeDrawable(attrs ...nux.Attr) ShapeDrawable {
	attr := nux.MergeAttrs(attrs...)
	me := &shapeDrawable{
		drawables: attr,
		paint:     nux.NewPaint(),
	}
	me.getShapeAndSet()
	return me
}

type shapeDrawable struct {
	x         int32
	y         int32
	width     int32
	height    int32
	paint     nux.Paint
	drawables nux.Attr
	shape     *Shape
}

func (me *shapeDrawable) getShapeAndSet() {
	if state := me.drawables.GetString("state", "normal"); me.drawables.Has(state) {
		d := me.drawables.Get(state, nil)
		if d != nil {
			switch t := d.(type) {
			case *Shape:
				me.shape = t
			case nux.Attr:
				shadow := t.GetAttr("shadow", nux.Attr{})
				shape := &Shape{
					Name:         t.GetString("name", "rect"),
					Solid:        t.GetColor("solid", 0),
					Stroke:       t.GetColor("stroke", 0),
					StrokeWidth:  t.GetDimen("strokeWidth", "1px").Value(),
					CornerRadius: t.GetDimen("cornerRadius", "0").Value(),
					ShadowColor:  shadow.GetColor("color", 0),
					ShadowX:      shadow.GetDimen("x", "0").Value(),
					ShadowY:      shadow.GetDimen("y", "0").Value(),
					ShadowBlur:   shadow.GetDimen("blur", "0").Value(),
				}
				me.drawables[state] = shape
				me.shape = shape
			default:
				log.Fatal("nuxui", "unknow shape of %T", d)
			}
		}
	}
}

func (me *shapeDrawable) SetState(state nux.Attr) {
	me.drawables = nux.MergeAttrs(me.drawables, state)
	me.getShapeAndSet()
}

func (me *shapeDrawable) State() nux.Attr {
	return me.drawables
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

func (me *shapeDrawable) SetShape(shape Shape) {
	me.shape = &shape
}

func (me *shapeDrawable) Shape() Shape {
	if me.shape == nil {
		return Shape{}
	}
	return *me.shape
}

func (me *shapeDrawable) Draw(canvas nux.Canvas) {
	if me.shape == nil {
		return
	}

	if me.shape.ShadowColor != 0 {
		me.paint.SetShadow(me.shape.ShadowColor, me.shape.ShadowX, me.shape.ShadowY, me.shape.ShadowBlur)
	}

	if me.shape.Solid != 0 {
		me.paint.SetColor(me.shape.Solid)
		me.paint.SetStyle(nux.PaintStyle_Fill)
		if me.shape.CornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.shape.CornerRadius, me.paint)
		} else {
			canvas.DrawRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.paint)
		}
	}

	if me.shape.Stroke != 0 {
		me.paint.SetColor(me.shape.Stroke)
		me.paint.SetStyle(nux.PaintStyle_Stroke)
		me.paint.SetWidth(me.shape.StrokeWidth)
		if me.shape.CornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.shape.CornerRadius, me.paint)
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
