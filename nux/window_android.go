// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
)

type nativeWindow struct {
	act android.Activity
}

func newNativeWindow(attr Attr) *nativeWindow {
	// darwin.SetWindowEventHandler(nativeWindowEventHandler)

	// width, height := measureWindowSize(attr.GetDimen("width", "50%"), attr.GetDimen("height", "50%"))
	me := &nativeWindow{
		// ptr: darwin.NewNSWindow(width, height),
	}
	// me.SetTitle(attr.GetString("title", ""))

	// runtime.SetFinalizer(me, freeWindow)
	return me
}

func freeWindow(me *nativeWindow) {
}

func (me *nativeWindow) Center() {
}

func (me *nativeWindow) Show() {
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	return 100, 100
}

func (me *nativeWindow) Title() string {
	return ""
}

func (me *nativeWindow) SetTitle(title string) {
}

func (me *nativeWindow) lockCanvas() Canvas {
	return newCanvas(0)
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	// canvas.Flush()
}

func (me *nativeWindow) draw(canvas Canvas, decor Widget) {
	if decor != nil {
		if f, ok := decor.(Draw); ok {
			canvas.Save()
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
		}
	}
}

func nativeWindowEventHandler(event any) bool {

	return false
}
