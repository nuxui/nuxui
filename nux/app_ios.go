// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/ios"
)

type nativeApp struct {
	ptr ios.UIApplication
}

func createNativeApp_() *nativeApp {
	me := &nativeApp{}
	ios.SetUIApplicationDelegate(me)
	me.ptr = ios.UIApplication_SharedApplication()
	return me
}

func (me *nativeApp) run() {
	me.ptr.Run()
}

func (me *nativeApp) terminate() {
}

func (me *nativeApp) WillFinishLaunchingWithOptions(launchOptions map[string]uintptr) bool {
	// log.I("nuxui", "WillFinishLaunchingWithOptions")
	return true
}

func (me *nativeApp) DidFinishLaunchingWithOptions(launchOptions map[string]uintptr) bool {
	// log.I("nuxui", "DidFinishLaunchingWithOptions")
	// theApp.onDidFinishLaunch()

	theApp.createMainWindow()
	theApp.mainWindow.Center()
	theApp.mainWindow.Show()

	theApp.windowPrepared <- struct{}{}

	return true
}

func runOnUI(callback func()) {
	if IsMainThread() {
		callback()
	} else {
		ios.BackToUI(callback)
	}
}

func invalidateRectAsync_(dirtRect *Rect) {
	theApp.mainWindow.native().ptr.InvalidateRectAsync(0, 0, 0, 0)
}

func startTextInput() {
	// darwin.NSApp().KeyWindow().StartTextInputAsync()
}

func stopTextInput() {
	// darwin.NSApp().KeyWindow().StopTextInputAsync()
}

func currentThreadID_() uint64 {
	return ios.CurrentThreadID()
}

func screenSize() (width, height int32) {
	return ios.UIScreen_MainScreenSize()
}
