// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/nux"
)

type Button interface {
	Label
}

type button label

func NewButton(attr nux.Attr) Button {
	btnattr := nux.Attr{
		"selectable": false,
		"clickable":  true,
	}
	nux.MergeAttrs(btnattr, attr)
	me := NewLabel(btnattr)
	return Button(me)
}
