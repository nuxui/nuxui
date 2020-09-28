// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"math"
	"time"

	"github.com/nuxui/nuxui/log"
)

type Event interface {
	Time() time.Time
	Type() EventType
	Action() EventAction
	Window() Window //TODO:: if global save a event, then window will not release

	// pointer event
	Pointer() int64
	Kind() Kind
	X() float32
	Y() float32
	ScreenX() float32
	ScreenY() float32
	ScrollX() float32
	ScrollY() float32
	Pressure() float32
	Stage() int32
	IsPrimary() bool
	Distance(x, y float32) float32

	// key event
	KeyCode() KeyCode
	Repeat() bool
	Modifiers() (none, capslock, shift, control, alt, super bool)
	Rune() rune // rune or 0
	Text() string

	Data() interface{}
}

type event struct {
	time   time.Time
	etype  EventType
	action EventAction
	window Window

	pointer  int64
	kind     Kind
	x        float32
	y        float32
	screenX  float32
	screenY  float32
	scrollX  float32
	scrollY  float32
	pressure float32
	stage    int32
	button   MouseButton

	keyCode       KeyCode
	repeat        bool
	modifierFlags uint32
	characters    string
	text          string

	data interface{}
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

func (me *event) Pointer() int64 {
	return me.pointer
}

func (me *event) Kind() Kind {
	return me.kind
}

func (me *event) X() float32 {
	return me.x
}

func (me *event) Y() float32 {
	return me.y
}

func (me *event) ScreenX() float32 {
	return me.screenX
}

func (me *event) ScreenY() float32 {
	return me.screenY
}

func (me *event) ScrollX() float32 {
	if me.action != Action_Scroll {
		log.E("nuxui", "obtain scroll at no wheel")
	}
	return me.scrollX
}

func (me *event) ScrollY() float32 {
	if me.action != Action_Scroll {
		log.E("nuxui", "obtain scroll at no wheel")
	}
	return me.scrollY
}

func (me *event) Pressure() float32 {
	return me.pressure
}

func (me *event) Stage() int32 {
	return me.stage
}

func (me *event) IsPrimary() bool {
	if me.kind == Kind_Mouse {
		return me.button == MB_Left
	}
	return true // TODO:: multi finger
}

func (me *event) Distance(x, y float32) float32 {
	dx := me.x - x
	dy := me.y - y
	return float32(math.Sqrt(float64(dx)*float64(dx) + float64(dy)*float64(dy)))
}

func (me *event) KeyCode() KeyCode {
	return me.keyCode
}

func (me *event) Repeat() bool {
	return me.repeat
}

func (me *event) Modifiers() (none, capslock, shift, control, alt, super bool) {
	capslock = me.modifierFlags&Mod_CapsLock == Mod_CapsLock
	shift = me.modifierFlags&Mod_Shift == Mod_Shift
	control = me.modifierFlags&Mod_Control == Mod_Control
	alt = me.modifierFlags&Mod_Alt == Mod_Alt
	super = me.modifierFlags&Mod_Super == Mod_Super
	none = !(capslock || shift || control || alt || super)
	return
}

func (me *event) Rune() rune {
	r := []rune(me.characters)
	if len(r) > 0 {
		return r[0]
	}
	return 0
}

func (me *event) Text() string {
	return me.text
}

func (me *event) Data() interface{} {
	return me.data
}

func (me *event) String() string {
	switch me.etype {
	case Type_PointerEvent:
		return fmt.Sprintf("PointerEvent: {pointer=%d, button=%s:%d, action=%s, x=%.2f, y=%.2f}", me.pointer, me.button, me.button, me.action, me.x, me.y)
	case Type_KeyEvent:
		return fmt.Sprintf("KeyEvent: {keyCode=%s, action=%s, x=%.2f, y=%.2f, rune='%c'}", me.keyCode, me.action, me.x, me.y, me.Rune())
	}
	return fmt.Sprintf("Event: {type:%s, action:%s, x=%.2f, y=%.2f}", me.etype, me.action, me.x, me.y)
}

// TODO::
// func RegisterEventAction() EventAction {

// }
