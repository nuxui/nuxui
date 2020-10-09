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
	log.I("nuxui", "ui.Column %s Measure width=%s, height=%s", me.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Column %s Measure", me.ID())

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

	originWidth := width
	originHeight := height

	childrenMeasuredFlags := make([]byte, len(me.Children()))
	maxWidthChanged := false

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
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.PaddingLeft().Value() / 100.0 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.PaddingLeft().Value()
				// ok, wait until maxWidth measured
			}
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			r := me.PaddingRight().Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
			hPPxUsed += r
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.PaddingRight().Value() / 100.0 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.PaddingRight().Value()
				// ok, wait until maxWidth measured
			}
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			t := me.PaddingTop().Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
			vPPxUsed += t
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.PaddingTop().Value() / 100.0 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
				// ok
			case nux.Auto, nux.Unlimit:
				vPPt += me.PaddingTop().Value()
				// ok, wait until height measured
			}
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			b := me.PaddingBottom().Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
			vPPxUsed += b
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.PaddingBottom().Value() / 100.0 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
				// ok
			case nux.Auto, nux.Unlimit:
				vPPt += me.PaddingBottom().Value()
				// ok, wait until height measured
			}
		}
	}
	// measureDuration_1lun := log.Time()

	// hPPt, vPPt > 100 ?????
	// innerWidth, innerHeight < 0 ??????
	innerWidth = float32(nux.MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(nux.MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		// TODO:: layout changed, but child do not need measured
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			// nux.Unlimit ???????
			cms.Width = nux.MeasureSpec(0, nux.Unlimit)
			cms.Height = nux.MeasureSpec(0, nux.Unlimit)

			if cs.HasMargin() {
				switch cs.MarginTop().Mode() {
				case nux.Pixel:
					t := cs.MarginTop().Value()
					cms.Margin.Top = util.Roundi32(t)
					vPx += t
					vPxUsed += t
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cs.MarginTop().Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cms.Margin.Top = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						t := cs.MarginTop().Value() / 100 * innerHeight
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
						vPxUsed += t
						// ok
					case nux.Auto, nux.Unlimit:
						vPt += cs.MarginTop().Value()
						// ok, wait until percent grand total
					}
				}

				switch cs.MarginBottom().Mode() {
				case nux.Pixel:
					b := cs.MarginBottom().Value()
					cms.Margin.Bottom = util.Roundi32(b)
					vPx += b
					vPxUsed += b
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cs.MarginBottom().Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cms.Margin.Bottom = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						b := cs.MarginBottom().Value() / 100 * innerHeight
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						// ok
					case nux.Auto, nux.Unlimit:
						vPt += cs.MarginBottom().Value()
						// ok, wait until percent grand total
					}
				}
			}

			canMeasureHeight := true

			// do not add height to vPxUsed until height measured
			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cms, nux.Pixel)
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					vWt += cs.Height().Value()
					// ok, wait until weight grand total
					canMeasureHeight = false
				case nux.Auto, nux.Unlimit:
					cms.Height = nux.MeasureSpec(0, nux.Pixel)
					setRatioWidth(cs, cms, nux.Pixel)
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					vPt += cs.Height().Value()
					// ok, wait until percent grand total
					canMeasureHeight = false
				}
				// ok
			case nux.Ratio:
				if cs.Width().Mode() == nux.Ratio {
					log.Fatal("nuxui", "width and height size mode can not both Ratio")
				}
				// ok
			case nux.Auto, nux.Unlimit:
				h := innerHeight - vPxUsed
				cms.Height = nux.MeasureSpec(util.Roundi32(h), cs.Height().Mode())
				setRatioWidth(cs, cms, nux.Pixel)
				// ok
			}

			log.I("nuxui", "measureHorizontal 1 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, canMeasureHeight, 0)

			measuredFlags, hPx := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, canMeasureHeight, 0)

			// add height to vPxUsed after height measured
			if nux.MeasureSpecMode(cms.Height) == nux.Pixel {
				vPxUsed += float32(cms.Height)
			}

			// find max innerWidth
			if hPx > hPxMax {
				hPxMax = hPx
			}

			childrenMeasuredFlags[index] = measuredFlags
		}
	}

	// log.TimeEnd(measureDuration_1lun, "nuxui", "ui.Column Measure measureDuration_1lun %s", me.ID())

	// Use the maximum width found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)

		if nux.MeasureSpecValue(originWidth) != nux.MeasureSpecValue(width) {
			maxWidthChanged = true
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

		// all ok
	}

	vPxRemain = innerHeight - vPxUsed

	log.V("nuxui", "---- first traversal end  hPxMax = %f, vPxRemain = %f", hPxMax, vPxRemain)

	// Second traversal, start to calculate the weight size
	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			canMeasureHeight := true

			// Vertical weight only works for Pixel Mode
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				// Vertical weight only works for nux.Pixel Mode
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
						setRatioWidth(cs, cms, nux.Pixel)
					} else {
						cms.Height = nux.MeasureSpec(0, nux.Pixel)
						setRatioWidth(cs, cms, nux.Pixel)
					}
				}
			case nux.Auto, nux.Unlimit:
				if cs.Height().Mode() == nux.Percent {
					// wait all weight size finish measured
					canMeasureHeight = false
				}
			}

			if measuredFlags := childrenMeasuredFlags[index]; (measuredFlags&hMeasuredWidthComplete != hMeasuredWidthComplete || maxWidthChanged) && canMeasureHeight {
				// now width mode is must be nux.Pixel
				if nux.MeasureSpecMode(width) != nux.Pixel {
					log.Fatal("nuxui", "can not run here")
				}

				needMeasure := true
				if maxWidthChanged {
					if measuredFlags&hMeasuredWidthComplete == hMeasuredWidthComplete {
						log.I("nuxui", "hMeasuredWidthComplete child:%s hpx=%d", child.ID(), (cms.Width + cms.Margin.Left + cms.Margin.Right))
						if cms.Width+cms.Margin.Left+cms.Margin.Right == util.Roundi32(innerWidth) {
							needMeasure = false
						}
					}
				}

				if needMeasure {
					log.I("nuxui", "measureHorizontal 2 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, canMeasureHeight, measuredFlags)

					newMeasuredFlags, _ := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, canMeasureHeight, measuredFlags)

					childrenMeasuredFlags[index] = newMeasuredFlags
				}
			}

			// Accumulate child.height that was not accumulated before to get the total value of vPx
			if nux.MeasureSpecMode(cms.Height) == nux.Pixel {
				vPx += float32(cms.Height)
			}

		}
	}

	// height mode is nux.Auto, and child has percent size
	switch nux.MeasureSpecMode(height) {
	case nux.Auto, nux.Unlimit:
		if vPt < 0 || vPt > 100 {
			log.Fatal("nuxui", "children percent size is out of range, it should 0% ~ 100%")
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
						setRatioWidth(cs, cms, nux.Pixel)
					}

					if measuredFlags := childrenMeasuredFlags[index]; measuredFlags != hMeasuredWidthComplete {
						// now width mode is must be nux.Pixel
						log.I("nuxui", "measureHorizontal 3 width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, true, measuredFlags)
						newMeasuredFlags, _ := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, true, measuredFlags)

						if newMeasuredFlags != hMeasuredWidthComplete {
							log.Fatal("nuxui", "can not run here")
						}
					}
				}
			}
		}
	}

	// log.I("nuxui", "originWidth=%s, originHeight=%s", nux.MeasureSpecString(originWidth), nux.MeasureSpecString(originHeight))

	switch nux.MeasureSpecMode(originWidth) {
	case nux.Pixel:
		ms.Width = originWidth
	case nux.Unlimit:
		ms.Width = width
	case nux.Auto:
		if nux.MeasureSpecValue(width) > nux.MeasureSpecValue(originWidth) {
			ms.Width = nux.MeasureSpec(nux.MeasureSpecValue(originWidth), nux.Pixel)
		} else {
			ms.Width = width
		}
	}

	switch nux.MeasureSpecMode(originHeight) {
	case nux.Pixel:
		ms.Height = originHeight
	case nux.Unlimit:
		ms.Height = height
	case nux.Auto:
		if nux.MeasureSpecValue(height) > nux.MeasureSpecValue(originHeight) {
			ms.Height = nux.MeasureSpec(nux.MeasureSpecValue(originHeight), nux.Pixel)
		} else {
			ms.Height = height
		}
	}
}

