// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	"nuxui.org/nuxui/util"
)

type Column interface {
	nux.Parent
	nux.Size
	Visual
}

type column struct {
	*nux.WidgetParent
	*nux.WidgetSize
	*WidgetVisual

	align          *Align // TODO horizontal|center
	childrenWidth  float32
	childrenHeight float32
	clipChildren   int
}

func NewColumn(attr nux.Attr) Column {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &column{
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

func (me *column) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *column) Measure(width, height nux.MeasureDimen) {
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

	vPx, vMPt, vMWt, vHWt, vPxMax, hasMaxHeight, hasAutoWeightHeight, hPxMax, measuredWidthCount, measuredCompleteCount := me.measureChildren1(width, height, hPPx, hPPt, innerWidth, innerHeight, childrenMeasuredFlags)

	if vMPt > 100.0 {
		log.Fatal("nuxui", "the sum of percent size should at 0% ~ 100%")
	}

	switch width.Mode() {
	case nux.Pixel:
		frame.Width = width.Value()
		me.childrenWidth = hPxMax
	case nux.Auto, nux.Unlimit:
		if measuredWidthCount == me.ChildrenCount() {
			if width.Mode() == nux.Auto {
				iw := float32(width.Value())*(1-hPPt/100.0) - hPPx
				if hPxMax > iw {
					hPxMax = iw
				}
			}
			innerWidth = hPxMax
			w := (innerWidth + hPPx) / (1 - hPPt/100.0)
			width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel) //TODO:: no padding
			frame.Width = util.Roundi32(w)                       //TODO:: no padding
			me.childrenWidth = hPxMax
		}
	}

	switch height.Mode() {
	case nux.Pixel:
		frame.Height = height.Value()
		me.childrenHeight = vPx
	case nux.Auto, nux.Unlimit:
		if hasMaxHeight || measuredCompleteCount == me.ChildrenCount() {
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

	// if measure finished at first traversal
	if measuredCompleteCount == me.ChildrenCount() {
		if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
			measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
		}
		return
	}

	vPxRemain := innerHeight - vPx
	if vPxRemain < 0 {
		vPxRemain = 0
	}

	vPx, hasPerHeightWtPx, perHeightWtPx, vPxMax, hasMaxHeight, hPxMax, measuredCount, measuredCompleteCount := me.measureChildren2(width, height, hPPx, hPPt, innerWidth, innerHeight, vPx, vPxRemain, vMWt, vHWt, vMPt, vPxMax, hasMaxHeight, hasAutoWeightHeight, false, 0, hPxMax, childrenMeasuredFlags)

	if measuredCount != me.ChildrenCount() {
		log.Fatal("nuxui", "can not run here")
	}

	switch width.Mode() {
	// case nux.Pixel: handled
	case nux.Auto, nux.Unlimit:
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
		frame.Width = util.Roundi32(w)
		me.childrenWidth = hPxMax // TODO:: exclude margin
	}

	switch height.Mode() {
	// case nux.Pixel: handled
	case nux.Auto, nux.Unlimit:
		if hasMaxHeight {
			innerHeight = vPxMax
			h := (innerHeight + vPPx) / (1 - vPPt/100.0)
			height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
			frame.Height = util.Roundi32(h)
			me.childrenHeight = vPxMax
		} else {
			h := (vPx + vPPx) / (1 - vPPt/100.0)
			frame.Height = util.Roundi32(h)
			me.childrenHeight = vPx
		}
	}

	if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
	}

	if measuredCompleteCount == me.ChildrenCount() {
		return
	}

	// third measure
	vPx, _, _, _, _, hPxMax, _, measuredCompleteCount = me.measureChildren2(width, height, hPPx, hPPt, innerWidth, innerHeight, vPx, vPxRemain, vMWt, vHWt, vMPt, vPxMax, false, false, hasPerHeightWtPx, perHeightWtPx, hPxMax, childrenMeasuredFlags)

	if measuredCompleteCount == me.ChildrenCount() {
		return
	}

	log.Fatal("nuxui", "column %s not measure complete.", me.Info().ID)
}

