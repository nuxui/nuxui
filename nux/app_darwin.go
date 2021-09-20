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

void runApp();
void terminate();
void startTextInput();
void stopTextInput();
void setTextInputRect(float x, float y, float w, float h);
void invalidate();
void loop_event();
void backToUI();
*/
import "C"
import (
	"runtime"
	"time"

	"github.com/nuxui/nuxui/log"
)

var theApp = &application{
	event:              make(chan Event),
	eventWaitDone:      make(chan Event),
	eventDone:          make(chan struct{}),
	runOnUI:            make(chan func()),
	nativeLoopPrepared: make(chan struct{}),
	drawSignal:         make(chan struct{}, drawSignalSize),
}

const (
	drawSignalSize = 50
)

func init() {
	go func() {
		var i, l int
		for {
			<-theApp.drawSignal
			l = len(theApp.drawSignal)
			for i = 0; i != l; i++ {
				<-theApp.drawSignal
			}
			requestRedraw()
			time.Sleep(16 * time.Millisecond)
		}
	}()
}

func app() Application {
	return theApp
}

type application struct {
	manifest      Manifest
	window        *window
	event         chan Event
	eventWaitDone chan Event
	runOnUI       chan func()

	eventDone          chan struct{}
	nativeLoopPrepared chan struct{}
	drawSignal         chan struct{}
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

func (me *application) SendEventAndWait(event Event) {
	me.eventWaitDone <- event
	<-me.eventDone
}

func (me *application) Terminate() {
	C.terminate()
}

//export windowCreated
func windowCreated(windptr C.uintptr_t) {
	theApp.window.windptr = windptr

	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: Action_WindowCreated,
		window: theApp.findWindow(windptr),
	}

	theApp.handleEvent(e)
}

//export windowResized
func windowResized(windptr C.uintptr_t) {
	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: Action_WindowMeasured,
		window: theApp.findWindow(windptr),
	}

	theApp.handleEvent(e)
}

//export windowDraw
func windowDraw(windptr C.uintptr_t) {
	log.V("nuxui", "windowDraw")

	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: Action_WindowDraw,
		window: theApp.findWindow(windptr),
	}

	theApp.handleEvent(e)
}

func (me *application) findWindow(windptr C.uintptr_t) Window {
	if me.window.windptr == windptr {
		return me.window
	}
	return nil
}

func (me *application) RequestRedraw(widget Widget) {
	if l := len(theApp.drawSignal); l >= drawSignalSize {
		for i := 0; i != l-1; i++ {
			<-theApp.drawSignal
		}
	}
	theApp.drawSignal <- struct{}{}
}

func run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	C.runApp()
}

//export go_nativeLoopPrepared
func go_nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

func startTextInput() {
	C.startTextInput()
}

func stopTextInput() {
	C.stopTextInput()
}

func setTextInputRect(x, y, w, h float32) {
	C.setTextInputRect(C.float(x), C.float(y), C.float(w), C.float(h))
}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

