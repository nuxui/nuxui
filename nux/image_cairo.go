// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || (windows && cairo)

package nux

import (
	"nuxui.org/nuxui/nux/internal/cairo"
	"path/filepath"
	"runtime"
	"strings"
)

func createImage(path string) Image {
	path, _ = filepath.Abs(path)
	ext := strings.ToLower(filepath.Ext(path))

	var img *nativeImage
	switch ext {
	case ".png":
		img = &nativeImage{ptr: cairo.ImageSurfaceCreateFromPNG(path)}
	case ".jpg", ".jpeg":
		img = &nativeImage{ptr: cairo.ImageSurfaceCreateFromJPEG(path)}
	}

	runtime.SetFinalizer(img, freeImage)
	return img
}

func freeImage(img *nativeImage) {
	img.ptr.Destroy()
}

type nativeImage struct {
	ptr *cairo.Surface
}

func (me *nativeImage) PixelSize() (width, height int32) {
	if me.ptr == nil {
		return 0, 0
	}
	return cairo.ImageSurfaceGetWidth(me.ptr), cairo.ImageSurfaceGetHeight(me.ptr)
}

func (me *nativeImage) Draw(canvas Canvas) {
	canvas.native().cairo.SetSourceSurface(me.ptr, 0, 0)
	canvas.native().cairo.Paint()
}
