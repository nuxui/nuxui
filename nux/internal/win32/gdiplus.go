// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

// https://doxygen.reactos.org/d3/d39/gdiplusflat_8h_source.html
// https://doxygen.reactos.org/d1/d07/gdiplus__private_8h_source.html
// https://github.com/mono/libgdiplus/blob/main/src/GdiPlusFlat.h
// https://github.com/Alexpux/mingw-w64/blob/master/mingw-w64-headers/include/gdiplus/gdiplusflat.h
// https://github.com/Alexpux/mingw-w64/blob/master/mingw-w64-headers/include/gdiplus/gdiplusgraphics.h

import (
	"math"
	"syscall"
	"unsafe"
)

var (
	modgdiplus                                = syscall.NewLazyDLL("gdiplus.dll")
	procGdiplusShutdown                       = modgdiplus.NewProc("GdiplusShutdown")
	procGdiplusStartup                        = modgdiplus.NewProc("GdiplusStartup")
	procGdipCreateFromHDC                     = modgdiplus.NewProc("GdipCreateFromHDC")
	procGdipCreateFromHDC2                    = modgdiplus.NewProc("GdipCreateFromHDC2")
	procGdipCreateFromHWND                    = modgdiplus.NewProc("GdipCreateFromHWND")
	procGdipCreateFromHWNDICM                 = modgdiplus.NewProc("GdipCreateFromHWNDICM")
	procGdipDeleteGraphics                    = modgdiplus.NewProc("GdipDeleteGraphics")
	procGdipGetDC                             = modgdiplus.NewProc("GdipGetDC")
	procGdipReleaseDC                         = modgdiplus.NewProc("GdipReleaseDC")
	procGdipSaveGraphics                      = modgdiplus.NewProc("GdipSaveGraphics")
	procGdipRestoreGraphics                   = modgdiplus.NewProc("GdipRestoreGraphics")
	procGdipTranslateWorldTransform           = modgdiplus.NewProc("GdipTranslateWorldTransform")
	procGdipScaleWorldTransform               = modgdiplus.NewProc("GdipScaleWorldTransform")
	procGdipRotateWorldTransform              = modgdiplus.NewProc("GdipRotateWorldTransform")
	procGdipSetClipRect                       = modgdiplus.NewProc("GdipSetClipRect")
	procGdipSetCompositingMode                = modgdiplus.NewProc("GdipSetCompositingMode")
	procGdipSetRenderingOrigin                = modgdiplus.NewProc("GdipSetRenderingOrigin")
	procGdipSetCompositingQuality             = modgdiplus.NewProc("GdipSetCompositingQuality")
	procGdipSetSmoothingMode                  = modgdiplus.NewProc("GdipSetSmoothingMode")
	procGdipGetSmoothingMode                  = modgdiplus.NewProc("GdipGetSmoothingMode")
	procGdipSetPixelOffsetMode                = modgdiplus.NewProc("GdipSetPixelOffsetMode")
	procGdipSetInterpolationMode              = modgdiplus.NewProc("GdipSetInterpolationMode")
	procGdipSetTextRenderingHint              = modgdiplus.NewProc("GdipSetTextRenderingHint")
	procGdipGraphicsClear                     = modgdiplus.NewProc("GdipGraphicsClear")
	procGdipDrawLine                          = modgdiplus.NewProc("GdipDrawLine")
	procGdipDrawLineI                         = modgdiplus.NewProc("GdipDrawLineI")
	procGdipDrawArc                           = modgdiplus.NewProc("GdipDrawArc")
	procGdipDrawArcI                          = modgdiplus.NewProc("GdipDrawArcI")
	procGdipDrawBezier                        = modgdiplus.NewProc("GdipDrawBezier")
	procGdipDrawBezierI                       = modgdiplus.NewProc("GdipDrawBezierI")
	procGdipDrawRectangle                     = modgdiplus.NewProc("GdipDrawRectangle")
	procGdipDrawRectangleI                    = modgdiplus.NewProc("GdipDrawRectangleI")
	procGdipDrawEllipse                       = modgdiplus.NewProc("GdipDrawEllipse")
	procGdipDrawEllipseI                      = modgdiplus.NewProc("GdipDrawEllipseI")
	procGdipDrawPie                           = modgdiplus.NewProc("GdipDrawPie")
	procGdipDrawPieI                          = modgdiplus.NewProc("GdipDrawPieI")
	procGdipDrawPolygonI                      = modgdiplus.NewProc("GdipDrawPolygonI")
	procGdipDrawPolygon                       = modgdiplus.NewProc("GdipDrawPolygon")
	procGdipDrawPath                          = modgdiplus.NewProc("GdipDrawPath")
	procGdipDrawString                        = modgdiplus.NewProc("GdipDrawString")
	procGdipDrawImage                         = modgdiplus.NewProc("GdipDrawImage")
	procGdipDrawImageI                        = modgdiplus.NewProc("GdipDrawImageI")
	procGdipDrawImageRect                     = modgdiplus.NewProc("GdipDrawImageRect")
	procGdipDrawImageRectI                    = modgdiplus.NewProc("GdipDrawImageRectI")
	procGdipFillRectangle                     = modgdiplus.NewProc("GdipFillRectangle")
	procGdipFillRectangleI                    = modgdiplus.NewProc("GdipFillRectangleI")
	procGdipFillPolygon                       = modgdiplus.NewProc("GdipFillPolygon")
	procGdipFillPolygonI                      = modgdiplus.NewProc("GdipFillPolygonI")
	procGdipFillPath                          = modgdiplus.NewProc("GdipFillPath")
	procGdipFillEllipse                       = modgdiplus.NewProc("GdipFillEllipse")
	procGdipFillEllipseI                      = modgdiplus.NewProc("GdipFillEllipseI")
	procGdipMeasureString                     = modgdiplus.NewProc("GdipMeasureString")
	procGdipMeasureCharacterRanges            = modgdiplus.NewProc("GdipMeasureCharacterRanges")
	procGdipCreatePen1                        = modgdiplus.NewProc("GdipCreatePen1")
	procGdipCreatePen2                        = modgdiplus.NewProc("GdipCreatePen2")
	procGdipClonePen                          = modgdiplus.NewProc("GdipClonePen")
	procGdipDeletePen                         = modgdiplus.NewProc("GdipDeletePen")
	procGdipSetPenWidth                       = modgdiplus.NewProc("GdipSetPenWidth")
	procGdipGetPenWidth                       = modgdiplus.NewProc("GdipGetPenWidth")
	procGdipSetPenLineCap197819               = modgdiplus.NewProc("GdipSetPenLineCap197819")
	procGdipSetPenStartCap                    = modgdiplus.NewProc("GdipSetPenStartCap")
	procGdipSetPenEndCap                      = modgdiplus.NewProc("GdipSetPenEndCap")
	procGdipSetPenDashCap197819               = modgdiplus.NewProc("GdipSetPenDashCap197819")
	procGdipGetPenStartCap                    = modgdiplus.NewProc("GdipGetPenStartCap")
	procGdipGetPenEndCap                      = modgdiplus.NewProc("GdipGetPenEndCap")
	procGdipGetPenDashCap197819               = modgdiplus.NewProc("GdipGetPenDashCap197819")
	procGdipSetPenLineJoin                    = modgdiplus.NewProc("GdipSetPenLineJoin")
	procGdipGetPenLineJoin                    = modgdiplus.NewProc("GdipGetPenLineJoin")
	procGdipSetPenCustomStartCap              = modgdiplus.NewProc("GdipSetPenCustomStartCap")
	procGdipGetPenCustomStartCap              = modgdiplus.NewProc("GdipGetPenCustomStartCap")
	procGdipSetPenCustomEndCap                = modgdiplus.NewProc("GdipSetPenCustomEndCap")
	procGdipGetPenCustomEndCap                = modgdiplus.NewProc("GdipGetPenCustomEndCap")
	procGdipSetPenMiterLimit                  = modgdiplus.NewProc("GdipSetPenMiterLimit")
	procGdipGetPenMiterLimit                  = modgdiplus.NewProc("GdipGetPenMiterLimit")
	procGdipSetPenMode                        = modgdiplus.NewProc("GdipSetPenMode")
	procGdipGetPenMode                        = modgdiplus.NewProc("GdipGetPenMode")
	procGdipSetPenTransform                   = modgdiplus.NewProc("GdipSetPenTransform")
	procGdipGetPenTransform                   = modgdiplus.NewProc("GdipGetPenTransform")
	procGdipResetPenTransform                 = modgdiplus.NewProc("GdipResetPenTransform")
	procGdipMultiplyPenTransform              = modgdiplus.NewProc("GdipMultiplyPenTransform")
	procGdipTranslatePenTransform             = modgdiplus.NewProc("GdipTranslatePenTransform")
	procGdipScalePenTransform                 = modgdiplus.NewProc("GdipScalePenTransform")
	procGdipRotatePenTransform                = modgdiplus.NewProc("GdipRotatePenTransform")
	procGdipSetPenColor                       = modgdiplus.NewProc("GdipSetPenColor")
	procGdipGetPenColor                       = modgdiplus.NewProc("GdipGetPenColor")
	procGdipSetPenBrushFill                   = modgdiplus.NewProc("GdipSetPenBrushFill")
	procGdipGetPenBrushFill                   = modgdiplus.NewProc("GdipGetPenBrushFill")
	procGdipGetPenFillType                    = modgdiplus.NewProc("GdipGetPenFillType")
	procGdipGetPenDashStyle                   = modgdiplus.NewProc("GdipGetPenDashStyle")
	procGdipSetPenDashStyle                   = modgdiplus.NewProc("GdipSetPenDashStyle")
	procGdipGetPenDashOffset                  = modgdiplus.NewProc("GdipGetPenDashOffset")
	procGdipSetPenDashOffset                  = modgdiplus.NewProc("GdipSetPenDashOffset")
	procGdipGetPenDashCount                   = modgdiplus.NewProc("GdipGetPenDashCount")
	procGdipSetPenDashArray                   = modgdiplus.NewProc("GdipSetPenDashArray")
	procGdipGetPenDashArray                   = modgdiplus.NewProc("GdipGetPenDashArray")
	procGdipGetPenCompoundCount               = modgdiplus.NewProc("GdipGetPenCompoundCount")
	procGdipSetPenCompoundArray               = modgdiplus.NewProc("GdipSetPenCompoundArray")
	procGdipGetPenCompoundArray               = modgdiplus.NewProc("GdipGetPenCompoundArray")
	procGdipCloneBrush                        = modgdiplus.NewProc("GdipCloneBrush")
	procGdipDeleteBrush                       = modgdiplus.NewProc("GdipDeleteBrush")
	procGdipGetBrushType                      = modgdiplus.NewProc("GdipGetBrushType")
	procGdipCreateSolidFill                   = modgdiplus.NewProc("GdipCreateSolidFill")
	procGdipSetSolidFillColor                 = modgdiplus.NewProc("GdipSetSolidFillColor")
	procGdipGetSolidFillColor                 = modgdiplus.NewProc("GdipGetSolidFillColor")
	procGdipLoadImageFromFile                 = modgdiplus.NewProc("GdipLoadImageFromFile")
	procGdipSaveImageToFile                   = modgdiplus.NewProc("GdipSaveImageToFile")
	procGdipGetImageWidth                     = modgdiplus.NewProc("GdipGetImageWidth")
	procGdipGetImageHeight                    = modgdiplus.NewProc("GdipGetImageHeight")
	procGdipGetImageGraphicsContext           = modgdiplus.NewProc("GdipGetImageGraphicsContext")
	procGdipDisposeImage                      = modgdiplus.NewProc("GdipDisposeImage")
	procGdipCreateBitmapFromScan0             = modgdiplus.NewProc("GdipCreateBitmapFromScan0")
	procGdipCreateBitmapFromFile              = modgdiplus.NewProc("GdipCreateBitmapFromFile")
	procGdipCreateBitmapFromHBITMAP           = modgdiplus.NewProc("GdipCreateBitmapFromHBITMAP")
	procGdipCreateHBITMAPFromBitmap           = modgdiplus.NewProc("GdipCreateHBITMAPFromBitmap")
	procGdipCreateFontFromDC                  = modgdiplus.NewProc("GdipCreateFontFromDC")
	procGdipCreateFont                        = modgdiplus.NewProc("GdipCreateFont")
	procGdipDeleteFont                        = modgdiplus.NewProc("GdipDeleteFont")
	procGdipNewInstalledFontCollection        = modgdiplus.NewProc("GdipNewInstalledFontCollection")
	procGdipCreateFontFamilyFromName          = modgdiplus.NewProc("GdipCreateFontFamilyFromName")
	procGdipDeleteFontFamily                  = modgdiplus.NewProc("GdipDeleteFontFamily")
	procGdipCreateStringFormat                = modgdiplus.NewProc("GdipCreateStringFormat")
	procGdipDeleteStringFormat                = modgdiplus.NewProc("GdipDeleteStringFormat")
	procGdipStringFormatGetGenericTypographic = modgdiplus.NewProc("GdipStringFormatGetGenericTypographic")
	procGdipCreatePath                        = modgdiplus.NewProc("GdipCreatePath")
	procGdipDeletePath                        = modgdiplus.NewProc("GdipDeletePath")
	procGdipAddPathArc                        = modgdiplus.NewProc("GdipAddPathArc")
	procGdipAddPathArcI                       = modgdiplus.NewProc("GdipAddPathArcI")
	procGdipAddPathLine                       = modgdiplus.NewProc("GdipAddPathLine")
	procGdipAddPathLineI                      = modgdiplus.NewProc("GdipAddPathLineI")
	procGdipAddPathEllipse                    = modgdiplus.NewProc("GdipAddPathEllipse")
	procGdipAddPathEllipseI                   = modgdiplus.NewProc("GdipAddPathEllipseI")
	procGdipAddPathRectangle                  = modgdiplus.NewProc("GdipAddPathRectangle")
	procGdipAddPathBezier                     = modgdiplus.NewProc("GdipAddPathBezier")
	procGdipStartPathFigure                   = modgdiplus.NewProc("GdipStartPathFigure")
	procGdipClosePathFigure                   = modgdiplus.NewProc("GdipClosePathFigure")
	procGdipClosePathFigures                  = modgdiplus.NewProc("GdipClosePathFigures")
	procGdipSetPathFillMode                   = modgdiplus.NewProc("GdipSetPathFillMode")
	procGdipGetGenericFontFamilySerif         = modgdiplus.NewProc("GdipGetGenericFontFamilySerif")
	procGdipGetGenericFontFamilySansSerif     = modgdiplus.NewProc("GdipGetGenericFontFamilySansSerif")
	procGdipGetGenericFontFamilyMonospace     = modgdiplus.NewProc("GdipGetGenericFontFamilyMonospace")
	procGdipFlush                             = modgdiplus.NewProc("GdipFlush")
)

