// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/util"
)

func OnTap(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionTap, callback)
}

func OnTapDown(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionTapDown, callback)
}

func OnTapUp(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionTapUp, callback)
}

func OnTapCancel(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionTapCancel, callback)
}

func OnSecondaryTap(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionSecondaryTap, callback)
}

func OnSecondaryTapDown(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionSecondaryTapDown, callback)
}

func OnSecondaryTapUp(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionSecondaryTapUp, callback)
}

func OnSecondaryTapCancel(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionSecondaryTapCancel, callback)
}

func OnOtherTap(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionOtherTap, callback)
}

func OnOtherTapDown(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionOtherTapDown, callback)
}

func OnOtherTapUp(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionOtherTapUp, callback)
}

func OnOtherTapCancel(widget Widget, callback GestureCallback) {
	addTapCallback(widget, actionOtherTapCancel, callback)
}

/*
func RemoveTapGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, actionTap, callback)
}

func RemoveTapDownGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, actionTapDown, callback)
}

func RemoveTapUpGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, actionTapUp, callback)
}

func RemoveTapCancelGesture(widget Widget, callback GestureCallback) {
	removeTapCallback(widget, actionTapCancel, callback)
}*/

func addTapCallback(widget Widget, which int, callback GestureCallback) {
	widget = widget.Info().Self
	if r := GestureManager().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*tapGestureRecognizer)
		recognizer.addCallback(which, callback)
	} else {
		recognizer := newTapGestureRecognizer(widget)
		recognizer.addCallback(which, callback)
		GestureManager().AddGestureRecognizer(widget, recognizer)
	}
}

/*
func removeTapCallback(widget Widget, which int, callback GestureCallback) {
	if r := GestureManager().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		recognizer := r.(*tapGestureRecognizer)
		recognizer.removeCallback(which, callback)
	}
}
*/

const (
	actionTapDown = iota
	actionTapUp
	actionTap
	actionTapCancel
	actionSecondaryTapDown
	actionSecondaryTapUp
	actionSecondaryTap
	actionSecondaryTapCancel
	actionOtherTapDown
	actionOtherTapUp
	actionOtherTap
	actionOtherTapCancel
)

///////////////////////////// TapGestureRecognizer   /////////////////////////////

type TapGestureRecognizer interface {
	GestureRecognizer
}

type tapGestureRecognizer struct {
	target      Widget
	callbacks   map[int][]GestureCallback
	initEvent   PointerEvent
	timer       Timer
	state       GestureState
	triggerDown bool
	action      EventAction
}

func newTapGestureRecognizer(target Widget) *tapGestureRecognizer {
	if target == nil {
		log.Fatal("nuxui", "target can not be nil")
	}

	return &tapGestureRecognizer{
		callbacks:   map[int][]GestureCallback{},
		initEvent:   nil,
		target:      target,
		timer:       nil,
		state:       GestureState_Ready,
		triggerDown: false,
		action:      Action_None,
	}
}

func (me *tapGestureRecognizer) addCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	if cbs, ok := me.callbacks[which]; ok {
		for _, cb := range cbs {
			if util.SameFunc(cb, callback) {
				log.Fatal("nuxui", "The %s callback is already existed.", []string{"OnTapDown", "OnTapUp", "OnTap"}[which])
			}
		}
		me.callbacks[which] = append(cbs, callback)
	} else {
		me.callbacks[which] = []GestureCallback{callback}
	}

}

func (me *tapGestureRecognizer) removeCallback(which int, callback GestureCallback) {
	if callback == nil {
		return
	}

	if cbs, ok := me.callbacks[which]; ok {
		for i, cb := range cbs {
			if util.SameFunc(cb, callback) {
				me.callbacks[which] = append(cbs[:i], cbs[i+1:]...)
			}
		}
	}
}

func (me *tapGestureRecognizer) PointerAllowed(event PointerEvent) bool {
	for k, _ := range me.callbacks {
		if k <= actionTapCancel {
			if event.IsPrimary() {
				return true
			}
		} else if k <= actionSecondaryTapCancel {
			if event.Button() == ButtonSecondary {
				return true
			}
		} else if k > actionSecondaryTapCancel {
			if event.Button() > ButtonSecondary {
				return true
			}
		}
	}

	return false
}

