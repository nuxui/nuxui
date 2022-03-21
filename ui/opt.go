// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

type Radio Opt

func NewRadio(attr nux.Attr) Radio {
	a := nux.Attr{
		"theme": "radio",
	}
	nux.MergeAttrs(a, attr)
	return Radio(NewOpt(a))
}

type Opt interface {
	Label
	Optable
}

type opt struct {
	*label
	opted    bool
	listener OptChangedCallback
}

func NewOpt(attr nux.Attr) Opt {
	myattr := nux.Attr{
		"selectable": true,
		"clickable":  true,
	}
	nux.MergeAttrs(myattr, attr)
	me := &opt{
		opted: attr.GetBool("opted", false),
		label: NewLabel(myattr).(*label),
	}

	return me
}

func (me *opt) Mount() {
	log.I("nuxui", "opt Mount")
	me.label.Mount()
	nux.OnTap(me.Info().Self, me.onTap)
}

func (me *opt) onTap(detail nux.GestureDetail) {
	log.V("nuxui", "opt onTap")
	if !me.Disable() {
		me.opted = !me.opted
		me.doOptedChanged(true)
	}
}

func (me *opt) doOptedChanged(fromUser bool) {
	changed := false

	if me.Background() != nil {
		if me.opted {
			me.Background().AddState(nux.State_Opted)
		} else {
			me.Background().DelState(nux.State_Opted)
		}
		changed = true
	}
	if me.Foreground() != nil {
		if me.opted {
			me.Background().AddState(nux.State_Opted)
		} else {
			me.Background().DelState(nux.State_Opted)
		}
		changed = true
	}
	if me.iconLeft != nil {
		if s, ok := me.iconLeft.(nux.Stateable); ok {
			s.DelState(nux.State_Pressed)
			if me.opted {
				s.AddState(nux.State_Opted)
			} else {
				s.DelState(nux.State_Opted)
			}
			changed = true
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}

	if me.listener != nil {
		me.listener(me, me.opted, fromUser)
	}
}

func (me *opt) SetOpted(opted bool, fromUser bool) {
	if me.opted != opted {
		me.opted = opted
		me.doOptedChanged(fromUser)
	}
}

func (me *opt) Opted() bool {
	return me.opted
}

func (me *opt) SetOptChangedCallback(listener OptChangedCallback) {
	me.listener = listener
}
