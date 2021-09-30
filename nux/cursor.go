// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Cursor interface {
	X() float32
	Y() float32
}

func GetCursor() Cursor {
	// TODO:: implement
	return nil
}

// Retrieves the position of the mouse cursor, in screen coordinates.
func GetCursorScreenPosition() (x, y float32) {
	return getCursorScreenPosition()
}

// Retrieves the position of the mouse cursor, in window coordinates.
func GetCursorWindowPosition(window Window) (x, y float32) {
	x0, y0 := getCursorScreenPosition()
	return CursorPositionScreenToWindow(window, x0, y0)
}

func CursorPositionScreenToWindow(window Window, x0, y0 float32) (x, y float32) {
	return cursorPositionScreenToWindow(window, x0, y0)
}

func CursorPositionWindowToScreen(window Window, x0, y0 float32) (x, y float32) {
	return cursorPositionWindowToScreen(window, x0, y0)
}
