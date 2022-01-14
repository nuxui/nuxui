// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

type COLORREF uint32

type PAINTSTRUCT struct {
	HDC         uintptr
	Erase       int32
	RcPaint     RECT
	Restore     int32
	IncUpdate   int32
	rgbReserved [32]byte
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type WINDOWPOS struct {
	HWND            uintptr
	HWNDInsertAfter uintptr
	X               int32
	Y               int32
	Cx              int32
	Cy              int32
	Flags           uint32
}

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	IconSm        uintptr
}

type POINT struct {
	X int32
	Y int32
}

type MSG struct {
	HWND    uintptr
	Message uint32
	Wparam  uintptr
	Lparam  uintptr
	Time    uint32
	Pt      POINT
}

// type GpStatus int32
// type GpGraphics struct{}
// type GpPen struct{}
// type GpBrush struct{}
// type GpSolidFill struct{ GpBrush }
// type GpStringFormat struct{}
// type GpFont struct{}
// type GpFontFamily struct{}
// type GpFontCollection struct{}
// type GpRegion struct{}
// type GpPath struct{}
// type ARGB uint32
// type GpUnit int32
// type GpImage struct{}
// type GpBitmap GpImage
// type GpMatrix struct{}
// type GpCustomLineCap struct{}

// type GdiplusStartupInput struct {
// 	GdiplusVersion           uint32
// 	DebugEventCallback       uintptr
// 	SuppressBackgroundThread int32
// 	SuppressExternalCodecs   int32
// }

// type GdiplusStartupOutput struct {
// 	NotificationHook   uintptr
// 	NotificationUnhook uintptr
// }
