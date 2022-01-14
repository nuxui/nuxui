// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

/*
#cgo LDFLAGS: -landroid -llog

#include <android/configuration.h>
#include <android/input.h>
#include <android/keycodes.h>
#include <android/looper.h>
#include <android/native_activity.h>
#include <android/native_window.h>
#include <EGL/egl.h>
#include <jni.h>
#include <pthread.h>
#include <stdlib.h>
#include <string.h>

*/
import "C"
import (
	// "runtime"
	"time"
	"unsafe"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux/internal/callfn"
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

func (me *application) OnCreate(data interface{}) {
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

//export go_initWindow
func go_initWindow(activity C.jobject) {
	theApp.window.activity = activity
	// TODO set width height background drawable
}

//export go_callMain
func go_callMain(mainPC uintptr) {
	callfn.CallFn(mainPC)
}

//export go_onStart
func go_onStart(activity C.jobject) {
	log.V("nuxui", "onStart")
}

//export go_onResume
func go_onResume(activity C.jobject) {
	log.V("nuxui", "onResume")
}

//export go_onPause
func go_onPause(activity C.jobject) {
	log.V("nuxui", "onPause")
}

//export go_onStop
func go_onStop(activity C.jobject) {
	log.V("nuxui", "onStop")
}

//export go_onDestroy
func go_onDestroy(activity C.jobject) {
	log.V("nuxui", "go onDestroy")
}

//export go_onLowMemory
func go_onLowMemory(activity C.jobject) {
	log.V("nuxui", "onLowMemory")
}

//export go_onSaveInstanceState
func go_onSaveInstanceState(activity C.jobject, outSize *C.size_t) unsafe.Pointer {
	log.V("nuxui", "onSaveInstanceState")
	return nil
}

//export go_onConfigurationChanged
func go_onConfigurationChanged(activity C.jobject) {
	log.V("nuxui", "onConfigurationChanged")
}

//export go_onContentRectChanged
func go_onContentRectChanged(activity C.jobject, rect *C.ARect) {
	log.V("nuxui", "onContentRectChanged")
}

//------------------------ window events ----------------------------

//export go_onNativeWindowCreated
func go_onNativeWindowCreated(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	theApp.window.jnienv = jnienv
	theApp.window.activity = activity
	theApp.window.surfaceHolder = surfaceHolder
	windowAction(activity, Action_WindowCreated)
	log.V("nuxui", "windowAction Action_WindowCreated surfaceHolder=%d", surfaceHolder)
}

//export go_onNativeWindowResized
func go_onNativeWindowResized(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject, format, width, height C.int) {
	log.V("nuxui", "onNativeWindowResized surfaceHolder=%d", surfaceHolder)
	theApp.window.jnienv = jnienv
	theApp.window.activity = activity
	theApp.window.surfaceHolder = surfaceHolder
	theApp.window.width = int32(width)
	theApp.window.height = int32(height)
	windowAction(activity, Action_WindowMeasured)
}

//export go_onNativeWindowRedrawNeeded
func go_onNativeWindowRedrawNeeded(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	log.V("nuxui", "onNativeWindowRedrawNeeded surfaceHolder=%d", surfaceHolder)
	theApp.window.jnienv = jnienv
	theApp.window.activity = activity
	theApp.window.surfaceHolder = surfaceHolder
	windowAction(activity, Action_WindowDraw)
}

func windowAction(activity C.jobject, action EventAction) {
	e := &windowEvent{
		event: event{
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: action,
			window: theApp.window,
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
	// go func() {
	// 	theApp.runOnUI <- callback
	// }()
	// C.backToUI()
}

func requestRedraw() {
	log.V("nuxui", "requestRedraw invalidate")
	// C.invalidate()
}

//export go_onWindowFocusChanged
func go_onWindowFocusChanged(activity C.jobject, hasFocus C.int) {
	log.V("nuxui", "onWindowFocusChanged")
	// e := &windowEvent{
	// 	id:     1,
	// 	time:   time.Now(),
	// 	action: Action_WindowFocusGained,
	// 	window: theApp.findWindow(activity, nil),
	// }

	// if hasFocus == 0 {
	// 	e.action = Action_WindowFocusLost
	// }

	// theApp.SendEventAndWait(e)
}

//export go_onNativeWindowDestroyed
func go_onNativeWindowDestroyed(activity C.jobject, awindow *C.ANativeWindow) {
	log.V("nuxui", "onNativeWindowDestroyed")
	// e := &windowEvent{
	// 	id:     1,
	// 	time:   time.Now(),
	// 	action: Action_WindowDraw,
	// 	window: theApp.findWindow(activity, awindow),
	// }

	// theApp.SendEventAndWait(e)
}
