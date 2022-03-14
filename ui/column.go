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
	nux.Layout
	nux.Measure
	nux.Draw
	Visual
}

type column struct {
	*nux.WidgetParent
	*nux.WidgetSize
	*WidgetVisual

	align          *Align // TODO not nil
	childrenHeight float32
}

func NewColumn(attrs ...nux.Attr) Column {
	attr := nux.MergeAttrs(attrs...)
	me := &column{
		align: NewAlign(attr.GetAttr("align", nux.Attr{})),
	}
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetParent = nux.NewWidgetParent(me, attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *column) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *column) Measure(width, height int32) {
	log.I("nuxui", "ui.Column %s Measure width=%s, height=%s", me.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Column %s Measure", me.Info().ID)

	var hPxMax float32 // max horizontal pixel size
	var hPPt float32   // horizontal padding percent
	var hPPx float32   // horizontal padding pixel
	var vPPt float32   // vertical padding percent
	var vPPx float32   // vvertical padding pixel

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

	frame := me.Frame()
	me.childrenHeight = 0

	// 1. Calculate its own padding size
	if p := me.Padding(); p != nil {
		if p.Left.Value() != 0 {
			switch p.Left.Mode() {
			case nux.Pixel:
				l := p.Left.Value()
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					l := p.Left.Value() / 100.0 * float32(nux.MeasureSpecValue(width))
					frame.Padding.Left = util.Roundi32(l)
					hPPx += l
					// ok
				case nux.Auto, nux.Unlimit:
					hPPt += p.Left.Value()
					// ok, wait until maxWidth measured
				}
			}
		}

		if p.Right.Value() != 0 {
			switch p.Right.Mode() {
			case nux.Pixel:
				r := p.Right.Value()
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					r := p.Right.Value() / 100.0 * float32(nux.MeasureSpecValue(width))
					frame.Padding.Right = util.Roundi32(r)
					hPPx += r
					// ok
				case nux.Auto, nux.Unlimit:
					hPPt += p.Right.Value()
					// ok, wait until maxWidth measured
				}
			}
		}

		if p.Top.Value() != 0 {
			switch p.Top.Mode() {
			case nux.Pixel:
				t := p.Top.Value()
				frame.Padding.Top = util.Roundi32(t)
				vPPx += t
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					t := p.Top.Value() / 100.0 * float32(nux.MeasureSpecValue(height))
					frame.Padding.Top = util.Roundi32(t)
					vPPx += t
					// ok
				case nux.Auto, nux.Unlimit:
					vPPt += p.Top.Value()
					// ok, wait until height measured
				}
			}
		}

		if p.Bottom.Value() != 0 {
			switch p.Bottom.Mode() {
			case nux.Pixel:
				b := p.Bottom.Value()
				frame.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					b := p.Bottom.Value() / 100.0 * float32(nux.MeasureSpecValue(height))
					frame.Padding.Bottom = util.Roundi32(b)
					vPPx += b
					// ok
				case nux.Auto, nux.Unlimit:
					vPPt += p.Bottom.Value()
					// ok, wait until height measured
				}
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
				if cm.Top.Value() != 0 {
					switch cm.Top.Mode() {
					case nux.Pixel:
						t := cm.Top.Value()
						cf.Margin.Top = util.Roundi32(t)
						vPx += t
						vPxUsed += t
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							vWt += cm.Top.Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							cf.Margin.Top = 0
							// ok
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							t := cm.Top.Value() / 100 * innerHeight
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							vPxUsed += t
							// ok
						case nux.Auto, nux.Unlimit:
							vPt += cm.Top.Value()
							// ok, wait until percent grand total
						}
					}
				}

				if cm.Bottom.Value() != 0 {
					switch cm.Bottom.Mode() {
					case nux.Pixel:
						b := cm.Bottom.Value()
						cf.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							vWt += cm.Bottom.Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							cf.Margin.Bottom = 0
							// ok
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							b := cm.Bottom.Value() / 100 * innerHeight
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							vPxUsed += b
							// ok
						case nux.Auto, nux.Unlimit:
							vPt += cm.Bottom.Value()
							// ok, wait until percent grand total
						}
					}
				}
			}

			canMeasureHeight := true

			// do not add height to vPxUsed until height measured
			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cf, h, nux.Pixel)
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					vWt += cs.Height().Value()
					// ok, wait until weight grand total
					if me.ChildrenCount() > 1 {
						canMeasureHeight = false
					}
				case nux.Auto, nux.Unlimit:
					cf.Height = nux.MeasureSpec(0, nux.Pixel)
					setRatioWidth(cs, cf, 0, nux.Pixel)
					// ok
				}
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cf, h, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					vPt += cs.Height().Value()
					// ok, wait until percent grand total
					if me.ChildrenCount() > 1 {
						canMeasureHeight = false
					}
				}
				// ok
			case nux.Ratio:
				if cs.Width().Mode() == nux.Ratio {
					log.Fatal("nuxui", "width and height size mode can not both Ratio")
				}
				// ok
			case nux.Auto, nux.Unlimit:
				h := innerHeight - vPxUsed
				cf.Height = nux.MeasureSpec(util.Roundi32(h), cs.Height().Mode())
				setRatioWidth(cs, cf, h, nux.Pixel)
				// ok
			}

			// log.I("nuxui", "measureHorizontal 1 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, canMeasureHeight, 0)

			measuredFlags, hPx := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, canMeasureHeight, 0)

			// add height to vPxUsed after height measured
			if nux.MeasureSpecMode(cf.Height) == nux.Pixel {
				vPxUsed += float32(cf.Height)
			}

			// find max innerWidth
			if hPx > hPxMax {
				hPxMax = hPx
			}

			childrenMeasuredFlags[index] = measuredFlags
		}
	}

	// log.TimeEnd(measureDuration_1lun, "nuxui", "ui.Column Measure measureDuration_1lun %s", me.Info().ID)

	// Use the maximum width found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)

		if nux.MeasureSpecValue(originWidth) != nux.MeasureSpecValue(width) {
			maxWidthChanged = true
		}

		if p := me.Padding(); p != nil {
			if p.Left.Value() != 0 && p.Left.Mode() == nux.Percent {
				l := p.Left.Value() / 100 * w
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
			}

			if p.Right.Value() != 0 && p.Right.Mode() == nux.Percent {
				r := p.Right.Value() / 100 * w
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
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
			cf := cs.Frame()

			canMeasureHeight := true

			// Vertical weight only works for Pixel Mode
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				// Vertical weight only works for nux.Pixel Mode
				if cm := cs.Margin(); cm != nil {
					if cm.Top.Value() != 0 && cm.Top.Mode() == nux.Weight {
						t := cm.Top.Value() / vWt * vPxRemain
						cf.Margin.Top = util.Roundi32(t)
						vPx += t
					}

					if cm.Bottom.Value() != 0 && cm.Bottom.Mode() == nux.Weight {
						b := cm.Bottom.Value() / vWt * vPxRemain
						cf.Margin.Bottom = util.Roundi32(b)
						vPx += b
					}
				}

				if cs.Height().Mode() == nux.Weight {
					if vWt > 0 && vPxRemain > 0 {
						h := cs.Height().Value() / vWt * vPxRemain
						cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cf, h, nux.Pixel)
					} else {
						cf.Height = nux.MeasureSpec(0, nux.Pixel)
						setRatioWidth(cs, cf, 0, nux.Pixel)
					}
				}
			case nux.Auto, nux.Unlimit:
				if cs.Height().Mode() == nux.Percent {
					// wait all weight size finish measured
					if me.ChildrenCount() > 1 {
						canMeasureHeight = false
					}
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
						log.I("nuxui", "hMeasuredWidthComplete child:%s hpx=%d", child.Info().ID, (cf.Width + cf.Margin.Left + cf.Margin.Right))
						if cf.Width+cf.Margin.Left+cf.Margin.Right == util.Roundi32(innerWidth) {
							needMeasure = false
						}
					}
				}

				if needMeasure {
					// log.I("nuxui", "measureHorizontal 2 child:%s, width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", child.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, canMeasureHeight, measuredFlags)

					newMeasuredFlags, _ := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, child, canMeasureHeight, measuredFlags)

					childrenMeasuredFlags[index] = newMeasuredFlags
				}
			}

			// Accumulate child.height that was not accumulated before to get the total value of vPx
			if nux.MeasureSpecMode(cf.Height) == nux.Pixel {
				vPx += float32(cf.Height)
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
			if p := me.Padding(); p != nil {
				if p.Top.Value() != 0 && p.Top.Mode() == nux.Percent {
					t := p.Top.Value() / 100.0 * h
					frame.Padding.Top = util.Roundi32(t)
				}

				if p.Bottom.Value() != 0 && p.Bottom.Mode() == nux.Percent {
					b := p.Bottom.Value() / 100.0 * h
					frame.Padding.Bottom = util.Roundi32(b)
				}
			}
		}

		if vPt > 0 {
			for index, child := range me.Children() {
				if compt, ok := child.(nux.Component); ok {
					child = compt.Content()
				}

				if cs, ok := child.(nux.Size); ok {
					cf := cs.Frame()

					if cm := cs.Margin(); cm != nil {
						if cm.Top.Value() != 0 && cm.Top.Mode() == nux.Percent {
							t := cm.Top.Value() / 100.0 * innerHeight
							cf.Margin.Top = util.Roundi32(t)
						}

						if cm.Bottom.Value() != 0 && cm.Bottom.Mode() == nux.Percent {
							b := cm.Bottom.Value() / 100.0 * innerHeight
							cf.Margin.Bottom = util.Roundi32(b)
						}
					}

					if cs.Height().Mode() == nux.Percent {
						ch := cs.Height().Value() / 100 * innerHeight
						cf.Height = nux.MeasureSpec(util.Roundi32(ch), nux.Pixel)
						setRatioWidth(cs, cf, h, nux.Pixel)
					}

					if measuredFlags := childrenMeasuredFlags[index]; measuredFlags != hMeasuredWidthComplete {
						// now width mode is must be nux.Pixel
						// log.I("nuxui", "measureHorizontal 3 width:%s, height:%s, hPPx:%.2f, hPPt:%.2f, innerWidth:%.2f, canMeasureHeight:%t, measuredFlags:%d", nux.MeasureSpecString(width), nux.MeasureSpecString(height), hPPx, hPPt, innerWidth, true, measuredFlags)
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

	setNewWidth(frame, originWidth, width)
	setNewHeight(frame, originHeight, height)
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
		cf := cs.Frame()

		// if alrady measured complete, return
		if hasMeasuredFlags == hMeasuredWidthComplete {
			if hpx := cf.Width + cf.Margin.Left + cf.Margin.Right; hpx == nux.MeasureSpecValue(width) {
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
				cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cf, w, nux.Pixel)
				// ok
			case nux.Weight:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					// do not add width to hWt until width measured
				case nux.Auto, nux.Unlimit:
					// wait until maxWidth measured.
				}
				// ok
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cf, w, nux.Pixel)
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
		if cm := cs.Margin(); cm != nil {
			if hasMeasuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
				if cm.Left.Value() != 0 {
					switch cm.Left.Mode() {
					case nux.Pixel:
						l := cm.Left.Value()
						cf.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						measuredFlags |= hMeasuredMarginLeft
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							hWt += cm.Left.Value()
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
							log.V("nuxui", "child:%s Margin Left not measured", child.Info().ID)
							// ok
						}
						// ok
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							l := cm.Left.Value() / 100 * innerWidth
							cf.Margin.Left = util.Roundi32(l)
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
			}

			if hasMeasuredFlags&hMeasuredMarginRight != hMeasuredMarginRight {
				if cm.Right.Value() != 0 {
					switch cm.Right.Mode() {
					case nux.Pixel:
						r := cm.Right.Value()
						cf.Margin.Right = int32(r)
						hPx += r
						hPxUsed += r
						measuredFlags |= hMeasuredMarginRight
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							hWt += cm.Right.Value()
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
							log.V("nuxui", "child:%s Margin Right not measured", child.Info().ID)
							// ok
						}
						// ok
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							r := cm.Right.Value() / 100 * innerWidth
							cf.Margin.Right = util.Roundi32(r)
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
					wt := cs.Width().Value() + hWt
					w := cs.Width().Value() / wt * hPxRemain
					cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cf, w, nux.Pixel)
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
				cf.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode())
				setRatioHeight(cs, cf, hPxRemain, nux.Pixel)
				// ok
			}
		}

		log.V("nuxui", "child: %s, canMeasureWidth=%t, canMeasureHeight=%t", child.Info().ID, canMeasureWidth, canMeasureHeight)
		if canMeasureWidth && canMeasureHeight {
			if m, ok := child.(nux.Measure); ok {
				if hasMeasuredFlags&hMeasuredWidth != hMeasuredWidth {
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
					measuredFlags |= hMeasuredWidth
				}

				hPxRemain -= float32(cf.Width)
				hPx += float32(cf.Width)

				if cm := cs.Margin(); cm != nil {
					if hasMeasuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
						if cm.Left.Value() != 0 {
							switch cm.Left.Mode() {
							// case nux.Pixel: has measured
							case nux.Weight:
								switch nux.MeasureSpecMode(width) {
								case nux.Pixel:
									if hWt > 0 && hPxRemain > 0 {
										cf.Margin.Left = util.Roundi32(cm.Left.Value() / hWt * hPxRemain)
									} else {
										cf.Margin.Left = 0
									}
									hPx += float32(cf.Margin.Left)
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
					}

					if hasMeasuredFlags&hMeasuredMarginRight != hMeasuredMarginRight {
						if cm.Right.Value() != 0 {
							switch cm.Right.Mode() {
							// case nux.Pixel: has measured
							case nux.Weight:
								switch nux.MeasureSpecMode(width) {
								case nux.Pixel:
									if hWt > 0 && hPxRemain > 0 {
										cf.Margin.Right = util.Roundi32(cm.Right.Value() / hWt * hPxRemain)
									} else {
										cf.Margin.Right = 0
									}
									hPx += float32(cf.Margin.Right)
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
func (me *column) Layout(x, y, width, height int32) {
	log.D("nuxui", "column layout %d, %d, %d, %d", x, y, width, height)
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(width)
	var innerWidth float32 = float32(height)

	innerHeight -= float32(frame.Padding.Top + frame.Padding.Bottom)
	innerWidth -= float32(frame.Padding.Left + frame.Padding.Right)
	switch me.align.Vertical {
	case Top:
		t = 0
	case Center:
		t = innerHeight/2 - me.childrenHeight/2
	case Bottom:
		t = innerHeight - me.childrenHeight
	}
	t += float32(frame.Padding.Top)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		l = 0
		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()
			cm := cs.Margin()
			if cs.Width().Mode() == nux.Weight || (cm != nil && (cm.Left.Mode() == nux.Weight || cm.Right.Mode() == nux.Weight)) {
				l += float32(frame.Padding.Left + cf.Margin.Left)
			} else {
				switch me.align.Horizontal {
				case Left:
					l += 0
				case Center:
					l += innerWidth/2 - float32(cf.Width)/2
				case Right:
					l += innerWidth - float32(cf.Width)
				}
				l += float32(frame.Padding.Left + cf.Margin.Left)
			}

			t += float32(cf.Margin.Top)

			cf.X = x + int32(l)
			cf.Y = y + int32(t)
			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cf.X, cf.Y, cf.Width, cf.Height)
			}

			t += float32(cf.Height + cf.Margin.Bottom)
		}
	}
}

func (me *column) Draw(canvas nux.Canvas) {
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
