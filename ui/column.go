// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Column interface {
	nux.Parent
	nux.Size
	Visual
	nux.Creating
	nux.Layout
	nux.Measure
	nux.Draw
}

type column struct {
	nux.WidgetParent
	nux.WidgetBase
	nux.WidgetSize
	WidgetVisual

	Align          Align // TODO not nil
	childrenHeight float32
}

func NewColumn() Column {
	me := &column{
		Align: Align{Vertical: Top, Horizontal: Left},
	}
	me.WidgetParent.Owner = me
	me.WidgetSize.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.Owner = me
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

func (me *column) Creating(attr nux.Attr) {
	if attr == nil {
		attr = nux.Attr{}
	}

	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetParent.Creating(attr)
	me.WidgetVisual.Creating(attr)

	me.Align = *NewAlign(attr.GetAttr("align", nux.Attr{}))
}

func (me *column) onSizeChanged(widget nux.Widget) {

}
func (me *column) onVisualChanged(widget nux.Widget) {

}

func (me *column) Measure(width, height int32) {
	var hPxMax float32   // max horizontal size
	var hPPt float32     // horizontal padding percent
	var hPPx float32     // horizontal padding pixel
	var hPPxUsed float32 // horizontal padding size include percent size

	var vPPt float32
	var vPPx float32
	var vPPxUsed float32

	var vPx float32       // sum of vertical pixel size
	var vWt float32       // sum of children vertical weight
	var vPt float32       // sum of vertical percent size
	var vPxUsed float32   // sum of children vertical size include percent size
	var vPxRemain float32 //

	var innerWidth float32
	var innerHeight float32

	measuredIndex := map[int]struct{}{}
	widthChanged := false

	ms := me.MeasuredSize()
	me.childrenHeight = 0

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

			cms.Width = nux.MeasureSpec(0, nux.Unspec)
			cms.Height = nux.MeasureSpec(0, nux.Unspec)

			if cs.HasMargin() {
				switch cs.MarginTop().Mode() {
				case nux.Pixel:
					t := cs.MarginTop().Value()
					vPx += t
					vPxUsed += t
					cms.Margin.Top = util.Roundi32(t)
				case nux.Weight:
					if nux.MeasureSpecMode(height) == nux.Auto {
						cms.Margin.Top = 0
					} else {
						vWt += cs.MarginTop().Value()
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						t := cs.MarginTop().Value() / 100 * innerHeight
						vPx += t
						vPxUsed += t
						cms.Margin.Top = util.Roundi32(t)
					case nux.Auto:
						vPt += cs.MarginTop().Value()
					}
				}

				switch cs.MarginBottom().Mode() {
				case nux.Pixel:
					b := cs.MarginBottom().Value()
					cms.Margin.Bottom = util.Roundi32(b)
					vPx += b
					vPxUsed += b
				case nux.Weight:
					if nux.MeasureSpecMode(height) != nux.Auto {
						vWt += cs.MarginBottom().Value()
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						b := cs.MarginBottom().Value() / 100 * innerHeight
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
					case nux.Auto:
						vPt += cs.MarginBottom().Value()
					}
				}
			}

			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cms, h, nux.Pixel)
			case nux.Weight:
				// When height is auto, all vertical weights are 0px
				if nux.MeasureSpecMode(height) == nux.Auto {
					cms.Height = nux.MeasureSpec(0, nux.Pixel)
					setRatioWidth(cs, cms, 0, nux.Pixel)
				} else {
					vWt += cs.Height().Value()
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, h, nux.Pixel)
				case nux.Auto:
					vPt += cs.Height().Value()
				}
			case nux.Ratio:
				if cs.Width().Mode() == nux.Ratio {
					log.Fatal("nuxui", "width and height size mode can not both nux.Ratio, at least one is definited.")
				}
			case nux.Auto:
				h := innerHeight - vPxUsed
				cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Auto)
				setRatioWidth(cs, cms, h, nux.Pixel)
			case nux.Default, nux.Unspec:
				// nothing
			}

			measured, hPx := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, index, measuredIndex)

			// find max innerWidth
			if hPx > hPxMax {
				hPxMax = hPx
			}

			if nux.MeasureSpecMode(cms.Height) == nux.Pixel {
				vPxUsed += float32(cms.Height)
			}

			if measured {
			}
		}
	}

	// TODO:: Optimize the number of traverses
	// if len(measuredIndex) == me.ChildrenCount() {
	// 	if nux.MeasureSpecMode(width) == nux.Auto {

	// 	}
	// 	return
	// }

	// Use the maximum width found in the first traversal as the width in auto mode, and calculate the percent size
	if nux.MeasureSpecMode(width) == nux.Auto {
		oldWidth := width
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
		if oldWidth != width {
			widthChanged = true
		}

		if me.HasPadding() {
			if me.PaddingLeft().Mode() == nux.Percent {
				l := me.PaddingLeft().Value() / 100 * w
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
			}

			if me.PaddingRight().Mode() == nux.Percent {
				r := me.PaddingRight().Value() / 100 * w
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
			}

			hPPt = 0
		}

	}

	vPxRemain = innerHeight - vPxUsed

	// Second traversal, start to calculate the weight size
	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			// Vertical weight only works for nux.Pixel Mode
			if nux.MeasureSpecMode(height) == nux.Pixel {
				if cs.HasMargin() {
					if cs.MarginTop().Mode() == nux.Weight {
						t := cs.MarginTop().Value() / vWt * vPxRemain
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
					}

					if cs.MarginBottom().Mode() == nux.Weight {
						b := cs.MarginBottom().Value() / vWt * vPxRemain
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
					}
				}

				if cs.Height().Mode() == nux.Weight {
					if vWt > 0 && vPxRemain > 0 {
						h := cs.Height().Value() / vWt * vPxRemain
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cms, h, nux.Pixel)
					} else {
						cms.Height = nux.MeasureSpec(0, nux.Pixel)
						setRatioWidth(cs, cms, 0, nux.Pixel)
					}
				}

			}

			if _, ok := measuredIndex[index]; !ok || widthChanged {
				// now width mode is must be nux.Pixel
				measured, _ := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, index, measuredIndex)

				if measured {
				}

			}

		}
	}

	// height mode is nux.Auto, and child has percent size
	if nux.MeasureSpecMode(height) == nux.Auto {
		if vPt < 0 || vPt > 100 {
			log.Fatal("nuxui", "children percent size is out of range, it should 0% ~ 100%")
		}

		// Accumulate child.height that was not accumulated before to get the total value of vPx
		for _, child := range me.Children() {
			if compt, ok := child.(nux.Component); ok {
				child = compt.Content()
			}

			if cs, ok := child.(nux.Size); ok {
				cms := cs.MeasuredSize()
				if nux.MeasureSpecMode(cms.Height) == nux.Pixel {
					vPx += float32(cms.Height)
				}
			}
		}

		innerHeight = vPx / (1.0 - vPt/100.0)
		h := (innerHeight + vPPx) / (1.0 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)

		if vPPt > 0 {
			if me.HasPadding() {
				if me.PaddingTop().Mode() == nux.Percent {
					t := me.PaddingTop().Value() / 100.0 * h
					ms.Padding.Top = util.Roundi32(t)
				}

				if me.PaddingBottom().Mode() == nux.Percent {
					b := me.PaddingBottom().Value() / 100.0 * h
					ms.Padding.Bottom = util.Roundi32(b)
				}
			}
		}

		if vPt > 0 {
			for index, child := range me.Children() {
				if compt, ok := child.(nux.Component); ok {
					child = compt.Content()
				}

				if cs, ok := child.(nux.Size); ok {
					cms := cs.MeasuredSize()

					if vPt > 0 {
						if cs.HasMargin() {
							if cs.MarginTop().Mode() == nux.Percent {
								t := cs.MarginTop().Value() / 100.0 * innerHeight
								cms.Margin.Top = util.Roundi32(t)
							}

							if cs.MarginBottom().Mode() == nux.Percent {
								b := cs.MarginBottom().Value() / 100.0 * innerHeight
								cms.Margin.Bottom = util.Roundi32(b)
							}
						}

						if cs.Height().Mode() == nux.Percent {
							ch := cs.Height().Value() / 100 * innerHeight
							cms.Height = nux.MeasureSpec(util.Roundi32(ch), nux.Pixel)
							setRatioWidth(cs, cms, ch, nux.Pixel)
						}
					}

					if _, ok := measuredIndex[index]; !ok {
						// now width mode is must be nux.Pixel
						measured, _ := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, index, measuredIndex)

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

// return hpx in innerWidth
func (me *column) measureHorizontal(width, height int32, hPPx, hPPt, innerWidth float32,
	child nux.Widget, index int, measuredIndex map[int]struct{}) (measured bool, hpx float32) {
	var hWt float32
	var hPx float32
	var hPt float32
	var hPxUsed float32
	var hPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cms := cs.MeasuredSize()

		// 1. First determine the known size of height, add weight and percentage
		switch cs.Width().Mode() {
		case nux.Pixel:
			w := cs.Width().Value()
			hPxUsed += w
			cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
			setRatioHeight(cs, cms, w, nux.Pixel)
		case nux.Weight:
			hWt += cs.Width().Value()
		case nux.Percent:
			w := cs.Width().Value() / 100 * innerWidth
			cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
			setRatioHeight(cs, cms, w, nux.Pixel)
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				hPxUsed += w
			case nux.Auto:
				hPxUsed += w
				hPt += cs.Width().Value()
			}
		case nux.Ratio:
			// measured when measure height
		case nux.Auto:
			// measure later
		case nux.Default, nux.Unspec:
			// ignore
		}

		// 2. First determine the known size of height, add weight and percentage
		if cs.HasMargin() {
			switch cs.MarginLeft().Mode() {
			case nux.Pixel:
				l := cs.MarginLeft().Value()
				hPx += l
				hPxUsed += l
				cms.Margin.Left = util.Roundi32(l)
			case nux.Weight:
				hWt += cs.MarginLeft().Value()
			case nux.Percent:
				l := cs.MarginLeft().Value() / 100 * innerWidth
				cms.Margin.Left = util.Roundi32(l)
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					hPx += l
					hPxUsed += l
				case nux.Auto:
					hPxUsed += l
					hPt += cs.MarginLeft().Value()
				}

			}

			switch cs.MarginRight().Mode() {
			case nux.Pixel:
				r := cs.MarginRight().Value()
				hPx += r
				hPxUsed += r
				cms.Margin.Right = int32(r)
			case nux.Weight:
				hWt += cs.MarginRight().Value()
			case nux.Percent:
				r := cs.MarginRight().Value() / 100 * innerWidth
				cms.Margin.Right = util.Roundi32(r)
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					hPx += r
					hPxUsed += r
				case nux.Auto:
					hPxUsed += r
					hPt += cs.MarginRight().Value()
				}

			}
		}

		hPxRemain = innerWidth - hPxUsed

		if hPxRemain < 0 {
			hPxRemain = 0
		}

		// Start calculating, only auto and weight are left here
		switch cs.Width().Mode() {
		case nux.Auto:
			cms.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), nux.Auto)
			setRatioHeight(cs, cms, hPxRemain, nux.Pixel)
		case nux.Weight:
			w := cs.Width().Value() / hWt * hPxRemain
			cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
			setRatioHeight(cs, cms, w, nux.Pixel)
		}

		if nux.MeasureSpecMode(cms.Height) != nux.Unspec && nux.MeasureSpecMode(cms.Width) != nux.Unspec &&
			!(nux.MeasureSpecMode(width) == nux.Auto && cs.Width().Mode() == nux.Percent) {
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

				hPxRemain -= float32(cms.Width)
				if cs.HasMargin() && nux.MeasureSpecMode(width) == nux.Pixel {
					switch cs.MarginLeft().Mode() {
					case nux.Weight:
						if hWt > 0 && hPxRemain > 0 {
							cms.Margin.Left = util.Roundi32(cs.MarginLeft().Value() / hWt * hPxRemain)
						} else {
							cms.Margin.Left = 0
						}
						hPx += float32(cms.Margin.Left)
					}

					switch cs.MarginRight().Mode() {
					case nux.Weight:
						if hWt > 0 && hPxRemain > 0 {
							cms.Margin.Right = util.Roundi32(cs.MarginRight().Value() / hWt * hPxRemain)
						} else {
							cms.Margin.Right = 0
						}
						hPx += float32(cms.Margin.Right)
					}
				}

				hPx += float32(cms.Width)
			}

			if hPt > 0 {
				if hPt > 100 {
					log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
				}

				hPx = hPx / (1.0 - hPt/100.0)
			}

			return true, hPx
		} else {
			return false, 0
		}

	}

	return false, 0
}

