// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"strings"

	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

func setRatioHeightIfNeed(cs nux.Size, cms *nux.Frame, width float32, mode nux.Mode) {
	if cs.Height().Mode() == nux.Ratio {
		if cms.Width == 0 {
			cms.Height = 0
		} else {
			cms.Height = int32(nux.MeasureSpec(util.Roundi32(width/cs.Height().Value()), mode))
		}
	}
}

func setRatioWidthIfNeed(cs nux.Size, cms *nux.Frame, height float32, mode nux.Mode) {
	if cs.Width().Mode() == nux.Ratio {
		if cms.Height == 0 {
			cms.Width = 0
		} else {
			cms.Width = int32(nux.MeasureSpec(util.Roundi32(height*cs.Width().Value()), mode))
		}
	}
}

func setNewWidth(frame *nux.Frame, originWidth, newWidth nux.MeasureDimen) {
	switch originWidth.Mode() {
	case nux.Pixel:
		frame.Width = originWidth.Value()
	case nux.Unlimit:
		frame.Width = newWidth.Value()
	case nux.Auto:
		if newWidth.Value() > originWidth.Value() {
			frame.Width = originWidth.Value()
		} else {
			frame.Width = newWidth.Value()
		}
	}

}

func setNewHeight(frame *nux.Frame, originHeight, newHeight nux.MeasureDimen) {
	switch originHeight.Mode() {
	case nux.Pixel:
		frame.Height = originHeight.Value()
	case nux.Unlimit:
		frame.Height = newHeight.Value()
	case nux.Auto:
		if newHeight.Value() > originHeight.Value() {
			frame.Height = originHeight.Value()
		} else {
			frame.Height = newHeight.Value()
		}
	}
}

func mergedStateFromString(statestr string) (state uint32) {
	names := strings.Split(strings.TrimSpace(statestr), "|")
	for _, name := range names {
		state |= nux.StateFromString(name)
	}
	return
}

func measurePadding(width, height nux.MeasureDimen, p *nux.Padding, frame *nux.Frame, emsSize float32, measuredFlag byte) (hPPx, hPPt, vPPx, vPPt float32, measuredFlagOut byte) {
	if p != nil && measuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		if measuredFlag&flagMeasuredPaddingLeft != flagMeasuredPaddingLeft {
			if p.Left.Value() != 0 {
				switch p.Left.Mode() {
				case nux.Pixel:
					l := p.Left.Value()
					frame.Padding.Left = util.Roundi32(l)
					hPPx += l
					measuredFlag |= flagMeasuredPaddingLeft
				case nux.Ems:
					if emsSize > 0 {
						l := emsSize * p.Left.Value()
						frame.Padding.Left = util.Roundi32(l)
						hPPx += l
					}
					measuredFlag |= flagMeasuredPaddingLeft
				case nux.Percent:
					switch width.Mode() {
					case nux.Pixel:
						l := p.Left.Value() / 100.0 * float32(width.Value())
						frame.Padding.Left = util.Roundi32(l)
						hPPx += l
						measuredFlag |= flagMeasuredPaddingLeft
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured
						hPPt += p.Left.Value()
					}
				}
			} else {
				measuredFlag |= flagMeasuredPaddingLeft
			}
		}

		if measuredFlag&flagMeasuredPaddingRight != flagMeasuredPaddingRight {
			if p.Right.Value() != 0 {
				switch p.Right.Mode() {
				case nux.Pixel:
					r := p.Right.Value()
					frame.Padding.Right = util.Roundi32(r)
					hPPx += r
					measuredFlag |= flagMeasuredPaddingRight
				case nux.Ems:
					if emsSize > 0 {
						r := emsSize * p.Right.Value()
						frame.Padding.Right = util.Roundi32(r)
						hPPx += r
					}
					measuredFlag |= flagMeasuredPaddingRight
				case nux.Percent:
					switch width.Mode() {
					case nux.Pixel:
						r := p.Right.Value() / 100.0 * float32(width.Value())
						frame.Padding.Right = util.Roundi32(r)
						hPPx += r
						measuredFlag |= flagMeasuredPaddingRight
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured
						hPPt += p.Right.Value()
					}
				}
			} else {
				measuredFlag |= flagMeasuredPaddingRight
			}
		}

		if measuredFlag&flagMeasuredPaddingTop != flagMeasuredPaddingTop {
			if p.Top.Value() != 0 {
				switch p.Top.Mode() {
				case nux.Pixel:
					t := p.Top.Value()
					frame.Padding.Top = util.Roundi32(t)
					vPPx += t
					measuredFlag |= flagMeasuredPaddingTop
				case nux.Ems:
					if emsSize > 0 {
						t := emsSize * p.Top.Value()
						frame.Padding.Top = util.Roundi32(t)
						hPPx += t
					}
					measuredFlag |= flagMeasuredPaddingTop
				case nux.Percent:
					switch height.Mode() {
					case nux.Pixel:
						t := p.Top.Value() / 100.0 * float32(height.Value())
						frame.Padding.Top = util.Roundi32(t)
						vPPx += t
						measuredFlag |= flagMeasuredPaddingTop
					case nux.Auto, nux.Unlimit:
						// wait until height measured
						vPPt += p.Top.Value()
					}
				}
			} else {
				measuredFlag |= flagMeasuredPaddingTop
			}
		}

		if measuredFlag&flagMeasuredPaddingBottom != flagMeasuredPaddingBottom {
			if p.Bottom.Value() != 0 {
				switch p.Bottom.Mode() {
				case nux.Pixel:
					b := p.Bottom.Value()
					frame.Padding.Bottom = util.Roundi32(b)
					vPPx += b
					measuredFlag |= flagMeasuredPaddingBottom
				case nux.Ems:
					if emsSize > 0 {
						b := emsSize * p.Bottom.Value()
						frame.Padding.Bottom = util.Roundi32(b)
						hPPx += b
					}
					measuredFlag |= flagMeasuredPaddingBottom
				case nux.Percent:
					switch height.Mode() {
					case nux.Pixel:
						b := p.Bottom.Value() / 100.0 * float32(height.Value())
						frame.Padding.Bottom = util.Roundi32(b)
						vPPx += b
						measuredFlag |= flagMeasuredPaddingBottom
					case nux.Auto, nux.Unlimit:
						// wait until height measured
						vPPt += p.Bottom.Value()
					}
				}
			} else {
				measuredFlag |= flagMeasuredPaddingBottom
			}
		}
	} else {
		measuredFlag |= flagMeasuredPaddingComplete
	}

	measuredFlagOut = measuredFlag
	return
}

func measureChildrenZeroSize(parent nux.Parent) {
	for _, child := range parent.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cs.Frame().Clear()
		}

		if cm, ok := child.(nux.Measure); ok {
			cm.Measure(0, 0)
		}
	}
}
