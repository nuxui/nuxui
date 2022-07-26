// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/darwin"
)

func getCursorScreenPosition() (x, y float32) {
	return darwin.CursorScreenPosition()
}

func cursorPositionScreenToWindow(w Window, px, py float32) (x, y float32) {
	return darwin.CursorPositionScreenToWindow(w.native().ptr, px, py)
}

func cursorPositionWindowToScreen(w Window, px, py float32) (x, y float32) {
	return darwin.CursorPositionWindowToScreen(w.native().ptr, px, py)
}

type cursor struct {
	ptr darwin.NSCursor
}

func (me *cursor) Set() {
	me.ptr.Set()
}

func loadNativeCursor(c NativeCursor) *cursor {
	var ptr darwin.NSCursor
	switch c {
	case CursorArrow:
		ptr = darwin.NSCursor_ArrowCursor()
	case CursorIBeam:
		ptr = darwin.NSCursor_IBeamCursor()
	case CursorHand:
		ptr = darwin.NSCursor_PointingHandCursor()
	default:
		log.Fatal("nux", "unknown cursor type: %d", c)
	}
	return &cursor{ptr: ptr}
}
