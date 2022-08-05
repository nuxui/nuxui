// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
)

type nativeApp struct {
	ref android.Application
}

func createNativeApp_() *nativeApp {
	me := &nativeApp{
		ref: android.NuxApplication_instance(),
	}
	android.SetApplicationDelegate(me)
	runtime.SetFinalizer(me, freeApp)
	return me
}

func freeApp(app *nativeApp) {
}

func (me *nativeApp) run() {
	log.I("nuxui", "nativeApp run() 0 %d", currentThreadID())
	// theApp.windowPrepared <- struct{}{}
}

func (me *nativeApp) terminate() {
}

func (me *nativeApp) OnCreateWindow(activity android.Activity) {
	// theApp.createMainWindow()
}

func (me *nativeApp) OnConfigurationChanged(app android.Application, newConfig android.Configuration) {

}

func (me *nativeApp) OnCreate(app android.Application) {
	theApp.createMainWindow()
}

func (me *nativeApp) OnLowMemory(app android.Application) {

}

func (me *nativeApp) OnTerminate(app android.Application) {

}

func (me *nativeApp) OnTrimMemory(app android.Application, level int32) {

}

func runOnUI(callback func()) {
	if IsMainThread() {
		callback()
	} else {
		android.BackToUI(callback)
	}
}

func invalidateRectAsync_(dirtRect *Rect) {
	// TODO:: error for render radio options
	// darwin.NSApp().KeyWindow().InvalidateRectAsync(float32(rect.X), float32(rect.Y), float32(rect.Width), float32(rect.Height))
	// darwin.NSApp().KeyWindow().InvalidateRectAsync(0, 0, 0, 0)
}

func startTextInput() {
}

func stopTextInput() {
}

func currentThreadID_() uint64 {
	return android.GetTid()
}

func screenSize() (width, height int32) {
	// return darwin.NSScreen_frameSize()
	return 100, 100
}
