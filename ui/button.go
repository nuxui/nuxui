// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
)

type Button interface {
	Text
}

type button text

func NewButton(attrs ...nux.Attr) Button {
	me := NewText(attrs...)
	return me
}
