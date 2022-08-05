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
uint64_t currentThreadID();
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
