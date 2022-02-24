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
	return me
}

func NewColorDrawableWithColor(color nux.Color) ColorDrawable {
	return NewColorDrawable(nux.Attr{
		"state":  "normal",
		"normal": color,
	})
}

type colorDrawable struct {
	x         int32
	y         int32
	width     int32
	height    int32
	paint     nux.Paint
	drawables nux.Attr
}

func (me *colorDrawable) SetState(state nux.Attr) {
	nux.MergeAttrs(me.drawables, state)
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
	return me.drawables.GetColor(me.drawables.GetString("state", "normal"), 0)
}

func (me *colorDrawable) SetColor(color nux.Color) {
	me.drawables.Set("normal", color)
	me.drawables.Set("state", "normal")
}

func (me *colorDrawable) Draw(canvas nux.Canvas) {
	c := me.Color()
	if c != 0 {
		me.paint.SetColor(me.Color())
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
