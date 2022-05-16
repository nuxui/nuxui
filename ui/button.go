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
	if !attr.Has("theme") {
		attr = nux.MergeAttrs(button_theme(nux.ThemeLight), attr)
	}

	me := NewLabel(attr)
	return Button(me)
}
