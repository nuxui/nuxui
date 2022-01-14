// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import (
	"syscall"
	"unsafe"
)

var (
	moduser32                      = syscall.NewLazyDLL("user32.dll")
	procGetDC                      = moduser32.NewProc("GetDC")
	procReleaseDC                  = moduser32.NewProc("ReleaseDC")
	procBeginPaint                 = moduser32.NewProc("BeginPaint")
	procEndPaint                   = moduser32.NewProc("EndPaint")
	procSendMessageW               = moduser32.NewProc("SendMessageW")
	procCreateWindowExW            = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW             = moduser32.NewProc("DefWindowProcW")
	procDestroyWindow              = moduser32.NewProc("DestroyWindow")
	procDispatchMessageW           = moduser32.NewProc("DispatchMessageW")
	procGetClientRect              = moduser32.NewProc("GetClientRect")
	procGetWindowRect              = moduser32.NewProc("GetWindowRect")
	procGetKeyboardLayout          = moduser32.NewProc("GetKeyboardLayout")
	procGetKeyboardState           = moduser32.NewProc("GetKeyboardState")
	procGetKeyState                = moduser32.NewProc("GetKeyState")
	procGetMessageW                = moduser32.NewProc("GetMessageW")
	procLoadCursorW                = moduser32.NewProc("LoadCursorW")
	procLoadIconW                  = moduser32.NewProc("LoadIconW")
	procMoveWindow                 = moduser32.NewProc("MoveWindow")
	procPostMessageW               = moduser32.NewProc("PostMessageW")
	procPostQuitMessage            = moduser32.NewProc("PostQuitMessage")
	procRegisterClassExW           = moduser32.NewProc("RegisterClassExW")
	procShowWindow                 = moduser32.NewProc("ShowWindow")
	procUpdateWindow               = moduser32.NewProc("UpdateWindow")
	procToUnicodeEx                = moduser32.NewProc("ToUnicodeEx")
	procTranslateMessage           = moduser32.NewProc("TranslateMessage")
	procUnregisterClassW           = moduser32.NewProc("UnregisterClassW")
	procSetWindowLongPtrW          = moduser32.NewProc("SetWindowLongPtrW")
	procGetWindowLongPtrW          = moduser32.NewProc("GetWindowLongPtrW")
	procSetWindowTextW             = moduser32.NewProc("SetWindowTextW")
	procGetWindowTextW             = moduser32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW       = moduser32.NewProc("GetWindowTextLengthW")
	procGetLayeredWindowAttributes = moduser32.NewProc("GetLayeredWindowAttributes")
	procSetLayeredWindowAttributes = moduser32.NewProc("SetLayeredWindowAttributes")
	procGetCursorPos               = moduser32.NewProc("GetCursorPos")
	procScreenToClient             = moduser32.NewProc("ScreenToClient")
	procClientToScreen             = moduser32.NewProc("ClientToScreen")
	procSystemParametersInfoW      = moduser32.NewProc("SystemParametersInfoW")
	procInvalidateRect             = moduser32.NewProc("InvalidateRect")
	procRedrawWindow               = moduser32.NewProc("RedrawWindow")
)

func RGB(r, g, b byte) COLORREF {
	return COLORREF(r) | COLORREF(g)<<8 | COLORREF(b)<<16
}

func GET_X_LPARAM(lp uintptr) int32 {
	return int32(LOWORD(lp))
}

func GET_Y_LPARAM(lp uintptr) int32 {
	return int32(HIWORD(lp))
}

func GET_WHEEL_DELTA_WPARAM(lp uintptr) int16 {
	return int16(HIWORD(lp))
}

func LOWORD(l uintptr) uint16 {
	return uint16(uint32(l))
}

func HIWORD(l uintptr) uint16 {
	return uint16(uint32(l >> 16))
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-beginpaint
func BeginPaint(hwnd uintptr, ps *PAINTSTRUCT) (hdc uintptr, err error) {
	hdc, _, err = procBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(ps)))
	if hdc != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-endpaint
