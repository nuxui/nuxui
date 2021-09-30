// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

/*
#include <windows.h>
#include <windowsx.h>
*/
import "C"
import (
	"github.com/nuxui/nuxui/log"
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-geticoninfoexw

func getCursorScreenPosition() (x, y float32) {
	var p C.POINT
	if C.GetCursorPos(&p) > 0 {
		return float32(p.x), float32(p.y)
	}
	log.E("nux", "windows GetCursorPos faild.")
	return
}

func cursorPositionScreenToWindow(wind Window, x0, y0 float32) (x, y float32) {
	var p C.POINT
	p.x = C.LONG(x0)
	p.y = C.LONG(y0)
	if C.ScreenToClient(wind.(*window).windptr, &p) > 0 {
		return float32(p.x), float32(p.y)
	}
	log.E("nux", "windows ScreenToClient faild.")
	return
}

func cursorPositionWindowToScreen(wind Window, x0, y0 float32) (x, y float32) {
	var p C.POINT
	p.x = C.LONG(x0)
	p.y = C.LONG(y0)
	if C.ClientToScreen(wind.(*window).windptr, &p) > 0 {
		return float32(p.x), float32(p.y)
	}
	log.E("nux", "windows ClientToScreen faild.")
	return
}
