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
	runOnUI                  = make(chan func())
	theNSApplicationDelegate NSApplicationDelegate
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

func SetNSApplicationDelegate(delegate NSApplicationDelegate) {
	theNSApplicationDelegate = delegate
}

//------------------------------------------------------------

//export go_nux_app_delegate
func go_nux_app_delegate(action C.int, obj C.uintptr_t) C.int {
	if theNSApplicationDelegate == nil {
		return -1
	}

	switch action {
	case 1:
		theNSApplicationDelegate.ApplicationWillFinishLaunching(NSNotification(obj))
	case 2:
		theNSApplicationDelegate.ApplicationDidFinishLaunching(NSNotification(obj))
	case 3:
		theNSApplicationDelegate.ApplicationWillBecomeActive(NSNotification(obj))
	case 4:
		theNSApplicationDelegate.ApplicationDidBecomeActive(NSNotification(obj))
	case 5:
		theNSApplicationDelegate.ApplicationWillResignActive(NSNotification(obj))
	case 6:
		theNSApplicationDelegate.ApplicationDidResignActive(NSNotification(obj))
	case 7:
		return C.int(theNSApplicationDelegate.ApplicationShouldTerminate(NSApplication(obj)))
	case 8:
		if theNSApplicationDelegate.ApplicationShouldTerminateAfterLastWindowClosed(NSApplication(obj)) {
			return 1
		}
		return 0
	case 9:
		theNSApplicationDelegate.ApplicationWillTerminate(NSNotification(obj))
	case 10:
		theNSApplicationDelegate.ApplicationWillHide(NSNotification(obj))
	case 11:
		theNSApplicationDelegate.ApplicationDidHide(NSNotification(obj))
	case 12:
		theNSApplicationDelegate.ApplicationWillUnhide(NSNotification(obj))
	case 13:
		theNSApplicationDelegate.ApplicationDidUnhide(NSNotification(obj))
	}
	return 0
}

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
