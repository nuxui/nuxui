// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

type Stateable interface {
	AddState(state uint32) (changed bool)
	DelState(state uint32) (changed bool)
	State() uint32
	HasState() bool
}

func NewStateBase(stateChanged func()) *StateBase {
	return &StateBase{
		state:        State_Default,
		stateChanged: stateChanged,
	}
}

type StateBase struct {
	state        uint32
	stateChanged func()
}

func (me *StateBase) AddState(state uint32) bool {
	s := me.state
	s |= state
	if me.state != s {
		me.state = s
		if me.stateChanged != nil {
			me.stateChanged()
		}
		return true
	}
	return false
}

func (me *StateBase) DelState(state uint32) bool {
	s := me.state
	s &= ^state
	if me.state != s {
		me.state = s
		if me.stateChanged != nil {
			me.stateChanged()
		}
		return true
	}
	return false
}

func (me *StateBase) State() uint32 {
	return me.state
}

func (me *StateBase) HasState() bool {
	return false
}

const (
	State_Disabled uint32 = 1 << iota
	State_Pressed
	State_Focused
	State_Hovered
	State_Checked
	State_Visited        // state of visited links
	State_Invalid        // invalid of input
	State_Default uint32 = 0
)

func StateFromString(state string) uint32 {
	switch state {
	case "disabled":
		return State_Disabled
	case "pressed":
		return State_Pressed
	case "focused":
		return State_Focused
	case "hovered":
		return State_Hovered
	case "checked":
		return State_Checked
	case "visited":
		return State_Visited
	case "invalid":
		return State_Invalid
	case "default":
		return State_Default
	}
	log.Fatal("nuxui", "unknow state %s", state)
	return State_Default
}

func StateToString(state uint32) string {
	switch state {
	case State_Disabled:
		return "Disabled"
	case State_Pressed:
		return "Pressed"
	case State_Focused:
		return "Focused"
	case State_Hovered:
		return "Hovered"
	case State_Checked:
		return "Checked"
	case State_Visited:
		return "Visited"
	case State_Invalid:
		return "Invalid"
	case State_Default:
		return "Default"
	}
	log.Fatal("nuxui", "unknow state %d", state)
	return ""
}
