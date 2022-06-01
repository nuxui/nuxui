// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/linux"
	"nuxui.org/nuxui/nux/internal/linux/xlib"
)

type nativeApp struct {
	display *xlib.Display
}

func createNativeApp_() *nativeApp {
	if linux.GetLocale(linux.LC_CTYPE) == "C" {
		linux.SetLocale(linux.LC_CTYPE, "")
	}

	xlib.XInitThreads()
	xlib.XrmInitialize()
	if xlib.XSupportsLocale(){
		xlib.XSetLocaleModifiers("")
	}

	return &nativeApp{
		display: xlib.XOpenDisplay(xlib.GetDisplayName()),
	}
}

func (me *nativeApp) run() {
	defer xlib.XCloseDisplay(me.display)

	var event xlib.XEvent
	for {
		xlib.XNextEvent(me.display, &event)
		log.I("nuxui", "### event type %d, %T", event.Type(), event.Convert())

		if xlib.XFilterEvent(&event, xlib.None) {
			continue
		}
	}
}

func (me *nativeApp) terminate() {
}


func runOnUI(callback func()) {
	// if IsMainThread() {
	// 	callback()
	// } else {
	// 	darwin.BackToUI(callback)
	// }
}

func invalidateRectAsync_(dirtRect *Rect) {
	// TODO:: error for render radio options
	// darwin.NSApp().KeyWindow().InvalidateRectAsync(float32(rect.X), float32(rect.Y), float32(rect.Width), float32(rect.Height))
	// darwin.NSApp().KeyWindow().InvalidateRectAsync(0, 0, 0, 0)
}

func startTextInput() {
	// darwin.NSApp().KeyWindow().StartTextInputAsync()
}

func stopTextInput() {
	// darwin.NSApp().KeyWindow().StopTextInputAsync()
}

func currentThreadID_() uint64 {
	return linux.CurrentThreadID()
}