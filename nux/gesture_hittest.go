// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

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
	if event.Type() != Type_PointerEvent {
		return
	}

	var hitTestResult HitTestResult

	switch event.Action() {
	case Action_Down:
		if _, ok := me.hitTestResults[event.Pointer()]; ok {
			log.Fatal("nuxui", "hitTestResult is already exist for event %s", event)
		}

		hitTestResult = NewHitTestResult()
		me.hitTest(widget, 0, 0, hitTestResult, event)
		me.hitTestResults[event.Pointer()] = hitTestResult
		log.V("hitTest", "len = %d", len(me.hitTestResults[event.Pointer()].Results()))
	case Action_Up:
		hitTestResult = me.hitTestResults[event.Pointer()]
		delete(me.hitTestResults, event.Pointer())
	case Action_Scroll:
		log.E("nuxui", "can not run here")
	default:
		hitTestResult = me.hitTestResults[event.Pointer()]
	}

	if hitTestResult != nil || event.Action() == Action_Hover {
		me.dispatchEvent(event, hitTestResult)
	}
}

func (me *hitTestResultManager) handleScrollEvent(widget Widget, event ScrollEvent) {
	handleScrollEvent(event)
}

func (me *hitTestResultManager) dispatchEvent(event PointerEvent, hitTestResult HitTestResult) {
	if hitTestResult == nil {
		// route event
		return
	}

	for _, w := range hitTestResult.Results() {
		if rs := GestureManager().getGestureRecognizers(w); rs != nil {
			for _, r := range rs {
				r.HandlePointerEvent(event)
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

func (me *hitTestResultManager) hitTest(widget Widget, offsetX, offsetY int32, hitTestResult HitTestResult, event PointerEvent) {
	if s, ok := widget.(Size); ok {
		frame := s.Frame()

		if p, ok := widget.(Parent); ok {
			for _, child := range p.Children() {
				if compt, ok := child.(Component); ok {
					child = compt.Content()
				}

				me.hitTest(child, offsetX+s.ScrollX(), offsetY+s.ScrollY(), hitTestResult, event)
			}
		}

		log.V("nuxui", "hitTest %T, '%s', ex=%f,ey=%f, x=%d, y=%d, r=%d, b=%d", widget, widget.Info().ID, event.X(), event.Y(),
			frame.X, frame.Y, frame.X+frame.Width, frame.Y+frame.Height)
		// is event in widget
		if event.X() >= float32(frame.X) && event.X() <= float32(frame.X+frame.Width) &&
			event.Y() >= float32(frame.Y) && event.Y() <= float32(frame.Y+frame.Height) {
			log.V("nuxui", "event in %T, '%s'", widget, widget.Info().ID)
			if rs := GestureManager().getGestureRecognizers(widget); rs != nil {
				log.V("nuxui", "hitTestResult.Add %T, '%s'", widget, widget.Info().ID)
				hitTestResult.Add(widget)
			}
		}
	} else if c, ok := widget.(Component); ok {
		me.hitTest(c.Content(), offsetX, offsetY, hitTestResult, event)
	}
}
