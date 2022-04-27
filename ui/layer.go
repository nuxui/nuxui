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
	childrenWidth  float32
	childrenHeight float32
	clipChildren   int
}

func NewLayer(attr nux.Attr) Layer {
	me := &layer{}
	me.WidgetParent = nux.NewWidgetParent(me, attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *layer) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *layer) Measure(width, height nux.MeasureDimen) {
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

	hPxMax, vPxMax, measuredCompleteCount := me.measureChildren(width, height, innerWidth, innerHeight, childrenMeasuredFlags)

	switch width.Mode() {
	case nux.Pixel:
		frame.Width = width.Value()
		me.childrenWidth = hPxMax
	case nux.Auto, nux.Unlimit:
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

	switch height.Mode() {
	case nux.Pixel:
		frame.Height = height.Value()
		me.childrenHeight = vPxMax
	case nux.Auto, nux.Unlimit:
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

	if measuredCompleteCount == me.ChildrenCount() {
		if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
			measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
		}
		return
	}

	measuredCompleteCount = me.measureChildrenMargin(width, height, innerWidth, innerHeight, childrenMeasuredFlags)

	if measuredCompleteCount != me.ChildrenCount() {
		log.Fatal("nuxui", "can not run here")
	}
}

// TODO if child visible == gone , then skip
func (me *layer) measureChildren(width, height nux.MeasureDimen, innerWidth, innerHeight float32, childrenMeasuredFlags []uint8) (hPxMax, vPxMax float32, measuredCompleteCount int) {
	var hPx float32
	var hMWt float32 // horizontal margin weight
	var hMPt float32 // horizontal margin percent

	var vPx float32
	var vMWt float32 // vertical margin weight
	var vMPt float32 // vertical margin percent

	var measuredFlags uint8

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		hPx = 0
		vPx = 0
		hMWt = 0
		hMPt = 0
		vMWt = 0
		vMPt = 0

		measuredFlags = childrenMeasuredFlags[index]

		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			// clear last measured size
			cf.Clear()
			cf.Width = int32(nux.MeasureSpec(0, nux.Unlimit))
			cf.Height = int32(nux.MeasureSpec(0, nux.Unlimit))

			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
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
							t := cm.Top.Value() / 100.0 * innerHeight
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Auto, nux.Unlimit:
							vMPt += cm.Top.Value()
						}
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					switch cm.Bottom.Mode() {
					case nux.Pixel:
						b := cm.Bottom.Value()
						cf.Margin.Bottom = util.Roundi32(b)
						vPx += b
						measuredFlags |= flagMeasuredMarginBottom
						// ok
					case nux.Weight:
						vMWt += cm.Bottom.Value()
					case nux.Percent:
						switch height.Mode() {
						case nux.Pixel:
							b := cm.Bottom.Value() / 100.0 * innerHeight
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Auto, nux.Unlimit:
							vMPt += cm.Bottom.Value()
						}
					}
				}

				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
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
							l := cm.Left.Value() / 100.0 * innerWidth
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Auto, nux.Unlimit:
							hMPt += cm.Left.Value()
						}
					}
				}

				if measuredFlags&flagMeasuredMarginRight == flagMeasuredMarginRight {
					switch cm.Right.Mode() {
					case nux.Pixel:
						r := cm.Right.Value()
						cf.Margin.Right = int32(r)
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
							hMPt += cm.Right.Value()
						}

					}
				}
			} else {
				measuredFlags |= flagMeasuredMarginComplete
			}

			hPxRemain := innerWidth - hPx
			vPxRemain := innerHeight - vPx
			if hPxRemain < 0 {
				hPxRemain = 0
			}
			if vPxRemain < 0 {
				vPxRemain = 0
			}

			if measuredFlags&flagMeasuredWidth != flagMeasuredWidth {
				switch cs.Width().Mode() {
				case nux.Pixel:
					w := cs.Width().Value()
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), nux.Pixel))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				case nux.Weight:
					w := cs.Width().Value() / (hMWt + cs.Width().Value()) * hPxRemain
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				case nux.Percent:
					w := cs.Width().Value() / 100.0 * innerWidth
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
					setRatioHeightIfNeed(cs, cf, w, nux.Pixel)
				case nux.Ratio:
					if cs.Height().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both 'Ratio'")
					}
				case nux.Auto, nux.Unlimit:
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(hPxRemain), cs.Width().Mode()))
					setRatioHeightIfNeed(cs, cf, hPxRemain, nux.Pixel)
				}
			}

			if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
				switch cs.Height().Mode() {
				case nux.Pixel:
					h := cs.Height().Value()
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), nux.Pixel))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				case nux.Weight:
					h := cs.Height().Value() / (vMWt + cs.Height().Value()) * vPxRemain
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), height.Mode()))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				case nux.Percent:
					h := cs.Height().Value() / 100.0 * innerHeight
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), height.Mode()))
					setRatioWidthIfNeed(cs, cf, h, nux.Pixel)
				case nux.Ratio:
					if cs.Width().Mode() == nux.Ratio {
						log.Fatal("nux", "width and height size mode can not both nux.Ratio, at least one is definited.")
					}
				case nux.Auto, nux.Unlimit:
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(vPxRemain), cs.Height().Mode()))
					setRatioWidthIfNeed(cs, cf, vPxRemain, nux.Pixel)
				}
			}

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
							// TODO:: if debug
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

					hPx += float32(cf.Width)
					vPx += float32(cf.Height)
					hPxRemain -= float32(cf.Width)
					vPxRemain -= float32(cf.Height)
					if hPxRemain < 0 {
						hPxRemain = 0
					}
					if vPxRemain < 0 {
						vPxRemain = 0
					}
				} else {
					measuredFlags |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
					cf.Width = 0
					cf.Height = 0
				}

			}

			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					switch cm.Top.Mode() {
					// case nux.Pixel: has measured
					case nux.Weight:
						switch height.Mode() {
						case nux.Pixel:
							t := cm.Top.Value() / vMWt * vPxRemain
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					case nux.Percent:
						switch height.Mode() {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					switch cm.Bottom.Mode() {
					case nux.Weight:
						switch height.Mode() {
						case nux.Pixel:
							b := cm.Bottom.Value() / vMWt * vPxRemain
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					case nux.Percent:
						switch height.Mode() {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							// wait max height
						}
					}
				}

				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					switch cm.Left.Mode() {
					case nux.Weight:
						switch width.Mode() {
						case nux.Pixel:
							l := cm.Left.Value() / hMWt * hPxRemain
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					case nux.Percent:
						switch width.Mode() {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					}
				}

				if measuredFlags&flagMeasuredMarginRight == flagMeasuredMarginRight {
					switch cm.Right.Mode() {
					case nux.Weight:
						switch width.Mode() {
						case nux.Pixel:
							r := cm.Right.Value() / hMWt * hPxRemain
							cf.Margin.Right = util.Roundi32(r)
							hPx += r
							measuredFlags |= flagMeasuredMarginRight
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					case nux.Percent:
						switch width.Mode() {
						// case nux.Pixel: has measured
						case nux.Auto, nux.Unlimit:
							// wait max width
						}
					}
				}
			}

			hPx = hPx / (1.0 - hMPt/100.0)
			vPx = vPx / (1.0 - vMPt/100.0)

			if hPx > hPxMax {
				hPxMax = hPx
			}

			if vPx > vPxMax {
				vPxMax = vPx
			}

		} else {
			measuredFlags = flagMeasuredComplete
		}

		childrenMeasuredFlags[index] = measuredFlags

		if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
			measuredCompleteCount++
		}
	}

	return
}

