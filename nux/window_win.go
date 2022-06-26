// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	"syscall"
	"time"
	"unsafe"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
)

const windowClass = "nux_window_cls"

type nativeWindow struct {
	hwnd        uintptr
	preHdc      uintptr
	hdcBuffer   uintptr
	hBitMap     uintptr
	paintStruct win32.PAINTSTRUCT
	canvas      *canvas
}

func newNativeWindow(attr Attr) *nativeWindow {
	wcname, _ := syscall.UTF16PtrFromString(windowClass)
	cursor, err := win32.LoadCursor(0, win32.IDC_ARROW)
	if err != nil {
		log.E("nuxui", "error LoadCursor: %s", err.Error())
	}

	var wc win32.WNDCLASSEX
	wc = win32.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(wc)),
		Style:         0,
		LpfnWndProc:   syscall.NewCallback(nativeWindowEventHandler),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     theApp.native.ptr,
		HIcon:         0,
		HCursor:       cursor,
		HbrBackground: 0,
		LpszMenuName:  nil,
		LpszClassName: wcname,
		IconSm:        0,
	}

	_, err = win32.RegisterClassEx(&wc)
	if err != nil {
		log.Fatal("nuxui", "error RegisterClass: %s", err.Error())
	}

	width, height := measureWindowSize(attr.GetDimen("width", "50%"), attr.GetDimen("height", "50%"))
	title, _ := syscall.UTF16PtrFromString(attr.GetString("title", ""))

	hwnd, err := win32.CreateWindowEx(
		win32.WS_EX_CLIENTEDGE,
		wcname,
		title,
		win32.WS_OVERLAPPEDWINDOW,
		win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT,
		width,
		height,
		0,
		0,
		theApp.native.ptr,
		0,
	)

	if err != nil {
		log.Fatal("nuxui", "error CreateWindowEx: %s", err.Error())
	}

	me := &nativeWindow{
		hwnd: hwnd,
	}
	return me
}

func (me *nativeWindow) Center() {
}

func (me *nativeWindow) Show() {
	win32.ShowWindow(me.hwnd, win32.SW_SHOWDEFAULT)
	err := win32.UpdateWindow(me.hwnd)
	if err != nil {
		log.E("error UpdateWindow %s", err.Error())
	}
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	var rect win32.RECT
	err := win32.GetClientRect(me.hwnd, &rect)
	if err != nil {
		log.E("nuxui", "error GetClientRect %s", err.Error())
		return 0, 0
	}
	return rect.Right - rect.Left, rect.Bottom - rect.Top
}

func (me *nativeWindow) Title() string {
	title, err := win32.GetWindowText(me.hwnd)
	if err != nil {
		log.E("nux", "error GetWindowText: %s", err.Error())
		return ""
	}
	return title
}

func (me *nativeWindow) SetTitle(title string) {
	if err := win32.SetWindowText(me.hwnd, title); err != nil {
		log.E("nux", "error SetWindowText: %s", err.Error())
	}
}

func (me *nativeWindow) lockCanvas() Canvas {
	hdc, _ := win32.BeginPaint(me.hwnd, &me.paintStruct)
	// most time preHdc == hdc
	if me.preHdc != hdc {
		if me.preHdc != 0 {
			win32.DeleteObject(me.hBitMap)
			win32.DeleteDC(me.preHdc)
			win32.DeleteDC(me.hdcBuffer)
		}
		w, h := me.ContentSize()
		me.preHdc = hdc
		me.hdcBuffer, _ = win32.CreateCompatibleDC(hdc)
		me.hBitMap, _ = win32.CreateCompatibleBitmap(hdc, w, h)
		win32.SelectObject(me.hdcBuffer, me.hBitMap)
		me.canvas = canvasFromHDC(me.hdcBuffer)
	}
	return me.canvas
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	win32.EndPaint(me.hwnd, &me.paintStruct)
}

func (me *nativeWindow) draw(canvas Canvas, decor Widget) {
	if decor != nil {
		if f, ok := decor.(Draw); ok {
			w, h := me.ContentSize()
			win32.PatBlt(me.hdcBuffer, 0, 0, w, h, win32.WHITENESS)
			// canvas.Save()
			canvas.ClipRect(0, 0, float32(w), float32(h))
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			// canvas.Restore()
			canvas.Flush()
			win32.BitBlt(me.preHdc, 0, 0, w, h, me.hdcBuffer, 0, 0, win32.SRCCOPY)
		}
	}
}

func nativeWindowEventHandler(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
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
		handlePointerEvent(hwnd, msg, int32(win32.HIWORD(wParam)), win32.GET_X_LPARAM(lParam), win32.GET_Y_LPARAM(lParam))
	case win32.WM_MOUSEWHEEL,
		win32.WM_MOUSEHWHEEL:
		handleScrollEvent2(hwnd, msg, wParam, lParam)
	case win32.WM_KEYDOWN,
		win32.WM_KEYUP,
		win32.WM_SYSKEYDOWN,
		win32.WM_SYSKEYUP:
		handleKeyEvent(hwnd, msg, wParam, lParam)
	case win32.WM_CHAR:
		if win32.HIWORD(lParam) != 0 { // 0 is IME char, ignore
			keycode := win32.LOWORD(wParam)
			if (keycode >= 0x20 && keycode <= 0x7E) || keycode == 0x09 || keycode == 0x0A || keycode == 0x0B || keycode == 0x0D {
				handleTypingEvent(hwnd, msg, wParam, lParam)
			}
		}
	case win32.WM_IME_COMPOSITION:
		handleTypingEvent(hwnd, msg, wParam, lParam)
	case win32.WM_USER:
		backToUI()
	case win32.WM_CREATE:
	case win32.WM_PAINT:
		// log.I("nuxui", "WM_PAINT")
		theApp.window.draw()
	case win32.WM_SIZE:
		theApp.window.resize()
		// w, h := theApp.window.ContentSize()
		// log.I("nuxui", "WM_SIZE, w=%d, h=%d",w,h)
		InvalidateRect(0, 0, 0, 0)
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
	case win32.WM_SETCURSOR:
		if win32.LOWORD(lParam) == win32.HTCLIENT && theApp.native.cursor != 0 {
			win32.SetCursor(theApp.native.cursor)
			return 1
		}
		return win32.DefWindowProc(hwnd, msg, wParam, lParam)
	case win32.WM_CLOSE:
		win32.DestroyWindow(hwnd)
	case win32.WM_DESTROY:
		win32.PostQuitMessage(0)
	}
	return win32.DefWindowProc(hwnd, msg, wParam, lParam)
}

