// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type Align struct {
	Vertical   int
	Horizontal int
}

func NewAlign(attr nux.Attr) *Align {
	a := &Align{}
	a.Init(attr)
	return a
}

func (me *Align) Init(attr nux.Attr) {
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
		log.Fatal("nux", fmt.Sprintf("Unknow alignment '%s' for vertical, only support 'top', 'center', 'bottom'", vertical))
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
		log.Fatal("nux", fmt.Sprintf("Unknow alignment '%s' for horizontal, only support 'left', 'center', 'right'", horizontal))
	}
}
