// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"runtime"
	"syscall"

	"nuxui.org/nuxui/nux/internal/win32"
)

func CreateImage(path string) Image {
	me := &nativeImage{}
	str, _ := syscall.UTF16PtrFromString(path)
	win32.GdipLoadImageFromFile(str, &me.ptr)
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	win32.GdipDisposeImage(img.ptr)
}

type nativeImage struct {
	ptr *win32.GpImage
}

func (me *nativeImage) Size() (width, height int32) {
	if me.ptr == nil {
		return 0, 0
	}
	var w, h uint32
	win32.GdipGetImageWidth(me.ptr, &w)
	win32.GdipGetImageHeight(me.ptr, &h)
	return int32(w), int32(h)
}
