// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

type PanGestureCallback func(detail PanDetail)

type PanDetail interface {
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

func OnPanStart(widget Widget, callback PanGestureCallback) {
	addPanCallback(widget, _ACTION_PAN_START, callback)
}

func OnPanUpdate(widget Widget, callback PanGestureCallback) {
	addPanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func OnPanEnd(widget Widget, callback PanGestureCallback) {
	addPanCallback(widget, _ACTION_PAN_END, callback)
}

func RemovePanStartGesture(widget Widget, callback PanGestureCallback) {
	removePanCallback(widget, _ACTION_PAN_START, callback)
}

func RemovePanUpdateGesture(widget Widget, callback PanGestureCallback) {
	removePanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func RemovePanEndGesture(widget Widget, callback PanGestureCallback) {
	removePanCallback(widget, _ACTION_PAN_END, callback)
}

func addPanCallback(widget Widget, which int, callback PanGestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*panGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*panGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newPanGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureBinding().AddGestureRecognizer(widget, recognizer)
	}
}

func removePanCallback(widget Widget, which int, callback PanGestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*panGestureRecognizer)(nil)); r != nil {
		pan := r.(*panGestureRecognizer)
		pan.removeCallback(which, callback)
	} /*else {
		// ignore
		log.V("nuxui", "callback is not existed, maybe already removed.")
	}*/
}

///////////////////////////// PanGestureRecognizer   /////////////////////////////
const (
	// _ACTION_PAN_DOWN = iota
	_ACTION_PAN_START = iota
	_ACTION_PAN_UPDATE
	_ACTION_PAN_END
	_ACTION_PAN_CANCEL
)

type panDetail struct {
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

func (me *panDetail) Target() Widget   { return me.target }
func (me *panDetail) Time() time.Time  { return me.time }
func (me *panDetail) Kind() Kind       { return me.kind }
func (me *panDetail) X() float32       { return me.x }
func (me *panDetail) Y() float32       { return me.y }
func (me *panDetail) ScreenX() float32 { return me.screenX }
func (me *panDetail) ScreenY() float32 { return me.screenY }
func (me *panDetail) WindowX() float32 { return me.windowX }
func (me *panDetail) WindowY() float32 { return me.windowY }

type panGestureRecognizer struct {
	target    Widget
	callbacks [][]unsafe.Pointer
	initEvent Event
	state     GestureState
}

func newPanGestureRecognizer(target Widget) *panGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &panGestureRecognizer{
		target:    target,
		callbacks: [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		initEvent: nil,
		state:     GestureState_Ready,
	}
}

func (me *panGestureRecognizer) addCallback(which int, callback PanGestureCallback) {
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.callbacks[which] {
		if o == p {
			if true /*TODO:: debug*/ {
				log.Fatal("nuxui", fmt.Sprintf("The %s callback is already existed.", []string{"OnPanDown", "OnPanUp", "OnPan"}[which]))
			} else {
				return
			}
		}
	}

	me.callbacks[which] = append(me.callbacks[which], unsafe.Pointer(&callback))
}

func (me *panGestureRecognizer) removeCallback(which int, callback PanGestureCallback) {
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for i, o := range me.callbacks[which] {
		if o == p {
			me.callbacks[which] = append(me.callbacks[which][:i], me.callbacks[which][i+1:]...)
		}
	}
}

func (me *panGestureRecognizer) PointerAllowed(event Event) bool {
	if len(me.callbacks[_ACTION_PAN_START]) == 0 &&
		len(me.callbacks[_ACTION_PAN_UPDATE]) == 0 &&
		len(me.callbacks[_ACTION_PAN_END]) == 0 &&
		len(me.callbacks[_ACTION_PAN_CANCEL]) == 0 {
		return false
	}

	if event.IsPrimary() {
		return true
	}

	return false
}

func (me *panGestureRecognizer) HandlePointerEvent(event Event) {
	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		GestureArenaManager().Add(event.Pointer(), me)
		me.state = GestureState_Possible
	case Action_Move:
		if me.initEvent == nil || me.initEvent.Pointer() != event.Pointer() {
			return
		}

		if me.state == GestureState_Accepted {
			me.invokePanUpdate(event)
		} else if me.state == GestureState_Possible {
			if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
				GestureArenaManager().Resolve(event.Pointer(), me, true)
			}
		}

	case Action_Up:
		if me.state == GestureState_Accepted {
			me.invokePanEnd(event)
			me.reset()
		} else if me.state == GestureState_Possible {
			GestureArenaManager().Resolve(event.Pointer(), me, false)
		}
	}
}

func (me *panGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		// do not need PanCancel?
		// me.invokePanCancel()
	}
	me.reset()
}

func (me *panGestureRecognizer) AccpetGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.state = GestureState_Accepted
		me.invokePanStart(me.initEvent)
	}
}

func (me *panGestureRecognizer) reset() {
	me.initEvent = nil
	me.state = GestureState_Ready
}

func (me *panGestureRecognizer) eventToDetail(event Event) PanDetail {
	if event == nil {
		return &panDetail{target: me.target}
	}

	x := event.X()
	y := event.Y()
	if s, ok := me.target.(Size); ok {
		ms := s.MeasuredSize()
		x = event.X() - float32(ms.Position.X)
		y = event.Y() - float32(ms.Position.Y)
	}

	return &panDetail{
		target:  me.target,
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

func (me *panGestureRecognizer) invokePanStart(event Event) {
	log.V("nuxui", "invokePanStart")
	for _, c := range me.callbacks[_ACTION_PAN_START] {
		(*(*(func(PanDetail)))(c))(me.eventToDetail(event))
	}
}

func (me *panGestureRecognizer) invokePanUpdate(event Event) {
	log.V("nuxui", "invokePanUpdate")
	for _, c := range me.callbacks[_ACTION_PAN_UPDATE] {
		(*(*(func(PanDetail)))(c))(me.eventToDetail(event))
	}
}

func (me *panGestureRecognizer) invokePanEnd(event Event) {
	log.V("nuxui", "invokePanEnd")
	for _, c := range me.callbacks[_ACTION_PAN_END] {
		(*(*(func(PanDetail)))(c))(me.eventToDetail(event))
	}
}

func (me *panGestureRecognizer) invokePanCancel() {
	log.V("nuxui", "invokePanCancel")
	for _, c := range me.callbacks[_ACTION_PAN_CANCEL] {
		(*(*(func(PanDetail)))(c))(me.eventToDetail(nil))
	}
}