/*
vPx: sum of vertical pixel size
vMWt: sum of children vertical margin weight
vWt: sum of children vertical height weight
vPt: sum of vertical percent size
*/
func (me *column) measureChildren1(width, height nux.MeasureDimen, hPPx, hPPt, innerWidth, innerHeight float32, childrenMeasuredFlags []uint8) (vPx, vMPt, vMWt, vHWt, vPxMax float32, hasMaxHeight, hasAutoWeightHeight bool, hPxMax float32, measuredWidthCount, measuredCompleteCount int) {
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
				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					if cm.Top.Value() != 0 {
						switch cm.Top.Mode() {
						case nux.Pixel:
							t := cm.Top.Value()
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Weight:
							vMWt += cm.Top.Value()
						case nux.Percent:
							switch height.Mode() {
							case nux.Pixel:
								t := cm.Top.Value() / 100 * innerHeight
								cf.Margin.Top = util.Roundi32(t)
								vPx += t
								measuredFlags |= flagMeasuredMarginTop
							case nux.Auto, nux.Unlimit:
								// wait until percent grand total
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
							vPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Weight:
							vMWt += cm.Bottom.Value()
						case nux.Percent:
							switch height.Mode() {
							case nux.Pixel:
								b := cm.Bottom.Value() / 100 * innerHeight
								cf.Margin.Bottom = util.Roundi32(b)
								vPx += b
								measuredFlags |= flagMeasuredMarginBottom
							case nux.Auto, nux.Unlimit:
								// wait until percent grand total
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

			canMeasureHeight := true
			// forceMeasure := false

			// Do not add height to vPx and vPxUsed until height measured
			if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
				if cs.Height().Value() > 0 || cs.Height().Mode() == nux.Auto || cs.Height().Mode() == nux.Unlimit {
					switch cs.Height().Mode() {
					case nux.Pixel:
						h := cs.Height().Value()
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
						setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
						measuredFlags |= flagMeasuredHeight
					case nux.Weight:
						// wait until weight grand total and use height.Mode()
						vHWt += cs.Height().Value()
						canMeasureHeight = false
						switch height.Mode() {
						case nux.Auto, nux.Unlimit:
							hasAutoWeightHeight = true
						}
					case nux.Percent:
						switch height.Mode() {
						case nux.Pixel:
							h := cs.Height().Value() / 100 * innerHeight
							cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
							setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
							measuredFlags |= flagMeasuredHeight
						case nux.Auto, nux.Unlimit:
							h := cs.Height().Value() / 100 * innerHeight
							cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), height.Mode()))
							setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
						}
					case nux.Ratio:
						if cs.Width().Mode() == nux.Ratio {
							log.Fatal("nuxui", "width and height size mode can not both Ratio")
						}
					case nux.Auto, nux.Unlimit:
						h := innerHeight - vPx
						if h < 0 {
							h = 0
						}
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), cs.Height().Mode()))
						setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
					}
				} else {
					cf.Height = 0
					setRatioWidthIfNeed(cs, cf, 0, nux.Pixel)
					measuredFlags |= flagMeasuredHeight
				}
			}

			measuredFlags, hPx := me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, canMeasureHeight, child, measuredFlags)

			if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
				measuredCompleteCount++
			}

			if measuredFlags&flagMeasuredWidth == flagMeasuredWidth {
				measuredWidthCount++
			}

			if nux.MeasureDimen(cf.Height).Mode() == nux.Pixel {
				vPx += float32(cf.Height)

				// find max innerHeight in this mode
				if height.Mode() != nux.Pixel && cs.Height().Mode() == nux.Percent {
					hasMaxHeight = true
					if cs.Height().Value() > 0 {
						v := float32(cf.Height) / (cs.Height().Value() / 100.0)
						if v > vPxMax {
							vPxMax = v
						}
					}
				}
			}

			if hPx > hPxMax {
				hPxMax = hPx
			}
			if width.Mode() == nux.Auto {
				if hPxMax >= innerWidth {
					hPxMax = innerWidth
					width = nux.MeasureSpec(width.Value(), nux.Pixel)
				}
			}

			childrenMeasuredFlags[index] = measuredFlags
		} else {
			childrenMeasuredFlags[index] = flagMeasuredComplete
			measuredWidthCount++
			measuredCompleteCount++
		}
	}

	v := vPx / (1 - vMPt/100.0)
	if v > vPxMax {
		vPxMax = v
	}
	return
}

