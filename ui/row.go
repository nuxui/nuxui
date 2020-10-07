// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Row interface {
	nux.Parent
	nux.Size
	nux.Creating
	nux.Measure
	nux.Layout
	nux.Draw
	Visual
}

type row struct {
	nux.WidgetParent
	nux.WidgetSize
	WidgetVisual

	Align Align

	childrenWidth float32
}

func NewRow() Row {
	me := &row{
		Align: Align{Vertical: Top, Horizontal: Left},
	}
	me.WidgetParent.Owner = me
	me.WidgetSize.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.Owner = me
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

func (me *row) Creating(attr nux.Attr) {
	if attr == nil {
		attr = nux.Attr{}
	}

	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetParent.Creating(attr)
	me.WidgetVisual.Creating(attr)

	me.Align = *NewAlign(attr.GetAttr("align", nux.Attr{}))
}

func (me *row) onSizeChanged(widget nux.Widget) {

}
func (me *row) onVisualChanged(widget nux.Widget) {

}

func (me *row) Measure(width, height int32) {
	var vPxMax float32   // max horizontal size
	var vPPt float32     // horizontal padding percent
	var vPPx float32     // horizontal padding pixel
	var vPPxUsed float32 // horizontal padding size include percent size

	var hPPt float32
	var hPPx float32
	var hPPxUsed float32

	var hPx float32       // sum of vertical pixel size
	var hWt float32       // sum of children vertical weight
	var hPt float32       // sum of vertical percent size
	var hPxUsed float32   // sum of children vertical size include percent size
	var hPxRemain float32 //

	var innerWidth float32
	var innerHeight float32

	measuredIndex := map[int]struct{}{}
	heightChanged := false

	ms := me.MeasuredSize()
	me.childrenWidth = 0

	// 1. Calculate its own padding size
	if me.HasPadding() {
		switch me.PaddingLeft().Mode() {
		case nux.Pixel:
			l := me.PaddingLeft().Value()
			ms.Padding.Left = util.Roundi32(l)
			hPPx += l
			hPPxUsed += l
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.PaddingLeft().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
			case nux.Auto:
				hPPt += me.PaddingLeft().Value()
			}
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			r := me.PaddingRight().Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
			hPPxUsed += r
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.PaddingRight().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
			case nux.Auto:
				hPPt += me.PaddingRight().Value()
			}
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			t := me.PaddingTop().Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
			vPPxUsed += t
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.PaddingTop().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
			case nux.Auto:
				vPPt += me.PaddingTop().Value()
			}
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			b := me.PaddingBottom().Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
			vPPxUsed += b
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.PaddingBottom().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
			case nux.Auto:
				vPPt += me.PaddingBottom().Value()
			}
		}
	}

	innerWidth = float32(nux.MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(nux.MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			cms.Width = nux.MeasureSpec(0, nux.Unlimit)
			cms.Height = nux.MeasureSpec(0, nux.Unlimit)

			if cs.HasMargin() {
				switch cs.MarginLeft().Mode() {
				case nux.Pixel:
					l := cs.MarginLeft().Value()
					hPx += l
					hPxUsed += l
					cms.Margin.Left = util.Roundi32(l)
				case nux.Weight:
					if nux.MeasureSpecMode(width) == nux.Auto {
						cms.Margin.Left = 0
					} else {
						hWt += cs.MarginLeft().Value()
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						l := cs.MarginLeft().Value() / 100 * innerWidth
						hPx += l
						hPxUsed += l
						cms.Margin.Left = util.Roundi32(l)
					case nux.Auto:
						hPt += cs.MarginLeft().Value()
					}
				}

				switch cs.MarginRight().Mode() {
				case nux.Pixel:
					r := cs.MarginRight().Value()
					cms.Margin.Right = util.Roundi32(r)
					hPx += r
					hPxUsed += r
				case nux.Weight:
					if nux.MeasureSpecMode(width) != nux.Auto {
						hWt += cs.MarginRight().Value()
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						r := cs.MarginRight().Value() / 100 * innerWidth
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
						hPxUsed += r
					case nux.Auto:
						hPt += cs.MarginRight().Value()
					}
				}
			}

			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cms, w, nux.Pixel)
			case nux.Weight:
				// When height is auto, all vertical weights are 0px
				if nux.MeasureSpecMode(width) == nux.Auto {
					cms.Width = nux.MeasureSpec(0, nux.Pixel)
					setRatioHeight(cs, cms, 0, nux.Pixel)
				} else {
					hWt += cs.Width().Value()
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, w, nux.Pixel)
				case nux.Auto:
					hPt += cs.Width().Value()
				}
			case nux.Ratio:
				if cs.Height().Mode() == nux.Ratio {
					log.Fatal("nuxui", "width and height size mode can not both Ratio, at least one is definited.")
				}
			case nux.Auto:
				w := innerWidth - hPxUsed
				cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Auto)
				setRatioHeight(cs, cms, w, nux.Pixel)
			case nux.Unlimit:
				// nothing
			}

			measured, vPx := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, index, measuredIndex)

			// find max innerWidth
			if vPx > vPxMax {
				vPxMax = vPx
			}

			if nux.MeasureSpecMode(cms.Width) == nux.Pixel {
				hPxUsed += float32(cms.Width)
			}

			if measured {
			}
		}
	}

	// Use the maximum height found in the first traversal as the width in auto mode, and calculate the percent size
	if nux.MeasureSpecMode(height) == nux.Auto {
		oldHeight := height
		innerHeight = vPxMax
		h := (innerHeight + vPPx) / (1 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
		if oldHeight != height {
			heightChanged = true
		}

		if me.HasPadding() {
			if me.PaddingTop().Mode() == nux.Percent {
				t := me.PaddingTop().Value() / 100 * h
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
			}

			if me.PaddingBottom().Mode() == nux.Percent {
				b := me.PaddingBottom().Value() / 100 * h
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
			}

			vPPt = 0
		}

	}

	hPxRemain = innerWidth - hPxUsed

	// Second traversal, start to calculate the weight size
	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			// Vertical weight only works for Pixel Mode
			if nux.MeasureSpecMode(width) == nux.Pixel {
				if cs.HasMargin() {
					if cs.MarginLeft().Mode() == nux.Weight {
						l := cs.MarginLeft().Value() / hWt * hPxRemain
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
					}

					if cs.MarginRight().Mode() == nux.Weight {
						r := cs.MarginRight().Value() / hWt * hPxRemain
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
					}
				}

				if cs.Width().Mode() == nux.Weight {
					if hWt > 0 && hPxRemain > 0 {
						w := cs.Width().Value() / hWt * hPxRemain
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cms, w, nux.Pixel)
					} else {
						cms.Width = nux.MeasureSpec(0, nux.Pixel)
						setRatioHeight(cs, cms, 0, nux.Pixel)
					}
				}

			}

			if _, ok := measuredIndex[index]; !ok || heightChanged {
				// now width mode is must be Pixel
				measured, _ := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, index, measuredIndex)

				if measured {
				}

			}

		}
	}

	// width mode is Auto, and child has percent size
	if nux.MeasureSpecMode(width) == nux.Auto {
		if hPt < 0 || hPt > 100 {
			log.Fatal("nuxui", "children percent size is out of range, it should 0% ~ 100%")
		}

		// Accumulate child.height that was not accumulated before to get the total value of vPx
		for _, child := range me.Children() {
			if compt, ok := child.(nux.Component); ok {
				child = compt.Content()
			}

			if cs, ok := child.(nux.Size); ok {
				cms := cs.MeasuredSize()
				if nux.MeasureSpecMode(cms.Width) == nux.Pixel {
					hPx += float32(cms.Width)
				}
			}
		}

		innerWidth = hPx / (1.0 - hPt/100.0)
		w := (innerWidth + hPPx) / (1.0 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)

		if hPPt > 0 {
			if me.HasPadding() {
				if me.PaddingLeft().Mode() == nux.Percent {
					l := me.PaddingLeft().Value() / 100.0 * w
					ms.Padding.Left = util.Roundi32(l)
				}

				if me.PaddingRight().Mode() == nux.Percent {
					r := me.PaddingRight().Value() / 100.0 * w
					ms.Padding.Right = util.Roundi32(r)
				}
			}
		}

		if hPt > 0 {
			for index, child := range me.Children() {
				if compt, ok := child.(nux.Component); ok {
					child = compt.Content()
				}

				if cs, ok := child.(nux.Size); ok {
					cms := cs.MeasuredSize()

					if hPt > 0 {
						if cs.HasMargin() {
							if cs.MarginLeft().Mode() == nux.Percent {
								l := cs.MarginLeft().Value() / 100.0 * innerWidth
								cms.Margin.Left = util.Roundi32(l)
							}

							if cs.MarginRight().Mode() == nux.Percent {
								r := cs.MarginRight().Value() / 100.0 * innerWidth
								cms.Margin.Right = util.Roundi32(r)
							}
						}

						if cs.Width().Mode() == nux.Percent {
							cw := cs.Width().Value() / 100 * innerWidth
							cms.Height = nux.MeasureSpec(util.Roundi32(cw), nux.Pixel)
							setRatioHeight(cs, cms, cw, nux.Pixel)
						}
					}

					if _, ok := measuredIndex[index]; !ok {
						// now width mode is must be nux.Pixel
						measured, _ := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, index, measuredIndex)

						if measured {

						}

					}
				}
			}
		}
	}

	ms.Height = height
	ms.Width = width
}

