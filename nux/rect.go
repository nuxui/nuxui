// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"

	"github.com/nuxui/nuxui/util"
)

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                         Rect                            ///////////////////
/////////////////////////////////////////////////////////////////////////////////////////
type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

// type RectXY struct {
// 	Left   int32
// 	Top    int32
// 	Right  int32
// 	Bottom int32
// 	X      int32 // x in window
// 	Y      int32 // y in window
// }

func (me *Rect) Width() int32 {
	if me.Right < me.Left {
		return 0.0
	}
	return me.Right - me.Left
}

func (me *Rect) Height() int32 {
	if me.Bottom < me.Top {
		return 0
	}
	return me.Bottom - me.Top
}

func (me *Rect) Center() (x, y int32) {
	if me.Right < me.Left || me.Bottom < me.Top {
		return 0, 0
	}

	x = me.Left + util.Roundi32(float32(me.Right-me.Left)/2.0)
	y = me.Top + util.Roundi32(float32(me.Bottom-me.Top)/2.0)
	return x, y
}

func (me *Rect) String() string {
	return fmt.Sprintf("{Left: %d, Top: %d, Right: %d, Bottom: %d}", me.Left, me.Top, me.Right, me.Bottom)
}

/////////////////////////////////////////////////////////////////////////////////////////
/////////////                         RectF                            //////////////////
/////////////////////////////////////////////////////////////////////////////////////////

type RectF struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

func (me *RectF) Width() float32 {
	if me.Right < me.Left {
		return 0.0
	}
	return me.Right - me.Left
}

func (me *RectF) Height() float32 {
	if me.Bottom < me.Top {
		return 0
	}
	return me.Bottom - me.Top
}

func (me *RectF) Center() (x, y float32) {
	if me.Right < me.Left || me.Bottom < me.Top {
		return 0, 0
	}

	x = me.Left + (me.Right-me.Left)/2.0
	y = me.Top + (me.Bottom-me.Top)/2.0
	return x, y
}
