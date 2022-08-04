// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package nux

import (
	"nuxui.org/nuxui/nux/internal/android"
	"runtime"
)

func loadImageFromFile(filename string) Image {
	me := &nativeImage{
		ref: android.CreateBitmap(filename),
	}
	runtime.SetFinalizer(me, freeImage)
	return me
}

func freeImage(img *nativeImage) {
	img.ref.Recycle()
}

type nativeImage struct {
	ref android.Bitmap
}

func (me *nativeImage) PixelSize() (width, height int32) {
	return me.ref.GetWidth(), me.ref.GetHeight()
}

func (me *nativeImage) Draw(canvas Canvas) {
	w, h := me.PixelSize()
	canvas.native().ref.DrawBitmap(me.ref, 0, 0, float32(w), float32(h), 0)
}