// Responsible for determining the position of the widget align, margin
// TODO measure other mode dimen
func (me *column) Layout(dx, dy, left, top, right, bottom int32) {
	// log.V("nuxui", "column layout %d, %d, %d, %d, %d, %d", dx, dy, left, top, right, bottom)
	ms := me.MeasuredSize()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(bottom - top)
	var innerWidth float32 = float32(right - left)

	innerHeight -= float32(ms.Padding.Top + ms.Padding.Bottom)
	innerWidth -= float32(ms.Padding.Left + ms.Padding.Right)
	switch me.Align.Vertical {
	case Top:
		t = 0
	case Center:
		t = innerHeight/2 - me.childrenHeight/2
	case Bottom:
		t = innerHeight - me.childrenHeight
	}
	t += float32(ms.Padding.Top)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		l = 0
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			if cs.Width().Mode() == nux.Weight || (cs.MarginLeft().Mode() == nux.Weight || cs.MarginRight().Mode() == nux.Weight) {
				l += float32(ms.Padding.Left + cms.Margin.Left)
			} else {
				switch me.Align.Horizontal {
				case Left:
					l += 0
				case Center:
					l += innerWidth/2 - float32(cms.Width)/2
				case Right:
					l += innerWidth - float32(cms.Width)
				}
				l += float32(ms.Padding.Left + cms.Margin.Left)
			}

			t += float32(cms.Margin.Top)

			cms.Position.Left = util.Roundi32(l)
			cms.Position.Top = util.Roundi32(t)
			cms.Position.Right = util.Roundi32(l) + cms.Width
			cms.Position.Bottom = util.Roundi32(t) + cms.Height
			cms.Position.X = dx + int32(l)
			cms.Position.Y = dy + int32(t)
			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cms.Position.X, cms.Position.Y, cms.Position.Left, cms.Position.Top, cms.Position.Right, cms.Position.Bottom)
			}

			t += float32(cms.Height + cms.Margin.Bottom)
		}
	}
}

func (me *column) Draw(canvas nux.Canvas) {
	// t1 := time.Now()
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}
	// log.V("nuxui", "draw Background used time %d", time.Now().Sub(t1).Milliseconds())

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		// TODO if child visible == gone , then skip
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()
			if draw, ok := child.(nux.Draw); ok {
				canvas.Save()
				canvas.Translate(cms.Position.Left, cms.Position.Top)
				canvas.ClipRect(0, 0, cms.Width, cms.Height)
				draw.Draw(canvas)
				canvas.Restore()
			}
		}
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
