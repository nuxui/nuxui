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

import (
	"nuxui.org/nuxui/log"
)

type ApplicationDelegate interface {
	OnConfigurationChanged(app Application, newConfig Configuration)
	OnCreate(app Application)
	OnLowMemory(app Application)
	OnTerminate(app Application)
	OnTrimMemory(app Application, level int32)

	OnCreateWindow(activity Activity)
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
func go_NuxApplication_onCreate(application C.jobject) {
	log.I("nuxui", "go_NuxApplication_onCreate  ==  0")
	if applicationDelegate != nil {
		applicationDelegate.OnCreate(Application(application))
	}
	log.I("nuxui", "go_NuxApplication_onCreate  ==  1")
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
