// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	// "nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/linux"
	"nuxui.org/nuxui/nux/internal/linux/xlib"
	"unsafe"
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
	if xlib.XSupportsLocale() {
		xlib.XSetLocaleModifiers("")
	}

	return &nativeApp{
		display: xlib.XOpenDisplay(xlib.GetDisplayName()),
	}
}

func (me *nativeApp) init() {

}

func (me *nativeApp) run() {
	defer xlib.XCloseDisplay(me.display)

	theApp.mainWindow = newWindow(theApp.manifest.GetAttr("mainWindow", nil))
	theApp.mainWindow.Center()
	theApp.mainWindow.Show()

	theApp.windowPrepared <- struct{}{}

	var event xlib.XEvent
	for {
		xlib.XNextEvent(me.display, &event)
		// log.I("nuxui", "### event type %d, %T, %s", event.Type(),event.Convert(), event.Convert())
		// log.I("nuxui", "### event type %d, %T", event.Type(), event.Convert())

		if xlib.XFilterEvent(&event, xlib.None) {
			continue
		}

		if theApp.mainWindow.native().handleNativeEvent(event.Convert()) {

		}
	}
}

func (me *nativeApp) terminate() {
}

var chanRunOnUI = make(chan func())

func runOnUI(callback func()) {
	if IsMainThread() {
		callback()
	} else {
		go func() {
			chanRunOnUI <- callback
		}()

		event := &xlib.XClientMessageEvent{
			Type:        xlib.ClientMessage,
			Serial:      0,
			SendEvent:   1,
			Display:     theApp.native.display,
			Window:      theApp.mainWindow.native().window,
			MessageType: xlib.XInternAtom(theApp.native.display, "nux_user_backToUI", false),
			Format:      32,
		}
		xlib.XSendEvent(event.Display, event.Window, true, xlib.NoEventMask, (*xlib.XEvent)(unsafe.Pointer(event)))
	}
}

func backToUI() {
	callback := <-chanRunOnUI
	callback()
}

func invalidateRectAsync_(dirtRect *Rect) {
	nw := theApp.mainWindow.native()
	w, h := nw.ContentSize()
	xlib.XClearArea(nw.display, nw.window, 0, 0, w, h, true)
	// xlib.XClearArea(nw.display, nw.window, 0, 0, 0, 0, true)
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

func screenSize() (width, height int32) {
	display := theApp.native.display
	screenNum := xlib.XDefaultScreen(display)
	return xlib.XDisplayWidth(display, screenNum), xlib.XDisplayHeight(display, screenNum)
}
