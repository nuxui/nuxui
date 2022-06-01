// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
)

func getCursorScreenPosition() (x, y float32) {
	return darwin.CursorScreenPosition()
}

func cursorPositionScreenToWindow(w Window, px, py float32) (x, y float32) {
	return darwin.CursorPositionScreenToWindow(w.nativeWindow().ptr, px, py)
}

func cursorPositionWindowToScreen(w Window, px, py float32) (x, y float32) {
	return darwin.CursorPositionWindowToScreen(w.nativeWindow().ptr, px, py)
}
