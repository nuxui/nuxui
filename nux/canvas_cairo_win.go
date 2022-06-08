// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && cairo

package nux

import (
	"nuxui.org/nuxui/nux/internal/cairo"
	"runtime"
)

type canvas struct {
	cairo   *cairo.Cairo
	surface *cairo.Surface
}

func canvasFromHDC(hdcBuffer uintptr) *canvas {
	surface := cairo.Win32SurfaceCreate(hdcBuffer)
	me := &canvas{
		surface: surface,
		cairo:   cairo.Create(surface),
	}

	runtime.SetFinalizer(me, freeCanvas)
	return me
}

func freeCanvas(me *canvas) {
	me.cairo.Destroy()
	me.surface.Destroy()
}