func EndPaint(hwnd uintptr, ps *PAINTSTRUCT) {
	procEndPaint.Call(hwnd, uintptr(unsafe.Pointer(ps)))
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-defwindowprocw
func DefWindowProc(hwnd uintptr, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	lResult, _, _ = procDefWindowProcW.Call(hwnd, uintptr(uMsg), uintptr(wParam), uintptr(lParam))
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerclassexw
func RegisterClassEx(wc *WNDCLASSEX) (atom uint16, err error) {
	r0, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(wc)))
	atom = uint16(r0)
	if atom != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-createwindowexw
func CreateWindowEx(exstyle uint32, className *uint16, windowText *uint16, style uint32, x int32, y int32, width int32, height int32, parent uintptr, menu uintptr, hInstance uintptr, lpParam uintptr) (hwnd uintptr, err error) {
	hwnd, _, err = procCreateWindowExW.Call(uintptr(exstyle), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowText)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(parent), uintptr(menu), uintptr(hInstance), uintptr(lpParam))
	if hwnd != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-showwindow
func ShowWindow(hwnd uintptr, cmdshow int32) (wasvisible bool) {
	r0, _, _ := procShowWindow.Call(hwnd, uintptr(cmdshow))
	wasvisible = r0 != 0
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-updatewindow
func UpdateWindow(hwnd uintptr) (err error) {
	r0, _, err := procUpdateWindow.Call(hwnd)
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmessagew
func GetMessage(msg *MSG, hwnd uintptr, msgfiltermin uint32, msgfiltermax uint32) (ret int32, err error) {
	r0, _, err := procGetMessageW.Call(uintptr(unsafe.Pointer(msg)), hwnd, uintptr(msgfiltermin), uintptr(msgfiltermax))
	ret = int32(r0)
	if ret >= 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-translatemessage
func TranslateMessage(msg *MSG) (done bool) {
	r0, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	done = r0 != 0
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-dispatchmessagew
func DispatchMessage(msg *MSG) (ret int32) {
	r0, _, _ := procDispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	ret = int32(r0)
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadcursorw
func LoadCursor(hInstance uintptr, cursorName uintptr) (cursor uintptr, err error) {
	cursor, _, err = procLoadCursorW.Call(uintptr(hInstance), uintptr(cursorName))
	if cursor != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-destroywindow
func DestroyWindow(hwnd uintptr) (err error) {
	r0, _, err := procDestroyWindow.Call(hwnd)
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-postquitmessage
func PostQuitMessage(exitCode int32) {
	procPostQuitMessage.Call(uintptr(exitCode), 0, 0)
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagew
func SendMessage(hwnd uintptr, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	lResult, _, _ = procSendMessageW.Call(hwnd, uintptr(uMsg), uintptr(wParam), uintptr(lParam), 0, 0)
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowrect
func GetWindowRect(hwnd uintptr, rect *RECT) (err error) {
	r0, _, err := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(rect)))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getclientrect
func GetClientRect(hwnd uintptr, rect *RECT) (err error) {
	r0, _, err := procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(rect)))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getlayeredwindowattributes
func GetLayeredWindowAttributes(hwnd uintptr, pcrKey *uint32, pbAlpha *byte, pdwFlags *int32) (err error) {
	r0, _, err := procGetLayeredWindowAttributes.Call(hwnd,
		uintptr(unsafe.Pointer(pcrKey)),
		uintptr(unsafe.Pointer(pbAlpha)),
		uintptr(unsafe.Pointer(pdwFlags)))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setlayeredwindowattributes
func SetLayeredWindowAttributes(hwnd uintptr, pcrKey uint32, pbAlpha byte, pdwFlags int32) (err error) {
	r0, _, err := procSetLayeredWindowAttributes.Call(hwnd,
		uintptr(pcrKey),
		uintptr(pbAlpha),
		uintptr(pdwFlags))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowlongptrw
func SetWindowLong(hwnd uintptr, nIndex int, dwNewLong uintptr) (ret uintptr, err error) {
	ret, _, err = procSetWindowLongPtrW.Call(hwnd, uintptr(nIndex), dwNewLong)
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowlongptrw
func GetWindowLong(hwnd uintptr, nIndex int) (ret uintptr, err error) {
	ret, _, err = procGetWindowLongPtrW.Call(hwnd, uintptr(nIndex))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowtextw
func SetWindowText(hwnd uintptr, title string) (err error) {
	text, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	r0, _, err := procSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(text)))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextw
func GetWindowText(hwnd uintptr) (title string, err error) {
	textLen, _, _ := procGetWindowTextLengthW.Call(hwnd)
	textLen++
	buf := make([]uint16, textLen)
	len, _, err := procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(textLen))
	if len == textLen-1 {
		title = syscall.UTF16ToString(buf)
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextlengthw
func GetWindowTextLength(hwnd uintptr) (length int) {
	ret, _, _ := procGetWindowTextLengthW.Call(hwnd)
	length = int(ret)
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos
func GetCursorPos(point *POINT) (err error) {
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(point)))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-screentoclient
func ScreenToClient(hwnd uintptr, point *POINT) (err error) {
	ret, _, err := procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(point)))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-clienttoscreen
func ClientToScreen(hwnd uintptr, point *POINT) (err error) {
	ret, _, err := procClientToScreen.Call(hwnd, uintptr(unsafe.Pointer(point)))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-invalidaterect
func InvalidateRect(hwnd uintptr, rect *RECT, bErase int32) (err error) {
	ret, _, err := procInvalidateRect.Call(hwnd, uintptr(unsafe.Pointer(rect)), uintptr(bErase))
	if ret != 0 {
		err = nil
	}
	return
}

func SystemParametersInfoW(uiAction uint32, uiParam uint32, pvParam uintptr, fWinIni uint32) (err error) {
	ret, _, err := procSystemParametersInfoW.Call(
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(pvParam),
		uintptr(fWinIni),
	)
	if ret != 0 {
		err = nil
	}
	return
}

func RedrawWindow(hwnd uintptr, lprcUpdate *RECT, hrgnUpdate uintptr, flags uint32) (err error) {
	ret, _, err := procRedrawWindow.Call(
		hwnd,
		uintptr(unsafe.Pointer(lprcUpdate)),
		hrgnUpdate,
		uintptr(flags),
	)
	if ret != 0 {
		err = nil
	}
	return
}