// return hpx in innerWidth
func (me *column) measureHorizontal(width, height int32, hPPx, hPPt, innerWidth float32,
	child nux.Widget, canMeasureHeight bool, hasMeasuredFlags byte) (measuredFlags byte, hpxmax float32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd("nuxui", "ui.Column measureHorizontal", measuredDuration)
	var hWt float32
	var hPx float32
	var hPt float32
	var hPxUsed float32
	var hPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cms := cs.MeasuredSize()

		// if alrady measured complete, return
		if hasMeasuredFlags == hMeasuredWidthComplete {
			if hpx := cms.Width + cms.Margin.Left + cms.Margin.Right; hpx == nux.MeasureSpecValue(width) {
				return hasMeasuredFlags, float32(hpx)
			}
		}

		// 1. First determine the known size of width, add weight and percentage
		// do not add width to hPxUsed until width measured
		// do not add width to hPx until width measured
		if hasMeasuredFlags&hMeasuredWidth != hMeasuredWidth {
			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cms, nux.Pixel)
				measuredFlags |= hMeasuredWidth
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					hWt += cs.Width().Value()
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured.
				}
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, nux.Pixel)
					measuredFlags |= hMeasuredWidth
					// ok
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured
					// ok
				}
			case nux.Ratio:
				// was measured when measure height
				// ok
			case nux.Auto, nux.Unlimit:
				// measure later
				// ok
			}
		}

		// 2. First determine the known size of width, add weight and percentage
		if cs.HasMargin() {
			if hasMeasuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
				switch cs.MarginLeft().Mode() {
				case nux.Pixel:
					l := cs.MarginLeft().Value()
					cms.Margin.Left = util.Roundi32(l)
					hPx += l
					hPxUsed += l
					measuredFlags |= hMeasuredMarginLeft
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cs.MarginLeft().Value()
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						log.V("nuxui", "child:%s Margin Left not measured", child.ID())
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						l := cs.MarginLeft().Value() / 100 * innerWidth
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						measuredFlags |= hMeasuredMarginLeft
						// ok
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						// ok
					}

				}
			}

			if hasMeasuredFlags&hMeasuredMarginRight != hMeasuredMarginRight {
				switch cs.MarginRight().Mode() {
				case nux.Pixel:
					r := cs.MarginRight().Value()
					cms.Margin.Right = int32(r)
					hPx += r
					hPxUsed += r
					measuredFlags |= hMeasuredMarginRight
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cs.MarginRight().Value()
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						log.V("nuxui", "child:%s Margin Right not measured", child.ID())
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						r := cs.MarginRight().Value() / 100 * innerWidth
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
						hPxUsed += r
						measuredFlags |= hMeasuredMarginRight
						// ok
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						// ok
					}

				}
			}
		} else {
			measuredFlags |= hMeasuredMarginLeft
			measuredFlags |= hMeasuredMarginRight
		}

		hPxRemain = innerWidth - hPxUsed

		if hPxRemain < 0 {
			hPxRemain = 0
		}

		canMeasureWidth := true

		// Start calculating, only auto and weight are left here
		if hasMeasuredFlags&hMeasuredWidth != hMeasuredWidth {
			switch cs.Width().Mode() {
			// case nux.Pixel: has measured
			case nux.Weight:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / hWt * hPxRemain
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, nux.Pixel)
					measuredFlags |= hMeasuredWidth
					// ok
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured.
					canMeasureWidth = false
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				// case nux.Pixel: has measured
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured.
					canMeasureWidth = false
					// ok
				}
			// case nux.Ratio: has measured
			case nux.Auto, nux.Unlimit:
				cms.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode())
				setRatioHeight(cs, cms, nux.Pixel)
				// ok
			}
		}

		log.V("nuxui", "child: %s, canMeasureWidth=%t, canMeasureHeight=%t", child.ID(), canMeasureWidth, canMeasureHeight)
		if canMeasureWidth && canMeasureHeight {
			if m, ok := child.(nux.Measure); ok {
				if hasMeasuredFlags&hMeasuredWidth != hMeasuredWidth {
					log.V("nuxui", "child: %s, Measure width=%s, height=%s", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height))
					m.Measure(cms.Width, cms.Height)

					if cs.Width().Mode() == nux.Ratio {
						oldWidth := cms.Width
						cms.Width = nux.MeasureSpec(util.Roundi32(float32(cms.Height)*cs.Width().Value()), nux.Pixel)
						if oldWidth != cms.Width {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(cms.Width, cms.Height)
						}
					}

					if cs.Height().Mode() == nux.Ratio {
						oldHeight := cms.Height
						cms.Height = nux.MeasureSpec(util.Roundi32(float32(cms.Width)/cs.Height().Value()), nux.Pixel)
						if oldHeight != cms.Height {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(cms.Width, cms.Height)
						}
					}
					measuredFlags |= hMeasuredWidth
				}

				hPxRemain -= float32(cms.Width)
				hPx += float32(cms.Width)

				if cs.HasMargin() {
					if hasMeasuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
						switch cs.MarginLeft().Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(width) {
							case nux.Pixel:
								if hWt > 0 && hPxRemain > 0 {
									cms.Margin.Left = util.Roundi32(cs.MarginLeft().Value() / hWt * hPxRemain)
								} else {
									cms.Margin.Left = 0
								}
								hPx += float32(cms.Margin.Left)
								measuredFlags |= hMeasuredMarginLeft
								// ok
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured.
								// ok
							}
						case nux.Percent:
							switch nux.MeasureSpecMode(width) {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured.
								// ok
							}
						}
					}

					if hasMeasuredFlags&hMeasuredMarginRight != hMeasuredMarginRight {
						switch cs.MarginRight().Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(width) {
							case nux.Pixel:
								if hWt > 0 && hPxRemain > 0 {
									cms.Margin.Right = util.Roundi32(cs.MarginRight().Value() / hWt * hPxRemain)
								} else {
									cms.Margin.Right = 0
								}
								hPx += float32(cms.Margin.Right)
								measuredFlags |= hMeasuredMarginRight
								// ok
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured.
								// ok
							}
						case nux.Percent:
							switch nux.MeasureSpecMode(width) {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured.
								// ok
							}
						}
					}
				}

			}

			// parent width is auto or unlimit
			if hPt > 0 {
				if hPt > 100 {
					log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
				}

				hPx = hPx / (1.0 - hPt/100.0)
			}

			return measuredFlags, hPx
		}
		return measuredFlags, 0
	}

	return measuredFlags, 0
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
			cms.Position.Right = cms.Position.Left + cms.Width
			cms.Position.Bottom = cms.Position.Top + cms.Height
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
