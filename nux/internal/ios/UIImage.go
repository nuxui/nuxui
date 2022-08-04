// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics -framework UIKit -framework CoreText -framework GLKit

#import <QuartzCore/QuartzCore.h>
#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

uintptr_t nux_UIImage_imageNamed(char* name){
	UIImage *image = [UIImage imageNamed:[NSString stringWithUTF8String:name]];
	return (uintptr_t)image;
}

void nux_UIImage_size(uintptr_t uiimage, CGFloat* outW, CGFloat* outH){
	CGSize s = [(UIImage*)uiimage size];
	if (outW) { *outW = s.width; };
	if (outH) { *outH = s.height; };
}

CGImageRef nux_UIImage_CGImage(uintptr_t uiimage){
	return [(UIImage*)uiimage CGImage];
}
*/
import "C"
import (
	"unsafe"
)

func UIImage_ImageNamed(filename string) UIImage {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return UIImage(C.nux_UIImage_imageNamed(cstr))
}

func (me UIImage) Size() (width, height float32) {
	var outW, outH C.CGFloat
	C.nux_UIImage_size(C.uintptr_t(me), &outW, &outH)
	return float32(outW), float32(outH)
}

func (me UIImage) CGImage() CGImageRef {
	return CGImageRef(C.nux_UIImage_CGImage(C.uintptr_t(me)))
}
