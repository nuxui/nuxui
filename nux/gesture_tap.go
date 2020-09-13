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

func OnTap(widget Widget, callback func(Widget)) {
	addTapCallback(widget, _ACTION_TAP, callback)
}

func OnTapDown(widget Widget, callback func(Widget)) {
	addTapCallback(widget, _ACTION_TAP_DOWN, callback)
}

func OnTapUp(widget Widget, callback func(Widget)) {
	addTapCallback(widget, _ACTION_TAP_UP, callback)
}

func OnTapCancel(widget Widget, callback func(Widget)) {
	addTapCallback(widget, _ACTION_TAP_CANCEL, callback)
}

// widget will auto clear all gesture when destory
// func RemoveTapGesture(widget Widget, callback func(Widget)) {
// 	removeTapCallback(widget, _ACTION_TAP, callback)
// }

// func RemoveTapDownGesture(widget Widget, callback func(Widget)) {
// 	removeTapCallback(widget, _ACTION_TAP_DOWN, callback)
// }

// func RemoveTapUpGesture(widget Widget, callback func(Widget)) {
// 	removeTapCallback(widget, _ACTION_TAP_UP, callback)
// }

// func RemoveTapCancelGesture(widget Widget, callback func(Widget)) {
// 	removeTapCallback(widget, _ACTION_TAP_CANCEL, callback)
// }

func addTapCallback(widget Widget, which int, callback func(Widget)) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		tap := r.(*tapGestureRecognizer)
		tap.addCallback(which, callback)
	} else {
		tap := newTapGestureRecognizer(widget).(*tapGestureRecognizer)
		tap.addCallback(which, callback)
		h := NewGestureHandler()
		h.AddGestureRecoginer(tap)
		GestureBinding().AddGestureHandler(widget, h)
	}
}

func removeTapCallback(widget Widget, which int, callback func(Widget)) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*tapGestureRecognizer)(nil)); r != nil {
		tap := r.(*tapGestureRecognizer)
		tap.removeCallback(which, callback)
	} /*else {
		// ignore
		log.V("nux", "callback is not existed, maybe already removed.")
	}*/
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
	initEvent   PointerEvent
	target      Widget
	timer       Timer
	state       GestureState
	triggerDown bool
}

func newTapGestureRecognizer(target Widget) TapGestureRecognizer {
	if target == nil {
		log.Fatal("nux", "target can not be nil")
	}

	return &tapGestureRecognizer{
		callbacks: [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		// initEvent:   PointerEvent{Pointer: 0},
		target:      target,
		state:       GestureState_Ready,
		triggerDown: false,
	}
}

func (me *tapGestureRecognizer) addCallback(which int, callback func(Widget)) {
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.callbacks[which] {
		if o == p {
			if true /*TODO:: debug*/ {
				log.Fatal("nux", fmt.Sprintf("The %s callback is already existed.", []string{"OnTapDown", "OnTapUp", "OnTap"}[which]))

			} else {
				return
			}
		}
	}

	me.callbacks[which] = append(me.callbacks[which], unsafe.Pointer(&callback))
}

func (me *tapGestureRecognizer) removeCallback(which int, callback func(Widget)) {
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

func (me *tapGestureRecognizer) PointerAllowed(event PointerEvent) bool {
	if len(me.callbacks[_ACTION_TAP]) == 0 &&
		len(me.callbacks[_ACTION_TAP_DOWN]) == 0 &&
		len(me.callbacks[_ACTION_TAP_UP]) == 0 &&
		len(me.callbacks[_ACTION_TAP_CANCEL]) == 0 {
		return false
	}

	return true
}

func (me *tapGestureRecognizer) HandlePointerEvent(event PointerEvent) {
	if event.IsPrimary() {
		switch event.Action() {
		case Action_Down:
			log.V("nux", "IsPrimaryButton Action_Down")
			me.initEvent = event
			GestureArenaManager().Add(event.Pointer(), me)

			p := me.initEvent.Pointer()
			log.V("nux", "IsPrimaryButton Action_Down NewTimer")
			me.timer = NewTimerBackToUI(DOWN_DELAY*time.Millisecond, func() {
				me.invokeTapDown(p)
			})
		case Action_Up:
			GestureArenaManager().Resolve(me.initEvent.Pointer(), me, me.initEvent.Pointer() == event.Pointer())
		case Action_Move:
			log.V("nux", "IsPrimaryButton Action_Move")
			if me.state == GestureState_Possible {
				if me.initEvent.Distance(event.X(), event.Y()) >= MIN_PAN_DISTANCE {
					GestureArenaManager().Resolve(event.Pointer(), me, false)
				}
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
		me.timer.Stop()
		me.timer = nil
	}
	me.invokeTapDown(pointer)
	me.invokeTapUpAndTap(pointer)
	me.reset()
}

func (me *tapGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Stop()
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
	for _, c := range me.callbacks[_ACTION_TAP_CANCEL] {
		(*(*(func(Widget)))(c))(me.target)
	}
}

func (me *tapGestureRecognizer) invokeTapUpAndTap(pointer int64) {
	if me.initEvent.Pointer() != pointer {
		return
	}

	for _, c := range me.callbacks[_ACTION_TAP_UP] {
		(*(*(func(Widget)))(c))(me.target)
	}

	for _, c := range me.callbacks[_ACTION_TAP] {
		(*(*(func(Widget)))(c))(me.target)
	}
}
