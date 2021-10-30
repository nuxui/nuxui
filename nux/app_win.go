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
void setTextInputRect(HWND hwnd, LONG x, LONG y, LONG w, LONG h);
*/
import "C"
import (
	"runtime"
	"syscall"
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

func init() {
	runtime.LockOSThread()
	timerLoopInstance.init()

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
	manifest           Manifest
	window             *window
	runOnUI            chan func()
	nativeLoopPrepared chan struct{}
	drawSignal         chan struct{}
}

func (me *application) OnCreate(data interface{}) {
}

func (me *application) MainWindow() Window {
	return me.window
}

func (me *application) Manifest() Manifest {
	return me.manifest
}

func (me *application) Terminate() {
	// C.terminate()
}

//export go_windowAction
func go_windowAction(windptr C.HWND, msg C.UINT) {
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

func (me *application) findWindow(windptr C.HWND) *window {
	if me.window.windptr == windptr {
		return me.window
	}
	return nil
}

func (me *application) RequestRedraw(widget Widget) {
	if w := GetWidgetWindow(widget); w != nil {
		C.InvalidateRect(w.(*window).windptr, nil, 0)
	}
}

func run() {
	defer runtime.UnlockOSThread()
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
	if theApp.window == nil {
		log.E("nuxui", "the application not ready")
		return
	}
	C.setTextInputRect(theApp.window.windptr, C.LONG(x), C.LONG(y), C.LONG(w), C.LONG(h))

}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

//export go_mouseEvent
func go_mouseEvent(windptr C.HWND, etype C.UINT, buttonNumber, x, y C.int) {
	log.V("nuxui", "go_mouseEvent etype=%d, x=%d, y=%d", etype, x, y)

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
	case C.WM_MOUSEMOVE:
		e.event.action = Action_Hover
		e.button = MB_None
		e.pointer = 0
	case C.WM_LBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Left
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.WM_LBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.WM_RBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Right
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.WM_RBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.WM_MBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Middle
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.WM_MBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Middle
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case C.WM_XBUTTONDOWN:
		e.event.action = Action_Down
		switch buttonNumber {
		case 1:
			e.button = MB_X1
		case 2:
			e.button = MB_X2
		default:
			e.button = MouseButton(buttonNumber)
		}
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case C.WM_XBUTTONUP:
		e.event.action = Action_Up
		switch buttonNumber {
		case 1:
			e.button = MB_X1
		case 2:
			e.button = MB_X2
		default:
			e.button = MouseButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	}

	theApp.handleEvent(e)

}

//export go_scrollEvent
func go_scrollEvent(windptr C.HWND, x, y, scrollX, scrollY C.float) {
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

//export go_typeEvent
func go_typeEvent(windptr C.HWND, etype uint, wParam C.WPARAM, lParam C.LPARAM) {
	e := &typeEvent{
		event: event{
			window: theApp.findWindow(windptr),
			time:   time.Now(),
			etype:  Type_TypeEvent,
		},
	}

	if etype == C.WM_CHAR {
		e.action = Action_Input
		buf := make([]uint16, 2)
		buf[0] = uint16(wParam & 0xffff)
		buf[1] = 0 // uint16((wParam >> 16) & 0xffff)
		e.text = syscall.UTF16ToString(buf)
		log.V("nux", "typing event WM_CHAR 2 : %s", syscall.UTF16ToString(buf))
	} else {
		var hIMC C.HIMC = C.ImmGetContext(windptr)
		if hIMC == nil {
			log.E("nuxui", "ImmGetContext faild.")
			return
		}

		if lParam&C.GCS_CURSORPOS == C.GCS_CURSORPOS {
			e.location = int32(C.DWORD(C.ImmGetCompositionStringW(hIMC, C.GCS_CURSORPOS, nil, 0)))
		}

		if lParam&C.GCS_COMPSTR == C.GCS_COMPSTR {
			textLen := C.DWORD(C.ImmGetCompositionStringW(hIMC, C.GCS_COMPSTR, nil, 0) + 1)
			buf := make([]uint16, textLen)
			C.ImmGetCompositionStringW(hIMC, C.GCS_COMPSTR, (C.LPVOID)(&buf[0]), textLen)
			log.V("nux", "typing event GCS_COMPSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Preedit
			e.text = syscall.UTF16ToString(buf)

		}

		if lParam&C.GCS_RESULTSTR == C.GCS_RESULTSTR {
			textLen := C.DWORD(C.ImmGetCompositionStringW(hIMC, C.GCS_RESULTSTR, nil, 0) + 1)
			buf := make([]uint16, textLen)
			C.ImmGetCompositionStringW(hIMC, C.GCS_RESULTSTR, (C.LPVOID)(&buf[0]), textLen)
			log.V("nux", "typing event GCS_RESULTSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Input
			e.text = syscall.UTF16ToString(buf)
		}

		C.ImmReleaseContext(windptr, hIMC)
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
	C.SendMessage(theApp.MainWindow().(*window).windptr, C.WM_USER, 0, 0)
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
