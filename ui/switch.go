// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "nuxui.org/nuxui/nux"

type Switch Check

func NewSwitch(attr nux.Attr) Switch {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := NewCheck(attr)
	return Switch(me)
}

type switcher check
