// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	modgdiplus                            = syscall.NewLazyDLL("gdiplus.dll")
	gdiplusShutdown                       = modgdiplus.NewProc("GdiplusShutdown")
	gdiplusStartup                        = modgdiplus.NewProc("GdiplusStartup")
	gdipCreateFromHDC                     = modgdiplus.NewProc("GdipCreateFromHDC")
	gdipCreateFromHDC2                    = modgdiplus.NewProc("GdipCreateFromHDC2")
	gdipCreateFromHWND                    = modgdiplus.NewProc("GdipCreateFromHWND")
	gdipCreateFromHWNDICM                 = modgdiplus.NewProc("GdipCreateFromHWNDICM")
	gdipDeleteGraphics                    = modgdiplus.NewProc("GdipDeleteGraphics")
	gdipGetDC                             = modgdiplus.NewProc("GdipGetDC")
	gdipReleaseDC                         = modgdiplus.NewProc("GdipReleaseDC")
	gdipSaveGraphics                      = modgdiplus.NewProc("GdipSaveGraphics")
	gdipRestoreGraphics                   = modgdiplus.NewProc("GdipRestoreGraphics")
	gdipTranslateWorldTransform           = modgdiplus.NewProc("GdipTranslateWorldTransform")
	gdipScaleWorldTransform               = modgdiplus.NewProc("GdipScaleWorldTransform")
	gdipRotateWorldTransform              = modgdiplus.NewProc("GdipRotateWorldTransform")
	gdipSetClipRect                       = modgdiplus.NewProc("GdipSetClipRect")
	gdipSetCompositingMode                = modgdiplus.NewProc("GdipSetCompositingMode")
	gdipSetRenderingOrigin                = modgdiplus.NewProc("GdipSetRenderingOrigin")
	gdipSetCompositingQuality             = modgdiplus.NewProc("GdipSetCompositingQuality")
	gdipSetSmoothingMode                  = modgdiplus.NewProc("GdipSetSmoothingMode")
	gdipGetSmoothingMode                  = modgdiplus.NewProc("GdipGetSmoothingMode")
	gdipSetPixelOffsetMode                = modgdiplus.NewProc("GdipSetPixelOffsetMode")
	gdipSetInterpolationMode              = modgdiplus.NewProc("GdipSetInterpolationMode")
	gdipSetTextRenderingHint              = modgdiplus.NewProc("GdipSetTextRenderingHint")
	gdipGraphicsClear                     = modgdiplus.NewProc("GdipGraphicsClear")
	gdipDrawLine                          = modgdiplus.NewProc("GdipDrawLine")
	gdipDrawLineI                         = modgdiplus.NewProc("GdipDrawLineI")
	gdipDrawArc                           = modgdiplus.NewProc("GdipDrawArc")
	gdipDrawArcI                          = modgdiplus.NewProc("GdipDrawArcI")
	gdipDrawBezier                        = modgdiplus.NewProc("GdipDrawBezier")
	gdipDrawBezierI                       = modgdiplus.NewProc("GdipDrawBezierI")
	gdipDrawRectangle                     = modgdiplus.NewProc("GdipDrawRectangle")
	gdipDrawRectangleI                    = modgdiplus.NewProc("GdipDrawRectangleI")
	gdipDrawEllipse                       = modgdiplus.NewProc("GdipDrawEllipse")
	gdipDrawEllipseI                      = modgdiplus.NewProc("GdipDrawEllipseI")
	gdipDrawPie                           = modgdiplus.NewProc("GdipDrawPie")
	gdipDrawPieI                          = modgdiplus.NewProc("GdipDrawPieI")
	gdipDrawPolygonI                      = modgdiplus.NewProc("GdipDrawPolygonI")
	gdipDrawPolygon                       = modgdiplus.NewProc("GdipDrawPolygon")
	gdipDrawPath                          = modgdiplus.NewProc("GdipDrawPath")
	gdipDrawString                        = modgdiplus.NewProc("GdipDrawString")
	gdipDrawImage                         = modgdiplus.NewProc("GdipDrawImage")
	gdipDrawImageI                        = modgdiplus.NewProc("GdipDrawImageI")
	gdipDrawImageRect                     = modgdiplus.NewProc("GdipDrawImageRect")
	gdipDrawImageRectI                    = modgdiplus.NewProc("GdipDrawImageRectI")
	gdipFillRectangle                     = modgdiplus.NewProc("GdipFillRectangle")
	gdipFillRectangleI                    = modgdiplus.NewProc("GdipFillRectangleI")
	gdipFillPolygon                       = modgdiplus.NewProc("GdipFillPolygon")
	gdipFillPolygonI                      = modgdiplus.NewProc("GdipFillPolygonI")
	gdipFillPath                          = modgdiplus.NewProc("GdipFillPath")
	gdipFillEllipse                       = modgdiplus.NewProc("GdipFillEllipse")
	gdipFillEllipseI                      = modgdiplus.NewProc("GdipFillEllipseI")
	gdipMeasureString                     = modgdiplus.NewProc("GdipMeasureString")
	gdipMeasureCharacterRanges            = modgdiplus.NewProc("GdipMeasureCharacterRanges")
	gdipCreatePen1                        = modgdiplus.NewProc("GdipCreatePen1")
	gdipCreatePen2                        = modgdiplus.NewProc("GdipCreatePen2")
	gdipClonePen                          = modgdiplus.NewProc("GdipClonePen")
	gdipDeletePen                         = modgdiplus.NewProc("GdipDeletePen")
	gdipSetPenWidth                       = modgdiplus.NewProc("GdipSetPenWidth")
	gdipGetPenWidth                       = modgdiplus.NewProc("GdipGetPenWidth")
	gdipSetPenLineCap197819               = modgdiplus.NewProc("GdipSetPenLineCap197819")
	gdipSetPenStartCap                    = modgdiplus.NewProc("GdipSetPenStartCap")
	gdipSetPenEndCap                      = modgdiplus.NewProc("GdipSetPenEndCap")
	gdipSetPenDashCap197819               = modgdiplus.NewProc("GdipSetPenDashCap197819")
	gdipGetPenStartCap                    = modgdiplus.NewProc("GdipGetPenStartCap")
	gdipGetPenEndCap                      = modgdiplus.NewProc("GdipGetPenEndCap")
	gdipGetPenDashCap197819               = modgdiplus.NewProc("GdipGetPenDashCap197819")
	gdipSetPenLineJoin                    = modgdiplus.NewProc("GdipSetPenLineJoin")
	gdipGetPenLineJoin                    = modgdiplus.NewProc("GdipGetPenLineJoin")
	gdipSetPenCustomStartCap              = modgdiplus.NewProc("GdipSetPenCustomStartCap")
	gdipGetPenCustomStartCap              = modgdiplus.NewProc("GdipGetPenCustomStartCap")
	gdipSetPenCustomEndCap                = modgdiplus.NewProc("GdipSetPenCustomEndCap")
	gdipGetPenCustomEndCap                = modgdiplus.NewProc("GdipGetPenCustomEndCap")
	gdipSetPenMiterLimit                  = modgdiplus.NewProc("GdipSetPenMiterLimit")
	gdipGetPenMiterLimit                  = modgdiplus.NewProc("GdipGetPenMiterLimit")
	gdipSetPenMode                        = modgdiplus.NewProc("GdipSetPenMode")
	gdipGetPenMode                        = modgdiplus.NewProc("GdipGetPenMode")
	gdipSetPenTransform                   = modgdiplus.NewProc("GdipSetPenTransform")
	gdipGetPenTransform                   = modgdiplus.NewProc("GdipGetPenTransform")
	gdipResetPenTransform                 = modgdiplus.NewProc("GdipResetPenTransform")
	gdipMultiplyPenTransform              = modgdiplus.NewProc("GdipMultiplyPenTransform")
	gdipTranslatePenTransform             = modgdiplus.NewProc("GdipTranslatePenTransform")
	gdipScalePenTransform                 = modgdiplus.NewProc("GdipScalePenTransform")
	gdipRotatePenTransform                = modgdiplus.NewProc("GdipRotatePenTransform")
	gdipSetPenColor                       = modgdiplus.NewProc("GdipSetPenColor")
	gdipGetPenColor                       = modgdiplus.NewProc("GdipGetPenColor")
	gdipSetPenBrushFill                   = modgdiplus.NewProc("GdipSetPenBrushFill")
	gdipGetPenBrushFill                   = modgdiplus.NewProc("GdipGetPenBrushFill")
	gdipGetPenFillType                    = modgdiplus.NewProc("GdipGetPenFillType")
	gdipGetPenDashStyle                   = modgdiplus.NewProc("GdipGetPenDashStyle")
	gdipSetPenDashStyle                   = modgdiplus.NewProc("GdipSetPenDashStyle")
	gdipGetPenDashOffset                  = modgdiplus.NewProc("GdipGetPenDashOffset")
	gdipSetPenDashOffset                  = modgdiplus.NewProc("GdipSetPenDashOffset")
	gdipGetPenDashCount                   = modgdiplus.NewProc("GdipGetPenDashCount")
	gdipSetPenDashArray                   = modgdiplus.NewProc("GdipSetPenDashArray")
	gdipGetPenDashArray                   = modgdiplus.NewProc("GdipGetPenDashArray")
	gdipGetPenCompoundCount               = modgdiplus.NewProc("GdipGetPenCompoundCount")
	gdipSetPenCompoundArray               = modgdiplus.NewProc("GdipSetPenCompoundArray")
	gdipGetPenCompoundArray               = modgdiplus.NewProc("GdipGetPenCompoundArray")
	gdipCloneBrush                        = modgdiplus.NewProc("GdipCloneBrush")
	gdipDeleteBrush                       = modgdiplus.NewProc("GdipDeleteBrush")
	gdipGetBrushType                      = modgdiplus.NewProc("GdipGetBrushType")
	gdipCreateSolidFill                   = modgdiplus.NewProc("GdipCreateSolidFill")
	gdipSetSolidFillColor                 = modgdiplus.NewProc("GdipSetSolidFillColor")
	gdipGetSolidFillColor                 = modgdiplus.NewProc("GdipGetSolidFillColor")
	gdipLoadImageFromFile                 = modgdiplus.NewProc("GdipLoadImageFromFile")
	gdipSaveImageToFile                   = modgdiplus.NewProc("GdipSaveImageToFile")
	gdipGetImageWidth                     = modgdiplus.NewProc("GdipGetImageWidth")
	gdipGetImageHeight                    = modgdiplus.NewProc("GdipGetImageHeight")
	gdipGetImageGraphicsContext           = modgdiplus.NewProc("GdipGetImageGraphicsContext")
	gdipDisposeImage                      = modgdiplus.NewProc("GdipDisposeImage")
	gdipCreateBitmapFromScan0             = modgdiplus.NewProc("GdipCreateBitmapFromScan0")
	gdipCreateBitmapFromFile              = modgdiplus.NewProc("GdipCreateBitmapFromFile")
	gdipCreateBitmapFromHBITMAP           = modgdiplus.NewProc("GdipCreateBitmapFromHBITMAP")
	gdipCreateHBITMAPFromBitmap           = modgdiplus.NewProc("GdipCreateHBITMAPFromBitmap")
	gdipCreateFontFromDC                  = modgdiplus.NewProc("GdipCreateFontFromDC")
	gdipCreateFont                        = modgdiplus.NewProc("GdipCreateFont")
	gdipDeleteFont                        = modgdiplus.NewProc("GdipDeleteFont")
	gdipNewInstalledFontCollection        = modgdiplus.NewProc("GdipNewInstalledFontCollection")
	gdipCreateFontFamilyFromName          = modgdiplus.NewProc("GdipCreateFontFamilyFromName")
	gdipDeleteFontFamily                  = modgdiplus.NewProc("GdipDeleteFontFamily")
	gdipCreateStringFormat                = modgdiplus.NewProc("GdipCreateStringFormat")
	gdipDeleteStringFormat                = modgdiplus.NewProc("GdipDeleteStringFormat")
	gdipStringFormatGetGenericTypographic = modgdiplus.NewProc("GdipStringFormatGetGenericTypographic")
	gdipCreatePath                        = modgdiplus.NewProc("GdipCreatePath")
	gdipDeletePath                        = modgdiplus.NewProc("GdipDeletePath")
	gdipAddPathArc                        = modgdiplus.NewProc("GdipAddPathArc")
	gdipAddPathArcI                       = modgdiplus.NewProc("GdipAddPathArcI")
	gdipAddPathLine                       = modgdiplus.NewProc("GdipAddPathLine")
	gdipAddPathLineI                      = modgdiplus.NewProc("GdipAddPathLineI")
	gdipClosePathFigure                   = modgdiplus.NewProc("GdipClosePathFigure")
	gdipClosePathFigures                  = modgdiplus.NewProc("GdipClosePathFigures")
	gdipGetGenericFontFamilySerif         = modgdiplus.NewProc("GdipGetGenericFontFamilySerif")
	gdipGetGenericFontFamilySansSerif     = modgdiplus.NewProc("GdipGetGenericFontFamilySansSerif")
	gdipGetGenericFontFamilyMonospace     = modgdiplus.NewProc("GdipGetGenericFontFamilyMonospace")
	gdipFlush                             = modgdiplus.NewProc("GdipFlush")
)

