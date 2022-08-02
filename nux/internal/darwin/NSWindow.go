// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Cocoa/Cocoa.h>

uintptr_t nux_NewNSWindow(CGFloat width, CGFloat height);
char*     nux_NSWindow_Title(uintptr_t window);
void      nux_NSWindow_SetTitle(uintptr_t window, char* title);
float     nux_NSWindow_Alpha(uintptr_t window);
void      nux_NSWindow_SetAlpha(uintptr_t window, float alpha);
void      nux_NSWindow_Size(uintptr_t window, int32_t *width, int32_t *height);
void      nux_NSWindow_ContentSize(uintptr_t window, int32_t *width, int32_t *height);
void      nux_NSWindow_Center(uintptr_t window);
void      nux_NSWindow_MakeKeyAndOrderFront(uintptr_t window);
void      nux_NSWindow_SetContentView(uintptr_t window, uintptr_t view);
void      nux_NSWindow_StartTextInput_async(uintptr_t window);
void      nux_NSWindow_StopTextInput_async(uintptr_t window);
void      nux_NSWindow_InvalidateRect_async(uintptr_t window, CGFloat x, CGFloat y, CGFloat width, CGFloat height);
void      nux_NSWindow_SetTextInputRect_async(uintptr_t window, CGFloat x, CGFloat y, CGFloat width, CGFloat height);
*/
import "C"

import (
	"unsafe"
)

const (
	NSWindowStyleMaskBorderless             NSWindowStyleMask = C.NSWindowStyleMaskBorderless
	NSWindowStyleMaskTitled                 NSWindowStyleMask = C.NSWindowStyleMaskTitled
	NSWindowStyleMaskClosable               NSWindowStyleMask = C.NSWindowStyleMaskClosable
	NSWindowStyleMaskMiniaturizable         NSWindowStyleMask = C.NSWindowStyleMaskMiniaturizable
	NSWindowStyleMaskResizable              NSWindowStyleMask = C.NSWindowStyleMaskResizable
	NSWindowStyleMaskUtilityWindow          NSWindowStyleMask = C.NSWindowStyleMaskUtilityWindow
	NSWindowStyleMaskDocModalWindow         NSWindowStyleMask = C.NSWindowStyleMaskDocModalWindow
	NSWindowStyleMaskNonactivatingPanel     NSWindowStyleMask = C.NSWindowStyleMaskNonactivatingPanel
	NSWindowStyleMaskUnifiedTitleAndToolbar NSWindowStyleMask = C.NSWindowStyleMaskUnifiedTitleAndToolbar
	NSWindowStyleMaskFullScreen             NSWindowStyleMask = C.NSWindowStyleMaskFullScreen
	NSWindowStyleMaskFullSizeContentView    NSWindowStyleMask = C.NSWindowStyleMaskFullSizeContentView
	NSWindowStyleMaskHUDWindow              NSWindowStyleMask = C.NSWindowStyleMaskHUDWindow
)

const (
	Event_WindowDidResize = iota
	Event_WindowDrawRect
)

var windowEventHandler func(any) bool

func NewNSWindow(width, height int32) NSWindow {
	return NSWindow(C.nux_NewNSWindow(C.CGFloat(width), C.CGFloat(height)))
}

func (me NSWindow) Alpha() float32 {
	return float32(C.nux_NSWindow_Alpha(C.uintptr_t(me)))
}

func (me NSWindow) SetAlpha(alpha float32) {
	C.nux_NSWindow_SetAlpha(C.uintptr_t(me), C.float(alpha))
}

func (me NSWindow) Title() string {
	return C.GoString(C.nux_NSWindow_Title(C.uintptr_t(me)))
}

func (me NSWindow) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.nux_NSWindow_SetTitle(C.uintptr_t(me), ctitle)
}

func (me NSWindow) Size() (width, height int32) {
	var w, h C.int32_t
	C.nux_NSWindow_Size(C.uintptr_t(me), &w, &h)
	return int32(w), int32(h)
}

func (me NSWindow) ContentSize() (width, height int32) {
	var w, h C.int32_t
	C.nux_NSWindow_ContentSize(C.uintptr_t(me), &w, &h)
	return int32(w), int32(h)
}

func (me NSWindow) Center() {
	C.nux_NSWindow_Center(C.uintptr_t(me))
}

func (me NSWindow) MakeKeyAndOrderFront() {
	C.nux_NSWindow_MakeKeyAndOrderFront(C.uintptr_t(me))
}

func (me NSWindow) SetContentView(view NSView) {
	C.nux_NSWindow_SetContentView(C.uintptr_t(me), C.uintptr_t(view))
}

func (me NSWindow) StartTextInputAsync() {
	C.nux_NSWindow_StartTextInput_async(C.uintptr_t(me))
}

func (me NSWindow) StopTextInputAsync() {
	C.nux_NSWindow_StopTextInput_async(C.uintptr_t(me))
}

func (me NSWindow) InvalidateRectAsync(x, y, width, height float32) {
	C.nux_NSWindow_InvalidateRect_async(C.uintptr_t(me), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
}

func (me NSWindow) SetTextInputRectAsync(x, y, width, height float32) {
	C.nux_NSWindow_SetTextInputRect_async(C.uintptr_t(me), C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
}

func SetWindowEventHandler(handler func(any) bool) {
	windowEventHandler = handler
}

//export go_nux_window_sendEvent
func go_nux_window_sendEvent(event C.uintptr_t) int {
	if windowEventHandler != nil && windowEventHandler(NSEvent(event)) {
		return 1
	}
	return 0
}

//export go_nux_windowDidResize
func go_nux_windowDidResize(window C.uintptr_t) {
	if windowEventHandler != nil {
		windowEventHandler(&WindowEvent{
			Window: NSWindow(window),
			Type:   Event_WindowDidResize,
		})
	}
}

//export go_nux_windowDrawRect
func go_nux_windowDrawRect(window C.uintptr_t) {
	if windowEventHandler != nil {
		windowEventHandler(&WindowEvent{
			Window: NSWindow(window),
			Type:   Event_WindowDrawRect,
		})
	}
}

//export go_nux_windowTypingEvent
func go_nux_windowTypingEvent(window C.uintptr_t, chars *C.char, action, location, length C.int) {
	if windowEventHandler != nil {
		windowEventHandler(&TypingEvent{
			Window:   NSWindow(window),
			Text:     C.GoString(chars),
			Action:   int32(action),
			Location: int32(location),
			Length:   int32(length),
		})
	}
}
