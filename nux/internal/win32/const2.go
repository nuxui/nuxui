// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

const LANG_NEUTRAL = 0x00

// Unit
const (
	UnitWorld      = 0 // 0 -- World coordinate (non-physical unit)
	UnitDisplay    = 1 // 1 -- Variable -- for PageTransform only
	UnitPixel      = 2 // 2 -- Each unit is one device pixel.
	UnitPoint      = 3 // 3 -- Each unit is a printer's point, or 1/72 inch.
	UnitInch       = 4 // 4 -- Each unit is 1 inch.
	UnitDocument   = 5 // 5 -- Each unit is 1/300 inch.
	UnitMillimeter = 6 // 6 -- Each unit is 1 millimeter.
)

const (
	AlphaMask = 0xff000000
	RedMask   = 0x00ff0000
	GreenMask = 0x0000ff00
	BlueMask  = 0x000000ff
)

// QualityMode
const (
	QualityModeInvalid = iota - 1
	QualityModeDefault
	QualityModeLow  // Best performance
	QualityModeHigh // Best rendering quality
)

// Alpha Compositing mode
const (
	CompositingModeSourceOver = iota // 0
	CompositingModeSourceCopy        // 1
)

// Alpha Compositing quality
const (
	CompositingQualityInvalid = iota + QualityModeInvalid
	CompositingQualityDefault
	CompositingQualityHighSpeed
	CompositingQualityHighQuality
	CompositingQualityGammaCorrected
	CompositingQualityAssumeLinear
)

// InterpolationMode
const (
	InterpolationModeInvalid = iota + QualityModeInvalid
	InterpolationModeDefault
	InterpolationModeLowQuality
	InterpolationModeHighQuality
	InterpolationModeBilinear
	InterpolationModeBicubic
	InterpolationModeNearestNeighbor
	InterpolationModeHighQualityBilinear
	InterpolationModeHighQualityBicubic
)

type GpSmoothingMode int32

const (
	SmoothingModeInvalid      GpSmoothingMode = QualityModeInvalid
	SmoothingModeDefault      GpSmoothingMode = 0
	SmoothingModeHighSpeed    GpSmoothingMode = 1
	SmoothingModeHighQuality  GpSmoothingMode = 2
	SmoothingModeNone         GpSmoothingMode = 3
	SmoothingModeAntiAlias8x4 GpSmoothingMode = 4
	SmoothingModeAntiAlias    GpSmoothingMode = 4
	SmoothingModeAntiAlias8x8 GpSmoothingMode = 5
)

type GpFlushIntention int32

const (
	FlushIntentionFlush GpFlushIntention = 0
	FlushIntentionSync  GpFlushIntention = 1
)

// Pixel Format Mode
const (
	PixelOffsetModeInvalid = iota + QualityModeInvalid
	PixelOffsetModeDefault
	PixelOffsetModeHighSpeed
	PixelOffsetModeHighQuality
	PixelOffsetModeNone // No pixel offset
	PixelOffsetModeHalf // Offset by -0.5, -0.5 for fast anti-alias perf
)

// Text Rendering Hint
const (
	TextRenderingHintSystemDefault            = iota // Glyph with system default rendering hint
	TextRenderingHintSingleBitPerPixelGridFit        // Glyph bitmap with hinting
	TextRenderingHintSingleBitPerPixel               // Glyph bitmap without hinting
	TextRenderingHintAntiAliasGridFit                // Glyph anti-alias bitmap with hinting
	TextRenderingHintAntiAlias                       // Glyph anti-alias bitmap without hinting
	TextRenderingHintClearTypeGridFit                // Glyph CT bitmap with hinting
)

// Fill mode constants
const (
	FillModeAlternate = iota // 0
	FillModeWinding          // 1
)

// BrushType
const (
	BrushTypeSolidColor GpBrushType = iota
	BrushTypeHatchFill
	BrushTypeTextureFill
	BrushTypePathGradient
	BrushTypeLinearGradient
)

// LineCap
const (
	LineCapFlat GpLineCap = iota
	LineCapSquare
	LineCapRound
	LineCapTriangle
	LineCapNoAnchor
	LineCapSquareAnchor
	LineCapRoundAnchor
	LineCapDiamondAnchor
	LineCapArrowAnchor
	LineCapCustom
	LineCapAnchorMask
)

// LineJoin
const (
	LineJoinMiter GpLineJoin = iota
	LineJoinBevel
	LineJoinRound
	LineJoinMiterClipped
)

// DashCap
const (
	DashCapFlat GpDashCap = iota
	DashCapRound
	DashCapTriangle
)

// DashStyle
const (
	DashStyleSolid GpDashStyle = iota
	DashStyleDash
	DashStyleDot
	DashStyleDashDot
	DashStyleDashDotDot
	DashStyleCustom
)

// PenAlignment
const (
	PenAlignmentCenter GpPenAlignment = iota
	PenAlignmentInset
)

// MatrixOrder
const (
	MatrixOrderPrepend GpMatrixOrder = iota
	MatrixOrderAppend
)

// PenType
const (
	PenTypeSolidColor GpPenType = iota
	PenTypeHatchFill
	PenTypeTextureFill
	PenTypePathGradient
	PenTypeLinearGradient
	PenTypeUnknown
)
