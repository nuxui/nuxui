// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

/*
#import <QuartzCore/QuartzCore.h>
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>

#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa

void      nux_BackToUI();

uintptr_t nux_NSApp();
uintptr_t nux_NSApp_SharedApplication();
uintptr_t nux_NSApp_KeyWindow(uintptr_t app);
void      nux_NSApp_Run(uintptr_t app);
void      nux_NSApp_Terminate(uintptr_t app);
*/
import "C"

var (
	runOnUI = make(chan func())
)

func NSApp_SharedApplication() NSApplication {
	return NSApplication(C.nux_NSApp_SharedApplication())
}

func NSApp() NSApplication {
	return NSApplication(C.nux_NSApp())
}

func (me NSApplication) KeyWindow() NSWindow {
	return NSWindow(C.nux_NSApp_KeyWindow(C.uintptr_t(me)))
}

func (me NSApplication) Run() {
	C.nux_NSApp_Run(C.uintptr_t(me))
}

func (me NSApplication) Terminate() {
	C.nux_NSApp_Terminate(C.uintptr_t(me))
}

//------------------------------------------------------------

//export go_nux_app_sendEvent
func go_nux_app_sendEvent(event C.uintptr_t) int {
	// TODO:: app event
	return 0
}

func BackToUI(callback func()) {
	go func() {
		runOnUI <- callback
	}()
	C.nux_BackToUI()
}

//export go_nux_backToUI
func go_nux_backToUI() {
	callback := <-runOnUI
	callback()
}
