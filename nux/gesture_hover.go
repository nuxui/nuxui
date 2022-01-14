// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"time"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

func OnHover(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, _ACTION_HOVER, callback)
}

func OnHoverEnter(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, _ACTION_HOVER_ENTER, callback)
}

func OnHoverExit(widget Widget, callback GestureCallback) {
	addHoverCallback(widget, _ACTION_HOVER_EXIT, callback)
}

// widget will auto clear all gesture when destroy
func RemoveHoverGesture(widget Widget, callback GestureCallback) {
	removeHoverCallback(widget, _ACTION_HOVER, callback)
}

func RemoveHoverEnterGesture(widget Widget, callback GestureCallback) {
	removeHoverCallback(widget, _ACTION_HOVER_ENTER, callback)
}

func RemoveHoverExitGesture(widget Widget, callback GestureCallback) {
	removeHoverCallback(widget, _ACTION_HOVER_EXIT, callback)
}

func addHoverCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*hoverGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*hoverGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newHoverGestureRecognizer(widget).(*hoverGestureRecognizer)
		recognizer.addCallback(which, callback)
		GestureBinding().AddGestureRecognizer(widget, recognizer)
	}
}

func removeHoverCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*hoverGestureRecognizer)(nil)); r != nil {
		hover := r.(*hoverGestureRecognizer)
		hover.removeCallback(which, callback)
	}
}

const (
	_ACTION_HOVER_ENTER = iota
	_ACTION_HOVER_EXIT
	_ACTION_HOVER
)

///////////////////////////// HoverGestureRecognizer   /////////////////////////////

type HoverGestureRecognizer interface {
	GestureRecognizer
}

type hoverGestureRecognizer struct {
	callbacks   [][]unsafe.Pointer
	initEvent   PointerEvent
	target      Widget
	timer       Timer
	state       GestureState
	triggerDown bool
}

func newHoverGestureRecognizer(target Widget) HoverGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &hoverGestureRecognizer{
		callbacks: [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		// initEvent:   PointerEvent{Pointer: 0},
		target:      target,
		state:       GestureState_Ready,
		triggerDown: false,
	}
}

func (me *hoverGestureRecognizer) addCallback(which int, callback GestureCallback) {
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

func (me *hoverGestureRecognizer) removeCallback(which int, callback GestureCallback) {
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

func (me *hoverGestureRecognizer) PointerAllowed(event PointerEvent) bool {
	if len(me.callbacks[_ACTION_HOVER]) == 0 &&
		len(me.callbacks[_ACTION_HOVER_ENTER]) == 0 &&
		len(me.callbacks[_ACTION_HOVER_EXIT]) == 0 &&
		len(me.callbacks[_ACTION_HOVER_EXIT]) == 0 {
		return false
	}

	return true
}

func (me *hoverGestureRecognizer) HandlePointerEvent(event PointerEvent) {
	if event.IsPrimary() {
		switch event.Action() {
		case Action_Down:
			log.V("nuxui", "IsPrimaryButton Action_Down")
			me.initEvent = event
			GestureArenaManager().Add(event.Pointer(), me)

			p := me.initEvent.Pointer()
			log.V("nuxui", "IsPrimaryButton Action_Down NewTimer")
			me.timer = NewTimerBackToUI(GESTURE_DOWN_DELAY*time.Millisecond, func() {
				me.invokeHoverEnter(p)
			})
		case Action_Up:
			GestureArenaManager().Resolve(me.initEvent.Pointer(), me, me.initEvent.Pointer() == event.Pointer())
		case Action_Move:
			log.V("nuxui", "IsPrimaryButton Action_Move")
			if me.state == GestureState_Possible {
				if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
					GestureArenaManager().Resolve(event.Pointer(), me, false)
				}
			}
		}
	}
}

func (me *hoverGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.invokeHoverExit()
	}
	me.reset()
}

func (me *hoverGestureRecognizer) AccpetGesture(pointer int64) {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}
	me.invokeHoverEnter(pointer)
	me.invokeHover(pointer)
	me.reset()
}

func (me *hoverGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	me.state = GestureState_Ready
	me.triggerDown = false
	me.initEvent = nil
}

func (me *hoverGestureRecognizer) invokeHoverEnter(pointer int64) {
	if me.triggerDown {
		return
	}

	me.triggerDown = true
	me.state = GestureState_Possible
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_HOVER_ENTER] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *hoverGestureRecognizer) invokeHoverExit() {
	for _, c := range me.callbacks[_ACTION_HOVER_EXIT] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *hoverGestureRecognizer) invokeHover(pointer int64) {
	for _, c := range me.callbacks[_ACTION_HOVER] {
		(*(*(func(Widget)))(c))(me.target)
	}
}
