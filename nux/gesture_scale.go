// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

func OnScaleStart(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_SCALE_START, callback)
}

func OnScaleUpdate(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_SCALE_UPDATE, callback)
}

func OnScaleEnd(widget Widget, callback GestureCallback) {
	addPanCallback(widget, _ACTION_SCALE_END, callback)
}

const (
	_ACTION_SCALE_START = iota
	_ACTION_SCALE_UPDATE
	_ACTION_SCALE_END
)
