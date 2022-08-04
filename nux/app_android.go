// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
)

type nativeApp struct {
	ref android.Application
}

func createNativeApp_() *nativeApp {
	me := &nativeApp{
		// ref: android.NuxApplication_instance(),
	}
	runtime.SetFinalizer(me, freeApp)
	return me
}

func freeApp(app *nativeApp) {
	// android.DeleteGlobalRef(android.JObject(app.ref))
}

func (me *nativeApp) run() {
	// theApp.windowPrepared <- struct{}{}
}

func (me *nativeApp) terminate() {
}

func runOnUI(callback func()) {
	// if IsMainThread() {
	// 	callback()
	// } else {
	// 	android.BackToUI(callback)
	// }
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
