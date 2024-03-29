// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "nuxui.org/nuxui/log"

type EventAction int
type EventType int
type Kind int
type PointerButton int

const (
	Type_None EventType = iota
	Type_WindowEvent
	Type_TypingEvent
	Type_AppExit
	Type_KeyEvent
	Type_PointerEvent
	Type_ScrollEvent
	Type_BackToUI
)

const (
	Action_None EventAction = iota
	Action_Down
	// Action_Move
	Action_Drag
	Action_Up
	Action_Hover
	Action_Scroll
	Action_Pressure

	// typing
	Action_Input
	Action_Preedit

	Action_WindowCreated
	Action_WindowMeasured
	Action_WindowDraw
	Action_WindowFocusGained
	Action_WindowFocusLost
	Action_WindowActived
	Action_WindowHide
)

const (
	Kind_None           Kind = 0 + iota //
	Kind_Mouse                          // A mouse-based pointer device.
	Kind_Touch                          // A touch-based pointer device.
	Kind_Stylus                         // Kind_Stylus A pointer device with a stylus that has been inverted.
	Kind_InvertedStylus                 // A pointer device with a stylus that has been inverted.
	Kind_Tablet                         // An tablet pointer device, same with mouse?
	Kind_Other                          // An other unknown pointer device.
)

const (
	ButtonNone PointerButton = iota
	ButtonPrimary
	ButtonSecondary
	ButtonMiddle
	ButtonX1
	ButtonX2
	ButtonOther
)

func (me Kind) String() string {
	switch me {
	case Kind_None:
		return "Kind_None"
	case Kind_Mouse:
		return "Kind_Mouse"
	case Kind_Touch:
		return "Kind_Touch"
	case Kind_Stylus:
		return "Kind_Stylus"
	case Kind_InvertedStylus:
		return "Kind_InvertedStylus"
	case Kind_Tablet:
		return "Kind_Tablet"
	case Kind_Other:
		return "Kind_Other"
	}
	log.Fatal("nuxui", "Kind %d not handled in switch", me)
	return ""
}

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
	// case Action_Move:
	// 	return "Action_Move"
	case Action_Drag:
		return "Action_Drag"
	case Action_Up:
		return "Action_Up"
	case Action_Hover:
		return "Action_Hover"
	case Action_Scroll:
		return "Action_Scroll"
	case Action_Input:
		return "Action_Input"
	case Action_Preedit:
		return "Action_Preedit"
	case Action_Pressure:
		return "Action_Pressure"
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
	case Action_WindowHide:
		return "Action_WindowHide"
	}
	log.Fatal("nuxui", "EventAction %d not handled in switch", me)
	return ""
}

func (me PointerButton) String() string {
	switch me {
	case ButtonPrimary:
		return "ButtonPrimary"
	case ButtonSecondary:
		return "ButtonSecondary"
	case ButtonMiddle:
		return "ButtonMiddle"
	case ButtonX1:
		return "ButtonX1"
	case ButtonX2:
		return "ButtonX2"
	case ButtonOther:
		return "ButtonOther"
	}
	log.Fatal("nuxui", "unknow PointerButton %d", me)
	return ""
}
