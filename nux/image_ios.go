// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ios

package nux

/*
#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>

CGImageRef nux_loadImageFromFile(char* name){
	UIImage *loadImage = [UIImage imageNamed:[NSString stringWithUTF8String:name]];
	CGImageRef cgimage=loadImage.CGImage;
	return cgimage;
}
*/
import "C"
import (
	"path/filepath"
	"runtime"
	"unsafe"
)

func loadImageFromFile(path string) Image {
	path, _ = filepath.Abs(path)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	me := &nativeImage{
		ptr: C.nux_loadImageFromFile(cpath),
	}
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	C.CGImageRelease(img.ptr)
}

type nativeImage struct {
	ptr C.CGImageRef
}

func (me *nativeImage) Size() (width, height int32) {
	return int32(C.CGImageGetWidth(me.ptr)), int32(C.CGImageGetHeight(me.ptr))
}
