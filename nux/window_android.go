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

type nativeWindow struct {
	act android.Activity
}

func newNativeWindow(attr Attr) *nativeWindow {
	log.I("nuxui", "newNativeWindow")

	// width, height := measureWindowSize(attr.GetDimen("width", "50%"), attr.GetDimen("height", "50%"))
	me := &nativeWindow{
		// ptr: darwin.NewNSWindow(width, height),
	}
	android.SetActivityDelegate(me)
	// me.SetTitle(attr.GetString("title", ""))

	runtime.SetFinalizer(me, freeWindow)
	return me
}

func freeWindow(me *nativeWindow) {
}

func (me *nativeWindow) OnCreate(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnCreate")
}

func (me *nativeWindow) OnStart(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnStart")
}

func (me *nativeWindow) OnRestart(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnRestart")
}

func (me *nativeWindow) OnResume(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnResume")
}

func (me *nativeWindow) OnPause(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnPause")
}

func (me *nativeWindow) OnStop(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnStop")
}

func (me *nativeWindow) OnDestroy(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnDestroy")
}

func (me *nativeWindow) OnSurfaceCreated(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceCreated")
}

func (me *nativeWindow) OnSurfaceChanged(activity android.Activity, surfaceHolder android.SurfaceHolder, format, width, height int32) {
	log.I("nuxui", "nativeWindow OnSurfaceChanged")
}

func (me *nativeWindow) OnSurfaceRedrawNeeded(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceRedrawNeeded")
}

func (me *nativeWindow) OnSurfaceDestroyed(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceDestroyed")
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

func nativeActivityEventHandler(event any) bool {

	return false
}
