// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

type nativeApp struct {
}

func createNativeApp_() *nativeApp {
	return &nativeApp{}
}

func (me *nativeApp) init() {

}

func (me *nativeApp) run() {
	alive := make(chan bool)
	<-alive
}

func (me *nativeApp) terminate() {
}

func runOnUI(callback func()) {

}

func invalidateRectAsync_(dirtRect *Rect) {
}

func startTextInput() {
}

func stopTextInput() {
}

func currentThreadID_() uint64 {
	return 0
}

func screenSize() (width, height int32) {
	return 800, 600
}
