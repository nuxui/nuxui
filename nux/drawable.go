// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type DrawableCreator func(Widget, ...Attr) Drawable

type Drawable interface {
	Size() (width, height int32)
	SetBounds(x, y, width, height int32)
	Draw(canvas Canvas)
	Equal(Drawable) bool
}
