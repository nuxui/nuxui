// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Layer interface {
	nux.Parent
	nux.Size
	Visual
}

type layer struct {
	nux.WidgetParent
	nux.WidgetBase
	nux.WidgetSize
	WidgetVisual
}

func NewLayer() Layer {
	me := &layer{}
	me.WidgetParent.Owner = me
	me.WidgetSize.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.Owner = me
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

func (me *layer) Creating(attr nux.Attr) {
	if attr == nil {
		attr = nux.Attr{}
	}

	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetParent.Creating(attr)
	me.WidgetVisual.Creating(attr)
}

func (me *layer) onSizeChanged(widget nux.Widget) {

}
func (me *layer) onVisualChanged(widget nux.Widget) {

}

func (me *layer) Measure(width, height int32) {
	log.I("nuxui", "layer:%s Measure width=%s, height=%s", me.ID(), nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Layer %s Measure", me.ID())

	ms := me.MeasuredSize()
	// childrenMeasuredFlags := map[int]byte{}
	childrenMeasuredFlags := make([]byte, len(me.Children()))
	if nux.MeasureSpecMode(width) == nux.Pixel && nux.MeasureSpecMode(height) == nux.Pixel {
		log.I("nuxui", "layer:%s  measure once", me.ID())
		me.measure(width, height, childrenMeasuredFlags)
	} else {
		w, h := me.measure(width, height, childrenMeasuredFlags)
		log.I("nuxui", "layer:%s  measure first w=%s, h=%s", me.ID(), nux.MeasureSpecString(w), nux.MeasureSpecString(h))
		if w != width || h != height {
			w, h = me.measure(w, h, childrenMeasuredFlags)
			log.I("nuxui", "layer:%s  measure second w=%s, h=%s", me.ID(), nux.MeasureSpecString(w), nux.MeasureSpecString(h))
		}
		ms.Width = w
		ms.Height = h
	}
}

func (me *layer) measure(width, height int32, childrenMeasuredFlags []byte) (outWidth, outHeight int32) {
	var hPxMax float32   // max horizontal size
	var hPPt float32     // horizontal padding percent
	var hPPx float32     // horizontal padding pixel
	var hPPxUsed float32 // horizontal padding size include percent size

	var vPxMax float32
	var vPPt float32
	var vPPx float32
	var vPPxUsed float32

	var innerWidth float32
	var innerHeight float32

	ms := me.MeasuredSize()

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
	} else {
		log.V("nuxui", "layer has no padding")
	}

	innerWidth = float32(nux.MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(nux.MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	var hPxUsed float32
	var hPx float32
	var hWt float32
	var hPt float32

	var vPxUsed float32
	var vPx float32
	var vWt float32
	var vPt float32

	var measuredFlags byte

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		hPxUsed = 0
		hPx = 0
		hWt = 0
		hPt = 0

		vPxUsed = 0
		vPx = 0
		vWt = 0
		vPt = 0

		measuredFlags = childrenMeasuredFlags[index]

		log.D("nuxui", "layer measure child:%s, measuredFlags=%.8b", child.ID(), measuredFlags)

		// TODO if child visible == gone , then skip

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			if measuredFlags == 0 {
				cms.Width = nux.MeasureSpec(0, nux.Unlimit)
				cms.Height = nux.MeasureSpec(0, nux.Unlimit)
			}

			if cs.HasMargin() {
				if measuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
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
							log.V("nuxui", "child:%s, ====", child.ID())
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait max height measured
						}
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
							vPt += cs.MarginTop().Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					vPx += float32(cms.Margin.Top)
					vPxUsed += float32(cms.Margin.Top)
				}

				if measuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
					switch cs.MarginBottom().Mode() {
					case nux.Pixel:
						b := cs.MarginBottom().Value()
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						measuredFlags |= hMeasuredMarginBottom
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							vWt += cs.MarginBottom().Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait max height measured
						}
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
							vPt += cs.MarginBottom().Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					vPx += float32(cms.Margin.Bottom)
					vPxUsed += float32(cms.Margin.Bottom)
				}

				if measuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
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
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait until max width measured.
						}
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
							hPt += cs.MarginLeft().Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					hPx += float32(cms.Margin.Left)
					hPxUsed += float32(cms.Margin.Left)
				}

				if measuredFlags&hMeasuredMarginRight == hMeasuredMarginRight {
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
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait until max width measured.
						}
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
							hPt += cs.MarginRight().Value()
							// ok, wait until percent grand total
						}

					}
				} else {
					hPx += float32(cms.Margin.Right)
					hPxUsed += float32(cms.Margin.Right)
				}
			} else {
				measuredFlags |= hMeasuredMarginLeft
				measuredFlags |= hMeasuredMarginRight
				measuredFlags |= hMeasuredMarginTop
				measuredFlags |= hMeasuredMarginBottom
			}

			hPxRemain := innerWidth - hPxUsed
			vPxRemain := innerHeight - vPxUsed
			if hPxRemain < 0 {
				hPxRemain = 0
			}
			if vPxRemain < 0 {
				vPxRemain = 0
			}

			canMeasureWidth := true
			canMeasureHeight := true

			if measuredFlags&hMeasuredWidth != hMeasuredWidth {
				switch cs.Width().Mode() {
				case nux.Pixel:
					w := cs.Width().Value()
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, nux.Pixel)
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						wt := hWt + cs.Width().Value()
						w := cs.Width().Value() / wt * hPxRemain
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cms, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						canMeasureWidth = false
						// ok, wait max width measured
					}
				case nux.Percent:
					log.V("nuxui", "child:%s parent,px width -> child percent width 0", child.ID())
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						log.V("nuxui", "child:%s parent,px width -> child percent width", child.ID())
						w := cs.Width().Value() / 100 * innerWidth
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cms, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.Width().Value()
						if hPx > 0 {
							_innerWidth := hPx / (1.0 - hPt/100.0)
							w := cs.Width().Value() / 100 * _innerWidth
							log.V("nuxui", "child:%s parent,auto width -> child percent width %f", child.ID(), w)
							cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
							setRatioWidth(cs, cms, nux.Pixel)
							// ok
						} else {
							canMeasureWidth = false
							log.V("nuxui", "child:%s parent,auto width -> child percent width 2", child.ID())
							//ok, marginLeft width marginRight has no definite value, wait max height measured
						}
					}
				case nux.Ratio:
					// ok, measured when measure height
					if cs.Height().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both 'Ratio'")
					}
				case nux.Auto, nux.Unlimit:
					cms.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode())
					setRatioHeight(cs, cms, nux.Pixel)
					// ok
				}
			}

			if measuredFlags&hMeasuredHeight != hMeasuredHeight {
				switch cs.Height().Mode() {
				case nux.Pixel:
					h := cs.Height().Value()
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, nux.Pixel)
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						wt := vWt + cs.Height().Value()
						h := cs.Height().Value() / wt * vPxRemain
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cms, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						canMeasureHeight = false
						// ok wait max height measured
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
						if vPx > 0 {
							_innerHeight := vPx / (1.0 - vPt/100.0)
							h := cs.Height().Value() / 100 * _innerHeight
							cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
							setRatioWidth(cs, cms, nux.Pixel)
							// ok
						} else {
							canMeasureHeight = false
							//ok, marginTop height marginBottom has no definite value, wait max height measured
						}
					}
				case nux.Ratio:
					if cs.Width().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
					}
				case nux.Auto, nux.Unlimit:
					cms.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode())
					setRatioWidth(cs, cms, nux.Pixel)
					// ok
				}
			}

			if canMeasureWidth && canMeasureHeight {
				if m, ok := child.(nux.Measure); ok {
					if measuredFlags&hMeasuredWidth != hMeasuredWidth || measuredFlags&hMeasuredHeight != hMeasuredHeight {
						log.I("nuxui", "child:%s call measure", child.ID())
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
					}

					measuredFlags |= hMeasuredWidth
					measuredFlags |= hMeasuredHeight

					hPx += float32(cms.Width)
					vPx += float32(cms.Height)
					hPxRemain -= float32(cms.Width)
					vPxRemain -= float32(cms.Height)
				}
			} // if canMeasureWidth && canMeasureHeight

			if cs.HasMargin() {
				if measuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
					switch cs.MarginTop().Mode() {
					// case nux.Pixel: has measured
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							t := cs.MarginTop().Value() / vWt * vPxRemain
							cms.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= hMeasuredMarginTop
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							if vPx > 0 {
								_innerHeight := vPx / (1.0 - vPt/100.0)
								t := cs.MarginTop().Value() / 100 * _innerHeight
								cms.Margin.Top = util.Roundi32(t)
								vPx += t
								// ok
							} else {
								//ok, marginTop height marginBottom has no definite value, wait max height measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
					switch cs.MarginBottom().Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							b := cs.MarginBottom().Value() / vWt * vPxRemain
							cms.Margin.Bottom = util.Roundi32(b)
							vPx += b
							measuredFlags |= hMeasuredMarginBottom
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							if vPx > 0 {
								_innerHeight := vPx / (1.0 - vPt/100.0)
								b := cs.MarginBottom().Value() / 100 * _innerHeight
								cms.Margin.Bottom = util.Roundi32(b)
								vPx += b
								// ok
							} else {
								//ok, marginTop height marginBottom has no definite value, wait max height measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
					switch cs.MarginLeft().Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							l := cs.MarginLeft().Value() / hWt * hPxRemain
							cms.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= hMeasuredMarginLeft
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							if hPx > 0 {
								_innerWidth := hPx / (1.0 - hPt/100.0)
								l := cs.MarginLeft().Value() / 100 * _innerWidth
								cms.Margin.Left = util.Roundi32(l)
								hPx += l
								// ok
							} else {
								//ok, marginLeft width marginRight has no definite value, wait max width measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginRight == hMeasuredMarginRight {
					switch cs.MarginRight().Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							r := cs.MarginRight().Value() / hWt * hPxRemain
							cms.Margin.Right = util.Roundi32(r)
							hPx += r
							measuredFlags |= hMeasuredMarginRight
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							if hPx > 0 {
								_innerWidth := hPx / (1.0 - hPt/100.0)
								r := cs.MarginRight().Value() / 100 * _innerWidth
								cms.Margin.Right = util.Roundi32(r)
								hPx += r
								// ok
							} else {
								//ok, marginLeft width marginRight has no definite value, wait max width measured
							}
						}
					}
				}
			}

			hPx = hPx / (1.0 - hPt/100.0)
			vPx = vPx / (1.0 - vPt/100.0)

			if hPx > hPxMax {
				hPxMax = hPx
			}

			if vPx > vPxMax {
				vPxMax = vPx
			}

			log.V("nuxui", "layer child:%s measured size= %s", child.ID(), cms)
		} // if cs, ok := child.(nux.Size); ok

		childrenMeasuredFlags[index] = measuredFlags
	} // for range children

	// Use the maximum width found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)

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
		}
	}

	switch nux.MeasureSpecMode(height) {
	case nux.Auto, nux.Unlimit:
		innerHeight = vPxMax
		h := (innerHeight + vPPx) / (1 - vPPt/100.0)
		height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)

		if me.HasPadding() {
			if me.PaddingTop().Mode() == nux.Percent {
				t := me.PaddingTop().Value() / 100 * h
				ms.Padding.Top = util.Roundi32(t)
				hPPx += t
				hPPxUsed += t
			}

			if me.PaddingBottom().Mode() == nux.Percent {
				b := me.PaddingBottom().Value() / 100 * h
				ms.Padding.Bottom = util.Roundi32(b)
				hPPx += b
				hPPxUsed += b
			}
		}
	}

	outWidth = width
	outHeight = height

	return
}

func (me *layer) Layout(dx, dy, left, top, right, bottom int32) {
	ms := me.MeasuredSize()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(bottom - top)
	var innerWidth float32 = float32(right - left)

	innerHeight -= float32(ms.Padding.Top + ms.Padding.Bottom)
	innerWidth -= float32(ms.Padding.Left + ms.Padding.Right)
	// t += float32(ms.Padding.Top)
	// l += float32(ms.Padding.Left)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		l = float32(ms.Padding.Left)
		t = float32(ms.Padding.Top)
		// TODO if child visible == gone , then skip
		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

			l += float32(cms.Margin.Left)
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
		}
	}
}

func (me *layer) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

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
