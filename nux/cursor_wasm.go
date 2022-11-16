// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

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
}

func (me *cursor) Set() {
}

func loadNativeCursor(c NativeCursor) *cursor {
	return &cursor{}
}
