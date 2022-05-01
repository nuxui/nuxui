// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
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

func NewShape(attr nux.Attr) *Shape {
	shadow := attr.GetAttr("shadow", nux.Attr{})
	return &Shape{
		Name:         attr.GetString("name", "rect"),
		Solid:        attr.GetColor("solid", 0),
		Stroke:       attr.GetColor("stroke", 0),
		StrokeWidth:  attr.GetDimen("strokeWidth", "1px").Value(),
		CornerRadius: attr.GetDimen("cornerRadius", "0").Value(),
		ShadowColor:  shadow.GetColor("color", 0),
		ShadowX:      shadow.GetDimen("x", "0").Value(),
		ShadowY:      shadow.GetDimen("y", "0").Value(),
		ShadowBlur:   shadow.GetDimen("blur", "0").Value(),
	}
}

type ShapeDrawable interface {
	nux.Drawable
	SetShape(shape Shape)
	Shape() Shape
}

func NewShapeDrawable(attr nux.Attr) ShapeDrawable {
	me := &shapeDrawable{
		state: nux.State_Default,
		paint: nux.NewPaint(),
	}

	if states := attr.GetAttrArray("states", nil); states != nil {
		me.states = map[uint32]*Shape{}

		for _, state := range states {
			if s := state.GetString("state", ""); s != "" {
				if state.Has("shape") {
					me.states[mergedStateFromString(s)] = NewShape(state.GetAttr("shape", nil))
				} else {
					log.E("nuxui", "the state need a shape field")
				}
			} else {
				log.E("nuxui", "the state need a state field")
			}
		}
	} else {
		if shape := attr.GetAttr("shape", nil); shape != nil {
			me.shape = NewShape(shape)
		}
	}

	me.applyState()
	return me
}

type shapeDrawable struct {
	x      int32
	y      int32
	width  int32
	height int32
	paint  nux.Paint
	shape  *Shape
	state  uint32
	states map[uint32]*Shape
}

func (me *shapeDrawable) applyState() {
	if me.states != nil {
		if s, ok := me.states[me.state]; ok {
			me.shape = s
		}
	}
}

func (me *shapeDrawable) AddState(state uint32) {
	s := me.state
	s |= state
	me.state = s
	me.applyState()
}

func (me *shapeDrawable) DelState(state uint32) {
	s := me.state
	s &= ^state
	me.state = s
	me.applyState()
}

func (me *shapeDrawable) State() uint32 {
	return me.state
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
