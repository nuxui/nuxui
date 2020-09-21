// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type GestureRecognizer interface {
	GestureArenaMember
	PointerAllowed(event Event) bool
	HandlePointerEvent(event Event)
	//TODO:: Clear(widget Widget) // clear callbacks of widget
}
