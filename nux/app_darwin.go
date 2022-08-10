// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
)

type nativeApp struct {
	ptr darwin.NSApplication
}

func createNativeApp_() *nativeApp {
	me := &nativeApp{}
	darwin.SetNSApplicationDelegate(me)
	me.ptr = darwin.NSApp_SharedApplication()
	return me
}

func (me *nativeApp) init() {

}

func (me *nativeApp) run() {
	theApp.mainWindow = newWindow(theApp.manifest.GetAttr("mainWindow", nil))
	theApp.mainWindow.Center()
	theApp.mainWindow.Show()

	theApp.windowPrepared <- struct{}{}

	me.ptr.Run()
}

func (me *nativeApp) terminate() {
	me.ptr.Terminate()
}

func (me *nativeApp) ApplicationWillFinishLaunching(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillFinishLaunching")
}

func (me *nativeApp) ApplicationDidFinishLaunching(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationDidFinishLaunching")
	theApp.mainWindow.mountWidget()
}

func (me *nativeApp) ApplicationWillBecomeActive(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillBecomeActive")
}

func (me *nativeApp) ApplicationDidBecomeActive(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationDidBecomeActive")
}

func (me *nativeApp) ApplicationWillResignActive(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillResignActive")
}

func (me *nativeApp) ApplicationDidResignActive(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationDidResignActive")
}

func (me *nativeApp) ApplicationShouldTerminate(sender darwin.NSApplication) darwin.NSApplicationTerminateReply {
	// log.I("nuxui", "ApplicationShouldTerminate")
	return darwin.NSTerminateNow
}

func (me *nativeApp) ApplicationShouldTerminateAfterLastWindowClosed(sender darwin.NSApplication) bool {
	// log.I("nuxui", "ApplicationShouldTerminateAfterLastWindowClosed")
	return true
}

func (me *nativeApp) ApplicationWillTerminate(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillTerminate")
}

func (me *nativeApp) ApplicationWillHide(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillHide")
}

func (me *nativeApp) ApplicationDidHide(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationDidHide")
}

func (me *nativeApp) ApplicationWillUnhide(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationWillUnhide")
}

func (me *nativeApp) ApplicationDidUnhide(notification darwin.NSNotification) {
	// log.I("nuxui", "ApplicationDidUnhide")
}

func runOnUI(callback func()) {
	if IsMainThread() {
		callback()
	} else {
		darwin.BackToUI(callback)
	}
}

func invalidateRectAsync_(dirtRect *Rect) {
	// TODO:: error for render radio options
	// darwin.NSApp().KeyWindow().InvalidateRectAsync(float32(rect.X), float32(rect.Y), float32(rect.Width), float32(rect.Height))
	darwin.NSApp().KeyWindow().InvalidateRectAsync(0, 0, 0, 0)
}

func startTextInput() {
	darwin.NSApp().KeyWindow().StartTextInputAsync()
}

func stopTextInput() {
	darwin.NSApp().KeyWindow().StopTextInputAsync()
}

func currentThreadID_() uint64 {
	return darwin.CurrentThreadID()
}

func screenSize() (width, height int32) {
	return darwin.NSScreen_frameSize()
}
