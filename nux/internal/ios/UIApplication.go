// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package ios

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION -DGLES_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework CoreGraphics -framework UIKit -framework CoreText -framework GLKit -framework UniformTypeIdentifiers -framework QuartzCore

#include <stdint.h>
#include <pthread.h>
#include <UIKit/UIDevice.h>
#import <GLKit/GLKit.h>
#import <UIKit/UIKit.h>
#import <CoreText/CoreText.h>
#import <CoreGraphics/CoreGraphics.h>

void      nux_BackToUI();

void nux_UIApplication_Run();
uintptr_t nux_UIApplication_sharedApplication();
*/
import "C"

var (
	runOnUI                  = make(chan func())
	theUIApplicationDelegate UIApplicationDelegate
)

func UIApplication_SharedApplication() UIApplication {
	return UIApplication(C.nux_UIApplication_sharedApplication())
}

func (me UIApplication) Run() {
	C.nux_UIApplication_Run()
}

func SetUIApplicationDelegate(delegate UIApplicationDelegate) {
	theUIApplicationDelegate = delegate
}

//export go_nux_app_delegate
func go_nux_app_delegate(action C.int, obj C.uintptr_t) C.int {
	if theUIApplicationDelegate == nil {
		return -1 // not handled
	}

	switch action {
	case 1:
		theUIApplicationDelegate.WillFinishLaunchingWithOptions(nil)
	case 2:
		theUIApplicationDelegate.DidFinishLaunchingWithOptions(nil)
	}

	return -1 // not handled
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