func (me *column) measureChildren2(width, height nux.MeasureDimen, hPPx, hPPt, innerWidth, innerHeight float32,
	vPx, vPxRemain, vMWt, vHWt, vMPt, vPxMax float32, hasMaxHeight, hasAutoWeightHeight, hasPerHeightWtPx bool, perHeightWtPx, hPxMax float32, childrenMeasuredFlags []uint8) (vPxOut float32, hasPerHeightWtPxOut bool, perHeightWtPxOut float32, vPxMaxOut float32, hasMaxHeightOut bool, hPxOut float32, measuredCount, measuredCompleteCount int) {
	var measuredFlags uint8
	var vHWtPx, hPx float32 = 0, 0

	if hasMaxHeight && vMPt > 0 {
		vPxRemain -= vPxMax * vMPt
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
				if cs.Height().Value() > 0 {
					switch cs.Height().Mode() {
					// case nux.Pixel: already handled
					// case nux.Percent: already handled
					// case nux.Ratio: already handled
					// case nux.Auto, nux.Unlimit: already handled
					case nux.Weight:
						h := vPxRemain * cs.Height().Value() / (vMWt + vHWt)
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), height.Mode()))
						setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
					}
				}
			}

			flag := flagMeasured | flagMeasuredHeight | flagMeasuredHorizontalComplete
			if measuredFlags&flag != flag {
				measuredFlags, hPx = me.measureHorizontal(width, height, hPPx, hPPt, innerWidth, true, child, measuredFlags)

				if hPx > hPxMax {
					hPxMax = hPx
				}

				if nux.MeasureDimen(cf.Height).Mode() == nux.Pixel {
					vPx += float32(cf.Height)

					if cs.Height().Mode() == nux.Weight {
						vHWtPx += float32(cf.Height)
					}
				}
			}

			// wait third measure after auto weight height mode measured max height
			if !hasAutoWeightHeight {
				if cm := cs.Margin(); cm != nil {
					if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
						if cm.Top.Value() != 0 {
							switch cm.Top.Mode() {
							// case nux.Pixel: already handled
							case nux.Weight:
								var t float32
								if hasPerHeightWtPx {
									t = cm.Top.Value() * perHeightWtPx
								} else {
									t = cm.Top.Value() / (vMWt + vHWt) * vPxRemain
								}
								vPx += t
								cf.Margin.Top = util.Roundi32(t)
								measuredFlags |= flagMeasuredMarginTop
							case nux.Percent:
								if hasMaxHeight {
									t := vPxMax * cs.Margin().Top.Value()
									vPx += t
									cf.Margin.Top = util.Roundi32(t)
									measuredFlags |= flagMeasuredMarginTop
								}
							}
						} else {
							measuredFlags |= flagMeasuredMarginTop
						}
					}

					if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
						if cm.Bottom.Value() != 0 {
							switch cm.Bottom.Mode() {
							// case nux.Pixel: // already measured
							case nux.Weight:
								var b float32
								if hasPerHeightWtPx {
									b = cm.Bottom.Value() * perHeightWtPx
								} else {
									b = cm.Bottom.Value() / (vMWt + vHWt) * vPxRemain
								}
								vPx += b
								cf.Margin.Bottom = util.Roundi32(b)
								measuredFlags |= flagMeasuredMarginBottom
							case nux.Percent:
								if hasMaxHeight {
									b := vPxMax * cs.Margin().Bottom.Value()
									vPx += b
									cf.Margin.Bottom = util.Roundi32(b)
									measuredFlags |= flagMeasuredMarginBottom
								}
							}
						} else {
							measuredFlags |= flagMeasuredMarginBottom
						}
					}
				} else {
					measuredFlags |= flagMeasuredMarginTop | flagMeasuredMarginBottom
				}
			}

			if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
				measuredCompleteCount++
			}

			if measuredFlags&flagMeasured == flagMeasured {
				measuredCount++
			}

			if width.Mode() == nux.Auto {
				if hPxMax >= innerWidth {
					hPxMax = innerWidth
					width = nux.MeasureSpec(width.Value(), nux.Pixel)
				}
			}

			childrenMeasuredFlags[index] = measuredFlags
		} else {
			childrenMeasuredFlags[index] = flagMeasuredComplete
			measuredCompleteCount++
			measuredCount++
		}
	}

	if hasAutoWeightHeight {
		hasPerHeightWtPx = true
		if vHWt > 0 {
			perHeightWtPx = vHWtPx / vHWt
		} else {
			perHeightWtPx = 0
		}
		hasMaxHeight = true
		vPxMax = (vPx + vMWt*perHeightWtPx) / (1 - vMPt)
	}

	return vPx, hasPerHeightWtPx, perHeightWtPx, vPxMax, hasMaxHeight, hPxMax, measuredCount, measuredCompleteCount
}

