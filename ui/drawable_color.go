// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

type ColorDrawable interface {
	nux.Drawable
	Color() nux.Color
	SetColor(nux.Color)
}

func NewColorDrawable(attr nux.Attr) ColorDrawable {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &colorDrawable{
		paint: nux.NewPaint(),
	}

	me.StateBase = nux.NewStateBase(me.applyState)

	if states := attr.GetAttrArray("states", nil); states != nil {
		me.states = map[uint32]nux.Color{}

		for _, state := range states {
			if s := state.GetString("state", ""); s != "" {
				if state.Has("color") {
					me.states[mergedStateFromString(s)] = state.GetColor("color", nux.Transparent)
				} else {
					log.E("nuxui", "the state need a color field")
				}
			} else {
				log.E("nuxui", "the state need a state field")
			}
		}
	} else {
		me.paint.SetColor(attr.GetColor("color", nux.Transparent))
	}

	me.applyState()
	return me
}

func NewColorDrawableWithColor(color nux.Color) ColorDrawable {
	return NewColorDrawable(nux.Attr{"color": color})
}

type colorDrawable struct {
	*nux.StateBase

	x      int32
	y      int32
	width  int32
	height int32
	paint  nux.Paint
	states map[uint32]nux.Color
	state  uint32
}

func (me *colorDrawable) applyState() {
	if me.states != nil {
		if color, ok := me.states[me.State()]; ok {
			me.paint.SetColor(color)
		}
	}
}

func (me *colorDrawable) HasState() bool {
	return !(me.states == nil || len(me.states) == 0 || (len(me.states) == 1 && me.State() == nux.State_Default))
}

func (me *colorDrawable) Size() (width, height int32) {
	return me.width, me.height
}

func (me *colorDrawable) SetBounds(x, y, width, height int32) {
	me.x = x
	me.y = y
	me.width = width
	me.height = height
}

func (me *colorDrawable) Color() nux.Color {
	return me.paint.Color()
}

func (me *colorDrawable) SetColor(color nux.Color) {
	me.paint.SetColor(color)
}

func (me *colorDrawable) Draw(canvas nux.Canvas) {
	if me.paint.Color() != 0 {
		canvas.DrawRect(float32(me.x), float32(me.y), float32(me.width), float32(me.height), me.paint)
	}
}

func (me *colorDrawable) Equal(drawable nux.Drawable) bool {
	if c, ok := drawable.(*colorDrawable); ok {
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
