// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

package nux

import ()

func loadImageFromFile(filename string) Image {
	me := &nativeImage{}
	return me
}

type nativeImage struct {
}

func (me *nativeImage) PixelSize() (width, height int32) {
	return 400, 400
}

func (me *nativeImage) Draw(canvas Canvas) {
}
