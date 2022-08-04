// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_createBitmap(char* fileName);
jint nux_bitmap_getWidth(jobject bitmap);
jint nux_bitmap_getHeight(jobject bitmap);
void nux_bitmap_recycle(jobject bitmap);
*/
import "C"

import "unsafe"

func CreateBitmap(filename string) Bitmap {
	cpath := C.CString(filename)
	defer C.free(unsafe.Pointer(cpath))
	return Bitmap(C.nux_createBitmap(cpath))
}

func (me Bitmap) GetWidth() int32 {
	return int32(C.nux_bitmap_getWidth(C.jobject(me)))
}

func (me Bitmap) GetHeight() int32 {
	return int32(C.nux_bitmap_getHeight(C.jobject(me)))
}

func (me Bitmap) Recycle() {
	C.nux_bitmap_recycle(C.jobject(me))
}
