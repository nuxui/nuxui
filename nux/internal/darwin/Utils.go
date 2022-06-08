// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>
#include <pthread.h>

void nux_CursorScreenPosition(CGFloat *outX, CGFloat *outY) {
  if (outX) { *outX = [NSEvent mouseLocation].x; };
  if (outY) { *outY = [NSScreen mainScreen].frame.size.height - [NSEvent mouseLocation].y; };
}

void nux_CursorPositionScreenToWindow(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY) {
  NSWindow *w = (NSWindow *)window;
  if (outX) { *outX = w.frame.origin.x + x; };
  if (outY) { *outY = [NSScreen mainScreen].frame.size.height - ([w contentView].bounds.size.height - y + w.frame.origin.y); };
}

void nux_CursorPositionWindowToScreen(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY) {
  NSWindow *w = (NSWindow *)window;
  if (outX) { *outX = x - w.frame.origin.x; };
  if (outY) { *outY = [w contentView].bounds.size.height - (([NSScreen mainScreen].frame.size.height - y) - w.frame.origin.y); };
}

uint64 nux_CurrentThreadID() {
  uint64 id;
  if (pthread_threadid_np(pthread_self(), &id)) { abort(); };
  return id;
}

void nux_NSObject_release(uintptr_t ptr){
	[((NSObject*)ptr) release];
}
*/
import "C"

func CurrentThreadID() uint64 {
	return uint64(C.nux_CurrentThreadID())
}

func CursorScreenPosition() (x, y float32) {
	var outX, outY C.CGFloat
	C.nux_CursorScreenPosition(&outX, &outY)
	return float32(outX), float32(outY)
}

func CursorPositionScreenToWindow(window NSWindow, px, py float32) (x, y float32) {
	var outX, outY C.CGFloat
	C.nux_CursorPositionScreenToWindow(C.uintptr_t(window), C.CGFloat(px), C.CGFloat(py), &outX, &outY)
	return float32(outX), float32(outY)
}

func CursorPositionWindowToScreen(window NSWindow, px, py float32) (x, y float32) {
	var outX, outY C.CGFloat
	C.nux_CursorPositionWindowToScreen(C.uintptr_t(window), C.CGFloat(px), C.CGFloat(py), &outX, &outY)
	return float32(outX), float32(outY)
}

func NSObject_release(ptr uintptr) {
	C.nux_NSObject_release(C.uintptr_t(ptr))
}
