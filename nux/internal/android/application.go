// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_NuxApplication_instance();
*/
import "C"

type ApplicationDelegate interface {
	OnConfigurationChanged(app Application, newConfig Configuration)
	OnCreate(app Application)
	OnLowMemory(app Application)
	OnTerminate(app Application)
	OnTrimMemory(app Application, level int32)
}

func SetApplicationDelegate(delegate ApplicationDelegate) {
	applicationDelegate = delegate
}

var (
	applicationDelegate ApplicationDelegate
)

func NuxApplication_instance() Application {
	return Application(C.nux_NuxApplication_instance())
}

//export go_NuxApplication_onCreate
func go_NuxApplication_onCreate(application C.jobject, density C.jfloat, densityDpi C.jint, scaledDensity C.jfloat, widthPixels C.jint, heightPixels C.jint, xdpi C.jfloat, ydpi C.jfloat) {
	_displayMetrics = &displayMetrics{
		WidthPixels:   int32(widthPixels),
		HeightPixels:  int32(heightPixels),
		Xdpi:          float32(xdpi),
		Ydpi:          float32(ydpi),
		ScaledDensity: float32(scaledDensity),
		Density:       float32(density),
		DensityDpi:    int32(densityDpi),
	}

	if applicationDelegate != nil {
		applicationDelegate.OnCreate(Application(application))
	}
}

//export go_NuxApplication_onConfigurationChanged
func go_NuxApplication_onConfigurationChanged(application C.jobject, newConfig C.jobject) {
	if applicationDelegate != nil {
		applicationDelegate.OnConfigurationChanged(Application(application), Configuration(newConfig))
	}
}

//export go_NuxApplication_onLowMemory
func go_NuxApplication_onLowMemory(application C.jobject) {
	if applicationDelegate != nil {
		applicationDelegate.OnLowMemory(Application(application))
	}
}

//export go_NuxApplication_onTerminate
func go_NuxApplication_onTerminate(application C.jobject) {
	if applicationDelegate != nil {
		applicationDelegate.OnTerminate(Application(application))
	}
}

//export go_NuxApplication_onTrimMemory
func go_NuxApplication_onTrimMemory(application C.jobject, level C.jint) {
	if applicationDelegate != nil {
		applicationDelegate.OnTrimMemory(Application(application), int32(level))
	}
}
