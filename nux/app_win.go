// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

/*
#include <windows.h>
#include <windowsx.h>
int win32_main();
void invalidate(HWND hwnd);
*/
import "C"
import (
	"syscall"
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
	// C.terminate()
}

//export windowAction
func windowAction(windptr C.HWND, msg C.UINT) {
	action := Action_WindowCreated
	switch msg {
	case C.WM_CREATE:
		theApp.window.windptr = windptr
		action = Action_WindowCreated
	case C.WM_PAINT:
		action = Action_WindowDraw
	case C.WM_SIZE:
		action = Action_WindowMeasured
	default:
		log.Fatal("nux", "can not run here.")
	}

	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: action,
		window: theApp.findWindow(windptr),
	}

	theApp.handleEvent(e)
}

func (me *application) findWindow(windptr C.HWND) Window {
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
	// winMain()
	C.win32_main()
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
func go_mouseEvent(windptr C.HWND, etype C.UINT, x, y C.int) {
	log.V("nuxui", "go_mouseEvent etype=%d, x=%d, y=%d", etype, x, y)

	// theApp.handleEvent(e)

}

//export go_scrollEvent
func go_scrollEvent(windptr C.HWND, scrollX, scrollY C.double) {
	log.V("nuxui", "go_scrollEvent scrollX=%f, scrollY=%f", scrollX, scrollY)

	// theApp.handleEvent(e)

}

var lastModifierKeyEvent map[KeyCode]bool = map[KeyCode]bool{}

//export go_drawEvent
func go_drawEvent(windptr C.HWND) {
	e := &event{
		window: theApp.findWindow(windptr),
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: Action_WindowDraw,
	}

	theApp.handleEvent(e)
}

//export go_keyEvent
func go_keyEvent(windptr C.HWND, etype uint, keyCode C.UINT32, modifierFlags uint, repeat byte, chars *C.char) {
	log.V("nux", "etype=0x%X keycode = 0x%X %s, modifierFlags=0x%X", etype, keyCode, KeyCode(keyCode), modifierFlags)
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
	case C.WM_KEYDOWN, C.WM_SYSKEYDOWN:
		e.event.action = Action_Down
	case C.WM_KEYUP, C.WM_SYSKEYUP:
		e.event.action = Action_Up
	}

	theApp.handleEvent(e)
}

//export go_typingEvent
func go_typingEvent(windptr C.HWND) {
	var hIMC C.HIMC = C.ImmGetContext(windptr)
	// if hIMC == 0 {
	// 	// error
	// }
	textLen := C.DWORD(C.ImmGetCompositionStringW(hIMC, C.GCS_RESULTSTR, nil, 0) + 1)
	buf := make([]uint16, textLen)
	C.ImmGetCompositionStringW(hIMC, C.GCS_RESULTSTR, (C.LPVOID)(&buf[0]), textLen)
	C.ImmReleaseContext(windptr, hIMC)
	log.V("nux", "typing event: %s", syscall.UTF16ToString(buf))
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
	// C.backToUI()
}

func requestRedraw() {
	log.V("nuxui", "requestRedraw invalidate")
	w := theApp.MainWindow().(*window)
	C.invalidate(w.windptr)
}

func convertModifierFlags(flags uint) uint32 {

	return 0
}

func convertVirtualKeyCode(vkcode uint16) KeyCode {
	return 0
}
