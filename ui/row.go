// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	"nuxui.org/nuxui/util"
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

	align          *Align
	childrenWidth  float32
	childrenHeight float32
	clipChildren   int
}

func NewRow(attr nux.Attr) Row {
	me := &row{
		align: NewAlign(attr.GetAttr("align", nux.Attr{})),
	}

	if attr.Has("clipChildren") {
		v := attr.GetBool("clipChildren", false)
		if v {
			me.clipChildren = clipChildrenYes
		} else {
			me.clipChildren = clipChildrenNo
		}
	} else {
		me.clipChildren = clipChildrenAuto
	}

	me.WidgetParent = nux.NewWidgetParent(me, attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *row) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *row) Measure(width, height nux.MeasureDimen) {
	frame := me.Frame()

	me.childrenWidth = 0
	me.childrenHeight = 0
	childrenMeasuredFlags := make([]uint8, len(me.Children()))

	hPPx, hPPt, vPPx, vPPt, paddingMeasuredFlag := measurePadding(width, height, me.Padding(), frame, -1, 0)
	if hPPt >= 100.0 || vPPt >= 100.0 {
		log.Fatal("nuxui", "padding percent size should at 0% ~ 100%")
	}

	innerWidth := float32(width.Value())*(1.0-hPPt/100.0) - hPPx
	innerHeight := float32(height.Value())*(1.0-vPPt/100.0) - vPPx

	if innerWidth <= 0 || innerHeight <= 0 {
		measureChildrenZeroSize(me)
		frame.Width = util.Roundi32(hPPx / (1 - hPPt/100.0))
		frame.Height = util.Roundi32(vPPx / (1 - vPPt/100.0))
		return
	}

	hPx, hMPt, hMWt, hHWt, hPxMax, hasMaxWidth, hasAutoWeightWidth, vPxMax, measuredHeightCount, measuredCompleteCount := me.measureChildren1(width, height, hPPx, hPPt, innerWidth, innerHeight, childrenMeasuredFlags)

	if hMPt > 100.0 {
		log.Fatal("nuxui", "the sum of percent size should at 0% ~ 100%")
	}

	switch height.Mode() {
	case nux.Pixel:
		frame.Height = height.Value()
		me.childrenHeight = vPxMax
	case nux.Auto, nux.Unlimit:
		if measuredHeightCount == me.ChildrenCount() {
			if height.Mode() == nux.Auto {
				ih := float32(height.Value())*(1-vPPt/100.0) - vPPx
				if vPxMax > ih {
					vPxMax = ih
				}
			}
			innerHeight = vPxMax
			h := (innerHeight + vPPx) / (1 - vPPt/100.0)
			height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
			frame.Height = util.Roundi32(h)
			me.childrenHeight = vPxMax
		}
	}

	switch width.Mode() {
	case nux.Pixel:
		frame.Width = width.Value()
		me.childrenWidth = hPx
	case nux.Auto, nux.Unlimit:
		if hasMaxWidth || measuredCompleteCount == me.ChildrenCount() {
			if width.Mode() == nux.Auto {
				iw := float32(width.Value())*(1-hPPt/100.0) - hPPx
				if hPxMax > iw {
					hPxMax = iw
				}
			}
			innerWidth = hPxMax
			w := (innerWidth + hPPx) / (1 - hPPt/100.0)
			width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
			frame.Width = util.Roundi32(w)
			me.childrenWidth = hPxMax
		}
	}

	// if measure finished at first traversal
	if measuredCompleteCount == me.ChildrenCount() {
		if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
			measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
		}
		return
	}

	hPxRemain := innerWidth - hPx
	if hPxRemain < 0 {
		hPxRemain = 0
	}

	hPx, hasPerWidthWtPx, perWidthWtPx, hPxMax, hasMaxWidth, vPxMax, measuredCount, measuredCompleteCount := me.measureChildren2(width, height, hPPx, hPPt, innerWidth, innerHeight, hPx, hPxRemain, hMWt, hHWt, hMPt, hPxMax, hasMaxWidth, hasAutoWeightWidth, false, 0, vPxMax, childrenMeasuredFlags)

	if measuredCount != me.ChildrenCount() {
		log.Fatal("nuxui", "can not run here")
	}

	switch height.Mode() {
	// case nux.Pixel: handled
	case nux.Auto, nux.Unlimit:
		innerHeight = vPxMax
		h := (innerHeight + vPPx) / (1 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
		frame.Height = util.Roundi32(h)
		me.childrenHeight = vPxMax // TODO:: exclude margin
	}

	switch width.Mode() {
	// case nux.Pixel: handled
	case nux.Auto, nux.Unlimit:
		if hasMaxWidth {
			innerWidth = hPxMax
			w := (innerWidth + hPPx) / (1 - hPPt/100.0)
			width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
			frame.Width = util.Roundi32(w)
			me.childrenWidth = hPxMax
		} else {
			w := (hPx + hPPx) / (1 - hPPt/100.0)
			frame.Width = util.Roundi32(w)
			me.childrenWidth = hPx
		}
	}

	if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
	}

	if measuredCompleteCount == me.ChildrenCount() {
		return
	}

	// third measure
	hPx, _, _, _, _, vPxMax, _, measuredCompleteCount = me.measureChildren2(width, height, hPPx, hPPt, innerWidth, innerHeight, hPx, hPxRemain, hMWt, hHWt, hMPt, hPxMax, false, false, hasPerWidthWtPx, perWidthWtPx, vPxMax, childrenMeasuredFlags)

	if measuredCompleteCount == me.ChildrenCount() {
		return
	}

	log.Fatal("nuxui", "row %s not measure complete.", me.Info().ID)
}

