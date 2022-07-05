// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import "unsafe"

type DWORD uint32
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

// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew
// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.14393.0/um/commdlg.h
type OPENFILENAME struct { // OPENFILENAMEW
	StructSize      uint32         //DWORD
	Owner           uintptr        //HWND
	Instance        uintptr        //HINSTANCE
	Filter          *uint16        //LPCWSTR
	CustomFilter    *uint16        //LPWSTR
	MaxCustomFilter uint32         //DWORD
	FilterIndex     uint32         //DWORD
	File            *uint16        //LPWSTR
	MaxFile         uint32         //DWORD
	FileTitle       *uint16        //LPWSTR
	MaxFileTitle    uint32         //DWORD
	InitialDir      *uint16        //LPCWSTR
	Title           *uint16        //LPCWSTR
	Flags           uint32         //DWORD
	FileOffset      uint16         //WORD
	FileExtension   uint16         //WORD
	DefExt          *uint16        //LPCWSTR
	CustData        uintptr        //LPARAM
	FnHook          uintptr        //LPOFNHOOKPROC
	TemplateName    *uint16        //LPCWSTR
	PvReserved      unsafe.Pointer //void*
	DwReserved      uint32         //DWORD
	FlagsEx         uint32         //DWORD
}

// https://docs.microsoft.com/en-us/windows/win32/api/shlobj_core/ns-shlobj_core-browseinfow
// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.14393.0/um/ShlObj_core.h
type BROWSEINFO struct { //BROWSEINFOW
	Owner       uintptr //HWND
	Root        uintptr //PCIDLIST_ABSOLUTE
	DisplayName *uint16 //LPWSTR
	Title       *uint16 //LPCWSTR
	Flags       uint32  //UINT
	Fn          uintptr //BFFCALLBACK
	Param       uintptr //LPARAM
	Image       int32   //int
}

// https://docs.microsoft.com/en-us/windows/win32/api/shtypes/ns-shtypes-shitemid
type SHITEMID struct {
	CB   uint16
	ABID [1]byte
}

// https://docs.microsoft.com/en-us/windows/win32/api/shtypes/ns-shtypes-itemidlist
type ITEMIDLIST struct {
	ID SHITEMID
}

// ------------------------ GDI+ -------------------------

type GpState int32
type GpStatus int32
type GpGraphics struct{}
type GpPen struct{}
type GpBrush struct{}
type GpStringFormat struct{}
type GpFont struct{}
type GpFontFamily struct{}
type GpFontCollection struct{}
type GpRegion struct{}
type GpPath struct{}
type ARGB uint32
type GpUnit int32
type GpImage struct{}
type GpBitmap GpImage
type GpMatrix struct{}
type GpCustomLineCap struct{}

// Windows Types
type HANDLE uintptr
type HPALETTE uintptr
type HBITMAP uintptr
type HDC uintptr
type HWND uintptr

// Enum types
type GpBrushType int32
type GpPenType int32
type GpLineCap int32
type GpLineJoin int32
type GpDashCap int32
type GpDashStyle int32
type GpPenAlignment int32
type GpMatrixOrder int32
type GpCombineMode int32

type BrushType GpBrushType
type PenType GpPenType
type LineCap GpLineCap
type LineJoin GpLineJoin
type DashCap GpDashCap
type DashStyle GpDashStyle
type PenAlignment GpPenAlignment
type MatrixOrder GpMatrixOrder

type GdiplusStartupInput struct {
	GdiplusVersion           uint32
	DebugEventCallback       uintptr
	SuppressBackgroundThread int32
	SuppressExternalCodecs   int32
}

type GdiplusStartupOutput struct {
	NotificationHook   uintptr
	NotificationUnhook uintptr
}

type RectF struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type PointF struct {
	X float32
	Y float32
}

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

type Point struct {
	X int32
	Y int32
}

// type EncoderParameter struct {
// 	Guid           ole.GUID
// 	NumberOfValues uint32
// 	TypeAPI        uint32
// 	Value          uintptr
// }

// type EncoderParameters struct {
// 	Count     uint32
// 	Parameter [1]EncoderParameter
// }

// In-memory pixel data formats:
// bits 0-7 = format index
// bits 8-15 = pixel size (in bits)
// bits 16-23 = flags
// bits 24-31 = reserved

type PixelFormat int32

