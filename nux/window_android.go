// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
	"time"
)

type nativeWindow struct {
	act           android.Activity
	surfaceHolder android.SurfaceHolder
	penddingAttr  Attr
}

func newNativeWindow(attr Attr) *nativeWindow {
	me := &nativeWindow{
		penddingAttr: attr,
	}
	android.SetActivityDelegate(me)
	runtime.SetFinalizer(me, freeWindow)
	return me
}

func freeWindow(me *nativeWindow) {
}

func (me *nativeWindow) OnCreate(activity android.Activity) {
	log.I("nuxui", "nativeWindow OnCreate")
	// init with me.penddingAttr
	me.act = android.Activity(android.NewGlobalRef(android.JObject(activity)))
	if title := me.penddingAttr.GetString("title", ""); title != "" {
		activity.SetTitle(title)
	}

	// theApp.mainWindow.mountWidget()
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
	// theApp.mainWindow.ejectWidget()
	android.DeleteGlobalRef(android.JObject(me.act))
}

func (me *nativeWindow) OnSurfaceCreated(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceCreated")
	// init with me.penddingAttr
	me.surfaceHolder = android.SurfaceHolder(android.NewGlobalRef(android.JObject(surfaceHolder)))

	me.penddingAttr = nil
}

func (me *nativeWindow) OnSurfaceChanged(activity android.Activity, surfaceHolder android.SurfaceHolder, format, width, height int32) {
	log.I("nuxui", "nativeWindow OnSurfaceChanged")
	theApp.mainWindow.resize()
}

func (me *nativeWindow) OnSurfaceRedrawNeeded(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceRedrawNeeded")
	theApp.mainWindow.draw()
}

func (me *nativeWindow) OnSurfaceDestroyed(activity android.Activity, surfaceHolder android.SurfaceHolder) {
	log.I("nuxui", "nativeWindow OnSurfaceDestroyed")
}

func (me *nativeWindow) OnTouch(activity android.Activity, event android.MotionEvent) bool {
	// log.I("nuxui", "nativeWindow OnTouch")
	if event.GetPointerCount() == 1 {
		return handlePointerEvent(event)
	}
	return false
}

func (me *nativeWindow) Center() {
}

func (me *nativeWindow) Show() {
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	dm := android.GetDisplayMetrics()
	return dm.WidthPixels, dm.HeightPixels
}

func (me *nativeWindow) Title() string {
	return ""
}

func (me *nativeWindow) SetTitle(title string) {
}

func (me *nativeWindow) lockCanvas() Canvas {
	return newCanvas(me.surfaceHolder.LockCanvas())
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	me.surfaceHolder.UnlockCanvas(canvas.native().ref)
	android.DeleteGlobalRef(android.JObject(canvas.native().ref))
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

var lastMouseEvent map[PointerButton]PointerEvent = map[PointerButton]PointerEvent{}

func handlePointerEvent(mevent android.MotionEvent) bool {
	x, y := mevent.GetXY(0)
	action := mevent.GetActionMasked()

	e := &pointerEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  ButtonPrimary,
		kind:    Kind_Touch,
		x:       x,
		y:       y,
	}

	switch action {
	case android.ACTION_DOWN, android.ACTION_POINTER_DOWN:
		e.action = Action_Down
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case android.ACTION_MOVE:
		e.action = Action_Drag
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case android.ACTION_UP, android.ACTION_POINTER_UP:
		e.action = Action_Up
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	case android.ACTION_CANCEL:
		e.action = Action_Up
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	}

	return App().MainWindow().handlePointerEvent(e)
}
