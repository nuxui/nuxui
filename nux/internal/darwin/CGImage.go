// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

CGImageRef CGImageSourceCreateImageAtIndex_(char* name){
    NSURL *url = [NSURL fileURLWithPath: [NSString stringWithUTF8String:name]];
    NSError *err;
    if ([url checkResourceIsReachableAndReturnError:&err] == NO){
        NSLog(@"%@", err);
    }
    CGImageSourceRef source = CGImageSourceCreateWithURL((__bridge CFURLRef)url, nil);
    CGImageRef ref = CGImageSourceCreateImageAtIndex(source, 0, nil);
	CFRelease(source);
	return ref;
}
*/
import "C"

import (
	"path/filepath"
	"unsafe"
)

func CGImageSourceCreateImageAtIndex(filename string) CGImage {
	path, _ := filepath.Abs(filename)
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return CGImage(C.CGImageSourceCreateImageAtIndex_(cpath))
}

func CGImageRelease(image CGImage) {
	C.CGImageRelease(C.CGImageRef(image))
}

func CGImageGetSize(image CGImage) (width, height int32) {
	return int32(C.CGImageGetWidth(C.CGImageRef(image))), int32(C.CGImageGetHeight(C.CGImageRef(image)))
}
