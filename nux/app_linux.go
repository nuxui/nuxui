// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

/*
#cgo pkg-config: x11

#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <X11/Xlib.h>

void run();
void invalidate(Display *display, Window window);
void setTextInputRect(short x, short y);
void runOnUI(Display *display, Window window);
int isMainThread();
*/
import "C"
import (
	"runtime"
	"time"

	"nuxui.org/nuxui/log"
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
	runtime.LockOSThread()
	timerLoopInstance.init()

	go func() {
		<-theApp.nativeLoopPrepared
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

func (me *application) OnCreate(data any) {
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
	// C.terminate()
}

//export go_windowCreated
func go_windowCreated(windptr C.Window, display *C.Display, visual *C.Visual) {
	theApp.window.windptr = windptr
	theApp.window.display = display
	theApp.window.visual = visual
	windowAction(windptr, Action_WindowCreated)
}

//export go_windowResized
func go_windowResized(windptr C.Window, width, height int32) {
	w := theApp.findWindow(windptr)
	if width != w.width || height != w.height {
		windowAction(windptr, Action_WindowMeasured)
	}
}

//export go_windowDraw
func go_windowDraw(windptr C.Window) {
	windowAction(windptr, Action_WindowDraw)
}

func windowAction(windptr C.Window, action EventAction) {
	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: action,
		window: theApp.findWindow(windptr),
	}

	theApp.handleEvent(e)
}

func (me *application) findWindow(windptr C.Window) *window {
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
	C.run()
}

//export go_nativeLoopPrepared
func go_nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

func startTextInput() {
	// C.startTextInput()
}

func stopTextInput() {
	// C.stopTextInput()
}

func setTextInputRect(x, y, w, h float32) {
	C.setTextInputRect(C.short(x), C.short(y))
}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

//export go_mouseEvent
func go_mouseEvent(windptr C.Window, etype int32, x, y float32, buttonNumber uint32) {
	log.V("nuxui", "go_mouseEvent x=%f, y=%f, buttonNumber=%d", x, y, buttonNumber)
	e := &pointerEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  MB_None,
		kind:    Kind_Mouse,
		x:       float32(x),
		y:       float32(y),
		// pressure: float32(pressure),
		// stage:    int32(stage),
	}

	switch etype {
	case C.MotionNotify:
		e.action = Action_Move
	case C.ButtonPress:
		e.action = Action_Down
	case C.ButtonRelease:
		e.action = Action_Up
	}

	switch buttonNumber {
	case C.Button1:
		e.button = MB_Left
	case C.Button2:
		e.button = MB_Right
	case C.Button3:
		e.button = MB_Middle
	default: // > Button7
		// MB_X1 = MB_Middle = 2
		if buttonNumber == 2 {
			e.button = MB_X1
		} else if buttonNumber == 8 {
			e.button = MB_X2
		}
	}

	theApp.handleEvent(e)
}

//export go_scrollEvent
func go_scrollEvent(windptr C.Window, x, y, scrollX, scrollY float32) {
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
func go_drawEvent(windptr C.Window) {
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
func go_keyEvent(windptr C.Window, etype uint, keyCode uint16, modifierFlags uint, repeat byte, chars *C.char) {
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
	case C.KeyPress:
		e.action = Action_Down
	case C.KeyRelease:
		e.action = Action_Up
	}

	theApp.handleEvent(e)
}

//export go_typeEvent
func go_typeEvent(action C.int, chars *C.char, length C.int, location C.int) {
	// 0 = Action_Input, 1 = Action_Preedit,
	act := Action_Input
	if action == 1 {
		act = Action_Preedit
	}

	e := &typeEvent{
		event: event{
			window: theApp.MainWindow(),
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
}

func runOnUI(callback func()) {
	go func() {
		theApp.runOnUI <- callback
	}()
	w := theApp.MainWindow().(*window)
	C.runOnUI(w.display, w.windptr)
}

func requestRedraw() {
	w := theApp.MainWindow().(*window)
	C.invalidate(w.display, w.windptr)
}

func convertModifierFlags(flags uint) uint32 {
	// TODO::
	return 0
}

func isMainThread() bool {
	return C.isMainThread() > 0
}
