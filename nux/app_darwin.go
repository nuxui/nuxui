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
	return &nativeApp{
		ptr: darwin.NSApp_SharedApplication(),
	}
}

func (me *nativeApp) run() {
	me.ptr.Run()
}

func (me *nativeApp) terminate() {
	me.ptr.Terminate()
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
