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
	ms := me.MeasuredSize()
	measuredIndex := map[int]struct{}{}
	if nux.MeasureSpecMode(width) == nux.Pixel && nux.MeasureSpecMode(height) == nux.Pixel {
		me.measure(width, height, measuredIndex)
	} else {
		w, h := me.measure(width, height, measuredIndex)
		w, h = me.measure(w, h, measuredIndex)
		ms.Width = w
		ms.Height = h
	}
}

func (me *layer) measure(width, height int32, measuredIndex map[int]struct{}) (outWidth, outHeight int32) {
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

	if me.HasPadding() {
		switch me.PaddingLeft().Mode() {
		case nux.Pixel:
			l := me.PaddingLeft().Value()
			ms.Padding.Left = util.Roundi32(l)
			hPPx += l
			hPPxUsed += l
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.PaddingLeft().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
			case nux.Auto:
				hPPt += me.PaddingLeft().Value()
			}
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			r := me.PaddingRight().Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
			hPPxUsed += r
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.PaddingRight().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
			case nux.Auto:
				hPPt += me.PaddingRight().Value()
			}
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			t := me.PaddingTop().Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
			vPPxUsed += t
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.PaddingTop().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
			case nux.Auto:
				vPPt += me.PaddingTop().Value()
			}
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			b := me.PaddingBottom().Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
			vPPxUsed += b
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.PaddingBottom().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
			case nux.Auto:
				vPPt += me.PaddingBottom().Value()
			}
		}
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

		// TODO if child visible == gone , then skip

		if cs, ok := child.(nux.Size); ok {
			cms := cs.MeasuredSize()
			if _, ok := measuredIndex[index]; !ok {
				cms.Width = nux.MeasureSpec(0, nux.Unlimit)
				cms.Height = nux.MeasureSpec(0, nux.Unlimit)
			}

			if cs.HasMargin() {
				switch cs.MarginTop().Mode() {
				case nux.Pixel:
					t := cs.MarginTop().Value()
					vPx += t
					vPxUsed += t
					cms.Margin.Top = util.Roundi32(t)
				case nux.Weight:
					vWt += cs.MarginTop().Value()
				case nux.Percent:
					t := cs.MarginTop().Value() / 100 * innerHeight
					cms.Margin.Top = util.Roundi32(t)
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vPx += t
						vPxUsed += t
					case nux.Auto:
						vPxUsed += t
						vPt += cs.MarginTop().Value()
					}

				}

				switch cs.MarginBottom().Mode() {
				case nux.Pixel:
					b := cs.MarginBottom().Value()
					vPx += b
					vPxUsed += b
					cms.Margin.Bottom = int32(b)
				case nux.Weight:
					vWt += cs.MarginBottom().Value()
				case nux.Percent:
					b := cs.MarginBottom().Value() / 100 * innerHeight
					cms.Margin.Bottom = util.Roundi32(b)
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						vPx += b
						vPxUsed += b
					case nux.Auto:
						vPxUsed += b
						vPt += cs.MarginBottom().Value()
					}

				}

				switch cs.MarginLeft().Mode() {
				case nux.Pixel:
					l := cs.MarginLeft().Value()
					hPx += l
					hPxUsed += l
					cms.Margin.Left = util.Roundi32(l)
				case nux.Weight:
					hWt += cs.MarginLeft().Value()
				case nux.Percent:
					l := cs.MarginLeft().Value() / 100.0 * innerWidth
					cms.Margin.Left = util.Roundi32(l)
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hPx += l
						hPxUsed += l
					case nux.Auto:
						hPxUsed += l
						hPt += cs.MarginLeft().Value()
					}

				}

				switch cs.MarginRight().Mode() {
				case nux.Pixel:
					r := cs.MarginRight().Value()
					hPx += r
					hPxUsed += r
					cms.Margin.Right = int32(r)
				case nux.Weight:
					hWt += cs.MarginRight().Value()
				case nux.Percent:
					r := cs.MarginRight().Value() / 100.0 * innerWidth
					cms.Margin.Right = util.Roundi32(r)
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
						hPx += r
						hPxUsed += r
					case nux.Auto:
						hPxUsed += r
						hPt += cs.MarginRight().Value()
					}

				}
			}

			hPxRemain := innerWidth - hPxUsed
			vPxRemain := innerHeight - vPxUsed

			measured := false
			if _, ok := measuredIndex[index]; !ok {
				switch cs.Width().Mode() {
				case nux.Pixel:
					w := cs.Width().Value()
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, w, nux.Pixel)
				case nux.Weight:
					hWt += cs.Width().Value()
					if nux.MeasureSpecMode(width) == nux.Pixel {
						w := cs.Width().Value() / hWt * hPxRemain
						hWt -= cs.Width().Value()
						// hPxRemain -= w
						cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
						setRatioHeight(cs, cms, w, nux.Pixel)
					}
				case nux.Percent:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = nux.MeasureSpec(util.Roundi32(w), nux.Pixel)
					setRatioHeight(cs, cms, w, nux.Pixel)
					switch nux.MeasureSpecMode(width) {
					case nux.Pixel:
					case nux.Auto:
						hPt += cs.Width().Value()
					}
				case nux.Ratio:
					// measured when measure height
					if cs.Height().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
					}
				case nux.Auto:
					cms.Width = nux.MeasureSpec(util.Roundi32(hPxRemain), nux.Auto)
					setRatioHeight(cs, cms, hPxRemain, nux.Pixel)
				case nux.Unlimit:
					// ignore
				}

				switch cs.Height().Mode() {
				case nux.Pixel:
					h := cs.Height().Value()
					cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
					setRatioWidth(cs, cms, h, nux.Pixel)
				case nux.Weight:
					vWt += cs.Height().Value()
					if nux.MeasureSpecMode(height) == nux.Pixel {
						h := cs.Height().Value() / vWt * vPxRemain
						vWt -= cs.Height().Value()
						// vPxRemain -= h
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cms, h, nux.Pixel)
					}
				case nux.Percent:
					switch nux.MeasureSpecMode(height) {
					case nux.Pixel:
						h := cs.Height().Value() / 100 * innerHeight
						cms.Height = nux.MeasureSpec(util.Roundi32(h), nux.Pixel)
						setRatioWidth(cs, cms, h, nux.Pixel)
					case nux.Auto:
						vPt += cs.Height().Value()
					}
				case nux.Ratio:
					if cs.Width().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
					}
				case nux.Auto:
					cms.Height = nux.MeasureSpec(util.Roundi32(vPxRemain), nux.Auto)
					setRatioWidth(cs, cms, vPxRemain, nux.Pixel)
				case nux.Unlimit:
					// nothing
				}

				if nux.MeasureSpecMode(cms.Height) != nux.Unlimit && nux.MeasureSpecMode(cms.Width) != nux.Unlimit {
					if m, ok := child.(nux.Measure); ok {
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

						measured = true
						measuredIndex[index] = struct{}{}

						hPx += float32(cms.Width)
						vPx += float32(cms.Height)
						hPxRemain -= float32(cms.Width)
						vPxRemain -= float32(cms.Height)

						hPx = hPx / (1.0 - hPt/100.0)
						vPx = vPx / (1.0 - vPt/100.0)

						if hPx > hPxMax {
							hPxMax = hPx
						}

						if vPx > vPxMax {
							vPxMax = vPx
						}
					}
				}
			} else {
				measured = true
				hPx += float32(cms.Width)
				vPx += float32(cms.Height)
				hPxRemain -= float32(cms.Width)
				vPxRemain -= float32(cms.Height)

				hPx = hPx / (1.0 - hPt/100.0)
				vPx = vPx / (1.0 - vPt/100.0)

				if hPx > hPxMax {
					hPxMax = hPx
				}

				if vPx > vPxMax {
					vPxMax = vPx
				}
			}

			if measured {
				if cs.HasMargin() {
					switch cs.MarginTop().Mode() {
					case nux.Weight:
						if nux.MeasureSpecMode(height) == nux.Pixel {
							t := cs.MarginTop().Value() / vWt * vPxRemain
							cms.Margin.Top = util.Roundi32(t)
						}
					case nux.Percent:
						if nux.MeasureSpecMode(height) == nux.Auto {
							t := cs.MarginTop().Value() / 100 * vPx
							cms.Margin.Top = util.Roundi32(t)
						}
					}

					switch cs.MarginBottom().Mode() {
					case nux.Weight:
						if nux.MeasureSpecMode(height) == nux.Pixel {
							t := cs.MarginBottom().Value() / vWt * vPxRemain
							cms.Margin.Bottom = util.Roundi32(t)
						}
					case nux.Percent:
						if nux.MeasureSpecMode(height) == nux.Auto {
							b := cs.MarginBottom().Value() / 100.0 * vPx
							cms.Margin.Bottom = util.Roundi32(b)
						}
					}

					switch cs.MarginLeft().Mode() {
					case nux.Weight:
						if nux.MeasureSpecMode(width) == nux.Pixel {
							l := cs.MarginLeft().Value() / hWt * hPxRemain
							cms.Margin.Left = util.Roundi32(l)
						}
					case nux.Percent:
						if nux.MeasureSpecMode(width) == nux.Auto {
							l := cs.MarginLeft().Value() / 100.0 * hPx
							cms.Margin.Left = util.Roundi32(l)
						}
					}

					switch cs.MarginRight().Mode() {
					case nux.Weight:
						if nux.MeasureSpecMode(width) == nux.Pixel {
							r := cs.MarginRight().Value() / hWt * hPxRemain
							cms.Margin.Right = util.Roundi32(r)
						}
					case nux.Percent:
						if nux.MeasureSpecMode(width) == nux.Auto {
							r := cs.MarginRight().Value() / 100.0 * hPx
							cms.Margin.Right = util.Roundi32(r)
						}
					}
				}
			}
		}
	}

	if nux.MeasureSpecMode(width) == nux.Auto {
		outWidth = util.Roundi32((hPxMax + vPPx) / (1.0 - vPPt/100.0))
	} else {
		outWidth = width
	}

	if nux.MeasureSpecMode(height) == nux.Auto {
		outHeight = util.Roundi32((vPxMax + vPPx) / (1.0 - vPPt/100.0))
	} else {
		outHeight = height
	}

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
