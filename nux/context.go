// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Context interface {
	Window() Window
	OpenIM() bool
	CloseIM() bool
}

type context struct {
	window Window
}

func (me *context) Window() Window {
	return me.window
}

func (me *context) OpenIM() bool {
	return true
}

func (me *context) CloseIM() bool {
	return true
}
