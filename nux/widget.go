// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
)

type Creator func(Context, ...Attr) Widget

type Widget interface {
	ID() string
	SetID(string)
	Parent() Parent
	AssignParent(parent Parent)
	// TODO:: attachToWindow?
}

type WidgetBase struct {
	id     string
	parent Parent
}

func NewWidgetBase(ctx Context, owner Widget, attrs ...Attr) *WidgetBase {
	attr := Attr{}
	if len(attrs) > 0 {
		attr = attrs[0]
	}
	return &WidgetBase{
		id: attr.GetString("id", ""),
	}
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
		log.Fatal("nuxui", "The parent of widget '%s' is already assigned, can not assign again.", me.ID())
	}
}

type Template interface {
	Template() string
}

type Render interface {
	Render() Widget
}

type viewfuncs interface {
	Measure
	Layout
	Draw
}
