// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && !cairo

package nux

import (
	"path/filepath"
	"runtime"
	"syscall"

	"nuxui.org/nuxui/nux/internal/win32"
)

func loadImageFromFile(path string) Image {
	path, _ = filepath.Abs(path)
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

func (me *nativeImage) PixelSize() (width, height int32) {
	if me.ptr == nil {
		return 0, 0
	}
	var w, h uint32
	win32.GdipGetImageWidth(me.ptr, &w)
	win32.GdipGetImageHeight(me.ptr, &h)
	return int32(w), int32(h)
}

func (me *nativeImage) Draw(canvas Canvas) {
	w, h := me.PixelSize()
	win32.GdipDrawImageRect(canvas.native().ptr, me.ptr, 0, 0, float32(w), float32(h))
}