var (
	gdiplustoken uintptr
)

func GdiplusShutdown() {
	gdiplusShutdown.Call(uintptr(unsafe.Pointer(&gdiplustoken)))
}

func GdiplusStartup(input *GdiplusStartupInput, output *GdiplusStartupOutput) GpStatus {
	var ret uintptr
	for i := 0; i != 255; i++ {
		ret, _, _ = gdiplusStartup.Call(
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
	ret, _, _ := gdipCreateFromHDC.Call(hdc, uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHDC2(hdc HDC, hDevice HANDLE, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHDC2.Call(
		uintptr(hdc),
		uintptr(hDevice),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWND(hwnd uintptr, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHWND.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipCreateFromHWNDICM(hwnd uintptr, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipCreateFromHWNDICM.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipDeleteGraphics(graphics *GpGraphics) GpStatus {
	ret, _, _ := gdipDeleteGraphics.Call(uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipFlush(graphics *GpGraphics, intention GpFlushIntention) GpStatus {
	ret, _, _ := gdipFlush.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(intention),
	)
	return GpStatus(ret)
}

func GdipGetDC(graphics *GpGraphics, hdc *uintptr) GpStatus {
	ret, _, _ := gdipGetDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(hdc)))
	return GpStatus(ret)
}

func GdipReleaseDC(graphics *GpGraphics, hdc uintptr) GpStatus {
	ret, _, _ := gdipReleaseDC.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hdc))
	return GpStatus(ret)
}

func GdipSaveGraphics(graphics *GpGraphics, state *GpState) GpStatus {
	ret, _, _ := gdipSaveGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(state)))
	return GpStatus(ret)
}

func GdipRestoreGraphics(graphics *GpGraphics, state GpState) GpStatus {
	ret, _, _ := gdipRestoreGraphics.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(state))
	return GpStatus(ret)
}

func GdipTranslateWorldTransform(graphics *GpGraphics, x float32, y float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipTranslateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipScaleWorldTransform(graphics *GpGraphics, x float32, y float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipScaleWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipRotateWorldTransform(graphics *GpGraphics, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipRotateWorldTransform.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(angle)),
		uintptr(order),
	)
	return GpStatus(ret)
}

func GdipSetClipRect(graphics *GpGraphics, x float32, y float32, width float32, height float32, mode GpCombineMode) GpStatus {
	ret, _, _ := gdipSetClipRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(mode),
	)
	return GpStatus(ret)
}

// func GdipMeasureString(graphics *GpGraphics, text string, font *GpFont, layout *RectF, size *RectF) GpStatus {
// 	str, err := syscall.UTF16FromString(text)
// 	l := len(str)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	ret, _, _ := gdipMeasureString.Call(
// 		uintptr(unsafe.Pointer(graphics)),
// 		uintptr(unsafe.Pointer(&str[0])),
// 		uintptr(l),
// 		uintptr(unsafe.Pointer(font)),
// 		uintptr(0),
// 		uintptr(unsafe.Pointer(layout)),
// 		uintptr(0),
// 		uintptr(unsafe.Pointer(size)),
// 		uintptr(0),
// 		uintptr(0),
// 	)
// 	return GpStatus(ret)
// }

func GdipSetCompositingMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetCompositingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetRenderingOrigin(graphics *GpGraphics, x, y int32) GpStatus {
	ret, _, _ := gdipSetRenderingOrigin.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipSetCompositingQuality(graphics *GpGraphics, quality int32) GpStatus {
	ret, _, _ := gdipSetCompositingQuality.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(quality))
	return GpStatus(ret)
}

func GdipSetInterpolationMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetInterpolationMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetPixelOffsetMode(graphics *GpGraphics, mode int32) GpStatus {
	ret, _, _ := gdipSetPixelOffsetMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipGetSmoothingMode(graphics *GpGraphics, mode *GpSmoothingMode) GpStatus {
	ret, _, _ := gdipGetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(mode)))
	return GpStatus(ret)
}
func GdipSetSmoothingMode(graphics *GpGraphics, mode GpSmoothingMode) GpStatus {
	ret, _, _ := gdipSetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return GpStatus(ret)
}

