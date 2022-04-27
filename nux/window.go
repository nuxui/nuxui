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
	LockCanvas() Canvas
	UnlockCanvas(Canvas)
	Decor() Widget
	Size() (width, height int32)
	ContentSize() (width, height int32)
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
		mainWidgetCreator := FindRegistedWidgetCreator(main)
		widgetTree := mainWidgetCreator(Attr{})
		me.decor.AddChild(widgetTree)
		mountWidget(me.decor, nil)
	}

}

func (me *window) Measure(width, height MeasureDimen) {
	if me.decor == nil {
		return
	}

	if s, ok := me.decor.(Size); ok {
		f := s.Frame()
		if f.Width == width.Value() && f.Height == height.Value() {
			// return
		}

		f.Width = width.Value()
		f.Height = height.Value()
	}

	if f, ok := me.decor.(Measure); ok {
		f.Measure(width, height)
	}
}

func (me *window) Layout(x, y, width, height int32) {
	if me.decor == nil {
		return
	}

	if f, ok := me.decor.(Layout); ok {
		f.Layout(x, y, width, height)
	}
}
