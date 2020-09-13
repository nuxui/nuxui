// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"

	"github.com/nuxui/nuxui/log"
)

type Creator func() Widget

type Widget interface {
	ID() string
	SetID(string)
	Parent() Parent
	AssignParent(parent Parent)
}

type WidgetBase struct {
	id     string
	parent Parent
}

func (me *WidgetBase) SetID(id string) {
	me.id = id
}

func (me *WidgetBase) ID() string {
	return me.id
}

func (me *WidgetBase) Parent() Parent {
	return me.parent
}

// parent can be nil, may be remove form parent
func (me *WidgetBase) AssignParent(parent Parent) {
	if me.parent == nil {
		me.parent = parent
	} else {
		log.Fatal("nux", fmt.Sprintf("The parent of widget '%s' is already assigned, can not assign again.", me.ID()))
	}
}

func (me *WidgetBase) Creating(attr Attr) {
	me.id = attr.GetString("id", "")
}

type Template interface {
	Template() string
}

type Render interface {
	Render() Widget
}
