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

// TODO:: widget addMixins auto remove when onDestory,
// TODO:: When add, judge whether callback has been added or use map[callback]struct{}
func OnLongPress(widget Widget, callback func(Widget)) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS, callback)
}

func OnLongPressDown(widget Widget, callback func(Widget)) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_DOWN, callback)
}

func OnLongPressMove(widget Widget, callback func(Widget)) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_DOWN, callback)
}

func OnLongPressUp(widget Widget, callback func(Widget)) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_UP, callback)
}

func OnLongPressCancel(widget Widget, callback func(Widget)) {
	addLongPressCallback(widget, _ACTION_LONG_PRESS_CANCEL, callback)
}

// func RemoveLongPressGesture(widget Widget, callback func(Widget)) {
// 	removeLongPressCallback(widget, _ACTION_LONG_PRESS, callback)
// }

// func RemoveLongPressDownGesture(widget Widget, callback func(Widget)) {
// 	removeLongPressCallback(widget, _ACTION_LONG_PRESS_DOWN, callback)
// }

// func RemoveLongPressMoveGesture(widget Widget, callback func(Widget)) {
// 	removeLongPressCallback(widget, _ACTION_LONG_PRESS_Move, callback)
// }

// func RemoveLongPressUpGesture(widget Widget, callback func(Widget)) {
// 	removeLongPressCallback(widget, _ACTION_LONG_PRESS_UP, callback)
// }

// func RemoveLongPressCancelGesture(widget Widget, callback func(Widget)) {
// 	removeLongPressCallback(widget, _ACTION_LONG_PRESS_CANCEL, callback)
// }

func addLongPressCallback(widget Widget, which int, callback func(Widget)) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*longPressGestureRecognizer)(nil)); r != nil {
		press := r.(*longPressGestureRecognizer)
		press.addCallback(which, callback)
	} else {
		press := newLongPressGestureRecognizer(widget)
		press.addCallback(which, callback)
		h := NewGestureHandler()
		h.AddGestureRecoginer(press)
		GestureBinding().AddGestureHandler(widget, h)
	}
}

func removeLongPressCallback(widget Widget, which int, callback func(Widget)) {
	if r := GestureBinding().FindGestureRecognizer(widget, (*longPressGestureRecognizer)(nil)); r != nil {
		press := r.(*longPressGestureRecognizer)
		press.removeCallback(which, callback)
	} else {
		if true /*TODO debug*/ {
			log.Fatal("nux", "callback is not existed, maybe already removed.")
		}
	}
}

const (
	_ACTION_LONG_PRESS_DOWN = iota
	_ACTION_LONG_PRESS
	_ACTION_LONG_PRESS_UP
	_ACTION_LONG_PRESS_Move
	_ACTION_LONG_PRESS_CANCEL
)

const _LONG_PRESS_TIMEOUT = 500

type longPressGestureRecognizer struct {
	callbacks [][]unsafe.Pointer
	initEvent PointerEvent
	target    Widget
	timer     Timer
	state     GestureState
	accepted  bool
}

func newLongPressGestureRecognizer(target Widget) *longPressGestureRecognizer {
	if target == nil {
		log.Fatal("nux", "target can not be nil")
	}

	return &longPressGestureRecognizer{
		callbacks: [][]unsafe.Pointer{[]unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}, []unsafe.Pointer{}},
		// initEvent: PointerEvent{Pointer: 0},
		target:   target,
		state:    GestureState_Ready,
		accepted: false,
	}
}

func (me *longPressGestureRecognizer) addCallback(which int, callback func(Widget)) {
	log.V("nux", "longPressGestureRecognizer addCallback")
	if callback == nil {
		return
	}

	p := unsafe.Pointer(&callback)
	for _, o := range me.callbacks[which] {
		if o == p {
			log.Fatal("nux", fmt.Sprintf("The %s callback is already existed.", []string{"OnLongPressDown", "OnLongPressUp", "OnLongPress"}[which]))

		}
	}
	log.V("nux", "longPressGestureRecognizer addCallback end")
	me.callbacks[which] = append(me.callbacks[which], unsafe.Pointer(&callback))
}

func (me *longPressGestureRecognizer) removeCallback(which int, callback func(Widget)) {
	if callback == nil {
		return
	}

	log.V("nux", "longPressGestureRecognizer removeCallback end")
	p := unsafe.Pointer(&callback)
	for i, o := range me.callbacks[which] {
		if o == p {
			me.callbacks[which] = append(me.callbacks[which][:i], me.callbacks[which][i+1:]...)
		}
	}
}

func (me *longPressGestureRecognizer) PointerAllowed(event PointerEvent) bool {
	if len(me.callbacks[_ACTION_LONG_PRESS]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_DOWN]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_UP]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_Move]) == 0 &&
		len(me.callbacks[_ACTION_LONG_PRESS_CANCEL]) == 0 {
		return false
	}

	return true
}

func (me *longPressGestureRecognizer) HandlePointerEvent(event PointerEvent) {
	log.V("nux", "### longPressGestureRecognizer HandlePointerEvent")

	if event.IsPrimary() {
		switch event.Action() {
		case Action_Down:
			me.initEvent = event
			GestureArenaManager().Add(event.Pointer(), me)

			p := me.initEvent.Pointer()
			me.timer = NewTimerBackToUI(_LONG_PRESS_TIMEOUT*time.Millisecond, func() {
				me.doLongPressDown(p)
			})
		case Action_Move:
			if me.state == GestureState_Accepted {
				me.invokeLongPressMove(event.Pointer())
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
}

func (me *longPressGestureRecognizer) RejectGesture(pointer int64) {
	if me.state == GestureState_Possible {
		me.invokeLongPressUCancel(pointer)
	}
	me.reset()
}

func (me *longPressGestureRecognizer) AccpetGesture(pointer int64) {
	if me.timer != nil {
		me.timer.Stop()
		me.timer = nil
	}
	me.state = GestureState_Accepted
	me.invokeLongPressDown(pointer)
}

func (me *longPressGestureRecognizer) reset() {
	if me.timer != nil {
		me.timer.Stop()
		me.timer = nil
	}

	me.state = GestureState_Ready
	me.initEvent = nil
}

func (me *longPressGestureRecognizer) doLongPressDown(pointer int64) {
	log.V("nux", "doLongPressDown")
	me.state = GestureState_Possible
	GestureArenaManager().Resolve(pointer, me, true)
}

func (me *longPressGestureRecognizer) invokeLongPressDown(pointer int64) {
	log.V("nux", "invokeLongPressDown")
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_LONG_PRESS_DOWN] {
			(*(*(func(Widget)))(c))(me.target)
		}

		for _, c := range me.callbacks[_ACTION_LONG_PRESS] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *longPressGestureRecognizer) invokeLongPressMove(pointer int64) {
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_LONG_PRESS_Move] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *longPressGestureRecognizer) invokeLongPressUp(pointer int64) {
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_LONG_PRESS_UP] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}

func (me *longPressGestureRecognizer) invokeLongPressUCancel(pointer int64) {
	if me.initEvent.Pointer() == pointer {
		for _, c := range me.callbacks[_ACTION_LONG_PRESS_CANCEL] {
			(*(*(func(Widget)))(c))(me.target)
		}
	}
}
