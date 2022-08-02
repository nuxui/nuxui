// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/ios"
	"runtime"
)

func loadImageFromFile(filename string) Image {
	img := ios.UIImage_ImageNamed(filename)
	me := &nativeImage{
		ptr: ios.CGImageCreateCopy(img.CGImage()),
	}

	// crash ??
	// me := &nativeImage{
	// 	ptr: ios.CGImageSourceCreateImageAtIndex(filename),
	// }
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	ios.CGImageRelease(img.ptr)
}

type nativeImage struct {
	ptr ios.CGImageRef
}

func (me *nativeImage) PixelSize() (width, height int32) {
	return ios.CGImageGetSize(me.ptr)
}

func (me *nativeImage) Draw(canvas Canvas) {
	w, h := me.PixelSize()
	ios.CGContextDrawImage(canvas.native().ctx, ios.CGRectMake(0, 0, float32(w), float32(h)), me.ptr)
}
