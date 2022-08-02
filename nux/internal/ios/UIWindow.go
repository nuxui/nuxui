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

uintptr_t nux_NewUIWindow(CGFloat width, CGFloat height);
void nux_UIWindow_makeKeyAndVisible(uintptr_t nuxwindow);
CGRect nux_UIWindow_frame(uintptr_t nuxwindow);
void nux_UIWindow_InvalidateRect_async(uintptr_t nuxwindow, CGFloat x, CGFloat y, CGFloat width, CGFloat height);
*/
import "C"

const (
	Event_WindowDidResize = iota
	Event_WindowDrawRect
)

var windowEventHandler func(any) bool

func NewUIWindow(width, height int32) UIWindow {
	return UIWindow(C.nux_NewUIWindow(C.CGFloat(width), C.CGFloat(height)))
}

func (me UIWindow) MakeKeyAndVisible() {
	C.nux_UIWindow_makeKeyAndVisible(C.uintptr_t(me))
}

func (me UIWindow) Frame() CGRect {
	return CGRect(C.nux_UIWindow_frame(C.uintptr_t(me)))
}

func (me UIWindow) InvalidateRectAsync(x, y, width, height float32) {
	C.nux_UIWindow_InvalidateRect_async(C.uintptr_t(me), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
}

func SetWindowEventHandler(handler func(any) bool) {
	windowEventHandler = handler
}

//export go_nux_window_sendEvent
func go_nux_window_sendEvent(event C.uintptr_t) int {
	if windowEventHandler != nil && windowEventHandler(UIEvent(event)) {
		return 1
	}
	return 0
}

//export go_nux_windowDidLoad
func go_nux_windowDidLoad(window C.uintptr_t) {
	if windowEventHandler != nil {
		windowEventHandler(&WindowEvent{
			Window: UIWindow(window),
			Type:   Event_WindowDidResize,
		})
	}
}

//export go_nux_windowDrawRect
func go_nux_windowDrawRect(window C.uintptr_t) {
	if windowEventHandler != nil {
		windowEventHandler(&WindowEvent{
			Window: UIWindow(window),
			Type:   Event_WindowDrawRect,
		})
	}
}
