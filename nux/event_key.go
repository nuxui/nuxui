// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type KeyCode int

type KeyEvent interface {
	Event
	KeyCode() KeyCode
}

type keyEvent struct {
	event
	keyCode KeyCode
}

func (me *keyEvent) KeyCode() KeyCode {
	return me.keyCode
}
