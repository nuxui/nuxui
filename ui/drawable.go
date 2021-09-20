// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

type Drawable interface {
	//TODO  Draw(Rect, Canvas)
	Width() int32
	Height() int32
	Draw(nux.Canvas)
	Equal(Drawable) bool
}
