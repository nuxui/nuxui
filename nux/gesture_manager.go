// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
)

var gestureManagerInstance = &gestureManager{
	hitTestResults: map[int64]HitTestResult{},
}

type gestureManager struct {
	hitTestResults map[int64]HitTestResult
}

func (me *gestureManager) handlePointerEvent(widget Widget, event PointerEvent) {
	if event.Type() != Type_PointerEvent {
		return
	}

	var hitTestResult HitTestResult

	switch event.Action() {
	case Action_Down:
		if _, ok := me.hitTestResults[event.Pointer()]; ok {
			log.Fatal("hitTest", "hitTestResult is already exist for event %s", event)
		}

		hitTestResult = NewHitTestResult()
		me.hitTest(widget, 0, 0, hitTestResult, event)
		me.hitTestResults[event.Pointer()] = hitTestResult
		// log.V("hitTest", "len = %d", len(me.hitTestResults[event.Pointer()].Results()))
	case Action_Up:
		hitTestResult = me.hitTestResults[event.Pointer()]
		delete(me.hitTestResults, event.Pointer())
	case Action_Scroll:
		log.E("gestureManager", "can not run here")
	default:
		hitTestResult = me.hitTestResults[event.Pointer()]
	}

	if hitTestResult != nil || event.Action() == Action_Hover {
		me.dispatchEvent(event, hitTestResult)
	}
}

func (me *gestureManager) handleScrollEvent(widget Widget, event ScrollEvent) {
	handleScrollEvent(event)
}

func (me *gestureManager) dispatchEvent(event PointerEvent, hitTestResult HitTestResult) {
	if hitTestResult == nil {
		// route event
		return
	}

	for _, w := range hitTestResult.Results() {
		if h := GestureBinding().FindGestureHandler(w); h != nil {
			h.HandlePointerEvent(event)
		}

		//TODO:: if widget is not translucent, then prevent the event from passing on
		// if v, ok := w.(Visual); ok && !v.Translucent() {
		// 	return
		// }
	}
}

func (me *gestureManager) hitTest(widget Widget, offsetX, offsetY int32, hitTestResult HitTestResult, event PointerEvent) {
	if s, ok := widget.(Size); ok {
		ms := s.MeasuredSize()

		if p, ok := widget.(Parent); ok {
			for _, child := range p.Children() {
				if compt, ok := child.(Component); ok {
					child = compt.Content()
				}

				me.hitTest(child, offsetX+s.ScrollX(), offsetY+s.ScrollY(), hitTestResult, event)
			}
		}

		// log.V("hitTest", "id = %s", widget.ID())

		// is event in widget
		if event.X() >= float32(ms.Position.X+offsetX) && event.X() <= float32(ms.Position.X+offsetX+ms.Width) &&
			event.Y() >= float32(ms.Position.Y+offsetY) && event.Y() <= float32(ms.Position.Y+offsetY+ms.Height) {
			if h := GestureBinding().FindGestureHandler(widget); h != nil {
				hitTestResult.Add(widget)
			}
		}
	}
}
