// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
	"runtime"
)

func createImage(filename string) Image {
	me := &nativeImage{
		ptr: darwin.CGImageSourceCreateImageAtIndex(filename),
	}
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	darwin.CGImageRelease(img.ptr)
}

type nativeImage struct {
	ptr darwin.CGImage
}

func (me *nativeImage) PixelSize() (width, height int32) {
	return darwin.CGImageGetSize(me.ptr)
}

func (me *nativeImage) Draw(canvas Canvas) {
	w, h := me.PixelSize()
	darwin.CGContextDrawImage(canvas.native().ctx, darwin.CGRectMake(0, 0, float32(w), float32(h)), me.ptr)
}
