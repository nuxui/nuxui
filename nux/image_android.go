// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

/*
#include <jni.h>
jobject createImage(char* fileName);
jint bitmap_getWidth(jobject bitmap);
jint bitmap_getHeight(jobject bitmap);
void bitmap_recycle(jobject bitmap);
*/
import "C"
import (
	"runtime"
	// "unsafe"
)

func CreateImage(path string) Image {
	chars := ([]byte)(path)

	me := &nativeImage{
		ptr: C.createImage((*C.char)(&chars[0])),
	}
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	C.bitmap_recycle(img.ptr)
}

type nativeImage struct {
	ptr C.jobject
}

func (me *nativeImage) Size() (width, height int32) {
	return int32(C.bitmap_getWidth(me.ptr)), int32(C.bitmap_getHeight(me.ptr))
}
