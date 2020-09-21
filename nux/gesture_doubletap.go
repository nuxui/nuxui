// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

func OnDoubleTap(widget Widget, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*doubleTapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*doubleTapGestureRecognizer)
		recognizer.addCallback(_ACTION_DOUBLETAP, callback)
	} else {
		recognizer := newDoubleTapGestureRecognizer(widget)
		recognizer.addCallback(_ACTION_DOUBLETAP, callback)
		GestureBinding().AddGestureRecognizer(widget, recognizer)
	}
}

// widget will auto clear all gesture when destory
func RemoveDoubleTapGesture(widget Widget, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*doubleTapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*doubleTapGestureRecognizer)
		recognizer.removeCallback(_ACTION_DOUBLETAP, callback)
	}
}

const (
	_ACTION_DOUBLETAP = iota
)

///////////////////////////// DoubleTapGestureRecognizer   /////////////////////////////

type DoubleTapGestureRecognizer interface {
	GestureRecognizer
}

type doubleTapGestureRecognizer struct {
	callbacks          [][]unsafe.Pointer
	rejectFirstPointer int64
	firstTap           Event
	secondTap          Event
	target             Widget
	timer              Timer
}

func newDoubleTapGestureRecognizer(target Widget) *doubleTapGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &doubleTapGestureRecognizer{
		callbacks:          [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		rejectFirstPointer: 0,
		firstTap:           nil,
		secondTap:          nil,
		target:             target,
		timer:              nil,
	}
}

func (me *doubleTapGestureRecognizer) addCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.callbacks[which] {
		if o == p {
			if true /*TODO:: debug*/ {
				log.Fatal("nuxui", fmt.Sprintf("The %s callback is already existed.", []string{"OnTap2Down", "OnTap2Up", "OnTap2"}[which]))

			} else {
				return
			}
		}
	}

	me.callbacks[which] = append(me.callbacks[which], unsafe.Pointer(&callback))
}

func (me *doubleTapGestureRecognizer) removeCallback(which int, callback GestureCallback) {
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

func (me *doubleTapGestureRecognizer) PointerAllowed(event Event) bool {
	if len(me.callbacks[_ACTION_DOUBLETAP]) == 0 {
		return false
	}

	if event.IsPrimary() {
		return true
	}

	return false
}

func (me *doubleTapGestureRecognizer) HandlePointerEvent(event Event) {
	switch event.Action() {
	case Action_Down:
		GestureArenaManager().Add(event.Pointer(), me)
		// log.V("nuxui", "#########  double tap ======== %d", event.Time().UnixNano())
	case Action_Up:
		if me.rejectFirstPointer == event.Pointer() {
			GestureArenaManager().Resolve(event.Pointer(), me, false)
		} else if me.firstTap == nil {
			if GestureArenaManager().Hold(event.Pointer()) {
				me.firstTap = event
				pointer := event.Pointer()
				me.timer = NewTimerBackToUI(GESTURE_DOUBLETAP_TIMEOUT, func() {
					GestureArenaManager().Resolve(pointer, me, false)
				})
			}
		} else {
			me.secondTap = event
			GestureArenaManager().Resolve(me.firstTap.Pointer(), me, true)
			GestureArenaManager().Resolve(event.Pointer(), me, true)
		}
	case Action_Move:
		if me.firstTap != nil {
			if me.firstTap.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
				GestureArenaManager().Resolve(event.Pointer(), me, false)
			}
		}
	}
}

func (me *doubleTapGestureRecognizer) RejectGesture(pointer int64) {
	if me.firstTap == nil {
		me.rejectFirstPointer = pointer
	} else {
		me.reset()
	}
}

func (me *doubleTapGestureRecognizer) AccpetGesture(pointer int64) {
	if me.firstTap != nil && me.secondTap != nil && me.secondTap.Pointer() == pointer {
		me.invokeDoubleTap()
		me.reset()
	}
}

func (me *doubleTapGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	me.firstTap = nil
	me.secondTap = nil
	me.rejectFirstPointer = 0
}

func (me *doubleTapGestureRecognizer) invokeDoubleTap() {
	// log.V("nuxui", "#########  double tap invokeDoubleTap")
	for _, c := range me.callbacks[_ACTION_DOUBLETAP] {
		(*(*(func(Widget)))(c))(me.target)
	}
}
