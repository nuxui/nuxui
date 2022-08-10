// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#cgo LDFLAGS: -landroid -llog

#include <jni.h>
#include <stdlib.h>
#include <unistd.h>
#include <pthread.h>

void  nux_BackToUI();
*/
import "C"

import (
	"nuxui.org/nuxui/nux/internal/android/callfn"
)

var (
	runOnUI = make(chan func())
)

//export go_nux_callMain
func go_nux_callMain(mainPC uintptr) {
	callfn.CallFn(mainPC)
}

func BackToUI(callback func()) {
	go func() {
		runOnUI <- callback
	}()
	C.nux_BackToUI()
}

//export go_nux_backToUI
func go_nux_backToUI() {
	callback := <-runOnUI
	callback()
}

func GetTid() uint64 {
	return uint64(C.gettid())
}
