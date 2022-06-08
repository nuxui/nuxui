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
	"unicode"
)

func CreateImage(path string) Image {
	path, _ = filepath.Abs(path)
	ext := strings.ToLowerSpecial(unicode.TurkishCase, filepath.Ext(path))

	var img *nativeImage
	switch ext {
	case ".png":
		img = &nativeImage{ptr: cairo.ImageSurfaceCreateFromPNG(path)}
	case ".svg":
		img = &nativeImage{ptr: cairo.SVGSurfaceCreate(path, 400, 400)}
		cairo.SVGSurfaceSetDocumentUnit(img.ptr, cairo.CAIRO_SVG_UNIT_PX)
		return img
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

func (me *nativeImage) Size() (width, height int32) {
	if me.ptr == nil {
		return 0, 0
	}
	return cairo.ImageSurfaceGetWidth(me.ptr), cairo.ImageSurfaceGetHeight(me.ptr)
}
