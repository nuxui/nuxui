// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

const (
	windowActionCreated = iota
	windowActionMeasure
	windowActionDraw
	windowActionDestroy
)

type Window interface {
	ID() uint64
	// Surface() Surface
	LockCanvas() (Canvas, error)
	UnlockCanvas() error
	// OnInputEvent()
	Decor() Widget
	Width() int32
	Height() int32
	ContentWidth() int32
	ContentHeight() int32
	Alpha() float32
	SetAlpha(alpha float32)
	Title() string
	SetTitle(title string)
	SetDelegate(delegate WindowDelegate)
	Delegate() WindowDelegate

	// private methods
	handlePointerEvent(event PointerEvent)
	handleScrollEvent(event ScrollEvent)
	handleKeyEvent(event KeyEvent)
	handleTypeEvent(event TypeEvent)
	requestFocus(widget Widget)
}

type WindowDelegate interface {
	AssignOwner(owner Window)
	Owner() Window
}

type windowDelegate_HandlePointerEvent interface {
	HandlePointerEvent(event PointerEvent)
}

func NewWindow(attr Attr) Window {
	return newWindow(attr)
}

func (me *window) OnCreate() {
	main := App().Manifest().Main()
	if main == "" {
		log.Fatal("nuxui", "no main widget found.")
	} else {
		mainWidgetCreator := FindRegistedWidgetCreatorByName(main)
		ctx := &context{}
		// widgetTree := RenderWidget(mainWidgetCreator(ctx, Attr{}))
		log.I("nuxui", "======  mainWidgetCreator")
		widgetTree := mainWidgetCreator(ctx, Attr{})
		log.I("nuxui", "======  decor AddChild  widgetTree begin")
		me.decor.AddChild(widgetTree)
		log.I("nuxui", "======  decor AddChild  widgetTree end")
		mountWidget(me.decor, nil)

	}

}

func (me *window) CreateDecor(ctx Context, attr Attr) Widget {
	creator := FindRegistedWidgetCreatorByName("github.com/nuxui/nuxui/ui.Layer")
	w := creator(ctx, attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

func (me *window) Measure(width, height int32) {
	if me.decor == nil {
		return
	}

	if s, ok := me.decor.(Size); ok {
		if s.MeasuredSize().Width == width && s.MeasuredSize().Height == height {
			// return
		}

		s.MeasuredSize().Width = width
		s.MeasuredSize().Height = height
	}

	me.surfaceResized = true

	if f, ok := me.decor.(Measure); ok {
		f.Measure(width, height)
	}
}

func (me *window) Layout(dx, dy, left, top, right, bottom int32) {
	if me.decor == nil {
		return
	}

	if f, ok := me.decor.(Layout); ok {
		f.Layout(dx, dy, left, top, right, bottom)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type decorGestureHandler struct {
}

func (me *decorGestureHandler) AddGestureRecoginer(recognizer GestureRecognizer) {
}

func (me *decorGestureHandler) RemoveGestureRecoginer(recognizer GestureRecognizer) {
}

func (me *decorGestureHandler) HandlePointerEvent(event PointerEvent) {
	switch event.Action() {
	case Action_Down:
		GestureArenaManager().Close(event.Pointer())
	case Action_Up:
		GestureArenaManager().Sweep(event.Pointer())
	}
}

func (me *decorGestureHandler) FindGestureRecognizer(recognizer GestureRecognizer) GestureRecognizer {
	return nil
}
