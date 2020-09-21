// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android

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
	"time"
	"unsafe"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux/internal/callfn"
)

var theApp = &application{
	event:              make(chan Event),
	eventWaitDone:      make(chan Event),
	eventDone:          make(chan struct{}),
	nativeLoopPrepared: make(chan struct{}),
}

func app() Application {
	return theApp
}

type application struct {
	manifest           Manifest
	window             *window
	event              chan Event
	eventWaitDone      chan Event
	eventDone          chan struct{}
	nativeLoopPrepared chan struct{}
}

func (me *application) Creating(attr Attr) {
	if me.manifest == nil {
		me.manifest = NewManifest()
	}

	if c, ok := me.manifest.(Creating); ok {
		c.Creating(attr.GetAttr("manifest", Attr{}))
	}

	if me.window == nil {
		me.window = newWindow()
	}

	me.window.Creating(attr)
}

func (me *application) Created(data interface{}) {
	if c, ok := me.manifest.(AnyCreated); ok {
		c.Created(data)
	}
}

func (me *application) MainWindow() Window {
	return me.window
}

func (me *application) Manifest() Manifest {
	return me.manifest
}

func (me *application) SendEvent(event Event) {
	me.event <- event
}

func (me *application) sendEventAndWaitDone(event Event) {
	me.eventWaitDone <- event
	<-me.eventDone
}

func (me *application) findWindow(activity *C.ANativeActivity, awindow *C.ANativeWindow) Window {
	if me.window.actptr == nil || me.window.windptr == nil {
		return nil
	}

	if me.window.actptr == activity || me.window.windptr == awindow {
		return me.window
	}

	return nil
}

func run() {
	go func() {
		<-theApp.nativeLoopPrepared
		theApp.loop()
	}()
}

//export nativeLoopPrepared
func nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

//export initWindow
func initWindow(activity *C.ANativeActivity) {
	theApp.window.actptr = activity
	// TODO set width height background drawable
}

//export callMain
func callMain(mainPC uintptr) {
	callfn.CallFn(mainPC)
}

//export onInputEvent
func onInputEvent(event *C.AInputEvent) int32 {
	log.V("nuxui", "onInputEvent")
	// log.V("nuxui", "onInputEvent type %d", C.AInputEvent_getType(event))
	// // https://developer.android.com/ndk/reference/group/input#group___input_1gaac34dfe6c6b73b43a4656c9dce041034
	// switch C.AInputEvent_getType(event) {
	// case C.AINPUT_EVENT_TYPE_KEY:
	// case C.AINPUT_EVENT_TYPE_MOTION:
	// 	// events := convertMouseMotionEvent(event)
	// 	// wind := GetWindowByID(0)
	// 	log.V("nuxui", "AINPUT_EVENT_TYPE_MOTION")
	// 	for _, e := range events {
	// 		log.V("nuxui", "DispatchPointerEvent %s", e)
	// 		// ui.DispatchPointerEvent(rootWidget, e)
	// 	}
	// 	// case C.AINPUT_EVENT_TYPE_FOCUS:
	// }
	// // drawFrame()
	// log.V("nuxui", "onInputEvent end")

	e := &inputEvent{
		id:     1,
		time:   time.Now(),
		action: 0,
	}

	theApp.SendEvent(e)
	return 0
}

//export onStart
func onStart(activity *C.ANativeActivity) {
	log.V("nuxui", "onStart")
}

//export onResume
func onResume(activity *C.ANativeActivity) {
	log.V("nuxui", "onResume")
}

//export onPause
func onPause(activity *C.ANativeActivity) {
	log.V("nuxui", "onPause")
}

//export onStop
func onStop(activity *C.ANativeActivity) {
	log.V("nuxui", "onStop")
}

//export onDestroy
func onDestroy(activity *C.ANativeActivity) {
	log.V("nuxui", "go onDestroy")
	e := &event{
		id:    1,
		time:  time.Now(),
		etype: Type_AppExit,
	}
	theApp.sendEventAndWaitDone(e)
	log.V("nuxui", "go onDestroy end")
}

//export onLowMemory
func onLowMemory(activity *C.ANativeActivity) {
	log.V("nuxui", "onLowMemory")
}

//export onSaveInstanceState
func onSaveInstanceState(activity *C.ANativeActivity, outSize *C.size_t) unsafe.Pointer {
	log.V("nuxui", "onSaveInstanceState")
	return nil
}

//export onConfigurationChanged
func onConfigurationChanged(activity *C.ANativeActivity) {
	log.V("nuxui", "onConfigurationChanged")
}

//export onContentRectChanged
func onContentRectChanged(activity *C.ANativeActivity, rect *C.ARect) {
	log.V("nuxui", "onContentRectChanged")
}

//------------------------ window events ----------------------------

//export onNativeWindowCreated
func onNativeWindowCreated(activity *C.ANativeActivity, awindow *C.ANativeWindow) {
	log.V("nuxui", "onNativeWindowCreated")
	theApp.window.actptr = activity
	theApp.window.windptr = awindow

	// on activity created, call window creating, then get window attrs to init window
	// send created
	e := &windowEvent{
		id:     1,
		time:   time.Now(),
		action: Action_WindowCreated,
		window: theApp.findWindow(activity, awindow),
	}

	theApp.sendEventAndWaitDone(e)
}

//export onNativeWindowResized
func onNativeWindowResized(activity *C.ANativeActivity, awindow *C.ANativeWindow) {
	log.V("nuxui", "onNativeWindowResized")
	e := &windowEvent{
		id:     1,
		time:   time.Now(),
		action: Action_WindowMeasured,
		window: theApp.findWindow(activity, awindow),
	}

	theApp.sendEventAndWaitDone(e)
}

//export onNativeWindowRedrawNeeded
func onNativeWindowRedrawNeeded(activity *C.ANativeActivity, awindow *C.ANativeWindow) {
	log.V("nuxui", "onNativeWindowRedrawNeeded")
	e := &windowEvent{
		id:     1,
		time:   time.Now(),
		action: Action_WindowDraw,
		window: theApp.findWindow(activity, awindow),
	}

	theApp.sendEventAndWaitDone(e)
}

//export onWindowFocusChanged
func onWindowFocusChanged(activity *C.ANativeActivity, hasFocus C.int) {
	log.V("nuxui", "onWindowFocusChanged")
	e := &windowEvent{
		id:     1,
		time:   time.Now(),
		action: Action_WindowFocusGained,
		window: theApp.findWindow(activity, nil),
	}

	if hasFocus == 0 {
		e.action = Action_WindowFocusLost
	}

	theApp.sendEventAndWaitDone(e)
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(activity *C.ANativeActivity, awindow *C.ANativeWindow) {
	log.V("nuxui", "onNativeWindowDestroyed")
	e := &windowEvent{
		id:     1,
		time:   time.Now(),
		action: Action_WindowDraw,
		window: theApp.findWindow(activity, awindow),
	}

	theApp.sendEventAndWaitDone(e)
}
