// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
)

type CheckBox interface {
	Label
}

type checkbox label

func NewCheckBox(attrs ...nux.Attr) CheckBox {
	a := nux.MergeAttrs(attrs...)
	me := NewLabel(a)
	return CheckBox(me)
}
