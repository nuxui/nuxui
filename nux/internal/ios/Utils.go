// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#include <stdint.h>
#include <pthread.h>
#include <UIKit/UIDevice.h>
#import <GLKit/GLKit.h>
#import <UIKit/UIKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

uint64_t nux_CurrentThreadID() {
	uint64_t id;
	if (pthread_threadid_np(pthread_self(), &id)) {
		abort();
	}
	return id;
}

void nux_UIScreen_mainScreenSize(CGFloat *width, CGFloat *height){
  CGSize size = [UIScreen mainScreen].bounds.size;
  if (width) { *width = size.width; };
  if (height) { *width = size.height; };
}

void nux_NSObject_release(uintptr_t ptr){
	[((NSObject*)ptr) release];
}
*/
import "C"

func UIScreen_MainScreenSize() (width, height int32) {
	var w, h C.CGFloat
	C.nux_UIScreen_mainScreenSize(&w, &h)
	return int32(w), int32(h)
}

func NSObject_release(ptr uintptr) {
	C.nux_NSObject_release(C.uintptr_t(ptr))
}

func CurrentThreadID() uint64 {
	return uint64(C.nux_CurrentThreadID())
}
