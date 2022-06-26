// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
)

type HitTestResult interface {
	Add(widget Widget)
	Remove(widget Widget)
	Contain(widget Widget) bool
	Results() []Widget
}

type hitTestResult struct {
	widgets []Widget
}

func NewHitTestResult() HitTestResult {
	h := &hitTestResult{
		widgets: []Widget{},
	}
	return h
}

func (me *hitTestResult) Add(widget Widget) {
	if debug_hittest {
		for _, w := range me.widgets {
			if w == widget {
				log.Fatal("hitTest", "The widget is already added.")
			}
		}
	}

	me.widgets = append(me.widgets, widget)
}

func (me *hitTestResult) Remove(widget Widget) {
	for i, w := range me.widgets {
		if w == widget {
			me.widgets = append(me.widgets[:i], me.widgets[i+1:]...)
			return
		}
	}

	if debug_hittest {
		log.Fatal("hitTest", "The widget is not exist.")
	}
}

func (me *hitTestResult) Contain(widget Widget) bool {
	for _, w := range me.widgets {
		if w == widget {
			return true
		}
	}
	return false
}

func (me *hitTestResult) Results() []Widget {
	return me.widgets
}

//------------------ hitTestResultManager --------------------------------

var hitTestResultManagerInstance *hitTestResultManager = &hitTestResultManager{hitTestResults: map[int64]HitTestResult{}}

type hitTestResultManager struct {
	hitTestResults map[int64]HitTestResult
}

func (me *hitTestResultManager) handlePointerEvent(widget Widget, event PointerEvent) {
	var hitTestResult HitTestResult

	switch event.Action() {
	case Action_Down:
		if _, ok := me.hitTestResults[event.Pointer()]; ok {
			log.Fatal("nuxui", "hitTestResult is already exist for event %s", event)
		}

		hitTestResult = NewHitTestResult()
		me.hitTest(widget, hitTestResult, event)
		me.hitTestResults[event.Pointer()] = hitTestResult
	case Action_Up:
		hitTestResult = me.hitTestResults[event.Pointer()]
		//TODO:: delete(me.hitTestResults, event.Pointer())
	case Action_Scroll:
		log.E("nuxui", "can not run here")
	default:
		hitTestResult = me.hitTestResults[event.Pointer()]
	}

	if event.Action() == Action_Hover {
		me.handleHoverEvent(widget, event)
	} else if hitTestResult != nil {
		me.dispatchEvent(event, hitTestResult)
	}
}

func (me *hitTestResultManager) handleScrollEvent(widget Widget, event ScrollEvent) {
	handleScrollEvent(event)
}

func (me *hitTestResultManager) handleHoverEvent(widget Widget, event PointerEvent) bool {
	if c, ok := widget.(Component); ok {
		return me.handleHoverEvent(c.Content(), event)
	}

	if s, ok := widget.(Size); ok {
		frame := s.Frame()
		if event.X() >= float32(frame.X) && event.X() <= float32(frame.X+frame.Width) &&
			event.Y() >= float32(frame.Y) && event.Y() <= float32(frame.Y+frame.Height) {

			if p, ok := widget.(Parent); ok {
				for _, child := range p.Children() {
					if me.handleHoverEvent(child, event) {
						return true
					}
				}
			}

			if HoverGestureManager().existHoverCallback(widget) {
				HoverGestureManager().invokeHoverEvent(widget, event)
				return true
			}
		}
	}
	return false
}

func (me *hitTestResultManager) dispatchEvent(event PointerEvent, hitTestResult HitTestResult) {
	for _, w := range hitTestResult.Results() {
		if rs := GestureManager().getGestureRecognizers(w); rs != nil {
			for _, r := range rs {
				if r.PointerAllowed(event) {
					r.HandleAllowedPointer(event)
				}
			}
		}
	}

	switch event.Action() {
	case Action_Down:
		GestureArenaManager().Close(event.Pointer())
	case Action_Up:
		GestureArenaManager().Sweep(event.Pointer())
	}
}

type HitTestable interface {
	HitTest(widget Widget, hitTestResult HitTestResult, event PointerEvent)
}

func (me *hitTestResultManager) hitTest(widget Widget, hitTestResult HitTestResult, event PointerEvent) {
	if hit, ok := widget.(HitTestable); ok {
		hit.HitTest(widget, hitTestResult, event)
	} else {
		if c, ok := widget.(Component); ok {
			me.hitTest(c.Content(), hitTestResult, event)
		}

		if s, ok := widget.(Size); ok {
			frame := s.Frame()
			if event.X() >= float32(frame.X) && event.X() <= float32(frame.X+frame.Width) &&
				event.Y() >= float32(frame.Y) && event.Y() <= float32(frame.Y+frame.Height) {

				if p, ok := widget.(Parent); ok {
					for _, child := range p.Children() {
						me.hitTest(child, hitTestResult, event)
					}
				}

				if rs := GestureManager().getGestureRecognizers(widget); rs != nil {
					hitTestResult.Add(widget)
				}
			}
		}
	}
}