func (me *tapGestureRecognizer) HandleAllowedPointer(event PointerEvent) {
	switch event.Action() {
	case Action_Down:
		me.action = Action_Down
		me.initEvent = event
		pointer := event.Pointer()
		GestureArenaManager().Add(pointer, me)

		if event.Kind() == Kind_Touch {
			me.timer = NewTimerBackToUI(GESTURE_DOWN_DELAY, func() {
				me.invokeTapDown(pointer)
			})
		} else {
			// log.I("nuxui", "first down second=%t", event.Button() == ButtonSecondary)
			me.invokeTapDown(pointer)
		}
	case Action_Up:
		// do not accept proactive, wait GestureArea sweep
		// log.I("nuxui", "Action_Up accept=%t", me.state == GestureState_Accepted)
		me.action = Action_Up
		if me.state == GestureState_Accepted {
			me.invokeTapUpAndTap(event.Pointer())
		}

	case Action_Drag:
		me.action = Action_Drag
		if me.state == GestureState_Possible {
			if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
				GestureArenaManager().Resolve(event.Pointer(), me, false)
			}
		}
	}
}

func (me *tapGestureRecognizer) RejectGesture(pointer int64) {
	// log.I("nuxui", "tap gesture reject %s", me.target.Info().ID)
	if me.state == GestureState_Possible {
		me.invokeTapCancel()
	}

	me.reset()
}

func (me *tapGestureRecognizer) AccpetGesture(pointer int64) {
	// log.I("nuxui", "tap gesture accept %s", me.target.Info().ID)
	if me.timer != nil {
		r := me.timer.Running()
		me.timer.Cancel()
		me.timer = nil

		if r {
			me.invokeTapDown(pointer)
		}
	}

	me.state = GestureState_Accepted
	if me.action == Action_Up {
		me.invokeTapUpAndTap(pointer)
	}
}

func (me *tapGestureRecognizer) Clear(widget Widget) {
	if me.target != widget {
		log.Fatal("nuxui", "target is not matched")
	}
	if me.initEvent != nil {
		GestureArenaManager().Resolve(me.initEvent.Pointer(), me, false)
	}
	me.reset()
	me.target = nil
	me.callbacks = map[int][]GestureCallback{}
}

func (me *tapGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Cancel()
		me.timer = nil
	}

	me.state = GestureState_Ready
	me.action = Action_None
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
		which := actionOtherTapDown
		switch me.initEvent.Button() {
		case ButtonPrimary:
			which = actionTapDown
		case ButtonSecondary:
			which = actionSecondaryTapDown
		}
		for _, cb := range me.callbacks[which] {
			cb(pointerEventToDetail(me.initEvent, me.target))
		}
	}
}

func (me *tapGestureRecognizer) invokeTapCancel() {
	if me.triggerDown {
		which := actionOtherTapCancel
		switch me.initEvent.Button() {
		case ButtonPrimary:
			which = actionTapCancel
		case ButtonSecondary:
			which = actionSecondaryTapCancel
		}
		for _, cb := range me.callbacks[which] {
			cb(pointerEventToDetail(nil, me.target))
		}
	}
}

func (me *tapGestureRecognizer) invokeTapUpAndTap(pointer int64) {
	if !me.triggerDown || me.initEvent.Pointer() != pointer {
		return
	}

	whichUp := actionOtherTapUp
	whichTap := actionOtherTap
	switch me.initEvent.Button() {
	case ButtonPrimary:
		whichUp = actionTapUp
		whichTap = actionTap
	case ButtonSecondary:
		whichUp = actionSecondaryTapUp
		whichTap = actionSecondaryTap
	}

	for _, cb := range me.callbacks[whichUp] {
		cb(pointerEventToDetail(me.initEvent, me.target))
	}

	for _, cb := range me.callbacks[whichTap] {
		cb(pointerEventToDetail(me.initEvent, me.target))
	}

	me.reset()
}