func GdipSetTextRenderingHint(graphics *GpGraphics, hint int32) GpStatus {
	ret, _, _ := gdipSetTextRenderingHint.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hint))
	return GpStatus(ret)
}

func GdipGraphicsClear(graphics *GpGraphics, color ARGB) GpStatus {
	ret, _, _ := gdipGraphicsClear.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipDrawLine(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 float32) GpStatus {
	ret, _, _ := gdipDrawLine.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipDrawLineI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := gdipDrawLineI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipDrawArc(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawArc.Call(
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
	ret, _, _ := gdipDrawArcI.Call(
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
	ret, _, _ := gdipDrawBezier.Call(
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
	ret, _, _ := gdipDrawBezierI.Call(
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
	ret, _, _ := gdipDrawRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawRectangleI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawEllipse(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipDrawEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawEllipseI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipDrawPie(graphics *GpGraphics, pen *GpPen, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipDrawPie.Call(
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
	ret, _, _ := gdipDrawPieI.Call(
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
	ret, _, _ := gdipDrawPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPolygonI(graphics *GpGraphics, pen *GpPen, points *Point, count int32) GpStatus {
	ret, _, _ := gdipDrawPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipDrawPath(graphics *GpGraphics, pen *GpPen, path *GpPath) GpStatus {
	ret, _, _ := gdipDrawPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDrawString(graphics *GpGraphics, text *uint16, length int32, font *GpFont, layoutRect *RectF, stringFormat *GpStringFormat, brush *GpBrush) GpStatus {
	ret, _, _ := gdipDrawString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipDrawImage(graphics *GpGraphics, image *GpImage, x, y float32) GpStatus {
	ret, _, _ := gdipDrawImage.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)))
	return GpStatus(ret)
}

func GdipDrawImageI(graphics *GpGraphics, image *GpImage, x, y int32) GpStatus {
	ret, _, _ := gdipDrawImageI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y))
	return GpStatus(ret)
}

func GdipDrawImageRect(graphics *GpGraphics, image *GpImage, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipDrawImageRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipDrawImageRectI(graphics *GpGraphics, image *GpImage, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipDrawImageRectI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillRectangle(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipFillRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillRectangleI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipFillRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillEllipse(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) GpStatus {
	ret, _, _ := gdipFillEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return GpStatus(ret)
}

func GdipFillEllipseI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) GpStatus {
	ret, _, _ := gdipFillEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return GpStatus(ret)
}

func GdipFillPolygon(graphics *GpGraphics, brush *GpBrush, points *PointF, count int32, fillMode int32) GpStatus {
	ret, _, _ := gdipFillPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillPolygonI(graphics *GpGraphics, brush *GpBrush, points *Point, count int32, fillMode int32) GpStatus {
	ret, _, _ := gdipFillPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return GpStatus(ret)
}

func GdipFillPath(graphics *GpGraphics, brush *GpBrush, path *GpPath) GpStatus {
	ret, _, _ := gdipFillPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipMeasureString(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) GpStatus {

	ret, _, _ := gdipMeasureString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
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

	ret, _, _ := gdipMeasureCharacterRanges.Call(
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

// Pen
func GdipCreatePen1(color ARGB, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := gdipCreatePen1.Call(
		uintptr(color),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))

	return GpStatus(ret)
}

func GdipCreatePen2(brush *GpBrush, width float32, unit GpUnit, pen **GpPen) GpStatus {
	ret, _, _ := gdipCreatePen2.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipClonePen(pen *GpPen, clonepen **GpPen) GpStatus {
	ret, _, _ := gdipClonePen.Call(uintptr(unsafe.Pointer(pen)), uintptr(unsafe.Pointer(clonepen)))
	return GpStatus(ret)
}

func GdipDeletePen(pen *GpPen) GpStatus {
	ret, _, _ := gdipDeletePen.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}

func GdipSetPenWidth(pen *GpPen, width float32) GpStatus {
	ret, _, _ := gdipSetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(width)))
	return GpStatus(ret)
}

func GdipGetPenWidth(pen *GpPen, width *float32) GpStatus {
	var penWidth uint32
	ret, _, _ := gdipGetPenWidth.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&penWidth)))
	*width = math.Float32frombits(penWidth)
	return GpStatus(ret)
}

func GdipSetPenLineCap197819(pen *GpPen, startCap, endCap GpLineCap, dashCap GpDashCap) GpStatus {
	ret, _, _ := gdipSetPenLineCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap),
		uintptr(endCap),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipSetPenStartCap(pen *GpPen, startCap GpLineCap) GpStatus {
	ret, _, _ := gdipSetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(startCap))
	return GpStatus(ret)
}
func GdipSetPenEndCap(pen *GpPen, endCap GpLineCap) GpStatus {
	ret, _, _ := gdipSetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(endCap))
	return GpStatus(ret)
}
func GdipSetPenDashCap197819(pen *GpPen, dashCap GpDashCap) GpStatus {
	ret, _, _ := gdipSetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashCap))
	return GpStatus(ret)
}
func GdipGetPenStartCap(pen *GpPen, startCap *GpLineCap) GpStatus {
	ret, _, _ := gdipGetPenStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(startCap)))
	return GpStatus(ret)
}
func GdipGetPenEndCap(pen *GpPen, endCap *GpLineCap) GpStatus {
	ret, _, _ := gdipGetPenEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(endCap)))
	return GpStatus(ret)
}
func GdipGetPenDashCap197819(pen *GpPen, dashCap *GpDashCap) GpStatus {
	ret, _, _ := gdipGetPenDashCap197819.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashCap)))
	return GpStatus(ret)
}
func GdipSetPenLineJoin(pen *GpPen, lineJoin GpLineJoin) GpStatus {
	ret, _, _ := gdipSetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(lineJoin))
	return GpStatus(ret)
}
func GdipGetPenLineJoin(pen *GpPen, lineJoin *GpLineJoin) GpStatus {
	ret, _, _ := gdipGetPenLineJoin.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(lineJoin)))
	return GpStatus(ret)
}
func GdipSetPenCustomStartCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := gdipSetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomStartCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := gdipGetPenCustomStartCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenCustomEndCap(pen *GpPen, customCap *GpCustomLineCap) GpStatus {
	ret, _, _ := gdipSetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipGetPenCustomEndCap(pen *GpPen, customCap **GpCustomLineCap) GpStatus {
	ret, _, _ := gdipGetPenCustomEndCap.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(customCap)))
	return GpStatus(ret)
}
func GdipSetPenMiterLimit(pen *GpPen, miterLimit float32) GpStatus {
	ret, _, _ := gdipSetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(miterLimit)))
	return GpStatus(ret)
}
func GdipGetPenMiterLimit(pen *GpPen, miterLimit *float32) GpStatus {
	var iMiterLimit uint32
	ret, _, _ := gdipGetPenMiterLimit.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iMiterLimit)))
	*miterLimit = math.Float32frombits(iMiterLimit)
	return GpStatus(ret)
}
func GdipSetPenMode(pen *GpPen, penMode GpPenAlignment) GpStatus {
	ret, _, _ := gdipSetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(penMode))
	return GpStatus(ret)
}
func GdipGetPenMode(pen *GpPen, penMode *GpPenAlignment) GpStatus {
	ret, _, _ := gdipGetPenMode.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penMode)))
	return GpStatus(ret)
}
func GdipSetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := gdipSetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipGetPenTransform(pen *GpPen, matrix *GpMatrix) GpStatus {
	ret, _, _ := gdipGetPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)))
	return GpStatus(ret)
}
func GdipResetPenTransform(pen *GpPen) GpStatus {
	ret, _, _ := gdipResetPenTransform.Call(uintptr(unsafe.Pointer(pen)))
	return GpStatus(ret)
}
func GdipMultiplyPenTransform(pen *GpPen, matrix *GpMatrix, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipMultiplyPenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(matrix)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipTranslatePenTransform(pen *GpPen, dx, dy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipTranslatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(dx)),
		uintptr(math.Float32bits(dy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipScalePenTransform(pen *GpPen, sx, sy float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipScalePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(sx)),
		uintptr(math.Float32bits(sy)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipRotatePenTransform(pen *GpPen, angle float32, order GpMatrixOrder) GpStatus {
	ret, _, _ := gdipRotatePenTransform.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(angle)),
		uintptr(order))
	return GpStatus(ret)
}
func GdipSetPenColor(pen *GpPen, argb ARGB) GpStatus {
	ret, _, _ := gdipSetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(argb))
	return GpStatus(ret)
}
func GdipGetPenColor(pen *GpPen, argb *ARGB) GpStatus {
	ret, _, _ := gdipGetPenColor.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(argb)))
	return GpStatus(ret)
}
func GdipSetPenBrushFill(pen *GpPen, brush *GpBrush) GpStatus {
	ret, _, _ := gdipSetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenBrushFill(pen *GpPen, brush **GpBrush) GpStatus {
	ret, _, _ := gdipGetPenBrushFill.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}
func GdipGetPenFillType(pen *GpPen, penType *GpPenType) GpStatus {
	ret, _, _ := gdipGetPenFillType.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(penType)))
	return GpStatus(ret)
}
func GdipGetPenDashStyle(pen *GpPen, dashStyle *GpDashStyle) GpStatus {
	ret, _, _ := gdipGetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dashStyle)))
	return GpStatus(ret)
}
func GdipSetPenDashStyle(pen *GpPen, dashStyle GpDashStyle) GpStatus {
	ret, _, _ := gdipSetPenDashStyle.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(dashStyle))
	return GpStatus(ret)
}
func GdipGetPenDashOffset(pen *GpPen, offset *float32) GpStatus {
	var iOffset uint32
	ret, _, _ := gdipGetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(&iOffset)))
	*offset = math.Float32frombits(iOffset)
	return GpStatus(ret)
}
func GdipSetPenDashOffset(pen *GpPen, offset float32) GpStatus {
	ret, _, _ := gdipSetPenDashOffset.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(offset)))
	return GpStatus(ret)
}
func GdipGetPenDashCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := gdipGetPenDashCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}
func GdipSetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipSetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}
func GdipGetPenDashArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipGetPenDashArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundCount(pen *GpPen, count *int32) GpStatus {
	ret, _, _ := gdipGetPenCompoundCount.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(count)))
	return GpStatus(ret)
}

func GdipSetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipSetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

func GdipGetPenCompoundArray(pen *GpPen, dash *float32, count int32) GpStatus {
	ret, _, _ := gdipGetPenCompoundArray.Call(
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(dash)),
		uintptr(count))
	return GpStatus(ret)
}

// Brush

func GdipCloneBrush(brush *GpBrush, clone **GpBrush) GpStatus {
	ret, _, _ := gdipCloneBrush.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(clone)))
	return GpStatus(ret)
}

func GdipDeleteBrush(brush *GpBrush) GpStatus {
	ret, _, _ := gdipDeleteBrush.Call(uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipGetBrushType(brush *GpBrush, brushType *GpBrushType) GpStatus {
	ret, _, _ := gdipGetBrushType.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(brushType)))
	return GpStatus(ret)
}

// Solid Brush

// func GdipCreateSolidFill(color ARGB, brush **GpSolidFill) GpStatus {
// 	ret, _, _ := gdipCreateSolidFill.Call(
// 		uintptr(color),
// 		uintptr(unsafe.Pointer(brush)))
// 	return GpStatus(ret)
// }

func GdipCreateSolidFill(color ARGB, brush **GpBrush) GpStatus {
	ret, _, _ := gdipCreateSolidFill.Call(
		uintptr(color),
		uintptr(unsafe.Pointer(brush)))
	return GpStatus(ret)
}

