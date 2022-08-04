// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_new_Paint();
void nux_paint_setColor(jobject paint, uint32_t color);
jint nux_paint_getColor(jobject paint);
void nux_paint_setTextSize(jobject paint, jfloat textSize);
void nux_paint_setStyle(jobject paint, jint style);
void nux_paint_setAntiAlias(jobject paint, jboolean aa);
void nux_paint_measureText(jobject paint, char* text, jint width, jint *outWidth, jint* outHeight);
*/
import "C"

func NewPaint() Paint {
	return Paint(C.nux_new_Paint())
}

func (me Paint) SetColor(color uint32) {
	C.nux_paint_setColor(C.jobject(me), C.uint32_t(color))
}

func (me Paint) GetColor() uint32 {
	return uint32(C.nux_paint_getColor(C.jobject(me)))
}
