// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

// TODO not a good name

type ComptHelper interface {
	Component() Widget
	AssignComponent(compt Widget)
}

type ComptHelperPart struct {
	compt Widget
}

func (me *ComptHelperPart) Component() Widget {
	if me.compt == nil {
		log.Fatal("nux", "The component is nil, should set it before use")
	}
	return me.compt
}

func (me *ComptHelperPart) AssignComponent(compt Widget) {
	if me.compt != nil {
		log.Fatal("nux", "The component is already assigned.")
	}
	me.compt = compt
}
