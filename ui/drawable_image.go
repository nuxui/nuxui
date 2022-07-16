// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

type ImageDrawable interface {
	nux.Drawable
}

// TODO:: if src exised, duplicate image
func NewImageDrawable(attr nux.Attr) ImageDrawable {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &imageDrawable{
		state: nux.State_Default,
		image: nil,
	}

	if states := attr.GetAttrArray("states", nil); states != nil {
		me.states = map[uint32]any{}

		for _, state := range states {
			if s := state.GetString("state", ""); s != "" {
				if state.Has("src") {
					me.states[mergedStateFromString(s)] = state.GetString("src", "")
				} else {
					log.E("nuxui", "the state need a src field")
				}
			} else {
				log.E("nuxui", "the state need a state field")
			}
		}
	} else {
		if src := attr.GetString("src", ""); src != "" {
			me.image, _ = nux.CreateImage(src)
		}
	}

	me.applyState()

	return me
}

func NewImageDrawableWithResource(src string) ImageDrawable {
	return NewImageDrawable(nux.Attr{
		"src": src,
	})
}

type imageDrawable struct {
	x      int32
	y      int32
	width  int32
	height int32
	image  nux.Image
	states map[uint32]any
	state  uint32
	scaleX float32
	scaleY float32
}

func (me *imageDrawable) applyState() {
	if me.states != nil {
		if i, ok := me.states[me.state]; ok {
			switch t := i.(type) {
			case nux.Image:
				me.image = t
			case string:
				me.image, _ = nux.CreateImage(t)
				me.states[me.state] = me.image
			default:
				log.Fatal("nuxui", "unknow image of %T:%s", t, t)
			}
		}
	}
}

func (me *imageDrawable) AddState(state uint32) {
	s := me.state
	s |= state
	me.state = s
	me.applyState()
}

func (me *imageDrawable) DelState(state uint32) {
	s := me.state
	s &= ^state
	me.state = s
	me.applyState()
}

func (me *imageDrawable) State() uint32 {
	return me.state
}

func (me *imageDrawable) Size() (width, height int32) {
	if me.image != nil {
		return me.image.PixelSize()
	}
	return 0, 0
}

func (me *imageDrawable) HasState() bool {
	return me.states == nil || len(me.states) == 0 || (len(me.states) == 1 && me.state == nux.State_Default)
}

func (me *imageDrawable) SetBounds(x, y, width, height int32) {
	me.x = x
	me.y = y
	me.width = width
	me.height = height
	w, h := me.Size()
	me.scaleX = float32(width) / float32(w)
	me.scaleY = float32(height) / float32(h)
}

func (me *imageDrawable) Draw(canvas nux.Canvas) {
	if me.image != nil {
		canvas.Save()
		canvas.Translate(float32(me.x), float32(me.y))
		canvas.Scale(me.scaleX, me.scaleY)
		canvas.DrawImage(me.image)
		canvas.Restore()
	}
}

func (me *imageDrawable) Equal(drawable nux.Drawable) bool {
	// TODO::
	return false
}
