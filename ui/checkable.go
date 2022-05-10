// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "nuxui.org/nuxui/nux"

type CheckChangedCallback func(widget CheckableWidget, checked bool, fromUser bool)

type Checkable interface {
	SetChecked(checked bool, fromUser bool)
	Checked() bool
	SetCheckChangedCallback(listener CheckChangedCallback)
}

type CheckableWidget interface {
	nux.Widget
	Checkable
}