var lastMouseEvent map[PointerButton]PointerEvent = map[PointerButton]PointerEvent{}

func handlePointerEvent(hwnd uintptr, etype uint32, buttonNumber, x, y int32) bool {
	e := &pointerEvent{
		event: event{
			window: App().MainWindow(),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  ButtonNone,
		kind:    Kind_Mouse,
		x:       float32(x),
		y:       float32(y),
		// pressure: float32(pressure),
		// stage:    int32(stage),
	}

	switch etype {
	case win32.WM_MOUSEMOVE:
		e.event.action = Action_Hover
		e.button = ButtonNone
		e.pointer = 0
	case win32.WM_LBUTTONDOWN:
		e.event.action = Action_Down
		e.button = ButtonPrimary
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_LBUTTONUP:
		e.event.action = Action_Up
		e.button = ButtonPrimary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_RBUTTONDOWN:
		e.event.action = Action_Down
		e.button = ButtonSecondary
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_RBUTTONUP:
		e.event.action = Action_Up
		e.button = ButtonSecondary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_MBUTTONDOWN:
		e.event.action = Action_Down
		e.button = ButtonMiddle
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_MBUTTONUP:
		e.event.action = Action_Up
		e.button = ButtonMiddle
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case win32.WM_XBUTTONDOWN:
		e.event.action = Action_Down
		switch buttonNumber {
		case 1:
			e.button = ButtonX1
		case 2:
			e.button = ButtonX2
		default:
			e.button = PointerButton(buttonNumber)
		}
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case win32.WM_XBUTTONUP:
		e.event.action = Action_Up
		switch buttonNumber {
		case 1:
			e.button = ButtonX1
		case 2:
			e.button = ButtonX2
		default:
			e.button = PointerButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	}

	return App().MainWindow().handlePointerEvent(e)
}

func handleScrollEvent2(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) bool {
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
			window: App().MainWindow(),
			time:   time.Now(),
			etype:  Type_ScrollEvent,
			action: Action_Scroll,
		},
		x:       float32(p.X),
		y:       float32(p.Y),
		scrollX: float32(scrollX),
		scrollY: float32(scrollY),
	}

	return App().MainWindow().handleScrollEvent(e)
}

var lastModifierKeyEvent map[KeyCode]bool = map[KeyCode]bool{}

func handleKeyEvent(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) bool {
	e := &keyEvent{
		event: event{
			window: App().MainWindow(),
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

	return App().MainWindow().handleKeyEvent(e)
}

func handleTypingEvent(hwnd uintptr, msg uint32, wParam uintptr, lParam uintptr) bool {
	e := &typingEvent{
		event: event{
			window: App().MainWindow(),
			time:   time.Now(),
			etype:  Type_TypingEvent,
		},
	}

	if msg == win32.WM_CHAR {
		e.action = Action_Input
		buf := make([]uint16, 2)
		buf[0] = uint16(wParam & 0xffff)
		buf[1] = 0 // uint16((wParam >> 16) & 0xffff)
		e.text = syscall.UTF16ToString(buf)
		// log.V("nux", "typing event WM_CHAR 2 : %s", syscall.UTF16ToString(buf))
	} else {
		hIMC := win32.ImmGetContext(hwnd)

		if lParam&win32.GCS_CURSORPOS == win32.GCS_CURSORPOS {
			e.location = int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_CURSORPOS, 0, 0))
		}

		if lParam&win32.GCS_COMPSTR == win32.GCS_COMPSTR {
			textLen := int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_COMPSTR, 0, 0) + 1)
			buf := make([]uint16, textLen)
			win32.ImmGetCompositionStringW(hIMC, win32.GCS_COMPSTR, uintptr(unsafe.Pointer((&buf[0]))), textLen)
			// log.V("nux", "typing event GCS_COMPSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Preedit
			e.text = syscall.UTF16ToString(buf)

		}

		if lParam&win32.GCS_RESULTSTR == win32.GCS_RESULTSTR {
			textLen := int32(win32.ImmGetCompositionStringW(hIMC, win32.GCS_RESULTSTR, 0, 0) + 1)
			buf := make([]uint16, textLen)
			win32.ImmGetCompositionStringW(hIMC, win32.GCS_RESULTSTR, uintptr(unsafe.Pointer((&buf[0]))), textLen)
			// log.V("nux", "typing event GCS_RESULTSTR: %s", syscall.UTF16ToString(buf))
			e.action = Action_Input
			e.text = syscall.UTF16ToString(buf)
		}

		win32.ImmReleaseContext(hwnd, hIMC)
	}

	return App().MainWindow().handleTypingEvent(e)
}
