// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

const (
	windowActionCreating = iota
	windowActionCreated  = iota
	windowActionMeasure  = iota
	windowActionDraw     = iota
	windowActionDestory  = iota
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
}

func NewWindow() Window {
	return newWindow()
}
