// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build release

package nux

const (
	debug_attr       = false
	debug_size       = false
	debug_register   = false
	debug_hittest    = false
	debug_gesture    = false
	debug_mainthread = false
	debug_event      = false
)

func DebugCheckMainThread() {}