var (
	gdiplustoken uintptr
)

func GdiplusShutdown() {
	procGdiplusShutdown.Call(uintptr(unsafe.Pointer(&gdiplustoken)))
}

func GdiplusStartup(input *GdiplusStartupInput, output *GdiplusStartupOutput) GpStatus {
	var ret uintptr
	for i := 0; i != 255; i++ {
		ret, _, _ = procGdiplusStartup.Call(
			uintptr(unsafe.Pointer(&gdiplustoken)),
			uintptr(unsafe.Pointer(input)),
			uintptr(unsafe.Pointer(output)))
		if GpStatus(ret) == Ok {
			break
		}
		input.GdiplusVersion = uint32(i)
	}

	return GpStatus(ret)
}

// Graphics
func GdipCreateFromHDC(hdc uintptr, graphics **GpGraphics) GpStatus {
	ret, _, _ := procGdipCreateFromHDC.Call(hdc, uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHDC2(hdc HDC, hDevice HANDLE, graphics **GpGraphics) GpStatus {
	ret, _, _ := procGdipCreateFromHDC2.Call(
		uintptr(hdc),
		uintptr(hDevice),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWND(hwnd uintptr, graphics **GpGraphics) GpStatus {
	ret, _, _ := procGdipCreateFromHWND.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWNDICM(hwnd uintptr, graphics **GpGraphics) GpStatus {
	ret, _, _ := procGdipCreateFromHWNDICM.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipDeleteGraphics(graphics *GpGraphics) GpStatus {
	ret, _, _ := procGdipDeleteGraphics.Call(uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipFlush(graphics *GpGraphics, intention GpFlushIntention) GpStatus {
	ret, _, _ := procGdipFlush.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(intention),
	)
	return GpStatus(ret)
}

func GdipGetDC(graphics *GpGraphics, hdc *uintptr) GpStatus {
	ret, _, _ := procGdipGetDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(hdc)))
	return GpStatus(ret)
}

func GdipReleaseDC(graphics *GpGraphics, hdc uintptr) GpStatus {
	ret, _, _ := procGdipReleaseDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hdc))
	return GpStatus(ret)
}

func GdipSaveGraphics(graphics *GpGraphics, state *GpState) GpStatus {
	ret, _, _ := procGdipSaveGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(state)))
	return GpStatus(ret)
}

func GdipRestoreGraphics(graphics *GpGraphics, state GpState) GpStatus {
	ret, _, _ := procGdipRestoreGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(state))
	return GpStatus(ret)
}

