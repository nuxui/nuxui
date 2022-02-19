// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/util"
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
	if r := GestureManager().FindGestureRecognizer(widget, (*hoverGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*hoverGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newHoverGestureRecognizer(widget).(*hoverGestureRecognizer)
		recognizer.addCallback(which, callback)
		GestureManager().AddGestureRecognizer(widget, recognizer)
	}
}

func removeHoverCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureManager().FindGestureRecognizer(widget, (*hoverGestureRecognizer)(nil)); r != nil {
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
	callbacks   [][]GestureCallback
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
		callbacks:   [][]GestureCallback{[]GestureCallback{}, []GestureCallback{}, []GestureCallback{}},
		initEvent:   nil,
		target:      target,
		state:       GestureState_Ready,
		triggerDown: false,
	}
}

func (me *hoverGestureRecognizer) addCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	for _, cb := range me.callbacks[which] {
		if util.SameFunc(cb, callback) {
			log.Fatal("nuxui", "The %s callback is already existed.", []string{"OnPanDown", "OnPanUp", "OnPan"}[which])
		}
	}

	me.callbacks[which] = append(me.callbacks[which], callback)
}

func (me *hoverGestureRecognizer) removeCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	for i, cb := range me.callbacks[which] {
		if util.SameFunc(cb, callback) {
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

			log.V("nuxui", "IsPrimaryButton Action_Down NewTimer")
			me.timer = NewTimerBackToUI(GESTURE_DOWN_DELAY, func() {
				me.invokeHoverEnter(event)
			})
		case Action_Move:
			log.V("nuxui", "IsPrimaryButton Action_Move")
			if me.state == GestureState_Possible {
				me.invokeHover(event)
			}
		case Action_Up:
			GestureArenaManager().Resolve(me.initEvent.Pointer(), me, me.initEvent.Pointer() == event.Pointer())
		}
	}
}

func (me *hoverGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.invokeHoverExit(me.initEvent)
	}
	me.reset()
}

func (me *hoverGestureRecognizer) AccpetGesture(pointer int64) {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}
	me.invokeHoverEnter(me.initEvent)
	me.invokeHover(me.initEvent)
	me.reset()
}

func (me *hoverGestureRecognizer) Clear(widget Widget) {
	me.reset()
	me.callbacks = [][]GestureCallback{[]GestureCallback{}, []GestureCallback{}, []GestureCallback{}}
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

func (me *hoverGestureRecognizer) invokeHoverEnter(event PointerEvent) {
	if me.triggerDown {
		return
	}

	me.triggerDown = true
	me.state = GestureState_Possible
	for _, cb := range me.callbacks[_ACTION_HOVER_ENTER] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *hoverGestureRecognizer) invokeHoverExit(event PointerEvent) {
	for _, cb := range me.callbacks[_ACTION_HOVER_EXIT] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *hoverGestureRecognizer) invokeHover(event PointerEvent) {
	for _, cb := range me.callbacks[_ACTION_HOVER] {
		cb(pointerEventToDetail(event, me.target))
	}
}
