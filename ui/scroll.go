// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
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

func NewScroll(context nux.Context, attrs ...nux.Attr) Scroll {
	me := &scroll{}
	me.WidgetParent = nux.NewWidgetParent(context, me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(context, me, attrs...)
	me.WidgetVisual = NewWidgetVisual(context, me, attrs...)
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

func (me *scroll) onSizeChanged(widget nux.Widget) {

}
func (me *scroll) onVisualChanged(widget nux.Widget) {

}

func (me *scroll) Measure(width, height int32) {
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
				// ok, wait until width measured
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
				// ok, wait until width measured
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

	innerWidth = float32(nux.MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(nux.MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	if len(me.Children()) == 1 {
		child := me.Children()[0]
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()

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
						// ok, cms.Margin.Top = 0, do nothing
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
						// ok, cms.Margin.Bottom = 0, do nothing
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
						// ok, cms.Margin.Left = 0, do nothing
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
					cms.Margin.Right = int32(r)
					hPx += r
					hPxUsed += r
					// ok
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hWt += cs.MarginRight().Value()
						// ok, wait until weight grand total
					case nux.Auto, nux.Unlimit:
						// ok, cms.Margin.Right = 0, do nothing
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

			hPxRemain := innerWidth - hPxUsed
			vPxRemain := innerHeight - vPxUsed

			switch cs.Width().Mode() {
			case nux.Pixel:
				w := cs.Width().Value()
				cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				setRatioHeight(cs, cms, w, nux.Pixel)
				// ok
			case nux.Weight:
				wt := hWt + cs.Width().Value()
				w := cs.Width().Value() / wt * hPxRemain
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
				case nux.Auto, nux.Unlimit:
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Unlimit)
				}
				setRatioHeight(cs, cms, w, nux.Pixel)
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, w, nux.Pixel)
					// ok
				case nux.Auto, nux.Unlimit:
					if hPx > 0 {
						hPt += cs.Width().Value() // ?????????? did width add to hPt ok?
						_innerWidth := hPx / (1.0 - hPt/100.0)
						w := cs.Width().Value() / 100 * _innerWidth
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cms, w, nux.Pixel)
					} else {
						w := cs.Width().Value() / 100 * innerWidth
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Unlimit)
						setRatioHeight(cs, cms, w, nux.Pixel)
					}
				}
			case nux.Ratio:
				// ok, measured when measure height
				if cs.Height().Mode() == nux.Ratio {
					log.Fatal("nux", "width and height size mode can not both 'Ratio'")
				}
			case nux.Auto, nux.Unlimit:
				cms.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), nux.Unlimit)
				setRatioHeight(cs, cms, hPxRemain, nux.Pixel)
				// ok
			}

			switch cs.Height().Mode() {
			case nux.Pixel:
				h := cs.Height().Value()
				cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				setRatioWidth(cs, cms, h, nux.Pixel)
				// ok
			case nux.Weight:
				wt := vWt + cs.Height().Value()
				h := cs.Height().Value() / wt * vPxRemain
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
				case nux.Auto, nux.Unlimit:
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Unlimit)
				}
				setRatioWidth(cs, cms, h, nux.Pixel)
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := cs.Height().Value() / 100 * innerHeight
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, h, nux.Pixel)
				case nux.Auto, nux.Unlimit:
					if vPx > 0 {
						vPt += cs.Height().Value()
						_innerHeight := vPx / (1.0 - vPt/100.0)
						h := cs.Height().Value() / 100 * _innerHeight
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cms, h, nux.Pixel)
					} else {
						h := cs.Height().Value() / 100 * innerHeight
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Unlimit)
						setRatioWidth(cs, cms, h, nux.Pixel)
					}
					// ok
				}
			case nux.Ratio:
				if cs.Width().Mode() == nux.Ratio {
					log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
				}
			case nux.Auto, nux.Unlimit:
				cms.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), nux.Unlimit)
				setRatioWidth(cs, cms, vPxRemain, nux.Pixel)
				// ok
			}

			if m, ok := child.(nux.Measure); ok {
				m.Measure(cms.Width, cms.Height)

				hPx += float32(cms.Width)
				vPx += float32(cms.Height)
				hPxRemain -= float32(cms.Width)
				vPxRemain -= float32(cms.Height)
			} else {
				// cms.Width = 0
				// cms.Height = 0
			}

			if cs.HasMargin() {
				switch cs.MarginTop().Mode() {
				// case nux.Pixel: has measured
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						t := cs.MarginTop().Value() / vWt * vPxRemain
						cms.Margin.Top = util.Roundi32(t)
						vPx += t
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Top = 0, do nothing
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
							// cms.Margin.Top = 0, do nothing
						}
					}
				}

				switch cs.MarginBottom().Mode() {
				case nux.Weight:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						b := cs.MarginBottom().Value() / vWt * vPxRemain
						cms.Margin.Bottom = util.Roundi32(b)
						vPx += b
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Bottom = 0, do nothing
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
							// cms.Margin.Bottom = 0, do nothing
						}
					}
				}

				switch cs.MarginLeft().Mode() {
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						l := cs.MarginLeft().Value() / hWt * hPxRemain
						cms.Margin.Left = util.Roundi32(l)
						hPx += l
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Left = 0, do nothing
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
							// cms.Margin.Left = 0, do nothing
						}
					}
				}

				switch cs.MarginRight().Mode() {
				case nux.Weight:
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						r := cs.MarginRight().Value() / hWt * hPxRemain
						cms.Margin.Right = util.Roundi32(r)
						hPx += r
					case nux.Auto, nux.Unlimit:
						// cms.Margin.Right = 0, do nothing
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
							// cms.Margin.Right = 0, do nothing
						}
					}
				}
			}

			hPx = hPx / (1.0 - hPt/100.0)
			vPx = vPx / (1.0 - vPt/100.0)

		}
	} // end one child

	switch nux.MeasureSpecMode(width) {
	case nux.Auto, nux.Unlimit:
		innerWidth = hPx
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
		innerHeight = vPx
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

	setNewWidth(ms, originWidth, width)
	setNewHeight(ms, originHeight, height)
}

func (me *scroll) Layout(dx, dy, left, top, right, bottom int32) {
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

func (me *scroll) Draw(canvas nux.Canvas) {
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
				canvas.Translate(float32(cms.Position.Left), float32(cms.Position.Top))
				canvas.ClipRect(0, 0, float32(cms.Width), float32(cms.Height))
				draw.Draw(canvas)
				canvas.Restore()
			}
		}
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
