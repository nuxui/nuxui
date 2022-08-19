// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

import (
	"nuxui.org/nuxui/log"
)

type displayMetrics struct {
	WidthPixels   int32
	HeightPixels  int32
	Xdpi          float32
	Ydpi          float32
	ScaledDensity float32
	Density       float32
	DensityDpi    int32
}

var _displayMetrics *displayMetrics

func GetDisplayMetrics() *displayMetrics {
	if _displayMetrics == nil {
		log.Fatal("nuxui", "application is not created yet")
	}

	return _displayMetrics
}
