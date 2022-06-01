// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
	"runtime"
)

func CreateImage(filename string) Image {
	me := &nativeImage{
		ref: darwin.CGImageSourceCreateImageAtIndex(filename),
	}
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	darwin.CGImageRelease(img.ref)
}

type nativeImage struct {
	ref darwin.CGImage
}

func (me *nativeImage) Size() (width, height int32) {
	return darwin.CGImageGetSize(me.ref)
}
