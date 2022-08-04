// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#cgo LDFLAGS: -llog

#include <jni.h>
#include <stdlib.h>
#include <unistd.h>
#include <pthread.h>

void nux_BackToUI();
*/
import "C"

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/android/callfn"
)

var (
	runOnUI = make(chan func())
)

//export go_nux_callMain
func go_nux_callMain(mainPC uintptr) {
	log.I("nuxui", "go_nux_callMain  == 0")
	callfn.CallFn(mainPC)
	log.I("nuxui", "go_nux_callMain  ==  1")
}

//export go_nux_surfaceCreated
func go_nux_surfaceCreated(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	// theApp.mainWindow.jnienv = jnienv
	// theApp.mainWindow.activity = activity
	// theApp.mainWindow.surfaceHolder = surfaceHolder
	// windowAction(activity, Action_WindowCreated)
	// log.V("nuxui", "go_surfaceCreated surfaceHolder=%d", surfaceHolder)
}

//export go_nux_surfaceChanged
func go_nux_surfaceChanged(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject, format, width, height C.int) {
	// log.V("nuxui", "go_surfaceChanged surfaceHolder=%d", surfaceHolder)
	// theApp.mainWindow.jnienv = jnienv
	// theApp.mainWindow.activity = activity
	// theApp.mainWindow.surfaceHolder = surfaceHolder
	// theApp.mainWindow.width = int32(width)
	// theApp.mainWindow.height = int32(height)
	// windowAction(activity, Action_WindowMeasured)
}

//export go_nux_surfaceRedrawNeeded
func go_nux_surfaceRedrawNeeded(jnienv *C.JNIEnv, activity C.jobject, surfaceHolder C.jobject) {
	// log.V("nuxui", "go_surfaceRedrawNeeded surfaceHolder=%d", surfaceHolder)
	// theApp.mainWindow.jnienv = jnienv
	// theApp.mainWindow.activity = activity
	// theApp.mainWindow.surfaceHolder = surfaceHolder
	// windowAction(activity, Action_WindowDraw)
}

// https://cs.android.com/android/platform/superproject/+/master:frameworks/base/core/java/android/view/MotionEvent.java
//export go_nux_onPointerEvent
func go_nux_onPointerEvent(deviceId, pointerId, action C.jint, x, y C.jfloat) {

}

func BackToUI(callback func()) {
	go func() {
		runOnUI <- callback
	}()
	C.nux_BackToUI()
}

//export go_nux_backToUI
func go_nux_backToUI() {
	// log.V("nuxui", "go_nux_backToUI ..........")
	callback := <-runOnUI
	callback()
}

func GetTid() uint64 {
	return uint64(C.gettid())
}
