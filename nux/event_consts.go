// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type EventAction int
type EventType int
type Kind int
type MouseButton int

const (
	Type_None EventType = iota
	Type_WindowEvent
	Type_TypingEvent
	Type_AppExit
	Type_KeyEvent
	Type_PointerEvent
	Type_BackToUI
)

const (
	Action_None EventAction = iota
	Action_Down
	Action_Move
	Action_Up
	Action_Hover
	Action_Scroll
	Action_Pressure

	Action_WindowCreating
	Action_WindowCreated
	Action_WindowMeasured
	Action_WindowDraw
	Action_WindowFocusGained
	Action_WindowFocusLost
	Action_WindowActived
)

const (
	Kind_None Kind = 0 + iota
	Kind_Mouse
	Kind_Touch
	Kind_Tablet
	Kind_Other
)

const (
	MB_None MouseButton = iota
	MB_Left
	MB_Middle
	MB_Right
	MB_X1
	MB_X2
	MB_Other
)

func (me EventType) String() string {
	switch me {
	case Type_None:
		return "Type_None "
	case Type_WindowEvent:
		return "Type_WindowEvent"
	case Type_TypingEvent:
		return "Type_TypingEvent"
	case Type_AppExit:
		return "Type_AppExit"
	case Type_KeyEvent:
		return "Type_KeyEvent"
	case Type_PointerEvent:
		return "Type_PointerEvent"
	}
	log.Fatal("nuxui", "EventType %d not handled in switch", me)
	return ""
}

func (me EventAction) String() string {
	switch me {
	case Action_None:
		return "Action_None"
	case Action_Down:
		return "Action_Down"
	case Action_Move:
		return "Action_Move"
	case Action_Up:
		return "Action_Up"
	case Action_Hover:
		return "Action_Hover"
	case Action_Scroll:
		return "Action_Scroll"
	case Action_Pressure:
		return "Action_Pressure"
	case Action_WindowCreating:
		return "Action_WindowCreating"
	case Action_WindowCreated:
		return "Action_WindowCreated"
	case Action_WindowMeasured:
		return "Action_WindowMeasured"
	case Action_WindowDraw:
		return "Action_WindowDraw"
	case Action_WindowFocusGained:
		return "Action_WindowFocusGained"
	case Action_WindowFocusLost:
		return "Action_WindowFocusLost"
	case Action_WindowActived:
		return "Action_WindowActived"
	}
	log.Fatal("nuxui", "EventAction %d not handled in switch", me)
	return ""
}
