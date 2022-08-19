// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_surfaceHolder_lockCanvas(jobject surfaceHolder);
void    nux_surfaceHolder_unlockCanvas(jobject surfaceHolder, jobject canvas);
void    nux_NuxActivity_invalidateAsync(jobject activity);
*/
import "C"

type ActivityDelegate interface {
	OnCreate(activity Activity)
	OnStart(activity Activity)
	OnRestart(activity Activity)
	OnResume(activity Activity)
	OnPause(activity Activity)
	OnStop(activity Activity)
	OnDestroy(activity Activity)

	OnSurfaceCreated(activity Activity, surfaceHolder SurfaceHolder)
	OnSurfaceChanged(activity Activity, surfaceHolder SurfaceHolder, format, width, height int32)
	OnSurfaceRedrawNeeded(activity Activity, surfaceHolder SurfaceHolder)
	OnSurfaceDestroyed(activity Activity, surfaceHolder SurfaceHolder)

	OnTouch(activity Activity, event MotionEvent) bool
}

func SetActivityDelegate(delegate ActivityDelegate) {
	activityDelegate = delegate
}

var (
	activityDelegate ActivityDelegate
)

func (me Activity) InvalidateAsync() {
	C.nux_NuxActivity_invalidateAsync(C.jobject(me))
}

//export go_NuxActivity_onCreate
func go_NuxActivity_onCreate(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnCreate(Activity(activity))
	}
}

//export go_NuxActivity_onStart
func go_NuxActivity_onStart(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnStart(Activity(activity))
	}
}

//export go_NuxActivity_onRestart
func go_NuxActivity_onRestart(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnRestart(Activity(activity))
	}
}

//export go_NuxActivity_onResume
func go_NuxActivity_onResume(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnResume(Activity(activity))
	}
}

//export go_NuxActivity_onPause
func go_NuxActivity_onPause(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnPause(Activity(activity))
	}
}

//export go_NuxActivity_onStop
func go_NuxActivity_onStop(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnStop(Activity(activity))
	}
}

//export go_NuxActivity_onDestroy
func go_NuxActivity_onDestroy(activity C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnDestroy(Activity(activity))
	}
}

//export go_NuxActivity_surfaceCreated
func go_NuxActivity_surfaceCreated(activity C.jobject, surfaceHolder C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnSurfaceCreated(Activity(activity), SurfaceHolder(surfaceHolder))
	}
}

//export go_NuxActivity_surfaceChanged
func go_NuxActivity_surfaceChanged(activity C.jobject, surfaceHolder C.jobject, format, width, height C.int) {
	if activityDelegate != nil {
		activityDelegate.OnSurfaceChanged(Activity(activity), SurfaceHolder(surfaceHolder), int32(format), int32(width), int32(height))
	}
}

//export go_NuxActivity_surfaceRedrawNeeded
func go_NuxActivity_surfaceRedrawNeeded(activity C.jobject, surfaceHolder C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnSurfaceRedrawNeeded(Activity(activity), SurfaceHolder(surfaceHolder))
	}
}

//export go_NuxActivity_surfaceDestroyed
func go_NuxActivity_surfaceDestroyed(activity C.jobject, surfaceHolder C.jobject) {
	if activityDelegate != nil {
		activityDelegate.OnSurfaceDestroyed(Activity(activity), SurfaceHolder(surfaceHolder))
	}
}

// https://cs.android.com/android/platform/superproject/+/master:frameworks/base/core/java/android/view/MotionEvent.java
//export go_NuxActivity_onTouch
func go_NuxActivity_onTouch(activity, motionEvent C.jobject) C.int {
	if activityDelegate != nil {
		if activityDelegate.OnTouch(Activity(activity), MotionEvent(motionEvent)) {
			return C.int(1)
		}
	}
	return C.int(0)
}

func (me Activity) SetTitle(title string) {

}

// return Canvas: GlobalRef
func (me SurfaceHolder) LockCanvas() Canvas {
	return Canvas(C.nux_surfaceHolder_lockCanvas(C.jobject(me)))
}

// canvas GlobalRef delete by caller
func (me SurfaceHolder) UnlockCanvas(canvas Canvas) {
	C.nux_surfaceHolder_unlockCanvas(C.jobject(me), C.jobject(canvas))
}
