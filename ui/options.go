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

	content := attr.GetAttr("content", nux.Attr{})
	nux.InflateLayoutAttr(me, content)
	me.bindchild()
	return me
}

func (me *options) bindchild() {
	if p, ok := me.Content().(nux.Parent); ok {
		for _, child := range p.Children() {
			if item, ok := child.(Optable); ok {
				item.SetOptChangedCallback(me.onOptedChanged)
			}
		}
	}
}

func (me *options) onOptedChanged(widget OptableWidget, isOpted bool, fromUser bool) {
	if me.radio && fromUser {
		if p, ok := me.Content().(nux.Parent); ok {
			for _, child := range p.Children() {
				if item, ok := child.(Optable); ok {
					if item != widget && item.Opted() {
						item.SetOpted(false, false)
					}
				}
			}
		}
	}
}

func (me *options) doOptedChanged() {

}
