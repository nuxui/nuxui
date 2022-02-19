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
	*nux.WidgetParent
	*nux.WidgetSize
	*WidgetVisual
}

func NewLayer(attrs ...nux.Attr) Layer {
	me := &layer{}
	me.WidgetParent = nux.NewWidgetParent(me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *layer) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *layer) Measure(width, height int32) {
	log.I("nuxui", "layer:%s Measure width=%s, height=%s", me.Info().ID, nux.MeasureSpecString(width), nux.MeasureSpecString(height))
	measureDuration := log.Time()
	log.I("nuxui", "layer:%s  Measure === 0 ", me.Info().ID)
	defer log.TimeEnd(measureDuration, "nuxui", "ui.Layer %s Measure", me.Info().ID)
	log.I("nuxui", "layer:%s  Measure === 1 ", me.Info().ID)

	originWidth := width
	originHeight := height

	frame := me.Frame()
	childrenMeasuredFlags := make([]byte, len(me.Children()))
	if nux.MeasureSpecMode(width) == nux.Pixel && nux.MeasureSpecMode(height) == nux.Pixel {
		log.I("nuxui", "layer:%s  measure once", me.Info().ID)
		me.measure(width, height, childrenMeasuredFlags)
	} else {
		w, h := me.measure(width, height, childrenMeasuredFlags)
		log.I("nuxui", "layer:%s  measure first w=%s, h=%s", me.Info().ID, nux.MeasureSpecString(w), nux.MeasureSpecString(h))
		if w != width || h != height {
			w, h = me.measure(w, h, childrenMeasuredFlags)
			log.I("nuxui", "layer:%s  measure second w=%s, h=%s", me.Info().ID, nux.MeasureSpecString(w), nux.MeasureSpecString(h))
		}
		width = w
		height = h
	}

	log.I("nuxui", "layer:%s  Measure === 2 ", me.Info().ID)
	setNewWidth(frame, originWidth, width)
	setNewHeight(frame, originHeight, height)
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

	frame := me.Frame()

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
			log.D("nuxui", "%T  nux.Component", child)
			child = compt.Content()
			log.D("nuxui", "%T  nux.Component 2", child)
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

		log.D("nuxui", "layer measure child:%s, measuredFlags=%.8b", child.Info().ID, measuredFlags)

		// TODO if child visible == gone , then skip

		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			if measuredFlags == 0 {
				cf.Width = nux.MeasureSpec(0, nux.Unlimit)
				cf.Height = nux.MeasureSpec(0, nux.Unlimit)
			}

			if cs.Margin() != nil {
				if measuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
					switch cs.Margin().Top.Mode() {
					case nux.Pixel:
						t := cs.Margin().Top.Value()
						cf.Margin.Top = util.Roundi32(t)
						vPx += t
						vPxUsed += t
						measuredFlags |= hMeasuredMarginTop
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							vWt += cs.Margin().Top.Value()
							log.V("nuxui", "child:%s, ====", child.Info().ID)
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait max height measured
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							t := cs.Margin().Top.Value() / 100 * innerHeight
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							vPxUsed += t
							measuredFlags |= hMeasuredMarginTop
							// ok
						case nux.Auto, nux.Unlimit:
							vPt += cs.Margin().Top.Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					vPx += float32(cf.Margin.Top)
					vPxUsed += float32(cf.Margin.Top)
				}

				if measuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
					switch cs.Margin().Bottom.Mode() {
					case nux.Pixel:
						b := cs.Margin().Bottom.Value()
						cf.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						measuredFlags |= hMeasuredMarginBottom
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							vWt += cs.Margin().Bottom.Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait max height measured
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							b := cs.Margin().Bottom.Value() / 100 * innerHeight
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							vPxUsed += b
							measuredFlags |= hMeasuredMarginBottom
							// ok
						case nux.Auto, nux.Unlimit:
							vPt += cs.Margin().Bottom.Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					vPx += float32(cf.Margin.Bottom)
					vPxUsed += float32(cf.Margin.Bottom)
				}

				if measuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
					switch cs.Margin().Left.Mode() {
					case nux.Pixel:
						l := cs.Margin().Left.Value()
						cf.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						measuredFlags |= hMeasuredMarginLeft
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							hWt += cs.Margin().Left.Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait until max width measured.
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							l := cs.Margin().Left.Value() / 100 * innerWidth
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							hPxUsed += l
							measuredFlags |= hMeasuredMarginLeft
							// ok
						case nux.Auto, nux.Unlimit:
							hPt += cs.Margin().Left.Value()
							// ok, wait until percent grand total
						}
					}
				} else {
					hPx += float32(cf.Margin.Left)
					hPxUsed += float32(cf.Margin.Left)
				}

				if measuredFlags&hMeasuredMarginRight == hMeasuredMarginRight {
					switch cs.Margin().Right.Mode() {
					case nux.Pixel:
						r := cs.Margin().Right.Value()
						cf.Margin.Right = int32(r)
						hPx += r
						hPxUsed += r
						measuredFlags |= hMeasuredMarginRight
						// ok
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							hWt += cs.Margin().Right.Value()
							// ok, wait until weight grand total
						case nux.Auto, nux.Unlimit:
							// ok, wait until max width measured.
						}
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							r := cs.Margin().Right.Value() / 100 * innerWidth
							cf.Margin.Right = util.Roundi32(r)
							hPx += r
							hPxUsed += r
							measuredFlags |= hMeasuredMarginRight
							// ok
						case nux.Auto, nux.Unlimit:
							hPt += cs.Margin().Right.Value()
							// ok, wait until percent grand total
						}

					}
				} else {
					hPx += float32(cf.Margin.Right)
					hPxUsed += float32(cf.Margin.Right)
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
					cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cf, w, nux.Pixel)
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						wt := hWt + cs.Width().Value()
						w := cs.Width().Value() / wt * hPxRemain
						cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cf, w, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						canMeasureWidth = false
						// ok, wait max width measured
					}
				case nux.Percent:
					log.V("nuxui", "child:%s parent,px width -> child percent width 0", child.Info().ID)
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						log.V("nuxui", "child:%s parent,px width -> child percent width", child.Info().ID)
						w := cs.Width().Value() / 100 * innerWidth
						cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cf, w, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.Width().Value()
						if hPx > 0 {
							_innerWidth := hPx / (1.0 - hPt/100.0)
							w := cs.Width().Value() / 100 * _innerWidth
							log.V("nuxui", "child:%s parent,auto width -> child percent width %f", child.Info().ID, w)
							cf.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
							setRatioHeight(cs, cf, w, nux.Pixel)
							// ok
						} else {
							canMeasureWidth = false
							log.V("nuxui", "child:%s parent,auto width -> child percent width 2", child.Info().ID)
							//ok, marginLeft width marginRight has no definite value, wait max height measured
						}
					}
				case nux.Ratio:
					// ok, measured when measure height
					if cs.Height().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both 'Ratio'")
					}
				case nux.Auto, nux.Unlimit:
					cf.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode())
					setRatioHeight(cs, cf, hPxRemain, nux.Pixel)
					// ok
				}
			}

			if measuredFlags&hMeasuredHeight != hMeasuredHeight {
				switch cs.Height().Mode() {
				case nux.Pixel:
					h := cs.Height().Value()
					cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cf, h, nux.Pixel)
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						wt := vWt + cs.Height().Value()
						h := cs.Height().Value() / wt * vPxRemain
						cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cf, h, nux.Pixel)
						// ok
					case nux.Auto, nux.Unlimit:
						canMeasureHeight = false
						// ok wait max height measured
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
						if vPx > 0 {
							_innerHeight := vPx / (1.0 - vPt/100.0)
							h := cs.Height().Value() / 100 * _innerHeight
							cf.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
							setRatioWidth(cs, cf, h, nux.Pixel)
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
					cf.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode())
					setRatioWidth(cs, cf, vPxRemain, nux.Pixel)
					// ok
				}
			}

			if canMeasureWidth && canMeasureHeight {
				if m, ok := child.(nux.Measure); ok {
					if measuredFlags&hMeasuredWidth != hMeasuredWidth || measuredFlags&hMeasuredHeight != hMeasuredHeight {
						log.I("nuxui", "child:%s call measure", child.Info().ID)
						m.Measure(cf.Width, cf.Height)

						if cs.Width().Mode() == nux.Ratio {
							oldWidth := cf.Width
							cf.Width = nux.MeasureSpec(util.Roundi32(float32(cf.Height)*cs.Width().Value()), nux.Pixel)
							if oldWidth != cf.Width {
								m.Measure(cf.Width, cf.Height)
							}
						}

						if cs.Height().Mode() == nux.Ratio {
							oldHeight := cf.Height
							cf.Height = nux.MeasureSpec(util.Roundi32(float32(cf.Width)/cs.Height().Value()), nux.Pixel)
							if oldHeight != cf.Height {
								m.Measure(cf.Width, cf.Height)
							}
						}
					}

					measuredFlags |= hMeasuredWidth
					measuredFlags |= hMeasuredHeight

					hPx += float32(cf.Width)
					vPx += float32(cf.Height)
					hPxRemain -= float32(cf.Width)
					vPxRemain -= float32(cf.Height)
				}
			} // if canMeasureWidth && canMeasureHeight

			if cs.Margin() != nil {
				if measuredFlags&hMeasuredMarginTop != hMeasuredMarginTop {
					switch cs.Margin().Top.Mode() {
					// case nux.Pixel: has measured
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							t := cs.Margin().Top.Value() / vWt * vPxRemain
							cf.Margin.Top = util.Roundi32(t)
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
								t := cs.Margin().Top.Value() / 100 * _innerHeight
								cf.Margin.Top = util.Roundi32(t)
								vPx += t
								// ok
							} else {
								//ok, marginTop height marginBottom has no definite value, wait max height measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginBottom != hMeasuredMarginBottom {
					switch cs.Margin().Bottom.Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(height) {
						case nux.Pixel:
							b := cs.Margin().Bottom.Value() / vWt * vPxRemain
							cf.Margin.Bottom = util.Roundi32(b)
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
								b := cs.Margin().Bottom.Value() / 100 * _innerHeight
								cf.Margin.Bottom = util.Roundi32(b)
								vPx += b
								// ok
							} else {
								//ok, marginTop height marginBottom has no definite value, wait max height measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginLeft != hMeasuredMarginLeft {
					switch cs.Margin().Left.Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							l := cs.Margin().Left.Value() / hWt * hPxRemain
							cf.Margin.Left = util.Roundi32(l)
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
								l := cs.Margin().Left.Value() / 100 * _innerWidth
								cf.Margin.Left = util.Roundi32(l)
								hPx += l
								// ok
							} else {
								//ok, marginLeft width marginRight has no definite value, wait max width measured
							}
						}
					}
				}

				if measuredFlags&hMeasuredMarginRight == hMeasuredMarginRight {
					switch cs.Margin().Right.Mode() {
					case nux.Weight:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							r := cs.Margin().Right.Value() / hWt * hPxRemain
							cf.Margin.Right = util.Roundi32(r)
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
								r := cs.Margin().Right.Value() / 100 * _innerWidth
								cf.Margin.Right = util.Roundi32(r)
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

			log.V("nuxui", "layer child:%s measured size= %s", child.Info().ID, cf)
		} // if cs, ok := child.(nux.Size); ok

		childrenMeasuredFlags[index] = measuredFlags
	} // for range children

	// Use the maximum width found in the first traversal as the width in auto mode, and calculate the percent size
	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPxMax
		w := (innerWidth + hPPx) / (1 - hPPt/100.0)
		width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)

		if me.Padding() != nil {
			if me.Padding().Left.Mode() == nux.Percent {
				l := me.Padding().Left.Value() / 100 * w
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
			}

			if me.Padding().Right.Mode() == nux.Percent {
				r := me.Padding().Right.Value() / 100 * w
				frame.Padding.Right = util.Roundi32(r)
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

		if me.Padding() != nil {
			if me.Padding().Top.Mode() == nux.Percent {
				t := me.Padding().Top.Value() / 100 * h
				frame.Padding.Top = util.Roundi32(t)
				hPPx += t
				hPPxUsed += t
			}

			if me.Padding().Bottom.Mode() == nux.Percent {
				b := me.Padding().Bottom.Value() / 100 * h
				frame.Padding.Bottom = util.Roundi32(b)
				hPPx += b
				hPPxUsed += b
			}
		}
	}

	outWidth = width
	outHeight = height

	return
}

func (me *layer) Layout(x, y, width, height int32) {
	frame := me.Frame()

	var l float32 = 0
	var t float32 = 0

	var innerHeight float32 = float32(width)
	var innerWidth float32 = float32(height)

	innerHeight -= float32(frame.Padding.Top + frame.Padding.Bottom)
	innerWidth -= float32(frame.Padding.Left + frame.Padding.Right)
	// t += float32(frame.Padding.Top)
	// l += float32(frame.Padding.Left)

	for _, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		l = float32(frame.Padding.Left)
		t = float32(frame.Padding.Top)
		// TODO if child visible == gone , then skip
		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			l += float32(cf.Margin.Left)
			t += float32(cf.Margin.Top)

			cf.X = x + int32(l)
			cf.Y = y + int32(t)

			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cf.X, cf.Y, cf.Width, cf.Height)
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

		if draw, ok := child.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}

}
