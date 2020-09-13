// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/util"
)

type Decor interface {
	Parent
	Size
}

type decor struct {
	WidgetParent
	WidgetBase
	WidgetSize
}

func NewDecor() Decor {
	me := &decor{}
	me.WidgetSize.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	return me
}

func (me *decor) Creating(attr Attr) {
	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetParent.Creating(attr)
}

func (me *decor) onSizeChanged(widget Widget) {

}
func (me *decor) onVisualChanged(widget Widget) {

}

func (me *decor) Measure(width, height int32) {
	ms := me.MeasuredSize()
	measuredIndex := map[int]struct{}{}
	if MeasureSpecMode(width) == Auto || MeasureSpecMode(height) == Auto {
		w, h := me.measure(width, height, measuredIndex)
		w, h = me.measure(w, h, measuredIndex)
		ms.Width = w
		ms.Height = h
	} else if MeasureSpecMode(width) == Pixel && MeasureSpecMode(height) == Pixel {
		me.measure(width, height, measuredIndex)
	}
}

func (me *decor) measure(width, height int32, measuredIndex map[int]struct{}) (outWidth, outHeight int32) {
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

	// 1. measure self padding size
	if me.HasPadding() {
		switch me.padding.Left.Mode() {
		case Pixel:
			l := me.padding.Left.Value()
			ms.Padding.Left = util.Roundi32(l)
			hPPx += l
			hPPxUsed += l
		case Percent:
			switch MeasureSpecMode(width) {
			case Pixel:
				l := me.padding.Left.Value() / 100 * float32(MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
				hPPxUsed += l
			case Auto:
				hPPt += me.padding.Left.Value()
			}
		}

		switch me.padding.Right.Mode() {
		case Pixel:
			r := me.padding.Right.Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
			hPPxUsed += r
		case Percent:
			switch MeasureSpecMode(width) {
			case Pixel:
				r := me.padding.Right.Value() / 100 * float32(MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
				hPPxUsed += r
			case Auto:
				hPPt += me.padding.Right.Value()
			}
		}

		switch me.padding.Top.Mode() {
		case Pixel:
			t := me.padding.Top.Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
			vPPxUsed += t
		case Percent:
			switch MeasureSpecMode(height) {
			case Pixel:
				t := me.padding.Top.Value() / 100 * float32(MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
				vPPxUsed += t
			case Auto:
				vPPt += me.padding.Top.Value()
			}
		}

		switch me.padding.Bottom.Mode() {
		case Pixel:
			b := me.padding.Bottom.Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
			vPPxUsed += b
		case Percent:
			switch MeasureSpecMode(height) {
			case Pixel:
				b := me.padding.Bottom.Value() / 100 * float32(MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
				vPPxUsed += b
			case Auto:
				vPPt += me.padding.Bottom.Value()
			}
		}
	}

	innerWidth = float32(MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight = float32(MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	var hPxUsed float32
	var hPx float32
	var hWt float32
	var hPt float32

	var vPxUsed float32
	var vPx float32
	var vWt float32
	var vPt float32

	for index, child := range me.Children() {
		if compt, ok := child.(Component); ok {
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

		if cs, ok := child.(Size); ok {
			cms := cs.MeasuredSize()
			if _, ok := measuredIndex[index]; !ok {
				cms.Width = MeasureSpec(0, Unspec)
				cms.Height = MeasureSpec(0, Unspec)
			}

			if cs.HasMargin() {
				switch cs.MarginTop().Mode() {
				case Pixel:
					t := cs.MarginTop().Value()
					vPx += t
					vPxUsed += t
					cms.Margin.Top = util.Roundi32(t)
				case Weight:
					vWt += cs.MarginTop().Value()
				case Percent:
					t := cs.MarginTop().Value() / 100 * innerHeight
					cms.Margin.Top = util.Roundi32(t)
					switch MeasureSpecMode(height) {
					case Pixel:
						vPx += t
						vPxUsed += t
					case Auto:
						vPxUsed += t
						vPt += cs.MarginTop().Value()
					}

				}

				switch cs.MarginBottom().Mode() {
				case Pixel:
					b := cs.MarginBottom().Value()
					vPx += b
					vPxUsed += b
					cms.Margin.Bottom = int32(b)
				case Weight:
					vWt += cs.MarginBottom().Value()
				case Percent:
					b := cs.MarginBottom().Value() / 100 * innerHeight
					cms.Margin.Bottom = util.Roundi32(b)
					switch MeasureSpecMode(height) {
					case Pixel:
						vPx += b
						vPxUsed += b
					case Auto:
						vPxUsed += b
						vPt += cs.MarginBottom().Value()
					}

				}

				switch cs.MarginLeft().Mode() {
				case Pixel:
					l := cs.MarginLeft().Value()
					hPx += l
					hPxUsed += l
					cms.Margin.Left = util.Roundi32(l)
				case Weight:
					hWt += cs.MarginLeft().Value()
				case Percent:
					l := cs.MarginLeft().Value() / 100.0 * innerWidth
					cms.Margin.Left = util.Roundi32(l)
					switch MeasureSpecMode(width) {
					case Pixel:
						hPx += l
						hPxUsed += l
					case Auto:
						hPxUsed += l
						hPt += cs.MarginLeft().Value()
					}

				}

				switch cs.MarginRight().Mode() {
				case Pixel:
					r := cs.MarginRight().Value()
					hPx += r
					hPxUsed += r
					cms.Margin.Right = int32(r)
				case Weight:
					hWt += cs.MarginRight().Value()
				case Percent:
					r := cs.MarginRight().Value() / 100.0 * innerWidth
					cms.Margin.Right = util.Roundi32(r)
					switch MeasureSpecMode(width) {
					case Pixel:
						hPx += r
						hPxUsed += r
					case Auto:
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
				case Pixel:
					w := cs.Width().Value()
					cms.Width = MeasureSpec(util.Roundi32(w), Pixel)
					setRatioHeight(cs, cms, w, Pixel)
				case Weight:
					hWt += cs.Width().Value()
					if MeasureSpecMode(width) == Pixel {
						w := cs.Width().Value() / hWt * hPxRemain
						hWt -= cs.Width().Value()
						// hPxRemain -= w
						cms.Width = MeasureSpec(util.Roundi32(w), Pixel)
						setRatioHeight(cs, cms, w, Pixel)
					}
				case Percent:
					w := cs.Width().Value() / 100 * innerWidth
					cms.Width = MeasureSpec(util.Roundi32(w), Pixel)
					setRatioHeight(cs, cms, w, Pixel)
					switch MeasureSpecMode(width) {
					case Pixel:
					case Auto:
						hPt += cs.Width().Value()
					}
				case Ratio:
					// measured when measure height
					if cs.Height().Mode() == Ratio {
						log.Fatal("nux", "width and height size mode can not both Ratio, at least one is definited.")
					}
				case Auto:
					cms.Width = MeasureSpec(util.Roundi32(hPxRemain), Auto)
					setRatioHeight(cs, cms, hPxRemain, Pixel)
				case Default, Unspec:
					// ignore
				}

				switch cs.Height().Mode() {
				case Pixel:
					h := cs.Height().Value()
					cms.Height = MeasureSpec(util.Roundi32(h), Pixel)
					setRatioWidth(cs, cms, h, Pixel)
				case Weight:
					vWt += cs.Height().Value()
					if MeasureSpecMode(height) == Pixel {
						h := cs.Height().Value() / vWt * vPxRemain
						vWt -= cs.Height().Value()
						// vPxRemain -= h
						cms.Height = MeasureSpec(util.Roundi32(h), Pixel)
						setRatioWidth(cs, cms, h, Pixel)
					}
				case Percent:
					switch MeasureSpecMode(height) {
					case Pixel:
						h := cs.Height().Value() / 100 * innerHeight
						cms.Height = MeasureSpec(util.Roundi32(h), Pixel)
						setRatioWidth(cs, cms, h, Pixel)
					case Auto:
						vPt += cs.Height().Value()
					}
				case Ratio:
					if cs.Width().Mode() == Ratio {
						log.Fatal("nux", "width and height size mode can not both Ratio, at least one is definited.")
					}
				case Auto:
					cms.Height = MeasureSpec(util.Roundi32(vPxRemain), Auto)
					setRatioWidth(cs, cms, vPxRemain, Pixel)
				case Default, Unspec:
					// nothing
				}

				if MeasureSpecMode(cms.Height) != Unspec && MeasureSpecMode(cms.Width) != Unspec {
					if m, ok := child.(Measure); ok {
						m.Measure(cms.Width, cms.Height)

						if cs.Width().Mode() == Ratio {
							oldWidth := cms.Width
							cms.Width = MeasureSpec(util.Roundi32(float32(cms.Height)*cs.Width().Value()), Pixel)
							if oldWidth != cms.Width {
								m.Measure(cms.Width, cms.Height)
							}
						}

						if cs.Height().Mode() == Ratio {
							oldHeight := cms.Height
							cms.Height = MeasureSpec(util.Roundi32(float32(cms.Width)/cs.Height().Value()), Pixel)
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
					case Weight:
						if MeasureSpecMode(height) == Pixel {
							t := cs.MarginTop().Value() / vWt * vPxRemain
							cms.Margin.Top = util.Roundi32(t)
						}
					case Percent:
						if MeasureSpecMode(height) == Auto {
							t := cs.MarginTop().Value() / 100 * vPx
							cms.Margin.Top = util.Roundi32(t)
						}
					}

					switch cs.MarginBottom().Mode() {
					case Weight:
						if MeasureSpecMode(height) == Pixel {
							t := cs.MarginBottom().Value() / vWt * vPxRemain
							cms.Margin.Bottom = util.Roundi32(t)
						}
					case Percent:
						if MeasureSpecMode(height) == Auto {
							b := cs.MarginBottom().Value() / 100.0 * vPx
							cms.Margin.Bottom = util.Roundi32(b)
						}
					}

					switch cs.MarginLeft().Mode() {
					case Weight:
						if MeasureSpecMode(width) == Pixel {
							l := cs.MarginLeft().Value() / hWt * hPxRemain
							cms.Margin.Left = util.Roundi32(l)
						}
					case Percent:
						if MeasureSpecMode(width) == Auto {
							l := cs.MarginLeft().Value() / 100.0 * hPx
							cms.Margin.Left = util.Roundi32(l)
						}
					}

					switch cs.MarginRight().Mode() {
					case Weight:
						if MeasureSpecMode(width) == Pixel {
							r := cs.MarginRight().Value() / hWt * hPxRemain
							cms.Margin.Right = util.Roundi32(r)
						}
					case Percent:
						if MeasureSpecMode(width) == Auto {
							r := cs.MarginRight().Value() / 100.0 * hPx
							cms.Margin.Right = util.Roundi32(r)
						}
					}
				}
			}
		}
	}

	if MeasureSpecMode(width) == Auto {
		outWidth = util.Roundi32((hPxMax + vPPx) / (1.0 - vPPt/100.0))
	} else {
		outWidth = width
	}

	if MeasureSpecMode(height) == Auto {
		outHeight = util.Roundi32((vPxMax + vPPx) / (1.0 - vPPt/100.0))
	} else {
		outHeight = height
	}

	return
}

// meature position
func (me *decor) Layout(dx, dy, left, top, right, bottom int32) {
	ms := me.MeasuredSize()

	// set frame for self
	ms.Position.Left = left
	ms.Position.Top = top
	ms.Position.Right = right
	ms.Position.Bottom = bottom
	ms.Position.X = dx
	ms.Position.Y = dy

	var l, t float32

	for _, child := range me.Children() {
		if compt, ok := child.(Component); ok {
			child = compt.Content()
		}

		l = float32(ms.Padding.Left)
		t = float32(ms.Padding.Top)
		// TODO if child visible == gone , then skip
		if cs, ok := child.(Size); ok {
			if cl, ok := child.(Layout); ok {
				cms := cs.MeasuredSize()

				l += float32(cms.Margin.Left)
				t += float32(cms.Margin.Top)

				cl.Layout(dx+int32(l), dy+int32(t), int32(l), int32(t), int32(l)+cms.Width, int32(t)+cms.Height)
			}
		}
	}
}

func (me *decor) Draw(canvas Canvas) {
	log.V("nux", "Decor Draw....")
	// if me.Background() != nil {
	// 	me.Background().Draw(canvas)
	// }
	for _, child := range me.Children() {
		if compt, ok := child.(Component); ok {
			child = compt.Content()
		}

		// TODO if child visible == gone , then skip
		if cs, ok := child.(Size); ok {
			cms := cs.MeasuredSize()
			if draw, ok := child.(Draw); ok {
				canvas.Save()
				canvas.Translate(cms.Position.Left, cms.Position.Top)
				canvas.ClipRect(0, 0, cms.Width, cms.Height)
				// t1 := time.Now()
				draw.Draw(canvas)
				// log.V("nux", "draw used time %d", time.Now().Sub(t1).Milliseconds())
				canvas.Restore()
			}
		}
	}

	// if me.Foreground() != nil {
	// 	me.Foreground().Draw(canvas)
	// }
}

func setRatioHeight(cs Size, cms *MeasuredSize, width float32, mode Mode) {
	if cs.Height().Mode() == Ratio {
		cms.Height = MeasureSpec(util.Roundi32(width/cs.Height().Value()), mode)
	}
}

func setRatioWidth(cs Size, cms *MeasuredSize, height float32, mode Mode) {
	if cs.Width().Mode() == Ratio {
		cms.Width = MeasureSpec(util.Roundi32(height*cs.Width().Value()), mode)
	}
}
