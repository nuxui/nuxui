// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/nux"
)

type Options interface {
	nux.Component
	Values() []string
	Selected() bool
	SetOnSelectionChanged(callback func(widget Options, fromUser bool))
}

type options struct {
	*nux.ComponentBase

	single             bool
	onSelectionChanged func(widget Options, fromUser bool)
}

func NewOptions(attr nux.Attr) Options {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &options{
		single: attr.GetBool("single", false),
	}
	me.ComponentBase = nux.NewComponentBase(me, attr)
	return me
}

func (me *options) OnMount() {
	me.initChildrenCallback(me.Content())
}

func (me *options) initChildrenCallback(widget nux.Widget) {
	if p, ok := widget.(nux.Parent); ok {
		for _, child := range p.Children() {
			me.initChildrenCallback(child)
		}
	} else {
		if check, ok := widget.(Checkable); ok {
			check.SetCheckChangedCallback(me.onCheckedChanged)
		}
	}
}

func (me *options) onCheckedChanged(widget CheckableWidget, checked bool, fromUser bool) {
	if me.single && fromUser {
		me.clearAllExcept(me.Content(), widget)
	}

	if fromUser { // Options only accpet it self 'fromUser'
		me.doSelectionChanged(fromUser)
	}
}

func (me *options) clearAllExcept(widget nux.Widget, except CheckableWidget) {
	if p, ok := widget.(nux.Parent); ok {
		for _, child := range p.Children() {
			me.clearAllExcept(child, except)
		}
	} else {
		if check, ok := widget.(Checkable); ok {
			if widget != except && check.Checked() {
				check.SetChecked(false, false)
			}
		}
	}
}

func (me *options) SetOnSelectionChanged(callback func(widget Options, fromUser bool)) {
	me.onSelectionChanged = callback
}

func (me *options) doSelectionChanged(fromUser bool) {
	if me.onSelectionChanged != nil {
		me.onSelectionChanged(me, fromUser)
	}
}

func (me *options) Values() []string {
	return me.getValues(me.Content())
}

func (me *options) getValues(widget nux.Widget) (values []string) {
	if p, ok := widget.(nux.Parent); ok {
		for _, child := range p.Children() {
			values = append(values, me.getValues(child)...)
		}
	} else {
		if check, ok := widget.(Checkable); ok {
			if check.Checked() {
				values = append(values, check.Value())
			}
		}
	}
	return
}

func (me *options) Selected() bool {
	return me.hasSelected(me.Content())
}

func (me *options) hasSelected(widget nux.Widget) bool {
	if p, ok := widget.(nux.Parent); ok {
		for _, child := range p.Children() {
			if me.hasSelected(child) {
				return true
			}
		}
	} else {
		if check, ok := widget.(Checkable); ok {
			if check.Checked() {
				return true
			}
		}
	}
	return false
}
