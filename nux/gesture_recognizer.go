// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"time"

	"github.com/nuxui/nuxui/log"
)

type GestureRecognizer interface {
	GestureArenaMember
	PointerAllowed(event Event) bool
	HandlePointerEvent(event Event)
	//TODO:: Clear(widget Widget) // clear callbacks of widget
}

type GestureCallback func(detail GestureDetail)

type GestureDetail interface {
	Target() Widget
	Time() time.Time
	Kind() Kind
	X() float32
	Y() float32
	ScreenX() float32
	ScreenY() float32
	WindowX() float32
	WindowY() float32
}

type gestureDetail struct {
	target  Widget
	time    time.Time
	kind    Kind
	x       float32
	y       float32
	screenX float32
	screenY float32
	windowX float32
	windowY float32
}

func (me *gestureDetail) Target() Widget   { return me.target }
func (me *gestureDetail) Time() time.Time  { return me.time }
func (me *gestureDetail) Kind() Kind       { return me.kind }
func (me *gestureDetail) X() float32       { return me.x }
func (me *gestureDetail) Y() float32       { return me.y }
func (me *gestureDetail) ScreenX() float32 { return me.screenX }
func (me *gestureDetail) ScreenY() float32 { return me.screenY }
func (me *gestureDetail) WindowX() float32 { return me.windowX }
func (me *gestureDetail) WindowY() float32 { return me.windowY }

func eventToDetail(event Event, target Widget) GestureDetail {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	if event == nil {
		return &gestureDetail{target: target}
	}

	x := event.X()
	y := event.Y()
	if s, ok := target.(Size); ok {
		ms := s.MeasuredSize()
		x = event.X() - float32(ms.Position.X)
		y = event.Y() - float32(ms.Position.Y)
	}

	return &gestureDetail{
		target:  target,
		time:    event.Time(),
		kind:    event.Kind(),
		x:       x,
		y:       y,
		screenX: event.ScreenX(),
		screenY: event.ScreenY(),
		windowX: event.X(),
		windowY: event.Y(),
	}
}
