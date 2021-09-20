// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/util"
)

func OnPanStart(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_START, callback)
}

func OnPanUpdate(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func OnPanEnd(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_END, callback)
}

func RemovePanStartGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_START, callback)
}

func RemovePanUpdateGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func RemovePanEndGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_END, callback)
}

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
	// _ACTION_PAN_DOWN = iota
	_ACTION_PAN_START = iota
	_ACTION_PAN_UPDATE
	_ACTION_PAN_END
	_ACTION_PAN_CANCEL
)

type panGestureRecognizer struct {
	target    Widget
	callbacks [][]GestureCallback
	initEvent PointerEvent
	state     GestureState
}

func newPanGestureRecognizer(target Widget) *panGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &panGestureRecognizer{
		target:    target,
		callbacks: [][]GestureCallback{[]GestureCallback{}, []GestureCallback{}, []GestureCallback{}, []GestureCallback{}},
		initEvent: nil,
		state:     GestureState_Ready,
	}
}

func (me *panGestureRecognizer) addCallback(which int, callback GestureCallback) {
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

func (me *panGestureRecognizer) removeCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	for i, cb := range me.callbacks[which] {
		if util.SameFunc(cb, callback) {
			me.callbacks[which] = append(me.callbacks[which][:i], me.callbacks[which][i+1:]...)
		}
	}
}

func (me *panGestureRecognizer) PointerAllowed(event PointerEvent) bool {
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

func (me *panGestureRecognizer) HandlePointerEvent(event PointerEvent) {
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

func (me *panGestureRecognizer) invokePanStart(event PointerEvent) {
	log.V("nuxui", "invokePanStart")
	for _, cb := range me.callbacks[_ACTION_PAN_START] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanUpdate(event PointerEvent) {
	log.V("nuxui", "invokePanUpdate")
	for _, cb := range me.callbacks[_ACTION_PAN_UPDATE] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanEnd(event PointerEvent) {
	log.V("nuxui", "invokePanEnd")
	for _, cb := range me.callbacks[_ACTION_PAN_END] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanCancel() {
	log.V("nuxui", "invokePanCancel")
	for _, cb := range me.callbacks[_ACTION_PAN_CANCEL] {
		cb(pointerEventToDetail(nil, me.target))
	}
}
