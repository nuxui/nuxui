// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"unsafe"

	"github.com/nuxui/nuxui/log"
)

func OnPanDown(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_DOWN, callback)
}

func OnPanStart(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_START, callback)
}

func OnPanUpdate(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func OnPanEnd(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_END, callback)
}

// func RemovePanDownGesture(widget Widget, callback GestureCallback) {
// 	removePanCallback(widget, _ACTION_PAN_DOWN, callback)
// }

// func RemovePanStartGesture(widget Widget, callback GestureCallback) {
// 	removePanCallback(widget, _ACTION_PAN_START, callback)
// }

// func RemovePanUpdateGesture(widget Widget, callback GestureCallback) {
// 	removePanCallback(widget, _ACTION_PAN_UPDATE, callback)
// }

// func RemovePanEndGesture(widget Widget, callback GestureCallback) {
// 	removePanCallback(widget, _ACTION_PAN_END, callback)
// }

func addPanCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*panGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*panGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newPanGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureBinding().AddGestureRecognizer(widget, recognizer)
	}
}

func removePanCallback(widget Widget, which int, callback GestureCallback) {
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
	_ACTION_PAN_DOWN = iota
	_ACTION_PAN_UPDATE
	_ACTION_PAN_START
	_ACTION_PAN_END
	_ACTION_PAN_CANCEL
)

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
		callbacks: [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		// initEvent: Event{Pointer: 0},
		state: GestureState_Ready,
	}
}

func (me *panGestureRecognizer) addCallback(which int, callback GestureCallback) {
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

func (me *panGestureRecognizer) removeCallback(which int, callback GestureCallback) {
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
	if len(me.callbacks[_ACTION_PAN_DOWN]) == 0 &&
		len(me.callbacks[_ACTION_PAN_START]) == 0 &&
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
			me.invokePanEnd()
			me.reset()
		} else if me.state == GestureState_Possible {
			GestureArenaManager().Resolve(event.Pointer(), me, false)
		}
	}
}

func (me *panGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.invokePanCancel()
	}
	me.reset()
}

func (me *panGestureRecognizer) AccpetGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.state = GestureState_Accepted
		me.invokePanDown(pointer)
		me.invokePanStart(pointer)
	}
}

func (me *panGestureRecognizer) reset() {
	me.initEvent = nil
	me.state = GestureState_Ready
}

func (me *panGestureRecognizer) invokePanDown(pointer int64) {
	log.V("nuxui", "invokePanDown")
	for _, c := range me.callbacks[_ACTION_PAN_DOWN] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *panGestureRecognizer) invokePanStart(pointer int64) {
	log.V("nuxui", "invokePanStart")
	for _, c := range me.callbacks[_ACTION_PAN_START] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *panGestureRecognizer) invokePanUpdate(event Event) {
	log.V("nuxui", "invokePanUpdate")
	for _, c := range me.callbacks[_ACTION_PAN_UPDATE] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *panGestureRecognizer) invokePanEnd() {
	log.V("nuxui", "invokePanEnd")
	for _, c := range me.callbacks[_ACTION_PAN_END] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *panGestureRecognizer) invokePanCancel() {
	log.V("nuxui", "invokePanCancel")
	for _, c := range me.callbacks[_ACTION_PAN_CANCEL] {
		(*(*(func(Widget)))(c))(me.target)
	}
}
