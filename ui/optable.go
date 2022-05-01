// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "nuxui.org/nuxui/nux"

type OptChangedCallback func(widget OptableWidget, opted bool, fromUser bool)

type Optable interface {
	SetOpted(opted bool, fromUser bool)
	Opted() bool
	SetOptChangedCallback(listener OptChangedCallback)
}

type OptableWidget interface {
	nux.Widget
	Optable
}
