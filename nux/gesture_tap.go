// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

func OnTap(widget Widget, callback GestureCallback) {
	addTapCallback(widget, _ACTION_TAP, callback)
}

func OnTapDown(widget Widget, callback GestureCallback) {
	addTapCallback(widget, _ACTION_TAP_DOWN, callback)
}

func OnTapUp(widget Widget, callback GestureCallback) {
	addTapCallback(widget, _ACTION_TAP_UP, callback)
}

func OnTapCancel(widget Widget, callback GestureCallback) {
	addTapCallback(widget, _ACTION_TAP_CANCEL, callback)
}

// TODO:: widget will auto clear all gesture when destory
func RemoveTapGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, _ACTION_TAP, callback)
}

func RemoveTapDownGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, _ACTION_TAP_DOWN, callback)
}

func RemoveTapUpGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, _ACTION_TAP_UP, callback)
}

func RemoveTapCancelGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, _ACTION_TAP_CANCEL, callback)
}

func addTapCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*tapGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newTapGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureBinding().AddGestureRecognizer(widget, recognizer)
	}
}

func removeTapCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*tapGestureRecognizer)
		recognizer.removeCallback(which, callback)
	}
}

const (
	_ACTION_TAP_DOWN = iota
	_ACTION_TAP_UP
	_ACTION_TAP
	_ACTION_TAP_CANCEL
)

///////////////////////////// TapGestureRecognizer   /////////////////////////////

type TapGestureRecognizer interface {
	GestureRecognizer
}

type tapGestureRecognizer struct {
	callbacks   [][]unsafe.Pointer
	initEvent   Event
	target      Widget
	timer       Timer
	state       GestureState
	triggerDown bool
}

func newTapGestureRecognizer(target Widget) *tapGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &tapGestureRecognizer{
		callbacks:   [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		initEvent:   nil,
		target:      target,
		timer:       nil,
		state:       GestureState_Ready,
		triggerDown: false,
	}
}

func (me *tapGestureRecognizer) addCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.callbacks[which] {
		if o == p {
			if true /*TODO:: debug*/ {
				log.Fatal("nuxui", "The %s callback is already existed.", []string{"OnTapDown", "OnTapUp", "OnTap"}[which])

			} else {
				return
			}
		}
	}

	me.callbacks[which] = append(me.callbacks[which], unsafe.Pointer(&callback))
}

func (me *tapGestureRecognizer) removeCallback(which int, callback GestureCallback) {
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

func (me *tapGestureRecognizer) PointerAllowed(event Event) bool {
	if len(me.callbacks[_ACTION_TAP]) == 0 &&
		len(me.callbacks[_ACTION_TAP_DOWN]) == 0 &&
		len(me.callbacks[_ACTION_TAP_UP]) == 0 &&
		len(me.callbacks[_ACTION_TAP_CANCEL]) == 0 {
		return false
	}

	if event.IsPrimary() {
		return true
	}

	return false
}

func (me *tapGestureRecognizer) HandlePointerEvent(event Event) {
	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		pointer := event.Pointer()
		GestureArenaManager().Add(pointer, me)

		if event.Kind() == Kind_Touch {
			me.timer = NewTimerBackToUI(GESTURE_DOWN_DELAY, func() {
				me.invokeTapDown(pointer)
			})
		} else {
			me.invokeTapDown(pointer)
		}
	case Action_Up:
		// do not accept proactive, wait GestureArea sweep
	case Action_Move:
		if me.state == GestureState_Possible {
			if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
				GestureArenaManager().Resolve(event.Pointer(), me, false)
			}
		}
	}
}

func (me *tapGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.invokeTapCancel()
	}
	me.reset()
}

func (me *tapGestureRecognizer) AccpetGesture(pointer int64) {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	if me.state == GestureState_Possible {
		me.invokeTapDown(pointer)
		me.invokeTapUpAndTap(pointer)
	}
	me.reset()
}

func (me *tapGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	me.state = GestureState_Ready
	me.triggerDown = false
	me.initEvent = nil
}

func (me *tapGestureRecognizer) invokeTapDown(pointer int64) {
	if me.triggerDown {
		return
	}

	me.triggerDown = true
	me.state = GestureState_Possible
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_TAP_DOWN] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *tapGestureRecognizer) invokeTapCancel() {
	if me.triggerDown {
		for _, c := range me.callbacks[_ACTION_TAP_CANCEL] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *tapGestureRecognizer) invokeTapUpAndTap(pointer int64) {
	if !me.triggerDown || me.initEvent.Pointer() != pointer {
		return
	}

	me.state = GestureState_Accepted

	for _, c := range me.callbacks[_ACTION_TAP_UP] {
		(*(*(func(Widget)))(c))(me.target)
	}

	for _, c := range me.callbacks[_ACTION_TAP] {
		(*(*(func(Widget)))(c))(me.target)
	}
}
