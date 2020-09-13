// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type EventAction int
type EventType int
type Kind int
type MouseButton int

const (
	Type_None EventType = iota
	Type_WindowEvent
	Type_TypingEvent
	Type_AppExit
	Type_InputEvent
)

const (
	Action_None EventAction = iota
	Action_Down
	Action_Move
	Action_Repeat
	Action_Up

	Action_WindowCreating
	Action_WindowCreated
	Action_WindowMeasured
	Action_WindowDraw
	Action_WindowFocusGained
	Action_WindowFocusLost
	Action_WindowActived
	// Action_Window
)

const (
	Kind_None Kind = 0 + iota
	Kind_Mouse
	Kind_Touch
	Kind_Pad
)

const (
	MB_None MouseButton = iota
	MB_Left
	MB_Middle
	MB_Right
	MB_X1
	MB_X2
)