//export go_mouseEvent
func go_mouseEvent(windptr C.uintptr_t, etype uint, x, y, screenX, screenY, scrollX, scrollY float32, buttonNumber int32, pressure float32, stage int32) {
	// log.V("nuxui", "go_mouseEvent x=%f, y=%f, screenX=%f, screenY=%f, scrollX=%f, scrollY=%f", x, y, screenX, screenY, scrollX, scrollY)
	e := &pointerEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer:  0,
		button:   MB_None,
		kind:     Kind_Mouse,
		x:        x,
		y:        y,
		screenX:  screenX,
		screenY:  screenY,
		pressure: pressure,
		stage:    stage,
	}

	switch etype {
	case C.NSEventTypeMouseMoved:
		e.event.action = Action_Hover
		e.button = MB_None
		e.pointer = 0
	case C.NSEventTypeLeftMouseDown:
		e.event.action = Action_Down
		e.button = MB_Left
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.NSEventTypeLeftMouseUp:
		e.event.action = Action_Up
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeLeftMouseDragged:
		e.event.action = Action_Drag
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeRightMouseDown:
		e.event.action = Action_Down
		e.button = MB_Right
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.NSEventTypeRightMouseUp:
		e.event.action = Action_Up
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeRightMouseDragged:
		e.event.action = Action_Drag
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeOtherMouseDown:
		e.event.action = Action_Down
		switch buttonNumber {
		case 2:
			e.button = MB_Middle
		case 3:
			e.button = MB_X1
		case 4:
			e.button = MB_X2
		default:
			e.button = MouseButton(buttonNumber)
		}
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.NSEventTypeOtherMouseUp:
		e.event.action = Action_Up
		switch buttonNumber {
		case 2:
			e.button = MB_Middle
		case 3:
			e.button = MB_X1
		case 4:
			e.button = MB_X2
		default:
			e.button = MouseButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeOtherMouseDragged:
		e.event.action = Action_Drag
		switch buttonNumber {
		case 2:
			e.button = MB_Middle
		case 3:
			e.button = MB_X1
		case 4:
			e.button = MB_X2
		default:
			e.button = MouseButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeScrollWheel:
		log.E("nux", "use go_scrollEvent")
	case C.NSEventTypePressure:
		// TODO:: stageTransition pressureBehavior NSPressureBehavior
		e.event.action = Action_Pressure
		e.button = MB_None
	}

	theApp.handleEvent(e)

}

//export go_scrollEvent
func go_scrollEvent(windptr C.uintptr_t, etype uint, x, y, screenX, screenY, scrollX, scrollY float32, buttonNumber int32, pressure float32, stage int32) {
	log.V("nuxui", "go_scrollEvent x=%f, y=%f, screenX=%f, screenY=%f, scrollX=%f, scrollY=%f", x, y, screenX, screenY, scrollX, scrollY)
	e := &scrollEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_ScrollEvent,
			action: Action_Scroll,
		},
		x:       x,
		y:       y,
		screenX: screenX,
		screenY: screenY,
		scrollX: scrollX,
		scrollY: scrollY,
	}

	theApp.handleEvent(e)

}

var lastModifierKeyEvent map[KeyCode]bool = map[KeyCode]bool{}

//export go_drawEvent
func go_drawEvent(windptr C.uintptr_t) {
	log.V("nuxui", "go_drawEvent -----")
	e := &event{
		window: theApp.findWindow(windptr),
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: Action_WindowDraw,
	}

	theApp.handleEvent(e)
}

//export go_keyEvent
func go_keyEvent(windptr C.uintptr_t, etype uint, keyCode uint16, modifierFlags uint, repeat byte, chars *C.char) {
	e := &keyEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_KeyEvent,
			action: Action_None,
		},
		keyCode:       convertVirtualKeyCode(keyCode),
		repeat:        false,
		modifierFlags: convertModifierFlags(modifierFlags),
		keyRune:       C.GoString(chars),
	}

	if repeat == 1 {
		e.repeat = true
	}

	switch etype {
	case C.NSEventTypeKeyDown:
		e.event.action = Action_Down
	case C.NSEventTypeKeyUp:
		e.event.action = Action_Up
	case C.NSEventTypeFlagsChanged:
		if down, ok := lastModifierKeyEvent[e.keyCode]; ok && down {
			lastModifierKeyEvent[e.keyCode] = false
			e.event.action = Action_Up
		} else {
			lastModifierKeyEvent[e.keyCode] = true
			e.event.action = Action_Down
		}
	}

	theApp.handleEvent(e)
}

//export go_typingEvent
func go_typingEvent(windptr C.uintptr_t, chars *C.char, action, location, length C.int) {
	// 0 = Action_Input, 1 = Action_Typing,
	act := Action_Input
	if action == 1 {
		act = Action_Typing
	}

	e := &keyEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_TypingEvent,
			action: act,
		},
		text:     C.GoString(chars),
		location: int32(location),
	}

	theApp.handleEvent(e)
}

//export go_backToUI
func go_backToUI() {
	log.V("nuxui", "go_backToUI ..........")
	select {
	case f := <-theApp.runOnUI:
		f()
		// case e := <-theApp.event:
		// case e := <-theApp.eventWaitDone:
	}
}

