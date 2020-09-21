// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type MixinCreator func() Mixin

type Mixin interface {
	ID() string
	SetID(id string)
	AssignOwner(owner Widget)
	Owner() Widget
}

type WidgetMixin struct {
	id    string
	owner Widget
}

func (me *WidgetMixin) ID() string {
	return me.id
}

func (me *WidgetMixin) SetID(id string) {
	me.id = id
}

// owner can be nil, may be remove from owner
func (me *WidgetMixin) AssignOwner(owner Widget) {
	if me.owner == owner {
		return
	}

	if me.owner != nil && owner != nil && me.owner != owner {
		log.Fatal("nuxui", "The mixin is already have a owner")
	}

	me.owner = owner
}

func (me *WidgetMixin) Owner() Widget {
	if me.owner == nil {
		log.Fatal("nuxui", "The owner should assign before use it")
	}
	return me.owner
}
