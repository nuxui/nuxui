// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

/*
#cgo LDFLAGS: -llog

#include <jni.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

void backToUI();
int isMainThread();
*/
import "C"
import (
	// "runtime"
	"time"
	// "unsafe"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/callfn"
)

var theApp = &application{
	runOnUI:            make(chan func()),
	nativeLoopPrepared: make(chan struct{}),
	drawSignal:         make(chan struct{}, drawSignalSize),
}

const (
	drawSignalSize = 50
)

func init() {
	// runtime.LockOSThread() // TODO:: RunOnJVM

	timerLoopInstance.init()

	go func() {
		<-theApp.nativeLoopPrepared
		// TODO:: paint not ok?
		log.V("nuxui", "<- nativeLoopPrepared")
		var i, l int
		for {
			<-theApp.drawSignal
			l = len(theApp.drawSignal)
			for i = 0; i != l; i++ {
				<-theApp.drawSignal
			}
			log.V("nuxui", "<-theApp.drawSignal requestRedraw")

			requestRedraw()
			time.Sleep(16 * time.Millisecond)
		}
	}()
}

func app() Application {
	return theApp
}

type application struct {
	manifest           Manifest
	window             *window
	runOnUI            chan func()
	nativeLoopPrepared chan struct{}
	drawSignal         chan struct{}
}

func (me *application) OnCreate(data any) {
}

func (me *application) MainWindow() Window {
	return me.window
}

func (me *application) Manifest() Manifest {
	return me.manifest
}

func (me *application) Terminate() {
}

func (me *application) RequestRedraw(widget Widget) {
	if l := len(theApp.drawSignal); l >= drawSignalSize {
		for i := 0; i != l-1; i++ {
			<-theApp.drawSignal
		}
	}
	theApp.drawSignal <- struct{}{}
}

func run() {
	// defer runtime.UnlockOSThread()
	// if tid := uint64(C.threadID()); tid != initThreadID {
	// 	log.Fatal("nuxui", "main called on thread %d, but init ran on %d", tid, initThreadID)
	// }
}

//export go_nativeLoopPrepared
func go_nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

//export go_callMain
func go_callMain(mainPC uintptr) {
	callfn.CallFn(mainPC)
}

var lastMouseEvent map[int]PointerEvent = map[int]PointerEvent{}

// https://cs.android.com/android/platform/superproject/+/master:frameworks/base/core/java/android/view/MotionEvent.java
//export go_onPointerEvent
func go_onPointerEvent(deviceId, pointerId, action C.jint, x, y C.jfloat) {
	pointerKey := int(pointerId)
	e := &pointerEvent{
		event: event{
			window: theApp.mainWindow,
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  MB_None,
		kind:    Kind_Touch,
		x:       float32(x),
		y:       float32(y),
	}

	switch action {
	case 0 /*MotionEvent.ACTION_DOWN*/ :
		e.action = Action_Down
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[pointerKey] = e
	case 1 /*MotionEvent.ACTION_UP*/ :
		e.action = Action_Up
		if v, ok := lastMouseEvent[pointerKey]; ok {
			e.pointer = v.Pointer()
		} else {
			// can not happend
		}
	case 2 /*MotionEvent.ACTION_MOVE*/ :
		e.action = Action_Drag
		if v, ok := lastMouseEvent[pointerKey]; ok {
			e.pointer = v.Pointer()
		}
	}

	theApp.handleEvent(e)

}

//------------------------ window events ----------------------------

//export go_surfaceCreated
func go_surfaceCreated(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	theApp.mainWindow.jnienv = jnienv
	theApp.mainWindow.activity = activity
	theApp.mainWindow.surfaceHolder = surfaceHolder
	windowAction(activity, Action_WindowCreated)
	log.V("nuxui", "go_surfaceCreated surfaceHolder=%d", surfaceHolder)
}

//export go_surfaceChanged
func go_surfaceChanged(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject, format, width, height C.int) {
	log.V("nuxui", "go_surfaceChanged surfaceHolder=%d", surfaceHolder)
	theApp.mainWindow.jnienv = jnienv
	theApp.mainWindow.activity = activity
	theApp.mainWindow.surfaceHolder = surfaceHolder
	theApp.mainWindow.width = int32(width)
	theApp.mainWindow.height = int32(height)
	windowAction(activity, Action_WindowMeasured)
}

//export go_surfaceRedrawNeeded
func go_surfaceRedrawNeeded(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	log.V("nuxui", "go_surfaceRedrawNeeded surfaceHolder=%d", surfaceHolder)
	theApp.mainWindow.jnienv = jnienv
	theApp.mainWindow.activity = activity
	theApp.mainWindow.surfaceHolder = surfaceHolder
	windowAction(activity, Action_WindowDraw)
}

func windowAction(activity C.jobject, action EventAction) {
	e := &windowEvent{
		event: event{
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: action,
			window: theApp.mainWindow,
		},
	}

	theApp.handleEvent(e)
}

func startTextInput() {
}

func stopTextInput() {
}

func setTextInputRect(x, y, w, h float32) {
}

//export go_backToUI
func go_backToUI() {
	log.V("nuxui", "go_backToUI ..........")
	callback := <-theApp.runOnUI
	callback()
}

func runOnUI(callback func()) {
	if isMainThread() {
		callback()
	} else {
		go func() {
			theApp.runOnUI <- callback
		}()
		C.backToUI()
	}
}

func requestRedraw() {
	log.V("nuxui", "requestRedraw invalidate")
	// C.invalidate()
}

func isMainThread() bool {
	if C.isMainThread() > 0 {
		return true
	}
	return false
}
