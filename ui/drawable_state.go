// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"reflect"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type StateDrawable interface {
	nux.Drawable
	SetState(state DrawableState)
	State() DrawableState
}

type DrawableState int

const (
	DrawableState_Normal DrawableState = iota
	DrawableState_Hover
	DrawableState_Focused
	DrawableState_Pressed
)

type stateDrawable struct {
	owner     nux.Widget
	state     DrawableState
	drawables []nux.Drawable
}

func NewStateDrawable(owner nux.Widget, attrs ...nux.Attr) StateDrawable {
	log.V("nuxui", "NewStateDrawable, %T, %s", owner, owner.Info().ID)

	attr := nux.MergeAttrs(attrs...)

	nux.OnTap(owner, func(detail nux.GestureDetail) {
		log.V("nuxui", "StateDrawable tapDown")
	})

	me := &stateDrawable{
		owner:     owner,
		drawables: []nux.Drawable{nil, nil, nil, nil},
	}
	nux.AddMixins(owner, me)

	if normal := attr.Get("normal", nil); normal != nil {
		me.drawables[DrawableState_Normal] = nux.InflateDrawable(owner, normal)
	}

	if hover := attr.Get("hover", nil); hover != nil {
		me.drawables[DrawableState_Hover] = nux.InflateDrawable(owner, hover)
	}

	if focused := attr.Get("focused", nil); focused != nil {
		me.drawables[DrawableState_Focused] = nux.InflateDrawable(owner, focused)
	}

	if pressed := attr.Get("pressed", nil); pressed != nil {
		me.drawables[DrawableState_Pressed] = nux.InflateDrawable(owner, pressed)
	}

	return me
}

func (me *stateDrawable) OnMount() {
	log.V("nuxui", "StateDrawable OnMount")

	nux.OnTap(me.owner, func(detail nux.GestureDetail) {
		log.V("nuxui", "StateDrawable tapDown")
	})
}

func (me *stateDrawable) Size() (width, height int32) {
	// return max width and height
	for _, d := range me.drawables {
		w, h := d.Size()
		if w > width {
			width = w
		}
		if h > height {
			height = h
		}
	}
	return
}

func (me *stateDrawable) SetBounds(x, y, width, height int32) {
	log.I("nuxui", "stateDrawable bounds x=%d, y=%d, w=%d, h=%d", x, y, width, height)
	for _, d := range me.drawables {
		if d != nil {
			d.SetBounds(x, y, width, height)
		}
	}
}

func (me *stateDrawable) Draw(canvas nux.Canvas) {
	me.drawables[me.state].Draw(canvas)
}

func (me *stateDrawable) Equal(drawable nux.Drawable) bool {
	return reflect.DeepEqual(me, drawable)
}

func (me *stateDrawable) SetState(state DrawableState) {
	if me.state != state {
		me.state = state
		// TODO change
	}
}

func (me *stateDrawable) State() DrawableState {
	return me.state
}
