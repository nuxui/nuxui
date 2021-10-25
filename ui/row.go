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
	nux.Measure
	nux.Layout
	nux.Draw
	Visual
}

type row struct {
	*nux.WidgetParent
	*nux.WidgetSize
	*WidgetVisual

	align *Align

	childrenWidth float32
}

func NewRow(context nux.Context, attrs ...nux.Attr) Row {
	attr := getAttr(attrs...)
	me := &row{
		align: NewAlign(attr.GetAttr("align", nux.Attr{})),
	}
	me.WidgetParent = nux.NewWidgetParent(context, me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(context, me, attrs...)
	me.WidgetVisual = NewWidgetVisual(context, me, attrs...)
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

func (me *row) onSizeChanged(widget nux.Widget) {

}
func (me *row) onVisualChanged(widget nux.Widget) {

}

func (me *row) Measure(width, height int32) {
	// log.I("nuxui", "ui.Row %s Measure width=%s, height=%s", me.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Row %s Measure", me.ID())

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

	originWidth := width
	originHeight := height

	childrenMeasuredFlags := make([]byte, len(me.Children()))
	maxHeightChanged := false

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
				switch cs.MarginLeft().Mode() {
				case nux.Pixel:
					l := cs.MarginLeft().Value()
					cms.Margin.Left = util.Roundi32(l)
					hPx += l
					hPxUsed += l
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cs.MarginLeft().Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cms.Margin.Left = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						l := cs.MarginLeft().Value() / 100 * innerWidth
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.MarginLeft().Value()
						// ok, wait until percent grand total
					}
				}

				switch cs.MarginRight().Mode() {
				case nux.Pixel:
					r := cs.MarginRight().Value()
					cms.Margin.Right = util.Roundi32(r)
					hPx += r
					hPxUsed += r
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cs.MarginRight().Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cms.Margin.Right = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						r := cs.MarginRight().Value() / 100 * innerWidth
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
						hPxUsed += r
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.MarginRight().Value()
						// ok, wait until percent grand total
					}
				}
			}

			canMeasureWidth := true

			// do not add width to hPxUsed until width measured
			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cms, w, nux.Pixel)
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					hWt += cs.Width().Value()
					// ok, wait until weight grand total
					if me.ChildrenCount() > 1 {
						canMeasureWidth = false
					}
				case nux.Auto, nux.Unlimit:
					cms.Width = nux.MeasureSpec(0, nux.Pixel)
					setRatioHeight(cs, cms, 0, nux.Pixel)
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, w, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					hPt += cs.Width().Value()
					// ok, wait until percent grand total
					if me.ChildrenCount() > 1 {
						canMeasureWidth = false
					}
				}
				// ok
			case nux.Ratio:
				if cs.Height().Mode() == nux.Ratio {
					log.Fatal("nuxui", "width and height size mode can not both Ratio")
				}
				// ok
			case nux.Auto, nux.Unlimit:
				w := innerWidth - hPxUsed
				cms.Width = nux.MeasureSpec(util.Roundi32(w), cs.Width().Mode())
				setRatioHeight(cs, cms, w, nux.Pixel)
				// ok
			}

			log.I("nuxui", "measureVertical 1 child:%s, width:%s, height:%s, vPPx:%.2f, vPPt:%.2f, innerHeight:%.2f, canMeasureWidth:%t, measuredFlags:%d", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height), vPPx, vPPt, innerHeight, canMeasureWidth, 0)

			measuredFlags, vPx := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, canMeasureWidth, 0)

			// add width to hPxUsed after width measured
			if nux.MeasureSpecMode(cms.Width) == nux.Pixel {
				hPxUsed += float32(cms.Width)
			}

			// find max innerHeight
			if vPx > vPxMax {
				vPxMax = vPx
			}

			childrenMeasuredFlags[index] = measuredFlags
		}
	}

	// log.TimeEnd(measureDuration_1lun, "nuxui", "ui.Column Measure measureDuration_1lun %s", me.ID())

	// Use the maximum height found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(height) {
	case nux.Auto, nux.Unlimit:
		innerHeight = vPxMax
		h := (innerHeight + vPPx) / (1 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)

		if nux.MeasureSpecValue(originHeight) != nux.MeasureSpecValue(height) {
			maxHeightChanged = true
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

		// all ok
	}

	hPxRemain = innerWidth - hPxUsed

	log.V("nuxui", "---- first traversal end  vPxMax = %f, hPxRemain = %f", vPxMax, hPxRemain)

	// Second traversal, start to calculate the weight size
	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			canMeasureWidth := true

			// Vertical weight only works for Pixel Mode
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				// Vertical weight only works for nux.Pixel Mode
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
			case nux.Auto, nux.Unlimit:
				if cs.Width().Mode() == nux.Percent {
					// wait all weight size finish measured
					if me.ChildrenCount() > 1 {
						canMeasureWidth = false
					}
				}
			}

			if measuredFlags := childrenMeasuredFlags[index]; (measuredFlags&hMeasuredHeightComplete != hMeasuredHeightComplete || maxHeightChanged) && canMeasureWidth {
				// now height mode is must be Pixel
				if nux.MeasureSpecMode(height) != nux.Pixel {
					log.Fatal("nuxui", "can not run here")
				}

				needMeasure := true
				if maxHeightChanged {
					if measuredFlags&hMeasuredHeightComplete == hMeasuredHeightComplete {
						log.I("nuxui", "hMeasuredHeightComplete child:%s vpx=%d", child.ID(), (cms.Height + cms.Margin.Top + cms.Margin.Bottom))
						if cms.Height+cms.Margin.Top+cms.Margin.Bottom == util.Roundi32(innerHeight) {
							needMeasure = false
						}
					}
				}

				if needMeasure {
					// log.I("nuxui", "measureVertical 2 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerHeight:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerHeight, canMeasureWidth, measuredFlags)

					newMeasuredFlags, _ := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, canMeasureWidth, measuredFlags)

					childrenMeasuredFlags[index] = newMeasuredFlags
				}
			}

			// Accumulate child.width that was not accumulated before to get the total value of vPx
			if nux.MeasureSpecMode(cms.Width) == nux.Pixel {
				hPx += float32(cms.Width)
			}

		}
	}

	// width mode is Auto, and child has percent size
	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		if hPt < 0 || hPt > 100 {
			log.Fatal("nuxui", "children percent size is out of range, it should 0% ~ 100%")
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
						cms.Width = nux.MeasureSpec(util.Roundi32(cw), nux.Pixel)
						setRatioHeight(cs, cms, w, nux.Pixel)
					}

					if measuredFlags := childrenMeasuredFlags[index]; measuredFlags != hMeasuredHeightComplete {
						// now width mode is must be nux.Pixel
						// log.I("nuxui", "measureVertical 3 width:%s, height:%s, vPPx:%.2f, vPPt:%.2f, innerHeight:%.2f, canMeasureHeight:%t, measuredFlags:%d", nux.MeasureSpecString(width), nux.MeasureSpecString(height), vPPx, vPPt, innerHeight, true, measuredFlags)
						newMeasuredFlags, _ := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, true, measuredFlags)

						if newMeasuredFlags != hMeasuredHeightComplete {
							log.Fatal("nuxui", "can not run here")
						}
					}
				}
			}
		}
	}

	// log.I("nuxui", "originWidth=%s, originHeight=%s", nux.MeasureSpecString(originWidth), nux.MeasureSpecString(originHeight))

	setNewWidth(ms, originWidth, width)
	setNewHeight(ms, originHeight, height)
}

