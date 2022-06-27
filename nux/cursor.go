// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type NativeCursor int

const (
	CursorArrow NativeCursor = iota
	CursorIBeam
	CursorWait
	CursorCrosshair
	CursorResizeWE
	CursorResizeNS
	CursorResizeNWSE
	CursorResizeNESW
	CursorHand
	CursorSizeAll
	CursorCustom
)

type Cursor interface {
	Set()
}

func LoadNativeCursor(c NativeCursor) Cursor {
	return loadNativeCursor(c)
}

// Retrieves the position of the mouse cursor, in screen coordinates.
func GetCursorScreenPosition() (x, y float32) {
	return getCursorScreenPosition()
}

// Retrieves the position of the mouse cursor, in window coordinates.
func GetCursorWindowPosition(window Window) (x, y float32) {
	px, py := getCursorScreenPosition()
	return CursorPositionScreenToWindow(window, px, py)
}

func CursorPositionScreenToWindow(window Window, px, py float32) (x, y float32) {
	return cursorPositionScreenToWindow(window, px, py)
}

func CursorPositionWindowToScreen(window Window, px, py float32) (x, y float32) {
	return cursorPositionWindowToScreen(window, px, py)
}
