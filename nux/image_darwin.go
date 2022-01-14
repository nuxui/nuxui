// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

CGImageRef createImage(char* name){
    NSURL *url = [NSURL fileURLWithPath: [NSString stringWithUTF8String:name]];
    NSError *err;
    if ([url checkResourceIsReachableAndReturnError:&err] == NO){
        NSLog(@"%@", err);
    }
    CGImageSourceRef source = CGImageSourceCreateWithURL((__bridge CFURLRef)url, nil);
    CGImageRef ref = CGImageSourceCreateImageAtIndex(source, 0, nil);
	return ref;
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func CreateImage(path string) Image {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	me := &nativeImage{
		ptr: C.createImage(cpath),
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