/*
hPx: sum of horizontal pixel size
hMWt: sum of children horizontal margin weight
hWt: sum of children horizontal height weight
hPt: sum of horizontal percent size
*/
func (me *row) measureChildren1(width, height nux.MeasureDimen, hPPx, hPPt, innerWidth, innerHeight float32, childrenMeasuredFlags []uint8) (hPx, hMPt, hMWt, hHWt, hPxMax float32, hasMaxWidth, hasAutoWeightWidth bool, vPxMax float32, measuredHeightCount, measuredCompleteCount int) {
	var measuredFlags uint8

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			// the child size mode will be Pixel after measured
			cf.Clear()
			cf.Width = int32(nux.MeasureSpec(0, nux.Unlimit))
			cf.Height = int32(nux.MeasureSpec(0, nux.Unlimit))

			measuredFlags = childrenMeasuredFlags[index]

			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					if cm.Left.Value() != 0 {
						switch cm.Left.Mode() {
						case nux.Pixel:
							l := cm.Left.Value()
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Weight:
							hMWt += cm.Left.Value()
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								l := cm.Left.Value() / 100 * innerWidth
								cf.Margin.Left = util.Roundi32(l)
								hPx += l
								measuredFlags |= flagMeasuredMarginLeft
							case nux.Auto, nux.Unlimit:
								// wait until percent grand total
								hMPt += cm.Left.Value()
							}
						}
					} else {
						measuredFlags |= flagMeasuredMarginLeft
					}
				}

				if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
					if cm.Right.Value() != 0 {
						switch cm.Right.Mode() {
						case nux.Pixel:
							r := cm.Right.Value()
							cf.Margin.Right = util.Roundi32(r)
							hPx += r
							measuredFlags |= flagMeasuredMarginRight
						case nux.Weight:
							hMWt += cm.Right.Value()
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								r := cm.Right.Value() / 100 * innerWidth
								cf.Margin.Right = util.Roundi32(r)
								hPx += r
								measuredFlags |= flagMeasuredMarginRight
							case nux.Auto, nux.Unlimit:
								// wait until percent grand total
								hMPt += cm.Right.Value()
							}
						}
					} else {
						measuredFlags |= flagMeasuredMarginRight
					}
				}
			} else {
				measuredFlags |= flagMeasuredMarginLeft | flagMeasuredMarginRight
			}

			canMeasureWidth := true
			// forceMeasure := false

			// Do not add height to vPx and vPxUsed until height measured
			if measuredFlags&flagMeasuredWidth != flagMeasuredWidth {
				if cs.Width().Value() > 0 || cs.Width().Mode() == nux.Auto || cs.Width().Mode() == nux.Unlimit {
					switch cs.Width().Mode() {
					case nux.Pixel:
						w := cs.Width().Value()
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
						setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
						measuredFlags |= flagMeasuredWidth
					case nux.Weight:
						// wait until weight grand total and use height.Mode()
						hHWt += cs.Width().Value()
						canMeasureWidth = false
						switch width.Mode() {
						case nux.Auto, nux.Unlimit:
							hasAutoWeightWidth = true
						}
					case nux.Percent:
						switch width.Mode() {
						case nux.Pixel:
							w := cs.Width().Value() / 100 * innerWidth
							cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
							setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
							measuredFlags |= flagMeasuredWidth
						case nux.Auto, nux.Unlimit:
							w := cs.Width().Value() / 100 * innerWidth
							cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
							setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
						}
					case nux.Ratio:
						if cs.Height().Mode() == nux.Ratio {
							log.Fatal("nuxui", "width and height size mode can not both Ratio")
						}
					case nux.Auto, nux.Unlimit:
						w := innerWidth - hPx
						if w < 0 {
							w = 0
						}
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), cs.Width().Mode()))
						setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
					}
				} else {
					cf.Width = 0
					setRatioHeightIfNeed(cs, cf, 0, nux.Pixel)
					measuredFlags |= flagMeasuredWidth
				}
			}

			measuredFlags, vPx := me.measureVertical(width, height, hPPx, hPPt, innerWidth, canMeasureWidth, child, measuredFlags)

			if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
				measuredCompleteCount++
			}

			if measuredFlags&flagMeasuredHeight == flagMeasuredHeight {
				measuredHeightCount++
			}

			if nux.MeasureDimen(cf.Width).Mode() == nux.Pixel {
				hPx += float32(cf.Width)

				// find max innerWidth in this mode
				if width.Mode() != nux.Pixel && cs.Width().Mode() == nux.Percent {
					hasMaxWidth = true
					if cs.Width().Value() > 0 {
						h := float32(cf.Width) / (cs.Width().Value() / 100.0)
						if h > hPxMax {
							hPxMax = h
						}
					}
				}
			}

			if vPx > vPxMax {
				vPxMax = vPx
			}
			if height.Mode() == nux.Auto {
				if vPxMax >= innerHeight {
					vPxMax = innerHeight
					height = nux.MeasureSpec(height.Value(), nux.Pixel)
				}
			}

			childrenMeasuredFlags[index] = measuredFlags
		} else {
			childrenMeasuredFlags[index] = flagMeasuredComplete
			measuredHeightCount++
			measuredCompleteCount++
		}
	}

	h := hPx / (1 - hMPt/100.0)
	if h > hPxMax {
		hPxMax = h
	}
	return
}

