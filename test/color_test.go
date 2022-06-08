// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

func TestDimen(t *testing.T) {
	color, err :=nux.ParseColor("#000022", 0xffffffff)
	log.I("test", "rgba: #%.8X, argb: #%.8X %s", color, color.ARGB(), err)
	color, err =nux.ParseColor("0x123456", 0xffffffff)
	log.I("test", "rgba: #%.8X, argb: #%.8X %s", color, color.ARGB(), err)
	color, err =nux.ParseColor("0x123456", 0xffffffff)
	log.I("test", "rgba: #%.8X, argb: #%.8X %s", color, color.ARGB(), err)
	log.I("test", "rgba: #%.8X, argb: #%.8X %s", (0xffffffff>>8), ((0xff<<24)), err)
}
