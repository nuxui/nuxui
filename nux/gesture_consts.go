// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "time"

type GestureCallback func(Widget)

type GestureState int

const (
	GestureState_Ready GestureState = 1 + iota
	GestureState_Possible
	GestureState_Accepted
	GestureState_Rejected
)

const (
	GESTURE_DOWN_DELAY                 = 100 * time.Millisecond
	GESTURE_DOWN2UP_DELAY              = 100
	GESTURE_LONG_PRESS_TIMEOUT         = 500 * time.Millisecond
	GESTURE_DOUBLETAP_TIMEOUT          = 300 * time.Millisecond
	GESTURE_MIN_PAN_DISTANCE   float32 = 10 /*use dp*/
)
