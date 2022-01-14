// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

type ImageDrawable interface {
	Drawable
}

type imageDrawable struct {
	image nux.Image
}

// TODO:: if src exised
func NewImageDrawable(src string) ImageDrawable {
	if src != "" {
		return &imageDrawable{
			image: nux.CreateImage(src),
		}
	}

	return &imageDrawable{}
}

func (me *imageDrawable) Size() (width, height int32) {
	if me.image == nil {
		return 0, 0
	}
	return me.image.Size()
}

func (me *imageDrawable) Draw(canvas nux.Canvas) {
	canvas.DrawImage(me.image)
}

func (me *imageDrawable) Equal(drawable Drawable) bool {
	return false
}
