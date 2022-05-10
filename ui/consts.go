// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

type Visible byte

const (
	Gone Visible = 0
	Show Visible = 1
	Hide Visible = 2

	Center = 0
	Left   = 1
	Top    = 2
	Right  = 3
	Bottom = 4
)

const (
	flagMeasured uint8 = 1 << iota // call Measure func
	flagMeasuredWidth
	flagMeasuredHeight
	flagMeasuredMarginLeft
	flagMeasuredMarginRight
	flagMeasuredMarginTop
	flagMeasuredMarginBottom
	flagMeasuredMarginComplete     = flagMeasuredMarginLeft | flagMeasuredMarginRight | flagMeasuredMarginTop | flagMeasuredMarginBottom
	flagMeasuredHorizontalComplete = flagMeasuredWidth | flagMeasuredMarginLeft | flagMeasuredMarginRight
	flagMeasuredVerticalComplete   = flagMeasuredHeight | flagMeasuredMarginTop | flagMeasuredMarginBottom
	flagMeasuredComplete           = flagMeasured | flagMeasuredWidth | flagMeasuredHeight | flagMeasuredMarginComplete
)

const (
	flagMeasuredPaddingLeft uint8 = 1 << iota
	flagMeasuredPaddingTop
	flagMeasuredPaddingRight
	flagMeasuredPaddingBottom
	flagMeasuredPaddingComplete = flagMeasuredPaddingLeft | flagMeasuredPaddingTop | flagMeasuredPaddingRight | flagMeasuredPaddingBottom
)

const (
	clipChildrenAuto = 0
	clipChildrenYes  = 1
	clipChildrenNo   = -1
)
