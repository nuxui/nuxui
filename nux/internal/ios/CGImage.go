// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION -DGLES_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics -framework UIKit -framework CoreText -framework GLKit -framework UniformTypeIdentifiers -framework QuartzCore -framework ImageIO

#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

CGImageRef nux_CGImageSourceCreateImageAtIndex(char* name){
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
	"unsafe"
)

func CGImageSourceCreateImageAtIndex(filename string) CGImageRef {
	cpath := C.CString(filename)
	defer C.free(unsafe.Pointer(cpath))
	return CGImageRef(C.nux_CGImageSourceCreateImageAtIndex(cpath))
}

func CGImageCreateCopy(image CGImageRef) CGImageRef {
	return CGImageRef(C.CGImageCreateCopy(C.CGImageRef(image)))
}

func CGImageRelease(image CGImageRef) {
	C.CGImageRelease(C.CGImageRef(image))
}

func CGImageGetSize(image CGImageRef) (width, height int32) {
	return int32(C.CGImageGetWidth(C.CGImageRef(image))), int32(C.CGImageGetHeight(C.CGImageRef(image)))
}