func (me *layer) measureChildrenMargin(width, height nux.MeasureDimen, innerWidth, innerHeight float32, childrenMeasuredFlags []uint8) (measuredCompleteCount int) {
	var hPxRemain float32
	var hMWt float32 // horizontal margin weight

	var vPxRemain float32
	var vMWt float32 // vertical margin weight

	var measuredFlags uint8

	for index, child := range me.Children() {
		if compt, ok := child.(nux.Component); ok {
			child = compt.Content()
		}

		hPxRemain = 0
		vPxRemain = 0
		hMWt = 0
		vMWt = 0

		measuredFlags = childrenMeasuredFlags[index]

		if cs, ok := child.(nux.Size); ok {
			cf := cs.Frame()

			if cm := cs.Margin(); cm != nil {
				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					switch cm.Top.Mode() {
					case nux.Weight:
						vMWt += cm.Top.Value()
					case nux.Percent:
						t := cm.Top.Value() / 100.0 * innerHeight
						cf.Margin.Top = util.Roundi32(t)
						measuredFlags |= flagMeasuredMarginTop
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					switch cm.Bottom.Mode() {
					case nux.Weight:
						vMWt += cm.Bottom.Value()
					case nux.Percent:
						b := cm.Bottom.Value() / 100.0 * innerHeight
						cf.Margin.Bottom = util.Roundi32(b)
						measuredFlags |= flagMeasuredMarginBottom
					}
				}

				vPxRemain = (innerHeight - float32(cf.Height+cf.Margin.Top+cf.Margin.Bottom))

				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					switch cm.Top.Mode() {
					case nux.Weight:
						if vPxRemain > 0 && vMWt > 0 {
							t := cm.Top.Value() / vMWt * vPxRemain
							cf.Margin.Top = util.Roundi32(t)
						}
						measuredFlags |= flagMeasuredMarginTop
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					switch cm.Bottom.Mode() {
					case nux.Weight:
						if vPxRemain > 0 && vMWt > 0 {
							b := cm.Bottom.Value() / vMWt * vPxRemain
							cf.Margin.Bottom = util.Roundi32(b)
						}
						measuredFlags |= flagMeasuredMarginBottom
					}
				}

				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					switch cm.Left.Mode() {
					case nux.Weight:
						hMWt += cm.Left.Value()
					case nux.Percent:
						l := cm.Left.Value() / 100.0 * innerWidth
						cf.Margin.Left = util.Roundi32(l)
						measuredFlags |= flagMeasuredMarginLeft
					}
				}

				if measuredFlags&flagMeasuredMarginRight == flagMeasuredMarginRight {
					switch cm.Right.Mode() {
					case nux.Weight:
						hMWt += cm.Right.Value()
					case nux.Percent:
						r := cm.Right.Value() / 100 * innerWidth
						cf.Margin.Right = util.Roundi32(r)
						measuredFlags |= flagMeasuredMarginRight
					}
				}

				hPxRemain = (innerWidth - float32(cf.Width+cf.Margin.Left+cf.Margin.Right))

				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					switch cm.Left.Mode() {
					case nux.Weight:
						if hPxRemain > 0 && hMWt > 0 {
							l := cm.Left.Value() / hMWt * hPxRemain
							cf.Margin.Left = util.Roundi32(l)
						}
						measuredFlags |= flagMeasuredMarginLeft
					}
				}

				if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
					switch cm.Right.Mode() {
					case nux.Weight:
						if hPxRemain > 0 && hMWt > 0 {
							r := cm.Right.Value() / hMWt * hPxRemain
							cf.Margin.Right = util.Roundi32(r)
						}
						measuredFlags |= flagMeasuredMarginRight
					}
				}

			} else {
				measuredFlags |= flagMeasuredMarginComplete
			}

		} else {
			measuredFlags = flagMeasuredComplete
		}

		childrenMeasuredFlags[index] = measuredFlags

		if measuredFlags&flagMeasuredComplete == flagMeasuredComplete {
			measuredCompleteCount++
		}
	}

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
