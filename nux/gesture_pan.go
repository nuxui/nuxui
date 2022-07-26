// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

func OnPanDown(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_DOWN, callback)
}

func OnPanUpdate(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func OnPanUp(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_UP, callback)
}

func OnPanCancel(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_PAN_CANCEL, callback)
}

func RemovePanDownGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_DOWN, callback)
}

func RemovePanUpdateGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_UPDATE, callback)
}

func RemovePanUpGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_UP, callback)
}

func RemovePanCancelGesture(widget Widget, callback GestureCallback) {
	removePanCallback(widget, _ACTION_PAN_CANCEL, callback)
}

func addPanCallback(widget Widget, which int, callback GestureCallback) {
	widget = widget.Info().Self
	if r := GestureManager().FindGestureRecognizer(widget, (*panGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*panGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newPanGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureManager().AddGestureRecognizer(widget, recognizer)
	}
}

func removePanCallback(widget Widget, which int, callback GestureCallback) {
	widget = widget.Info().Self
	if r := GestureManager().FindGestureRecognizer(widget, (*panGestureRecognizer)(nil)); r != nil {
		pan := r.(*panGestureRecognizer)
		pan.removeCallback(which, callback)
	}
}

///////////////////////////// PanGestureRecognizer   /////////////////////////////
const (
	_ACTION_PAN_DOWN = iota
	_ACTION_PAN_UPDATE
	_ACTION_PAN_UP
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
		callbacks: [][]GestureCallback{{}, {}, {}, {}},
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
	if len(me.callbacks[_ACTION_PAN_DOWN]) == 0 &&
		len(me.callbacks[_ACTION_PAN_UPDATE]) == 0 &&
		len(me.callbacks[_ACTION_PAN_UP]) == 0 &&
		len(me.callbacks[_ACTION_PAN_CANCEL]) == 0 {
		return false
	}

	if event.IsPrimary() {
		return true
	}

	return false
}

func (me *panGestureRecognizer) HandleAllowedPointer(event PointerEvent) {
	// log.I("nuxui", "panGestureRecognizer HandleAllowedPointer %s", event.Action())
	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		GestureArenaManager().Add(event.Pointer(), me)
		me.state = GestureState_Possible
	case Action_Drag:
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
			me.invokePanUp(event)
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
		me.invokePanDown(me.initEvent)
	}
}

func (me *panGestureRecognizer) Clear(widget Widget) {
	if me.initEvent != nil {
		GestureArenaManager().Resolve(me.initEvent.Pointer(), me, false)
	}
	me.reset()
	me.callbacks = [][]GestureCallback{{}, {}, {}, {}}
}

func (me *panGestureRecognizer) reset() {
	me.initEvent = nil
	me.state = GestureState_Ready
}

func (me *panGestureRecognizer) invokePanDown(event PointerEvent) {
	// log.V("nuxui", "invokePanDown")
	for _, cb := range me.callbacks[_ACTION_PAN_DOWN] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanUpdate(event PointerEvent) {
	// log.V("nuxui", "invokePanUpdate")
	for _, cb := range me.callbacks[_ACTION_PAN_UPDATE] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanUp(event PointerEvent) {
	// log.V("nuxui", "invokePanUp")
	for _, cb := range me.callbacks[_ACTION_PAN_UP] {
		cb(pointerEventToDetail(event, me.target))
	}
}

func (me *panGestureRecognizer) invokePanCancel() {
	// log.V("nuxui", "invokePanCancel")
	for _, cb := range me.callbacks[_ACTION_PAN_CANCEL] {
		cb(pointerEventToDetail(nil, me.target))
	}
}
