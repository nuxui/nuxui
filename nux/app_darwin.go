// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package nux

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Cocoa -framework OpenGL
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>
#include <pthread.h>
#include <stdlib.h>

void runApp(void);
*/
import "C"
import (
	"fmt"
	"time"

	"github.com/nuxui/nuxui/log"
)

var theApp = &application{
	event:              make(chan Event),
	eventWaitDone:      make(chan Event),
	eventDone:          make(chan struct{}),
	nativeLoopPrepared: make(chan struct{}),
}

func app() Application {
	return theApp
}

type application struct {
	manifest           Manifest
	window             *window
	event              chan Event
	eventWaitDone      chan Event
	eventDone          chan struct{}
	nativeLoopPrepared chan struct{}
}

func (me *application) Creating(attr Attr) {
	if me.manifest == nil {
		me.manifest = NewManifest()
	}

	if c, ok := me.manifest.(Creating); ok {
		c.Creating(attr.GetAttr("manifest", Attr{}))
	}

	if me.window == nil {
		me.window = newWindow()
	}

	me.window.Creating(attr)
}

func (me *application) Created(data interface{}) {
	if c, ok := me.manifest.(AnyCreated); ok {
		c.Created(data)
	}
}

func (me *application) MainWindow() Window {
	return me.window
}

func (me *application) Manifest() Manifest {
	return me.manifest
}

func (me *application) SendEvent(event Event) {
	me.event <- event
}

func (me *application) sendEventAndWaitDone(event Event) {
	me.eventWaitDone <- event
	<-me.eventDone
}

//export windowCreated
func windowCreated(windptr C.uintptr_t) {
	theApp.window.windptr = windptr

	e := &windowEvent{
		event: event{
			id:     1,
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: Action_WindowCreated,
		},
		window: theApp.findWindow(windptr),
	}

	theApp.sendEventAndWaitDone(e)
}

//export windowResized
func windowResized(windptr C.uintptr_t) {
	e := &windowEvent{
		event: event{
			id:     1,
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: Action_WindowMeasured,
		},
		window: theApp.findWindow(windptr),
	}

	theApp.sendEventAndWaitDone(e)
}

//export windowDraw
func windowDraw(windptr C.uintptr_t) {
	log.V("nux", "windowDraw")

	e := &windowEvent{
		event: event{
			id:     1,
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: Action_WindowDraw,
		},
		window: theApp.findWindow(windptr),
	}

	theApp.sendEventAndWaitDone(e)
}

//export onMouseDown
func onMouseDown(windptr C.uintptr_t, x C.float, y C.float) {
	log.V("nux", "go go go xxxxx onMouseDown x: %f, y: %f", float32(x), float32(y))
}

func (me *application) findWindow(windptr C.uintptr_t) Window {
	if me.window.windptr == windptr {
		return me.window
	}
	return nil
}

func run() {
	go func() {
		<-theApp.nativeLoopPrepared
		theApp.loop()
	}()

	C.runApp()
}

//export nativeLoopPrepared
func nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

func startTextInput() {
}

//export eventKey
func eventKey(runeVal int32, direction uint8, code uint16, flags uint32) {
	// var modifiers key.Modifiers
	// for _, mod := range mods {
	// 	if flags&mod.flags == mod.flags {
	// 		modifiers |= mod.mod
	// 	}
	// }

	// theApp.eventsIn <- key.Event{
	// 	Rune:      convRune(rune(runeVal)),
	// 	Code:      convVirtualKeyCode(code),
	// 	Modifiers: modifiers,
	// 	Direction: key.Direction(direction),
	// }
	convVirtualKeyCode(code)
}

