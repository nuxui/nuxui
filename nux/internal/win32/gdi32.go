// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
	"unsafe"
)

const (
	SRCCOPY        = 0x00CC0020
	SRCPAINT       = 0x00EE0086
	SRCAND         = 0x008800C6
	SRCINVERT      = 0x00660046
	SRCERASE       = 0x00440328
	NOTSRCCOPY     = 0x00330008
	NOTSRCERASE    = 0x001100A6
	MERGECOPY      = 0x00C000CA
	MERGEPAINT     = 0x00BB0226
	PATCOPY        = 0x00F00021
	PATPAINT       = 0x00FB0A09
	PATINVERT      = 0x005A0049
	DSTINVERT      = 0x00550009
	BLACKNESS      = 0x00000042
	WHITENESS      = 0x00FF0062
	NOMIRRORBITMAP = 0x80000000
	CAPTUREBLT     = 0x40000000
)

var (
	modgdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procGetDeviceCaps          = modgdi32.NewProc("GetDeviceCaps")
	procSaveDC                 = modgdi32.NewProc("SaveDC")
	procRestoreDC              = modgdi32.NewProc("RestoreDC")
	procSelectObject           = modgdi32.NewProc("SelectObject")
	procDeleteObject           = modgdi32.NewProc("DeleteObject")
	procPatBlt                 = modgdi32.NewProc("PatBlt")
	procBitBlt                 = modgdi32.NewProc("BitBlt")
	procPolygon                = modgdi32.NewProc("Polygon")
	procCreateCompatibleDC     = modgdi32.NewProc("CreateCompatibleDC")
	procDeleteDC               = modgdi32.NewProc("DeleteDC")
	procCreateCompatibleBitmap = modgdi32.NewProc("CreateCompatibleBitmap")
	procGetBoundsRect          = modgdi32.NewProc("GetBoundsRect")
)

func SaveDC(hdc uintptr) (nSavedDC int) {
	r0, _, _ := procSaveDC.Call(hdc)
	nSavedDC = int(r0)
	return
}

func RestoreDC(hdc uintptr, nSavedDC int) (err error) {
	r0, _, err := procRestoreDC.Call(hdc, uintptr(nSavedDC))
	if r0 != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-createcompatibledc
func CreateCompatibleDC(hdc uintptr) (ret uintptr, err error) {
	ret, _, err = procCreateCompatibleDC.Call(hdc)
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getdevicecaps
func GetDeviceCaps(hdc uintptr, index int32) int32 {
	ret, _, _ := procGetDeviceCaps.Call(hdc, uintptr(index))
	return int32(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-selectobject
func SelectObject(hdc uintptr, hGdiObj uintptr) (hGdiObjout uintptr, err error) {
	hGdiObjout, _, err = procSelectObject.Call(hdc, hGdiObj)
	if hGdiObjout != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deleteobject
func DeleteObject(hGdiObj uintptr) (err error) {
	ret, _, err := procDeleteObject.Call(hGdiObj)
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-createcompatiblebitmap
func CreateCompatibleBitmap(hdc uintptr, width int32, height int32) (hbitmap uintptr, err error) {
	hbitmap, _, err = procCreateCompatibleBitmap.Call(hdc, uintptr(width), uintptr(height))
	if hbitmap != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-patblt
func PatBlt(hdc uintptr, x, y, w, h, rop int32) (hbitmap uintptr, err error) {
	hbitmap, _, err = procPatBlt.Call(hdc, uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(rop))
	if hbitmap != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-bitblt
func BitBlt(hdc uintptr, x, y, w, h int32, hdcSrc uintptr, x1, x2, rop int32) (err error) {
	ret, _, err := procBitBlt.Call(
		hdc, uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		hdcSrc, uintptr(x1), uintptr(x2), uintptr(rop))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-polygon
func Polygon(hdc uintptr, point *POINT, cpt int32) (err error) {
	ret, _, err := procPolygon.Call(hdc, uintptr(unsafe.Pointer(point)), uintptr(cpt))
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-deletedc
func DeleteDC(hdc uintptr) (err error) {
	ret, _, err := procDeleteDC.Call(hdc)
	if ret != 0 {
		err = nil
	}
	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-getboundsrect
func GetBoundsRect(hdc uintptr, lprect *RECT, flags uint32) (sate uint32, err error) {
	ret, _, err := procGetBoundsRect.Call(hdc, uintptr(unsafe.Pointer(lprect)), uintptr(flags))
	if ret != 0 {
		sate = uint32(ret)
		err = nil
	}
	return
}