func GdipTranslateWorldTransform(graphics *GpGraphics, x float32, y float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipTranslateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipScaleWorldTransform(graphics *GpGraphics, x float32, y float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipScaleWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipRotateWorldTransform(graphics *GpGraphics, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipRotateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(angle)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipSetClipRect(graphics *GpGraphics, x float32, y float32, width float32, height float32, mode GpCombineMode) GpStatus {
	ret, _, _ := procGdipSetClipRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(mode),
	)
	return GpStatus(ret)
}

func GdipSetCompositingMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := procGdipSetCompositingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetRenderingOrigin(graphics *GpGraphics, x, y int32) GpStatus {
	ret, _, _ := procGdipSetRenderingOrigin.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipSetCompositingQuality(graphics *GpGraphics, quality int32) GpStatus {
	ret, _, _ := procGdipSetCompositingQuality.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(quality))
	return GpStatus(ret)
}

func GdipSetInterpolationMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := procGdipSetInterpolationMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetPixelOffsetMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := procGdipSetPixelOffsetMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipGetSmoothingMode(graphics *GpGraphics, mode *GpSmoothingMode) GpStatus {
	ret, _, _ := procGdipGetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(mode)))
	return GpStatus(ret)
}
func GdipSetSmoothingMode(graphics *GpGraphics, mode GpSmoothingMode) GpStatus {
	ret, _, _ := procGdipSetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetTextRenderingHint(graphics *GpGraphics, hint int32) GpStatus {
	ret, _, _ := procGdipSetTextRenderingHint.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hint))
	return GpStatus(ret)
}