func runOnUI(callback func()) {
	go func() {
		theApp.runOnUI <- callback
	}()
	C.backToUI()
}

func requestRedraw() {
	log.V("nuxui", "requestRedraw invalidate")
	C.invalidate()
}

func convertModifierFlags(flags uint) uint32 {
	var mods uint32 = 0
	if flags&C.NSEventModifierFlagShift == C.NSEventModifierFlagShift {
		mods |= Mod_Shift
	}
	if flags&C.NSEventModifierFlagControl == C.NSEventModifierFlagControl {
		mods |= Mod_Control
	}
	if flags&C.NSEventModifierFlagOption == C.NSEventModifierFlagOption {
		mods |= Mod_Alt
	}
	if flags&C.NSEventModifierFlagCommand == C.NSEventModifierFlagCommand {
		mods |= Mod_Super
	}
	if flags&C.NSEventModifierFlagCapsLock == C.NSEventModifierFlagCapsLock {
		mods |= Mod_CapsLock
	}

	return mods
}

func convertVirtualKeyCode(vkcode uint16) KeyCode {
	switch vkcode {
	case C.kVK_ANSI_A:
		return Key_A
	case C.kVK_ANSI_B:
		return Key_B
	case C.kVK_ANSI_C:
		return Key_C
	case C.kVK_ANSI_D:
		return Key_D
	case C.kVK_ANSI_E:
		return Key_E
	case C.kVK_ANSI_F:
		return Key_F
	case C.kVK_ANSI_G:
		return Key_G
	case C.kVK_ANSI_H:
		return Key_H
	case C.kVK_ANSI_I:
		return Key_I
	case C.kVK_ANSI_J:
		return Key_J
	case C.kVK_ANSI_K:
		return Key_K
	case C.kVK_ANSI_L:
		return Key_L
	case C.kVK_ANSI_M:
		return Key_M
	case C.kVK_ANSI_N:
		return Key_N
	case C.kVK_ANSI_O:
		return Key_O
	case C.kVK_ANSI_P:
		return Key_P
	case C.kVK_ANSI_Q:
		return Key_Q
	case C.kVK_ANSI_R:
		return Key_R
	case C.kVK_ANSI_S:
		return Key_S
	case C.kVK_ANSI_T:
		return Key_T
	case C.kVK_ANSI_U:
		return Key_U
	case C.kVK_ANSI_V:
		return Key_V
	case C.kVK_ANSI_W:
		return Key_W
	case C.kVK_ANSI_X:
		return Key_X
	case C.kVK_ANSI_Y:
		return Key_Y
	case C.kVK_ANSI_Z:
		return Key_Z
	case C.kVK_ANSI_1:
		return Key_1
	case C.kVK_ANSI_2:
		return Key_2
	case C.kVK_ANSI_3:
		return Key_3
	case C.kVK_ANSI_4:
		return Key_4
	case C.kVK_ANSI_5:
		return Key_5
	case C.kVK_ANSI_6:
		return Key_6
	case C.kVK_ANSI_7:
		return Key_7
	case C.kVK_ANSI_8:
		return Key_8
	case C.kVK_ANSI_9:
		return Key_9
	case C.kVK_ANSI_0:
		return Key_0
	case C.kVK_F1:
		return Key_F1
	case C.kVK_F2:
		return Key_F2
	case C.kVK_F3:
		return Key_F3
	case C.kVK_F4:
		return Key_F4
	case C.kVK_F5:
		return Key_F5
	case C.kVK_F6:
		return Key_F6
	case C.kVK_F7:
		return Key_F7
	case C.kVK_F8:
		return Key_F8
	case C.kVK_F9:
		return Key_F9
	case C.kVK_F10:
		return Key_F10
	case C.kVK_F11:
		return Key_F11
	case C.kVK_F12:
		return Key_F12
	case C.kVK_F13:
		return Key_F13
	case C.kVK_F14:
		return Key_F14
	case C.kVK_F15:
		return Key_F15
	case C.kVK_F16:
		return Key_F16
	case C.kVK_F17:
		return Key_F17
	case C.kVK_F18:
		return Key_F18
	case C.kVK_F19:
		return Key_F19
	case C.kVK_F20:
		return Key_F20

	case C.kVK_Return:
		return Key_Return
	case C.kVK_Tab:
		return Key_Tab
	case C.kVK_Space:
		return Key_Space
	case C.kVK_Delete:
		return Key_Delete
	case C.kVK_Escape:
		return Key_Escape
	case C.kVK_Command:
		return Key_Super
	case C.kVK_CapsLock:
		return Key_CapsLock
	case C.kVK_Option:
		return Key_AltLeft
	case C.kVK_RightOption:
		return Key_AltRight
	case C.kVK_Shift:
		return Key_ShiftLeft
	case C.kVK_RightShift:
		return Key_ShiftRight
	case C.kVK_Control:
		return Key_ControlLeft
	case C.kVK_RightControl:
		return Key_ControlRight
	case C.kVK_ANSI_Equal:
		return Key_Equal
	case C.kVK_ANSI_Minus:
		return Key_Minus
	case C.kVK_ANSI_LeftBracket:
		return Key_BracketLeft
	case C.kVK_ANSI_RightBracket:
		return Key_BracketRight
	case C.kVK_ANSI_Quote:
		return Key_Quote
	case C.kVK_ANSI_Semicolon:
		return Key_Semicolon
	case C.kVK_ANSI_Comma:
		return Key_Comma
	case C.kVK_ANSI_Slash:
		return Key_Slash
	case C.kVK_ANSI_Backslash:
		return Key_Backslash
	case C.kVK_ANSI_Period:
		return Key_Period
	case C.kVK_ANSI_Grave:
		return Key_Grave

	case C.kVK_ANSI_Keypad1:
		return Key_Pad1
	case C.kVK_ANSI_Keypad2:
		return Key_Pad2
	case C.kVK_ANSI_Keypad3:
		return Key_Pad3
	case C.kVK_ANSI_Keypad4:
		return Key_Pad4
	case C.kVK_ANSI_Keypad5:
		return Key_Pad5
	case C.kVK_ANSI_Keypad6:
		return Key_Pad6
	case C.kVK_ANSI_Keypad7:
		return Key_Pad7
	case C.kVK_ANSI_Keypad8:
		return Key_Pad8
	case C.kVK_ANSI_Keypad9:
		return Key_Pad9
	case C.kVK_ANSI_Keypad0:
		return Key_Pad0
	case C.kVK_ANSI_KeypadPlus:
		return Key_PadAdd
	case C.kVK_ANSI_KeypadMinus:
		return Key_PadSubtract
	case C.kVK_ANSI_KeypadMultiply:
		return Key_PadMultiply
	case C.kVK_ANSI_KeypadDivide:
		return Key_PadDivide
	case C.kVK_ANSI_KeypadEquals:
		return Key_PadEquals
	case C.kVK_ANSI_KeypadClear: // num lock
		return Key_PadClear
	case C.kVK_ANSI_KeypadDecimal:
		return Key_PadDecimal
	case C.kVK_ANSI_KeypadEnter:
		return Key_PadEnter
	case C.kVK_UpArrow:
		return Key_ArrowUp
	case C.kVK_DownArrow:
		return Key_ArrowDown
	case C.kVK_LeftArrow:
		return Key_ArrowLeft
	case C.kVK_RightArrow:
		return Key_ArrowRight
	case C.kVK_PageUp:
		return Key_PageUp
	case C.kVK_PageDown:
		return Key_PageDown
	case C.kVK_Home:
		return Key_Home
	case C.kVK_End:
		return Key_End

	case C.kVK_Help:
		return Key_Help
	case C.kVK_ForwardDelete:
		return Key_ForwardDelete

	case C.kVK_Mute:
		return Key_Mute
	case C.kVK_VolumeUp:
		return Key_VolumeUp
	case C.kVK_VolumeDown:
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
		return Key_Unknown
	}
}
