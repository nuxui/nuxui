// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

func OnLongPress(widget Widget, callback GestureCallback) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS, callback)
}

func OnLongPressMove(widget Widget, callback GestureCallback) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_Move, callback)
}

func OnLongPressUp(widget Widget, callback GestureCallback) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_UP, callback)
}

func RemoveLongPressGesture(widget Widget, callback GestureCallback) {
	removeLongPressCallback(widget, _ACTION_LONG_PRESS, callback)
}

func RemoveLongPressMoveGesture(widget Widget, callback GestureCallback) {
	removeLongPressCallback(widget, _ACTION_LONG_PRESS_Move, callback)
}

func RemoveLongPressUpGesture(widget Widget, callback GestureCallback) {
	removeLongPressCallback(widget, _ACTION_LONG_PRESS_UP, callback)
}

func addLongPressCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureManager().FindGestureRecognizer(widget, (*longPressGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*longPressGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newLongPressGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureManager().AddGestureRecognizer(widget, recognizer)
	}
}

func removeLongPressCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureManager().FindGestureRecognizer(widget, (*longPressGestureRecognizer)(nil)); r != nil {
		press := r.(*longPressGestureRecognizer)
		press.removeCallback(which, callback)
	} else {
		if debug_gesture {
			log.Fatal("nuxui", "callback is not existed, maybe already removed.")
		}
	}
}

const (
	_ACTION_LONG_PRESS = iota
	_ACTION_LONG_PRESS_UP
	_ACTION_LONG_PRESS_Move
)

type longPressGestureRecognizer struct {
	callbacks [][]GestureCallback
	initEvent PointerEvent
	target    Widget
	timer     Timer
	state     GestureState
	accepted  bool
}

func newLongPressGestureRecognizer(target Widget) *longPressGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &longPressGestureRecognizer{
		callbacks: [][]GestureCallback{[]GestureCallback{}, []GestureCallback{}, []GestureCallback{}},
		initEvent: nil,
		target:    target,
		timer:     nil,
		state:     GestureState_Ready,
		accepted:  false,
	}
}

func (me *longPressGestureRecognizer) addCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	for _, cb := range me.callbacks[which] {
		if util.SameFunc(cb, callback) {
			log.Fatal("nuxui", "The %s callback is already existed.", []string{"OnLongPress", "OnLongPressMove", "OnLongPressUp"}[which])

		}
	}
	me.callbacks[which] = append(me.callbacks[which], callback)
}

func (me *longPressGestureRecognizer) removeCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	for i, cb := range me.callbacks[which] {
		if util.SameFunc(cb, callback) {
			me.callbacks[which] = append(me.callbacks[which][:i], me.callbacks[which][i+1:]...)
		}
	}
}

func (me *longPressGestureRecognizer) PointerAllowed(event PointerEvent) bool {
	if len(me.callbacks[_ACTION_LONG_PRESS]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_UP]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_Move]) == 0 {
		return false
	}

	if event.IsPrimary() {
		return true
	}

	return false
}

func (me *longPressGestureRecognizer) HandlePointerEvent(event PointerEvent) {
	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		pointer := event.Pointer()
		GestureArenaManager().Add(pointer, me)
		me.timer = NewTimerBackToUI(GESTURE_LONG_PRESS_TIMEOUT, func() {
			me.state = GestureState_Possible
			GestureArenaManager().Resolve(pointer, me, true)
		})
	case Action_Move:
		if me.state == GestureState_Accepted {
			me.invokeLongPressMove(event.Pointer())
		} else {
			if me.initEvent != nil && me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
				GestureArenaManager().Resolve(event.Pointer(), me, false)
			}
		}
	case Action_Up:
		if me.state == GestureState_Accepted {
			me.invokeLongPressUp(event.Pointer())
			me.reset()
		} else {
			GestureArenaManager().Resolve(event.Pointer(), me, false)
		}
	}
}

func (me *longPressGestureRecognizer) RejectGesture(pointer int64) {
	me.reset()
}

func (me *longPressGestureRecognizer) AccpetGesture(pointer int64) {
	if me.state == GestureState_Possible {
		if me.timer != nil {
			me.timer.Cancel()
			me.timer = nil
		}
		me.invokeLongPress(pointer)
	} else {
		// only one gesture, but not long press
		me.reset()
	}
}

func (me *longPressGestureRecognizer) Clear(widget Widget) {
	if me.initEvent != nil {
		GestureArenaManager().Resolve(me.initEvent.Pointer(), me, false)
	}
	me.reset()
	me.callbacks = [][]GestureCallback{[]GestureCallback{}, []GestureCallback{}, []GestureCallback{}}
}

func (me *longPressGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	me.state = GestureState_Ready
	me.initEvent = nil
}

func (me *longPressGestureRecognizer) invokeLongPress(pointer int64) {
	// log.V("nuxui", "invokeLongPressDown")
	me.state = GestureState_Accepted
	if me.initEvent.Pointer() == pointer {
		for _, cb := range me.callbacks[_ACTION_LONG_PRESS] {
			cb(pointerEventToDetail(me.initEvent, me.target))
		}
	}
}

func (me *longPressGestureRecognizer) invokeLongPressMove(pointer int64) {
	if me.initEvent.Pointer() == pointer {
		for _, cb := range me.callbacks[_ACTION_LONG_PRESS_Move] {
			cb(pointerEventToDetail(me.initEvent, me.target))
		}
	}
}

func (me *longPressGestureRecognizer) invokeLongPressUp(pointer int64) {
	if me.initEvent.Pointer() == pointer {
		for _, cb := range me.callbacks[_ACTION_LONG_PRESS_UP] {
			cb(pointerEventToDetail(me.initEvent, me.target))
		}
	}
}