func GdipGraphicsClear(graphics *GpGraphics, color ARGB) GpStatus {
	ret, _, _ := procGdipGraphicsClear.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipDrawLine(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 float32) GpStatus {
	ret, _, _ := procGdipDrawLine.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipDrawLineI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := procGdipDrawLineI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipDrawArc(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipDrawArc.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawArcI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipDrawArcI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawBezier(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2, x3, y3, x4, y4 float32) GpStatus {
	ret, _, _ := procGdipDrawBezier.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)),
		uintptr(math.Float32bits(x3)),
		uintptr(math.Float32bits(y3)),
		uintptr(math.Float32bits(x4)),
		uintptr(math.Float32bits(y4)))
	return GpStatus(ret)
}

func GdipDrawBezierI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2, x3, y3, x4, y4 int32) GpStatus {
	ret, _, _ := procGdipDrawBezierI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2),
		uintptr(x3),
		uintptr(y3),
		uintptr(x4),
		uintptr(y4))
	return GpStatus(ret)
}

func GdipDrawRectangle(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipDrawRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawRectangleI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipDrawRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawEllipse(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipDrawEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawEllipseI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipDrawEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawPie(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipDrawPie.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawPieI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipDrawPieI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipDrawPolygon(graphics *GpGraphics, pen *GpPen, points *PointF, count int32) GpStatus {
	ret, _, _ := procGdipDrawPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPolygonI(graphics *GpGraphics, pen *GpPen, points *Point, count int32) GpStatus {
	ret, _, _ := procGdipDrawPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPath(graphics *GpGraphics, pen *GpPen, path *GpPath) GpStatus {
	ret, _, _ := procGdipDrawPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDrawImage(graphics *GpGraphics, image *GpImage, x, y float32) GpStatus {
	ret, _, _ := procGdipDrawImage.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)))
	return GpStatus(ret)
}

func GdipDrawImageI(graphics *GpGraphics, image *GpImage, x, y int32) GpStatus {
	ret, _, _ := procGdipDrawImageI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipDrawImageRect(graphics *GpGraphics, image *GpImage, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipDrawImageRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawImageRectI(graphics *GpGraphics, image *GpImage, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipDrawImageRectI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillRectangle(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipFillRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillRectangleI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipFillRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillEllipse(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipFillEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillEllipseI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipFillEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillPolygon(graphics *GpGraphics, brush *GpBrush, points *PointF, count int32, fillMode int32) GpStatus {
	ret, _, _ := procGdipFillPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillPolygonI(graphics *GpGraphics, brush *GpBrush, points *Point, count int32, fillMode int32) GpStatus {
	ret, _, _ := procGdipFillPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillPath(graphics *GpGraphics, brush *GpBrush, path *GpPath) GpStatus {
	ret, _, _ := procGdipFillPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipMeasureString(
	graphics *GpGraphics, text string, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) GpStatus {

	str, err := syscall.UTF16FromString(text)
	length := len(str)
	if err != nil {
		panic(err.Error())
	}

	ret, _, _ := procGdipMeasureString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(&str[0])),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(boundingBox)),
		uintptr(unsafe.Pointer(codepointsFitted)),
		uintptr(unsafe.Pointer(linesFilled)))
	return GpStatus(ret)
}

func GdipMeasureCharacterRanges(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, regionCount int32,
	regions **GpRegion) GpStatus {

	ret, _, _ := procGdipMeasureCharacterRanges.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(regionCount),
		uintptr(unsafe.Pointer(regions)))
	return GpStatus(ret)
}

func GdipDrawString(graphics *GpGraphics, text string, font *GpFont, layoutRect *RectF, stringFormat *GpStringFormat, brush *GpBrush) GpStatus {
	str, err := syscall.UTF16FromString(text)
	length := len(str)
	if err != nil {
		panic(err.Error())
	}

	ret, _, _ := procGdipDrawString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(&str[0])),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

// Pen
func GdipCreatePen1(color ARGB, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := procGdipCreatePen1.Call(
		uintptr(color),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))

	return GpStatus(ret)
}

func GdipCreatePen2(brush *GpBrush, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := procGdipCreatePen2.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipClonePen(pen *GpPen, clonepen **GpPen) GpStatus {
	ret, _, _ := procGdipClonePen.Call(uintptr(unsafe.Pointer(pen)), uintptr(unsafe.Pointer(clonepen)))
	return GpStatus(ret)
}

func GdipDeletePen(pen *GpPen) GpStatus {
	ret, _, _ := procGdipDeletePen.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipSetPenWidth(pen *GpPen, width float32) GpStatus {
	ret, _, _ := procGdipSetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(width)))
	return GpStatus(ret)
}

func GdipGetPenWidth(pen *GpPen, width *float32) GpStatus {
	var penWidth uint32
	ret, _, _ := procGdipGetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&penWidth)))
	*width = math.Float32frombits(penWidth)
	return GpStatus(ret)
}

func GdipSetPenLineCap197819(pen *GpPen, startCap, endCap GpLineCap, dashCap GpDashCap) GpStatus {
	ret, _, _ := procGdipSetPenLineCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap),
		uintptr(endCap),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipSetPenStartCap(pen *GpPen, startCap GpLineCap) GpStatus {
	ret, _, _ := procGdipSetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap))
	return GpStatus(ret)
}
func GdipSetPenEndCap(pen *GpPen, endCap GpLineCap) GpStatus {
	ret, _, _ := procGdipSetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(endCap))
	return GpStatus(ret)
}
func GdipSetPenDashCap197819(pen *GpPen, dashCap GpDashCap) GpStatus {
	ret, _, _ := procGdipSetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipGetPenStartCap(pen *GpPen, startCap *GpLineCap) GpStatus {
	ret, _, _ := procGdipGetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(startCap)))
	return GpStatus(ret)
}
func GdipGetPenEndCap(pen *GpPen, endCap *GpLineCap) GpStatus {
	ret, _, _ := procGdipGetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(endCap)))
	return GpStatus(ret)
}
func GdipGetPenDashCap197819(pen *GpPen, dashCap *GpDashCap) GpStatus {
	ret, _, _ := procGdipGetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashCap)))
	return GpStatus(ret)
}
func GdipSetPenLineJoin(pen *GpPen, lineJoin GpLineJoin) GpStatus {
	ret, _, _ := procGdipSetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(lineJoin))
	return GpStatus(ret)
}
func GdipGetPenLineJoin(pen *GpPen, lineJoin *GpLineJoin) GpStatus {
	ret, _, _ := procGdipGetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(lineJoin)))
	return GpStatus(ret)
}
func GdipSetPenCustomStartCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := procGdipSetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomStartCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := procGdipGetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenCustomEndCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := procGdipSetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomEndCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := procGdipGetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenMiterLimit(pen *GpPen, miterLimit float32) GpStatus {
	ret, _, _ := procGdipSetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(miterLimit)))
	return GpStatus(ret)
}
func GdipGetPenMiterLimit(pen *GpPen, miterLimit *float32) GpStatus {
	var iMiterLimit uint32
	ret, _, _ := procGdipGetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iMiterLimit)))
	*miterLimit = math.Float32frombits(iMiterLimit)
	return GpStatus(ret)
}
func GdipSetPenMode(pen *GpPen, penMode GpPenAlignment) GpStatus {
	ret, _, _ := procGdipSetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(penMode))
	return GpStatus(ret)
}
func GdipGetPenMode(pen *GpPen, penMode *GpPenAlignment) GpStatus {
	ret, _, _ := procGdipGetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penMode)))
	return GpStatus(ret)
}
func GdipSetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := procGdipSetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipGetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := procGdipGetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipResetPenTransform(pen *GpPen) GpStatus {
	ret, _, _ := procGdipResetPenTransform.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}
