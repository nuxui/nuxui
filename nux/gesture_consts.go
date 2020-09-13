// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type GestureState int

const (
	GestureState_Ready GestureState = 1 + iota
	GestureState_Possible
	GestureState_Accepted
)

const (
	DOWN_DELAY               = 100
	DOWN_TO_UP_DELAY         = 100
	MIN_PAN_DISTANCE float64 = 10 /*use dp*/
)
