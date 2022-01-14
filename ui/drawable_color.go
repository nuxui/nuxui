// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
)

type ColorDrawable interface {
	Drawable
	Color() nux.Color
	SetColor(nux.Color)
}

func NewColorDrawable() ColorDrawable {
	return &colorDrawable{color: nux.Transparent}
}

type colorDrawable struct {
	color nux.Color
}

func (me *colorDrawable) Color() nux.Color {
	return me.color
}

func (me *colorDrawable) SetColor(color nux.Color) {
	me.color = color
}

func (me *colorDrawable) Size() (width, height int32) {
	return 0, 0
}

func (me *colorDrawable) Draw(canvas nux.Canvas) {
	// t1 := time.Now()
	canvas.DrawColor(me.color)
	// log.V("nuxui", "draw DrawColor used time %d", time.Now().Sub(t1).Milliseconds())
}

func (me *colorDrawable) Equal(drawable Drawable) bool {
	if c, ok := drawable.(*colorDrawable); ok {
		if me.color == c.color {
			return true
		}
	}
	return false
}