func (me *row) measureChildren2(width, height nux.MeasureDimen, hPPx, hPPt, innerWidth, innerHeight float32,
	hPx, hPxRemain, hMWt, hHWt, hMPt, hPxMax float32, hasMaxWidth, hasAutoWeightWidth, hasPerWidthWtPx bool, perWidthWtPx, vPxMax float32, childrenMeasuredFlags []uint8) (hPxOut float32, hasPerWidthWtPxOut bool, perWidthWtPxOut float32, hPxMaxOut float32, hasMaxWidthOut bool, vPxOut float32, measuredCount, measuredCompleteCount int) {
	var measuredFlags uint8
	var hHWtPx, vPx float32 = 0, 0

	if hasMaxWidth && hMPt > 0 {
		hPxRemain -= hPxMax * hMPt
	}

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			measuredFlags = childrenMeasuredFlags[index]

			// Do not add height to vPx and vPxUsed until height measured
			if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
				if cs.Width().Value() > 0 {
					switch cs.Width().Mode() {
					// case nux.Pixel: already handled
					// case nux.Percent: already handled
					// case nux.Ratio: already handled
					// case nux.Auto, nux.Unlimit: already handled
					case nux.Weight:
						w := hPxRemain * cs.Width().Value() / (hMWt + hHWt)
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
						setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
					}
				}
			}

			flag := flagMeasured | flagMeasuredWidth | flagMeasuredVerticalComplete
			if measuredFlags&flag != flag {
				measuredFlags, vPx = me.measureVertical(width, height, hPPx, hPPt, innerHeight, true, child, measuredFlags)

				if vPx > vPxMax {
					vPxMax = vPx
				}

				if nux.MeasureDimen(cf.Width).Mode() == nux.Pixel {
					hPx += float32(cf.Width)

					if cs.Width().Mode() == nux.Weight {
						hHWtPx += float32(cf.Width)
					}
				}
			}

			// wait third measure after auto weight height mode measured max height
			if !hasAutoWeightWidth {
				if cm := cs.Margin(); cm != nil {
					if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
						if cm.Left.Value() != 0 {
							switch cm.Left.Mode() {
							// case nux.Pixel: already handled
							case nux.Weight:
								var l float32
								if hasPerWidthWtPx {
									l = cm.Left.Value() * perWidthWtPx
								} else {
									l = cm.Left.Value() / (hMWt + hHWt) * hPxRemain
								}
								hPx += l
								cf.Margin.Left = util.Roundi32(l)
								measuredFlags |= flagMeasuredMarginLeft
							case nux.Percent:
								if hasMaxWidth {
									l := hPxMax * cs.Margin().Left.Value()
									hPx += l
									cf.Margin.Left = util.Roundi32(l)
									measuredFlags |= flagMeasuredMarginLeft
								}
							}
						} else {
							measuredFlags |= flagMeasuredMarginLeft
						}
					}

					if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
						if cm.Right.Value() != 0 {
							switch cm.Right.Mode() {
							// case nux.Pixel: // already measured
							case nux.Weight:
								var r float32
								if hasPerWidthWtPx {
									r = cm.Right.Value() * perWidthWtPx
								} else {
									r = cm.Right.Value() / (hMWt + hHWt) * hPxRemain
								}
								hPx += r
								cf.Margin.Right = util.Roundi32(r)
								measuredFlags |= flagMeasuredMarginRight
							case nux.Percent:
								if hasMaxWidth {
									r := hPxMax * cs.Margin().Right.Value()
									hPx += r
									cf.Margin.Right = util.Roundi32(r)
									measuredFlags |= flagMeasuredMarginRight
								}
							}
						} else {
							measuredFlags |= flagMeasuredMarginRight
						}
					}
				} else {
					measuredFlags |= flagMeasuredMarginLeft | flagMeasuredMarginRight
				}
			}

			if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
				measuredCompleteCount++
			}

			if measuredFlags&flagMeasured == flagMeasured {
				measuredCount++
			}

			if height.Mode() == nux.Auto {
				if vPxMax >= innerHeight {
					vPxMax = innerHeight
					height = nux.MeasureSpec(height.Value(), nux.Pixel)
				}
			}

			childrenMeasuredFlags[index] = measuredFlags
		} else {
			childrenMeasuredFlags[index] = flagMeasuredComplete
			measuredCompleteCount++
			measuredCount++
		}
	}

	if hasAutoWeightWidth {
		hasPerWidthWtPx = true
		if hHWt > 0 {
			perWidthWtPx = hHWtPx / hHWt
		} else {
			perWidthWtPx = 0
		}
		hasMaxWidth = true
		hPxMax = (hPx + hMWt*perWidthWtPx) / (1 - hMPt)
	}

	return hPx, hasPerWidthWtPx, perWidthWtPx, hPxMax, hasMaxWidth, vPxMax, measuredCount, measuredCompleteCount
}

