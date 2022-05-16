// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/nux"
)

type Options interface {
	nux.Component
}

type options struct {
	*nux.ComponentBase

	radio bool
}

func NewOptions(attr nux.Attr) Options {
	me := &options{
		radio: attr.GetBool("radio", false),
	}
	me.ComponentBase = nux.NewComponentBase(me, attr)
	return me
}

func (me *options) Mount() {
	me.bindchild()
}

func (me *options) bindchild() {
	if p, ok := me.Content().(nux.Parent); ok {
		for _, child := range p.Children() {
			if item, ok := child.(Checkable); ok {
				item.SetCheckChangedCallback(me.onCheckedChanged)
			}
		}
	}
}

func (me *options) onCheckedChanged(widget CheckableWidget, checked bool, fromUser bool) {
	if me.radio && fromUser {
		if p, ok := me.Content().(nux.Parent); ok {
			for _, child := range p.Children() {
				if item, ok := child.(Checkable); ok {
					if item != widget && item.Checked() {
						item.SetChecked(false, false)
					}
				}
			}
		}
	}
}

func (me *options) doCheckedChanged() {

}
