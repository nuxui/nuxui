// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

type ImageDrawable interface {
	nux.Drawable
}

// TODO:: if src exised, duplicate image
func NewImageDrawable(attrs ...nux.Attr) ImageDrawable {
	attr := nux.MergeAttrs(attrs...)
	return NewImageDrawableWithResource(attr.GetString("src", ""))
}

func NewImageDrawableWithResource(src string) ImageDrawable {
	if src != "" {
		return &imageDrawable{
			image: nux.CreateImage(src),
		}
	}

	return &imageDrawable{}
}

type imageDrawable struct {
	x      int32
	y      int32
	width  int32
	height int32
	image  nux.Image
	state  nux.Attr
}

func (me *imageDrawable) SetState(state nux.Attr) {
	nux.MergeAttrs(me.state, state)
}

func (me *imageDrawable) State() nux.Attr {
	return me.state
}

func (me *imageDrawable) Size() (width, height int32) {
	if me.image == nil {
		return 0, 0
	}
	return me.image.Size()
}

func (me *imageDrawable) SetBounds(x, y, width, height int32) {
	me.x = x
	me.y = y
	me.width = width
	me.height = height
}

func (me *imageDrawable) Draw(canvas nux.Canvas) {
	canvas.Save()
	canvas.Translate(float32(me.x), float32(me.y))
	canvas.DrawImage(me.image)
	canvas.Restore()
}

func (me *imageDrawable) Equal(drawable nux.Drawable) bool {
	// TODO::
	return false
}
