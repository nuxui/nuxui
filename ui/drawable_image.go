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

// TODO:: if src exised, get cached image
func NewImageDrawable(attr nux.Attr) ImageDrawable {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &imageDrawable{
		image: nil,
	}

	me.StateBase = nux.NewStateBase(me.applyState)

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
			me.image, _ = nux.LoadImageFromFile(src)
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
	*nux.StateBase

	x      int32
	y      int32
	width  int32
	height int32
	image  nux.Image
	states map[uint32]any
	scaleX float32
	scaleY float32
}

func (me *imageDrawable) applyState() {
	if me.states != nil {
		if i, ok := me.states[me.State()]; ok {
			switch t := i.(type) {
			case nux.Image:
				me.image = t
			case string:
				var err error
				me.image, err = nux.LoadImageFromFile(t)
				if err != nil {
					log.Fatal("nuxui", "load image from file: %s", err.Error())
				}
				me.states[me.State()] = me.image
			default:
				log.Fatal("nuxui", "unknow image of %T:%s", t, t)
			}
		}
	}
}

func (me *imageDrawable) Size() (width, height int32) {
	if me.image != nil {
		return me.image.PixelSize()
	}
	return 0, 0
}

func (me *imageDrawable) HasState() bool {
	return !(me.states == nil || len(me.states) == 0 || (len(me.states) == 1 && me.State() == nux.State_Default))
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
