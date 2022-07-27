// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/linux/xlib"
)

func getCursorScreenPosition() (x, y float32) {
	// TODO::
	return
}

func cursorPositionScreenToWindow(wind Window, x0, y0 float32) (x, y float32) {
	// TODO::
	return
}

func cursorPositionWindowToScreen(wind Window, x0, y0 float32) (x, y float32) {
	// TODO::
	return
}

type cursor struct {
	c xlib.Cursor
}

func (me *cursor) Set() {
	w := theApp.MainWindow().native()
	xlib.XDefineCursor(w.display, w.window, me.c)
}

// https://tronche.com/gui/x/xlib/appendix/b/
func loadNativeCursor(c NativeCursor) *cursor {
	var shape xlib.CursorShape
	switch c {
	case CursorArrow:
		shape = xlib.XC_left_ptr
	case CursorIBeam:
		shape = xlib.XC_xterm
	case CursorWait:
		shape = xlib.XC_watch
	case CursorCrosshair:
		shape = xlib.XC_cross
	case CursorResizeWE:
		shape = xlib.XC_sb_h_double_arrow
	case CursorResizeNS:
		shape = xlib.XC_sb_v_double_arrow
	case CursorResizeNWSE:
		shape = xlib.XC_fleur
	case CursorResizeNESW:
		shape = xlib.XC_fleur
	case CursorFinger:
		shape = xlib.XC_hand2
	case CursorDrag:
		shape = 61 // TODO:: which is drag?
	case CursorHand:
		shape = xlib.XC_hand1
	default:
		log.Fatal("nux", "unknown cursor type: %d", c)
	}
	return &cursor{c: xlib.XCreateFontCursor(theApp.native.display, shape)}
}
