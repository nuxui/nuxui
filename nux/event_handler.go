// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type KeyEventHandler interface {
	OnKeyEvent(event KeyEvent) bool
}

type TypingEventHandler interface {
	OnTypingEvent(event TypingEvent) bool
}
