// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-geticoninfoexw

func getCursorScreenPosition() (x, y float32) {
	var p win32.POINT
	err := win32.GetCursorPos(&p)
	if err != nil {
		log.E("nux", "windows GetCursorPos faild: %s", err.Error())
		return
	}
	return float32(p.X), float32(p.Y)
}

func cursorPositionScreenToWindow(wind Window, x0, y0 float32) (x, y float32) {
	p := &win32.POINT{X: int32(x0), Y: int32(y0)}
	err := win32.ScreenToClient(wind.native().hwnd, p)
	if err != nil {
		log.E("nux", "windows ScreenToClient faild: %s", err.Error())
		return
	}
	return float32(p.X), float32(p.Y)
}

func cursorPositionWindowToScreen(wind Window, x0, y0 float32) (x, y float32) {
	p := &win32.POINT{X: int32(x0), Y: int32(y0)}
	err := win32.ClientToScreen(wind.native().hwnd, p)
	if err != nil {
		log.E("nux", "windows ClientToScreen faild: %s", err.Error())
		return
	}
	return float32(p.X), float32(p.Y)
}

type cursor struct {
	ptr uintptr
}

func (me *cursor) Set() {
	theApp.native.cursor = me.ptr
}

func loadNativeCursor(c NativeCursor) *cursor {
	var shape uintptr
	switch c {
	case CursorArrow:
		shape = win32.IDC_ARROW
	case CursorIBeam:
		shape = win32.IDC_IBEAM
	case CursorWait:
		shape = win32.IDC_WAIT
	case CursorCrosshair:
		shape = win32.IDC_CROSS
	case CursorResizeWE:
		shape = win32.IDC_SIZEWE
	case CursorResizeNS:
		shape = win32.IDC_SIZENS
	case CursorResizeNWSE:
		shape = win32.IDC_SIZENWSE
	case CursorResizeNESW:
		shape = win32.IDC_SIZENESW
	default:
		log.Fatal("nux", "unknown cursor type: %d", c)
	}

	c, err := win32.LoadCursor(0, shape)
	if err != nil {
		log.Fatal("nuxui", "error when LoadCursor: %s", err.Error())
	}
	return &cursor{ptr: c}
}
