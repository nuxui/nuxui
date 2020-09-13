// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type WindowEvent interface {
	Event
	Window() Window
}

type windowEvent struct {
	event
	window Window // TODO:: use window id
}

func (me *windowEvent) Window() Window {
	return me.window
}