// return hPx in innerWidth
func (me *column) measureHorizontal(width, height nux.MeasureDimen, hPPx, hPPt, innerWidth float32, canMeasureHeight bool, child nux.Widget, measuredFlags uint8) (measuredFlagsOut uint8, hPx float32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd("nuxui", "ui.Column measureHorizontal", measuredDuration)
	var hMWt float32 // margin weight size
	var hMPx float32 // margin pixel size
	var hMPt float32 // margin percent size
	var hHPx float32 // width pixel size
	var hPxRemain float32

	if cs, ok := child.(nux.Size); ok {
		cf := cs.Frame()

		// 1. determine the known size of width, add weight and percentage
		if cm := cs.Margin(); cm != nil {
			if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
				if cm.Left.Value() != 0 {
					switch cm.Left.Mode() {
					case nux.Pixel:
						l := cm.Left.Value()
						cf.Margin.Left = util.Roundi32(l)
						hMPx += l
						measuredFlags |= flagMeasuredMarginLeft
					case nux.Weight:
						hMWt += cm.Left.Value()
					case nux.Percent:
						switch width.Mode() {
						case nux.Pixel:
							l := cm.Left.Value() / 100 * innerWidth
							cf.Margin.Left = util.Roundi32(l)
							hMPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
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
						hMPx += r
						measuredFlags |= flagMeasuredMarginRight
					case nux.Weight:
						hMWt += cm.Right.Value()
					case nux.Percent:
						switch width.Mode() {
						case nux.Pixel:
							r := cm.Right.Value() / 100 * innerWidth
							cf.Margin.Right = util.Roundi32(r)
							hMPx += r
							measuredFlags |= flagMeasuredMarginRight
						case nux.Auto, nux.Unlimit:
							// wait until maxWidth measured.
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

		if hMPt < 0 || hMPt > 100 {
			log.Fatal("nuxui", "percent size out of range, it should between 0% ~ 100%.")
		}

		hPxRemain = innerWidth - hMPx

		if hPxRemain < 0 {
			hPxRemain = 0
		}

		// 2. First determine the known size of width, add weight and percentage
		// do not add width to hPx and hPxUsed until width measured
		if measuredFlags&flagMeasuredWidth != flagMeasuredWidth {
			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
				setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				hHPx = w
				measuredFlags |= flagMeasuredWidth
			case nux.Weight:
				switch width.Mode() {
				case nux.Pixel:
					w := cs.Width().Value() / (hMWt + cs.Width().Value()) * hPxRemain
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
					measuredFlags |= flagMeasuredWidth
					hHPx = w
				case nux.Auto, nux.Unlimit:
					w := cs.Width().Value() / (hMWt + cs.Width().Value()) * hPxRemain
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				}
			case nux.Percent:
				switch width.Mode() {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
					measuredFlags |= flagMeasuredWidth
					hHPx = w
				case nux.Auto, nux.Unlimit:
					w := cs.Width().Value() / 100 * innerWidth
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				}
			// case nux.Ratio:// handled when handle height size
			case nux.Auto, nux.Unlimit:
				cf.Width = int32(nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode()))
				setRatioHeightIfNeed(cs, cf, hPxRemain, nux.Pixel)
			}
		}

		if canMeasureHeight {
			if measuredFlags&flagMeasured != flagMeasured {
				if m, ok := child.(nux.Measure); ok {

					m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))

					if nux.MeasureDimen(cf.Width).Mode() != nux.Pixel ||
						nux.MeasureDimen(cf.Height).Mode() != nux.Pixel {
						log.Fatal("nuxui", "column %s the child %s(%T) measured not completed", me.Info().ID, child.Info().ID, child)
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

				hPxRemain -= float32(cf.Width)
				if hPxRemain < 0 {
					hPxRemain = 0
				}

			}
		}

		if measuredFlags&flagMeasuredWidth == flagMeasuredWidth {
			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					if cm.Left.Value() != 0 {
						switch cm.Left.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch width.Mode() {
							case nux.Pixel:
								if hMWt > 0 && hPxRemain > 0 {
									l := cm.Left.Value() / hMWt * hPxRemain
									cf.Margin.Left = util.Roundi32(l)
									hMPx += l
								} else {
									cf.Margin.Left = 0
								}
								measuredFlags |= flagMeasuredMarginLeft
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured, width mode will be Pixel
							}
						case nux.Percent:
							switch width.Mode() {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								if hMWt == 0 {
									lp := cm.Left.Value() / 100.0
									l := (hMPx + hHPx) / (1 - lp) * lp
									cf.Margin.Left = util.Roundi32(l)
									hMPx += l
								}
								// else wait until maxWidth measured, width mode will be Pixel
							}
						}
					}
				}

				if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
					if cm.Right.Value() != 0 {
						switch cm.Right.Mode() {
						// case nux.Pixel: has measured
						case nux.Weight:
							switch width.Mode() {
							case nux.Pixel:
								if hMWt > 0 && hPxRemain > 0 {
									r := cm.Right.Value() / hMWt * hPxRemain
									cf.Margin.Right = util.Roundi32(r)
									hMPx += r
								} else {
									cf.Margin.Right = 0
								}
								measuredFlags |= flagMeasuredMarginRight
							case nux.Auto, nux.Unlimit:
								// wait until maxWidth measured, width mode will be Pixel
							}
						case nux.Percent:
							switch width.Mode() {
							// case nux.Pixel: has measured
							case nux.Auto, nux.Unlimit:
								if hMWt == 0 {
									rp := cm.Right.Value() / 100.0
									r := (hMPx + hHPx) / (1 - rp) * rp
									cf.Margin.Right = util.Roundi32(r)
									hMPx += r
								}
								// else wait until maxWidth measured, width mode will be Pixel
							}
						}
					}
				}
			}

			if cs.Width().Mode() == nux.Percent {
				pw := float32(cf.Width) / (cs.Width().Value() / 100.0)
				if pw > hMPx+float32(cf.Width) {
					return measuredFlags, pw
				}
			}

			return measuredFlags, hMPx + float32(cf.Width)
		}

		return measuredFlags, hMPx + hHPx
	} else {
		measuredFlags = flagMeasuredComplete
	}

	return measuredFlags, 0
}

// Responsible for determining the position of the widget align, margin
// TODO measure other mode dimen
func (me *column) Layout(x, y, width, height int32) {
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	var l float32 = 0
	var t float32 = 0

	var innerWidth float32 = float32(width)
	var innerHeight float32 = float32(height)

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

	frame := me.Frame()
	clip := me.clipChildren == clipChildrenYes
	if me.clipChildren == clipChildrenAuto {
		clip = me.childrenHeight > float32(frame.Height)
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