func GdipMultiplyPenTransform(pen *GpPen, matrix *GpMatrix, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipMultiplyPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipTranslatePenTransform(pen *GpPen, dx, dy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipTranslatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(dx)),
		uintptr(math.Float32bits(dy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipScalePenTransform(pen *GpPen, sx, sy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipScalePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(sx)),
		uintptr(math.Float32bits(sy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipRotatePenTransform(pen *GpPen, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := procGdipRotatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(angle)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipSetPenColor(pen *GpPen, argb ARGB) GpStatus {
	ret, _, _ := procGdipSetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(argb))
	return GpStatus(ret)
}
func GdipGetPenColor(pen *GpPen, argb *ARGB) GpStatus {
	ret, _, _ := procGdipGetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(argb)))
	return GpStatus(ret)
}
func GdipSetPenBrushFill(pen *GpPen, brush *GpBrush) GpStatus {
	ret, _, _ := procGdipSetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenBrushFill(pen *GpPen, brush **GpBrush) GpStatus {
	ret, _, _ := procGdipGetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenFillType(pen *GpPen, penType *GpPenType) GpStatus {
	ret, _, _ := procGdipGetPenFillType.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penType)))
	return GpStatus(ret)
}
func GdipGetPenDashStyle(pen *GpPen, dashStyle *GpDashStyle) GpStatus {
	ret, _, _ := procGdipGetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashStyle)))
	return GpStatus(ret)
}

