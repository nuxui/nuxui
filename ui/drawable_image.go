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
	me := &imageDrawable{
		state: attr,
		image: map[string]nux.Image{},
	}
	me.createImagesIfNeed()
	return me
}

func NewImageDrawableWithResource(src string) ImageDrawable {
	return NewImageDrawable(nux.Attr{
		"default": src,
	})
}

type imageDrawable struct {
	x      int32
	y      int32
	width  int32
	height int32
	image  map[string]nux.Image
	state  nux.Attr
	scaleX float32
	scaleY float32
}

func (me *imageDrawable) createImagesIfNeed() {
	if state := me.state.GetString("state", "default"); me.state.Has(state) {
		if _, ok := me.image[state]; !ok {
			if src := me.state.GetString(state, ""); src != "" {
				me.image[state] = nux.CreateImage(src)
			}
		}
	} else {
		if _, ok := me.image["default"]; !ok {
			if src := me.state.GetString("src", ""); src != "" {
				me.image[state] = nux.CreateImage(src)
			}
		}
	}
}

func (me *imageDrawable) SetState(state nux.Attr) {
	nux.MergeAttrs(me.state, state)
	me.createImagesIfNeed()
}

func (me *imageDrawable) State() nux.Attr {
	return me.state
}

func (me *imageDrawable) Size() (width, height int32) {
	state := me.state.GetString("state", "default")
	if i, ok := me.image[state]; ok {
		return i.Size()
	}
	return 0, 0
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
	state := me.state.GetString("state", "default")
	if img, ok := me.image[state]; ok {
		canvas.Save()
		canvas.Translate(float32(me.x), float32(me.y))
		canvas.Scale(me.scaleX, me.scaleY)
		canvas.DrawImage(img)
		canvas.Restore()
	}
}

func (me *imageDrawable) Equal(drawable nux.Drawable) bool {
	// TODO::
	return false
}
