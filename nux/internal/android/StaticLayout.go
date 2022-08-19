// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_new_StaticLayout(char* text, jint width, jobject paint);
void nux_staticLayout_getSize(jobject staticLayout, jint* outWidth, jint* outHeight);
void nux_staticLayout_draw(jobject staticLayout, jobject canvas);
jint nux_staticLayout_getLineCount(jobject staticLayout);
*/
import "C"

import (
	"unsafe"
)

func NewStaticLayout(text string, width int32, paint Paint) StaticLayout {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	return StaticLayout(C.nux_new_StaticLayout(cstr, C.jint(width), C.jobject(paint)))
}

func (me StaticLayout) GetSize() (width, height int32) {
	w := C.jint(0)
	h := C.jint(0)
	C.nux_staticLayout_getSize(C.jobject(me), (*C.jint)(&w), (*C.jint)(&h))
	return int32(w), int32(h)
}

func (me StaticLayout) Draw(canvas Canvas) {
	C.nux_staticLayout_draw(C.jobject(me), C.jobject(canvas))
}