func GdipSetPenDashStyle(pen *GpPen, dashStyle GpDashStyle) GpStatus {
	ret, _, _ := procGdipSetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashStyle))
	return GpStatus(ret)
}
func GdipGetPenDashOffset(pen *GpPen, offset *float32) GpStatus {
	var iOffset uint32
	ret, _, _ := procGdipGetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iOffset)))
	*offset = math.Float32frombits(iOffset)
	return GpStatus(ret)
}
func GdipSetPenDashOffset(pen *GpPen, offset float32) GpStatus {
	ret, _, _ := procGdipSetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(offset)))
	return GpStatus(ret)
}
func GdipGetPenDashCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := procGdipGetPenDashCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}
func GdipSetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := procGdipSetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}
func GdipGetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := procGdipGetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := procGdipGetPenCompoundCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}

func GdipSetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := procGdipSetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := procGdipGetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

// Brush

func GdipCloneBrush(brush *GpBrush, clone **GpBrush) GpStatus {
	ret, _, _ := procGdipCloneBrush.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(clone)))
	return GpStatus(ret)
}

func GdipDeleteBrush(brush *GpBrush) GpStatus {
	ret, _, _ := procGdipDeleteBrush.Call(uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipGetBrushType(brush *GpBrush, brushType *GpBrushType) GpStatus {
	ret, _, _ := procGdipGetBrushType.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(brushType)))
	return GpStatus(ret)
}