// return vPx in innerHeight
func (me *row) measureVertical(width, height int32, vPPx, vPPt, innerHeight float32,
	child nux.Widget, index int, measuredIndex map[int]struct{}) (measured bool, hpx float32) {
	var vWt float32
	var vPx float32
	var vPt float32
	var vPxUsed float32
	var vPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cms := cs.MeasuredSize()

		// 1. First determine the known size of height, add weight and percentage
		switch cs.Height().Mode() {
		case nux.Pixel:
			h := cs.Height().Value()
			vPxUsed += h
			cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
			setRatioWidth(cs, cms, h, nux.Pixel)
		case nux.Weight:
			vWt += cs.Height().Value()
		case nux.Percent:
			h := cs.Height().Value() / 100 * innerHeight
			cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
			setRatioWidth(cs, cms, h, nux.Pixel)
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				vPxUsed += h
			case nux.Auto:
				vPxUsed += h
				vPt += cs.Height().Value()
			}
		case nux.Ratio:
			// measured when measure height
		case nux.Auto:
			// measure later
		case nux.Unlimit:
			// ignore
		}

		// 2. First determine the known size of height, add weight and percentage
		if cs.HasMargin() {
			switch cs.MarginTop().Mode() {
			case nux.Pixel:
				t := cs.MarginTop().Value()
				vPx += t
				vPxUsed += t
				cms.Margin.Top = util.Roundi32(t)
			case nux.Weight:
				vWt += cs.MarginTop().Value()
			case nux.Percent:
				t := cs.MarginTop().Value() / 100 * innerHeight
				cms.Margin.Top = util.Roundi32(t)
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					vPx += t
					vPxUsed += t
				case nux.Auto:
					vPxUsed += t
					vPt += cs.MarginTop().Value()
				}

			}

			switch cs.MarginBottom().Mode() {
			case nux.Pixel:
				b := cs.MarginBottom().Value()
				vPx += b
				vPxUsed += b
				cms.Margin.Bottom = int32(b)
			case nux.Weight:
				vWt += cs.MarginBottom().Value()
			case nux.Percent:
				b := cs.MarginBottom().Value() / 100 * innerHeight
				cms.Margin.Bottom = util.Roundi32(b)
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					vPx += b
					vPxUsed += b
				case nux.Auto:
					vPxUsed += b
					vPt += cs.MarginBottom().Value()
				}

			}
		}

		vPxRemain = innerHeight - vPxUsed

		if vPxRemain < 0 {
			vPxRemain = 0
		}

		// Start calculating, only auto and weight are left here
		switch cs.Height().Mode() {
		case nux.Auto:
			cms.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), nux.Auto)
			setRatioWidth(cs, cms, vPxRemain, nux.Pixel)
		case nux.Weight:
			h := cs.Height().Value() / vWt * vPxRemain
			cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
			setRatioWidth(cs, cms, h, nux.Pixel)
		}

		if nux.MeasureSpecMode(cms.Width) != nux.Unlimit && nux.MeasureSpecMode(cms.Height) != nux.Unlimit &&
			!(nux.MeasureSpecMode(height) == nux.Auto && cs.Height().Mode() == nux.Percent) {
			if m, ok := child.(nux.Measure); ok {
				measuredIndex[index] = struct{}{}
				m.Measure(cms.Width, cms.Height)

				if cs.Width().Mode() == nux.Ratio {
					oldWidth := cms.Width
					cms.Width = nux.MeasureSpec(util.Roundi32(float32(cms.Height)*cs.Width().Value()), nux.Pixel)
					if oldWidth != cms.Width {
						m.Measure(cms.Width, cms.Height)
					}
				}

				if cs.Height().Mode() == nux.Ratio {
					oldHeight := cms.Height
					cms.Height = nux.MeasureSpec(util.Roundi32(float32(cms.Width)/cs.Height().Value()), nux.Pixel)
					if oldHeight != cms.Height {
						m.Measure(cms.Width, cms.Height)
					}
				}

				vPxRemain -= float32(cms.Height)
				if cs.HasMargin() && nux.MeasureSpecMode(height) == nux.Pixel {
					switch cs.MarginTop().Mode() {
					case nux.Weight:
						if vWt > 0 && vPxRemain > 0 {
							cms.Margin.Top = util.Roundi32(cs.MarginTop().Value() / vWt * vPxRemain)
						} else {
							cms.Margin.Top = 0
						}
						vPx += float32(cms.Margin.Top)
					}

					switch cs.MarginBottom().Mode() {
					case nux.Weight:
						if vWt > 0 && vPxRemain > 0 {
							cms.Margin.Bottom = util.Roundi32(cs.MarginBottom().Value() / vWt * vPxRemain)
						} else {
							cms.Margin.Bottom = 0
						}
						vPx += float32(cms.Margin.Bottom)
					}
				}

				vPx += float32(cms.Height)
			}

			if vPt > 0 {
				if vPt > 100 {
					log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
				}

				vPx = vPx / (1.0 - vPt/100.0)
			}

			return true, vPx
		} else {
			return false, 0
		}
	}

	return false, 0
}

