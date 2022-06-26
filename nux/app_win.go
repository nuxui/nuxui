// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
)

type nativeApp struct {
	ptr    uintptr
	cursor uintptr
}

var (
	chanBackToUI         = make(chan func())
	gdiplusStartupInput  win32.GdiplusStartupInput
	gdiplusStartupOutput win32.GdiplusStartupOutput
)

func createNativeApp_() *nativeApp {
	win32.GdiplusStartup(&gdiplusStartupInput, nil)
	return &nativeApp{}
}

func (me *nativeApp) run() {
	defer win32.GdiplusShutdown()

	var msg win32.MSG
	var ret int32
	var err error
	for {
		ret, err = win32.GetMessage(&msg, 0, 0, 0)
		if ret > 0 {
			win32.TranslateMessage(&msg)
			win32.DispatchMessage(&msg)
		} else if ret == 0 { // quit
			break
		} else if err != nil {
			log.E("nuxui", "error GetMessage: %s", err.Error())
		}
	}
}

func (me *nativeApp) terminate() {
	win32.PostQuitMessage(0)
}

func invalidateRectAsync_(dirtRect *Rect) {
	win32.RedrawWindow(theApp.window.native().hwnd, nil, 0, win32.RDW_INVALIDATE)
}

func startTextInput() {
	// win32.startTextInput()
}

func stopTextInput() {
	// win32.stopTextInput()
}

func runOnUI(callback func()) {
	if IsMainThread() {
		callback()
	} else {
		go func() {
			chanBackToUI <- callback
		}()
		win32.SendMessage(theApp.window.native().hwnd, win32.WM_USER, 0, 0)
	}
}

func backToUI() {
	callback := <-chanBackToUI
	callback()
}

func currentThreadID_() uint64 {
	return uint64(win32.GetCurrentThreadId())
}

func screenSize() (width, height int32) {
	width, _ = win32.GetSystemMetrics(win32.SM_CXSCREEN)
	height, _ = win32.GetSystemMetrics(win32.SM_CYSCREEN)
	return
}