// return vPx in innerHeight
func (me *row) measureVertical(width, height int32, vPPx, vPPt, innerHeight float32,
	child nux.Widget, canMeasureWidth bool, hasMeasuredFlags byte) (measuredFlags byte, vpxmax float32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd("nuxui", "ui.Column measureVertical", measuredDuration)
	var vWt float32
	var vPx float32
	var vPt float32
	var vPxUsed float32
	var vPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cms := cs.MeasuredSize()

		// if alrady measured complete, return
		if hasMeasuredFlags == hMeasuredHeightComplete {
			if vpx := cms.Height + cms.Margin.Top + cms.Margin.Bottom; vpx == nux.MeasureSpecValue(height) {
				return hasMeasuredFlags, float32(vpx)
			}
		}

		// 1. First determine the known size of height, add weight and percentage
		// do not add height to hPxUsed until height measured
		// do not add height to hPx until height measured
		if hasMeasuredFlags&hMeasuredHeight != hMeasuredHeight {
			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cms, h, nux.Pixel)
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					// do not add height to vWt until width measured
				case nux.Auto, nux.Unlimit:
					// wait until max height measured.
				}
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, h, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					// wait until max height measured
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

		// 2. First determine the known size of height, add weight and percentage
		if cs.HasMargin() {
			if hasMeasuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
				switch cs.MarginTop().Mode() {
				case nux.Pixel:
					t := cs.MarginTop().Value()
					cms.Margin.Top = util.Roundi32(t)
					vPx += t
					vPxUsed += t
					measuredFlags |= hMeasuredMarginTop
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cs.MarginTop().Value()
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						log.V("nuxui", "child:%s Margin Top not measured", child.ID())
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						t := cs.MarginTop().Value() / 100 * innerHeight
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
						vPxUsed += t
						measuredFlags |= hMeasuredMarginTop
						// ok
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						// ok
					}

				}
			}

			if hasMeasuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
				switch cs.MarginBottom().Mode() {
				case nux.Pixel:
					b := cs.MarginBottom().Value()
					cms.Margin.Bottom = int32(b)
					vPx += b
					vPxUsed += b
					measuredFlags |= hMeasuredMarginBottom
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cs.MarginBottom().Value()
					case nux.Auto, nux.Unlimit:
						// wait until max height measured.
						log.V("nuxui", "child:%s Margin Right not measured", child.ID())
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						b := cs.MarginBottom().Value() / 100 * innerHeight
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						measuredFlags |= hMeasuredMarginBottom
						// ok
					case nux.Auto, nux.Unlimit:
						// wait until max height measured.
						// ok
					}

				}
			}
		} else {
			measuredFlags |= hMeasuredMarginTop
			measuredFlags |= hMeasuredMarginBottom
		}

		vPxRemain = innerHeight - vPxUsed

		if vPxRemain < 0 {
			vPxRemain = 0
		}

		canMeasureHeight := true

		// Start calculating, only auto and weight are left here
		if hasMeasuredFlags&hMeasuredHeight != hMeasuredHeight {
			switch cs.Height().Mode() {
			// case nux.Pixel: has measured
			case nux.Weight:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					wt := cs.Height().Value() + vWt
					h := cs.Height().Value() / wt * vPxRemain
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, h, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					// wait until max height measured.
					canMeasureHeight = false
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				// case nux.Pixel: has measured
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured.
					canMeasureHeight = false
					// ok
				}
			// case nux.Ratio: has measured
			case nux.Auto, nux.Unlimit:
				cms.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode())
				setRatioWidth(cs, cms, vPxRemain, nux.Pixel)
				// ok
			}
		}

		log.V("nuxui", "child: %s, canMeasureWidth=%t, canMeasureHeight=%t", child.ID(), canMeasureWidth, canMeasureHeight)
		if canMeasureWidth && canMeasureHeight {
			if m, ok := child.(nux.Measure); ok {
				if hasMeasuredFlags&hMeasuredHeight != hMeasuredHeight {
					// log.V("nuxui", "child: %s, Measure width=%s, height=%s", child.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height))
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
					measuredFlags |= hMeasuredHeight
				}

				vPxRemain -= float32(cms.Height)
				vPx += float32(cms.Height)

				if cs.HasMargin() {
					if hasMeasuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
						switch cs.MarginTop().Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(height) {
							case nux.Pixel:
								if vWt > 0 && vPxRemain > 0 {
									cms.Margin.Top = util.Roundi32(cs.MarginTop().Value() / vWt * vPxRemain)
								} else {
									cms.Margin.Top = 0
								}
								vPx += float32(cms.Margin.Top)
								measuredFlags |= hMeasuredMarginTop
								// ok
							case nux.Auto, nux.Unlimit:
								// wait until max height measured.
								// ok
							}
						case nux.Percent:
							switch nux.MeasureSpecMode(height) {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								// wait until max height measured.
								// ok
							}
						}
					}

					if hasMeasuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
						switch cs.MarginBottom().Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(height) {
							case nux.Pixel:
								if vWt > 0 && vPxRemain > 0 {
									cms.Margin.Bottom = util.Roundi32(cs.MarginBottom().Value() / vWt * vPxRemain)
								} else {
									cms.Margin.Bottom = 0
								}
								vPx += float32(cms.Margin.Bottom)
								measuredFlags |= hMeasuredMarginBottom
								// ok
							case nux.Auto, nux.Unlimit:
								// wait until max height measured.
								// ok
							}
						case nux.Percent:
							switch nux.MeasureSpecMode(height) {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								// wait until max height measured.
								// ok
							}
						}
					}
				}

			}

			// parent width is auto or unlimit
			if vPt > 0 {
				if vPt > 100 {
					log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
				}

				vPx = vPx / (1.0 - vPt/100.0)
			}

			return measuredFlags, vPx
		}
		return measuredFlags, 0
	}

	return measuredFlags, 0
}

// Responsible for determining the position of the widget align, margin...
// TODO measure other mode dimen
func (me *row) Layout(dx, dy, left, top, right, bottom int32) {
	log.V("nuxui", "row:%s layout %d, %d, %d, %d, %d, %d", me.ID(), dx, dy, left, top, right, bottom)
	ms := me.MeasuredSize()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(bottom - top)
	var innerWidth float32 = float32(right - left)

	innerHeight -= float32(ms.Padding.Top + ms.Padding.Bottom)
	innerWidth -= float32(ms.Padding.Left + ms.Padding.Right)
	switch me.align.Horizontal {
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
				switch me.align.Vertical {
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
