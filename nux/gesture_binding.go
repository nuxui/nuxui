// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

var gestureBindingInstance *gestureBinding = &gestureBinding{widgetGestureHandlers: map[Widget]GestureHandler{}}

func GestureBinding() *gestureBinding {
	return gestureBindingInstance
}

type gestureBinding struct {
	widgetGestureHandlers map[Widget]GestureHandler
}

func (me *gestureBinding) FindGestureHandler(widget Widget) GestureHandler {
	return me.widgetGestureHandlers[widget]
}

func (me *gestureBinding) AddGestureHandler(widget Widget, handler GestureHandler) {
	me.widgetGestureHandlers[widget] = handler
}

//ClearGestureHandler clear widget gesture handler when destory
func (me *gestureBinding) ClearGestureHandler(widget Widget) {
	delete(me.widgetGestureHandlers, widget)
}

func (me *gestureBinding) AddGestureRecognizer(widget Widget, recognizer GestureRecognizer) {
	if handler, ok := me.widgetGestureHandlers[widget]; ok {
		handler.AddGestureRecoginer(recognizer)
	} else {
		h := NewGestureHandler()
		h.AddGestureRecoginer(recognizer)
		me.widgetGestureHandlers[widget] = h
	}
}

// recognizer
func (me *gestureBinding) FindGestureRecognizer(widget Widget, recognizerType GestureRecognizer) GestureRecognizer {
	if h, ok := me.widgetGestureHandlers[widget]; ok {
		return h.FindGestureRecognizer(recognizerType)
	}

	return nil
}
