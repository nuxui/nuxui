// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>

jint   nux_motionevent_getPointerCount(jobject motionevent);
jfloat nux_motionevent_getX(jobject motionevent, jint index);
jfloat nux_motionevent_getY(jobject motionevent, jint index);
void   nux_motionevent_getXY(jobject motionevent, jint index, jfloat* outX, jfloat* outY);
jint   nux_motionevent_getActionMasked(jobject motionevent);
*/
import "C"

const (
	ACTION_DOWN = iota
	ACTION_UP
	ACTION_MOVE
	ACTION_CANCEL
	ACTION_OUTSIDE
	ACTION_POINTER_DOWN
	ACTION_POINTER_UP
	ACTION_HOVER_MOVE
	ACTION_SCROLL
	ACTION_HOVER_ENTER
	ACTION_HOVER_EXIT
	ACTION_BUTTON_PRESS
	ACTION_BUTTON_RELEASE // 12
)

func (me MotionEvent) GetPointerCount() int32 {
	return int32(C.nux_motionevent_getPointerCount(C.jobject(me)))
}

func (me MotionEvent) GetX(index int) float32 {
	return float32(C.nux_motionevent_getX(C.jobject(me), C.jint(index)))
}

func (me MotionEvent) GetY(index int) float32 {
	return float32(C.nux_motionevent_getY(C.jobject(me), C.jint(index)))
}

func (me MotionEvent) GetXY(index int) (x, y float32) {
	outX := C.jfloat(0)
	outY := C.jfloat(0)
	C.nux_motionevent_getXY(C.jobject(me), C.jint(index), &outX, &outY)
	return float32(outX), float32(outY)
}

func (me MotionEvent) GetActionMasked() int32 {
	return int32(C.nux_motionevent_getActionMasked(C.jobject(me)))
}
