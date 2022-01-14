// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS:
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>

void cursor_getScreenPosition(CGFloat* x, CGFloat* y);
void cursor_positionWindowToScreen(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY);
void cursor_positionScreenToWindow(uintptr_t window, CGFloat x, CGFloat y, CGFloat *outX, CGFloat *outY);
*/
import "C"

func getCursorScreenPosition() (x, y float32) {
	var outX, outY C.CGFloat
	C.cursor_getScreenPosition(&outX, &outY)
	return float32(outX), float32(outY)
}

func cursorPositionScreenToWindow(wind Window, x0, y0 float32) (x, y float32) {
	var outX, outY C.CGFloat
	C.cursor_positionScreenToWindow(wind.(*window).windptr, C.CGFloat(x0), C.CGFloat(y0), &outX, &outY)
	return float32(outX), float32(outY)
}

func cursorPositionWindowToScreen(wind Window, x0, y0 float32) (x, y float32) {
	var outX, outY C.CGFloat
	C.cursor_positionWindowToScreen(wind.(*window).windptr, C.CGFloat(x0), C.CGFloat(y0), &outX, &outY)
	return float32(outX), float32(outY)
}
