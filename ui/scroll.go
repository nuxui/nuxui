// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	"nuxui.org/nuxui/util"
)

// only one child
type Scroll interface {
	nux.Parent
	nux.Size
	Visual
}

type scroll struct {
	*nux.WidgetParent
	*nux.WidgetSize
	*WidgetVisual
	nux.WidgetBase
}

func NewScroll(attr nux.Attr) Scroll {
	me := &scroll{}
	me.WidgetParent = nux.NewWidgetParent(me, attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *scroll) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *scroll) Measure(width, height nux.MeasureDimen) {
	if len(me.Children()) > 1 {
		log.Fatal("nuxui", "ui.Scroll can only contain 1 child, now has %d children", len(me.Children()))
	}

	var hPPt float32     // horizontal padding percent
	var hPPx float32     // horizontal padding pixel
	var hPPxUsed float32 // horizontal padding size include percent size

	var vPPt float32
	var vPPx float32
	var vPPxUsed float32

	var innerWidth float32
	var innerHeight float32

	var hPxUsed float32
	var hPx float32
	var hWt float32
	var hPt float32

	var vPxUsed float32
	var vPx float32
	var vWt float32
	var vPt float32

	originWidth := width
	originHeight := height

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
			switch width.Mode() {
			case nux.Pixel:
				l := me.Padding().Left.Value() / 100.0 * float32(width.Value())
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.Padding().Left.Value()
				// ok, wait until width measured
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
			switch width.Mode() {
			case nux.Pixel:
				r := me.Padding().Right.Value() / 100.0 * float32(width.Value())
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
				// ok
			case nux.Auto, nux.Unlimit:
				hPPt += me.Padding().Right.Value()
				// ok, wait until width measured
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
			switch height.Mode() {
			case nux.Pixel:
				t := me.Padding().Top.Value() / 100.0 * float32(height.Value())
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
			switch height.Mode() {
			case nux.Pixel:
				b := me.Padding().Bottom.Value() / 100.0 * float32(height.Value())
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

	innerWidth = float32(width.Value())*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(height.Value())*(1.0-vPPt/100.0) - vPPx

	if len(me.Children()) == 1 {
		child := me.Children()[0]
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.Frame()

			if cs.Margin() != nil {
				switch cs.Margin().Top.Mode() {
				case nux.Pixel:
					t := cs.Margin().Top.Value()
					cms.Margin.Top = util.Roundi32(t)
					vPx += t
					vPxUsed += t
					// ok
				case nux.Weight:
					switch height.Mode() {
					case nux.Pixel:
						vWt += cs.Margin().Top.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						// ok, cms.Margin.Top = 0, do nothing
					}
				case nux.Percent:
					switch height.Mode() {
					case nux.Pixel:
						t := cs.Margin().Top.Value() / 100 * innerHeight
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
						vPxUsed += t
						// ok
					case nux.Auto, nux.Unlimit:
						vPt += cs.Margin().Top.Value()
						// ok, wait until percent grand total
					}
				}

				switch cs.Margin().Bottom.Mode() {
				case nux.Pixel:
					b := cs.Margin().Bottom.Value()
					cms.Margin.Bottom = util.Roundi32(b)
					vPx += b
					vPxUsed += b
					// ok
				case nux.Weight:
					switch height.Mode() {
					case nux.Pixel:
						vWt += cs.Margin().Bottom.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						// ok, cms.Margin.Bottom = 0, do nothing
					}
				case nux.Percent:
					switch height.Mode() {
					case nux.Pixel:
						b := cs.Margin().Bottom.Value() / 100 * innerHeight
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
						vPxUsed += b
						// ok
					case nux.Auto, nux.Unlimit:
						vPt += cs.Margin().Bottom.Value()
						// ok, wait until percent grand total
					}
				}

				switch cs.Margin().Left.Mode() {
				case nux.Pixel:
					l := cs.Margin().Left.Value()
					cms.Margin.Left = util.Roundi32(l)
					hPx += l
					hPxUsed += l
					// ok
				case nux.Weight:
					switch width.Mode() {
					case nux.Pixel:
						hWt += cs.Margin().Left.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						// ok, cms.Margin.Left = 0, do nothing
					}
				case nux.Percent:
					switch width.Mode() {
					case nux.Pixel:
						l := cs.Margin().Left.Value() / 100 * innerWidth
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
						hPxUsed += l
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.Margin().Left.Value()
						// ok, wait until percent grand total
					}
				}

				switch cs.Margin().Right.Mode() {
				case nux.Pixel:
					r := cs.Margin().Right.Value()
					cms.Margin.Right = int32(r)
					hPx += r
					hPxUsed += r
					// ok
				case nux.Weight:
					switch width.Mode() {
					case nux.Pixel:
						hWt += cs.Margin().Right.Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						// ok, cms.Margin.Right = 0, do nothing
					}
				case nux.Percent:
					switch width.Mode() {
					case nux.Pixel:
						r := cs.Margin().Right.Value() / 100 * innerWidth
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
						hPxUsed += r
						// ok
					case nux.Auto, nux.Unlimit:
						hPt += cs.Margin().Right.Value()
						// ok, wait until percent grand total
					}
				}
			}

			hPxRemain := innerWidth - hPxUsed
			vPxRemain := innerHeight - vPxUsed

			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cms.Width = util.Roundi32(w)
				setRatioHeightIfNeed(cs, cms, w, nux.Pixel)
				// ok
			case nux.Weight:
				wt := hWt + cs.Width().Value()
				w := cs.Width().Value() / wt * hPxRemain
				switch width.Mode() {
				case nux.Pixel:
					cms.Width = util.Roundi32(w)
				case nux.Auto, nux.Unlimit:
					cms.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Unlimit))
				}
				setRatioHeightIfNeed(cs, cms, w, nux.Pixel)
			case nux.Percent:
				switch width.Mode() {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = util.Roundi32(w)
					setRatioHeightIfNeed(cs, cms, w, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					if hPx > 0 {
						hPt += cs.Width().Value() // ?????????? did width add to hPt ok?
						_innerWidth := hPx / (1.0 - hPt/100.0)
						w := cs.Width().Value() / 100 * _innerWidth
						cms.Width = util.Roundi32(w)
						setRatioHeightIfNeed(cs, cms, w, nux.Pixel)
					} else {
						w := cs.Width().Value() / 100 * innerWidth
						cms.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Unlimit))
						setRatioHeightIfNeed(cs, cms, w, nux.Pixel)
					}
				}
			case nux.Ratio:
				// ok, measured when measure height
				if cs.Height().Mode() == nux.Ratio {
					log.Fatal("nux", "width and height size mode can not both 'Ratio'")
				}
			case nux.Auto, nux.Unlimit:
				cms.Width = int32(nux.MeasureSpec(util.Roundi32(hPxRemain), nux.Unlimit))
				setRatioHeightIfNeed(cs, cms, hPxRemain, nux.Pixel)
				// ok
			}

			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cms.Height = util.Roundi32(h)
				setRatioWidthIfNeed(cs, cms, h, nux.Pixel)
				// ok
			case nux.Weight:
				wt := vWt + cs.Height().Value()
				h := cs.Height().Value() / wt * vPxRemain
				switch height.Mode() {
				case nux.Pixel:
					cms.Height = util.Roundi32(h)
				case nux.Auto, nux.Unlimit:
					cms.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Unlimit))
				}
				setRatioWidthIfNeed(cs, cms, h, nux.Pixel)
			case nux.Percent:
				switch height.Mode() {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cms.Height = util.Roundi32(h)
					setRatioWidthIfNeed(cs, cms, h, nux.Pixel)
				case nux.Auto, nux.Unlimit:
					if vPx > 0 {
						vPt += cs.Height().Value()
						_innerHeight := vPx / (1.0 - vPt/100.0)
						h := cs.Height().Value() / 100 * _innerHeight
						cms.Height = util.Roundi32(h)
						setRatioWidthIfNeed(cs, cms, h, nux.Pixel)
					} else {
						h := cs.Height().Value() / 100 * innerHeight
						cms.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Unlimit))
						setRatioWidthIfNeed(cs, cms, h, nux.Pixel)
					}
					// ok
				}
			case nux.Ratio:
				if cs.Width().Mode() == nux.Ratio {
					log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
				}
			case nux.Auto, nux.Unlimit:
				cms.Height = int32(nux.MeasureSpec(util.Roundi32(vPxRemain), nux.Unlimit))
				setRatioWidthIfNeed(cs, cms, vPxRemain, nux.Pixel)
				// ok
			}

			if m, ok := child.(nux.Measure); ok {
				m.Measure(nux.MeasureDimen(cms.Width), nux.MeasureDimen(cms.Height))

				hPx += float32(cms.Width)
				vPx += float32(cms.Height)
				hPxRemain -= float32(cms.Width)
				vPxRemain -= float32(cms.Height)
			} else {
				// cms.Width = 0
				// cms.Height = 0
			}

			if cs.Margin() != nil {
				switch cs.Margin().Top.Mode() {
				// case nux.Pixel: has measured
				case nux.Weight:
					switch height.Mode() {
					case nux.Pixel:
						t := cs.Margin().Top.Value() / vWt * vPxRemain
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Top = 0, do nothing
					}
				case nux.Percent:
					switch height.Mode() {
					// case nux.Pixel: has measured
					case nux.Auto, nux.Unlimit:
						if vPx > 0 {
							_innerHeight := vPx / (1.0 - vPt/100.0)
							t := cs.Margin().Top.Value() / 100 * _innerHeight
							cms.Margin.Top = util.Roundi32(t)
							vPx += t
							// ok
						} else {
							// cms.Margin.Top = 0, do nothing
						}
					}
				}

				switch cs.Margin().Bottom.Mode() {
				case nux.Weight:
					switch height.Mode() {
					case nux.Pixel:
						b := cs.Margin().Bottom.Value() / vWt * vPxRemain
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Bottom = 0, do nothing
					}
				case nux.Percent:
					switch height.Mode() {
					// case nux.Pixel: has measured
					case nux.Auto, nux.Unlimit:
						if vPx > 0 {
							_innerHeight := vPx / (1.0 - vPt/100.0)
							b := cs.Margin().Bottom.Value() / 100 * _innerHeight
							cms.Margin.Bottom = util.Roundi32(b)
							vPx += b
							// ok
						} else {
							// cms.Margin.Bottom = 0, do nothing
						}
					}
				}

				switch cs.Margin().Left.Mode() {
				case nux.Weight:
					switch width.Mode() {
					case nux.Pixel:
						l := cs.Margin().Left.Value() / hWt * hPxRemain
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Left = 0, do nothing
					}
				case nux.Percent:
					switch width.Mode() {
					// case nux.Pixel: has measured
					case nux.Auto, nux.Unlimit:
						if hPx > 0 {
							_innerWidth := hPx / (1.0 - hPt/100.0)
							l := cs.Margin().Left.Value() / 100 * _innerWidth
							cms.Margin.Left = util.Roundi32(l)
							hPx += l
							// ok
						} else {
							// cms.Margin.Left = 0, do nothing
						}
					}
				}

				switch cs.Margin().Right.Mode() {
				case nux.Weight:
					switch width.Mode() {
					case nux.Pixel:
						r := cs.Margin().Right.Value() / hWt * hPxRemain
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Right = 0, do nothing
					}
				case nux.Percent:
					switch width.Mode() {
					// case nux.Pixel: has measured
					case nux.Auto, nux.Unlimit:
						if hPx > 0 {
							_innerWidth := hPx / (1.0 - hPt/100.0)
							r := cs.Margin().Right.Value() / 100 * _innerWidth
							cms.Margin.Right = util.Roundi32(r)
							hPx += r
							// ok
						} else {
							// cms.Margin.Right = 0, do nothing
						}
					}
				}
			}

			hPx = hPx / (1.0 - hPt/100.0)
			vPx = vPx / (1.0 - vPt/100.0)

		}
	} // end one child

	switch width.Mode() {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPx
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

	switch height.Mode() {
	case nux.Auto, nux.Unlimit:
		innerHeight = vPx
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

	setNewWidth(frame, originWidth, width)
	setNewHeight(frame, originHeight, height)
}

func (me *scroll) Layout(x, y, width, height int32) {
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
			cms := cs.Frame()

			l += float32(cms.Margin.Left)
			t += float32(cms.Margin.Top)

			cms.X = x + int32(l)
			cms.Y = y + int32(t)

			if cl, ok := child.(nux.Layout); ok {
				cl.Layout(cms.X, cms.Y, cms.Width, cms.Height)
			}
		}
	}
}

func (me *scroll) Draw(canvas nux.Canvas) {
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
