// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "math"

type PointerEvent interface {
	Event
	Pointer() int64
	Kind() Kind
	X() int32
	Y() int32
	RawX() int32
	RawY() int32
	Offset(x, y int32)
	Action() EventAction
	IsPrimary() bool
	Distance(x, y int32) float64
}

type pointerEvent struct {
	event
	pointer int64
	kind    Kind
	x       int32
	y       int32
	rawX    int32
	rawY    int32
}

func (me *pointerEvent) Pointer() int64 {
	return me.pointer
}

func (me *pointerEvent) Kind() Kind {
	return me.kind
}

func (me *pointerEvent) X() int32 {
	return me.x
}

func (me *pointerEvent) Y() int32 {
	return me.y
}

func (me *pointerEvent) RawX() int32 {
	return me.rawX
}

func (me *pointerEvent) RawY() int32 {
	return me.rawY
}

func (me *pointerEvent) Offset(x, y int32) {
	me.x -= x
	me.y -= y
}

func (me *pointerEvent) Action() EventAction {
	return me.action
}

func (me *pointerEvent) IsPrimary() bool {
	return false
}

func (me *pointerEvent) Distance(x, y int32) float64 {
	dx := me.x - x
	dy := me.y - y
	return math.Sqrt(float64(dx)*float64(dx) + float64(dy)*float64(dy))
}

type MotionEvent interface {
	PointerEvent
}

type MouseEvent interface {
	PointerEvent
	Button() MouseButton
	Value() int // wheel value
}

type mouseEvent struct {
	pointerEvent
	button MouseButton
	value  int
}

func (me *mouseEvent) Button() MouseButton {
	return me.button
}

func (me *mouseEvent) Value() int {
	return me.value
}

func (me *mouseEvent) IsPrimary() bool {
	return me.button == MB_Left
}
