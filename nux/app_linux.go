// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

/*


#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <X11/Xlib.h>

void run();
void invalidate(Display *display, Window window);
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
	// C.terminate()
}

//export go_windowCreated
func go_windowCreated(windptr C.uintptr_t, display *C.Display, visual *C.Visual) {
	theApp.window.windptr = windptr
	theApp.window.display = display
	theApp.window.visual = visual
	windowAction(windptr, Action_WindowCreated)
}

//export go_windowResized
func go_windowResized(windptr C.uintptr_t, width, height int32) {
	w := theApp.findWindow(windptr)
	if width != w.width || height != w.height {
		windowAction(windptr, Action_WindowMeasured)
	}
}

//export go_windowDraw
func go_windowDraw(windptr C.uintptr_t) {
	windowAction(windptr, Action_WindowDraw)
}

func windowAction(windptr C.uintptr_t, action EventAction) {
	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: action,
		window: theApp.findWindow(windptr),
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
	// C.setTextInputRect(C.float(x), C.float(y), C.float(w), C.float(h))
}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

//export go_mouseEvent
func go_mouseEvent(windptr C.uintptr_t, etype uint, x, y float32, buttonNumber int32, pressure float32, stage int32) {
	// log.V("nuxui", "go_mouseEvent x=%f, y=%f", x, y)
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

	theApp.handleEvent(e)

}

//export go_scrollEvent
func go_scrollEvent(windptr C.uintptr_t, x, y, scrollX, scrollY float32) {
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
		e.event.action = Action_Down
	case C.KeyRelease:
		e.event.action = Action_Up
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
	// TODO::
}

func requestRedraw() {
	w := theApp.MainWindow().(*window)
	C.invalidate(w.display, w.windptr)
}

func convertModifierFlags(flags uint) uint32 {
	// TODO::
	return 0
}
