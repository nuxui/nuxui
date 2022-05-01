// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

type Align struct {
	Vertical   int
	Horizontal int
}

// TODO:: center|left
func NewAlign(attr nux.Attr) *Align {
	me := &Align{}
	vertical := attr.GetString("vertical", "top")
	horizontal := attr.GetString("horizontal", "left")

	switch vertical {
	case "top":
		me.Vertical = Top
	case "center":
		me.Vertical = Center
	case "bottom":
		me.Vertical = Bottom
	default:
		me.Vertical = Top
		log.Fatal("nuxui", "Unknow alignment '%s' for vertical, only support 'top', 'center', 'bottom'", vertical)
	}

	switch horizontal {
	case "left":
		me.Horizontal = Left
	case "center":
		me.Horizontal = Center
	case "right":
		me.Horizontal = Right
	default:
		me.Horizontal = Left
		log.Fatal("nuxui", "Unknow alignment '%s' for horizontal, only support 'left', 'center', 'right'", horizontal)
	}
	return me
}
