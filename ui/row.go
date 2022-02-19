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

func NewRow(attrs ...nux.Attr) Row {
	attr := nux.MergeAttrs(attrs...)
	me := &row{
		align: NewAlign(attr.GetAttr("align", nux.Attr{})),
	}
	me.WidgetParent = nux.NewWidgetParent(me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *row) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *row) Measure(width, height int32) {
	// log.I("nuxui", "ui.Row %s Measure width=%s, height=%s", me.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Row %s Measure", me.Info().ID)

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

	frame := me.Frame()
	me.childrenWidth = 0

	// 1. Calculate its own padding size
	if me.Padding() != nil {
		switch me.Padding().Left.Mode() {
		case nux.Pixel:
			l := me.Padding().Left.Value()
			frame.Padding.Left = util.Roundi32(l)
			hPPx += l
			hPPxUsed += l
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.Padding().Left.Value() / 100.0 * float32(nux.MeasureSpecValue(width))
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.Padding().Left.Value()
				// ok, wait until maxWidth measured
			}
		}

		switch me.Padding().Right.Mode() {
		case nux.Pixel:
			r := me.Padding().Right.Value()
			frame.Padding.Right = util.Roundi32(r)
			hPPx += r
			hPPxUsed += r
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.Padding().Right.Value() / 100.0 * float32(nux.MeasureSpecValue(width))
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.Padding().Right.Value()
				// ok, wait until maxWidth measured
			}
		}

		switch me.Padding().Top.Mode() {
		case nux.Pixel:
			t := me.Padding().Top.Value()
			frame.Padding.Top = util.Roundi32(t)
			vPPx += t
			vPPxUsed += t
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.Padding().Top.Value() / 100.0 * float32(nux.MeasureSpecValue(height))
				frame.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
				// ok
			case nux.Auto, nux.Unlimit:
				vPPt += me.Padding().Top.Value()
				// ok, wait until height measured
			}
		}

		switch me.Padding().Bottom.Mode() {
		case nux.Pixel:
			b := me.Padding().Bottom.Value()
			frame.Padding.Bottom = util.Roundi32(b)
			vPPx += b
			vPPxUsed += b
			// ok
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.Padding().Bottom.Value() / 100.0 * float32(nux.MeasureSpecValue(height))
				frame.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
				// ok
			case nux.Auto, nux.Unlimit:
				vPPt += me.Padding().Bottom.Value()
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
			cf := cs.Frame()

			// nux.Unlimit ???????
			cf.Width = nux.MeasureSpec(0, nux.Unlimit)
			cf.Height = nux.MeasureSpec(0, nux.Unlimit)

			if cm := cs.Margin(); cm != nil {
				switch cm.Left.Mode() {
				case nux.Pixel:
					l := cm.Left.Value()
					cf.Margin.Left = util.Roundi32(l)
					hPx += l
					hPxUsed += l
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cm.Left.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cf.Margin.Left = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						l := cm.Left.Value() / 100 * innerWidth
						cf.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cm.Left.Value()
						// ok, wait until percent grand total
					}
				}

				switch cm.Right.Mode() {
				case nux.Pixel:
					r := cm.Right.Value()
					cf.Margin.Right = util.Roundi32(r)
					hPx += r
					hPxUsed += r
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cm.Right.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						cf.Margin.Right = 0
						// ok
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						r := cm.Right.Value() / 100 * innerWidth
						cf.Margin.Right = util.Roundi32(r)
						hPx += r
						hPxUsed += r
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cm.Right.Value()
						// ok, wait until percent grand total
					}
				}
			}

			canMeasureWidth := true

			// do not add width to hPxUsed until width measured
			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cf, w, nux.Pixel)
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
					cf.Width = nux.MeasureSpec(0, nux.Pixel)
					setRatioHeight(cs, cf, 0, nux.Pixel)
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cf, w, nux.Pixel)
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
				cf.Width = nux.MeasureSpec(util.Roundi32(w), cs.Width().Mode())
				setRatioHeight(cs, cf, w, nux.Pixel)
				// ok
			}

			log.I("nuxui", "measureVertical 1 child:%s, width:%s, height:%s, vPPx:%.2f, vPPt:%.2f, innerHeight:%.2f, canMeasureWidth:%t, measuredFlags:%d", child.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height), vPPx, vPPt, innerHeight, canMeasureWidth, 0)

			measuredFlags, vPx := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, canMeasureWidth, 0)

			// add width to hPxUsed after width measured
			if nux.MeasureSpecMode(cf.Width) == nux.Pixel {
				hPxUsed += float32(cf.Width)
			}

			// find max innerHeight
			if vPx > vPxMax {
				vPxMax = vPx
			}

			childrenMeasuredFlags[index] = measuredFlags
		}
	}

	// log.TimeEnd(measureDuration_1lun, "nuxui", "ui.Column Measure measureDuration_1lun %s", me.Info().ID)

	// Use the maximum height found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(height) {
	case nux.Auto, nux.Unlimit:
		innerHeight = vPxMax
		h := (innerHeight + vPPx) / (1 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)

		if nux.MeasureSpecValue(originHeight) != nux.MeasureSpecValue(height) {
			maxHeightChanged = true
		}

		if me.Padding() != nil {
			if me.Padding().Top.Mode() == nux.Percent {
				t := me.Padding().Top.Value() / 100 * h
				frame.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
			}

			if me.Padding().Bottom.Mode() == nux.Percent {
				b := me.Padding().Bottom.Value() / 100 * h
				frame.Padding.Bottom = util.Roundi32(b)
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
			cf := cs.Frame()

			canMeasureWidth := true

			// Vertical weight only works for Pixel Mode
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				// Vertical weight only works for nux.Pixel Mode
				if cm := cs.Margin(); cm != nil {
					if cm.Left.Mode() == nux.Weight {
						l := cm.Left.Value() / hWt * hPxRemain
						cf.Margin.Left = util.Roundi32(l)
						hPx += l
					}

					if cm.Right.Mode() == nux.Weight {
						r := cm.Right.Value() / hWt * hPxRemain
						cf.Margin.Right = util.Roundi32(r)
						hPx += r
					}
				}

				if cs.Width().Mode() == nux.Weight {
					if hWt > 0 && hPxRemain > 0 {
						w := cs.Width().Value() / hWt * hPxRemain
						cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cf, w, nux.Pixel)
					} else {
						cf.Width = nux.MeasureSpec(0, nux.Pixel)
						setRatioHeight(cs, cf, 0, nux.Pixel)
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
						log.I("nuxui", "hMeasuredHeightComplete child:%s vpx=%d", child.Info().ID, (cf.Height + cf.Margin.Top + cf.Margin.Bottom))
						if cf.Height+cf.Margin.Top+cf.Margin.Bottom == util.Roundi32(innerHeight) {
							needMeasure = false
						}
					}
				}

				if needMeasure {
					// log.I("nuxui", "measureVertical 2 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerHeight:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerHeight, canMeasureWidth, measuredFlags)

					newMeasuredFlags, _ := me.measureVertical(width, height, vPPx, vPPt, innerHeight, child, canMeasureWidth, measuredFlags)

					childrenMeasuredFlags[index] = newMeasuredFlags
				}
			}

			// Accumulate child.width that was not accumulated before to get the total value of vPx
			if nux.MeasureSpecMode(cf.Width) == nux.Pixel {
				hPx += float32(cf.Width)
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
			if me.Padding() != nil {
				if me.Padding().Left.Mode() == nux.Percent {
					l := me.Padding().Left.Value() / 100.0 * w
					frame.Padding.Left = util.Roundi32(l)
				}

				if me.Padding().Right.Mode() == nux.Percent {
					r := me.Padding().Right.Value() / 100.0 * w
					frame.Padding.Right = util.Roundi32(r)
				}
			}
		}

		if hPt > 0 {
			for index, child := range me.Children() {
				if compt, ok := child.(nux.Component); ok {
					child = compt.Content()
				}

				if cs, ok := child.(nux.Size); ok {
					cf := cs.Frame()

					if cm := cs.Margin(); cm != nil {
						if cm.Left.Mode() == nux.Percent {
							l := cm.Left.Value() / 100.0 * innerWidth
							cf.Margin.Left = util.Roundi32(l)
						}

						if cm.Right.Mode() == nux.Percent {
							r := cm.Right.Value() / 100.0 * innerWidth
							cf.Margin.Right = util.Roundi32(r)
						}
					}

					if cs.Width().Mode() == nux.Percent {
						cw := cs.Width().Value() / 100 * innerWidth
						cf.Width = nux.MeasureSpec(util.Roundi32(cw), nux.Pixel)
						setRatioHeight(cs, cf, w, nux.Pixel)
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

	setNewWidth(frame, originWidth, width)
	setNewHeight(frame, originHeight, height)
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
		cf := cs.Frame()

		// if alrady measured complete, return
		if hasMeasuredFlags == hMeasuredHeightComplete {
			if vpx := cf.Height + cf.Margin.Top + cf.Margin.Bottom; vpx == nux.MeasureSpecValue(height) {
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
				cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cf, h, nux.Pixel)
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
					cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cf, h, nux.Pixel)
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
		if cm := cs.Margin(); cm != nil {
			if hasMeasuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
				switch cm.Top.Mode() {
				case nux.Pixel:
					t := cm.Top.Value()
					cf.Margin.Top = util.Roundi32(t)
					vPx += t
					vPxUsed += t
					measuredFlags |= hMeasuredMarginTop
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cm.Top.Value()
					case nux.Auto, nux.Unlimit:
						// wait until maxWidth measured.
						log.V("nuxui", "child:%s Margin Top not measured", child.Info().ID)
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						t := cm.Top.Value() / 100 * innerHeight
						cf.Margin.Top = util.Roundi32(t)
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
				switch cm.Bottom.Mode() {
				case nux.Pixel:
					b := cm.Bottom.Value()
					cf.Margin.Bottom = int32(b)
					vPx += b
					vPxUsed += b
					measuredFlags |= hMeasuredMarginBottom
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vWt += cm.Bottom.Value()
					case nux.Auto, nux.Unlimit:
						// wait until max height measured.
						log.V("nuxui", "child:%s Margin Right not measured", child.Info().ID)
						// ok
					}
					// ok
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						b := cm.Bottom.Value() / 100 * innerHeight
						cf.Margin.Bottom = util.Roundi32(b)
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
					cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cf, h, nux.Pixel)
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
				cf.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode())
				setRatioWidth(cs, cf, vPxRemain, nux.Pixel)
				// ok
			}
		}

		log.V("nuxui", "child: %s, canMeasureWidth=%t, canMeasureHeight=%t", child.Info().ID, canMeasureWidth, canMeasureHeight)
		if canMeasureWidth && canMeasureHeight {
			if m, ok := child.(nux.Measure); ok {
				if hasMeasuredFlags&hMeasuredHeight != hMeasuredHeight {
					// log.V("nuxui", "child: %s, Measure width=%s, height=%s", child.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height))
					m.Measure(cf.Width, cf.Height)

					if cs.Width().Mode() == nux.Ratio {
						oldWidth := cf.Width
						cf.Width = nux.MeasureSpec(util.Roundi32(float32(cf.Height)*cs.Width().Value()), nux.Pixel)
						if oldWidth != cf.Width {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(cf.Width, cf.Height)
						}
					}

					if cs.Height().Mode() == nux.Ratio {
						oldHeight := cf.Height
						cf.Height = nux.MeasureSpec(util.Roundi32(float32(cf.Width)/cs.Height().Value()), nux.Pixel)
						if oldHeight != cf.Height {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(cf.Width, cf.Height)
						}
					}
					measuredFlags |= hMeasuredHeight
				}

				vPxRemain -= float32(cf.Height)
				vPx += float32(cf.Height)

				if cm := cs.Margin(); cm != nil {
					if hasMeasuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
						switch cm.Top.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(height) {
							case nux.Pixel:
								if vWt > 0 && vPxRemain > 0 {
									cf.Margin.Top = util.Roundi32(cm.Top.Value() / vWt * vPxRemain)
								} else {
									cf.Margin.Top = 0
								}
								vPx += float32(cf.Margin.Top)
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
						switch cm.Bottom.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch nux.MeasureSpecMode(height) {
							case nux.Pixel:
								if vWt > 0 && vPxRemain > 0 {
									cf.Margin.Bottom = util.Roundi32(cm.Bottom.Value() / vWt * vPxRemain)
								} else {
									cf.Margin.Bottom = 0
								}
								vPx += float32(cf.Margin.Bottom)
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
func (me *row) Layout(x, y, width, height int32) {
	log.V("nuxui", "row:%s layout %d, %d, %d, %d", me.Info().ID, x, y, width, height)
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(height)
	var innerWidth float32 = float32(width)

	innerHeight -= float32(frame.Padding.Top + frame.Padding.Bottom)
	innerWidth -= float32(frame.Padding.Left + frame.Padding.Right)
	switch me.align.Horizontal {
	case Left:
		l = 0
	case Center:
		l = innerWidth/2 - me.childrenWidth/2
	case Right:
		l = innerWidth - me.childrenWidth
	}
	l += float32(frame.Padding.Left)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		t = 0
		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()
			cm := cs.Margin()
			if cs.Height().Mode() == nux.Weight || (cm != nil && (cm.Top.Mode() == nux.Weight || cm.Bottom.Mode() == nux.Weight)) {
				t += float32(frame.Padding.Top + cf.Margin.Top)
			} else {
				switch me.align.Vertical {
				case Top:
					t += 0
				case Center:
					t += innerHeight/2 - float32(cf.Height)/2
				case Bottom:
					t += innerHeight - float32(cf.Height)
				}
				t += float32(frame.Padding.Top + cf.Margin.Top)
			}

			l += float32(cf.Margin.Left)

			cf.X = x + int32(l)
			cf.Y = y + int32(t)
			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cf.X, cf.Y, cf.Width, cf.Height)
			}

			l += float32(cf.Width + cf.Margin.Right)
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

		if draw, ok := child.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