// Responsible for determining the position of the widget align, margin...
// TODO measure other mode dimen
// layout parent setChildFrame
func (me *row) Layout(dx, dy, left, top, right, bottom int32) {
	ms := me.MeasuredSize()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(bottom - top)
	var innerWidth float32 = float32(right - left)

	innerHeight -= float32(ms.Padding.Top + ms.Padding.Bottom)
	innerWidth -= float32(ms.Padding.Left + ms.Padding.Right)
	switch me.Align.Horizontal {
	case Left:
		l = 0
	case Center:
		l = innerWidth/2 - me.childrenWidth/2
	case Right:
		l = innerWidth - me.childrenWidth
	}
	l += float32(ms.Padding.Left)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		t = 0
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			if cs.Height().Mode() == nux.Weight || (cs.HasMargin() && (cs.MarginTop().Mode() == nux.Weight || cs.MarginBottom().Mode() == nux.Weight)) {
				t += float32(ms.Padding.Top + cms.Margin.Top)
			} else {
				switch me.Align.Vertical {
				case Top:
					t += 0
				case Center:
					t += innerHeight/2 - float32(cms.Height)/2
				case Bottom:
					t += innerHeight - float32(cms.Height)
				}
				t += float32(ms.Padding.Top + cms.Margin.Top)
			}

			l += float32(cms.Margin.Left)

			cms.Position.Left = util.Roundi32(l)
			cms.Position.Top = util.Roundi32(t)
			cms.Position.Right = cms.Position.Left + cms.Width
			cms.Position.Bottom = cms.Position.Top + cms.Height
			cms.Position.X = dx + int32(l)
			cms.Position.Y = dy + int32(t)
			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cms.Position.X, cms.Position.Y, cms.Position.Left, cms.Position.Top, cms.Position.Right, cms.Position.Bottom)
			}

			l += float32(cms.Width + cms.Margin.Right)
		}
	}
}

func (me *row) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		// TODO if child visible == gonecscs , then skip
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()
			if d, ok := child.(nux.Draw); ok {
				canvas.Save()
				canvas.Translate(cms.Position.Left, cms.Position.Top)
				canvas.ClipRect(0, 0, cms.Width, cms.Height)
				d.Draw(canvas)
				canvas.Restore()
			}
		}
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
