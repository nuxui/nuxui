// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
)

// TODO
// child.BringToFront(){
// child.Parent().BringChildToFront(child)
// }

type Parent interface {
	Widget
	InsertChild(index int, child Widget)
	AddChild(child Widget) //TODO:: AddChild(child Widget...)
	Remove(index int)
	RemoveChild(child Widget) (index int)
	ReplaceChild(src, dest Widget) (index int)
	Children() []Widget // the order of Children is from bottom to top
	ChildrenCount() int
	ChildAt(index int) Widget
	Find(id string) Widget
}

type WidgetParent struct {
	Owner Parent
	WidgetBase
	children []Widget
}

func (me *WidgetParent) Creating(attr Attr) {
	if me.Owner == nil {
		log.Fatal("nuxui", "set Owner to WidgetParent before use")
	}
}

func (me *WidgetParent) InsertChild(index int, child Widget) {
	if index < 0 && index > len(me.children) {
		log.Fatal("nuxui", "index out of range when insert child to parent")
	} else {
		c := make([]Widget, 0, len(me.children)+1)
		c = append(c, me.children[:index]...)
		c = append(c, child)
		c = append(c, me.children[index:]...)
		child.AssignParent(me.Owner)
		me.children = c
	}
}

func (me *WidgetParent) AddChild(child Widget) {
	if child == nil {
		log.E("nuxui", "child should not be nil when add child to parent")
		return
	}

	if me.children == nil {
		me.children = make([]Widget, 0, 10)
	}

	child.AssignParent(me.Owner)
	me.children = append(me.children, child)
	// RequestLayout(me.Owner)
	// RequestRedraw(me.Owner)
}

func (me *WidgetParent) Remove(index int) {
	if index < 0 && index >= len(me.children) {
		log.Fatal("nuxui", "out of range")
	} else {
		me.children[index].AssignParent(nil)
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
		me.Remove(index)
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
		dest.AssignParent(me.Owner)
		me.children[index].AssignParent(nil)
		me.children[index] = dest
	} else {
		log.Fatal("nuxui", "child is not found in parent")
	}

	return index
}

func (me *WidgetParent) GetWidget(id string) Widget {
	// TODO:: need check me.children is nil?
	for _, child := range me.children {
		if child.ID() == id {
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
