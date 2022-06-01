// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"math"
	"time"

	"nuxui.org/nuxui/log"
)

//---------------  Readonly Event --------------------------------

type Event interface {
	Time() time.Time
	Type() EventType
	Action() EventAction
	Window() Window
	Data() any
}

type event struct {
	time   time.Time
	etype  EventType
	action EventAction
	window Window
	data   any
}

func (me *event) Time() time.Time {
	return me.time
}

func (me *event) Type() EventType {
	return me.etype
}

func (me *event) Action() EventAction {
	return me.action
}

func (me *event) Window() Window {
	return me.window
}

func (me *event) Data() any {
	return me.data
}

func (me *event) String() string {
	return fmt.Sprintf("Event: {type:%s, action:%s}", me.etype, me.action)
}

//---------------  WindowEvent --------------------------------

// WindowEvent
type WindowEvent interface {
	Event
	windowEvent()
}

type windowEvent struct {
	event
}

func (me *windowEvent) windowEvent() {}

//---------------  PointerEvent --------------------------------

type PointerEvent interface {
	Event

	Pointer() int64
	Kind() Kind
	X() float32
	Y() float32
	// Pressure() float32
	// Stage() int32
	IsPrimary() bool
	Distance(x, y float32) float32
	pointerEvent()
}

type pointerEvent struct {
	event

	pointer int64
	kind    Kind
	x       float32
	y       float32
	// pressure float32
	// stage    int32
	button MouseButton
}

func (me *pointerEvent) Pointer() int64 {
	return me.pointer
}

func (me *pointerEvent) Kind() Kind {
	return me.kind
}

func (me *pointerEvent) X() float32 {
	return me.x
}

func (me *pointerEvent) Y() float32 {
	return me.y
}

// func (me *pointerEvent) Pressure() float32 {
// 	return me.pressure
// }

// func (me *pointerEvent) Stage() int32 {
// 	return me.stage
// }

func (me *pointerEvent) IsPrimary() bool {
	if me.kind == Kind_Mouse {
		return me.button == MB_Left
	}
	return true // TODO:: multi finger
}

func (me *pointerEvent) Distance(x, y float32) float32 {
	dx := me.x - x
	dy := me.y - y
	return float32(math.Sqrt(float64(dx)*float64(dx) + float64(dy)*float64(dy)))
}

func (me *pointerEvent) pointerEvent() {}

func (me *pointerEvent) String() string {
	return fmt.Sprintf("PointerEvent: {pointer=%d, button=%s:%d, action=%s, x=%.2f, y=%.2f}", me.pointer, me.button, me.button, me.action, me.x, me.y)
}

//---------------  ScrollEvent --------------------------------

type ScrollEvent interface {
	Event

	X() float32
	Y() float32
	ScrollX() float32
	ScrollY() float32
	scrollEvent()
}

type scrollEvent struct {
	event

	x       float32
	y       float32
	scrollX float32
	scrollY float32
}

func (me *scrollEvent) X() float32 {
	return me.x
}

func (me *scrollEvent) Y() float32 {
	return me.y
}

func (me *scrollEvent) ScrollX() float32 {
	if me.action != Action_Scroll {
		log.E("nuxui", "obtain scroll at no wheel")
	}
	return me.scrollX
}

func (me *scrollEvent) ScrollY() float32 {
	if me.action != Action_Scroll {
		log.E("nuxui", "obtain scroll at no wheel")
	}
	return me.scrollY
}

func (me *scrollEvent) scrollEvent() {}

//---------------  KeyEvent --------------------------------

type KeyEvent interface {
	Event

	KeyCode() KeyCode
	Repeat() bool
	Modifiers() (none, capslock, shift, control, alt, super bool)
	Rune() rune // rune or 0
	keyEvent()
}

type keyEvent struct {
	event

	keyCode       KeyCode
	repeat        bool
	modifierFlags uint32
	keyRune       string
}

func (me *keyEvent) KeyCode() KeyCode {
	return me.keyCode
}

func (me *keyEvent) Repeat() bool {
	return me.repeat
}

func (me *keyEvent) Modifiers() (none, capslock, shift, control, alt, super bool) {
	capslock = me.modifierFlags&Mod_CapsLock == Mod_CapsLock
	shift = me.modifierFlags&Mod_Shift == Mod_Shift
	control = me.modifierFlags&Mod_Control == Mod_Control
	alt = me.modifierFlags&Mod_Alt == Mod_Alt
	super = me.modifierFlags&Mod_Super == Mod_Super
	none = !(capslock || shift || control || alt || super)
	return
}

func (me *keyEvent) Rune() rune {
	r := []rune(me.keyRune)
	if debug_event {
		if len(r) > 1 {
			log.W("nuxui", "KeyEvent original rune is %s length = %d, %d", me.keyRune, len(me.keyRune), len(r))
		}
	}
	if len(r) > 0 {
		return r[0]
	}
	return 0
}

func (me *keyEvent) keyEvent() {}

func (me *keyEvent) String() string {
	return fmt.Sprintf("KeyEvent: {keyCode=%s, action=%s, rune='%c'}", me.keyCode, me.action, me.Rune())
}

//---------------  TypingEvent --------------------------------

type TypingEvent interface {
	Event
	Text() string
	Location() int32
	typingEvent()
}

type typingEvent struct {
	event
	text     string
	location int32
}

func (me *typingEvent) Text() string {
	return me.text
}

func (me *typingEvent) Location() int32 {
	return me.location
}

func (me *typingEvent) typingEvent() {}
func (me *typingEvent) String() string {
	return fmt.Sprintf("TypingEvent: {action=%s, text=%s, location=%d}", me.action, me.text, me.location)
}
