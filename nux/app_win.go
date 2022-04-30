// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux/internal/win32"
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
	manifest           Manifest
	window             *window
	runOnUI            chan func()
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

func (me *application) Terminate() {
	// win32.terminate()
}

func windowAction(hwnd uintptr, msg uint32) {
	action := Action_WindowCreated
	switch msg {
	case win32.WM_CREATE:
		theApp.window.hwnd = hwnd
		action = Action_WindowCreated
	case win32.WM_PAINT:
		action = Action_WindowDraw
	case win32.WM_SIZE:
		action = Action_WindowMeasured
	default:
		log.Fatal("nux", "can not run here.")
	}

	e := &event{
		time:   time.Now(),
		etype:  Type_WindowEvent,
		action: action,
		window: theApp.findWindow(hwnd),
	}

	theApp.handleEvent(e)
}

func (me *application) findWindow(hwnd uintptr) *window {
	if me.window.hwnd == hwnd {
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

//export go_nativeLoopPrepared
func go_nativeLoopPrepared() {
	theApp.nativeLoopPrepared <- struct{}{}
}

const windowClass = "nux_window_cls"

var (
	theInstance uintptr
)

func run() {
	defer runtime.UnlockOSThread()
	mainThreadId = win32.GetCurrentThreadId()
	var in win32.GdiplusStartupInput
	var out win32.GdiplusStartupOutput
	fmt.Printf("GdiplusStartupInput size = %d\n", unsafe.Sizeof(in))
	fmt.Printf("GdiplusStartupOutput size = %d\n", unsafe.Sizeof(out))
	win32.GdiplusStartup(&in, nil)
	defer win32.GdiplusShutdown()

	var wc win32.WNDCLASSEX
	wcname, _ := syscall.UTF16PtrFromString(windowClass)
	title, _ := syscall.UTF16PtrFromString("title")
	cursor, err := win32.LoadCursor(0, win32.IDC_ARROW)
	if err != nil {
		fmt.Println("LoadCursor error %s", err.Error())
	}
	wc = win32.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(wc)),
		Style:         0,
		LpfnWndProc:   syscall.NewCallback(windowWndProc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     theInstance,
		HIcon:         0,
		HCursor:       cursor,
		HbrBackground: 0,
		LpszMenuName:  nil,
		LpszClassName: wcname,
		IconSm:        0,
	}

	_, err = win32.RegisterClassEx(&wc)
	if err != nil {
		fmt.Println("RegisterClass error %s", err.Error())
	}

	hwnd, err := win32.CreateWindowEx(
		win32.WS_EX_CLIENTEDGE,
		wcname,
		title,
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT,
		800,
		600,
		0,
		0,
		theInstance,
		0,
	)

	if err != nil {
		fmt.Println("CreateWindowEx error %s", err.Error())
	}

	win32.ShowWindow(hwnd, win32.SW_SHOWDEFAULT)
	err = win32.UpdateWindow(hwnd)
	if err != nil {
		fmt.Println("UpdateWindow error %s", err.Error())
	}

	var msg win32.MSG
	var ret int32
	for {
		ret, err = win32.GetMessage(&msg, 0, 0, 0)
		// fmt.Printf("msg = 0x%.4X\n", msg.Message)
		if ret > 0 {
			win32.TranslateMessage(&msg)
			win32.DispatchMessage(&msg)
		} else if ret == 0 { // quit
			break
		} else if err != nil {
			fmt.Println("GetMessage error %s", err.Error())
		}
	}

}

func windowWndProc(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	// fmt.Println("msg = %d", uMsg)
	switch msg {
	case win32.WM_NCHITTEST:
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_MOUSEMOVE,
		win32.WM_LBUTTONDOWN,
		win32.WM_LBUTTONUP,
		win32.WM_MBUTTONDOWN,
		win32.WM_MBUTTONUP,
		win32.WM_RBUTTONDOWN,
		win32.WM_RBUTTONUP,
		win32.WM_XBUTTONDOWN,
		win32.WM_XBUTTONUP:
		mouseEvent(hwnd, msg, int32(win32.HIWORD(wParam)), win32.GET_X_LPARAM(lParam), win32.GET_Y_LPARAM(lParam))
	case win32.WM_MOUSEWHEEL,
		win32.WM_MOUSEHWHEEL:
		mouseScrollEvent(hwnd, msg, wParam, lParam)
	case win32.WM_KEYDOWN,
		win32.WM_KEYUP,
		win32.WM_SYSKEYDOWN,
		win32.WM_SYSKEYUP:
		keyboardEvent(hwnd, msg, wParam, lParam)
	case win32.WM_CHAR:
		if win32.HIWORD(lParam) != 0 { // 0 is IME char, ignore
			keycode := win32.LOWORD(wParam)
			if (keycode >= 0x20 && keycode <= 0x7E) || keycode == 0x09 || keycode == 0x0A || keycode == 0x0B || keycode == 0x0D {
				imeTypeEvent(hwnd, msg, wParam, lParam)
			}
		}
	case win32.WM_IME_COMPOSITION:
		imeTypeEvent(hwnd, msg, wParam, lParam)
	case win32.WM_USER:
		backToUI()
	case win32.WM_CREATE:
		windowAction(hwnd, msg)
		go_nativeLoopPrepared()
	case win32.WM_PAINT:
		windowAction(hwnd, msg)
	case win32.WM_SIZE:
		windowAction(hwnd, msg)
	case win32.WM_NCLBUTTONDBLCLK,
		win32.WM_NCLBUTTONDOWN,
		win32.WM_NCLBUTTONUP,
		win32.WM_NCMBUTTONDBLCLK,
		win32.WM_NCMBUTTONDOWN,
		win32.WM_NCMBUTTONUP,
		win32.WM_NCMOUSEHOVER,
		win32.WM_NCMOUSELEAVE,
		win32.WM_NCMOUSEMOVE,
		win32.WM_NCRBUTTONDBLCLK,
		win32.WM_NCRBUTTONDOWN,
		win32.WM_NCRBUTTONUP,
		win32.WM_NCXBUTTONDBLCLK,
		win32.WM_NCXBUTTONDOWN,
		win32.WM_NCXBUTTONUP:
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_CAPTURECHANGED:
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_LBUTTONDBLCLK,
		win32.WM_MBUTTONDBLCLK,
		win32.WM_RBUTTONDBLCLK,
		win32.WM_XBUTTONDBLCLK:
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_MOUSEACTIVATE,
		win32.WM_MOUSEHOVER,
		win32.WM_MOUSELEAVE:
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_CLOSE:
		win32.DestroyWindow(hwnd)
	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

func startTextInput() {
	// win32.startTextInput()
}

func stopTextInput() {
	// win32.stopTextInput()
}

func setTextInputRect(x, y, w, h float32) {
	if theApp.window == nil {
		log.E("nuxui", "the application not ready")
		return
	}

	hwnd := theApp.window.hwnd
	himc := win32.ImmGetContext(hwnd)
	if himc > 0 {
		comp := &win32.COMPOSITIONFORM{}
		comp.Style = int32(win32.CFS_POINT)
		comp.CurrentPos.X = int32(x)
		comp.CurrentPos.Y = int32(y)
		win32.ImmSetCompositionWindow(himc, comp)
		win32.ImmReleaseContext(hwnd, himc)
	}

}

var lastMouseEvent map[MouseButton]PointerEvent = map[MouseButton]PointerEvent{}

func mouseEvent(hwnd uintptr, etype uint32, buttonNumber, x, y int32) {
	e := &pointerEvent{
		event: event{
			window: theApp.findWindow(hwnd),
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
	case win32.WM_MOUSEMOVE:
		e.event.action = Action_Hover
		e.button = MB_None
		e.pointer = 0
	case win32.WM_LBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Left
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_LBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Left
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_RBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Right
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_RBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Right
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_MBUTTONDOWN:
		e.event.action = Action_Down
		e.button = MB_Middle
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_MBUTTONUP:
		e.event.action = Action_Up
		e.button = MB_Middle
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_XBUTTONDOWN:
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
	case win32.WM_XBUTTONUP:
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

func mouseScrollEvent(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) {
	var p win32.POINT
	p.X = win32.GET_X_LPARAM(lParam)
	p.Y = win32.GET_Y_LPARAM(lParam)
	win32.ScreenToClient(hwnd, &p)

	var aMouseInfo [3]int32
	var lines int32 = 1
	if err := win32.SystemParametersInfoW(
		win32.SPI_GETWHEELSCROLLLINES, 0,
		uintptr(unsafe.Pointer(&aMouseInfo[0])), 0); err == nil {
		lines = aMouseInfo[0] // get lines number of scroll each time
	}

	var scrollX, scrollY float32
	if msg == win32.WM_MOUSEWHEEL {
		scrollY = float32(float32(win32.GET_WHEEL_DELTA_WPARAM(wParam)) * float32(lines) / float32(win32.WHEEL_DELTA))
	} else if msg == win32.WM_MOUSEHWHEEL {
		scrollX = -float32(float32(win32.GET_WHEEL_DELTA_WPARAM(wParam)) * float32(lines) / float32(win32.WHEEL_DELTA))
	}

	e := &scrollEvent{
		event: event{
			window: theApp.findWindow(hwnd),
			time:   time.Now(),
			etype:  Type_ScrollEvent,
			action: Action_Scroll,
		},
		x:       float32(p.X),
		y:       float32(p.Y),
		scrollX: float32(scrollX),
		scrollY: float32(scrollY),
	}

	theApp.handleEvent(e)

}

var lastModifierKeyEvent map[KeyCode]bool = map[KeyCode]bool{}

func keyboardEvent(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) {
	fmt.Println("keyboardEvent ######")
	e := &keyEvent{
		event: event{
			window: theApp.findWindow(hwnd),
			time:   time.Now(),
			etype:  Type_KeyEvent,
			action: Action_None,
		},
		keyCode: KeyCode(win32.LOWORD(wParam)),
		repeat:  false,
		// keyRune: chars,
		// modifierFlags: convertModifierFlags(modifierFlags),
	}

	// if repeat == 1 {
	// 	e.repeat = true
	// }

	switch msg {
	case win32.WM_KEYDOWN, win32.WM_SYSKEYDOWN:
		e.event.action = Action_Down
	case win32.WM_KEYUP, win32.WM_SYSKEYUP:
		e.event.action = Action_Up
	}

	theApp.handleEvent(e)
}

func imeTypeEvent(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) {
	fmt.Println("imeTypeEvent ######")
	e := &typeEvent{
		event: event{
			window: theApp.findWindow(hwnd),
			time:   time.Now(),
			etype:  Type_TypeEvent,
		},
	}

	if msg == win32.WM_CHAR {
		e.action = Action_Input
		buf := make([]uint16, 2)
		buf[0] = uint16(wParam & 0xffff)
		buf[1] = 0 // uint16((wParam >> 16) & 0xffff)
		e.text = syscall.UTF16ToString(buf)
		log.V("nux", "typing event WM_CHAR 2 : %s", syscall.UTF16ToString(buf))
	} else {
		hIMC := win32.ImmGetContext(hwnd)

		if lParam&win32.GCS_CURSORPOS == win32.GCS_CURSORPOS {
			e.location = int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_CURSORPOS, 0, 0))
		}

		if lParam&win32.GCS_COMPSTR == win32.GCS_COMPSTR {
			textLen := int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_COMPSTR, 0, 0) + 1)
			buf := make([]uint16, textLen)
			win32.ImmGetCompositionStringW(hIMC, win32.GCS_COMPSTR, uintptr(unsafe.Pointer((&buf[0]))), textLen)
			log.V("nux", "typing event GCS_COMPSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Preedit
			e.text = syscall.UTF16ToString(buf)

		}

		if lParam&win32.GCS_RESULTSTR == win32.GCS_RESULTSTR {
			textLen := int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_RESULTSTR, 0, 0) + 1)
			buf := make([]uint16, textLen)
			win32.ImmGetCompositionStringW(hIMC, win32.GCS_RESULTSTR, uintptr(unsafe.Pointer((&buf[0]))), textLen)
			log.V("nux", "typing event GCS_RESULTSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Input
			e.text = syscall.UTF16ToString(buf)
		}

		win32.ImmReleaseContext(hwnd, hIMC)
	}

	theApp.handleEvent(e)
}

func backToUI() {
	callback := <-theApp.runOnUI
	callback()
}

func runOnUI(callback func()) {
	go func() {
		theApp.runOnUI <- callback
	}()
	// win32.SendMessage(theApp.MainWindow().(*window).hwnd, win32.WM_USER, 0, 0)
	win32.SendMessage(theApp.MainWindow().(*window).hwnd, win32.WM_USER, 0, 0)
}

func requestRedraw() {
	log.V("nuxui", "requestRedraw invalidate")
	// w := theApp.MainWindow().(*window)
	// win32.SendMessage(w.hwnd, win32.WM_PAINT, 0, 0)
	// win32.InvalidateRect(theApp.window.hwnd, nil, 1)
	win32.RedrawWindow(theApp.window.hwnd, nil, 0, win32.RDW_INVALIDATE)
}

func convertModifierFlags(flags uint32) uint32 {

	return 0
}

func convertVirtualKeyCode(vkcode uint16) KeyCode {
	return 0
}

var mainThreadId uint32

func isMainThread() bool {
	return mainThreadId == win32.GetCurrentThreadId()
}