func GdipSetSolidFillColor(brush *GpBrush, color ARGB) GpStatus {
	ret, _, _ := gdipSetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(color))
	return GpStatus(ret)
}

func GdipGetSolidFillColor(brush *GpBrush, color *ARGB) GpStatus {
	ret, _, _ := gdipGetSolidFillColor.Call(
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(color)))
	return GpStatus(ret)
}

// Font
func GdipCreateFontFromDC(hdc uintptr, font **GpFont) GpStatus {
	ret, _, _ := gdipCreateFontFromDC.Call(
		hdc,
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipCreateFont(fontFamily *GpFontFamily, emSize float32, style int32, unit GpUnit, font **GpFont) GpStatus {
	ret, _, _ := gdipCreateFont.Call(
		uintptr(unsafe.Pointer(fontFamily)),
		uintptr(math.Float32bits(emSize)),
		uintptr(style),
		uintptr(unit),
		uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipDeleteFont(font *GpFont) GpStatus {
	ret, _, _ := gdipDeleteFont.Call(uintptr(unsafe.Pointer(font)))
	return GpStatus(ret)
}

func GdipNewInstalledFontCollection(fontCollection **GpFontCollection) GpStatus {
	ret, _, _ := gdipNewInstalledFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return GpStatus(ret)
}

func GdipCreateFontFamilyFromName(name *uint16, fontCollection *GpFontCollection, fontFamily **GpFontFamily) GpStatus {
	ret, _, _ := gdipCreateFontFamilyFromName.Call(
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

func GdipDeleteFontFamily(fontFamily *GpFontFamily) GpStatus {
	ret, _, _ := gdipDeleteFontFamily.Call(uintptr(unsafe.Pointer(fontFamily)))
	return GpStatus(ret)
}

// StringFormat

func GdipCreateStringFormat(formatAttributes int32, language uint16, format **GpStringFormat) GpStatus {
	ret, _, _ := gdipCreateStringFormat.Call(
		uintptr(formatAttributes),
		uintptr(language),
		uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipStringFormatGetGenericTypographic(format **GpStringFormat) GpStatus {
	ret, _, _ := gdipStringFormatGetGenericTypographic.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

func GdipDeleteStringFormat(format *GpStringFormat) GpStatus {
	ret, _, _ := gdipDeleteStringFormat.Call(uintptr(unsafe.Pointer(format)))
	return GpStatus(ret)
}

// Path

func GdipCreatePath(brushMode int32, path **GpPath) GpStatus {
	ret, _, _ := gdipCreatePath.Call(uintptr(brushMode), uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipDeletePath(path *GpPath) GpStatus {
	ret, _, _ := gdipDeletePath.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipAddPathArc(path *GpPath, x, y, width, height, startAngle, sweepAngle float32) GpStatus {
	ret, _, _ := gdipAddPathArc.Call(
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
	ret, _, _ := gdipAddPathArcI.Call(
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
	ret, _, _ := gdipAddPathLine.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return GpStatus(ret)
}

func GdipAddPathLineI(path *GpPath, x1, y1, x2, y2 int32) GpStatus {
	ret, _, _ := gdipAddPathLineI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return GpStatus(ret)
}

func GdipClosePathFigure(path *GpPath) GpStatus {
	ret, _, _ := gdipClosePathFigure.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

func GdipClosePathFigures(path *GpPath) GpStatus {
	ret, _, _ := gdipClosePathFigures.Call(uintptr(unsafe.Pointer(path)))
	return GpStatus(ret)
}

// Image

func GdipGetImageGraphicsContext(image *GpImage, graphics **GpGraphics) GpStatus {
	ret, _, _ := gdipGetImageGraphicsContext.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(graphics)))
	return GpStatus(ret)
}

func GdipLoadImageFromFile(filename *uint16, image **GpImage) GpStatus {
	ret, _, _ := gdipLoadImageFromFile.Call(
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(image)))
	return GpStatus(ret)
}

// func GdipSaveImageToFile(image *GpBitmap, filename *uint16, clsidEncoder *ole.GUID, encoderParams *EncoderParameters) GpStatus {
// 	ret, _, _ := gdipSaveImageToFile.Call(uintptr(unsafe.Pointer(image)),
// 		uintptr(unsafe.Pointer(filename)), uintptr(unsafe.Pointer(clsidEncoder)),
// 		uintptr(unsafe.Pointer(encoderParams)))
// 	return GpStatus(ret)
// }

func GdipGetImageWidth(image *GpImage, width *uint32) GpStatus {
	ret, _, _ := gdipGetImageWidth.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(width)))
	return GpStatus(ret)
}

func GdipGetImageHeight(image *GpImage, height *uint32) GpStatus {
	ret, _, _ := gdipGetImageHeight.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(height)))
	return GpStatus(ret)
}

func GdipDisposeImage(image *GpImage) GpStatus {
	ret, _, _ := syscall.Syscall(gdipDisposeImage.Addr(), 1,
		uintptr(unsafe.Pointer(image)),
		0,
		0)

	return GpStatus(ret)
}

// Bitmap

func GdipCreateBitmapFromFile(filename *uint16, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateBitmapFromFile.Addr(), 2,
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(bitmap)),
		0)

	return GpStatus(ret)
}

func GdipCreateBitmapFromHBITMAP(hbm HBITMAP, hpal HPALETTE, bitmap **GpBitmap) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateBitmapFromHBITMAP.Addr(), 3,
		uintptr(hbm),
		uintptr(hpal),
		uintptr(unsafe.Pointer(bitmap)))

	return GpStatus(ret)
}

func GdipCreateHBITMAPFromBitmap(bitmap *GpBitmap, hbmReturn *HBITMAP, background ARGB) GpStatus {
	ret, _, _ := syscall.Syscall(gdipCreateHBITMAPFromBitmap.Addr(), 3,
		uintptr(unsafe.Pointer(bitmap)),
		uintptr(unsafe.Pointer(hbmReturn)),
		uintptr(background))

	return GpStatus(ret)
}

func GdipCreateBitmapFromScan0(width, height, stride int32, format PixelFormat, scan0 *byte, bitmap **GpBitmap) GpStatus {
	ret, _, _ := gdipCreateBitmapFromScan0.Call(
		uintptr(width),
		uintptr(height),
		uintptr(stride),
		uintptr(format),
		uintptr(unsafe.Pointer(scan0)),
		uintptr(unsafe.Pointer(bitmap)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySerif(family **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilySerif.Call(
		uintptr(unsafe.Pointer(family)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilySansSerif(family **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilySansSerif.Call(
		uintptr(unsafe.Pointer(family)))
	return GpStatus(ret)
}

func GdipGetGenericFontFamilyMonospace(family **GpFontFamily) GpStatus {
	ret, _, _ := gdipGetGenericFontFamilyMonospace.Call(
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
