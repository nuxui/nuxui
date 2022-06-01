// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !release

package nux

import (
	"nuxui.org/nuxui/log"
)

const (
	debug_attr       = true
	debug_size       = true
	debug_register   = true
	debug_hittest    = true
	debug_gesture    = true
	debug_mainthread = true
	debug_event      = true
)

func DebugCheckMainThread() {
	if debug_mainthread && !IsMainThread() {
		log.Fatal("nuxui", "ui operation run out of main thread")
	}
}