// Solid Brush

// func GdipCreateSolidFill(color ARGB, brush **GpSolidFill) GpStatus {
// 	ret, _, _ := procGdipCreateSolidFill.Call(
// 		uintptr(color),
// 		uintptr(unsafe.Pointer(brush)))
// 	return GpStatus(ret)
// }

func GdipCreateSolidFill(color ARGB, brush **GpBrush) GpStatus {
	ret, _, _ := procGdipCreateSolidFill.Call(
		uintptr(color),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipSetSolidFillColor(brush *GpBrush, color ARGB) GpStatus {
	ret, _, _ := procGdipSetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipGetSolidFillColor(brush *GpBrush, color *ARGB) GpStatus {
	ret, _, _ := procGdipGetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(color)))
	return GpStatus(ret)
}

// Font
func GdipCreateFontFromDC(hdc uintptr, font **GpFont) GpStatus {
	ret, _, _ := procGdipCreateFontFromDC.Call(
		hdc,
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipCreateFont(fontFamily *GpFontFamily, emSize float32, style int32, unit GpUnit, font **GpFont) GpStatus {
	ret, _, _ := procGdipCreateFont.Call(
		uintptr(unsafe.Pointer(fontFamily)),
		uintptr(math.Float32bits(emSize)),
		uintptr(style),
		uintptr(unit),
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipDeleteFont(font *GpFont) GpStatus {
	ret, _, _ := procGdipDeleteFont.Call(uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipNewInstalledFontCollection(fontCollection **GpFontCollection) GpStatus {
	ret, _, _ := procGdipNewInstalledFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return GpStatus(ret)
}

func GdipCreateFontFamilyFromName(familyName string, fontCollection *GpFontCollection, fontFamily **GpFontFamily) GpStatus {
	cname, err := syscall.UTF16PtrFromString(familyName)
	if err != nil {
		panic(err)
	}
	ret, _, _ := procGdipCreateFontFamilyFromName.Call(
		uintptr(unsafe.Pointer(cname)),
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

func GdipDeleteFontFamily(fontFamily *GpFontFamily) GpStatus {
	ret, _, _ := procGdipDeleteFontFamily.Call(uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

// StringFormat

func GdipCreateStringFormat(formatAttributes int32, language uint16, format **GpStringFormat) GpStatus {
	ret, _, _ := procGdipCreateStringFormat.Call(
		uintptr(formatAttributes),
		uintptr(language),
		uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipStringFormatGetGenericTypographic(format **GpStringFormat) GpStatus {
	ret, _, _ := procGdipStringFormatGetGenericTypographic.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipDeleteStringFormat(format *GpStringFormat) GpStatus {
	ret, _, _ := procGdipDeleteStringFormat.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

// Path

func GdipCreatePath(brushMode GpFillMode, path **GpPath) GpStatus {
	ret, _, _ := procGdipCreatePath.Call(uintptr(brushMode), uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDeletePath(path *GpPath) GpStatus {
	ret, _, _ := procGdipDeletePath.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipAddPathArc(path *GpPath, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipAddPathArc.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipAddPathArcI(path *GpPath, x, y, width, height int32, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := procGdipAddPathArcI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return GpStatus(ret)
}

func GdipAddPathLine(path *GpPath, x1, y1, x2, y2 float32) GpStatus {
	ret, _, _ := procGdipAddPathLine.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipAddPathLineI(path *GpPath, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := procGdipAddPathLineI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipAddPathEllipseI(path *GpPath, x, y, width, height int32) GpStatus {
	ret, _, _ := procGdipAddPathEllipseI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
	)
	return GpStatus(ret)
}

func GdipAddPathEllipse(path *GpPath, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipAddPathEllipse.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
	)
	return GpStatus(ret)
}

func GdipAddPathRectangle(path *GpPath, x, y, width, height float32) GpStatus {
	ret, _, _ := procGdipAddPathRectangle.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
	)
	return GpStatus(ret)
}

func GdipAddPathBezier(path *GpPath, x1, y1, x2, y2, x3, y3, x4, y4 float32) GpStatus {
	ret, _, _ := procGdipAddPathBezier.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)),
		uintptr(math.Float32bits(x3)),
		uintptr(math.Float32bits(y3)),
		uintptr(math.Float32bits(x4)),
		uintptr(math.Float32bits(y4)),
	)
	return GpStatus(ret)
}

func GdipStartPathFigure(path *GpPath) GpStatus {
	ret, _, _ := procGdipStartPathFigure.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipClosePathFigure(path *GpPath) GpStatus {
	ret, _, _ := procGdipClosePathFigure.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipClosePathFigures(path *GpPath) GpStatus {
	ret, _, _ := procGdipClosePathFigures.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipSetPathFillMode(path *GpPath, fillmode GpFillMode) GpStatus {
	ret, _, _ := procGdipSetPathFillMode.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(fillmode),
	)
	return GpStatus(ret)
}

// Image

func GdipGetImageGraphicsContext(image *GpImage, graphics **GpGraphics) GpStatus {
	ret, _, _ := procGdipGetImageGraphicsContext.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipLoadImageFromFile(filename *uint16, image **GpImage) GpStatus {
	ret, _, _ := procGdipLoadImageFromFile.Call(
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(image)))
	return GpStatus(ret)
}

// func GdipSaveImageToFile(image *GpBitmap, filename *uint16, clsidEncoder *ole.GUID, encoderParams *EncoderParameters) GpStatus {
// 	ret, _, _ := procGdipSaveImageToFile.Call(uintptr(unsafe.Pointer(image)),
// 		uintptr(unsafe.Pointer(filename)), uintptr(unsafe.Pointer(clsidEncoder)),
// 		uintptr(unsafe.Pointer(encoderParams)))
// 	return GpStatus(ret)
// }

func GdipGetImageWidth(image *GpImage, width *uint32) GpStatus {
	ret, _, _ := procGdipGetImageWidth.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(width)))
	return GpStatus(ret)
}

func GdipGetImageHeight(image *GpImage, height *uint32) GpStatus {
	ret, _, _ := procGdipGetImageHeight.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(height)))
	return GpStatus(ret)
}

func GdipDisposeImage(image *GpImage) GpStatus {
	ret, _, _ := syscall.Syscall(procGdipDisposeImage.Addr(), 1,
		uintptr(unsafe.Pointer(image)),
		0,
		0)

	return GpStatus(ret)
}

// Bitmap

func GdipCreateBitmapFromFile(filename *uint16, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(procGdipCreateBitmapFromFile.Addr(), 2,
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(bitmap)),
		0)

	return GpStatus(ret)
}

func GdipCreateBitmapFromHBITMAP(hbm HBITMAP, hpal HPALETTE, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(procGdipCreateBitmapFromHBITMAP.Addr(), 3,
		uintptr(hbm),
		uintptr(hpal),
		uintptr(unsafe.Pointer(bitmap)))

	return GpStatus(ret)
}

func GdipCreateHBITMAPFromBitmap(bitmap *GpBitmap, hbmReturn *HBITMAP, background ARGB) GpStatus {
	ret, _, _ := syscall.Syscall(procGdipCreateHBITMAPFromBitmap.Addr(), 3,
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(hbmReturn)),
		uintptr(background))

	return GpStatus(ret)
}

func GdipCreateBitmapFromScan0(width, height, stride int32, format PixelFormat, scan0 *byte, bitmap **GpBitmap) GpStatus {
	ret, _, _ := procGdipCreateBitmapFromScan0.Call(
		uintptr(width),
		uintptr(height),
		uintptr(stride),
		uintptr(format),
		uintptr(unsafe.Pointer(scan0)),
		uintptr(unsafe.Pointer(bitmap)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySerif(family **GpFontFamily) GpStatus {
	ret, _, _ := procGdipGetGenericFontFamilySerif.Call(
		uintptr(unsafe.Pointer(family)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySansSerif(family **GpFontFamily) GpStatus {
	ret, _, _ := procGdipGetGenericFontFamilySansSerif.Call(
		uintptr(unsafe.Pointer(family)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilyMonospace(family **GpFontFamily) GpStatus {
	ret, _, _ := procGdipGetGenericFontFamilyMonospace.Call(
		uintptr(unsafe.Pointer(family)))
	return GpStatus(ret)
}

/*
func SavePNG(fileName string, newBMP win.HBITMAP) error {
	// HBITMAP
	var bmp *win.GpBitmap
	if win.GdipCreateBitmapFromHBITMAP(newBMP, 0, &bmp) != 0 {
		return fmt.Errorf("failed to create HBITMAP")
	}
	defer win.GdipDisposeImage((*GpImage)(bmp))
	clsid, err := ole.CLSIDFromString("{557CF406-1A04-11D3-9A73-0000F81EF32E}")
	if err != nil {
		return err
	}
	fname, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return err
	}
	if GdipSaveImageToFile(bmp, fname, clsid, nil) != 0 {
		return fmt.Errorf("failed to call PNG encoder")
	}
	return nil
}
*/