// return vPx in innerHeight
func (me *row) measureVertical(width, height nux.MeasureDimen, hPPx, hPPt, innerHeight float32, canMeasureWidth bool, child nux.Widget, measuredFlags uint8) (measuredFlagsOut uint8, vPx float32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd("nuxui", "ui.Row measureHorizontal", measuredDuration)
	var vMWt float32 // margin weight size
	var vMPx float32 // margin pixel size
	var vMPt float32 // margin percent size
	var vHPx float32 // width pixel size
	var vPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cf := cs.Frame()

		// 1. determine the known size of width, add weight and percentage
		if cm := cs.Margin(); cm != nil {
			if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
				if cm.Top.Value() != 0 {
					switch cm.Top.Mode() {
					case nux.Pixel:
						t := cm.Top.Value()
						cf.Margin.Top = util.Roundi32(t)
						vMPx += t
						measuredFlags |= flagMeasuredMarginTop
					case nux.Weight:
						vMWt += cm.Top.Value()
					case nux.Percent:
						switch height.Mode() {
						case nux.Pixel:
							t := cm.Top.Value() / 100 * innerHeight
							cf.Margin.Top = util.Roundi32(t)
							vMPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
							vMPt += cm.Top.Value()
						}
					}
				} else {
					measuredFlags |= flagMeasuredMarginTop
				}
			}

			if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
				if cm.Bottom.Value() != 0 {
					switch cm.Bottom.Mode() {
					case nux.Pixel:
						b := cm.Bottom.Value()
						cf.Margin.Bottom = util.Roundi32(b)
						vMPx += b
						measuredFlags |= flagMeasuredMarginBottom
					case nux.Weight:
						vMWt += cm.Bottom.Value()
					case nux.Percent:
						switch height.Mode() {
						case nux.Pixel:
							b := cm.Bottom.Value() / 100 * innerHeight
							cf.Margin.Bottom = util.Roundi32(b)
							vMPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
							vMPt += cm.Bottom.Value()
						}
					}
				} else {
					measuredFlags |= flagMeasuredMarginBottom
				}
			}
		} else {
			measuredFlags |= flagMeasuredMarginTop | flagMeasuredMarginBottom
		}

		if vMPt < 0 || vMPt > 100 {
			log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
		}

		vPxRemain = innerHeight - vMPx

		if vPxRemain < 0 {
			vPxRemain = 0
		}

		// 2. First determine the known size of width, add weight and percentage
		// do not add width to hPx and hPxUsed until width measured
		if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
				setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				vHPx = h
				measuredFlags |= flagMeasuredHeight
			case nux.Weight:
				switch height.Mode() {
				case nux.Pixel:
					h := cs.Height().Value() / (vMWt + cs.Height().Value()) * vPxRemain
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
					measuredFlags |= flagMeasuredHeight
					vHPx = h
				case nux.Auto, nux.Unlimit:
					h := cs.Height().Value() / (vMWt + cs.Height().Value()) * vPxRemain
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), height.Mode()))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				}
			case nux.Percent:
				switch height.Mode() {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
					measuredFlags |= flagMeasuredHeight
					vHPx = h
				case nux.Auto, nux.Unlimit:
					h := cs.Height().Value() / 100 * innerHeight
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), width.Mode()))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				}
			// case nux.Ratio:// handled when handle height size
			case nux.Auto, nux.Unlimit:
				cf.Height = int32(nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode()))
				setRatioWidthIfNeed(cs, cf, vPxRemain, nux.Pixel)
			}
		}

		if canMeasureWidth {
			if measuredFlags&flagMeasured != flagMeasured {
				if m, ok := child.(nux.Measure); ok {

					m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))

					if nux.MeasureDimen(cf.Width).Mode() != nux.Pixel ||
						nux.MeasureDimen(cf.Height).Mode() != nux.Pixel {
						log.Fatal("nuxui", "row %s the child %s(%T) measured not completed", me.Info().ID, child.Info().ID, child)
					}

					if cs.Width().Mode() == nux.Ratio {
						oldwidth := cf.Width
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Height)*cs.Width().Value()), nux.Pixel))
						if oldwidth != cf.Width {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					if cs.Height().Mode() == nux.Ratio {
						oldheight := cf.Height
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Width)/cs.Height().Value()), nux.Pixel))
						if oldheight != cf.Height {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					measuredFlags |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
				} else {
					measuredFlags |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
					cf.Width = 0
					cf.Height = 0
				}

				vPxRemain -= float32(cf.Height)
				if vPxRemain < 0 {
					vPxRemain = 0
				}

			}
		}

		if measuredFlags&flagMeasuredHeight == flagMeasuredHeight {
			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					if cm.Top.Value() != 0 {
						switch cm.Top.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch height.Mode() {
							case nux.Pixel:
								if vMWt > 0 && vPxRemain > 0 {
									t := cm.Top.Value() / vMWt * vPxRemain
									cf.Margin.Top = util.Roundi32(t)
									vMPx += t
								} else {
									cf.Margin.Top = 0
								}
								measuredFlags |= flagMeasuredMarginTop
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured, width mode will be Pixel
							}
						case nux.Percent:
							switch width.Mode() {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								if vMWt == 0 {
									tp := cm.Top.Value() / 100.0
									t := (vMPx + vHPx) / (1 - tp) * tp
									cf.Margin.Top = util.Roundi32(t)
									vMPx += t
								}
								// else wait until maxWidth measured, width mode will be Pixel
							}
						}
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					if cm.Bottom.Value() != 0 {
						switch cm.Bottom.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch height.Mode() {
							case nux.Pixel:
								if vMWt > 0 && vPxRemain > 0 {
									b := cm.Bottom.Value() / vMWt * vPxRemain
									cf.Margin.Bottom = util.Roundi32(b)
									vMPx += b
								} else {
									cf.Margin.Bottom = 0
								}
								measuredFlags |= flagMeasuredMarginBottom
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured, height mode will be Pixel
							}
						case nux.Percent:
							switch height.Mode() {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								if vMWt == 0 {
									bp := cm.Bottom.Value() / 100.0
									b := (vMPx + vHPx) / (1 - bp) * bp
									cf.Margin.Bottom = util.Roundi32(b)
									vMPx += b
								}
								// else wait until maxWidth measured, height mode will be Pixel
							}
						}
					}
				}
			}

			if cs.Height().Mode() == nux.Percent {
				ph := float32(cf.Height) / (cs.Height().Value() / 100.0)
				if ph > vMPx+float32(cf.Height) {
					return measuredFlags, ph
				}
			}

			return measuredFlags, vMPx + float32(cf.Height)
		}

		return measuredFlags, vMPx + vHPx
	} else {
		measuredFlags = flagMeasuredComplete
	}

	return measuredFlags, 0
}

// Responsible for determining the position of the widget align, margin...
// TODO measure other mode dimen
func (me *row) Layout(x, y, width, height int32) {
	log.W("nuxui", "row:%s layout %d, %d, %d, %d", me.Info().ID, x, y, width, height)
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

	frame := me.Frame()
	clip := me.clipChildren == clipChildrenYes
	if me.clipChildren == clipChildrenAuto {
		clip = me.childrenWidth > float32(frame.Width)
	}
	if clip {
		canvas.Save()
		canvas.ClipRect(float32(frame.X), float32(frame.Y), float32(frame.Width), float32(frame.Height))
	}

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if draw, ok := child.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if clip {
		canvas.Restore()
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
