// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Component interface {
	Widget
	SetContent(Widget)
	Content() Widget
}

func NewComponentBase(owner Widget, attr Attr) *ComponentBase {
	return &ComponentBase{
		info: &WidgetInfo{
			ID: attr.GetString("id", ""),
		},
	}
}

type ComponentBase struct {
	info    *WidgetInfo
	content Widget
}

func (me *ComponentBase) Info() *WidgetInfo {
	return me.info
}

func (me *ComponentBase) SetContent(content Widget) {
	me.content = content
	// content.AssignParent(me)
}

func (me *ComponentBase) Content() Widget {
	return me.content
}
