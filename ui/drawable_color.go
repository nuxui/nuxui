// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"reflect"

	"github.com/nuxui/nuxui/nux"
)

type ColorDrawable interface {
	nux.Drawable
	Color() nux.Color
	SetColor(nux.Color)
}

func NewColorDrawable(attrs ...nux.Attr) ColorDrawable {
	attr := nux.MergeAttrs(attrs...)
	me := &colorDrawable{
		paint:     nux.NewPaint(),
		drawables: attr,
	}
	me.getColorAndSet()
	return me
}

func NewColorDrawableWithColor(color nux.Color) ColorDrawable {
	me := &colorDrawable{
		paint:     nux.NewPaint(),
		drawables: nux.Attr{},
	}
	me.paint.SetColor(color)
	return me
}

type colorDrawable struct {
	x         int32
	y         int32
	width     int32
	height    int32
	paint     nux.Paint
	drawables nux.Attr
}

func (me *colorDrawable) getColorAndSet() {
	if state := me.drawables.GetString("state", "default"); me.drawables.Has(state) {
		me.paint.SetColor(me.drawables.GetColor(state, 0))
	} else {
		me.paint.SetColor(me.drawables.GetColor("color", 0))
	}
}

func (me *colorDrawable) SetState(state nux.Attr) {
	me.drawables = nux.MergeAttrs(me.drawables, state)
	me.getColorAndSet()
}

func (me *colorDrawable) State() nux.Attr {
	return me.drawables
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
	return reflect.DeepEqual(me, drawable)
	// if c, ok := drawable.(*colorDrawable); ok {
	// 	if me.paint.Color().Equal(c.paint.Color()) &&
	// 		me.left == c.left &&
	// 		me.top == c.top &&
	// 		me.right == c.right &&
	// 		me.bottom == c.bottom {
	// 		return true
	// 	}
	// }
	// return false
}
