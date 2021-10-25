// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

package nux

/*
#import <Carbon/Carbon.h> // for HIToolbox/Events.h
#import <Cocoa/Cocoa.h>

void runApp();
void terminate();
void startTextInput();
void stopTextInput();
void setTextInputRect(float x, float y, float w, float h);
void invalidate();
void backToUI();

uint64 threadID();
*/
import "C"
import (
	"runtime"
	"time"

	"github.com/nuxui/nuxui/log"
)

var theApp = &application{
	runOnUI:            make(chan func()),
	nativeLoopPrepared: make(chan struct{}),
	drawSignal:         make(chan struct{}, drawSignalSize),
}

const (
	drawSignalSize = 50
)

var initThreadID uint64

func init() {
	runtime.LockOSThread()
	initThreadID = uint64(C.threadID())

	timerLoopInstance.init()

	go func() {
		<-theApp.nativeLoopPrepared
		// TODO:: paint not ok?
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
	manifest           Manifest
	window             *window
	runOnUI            chan func()
	nativeLoopPrepared chan struct{}
	drawSignal         chan struct{}
}

func (me *application) OnCreate(data interface{}) {
	if c, ok := me.manifest.(OnCreate); ok {
		c.OnCreate()
	}
}

func (me *application) MainWindow() Window {
	return me.window
}

func (me *application) Manifest() Manifest {
	return me.manifest
}

func (me *application) Terminate() {
	C.terminate()
}

//export go_windowCreated
func go_windowCreated(windptr C.uintptr_t) {
	theApp.window.windptr = windptr
	windowAction(windptr, Action_WindowCreated)
}

//export go_windowResized
func go_windowResized(windptr C.uintptr_t) {
	windowAction(windptr, Action_WindowMeasured)
}

//export go_windowDraw
func go_windowDraw(windptr C.uintptr_t) {
	windowAction(windptr, Action_WindowDraw)
}

func windowAction(windptr C.uintptr_t, action EventAction) {
	e := &windowEvent{
		event: event{
			time:   time.Now(),
			etype:  Type_WindowEvent,
			action: action,
			window: theApp.findWindow(windptr),
		},
	}

	theApp.handleEvent(e)
}

func (me *application) findWindow(windptr C.uintptr_t) *window {
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
	defer runtime.UnlockOSThread()
	if tid := uint64(C.threadID()); tid != initThreadID {
		log.Fatal("nuxui", "main called on thread %d, but init ran on %d", tid, initThreadID)
	}

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
func go_mouseEvent(windptr C.uintptr_t, etype C.NSEventType, x, y C.CGFloat, buttonNumber C.NSInteger, pressure C.float, stage C.NSInteger) {
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
		x:        float32(x),
		y:        float32(y),
		pressure: float32(pressure),
		stage:    int32(stage),
	}

	switch etype {
	case C.NSEventTypeMouseMoved:
		e.action = Action_Hover
		e.button = MB_None
		e.pointer = 0
	case C.NSEventTypeLeftMouseDown:
		e.action = Action_Down
		e.button = MB_Left
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.NSEventTypeLeftMouseUp:
		e.action = Action_Up
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeLeftMouseDragged:
		e.action = Action_Drag
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeRightMouseDown:
		e.action = Action_Down
		e.button = MB_Right
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.NSEventTypeRightMouseUp:
		e.action = Action_Up
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeRightMouseDragged:
		e.action = Action_Drag
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.NSEventTypeOtherMouseDown:
		e.action = Action_Down
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
		e.action = Action_Up
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
		e.action = Action_Drag
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
	case C.NSEventTypePressure:
		// TODO:: stageTransition pressureBehavior NSPressureBehavior
		e.action = Action_Pressure
		e.button = MB_None
	}

	theApp.handleEvent(e)

}

//export go_scrollEvent
func go_scrollEvent(windptr C.uintptr_t, x, y, scrollX, scrollY C.CGFloat) {
	log.V("nuxui", "go_scrollEvent, x=%f, y=%f, scrollX=%f, scrollY=%f", x, y, scrollX, scrollY)
	e := &scrollEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_ScrollEvent,
			action: Action_Scroll,
		},
		x:       float32(x),
		y:       float32(y),
		scrollX: float32(scrollX),
		scrollY: float32(scrollY),
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
	if keyCode == 0x36 {
		keyCode = 0x37 // kVK_Command
	}
	e := &keyEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_KeyEvent,
			action: Action_None,
		},
		keyCode:       KeyCode(keyCode),
		repeat:        false,
		modifierFlags: convertModifierFlags(modifierFlags),
		keyRune:       C.GoString(chars),
	}

	if repeat == 1 {
		e.repeat = true
	}

	switch etype {
	case C.NSEventTypeKeyDown:
		e.action = Action_Down
	case C.NSEventTypeKeyUp:
		e.action = Action_Up
	case C.NSEventTypeFlagsChanged:
		if down, ok := lastModifierKeyEvent[e.keyCode]; ok && down {
			lastModifierKeyEvent[e.keyCode] = false
			e.action = Action_Up
		} else {
			lastModifierKeyEvent[e.keyCode] = true
			e.action = Action_Down
		}
	}

	theApp.handleEvent(e)
}

//export go_typeEvent
func go_typeEvent(windptr C.uintptr_t, chars *C.char, action, location, length C.int) {
	// 0 = Action_Input, 1 = Action_Preedit,
	act := Action_Input
	if action == 1 {
		act = Action_Preedit
	}

	e := &typeEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_TypeEvent,
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
	callback := <-theApp.runOnUI
	callback()

	// select {
	// case f := <-theApp.runOnUI:
	// 	f()
	// 	case e := <-theApp.event:
	// 	case e := <-theApp.eventWaitDone:
	// }
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
