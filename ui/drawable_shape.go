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
}

func NewShapeDrawable(attrs ...nux.Attr) ShapeDrawable {
	attr := nux.MergeAttrs(attrs...)
	me := &shapeDrawable{

		drawables: attr,
		paint:     nux.NewPaint(),
	}
	return me
}

type shapeDrawable struct {
	x         int32
	y         int32
	width     int32
	height    int32
	paint     nux.Paint
	drawables nux.Attr
}

func (me *shapeDrawable) SetState(state nux.Attr) {
	nux.MergeAttrs(me.drawables, state)
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

func (me *shapeDrawable) SetShape(shape *Shape) {
	me.drawables["normal"] = shape
}

func (me *shapeDrawable) Shape() *Shape {
	state := me.drawables.GetString("state", "normal")
	d := me.drawables.Get(state, nil)
	if d == nil {
		return &Shape{}
	}
	if s, ok := d.(*Shape); ok {
		return s
	}
	// if attr, ok := d.(nux.Attr);
	switch t := d.(type) {
	case *Shape:
		return t
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
		return shape
	default:
		log.Fatal("nuxui", "unknow shape of %T", d)
	}
	return nil
}

func (me *shapeDrawable) Draw(canvas nux.Canvas) {
	shape := me.Shape()

	if shape.ShadowColor != 0 {
		me.paint.SetShadow(shape.ShadowColor, shape.ShadowX, shape.ShadowY, shape.ShadowBlur)
	}

	if shape.Solid != 0 {
		me.paint.SetColor(shape.Solid)
		me.paint.SetStyle(nux.PaintStyle_Fill)
		if shape.CornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), shape.CornerRadius, me.paint)
		} else {
			canvas.DrawRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.paint)
		}
	}

	if shape.Stroke != 0 {
		me.paint.SetColor(shape.Stroke)
		me.paint.SetStyle(nux.PaintStyle_Stroke)
		me.paint.SetWidth(shape.StrokeWidth)
		if shape.CornerRadius > 0 {
			canvas.DrawRoundRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), shape.CornerRadius, me.paint)
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