// convVirtualKeyCode converts a Carbon/Cocoa virtual key code number
// into the standard keycodes used by the key package.
//
// To get a sense of the key map, see the diagram on
//	http://boredzo.org/blog/archives/2007-05-22/virtual-key-codes
func convVirtualKeyCode(vkcode uint16) KeyCode {
	log.V("keycode", "convVirtualKeyCaode %d", vkcode)
	switch vkcode {
	case C.kVK_ANSI_A:
		log.V("keycode", fmt.Sprintf("Key_A = %d", vkcode))
		return Key_A
	case C.kVK_ANSI_B:
		log.V("keycode", fmt.Sprintf("Key_B = %d", vkcode))
		return Key_B
	case C.kVK_ANSI_C:
		log.V("keycode", fmt.Sprintf("Key_C = %d", vkcode))
		return Key_C
	case C.kVK_ANSI_D:
		log.V("keycode", fmt.Sprintf("Key_D = %d", vkcode))
		return Key_D
	case C.kVK_ANSI_E:
		log.V("keycode", fmt.Sprintf("Key_E = %d", vkcode))
		return Key_E
	case C.kVK_ANSI_F:
		log.V("keycode", fmt.Sprintf("Key_F = %d", vkcode))
		return Key_F
	case C.kVK_ANSI_G:
		log.V("keycode", fmt.Sprintf("Key_G = %d", vkcode))
		return Key_G
	case C.kVK_ANSI_H:
		log.V("keycode", fmt.Sprintf("Key_H = %d", vkcode))
		return Key_H
	case C.kVK_ANSI_I:
		log.V("keycode", fmt.Sprintf("Key_I = %d", vkcode))
		return Key_I
	case C.kVK_ANSI_J:
		log.V("keycode", fmt.Sprintf("Key_J = %d", vkcode))
		return Key_J
	case C.kVK_ANSI_K:
		log.V("keycode", fmt.Sprintf("Key_K = %d", vkcode))
		return Key_K
	case C.kVK_ANSI_L:
		log.V("keycode", fmt.Sprintf("Key_L = %d", vkcode))
		return Key_L
	case C.kVK_ANSI_M:
		log.V("keycode", fmt.Sprintf("Key_M = %d", vkcode))
		return Key_M
	case C.kVK_ANSI_N:
		log.V("keycode", fmt.Sprintf("Key_N = %d", vkcode))
		return Key_N
	case C.kVK_ANSI_O:
		log.V("keycode", fmt.Sprintf("Key_O = %d", vkcode))
		return Key_O
	case C.kVK_ANSI_P:
		log.V("keycode", fmt.Sprintf("Key_P = %d", vkcode))
		return Key_P
	case C.kVK_ANSI_Q:
		log.V("keycode", fmt.Sprintf("Key_Q = %d", vkcode))
		return Key_Q
	case C.kVK_ANSI_R:
		log.V("keycode", fmt.Sprintf("Key_R = %d", vkcode))
		return Key_R
	case C.kVK_ANSI_S:
		log.V("keycode", fmt.Sprintf("Key_S = %d", vkcode))
		return Key_S
	case C.kVK_ANSI_T:
		log.V("keycode", fmt.Sprintf("Key_T = %d", vkcode))
		return Key_T
	case C.kVK_ANSI_U:
		log.V("keycode", fmt.Sprintf("Key_U = %d", vkcode))
		return Key_U
	case C.kVK_ANSI_V:
		log.V("keycode", fmt.Sprintf("Key_V = %d", vkcode))
		return Key_V
	case C.kVK_ANSI_W:
		log.V("keycode", fmt.Sprintf("Key_W = %d", vkcode))
		return Key_W
	case C.kVK_ANSI_X:
		log.V("keycode", fmt.Sprintf("Key_X = %d", vkcode))
		return Key_X
	case C.kVK_ANSI_Y:
		log.V("keycode", fmt.Sprintf("Key_Y = %d", vkcode))
		return Key_Y
	case C.kVK_ANSI_Z:
		log.V("keycode", fmt.Sprintf("Key_Z = %d", vkcode))
		return Key_Z
	case C.kVK_ANSI_1:
		log.V("keycode", fmt.Sprintf("Key_1 = %d", vkcode))
		return Key_1
	case C.kVK_ANSI_2:
		log.V("keycode", fmt.Sprintf("Key_2 = %d", vkcode))
		return Key_2
	case C.kVK_ANSI_3:
		log.V("keycode", fmt.Sprintf("Key_3 = %d", vkcode))
		return Key_3
	case C.kVK_ANSI_4:
		log.V("keycode", fmt.Sprintf("Key_4 = %d", vkcode))
		return Key_4
	case C.kVK_ANSI_5:
		log.V("keycode", fmt.Sprintf("Key_5 = %d", vkcode))
		return Key_5
	case C.kVK_ANSI_6:
		log.V("keycode", fmt.Sprintf("Key_6 = %d", vkcode))
		return Key_6
	case C.kVK_ANSI_7:
		log.V("keycode", fmt.Sprintf("Key_7 = %d", vkcode))
		return Key_7
	case C.kVK_ANSI_8:
		log.V("keycode", fmt.Sprintf("Key_8 = %d", vkcode))
		return Key_8
	case C.kVK_ANSI_9:
		log.V("keycode", fmt.Sprintf("Key_9 = %d", vkcode))
		return Key_9
	case C.kVK_ANSI_0:
		log.V("keycode", fmt.Sprintf("Key_0 = %d", vkcode))
		return Key_0
	case C.kVK_F1:
		log.V("keycode", fmt.Sprintf("Key_F1 = %d", vkcode))
		return Key_F1
	case C.kVK_F2:
		log.V("keycode", fmt.Sprintf("Key_F2 = %d", vkcode))
		return Key_F2
	case C.kVK_F3:
		log.V("keycode", fmt.Sprintf("Key_F3 = %d", vkcode))
		return Key_F3
	case C.kVK_F4:
		log.V("keycode", fmt.Sprintf("Key_F4 = %d", vkcode))
		return Key_F4
	case C.kVK_F5:
		log.V("keycode", fmt.Sprintf("Key_F5 = %d", vkcode))
		return Key_F5
	case C.kVK_F6:
		log.V("keycode", fmt.Sprintf("Key_F6 = %d", vkcode))
		return Key_F6
	case C.kVK_F7:
		log.V("keycode", fmt.Sprintf("Key_F7 = %d", vkcode))
		return Key_F7
	case C.kVK_F8:
		log.V("keycode", fmt.Sprintf("Key_F8 = %d", vkcode))
		return Key_F8
	case C.kVK_F9:
		log.V("keycode", fmt.Sprintf("Key_F9 = %d", vkcode))
		return Key_F9
	case C.kVK_F10:
		log.V("keycode", fmt.Sprintf("Key_F10 = %d", vkcode))
		return Key_F10
	case C.kVK_F11:
		log.V("keycode", fmt.Sprintf("Key_F11 = %d", vkcode))
		return Key_F11
	case C.kVK_F12:
		log.V("keycode", fmt.Sprintf("Key_F12 = %d", vkcode))
		return Key_F12
	case C.kVK_F13:
		log.V("keycode", fmt.Sprintf("Key_F13 = %d", vkcode))
		return Key_F13
	case C.kVK_F14:
		log.V("keycode", fmt.Sprintf("Key_F14 = %d", vkcode))
		return Key_F14
	case C.kVK_F15:
		log.V("keycode", fmt.Sprintf("Key_F15 = %d", vkcode))
		return Key_F15
	case C.kVK_F16:
		log.V("keycode", fmt.Sprintf("Key_F16 = %d", vkcode))
		return Key_F16
	case C.kVK_F17:
		log.V("keycode", fmt.Sprintf("Key_F17 = %d", vkcode))
		return Key_F17
	case C.kVK_F18:
		log.V("keycode", fmt.Sprintf("Key_F18 = %d", vkcode))
		return Key_F18
	case C.kVK_F19:
		log.V("keycode", fmt.Sprintf("Key_F19 = %d", vkcode))
		return Key_F19
	case C.kVK_F20:
		log.V("keycode", fmt.Sprintf("Key_F20 = %d", vkcode))
		return Key_F20

	case C.kVK_Return:
		log.V("keycode", fmt.Sprintf("Key_Return = %d", vkcode))
		return Key_Return
	case C.kVK_Tab:
		log.V("keycode", fmt.Sprintf("Key_Tab = %d", vkcode))
		return Key_Tab
	case C.kVK_Space:
		log.V("keycode", fmt.Sprintf("Key_Space = %d", vkcode))
		return Key_Space
	case C.kVK_Delete:
		log.V("keycode", fmt.Sprintf("Key_Delete = %d", vkcode))
		return Key_Delete
	case C.kVK_Escape:
		log.V("keycode", fmt.Sprintf("Key_Escape = %d", vkcode))
		return Key_Escape
	case C.kVK_Command:
		log.V("keycode", fmt.Sprintf("Key_Command = %d", vkcode))
		return Key_Command
	case C.kVK_CapsLock:
		log.V("keycode", fmt.Sprintf("Key_CapsLock = %d", vkcode))
		return Key_CapsLock
	case C.kVK_Option:
		log.V("keycode", fmt.Sprintf("Key_LAlt = %d", vkcode))
		return Key_LAlt
	case C.kVK_RightOption:
		log.V("keycode", fmt.Sprintf("Key_RAlt = %d", vkcode))
		return Key_RAlt
	case C.kVK_Shift:
		log.V("keycode", fmt.Sprintf("Key_LShift = %d", vkcode))
		return Key_LShift
	case C.kVK_RightShift:
		log.V("keycode", fmt.Sprintf("Key_RShift = %d", vkcode))
		return Key_RShift
	case C.kVK_Control:
		log.V("keycode", fmt.Sprintf("Key_LControl = %d", vkcode))
		return Key_LControl
	case C.kVK_RightControl:
		log.V("keycode", fmt.Sprintf("Key_RControl = %d", vkcode))
		return Key_RControl
	case C.kVK_ANSI_Equal:
		log.V("keycode", fmt.Sprintf("Key_Equal = %d", vkcode))
		return Key_Equal
	case C.kVK_ANSI_Minus:
		log.V("keycode", fmt.Sprintf("Key_Minus = %d", vkcode))
		return Key_Minus
	case C.kVK_ANSI_LeftBracket:
		log.V("keycode", fmt.Sprintf("Key_LBracket = %d", vkcode))
		return Key_LBracket
	case C.kVK_ANSI_RightBracket:
		log.V("keycode", fmt.Sprintf("Key_RBracket = %d", vkcode))
		return Key_RBracket
	case C.kVK_ANSI_Quote:
		log.V("keycode", fmt.Sprintf("Key_Quote = %d", vkcode))
		return Key_Quote
	case C.kVK_ANSI_Semicolon:
		log.V("keycode", fmt.Sprintf("Key_Semicolon = %d", vkcode))
		return Key_Semicolon
	case C.kVK_ANSI_Comma:
		log.V("keycode", fmt.Sprintf("Key_Comma = %d", vkcode))
		return Key_Comma
	case C.kVK_ANSI_Slash:
		log.V("keycode", fmt.Sprintf("Key_Slash = %d", vkcode))
		return Key_Slash
	case C.kVK_ANSI_Backslash:
		log.V("keycode", fmt.Sprintf("Key_Backslash = %d", vkcode))
		return Key_Backslash
	case C.kVK_ANSI_Period:
		log.V("keycode", fmt.Sprintf("Key_Period = %d", vkcode))
		return Key_Period
	case C.kVK_ANSI_Grave:
		log.V("keycode", fmt.Sprintf("Key_Grave = %d", vkcode))
		return Key_Grave

	case C.kVK_ANSI_Keypad1:
		log.V("keycode", fmt.Sprintf("Key_Pad1 = %d", vkcode))
		return Key_Pad1
	case C.kVK_ANSI_Keypad2:
		log.V("keycode", fmt.Sprintf("Key_Pad2 = %d", vkcode))
		return Key_Pad2
	case C.kVK_ANSI_Keypad3:
		log.V("keycode", fmt.Sprintf("Key_Pad3 = %d", vkcode))
		return Key_Pad3
	case C.kVK_ANSI_Keypad4:
		log.V("keycode", fmt.Sprintf("Key_Pad4 = %d", vkcode))
		return Key_Pad4
	case C.kVK_ANSI_Keypad5:
		log.V("keycode", fmt.Sprintf("Key_Pad5 = %d", vkcode))
		return Key_Pad5
	case C.kVK_ANSI_Keypad6:
		log.V("keycode", fmt.Sprintf("Key_Pad6 = %d", vkcode))
		return Key_Pad6
	case C.kVK_ANSI_Keypad7:
		log.V("keycode", fmt.Sprintf("Key_Pad7 = %d", vkcode))
		return Key_Pad7
	case C.kVK_ANSI_Keypad8:
		log.V("keycode", fmt.Sprintf("Key_Pad8 = %d", vkcode))
		return Key_Pad8
	case C.kVK_ANSI_Keypad9:
		log.V("keycode", fmt.Sprintf("Key_Pad9 = %d", vkcode))
		return Key_Pad9
	case C.kVK_ANSI_Keypad0:
		log.V("keycode", fmt.Sprintf("Key_Pad0 = %d", vkcode))
		return Key_Pad0
	case C.kVK_ANSI_KeypadPlus:
		log.V("keycode", fmt.Sprintf("Key_PadAdd = %d", vkcode))
		return Key_PadAdd
	case C.kVK_ANSI_KeypadMinus:
		log.V("keycode", fmt.Sprintf("Key_PadSubtract = %d", vkcode))
		return Key_PadSubtract
	case C.kVK_ANSI_KeypadMultiply:
		log.V("keycode", fmt.Sprintf("Key_PadMultiply = %d", vkcode))
		return Key_PadMultiply
	case C.kVK_ANSI_KeypadDivide:
		log.V("keycode", fmt.Sprintf("Key_PadDivide = %d", vkcode))
		return Key_PadDivide
	case C.kVK_ANSI_KeypadEquals:
		log.V("keycode", fmt.Sprintf("Key_PadEquals = %d", vkcode))
		return Key_PadEquals
	case C.kVK_ANSI_KeypadClear: // num lock
		log.V("keycode", fmt.Sprintf("Key_PadClear = %d", vkcode))
		return Key_PadClear
	case C.kVK_ANSI_KeypadDecimal:
		log.V("keycode", fmt.Sprintf("Key_PadDecimal = %d", vkcode))
		return Key_PadDecimal
	case C.kVK_ANSI_KeypadEnter:
		log.V("keycode", fmt.Sprintf("Key_PadEnter = %d", vkcode))
		return Key_PadEnter
	case C.kVK_UpArrow:
		log.V("keycode", fmt.Sprintf("Key_UpArrow = %d", vkcode))
		return Key_UpArrow
	case C.kVK_DownArrow:
		log.V("keycode", fmt.Sprintf("Key_DownArrow = %d", vkcode))
		return Key_DownArrow
	case C.kVK_LeftArrow:
		log.V("keycode", fmt.Sprintf("Key_LeftArrow = %d", vkcode))
		return Key_LeftArrow
	case C.kVK_RightArrow:
		log.V("keycode", fmt.Sprintf("Key_RightArrow = %d", vkcode))
		return Key_RightArrow
	case C.kVK_PageUp:
		log.V("keycode", fmt.Sprintf("Key_PageUp = %d", vkcode))
		return Key_PageUp
	case C.kVK_PageDown:
		log.V("keycode", fmt.Sprintf("Key_PageDown = %d", vkcode))
		return Key_PageDown
	case C.kVK_Home:
		log.V("keycode", fmt.Sprintf("Key_Home = %d", vkcode))
		return Key_Home
	case C.kVK_End:
		log.V("keycode", fmt.Sprintf("Key_End = %d", vkcode))
		return Key_End

	case C.kVK_Help:
		log.V("keycode", fmt.Sprintf("Key_Help = %d", vkcode))
		return Key_Help
	case C.kVK_ForwardDelete:
		log.V("keycode", fmt.Sprintf("Key_ForwardDelete = %d", vkcode))
		return Key_ForwardDelete

	case C.kVK_Mute:
		log.V("keycode", fmt.Sprintf("Key_Mute = %d", vkcode))
		return Key_Mute
	case C.kVK_VolumeUp:
		log.V("keycode", fmt.Sprintf("Key_VolumeUp = %d", vkcode))
		return Key_VolumeUp
	case C.kVK_VolumeDown:
		log.V("keycode", fmt.Sprintf("Key_VolumeDown = %d", vkcode))
		return Key_VolumeDown

	// 116: Keyboard Execute
	// 118: Keyboard Menu
	// 119: Keyboard Select
	// 120: Keyboard Stop
	// 121: Keyboard Again
	// 122: Keyboard Undo
	// 123: Keyboard Cut
	// 124: Keyboard Copy
	// 125: Keyboard Paste
	// 126: Keyboard Find
	// 130: Keyboard Locking Caps Lock
	// 131: Keyboard Locking Num Lock
	// 132: Keyboard Locking Scroll Lock
	// 133: Keyboard Comma
	// 134: Keyboard Equal Sign
	// ...: Bunch of stuff
	default:
		log.V("keycode", fmt.Sprintf("Key_Unknown = %d", vkcode))
		return Key_Unknown
	}
}
