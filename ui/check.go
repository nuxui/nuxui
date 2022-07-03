// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/nux"
)

type Radio Check

func NewRadio(attr nux.Attr) Radio {
	if attr == nil {
		attr = nux.Attr{}
	}
	a := nux.Attr{
		"theme": "radio",
	}
	return Radio(NewCheck(nux.MergeAttrs(a, attr)))
}

type Check interface {
	Label
	Checkable
}

type check struct {
	*label
	checked  bool
	listener CheckChangedCallback
	hasValue bool // if !hasValue, value is text
	value    string
}

func NewCheck(attr nux.Attr) Check {
	if attr == nil {
		attr = nux.Attr{}
	}
	myattr := nux.Attr{
		"selectable": true,
		"clickable":  true,
		"checkable":  true,
	}
	me := &check{
		checked:  attr.GetBool("checked", false),
		hasValue: attr.Has("value"),
		value:    attr.GetString("value", ""),
		label:    NewLabel(nux.MergeAttrs(myattr, attr)).(*label),
	}

	return me
}

func (me *check) OnMount() {
	me.label.OnMount()
	nux.OnTap(me.Info().Self, me.onTap)
}

func (me *check) onTap(detail nux.GestureDetail) {
	if !me.Disable() {
		me.checked = !me.checked
		me.doCheckedChanged(true)
	}
}

func (me *check) doCheckedChanged(fromUser bool) {
	changed := false

	if me.Background() != nil {
		if me.checked {
			me.Background().AddState(nux.State_Checked)
		} else {
			me.Background().DelState(nux.State_Checked)
		}
		changed = true
	}
	if me.Foreground() != nil {
		if me.checked {
			me.Background().AddState(nux.State_Checked)
		} else {
			me.Background().DelState(nux.State_Checked)
		}
		changed = true
	}
	if me.updateIconState(me.iconLeft, me.iconTop, me.iconRight, me.iconBottom) {
		changed = true
	}

	if changed {
		nux.RequestRedraw(me)
	}

	if me.listener != nil {
		me.listener(me, me.checked, fromUser)
	}
}

func (me *check) updateIconState(icons ...nux.Widget) (changed bool) {
	for _, icon := range icons {
		if icon != nil {
			if s, ok := icon.(nux.Stateable); ok {
				if me.checked {
					s.AddState(nux.State_Checked)
				} else {
					s.DelState(nux.State_Checked)
				}
				changed = true
			}
		}
	}
	return
}

func (me *check) SetChecked(checked bool, fromUser bool) {
	if me.checked != checked {
		me.checked = checked
		me.doCheckedChanged(fromUser)
	}
}

func (me *check) Checked() bool {
	return me.checked
}

func (me *check) SetValue(value string) {
	me.hasValue = true
	me.value = value
}

func (me *check) Value() string {
	if me.hasValue {
		return me.value
	}
	return me.Text()
}

func (me *check) SetCheckChangedCallback(listener CheckChangedCallback) {
	me.listener = listener
}
