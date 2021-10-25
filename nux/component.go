// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type Component interface {
	Widget
	SetContent(Widget)
	Content() Widget
}

func NewComponentBase(ctx Context, owner Widget, attrs ...Attr) *ComponentBase {
	attr := Attr{}
	if len(attrs) > 0 {
		attr = attrs[0]
	}
	return &ComponentBase{
		id: attr.GetString("id", ""),
	}
}

type ComponentBase struct {
	id      string
	parent  Parent
	content Widget
}

func (me *ComponentBase) SetID(id string) {
	me.id = id
}

func (me *ComponentBase) ID() string {
	return me.id
}

func (me *ComponentBase) Parent() Parent {
	return me.parent
}

// parent can be nil, may be remove form parent
func (me *ComponentBase) AssignParent(parent Parent) {
	if me.parent == nil {
		me.parent = parent
	} else {
		log.Fatal("nuxui", "The parent of widget '%s' is already assigned, can not assign again.", me.ID())
	}
}

func (me *ComponentBase) SetContent(content Widget) {
	me.content = content
	// content.AssignParent(me)
}

func (me *ComponentBase) Content() Widget {
	return me.content
}

// type Component interface {
// 	Parent
// 	Component() Widget
// 	Content() Widget
// }

// type component struct {
// 	WidgetParent
// 	component Widget
// 	content   Widget
// }

// /*
// NewComponent component, child not nil
// */
// func NewComponent(compt, content Widget) Component {
// 	if compt == nil {
// 		log.Fatal("nuxui", "component can not ne nil")
// 	}

// 	if content == nil {
// 		log.Fatal("nuxui", "child can not ne nil")
// 	}

// 	me := &component{
// 		component: compt,
// 		content:   content,
// 	}
// 	me.WidgetParent.Owner = me
// 	content.AssignParent(me)
// 	return me
// }

// func (me *component) Component() Widget {
// 	return me.component
// }

// func (me *component) Content() Widget {
// 	return me.content
// }
