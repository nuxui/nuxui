// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
)

// TODO
// child.BringToFront(){
// child.Parent().BringChildToFront(child)
// }

type Parent interface {
	Widget
	InsertChild(index int, child ...Widget)
	AddChild(child ...Widget)
	RemoveChildAt(index int)
	RemoveChild(child Widget) (index int)
	ReplaceChild(src, dest Widget) (index int)
	Children() []Widget // the order of Children is from bottom to top
	ChildrenCount() int
	ChildAt(index int) Widget
	Find(id string) Widget
}

type WidgetParent struct {
	owner Parent
	*WidgetBase
	children []Widget
}

func NewWidgetParent(owner Parent, attr Attr) *WidgetParent {
	me := &WidgetParent{
		owner:      owner,
		WidgetBase: NewWidgetBase(attr),
		children:   []Widget{},
	}

	return me
}

func (me *WidgetParent) InsertChild(index int, child ...Widget) {
	if index < 0 || index > len(me.children) {
		log.Fatal("nuxui", "index out of range when insert child to parent")
	} else {
		c := make([]Widget, len(me.children)+len(child))
		c = append(c, me.children[:index]...)
		c = append(c, child...)
		c = append(c, me.children[index:]...)
		me.children = c

		for _, c := range child {
			MountWidget(c, me.owner)
		}
	}
}

func (me *WidgetParent) AddChild(child ...Widget) {
	if child == nil {
		log.E("nuxui", "child should not be nil when add child to parent")
		return
	}

	for _, c := range child {
		MountWidget(c, me.owner)
	}

	me.children = append(me.children, child...)
	// RequestLayout(me.owner)
	// RequestRedraw(me.owner)
}

func (me *WidgetParent) RemoveChildAt(index int) {
	if index < 0 || index >= len(me.children) {
		log.Fatal("nuxui", "index out of range")
	} else {
		EjectChild(me.children[index])
		me.children = append(me.children[:index], me.children[index+1:]...)
	}
}

func (me *WidgetParent) RemoveChild(child Widget) int {
	index := -1
	for i, c := range me.children {
		if c == child {
			index = i
		}
	}

	if index != -1 {
		me.RemoveChildAt(index)
	} else {
		log.Fatal("nuxui", "child is not found in parent widget")
	}

	return index
}

func (me *WidgetParent) ReplaceChild(src, dest Widget) int {
	index := -1
	for i, child := range me.children {
		if child == src {
			index = i
		}
	}

	if index != -1 {
		EjectChild(me.children[index])
		MountWidget(dest, me.owner)
		me.children[index] = dest
	} else {
		log.Fatal("nuxui", "child is not found in parent")
	}

	return index
}

func (me *WidgetParent) GetWidget(id string) Widget {
	// TODO:: need check me.children is nil?
	for _, child := range me.children {
		if child.Info().ID == id {
			return child
		}
	}

	log.Fatal("nuxui", "The id '%s' was not found\n", id)
	return nil
}

func (me *WidgetParent) Children() []Widget {
	return me.children
}

func (me *WidgetParent) ChildrenCount() int {
	return len(me.children)
}

func (me *WidgetParent) ChildAt(index int) Widget {
	return nil
}

func (me *WidgetParent) Find(id string) Widget {
	return nil
}