const (
	PixelFormatIndexed        = 0x00010000 // Indexes into a palette
	PixelFormatGDI            = 0x00020000 // Is a GDI-supported format
	PixelFormatAlpha          = 0x00040000 // Has an alpha component
	PixelFormatPAlpha         = 0x00080000 // Pre-multiplied alpha
	PixelFormatExtended       = 0x00100000 // Extended color 16 bits/channel
	PixelFormatCanonical      = 0x00200000
	PixelFormatUndefined      = 0
	PixelFormatDontCare       = 0
	PixelFormat1bppIndexed    = (1 | (1 << 8) | PixelFormatIndexed | PixelFormatGDI)
	PixelFormat4bppIndexed    = (2 | (4 << 8) | PixelFormatIndexed | PixelFormatGDI)
	PixelFormat8bppIndexed    = (3 | (8 << 8) | PixelFormatIndexed | PixelFormatGDI)
	PixelFormat16bppGrayScale = (4 | (16 << 8) | PixelFormatExtended)
	PixelFormat16bppRGB555    = (5 | (16 << 8) | PixelFormatGDI)
	PixelFormat16bppRGB565    = (6 | (16 << 8) | PixelFormatGDI)
	PixelFormat16bppARGB1555  = (7 | (16 << 8) | PixelFormatAlpha | PixelFormatGDI)
	PixelFormat24bppRGB       = (8 | (24 << 8) | PixelFormatGDI)
	PixelFormat32bppRGB       = (9 | (32 << 8) | PixelFormatGDI)
	PixelFormat32bppARGB      = (10 | (32 << 8) | PixelFormatAlpha | PixelFormatGDI | PixelFormatCanonical)
	PixelFormat32bppPARGB     = (11 | (32 << 8) | PixelFormatAlpha | PixelFormatPAlpha | PixelFormatGDI)
	PixelFormat48bppRGB       = (12 | (48 << 8) | PixelFormatExtended)
	PixelFormat64bppARGB      = (13 | (64 << 8) | PixelFormatAlpha | PixelFormatCanonical | PixelFormatExtended)
	PixelFormat64bppPARGB     = (14 | (64 << 8) | PixelFormatAlpha | PixelFormatPAlpha | PixelFormatExtended)
	PixelFormat32bppCMYK      = (15 | (32 << 8))
	PixelFormatMax            = 16
)

func NewRect(x, y, width, height int32) *Rect {
	return &Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func NewRectF(x, y, width, height float32) *RectF {
	return &RectF{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (rect *Rect) Left() int32 {
	return rect.X
}

func (rect *Rect) Top() int32 {
	return rect.Y
}

func (rect *RectF) Left() float32 {
	return rect.X
}

func (rect *RectF) Top() float32 {
	return rect.Y
}

func (rect *Rect) Right() int32 {
	return rect.X + rect.Width
}

func (rect *Rect) Bottom() int32 {
	return rect.Y + rect.Height
}

func (rect *RectF) Right() float32 {
	return rect.X + rect.Width
}

func (rect *RectF) Bottom() float32 {
	return rect.Y + rect.Height
}

func (s GpStatus) String() string {
	switch s {
	case Ok:
		return "Ok"
	case GenericError:
		return "GenericError"
	case InvalidParameter:
		return "InvalidParameter"
	case OutOfMemory:
		return "OutOfMemory"
	case ObjectBusy:
		return "ObjectBusy"
	case InsufficientBuffer:
		return "InsufficientBuffer"
	case NotImplemented:
		return "NotImplemented"
	case Win32Error:
		return "Win32Error"
	case WrongState:
		return "WrongState"
	case Aborted:
		return "Aborted"
	case FileNotFound:
		return "FileNotFound"
	case ValueOverflow:
		return "ValueOverflow"
	case AccessDenied:
		return "AccessDenied"
	case UnknownImageFormat:
		return "UnknownImageFormat"
	case FontFamilyNotFound:
		return "FontFamilyNotFound"
	case FontStyleNotFound:
		return "FontStyleNotFound"
	case NotTrueTypeFont:
		return "NotTrueTypeFont"
	case UnsupportedGdiplusVersion:
		return "UnsupportedGdiplusVersion"
	case GdiplusNotInitialized:
		return "GdiplusNotInitialized"
	case PropertyNotFound:
		return "PropertyNotFound"
	case PropertyNotSupported:
		return "PropertyNotSupported"
	case ProfileNotFound:
		return "ProfileNotFound"
	}

	return "Unknown Status Value"
}
