// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

const (
	windowActionCreating = iota
	windowActionCreated
	windowActionMeasure
	windowActionDraw
	windowActionDestory
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
	handlePointerEvent(event Event)
}

type WindowDelegate interface {
	AssignOwner(owner Window)
	Owner() Window
}

type windowDelegate_HandlePointerEvent interface {
	HandlePointerEvent(event Event)
}

func NewWindow() Window {
	return newWindow()
}

type decorGestureHandler struct {
}

func (me *decorGestureHandler) AddGestureRecoginer(recognizer GestureRecognizer) {
}

func (me *decorGestureHandler) RemoveGestureRecoginer(recognizer GestureRecognizer) {
}

func (me *decorGestureHandler) HandlePointerEvent(event Event) {
	// log.V("decorGestureHandler", "####### decorGestureHandler HandlePointerEvent")
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
