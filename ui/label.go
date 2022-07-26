// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Label can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Label does not do this and uses the new AlignedText

import (
	"math"
	"time"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	"nuxui.org/nuxui/util"
)

type Label interface {
	nux.Widget
	nux.Size
	// nux.Stateable
	Visual

	Text() string
	SetText(text string)
}

type label struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	font               nux.Font
	fontLayout         nux.FontLayout
	text               string
	textColor          nux.Color
	textHighlightColor nux.Color
	paint              nux.Paint
	ellipsize          int
	iconLeft           nux.Widget
	iconTop            nux.Widget
	iconRight          nux.Widget
	iconBottom         nux.Widget
	align              *Align

	downTime    time.Time
	textOffsetX int32
	textOffsetY int32
	textWidth   int32
	textHeight  int32

	// state uint32
}

func NewLabel(attr nux.Attr) Label {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &label{
		text:               attr.GetString("text", ""),
		textColor:          attr.GetColor("textColor", nux.Black),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		align:              NewAlign(attr.GetAttr("align", nux.Attr{"horizontal": "center", "vertical": "center"})),
		font:               nux.NewFont(attr.GetAttr("font", nil)),
		fontLayout:         nux.NewFontLayout(),
		paint:              nux.NewPaint(),
	}

	if icon := attr.GetAttr("icon", nil); icon != nil {
		if iconLeft := icon.GetAttr("left", nil); iconLeft != nil {
			me.iconLeft = nux.InflateLayoutAttr(nil, iconLeft, nil)
		}
		if iconTop := icon.GetAttr("top", nil); iconTop != nil {
			me.iconTop = nux.InflateLayoutAttr(nil, iconTop, nil)
		}
		if iconRight := icon.GetAttr("right", nil); iconRight != nil {
			me.iconRight = nux.InflateLayoutAttr(nil, iconRight, nil)
		}
		if iconBottom := icon.GetAttr("bottom", nil); iconBottom != nil {
			me.iconBottom = nux.InflateLayoutAttr(nil, iconBottom, nil)
		}
	}

	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

func (me *label) OnMount() {
	nux.OnTapDown(me, me.onTapDown)
	nux.OnTapUp(me, me.onTapUp)
	nux.OnTapCancel(me, me.onTapUp)
	nux.OnTap(me, me.onTap)
	nux.OnHoverEnter(me, me.OnHoverEnter)
	nux.OnHoverExit(me, me.OnHoverExit)
}

func (me *label) OnEject() {
}

func (me *label) OnHoverEnter(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().AddState(nux.State_Hovered)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().AddState(nux.State_Hovered)
			changed = true
		}
		changed = changed || me.addIconState(me.iconLeft) || me.addIconState(me.iconTop) || me.addIconState(me.iconRight) || me.addIconState(me.iconBottom)
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) OnHoverExit(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().DelState(nux.State_Hovered)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().DelState(nux.State_Hovered)
			changed = true
		}
		changed = changed || me.delIconState(me.iconLeft) || me.delIconState(me.iconTop) || me.delIconState(me.iconRight) || me.delIconState(me.iconBottom)
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) onTapDown(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().AddState(nux.State_Pressed)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().AddState(nux.State_Pressed)
			changed = true
		}
		changed = changed || me.addIconState(me.iconLeft) || me.addIconState(me.iconTop) || me.addIconState(me.iconRight) || me.addIconState(me.iconBottom)
	}

	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) onTapUp(detail nux.GestureDetail) {
	changed := false
	if !me.Disable() {
		if me.Background() != nil {
			me.Background().DelState(nux.State_Pressed)
			changed = true
		}
		if me.Foreground() != nil {
			me.Foreground().DelState(nux.State_Pressed)
			changed = true
		}
		changed = changed || me.delIconState(me.iconLeft) || me.delIconState(me.iconTop) || me.delIconState(me.iconRight) || me.delIconState(me.iconBottom)
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) onTap(detail nux.GestureDetail) {
}

func (me *label) addIconState(icon nux.Widget) (changed bool) {
	if icon != nil {
		if s, ok := icon.(nux.Stateable); ok {
			s.AddState(nux.State_Pressed)
			changed = true
		}
	}
	return
}

func (me *label) delIconState(icon nux.Widget) (changed bool) {
	if icon != nil {
		if s, ok := icon.(nux.Stateable); ok {
			s.DelState(nux.State_Pressed)
			changed = true
		}
	}
	return
}

// func (me *label) AddState(state uint32) {
// 	s := me.state
// 	s |= state
// 	me.state = s
// 	// me.applyState()
// }

// func (me *label) DelState(state uint32) {
// 	s := me.state
// 	s &= ^state
// 	me.state = s
// 	// me.applyState()
// }

// func (me *label) State() uint32 {
// 	return me.state
// }

func (me *label) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *label) Text() string {
	return me.text
}

func (me *label) SetText(text string) {
	if me.text != text {
		me.text = text
		nux.RequestLayout(me)
	}
}

// TODO:: if only text
func (me *label) Measure(width, height nux.MeasureDimen) {
	frame := me.Frame()

	me.textWidth, me.textHeight = me.fontLayout.MeasureText(me.font, me.text, width.Value(), height.Value())
	txtW := float32(me.textWidth)
	txtH := float32(me.textHeight)

	hPPx, hPPt, vPPx, vPPt, paddingMeasuredFlag := measurePadding(width, height, me.Padding(), frame, float32(txtH), 0)
	if hPPt >= 100.0 || vPPt >= 100.0 {
		log.Fatal("nuxui", "padding percent size should at 0% ~ 100%")
	}

	innerWidth := float32(width.Value())*(1.0-hPPt/100.0) - hPPx
	innerHeight := float32(height.Value())*(1.0-vPPt/100.0) - vPPx

	hPxl, vPxl, hMWtl, vMWtl, hWtl, vWtl, hPtl, vPtl, ptWL, ptHL, hasL, fl := me.measureIconSize(me.iconLeft, width, height, innerWidth, innerHeight, txtH)
	hPxr, vPxr, hMWtr, vMWtr, hWtr, vWtr, hPtr, vPtr, ptWR, ptHR, hasR, fr := me.measureIconSize(me.iconRight, width, height, innerWidth, innerHeight, txtH)
	hPxt, vPxt, hMWtt, vMWtt, hWtt, vWtt, hPtt, vPtt, ptWT, ptHT, hasT, ft := me.measureIconSize(me.iconTop, width, height, innerWidth, innerHeight, txtH)
	hPxb, vPxb, hMWtb, vMWtb, hWtb, vWtb, hPtb, vPtb, ptWB, ptHB, hasB, fb := me.measureIconSize(me.iconBottom, width, height, innerWidth, innerHeight, txtH)

	if hasL {
		w, h := me.measureIconWeightSize(me.iconLeft, width, height, innerWidth-hPxl-hPxr-txtW, innerHeight-vPxl, hMWtl+hWtl+hMWtr+hWtr, vMWtl+vWtl, fl)
		if fl&flagMeasuredWidth != flagMeasuredWidth {
			hPxl += w
		}
		if fl&flagMeasuredHeight != flagMeasuredHeight {
			vPxl += h
		}
		fl |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
	}

	if hasR {
		w, h := me.measureIconWeightSize(me.iconRight, width, height, innerWidth-hPxl-hPxr-txtW, innerHeight-vPxr, hMWtl+hWtl+hMWtr+hWtr, vMWtr+vWtr, fr)
		if fr&flagMeasuredWidth != flagMeasuredWidth {
			hPxr += w
		}
		if fr&flagMeasuredHeight != flagMeasuredHeight {
			vPxr += h
		}
		fr |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
	}

	if hasT {
		w, h := me.measureIconWeightSize(me.iconTop, width, height, innerWidth-hPxt, innerHeight-vPxt-vPxb-txtH, hMWtt+hWtt, vMWtt+vWtt+vMWtb+vWtb, fl)
		if ft&flagMeasuredWidth != flagMeasuredWidth {
			hPxt += w
		}
		if ft&flagMeasuredHeight != flagMeasuredHeight {
			vPxt += h
		}
		ft |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
	}

	if hasB {
		w, h := me.measureIconWeightSize(me.iconBottom, width, height, innerWidth-hPxb, innerHeight-vPxt-vPxb-txtH, hMWtb+hWtb, vMWtt+vWtt+vMWtb+vWtb, fr)
		if fb&flagMeasuredWidth != flagMeasuredWidth {
			hPxb += w
		}
		if fb&flagMeasuredHeight != flagMeasuredHeight {
			vPxb += h
		}
		fb |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
	}

	checkPercents(hPtl, vPtl, hPtt, vPtt, hPtr, vPtr, hPtb, vPtb)

	// find max inner width or height
	hPxMax := util.MaxF(hPxl+hPxr+txtW, hPxt, hPxb, ptWL, ptWT, ptWR, ptWB, (hPxt / (1 - hPtt/100)), ((hPxl + hPxr + txtW) / (1 - (hPtl-hPtr)/100)), (hPxb / (1 - hPtb/100)))
	vPxMax := util.MaxF(vPxt+vPxb+txtH, vPxl, vPxr, ptHL, ptHT, ptHR, ptHB, (vPxl / (1 - vPtl/100)), ((vPxt + vPxb + txtH) / (1 - (vPtt+vPtb)/100)), (vPxr / (1 - vPtr/100)))

	switch width.Mode() {
	case nux.Auto:
		innerWidth = util.MinF(innerWidth, hPxMax)
	case nux.Unlimit:
		innerWidth = util.MaxF(innerWidth, hPxMax)
	}

	switch height.Mode() {
	case nux.Auto:
		innerHeight = util.MinF(innerHeight, vPxMax)
	case nux.Unlimit:
		innerHeight = util.MaxF(innerHeight, vPxMax)
	}

	if fl&flagMeasuredComplete != flagMeasuredComplete {
		fl = me.measureIconMarginSize(me.iconLeft, width, height, innerWidth, innerHeight, innerWidth-hPxl-hPxr-txtW, innerHeight-vPxl, hMWtl+hMWtr, vMWtl, fl)
	}
	if fr&flagMeasuredComplete != flagMeasuredComplete {
		fr = me.measureIconMarginSize(me.iconRight, width, height, innerWidth, innerHeight, innerWidth-hPxl-hPxr-txtW, innerHeight-vPxr, hMWtl+hMWtr, vMWtr, fr)
	}
	if ft&flagMeasuredComplete != flagMeasuredComplete {
		ft = me.measureIconMarginSize(me.iconTop, width, height, innerWidth, innerHeight, innerWidth-hPxt, innerHeight-vPxt-txtH-vPxb, hMWtt, vMWtt+vMWtb, ft)
	}
	if fb&flagMeasuredComplete != flagMeasuredComplete {
		fb = me.measureIconMarginSize(me.iconBottom, width, height, innerWidth, innerHeight, innerWidth-hPxb, innerHeight-vPxt-txtH-vPxb, hMWtb, vMWtt+vMWtb, fb)
	}

	if fl&fr&ft&fb&flagMeasuredComplete != flagMeasuredComplete {
		log.Fatal("nuxui", "can not run here")
	}

	if width.Mode() == nux.Pixel {
		frame.Width = width.Value()
	} else {
		w := (innerWidth + hPPx) / (1.0 - hPPt/100.0)
		frame.Width = int32(math.Ceil(float64(w)))
	}

	if height.Mode() == nux.Pixel {
		frame.Height = height.Value()
	} else {
		h := (innerHeight + vPPx) / (1.0 - vPPt/100.0)
		frame.Height = int32(math.Ceil(float64(h)))
	}

	if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		measurePadding(width, height, me.Padding(), frame, txtH, paddingMeasuredFlag)
	}

	// log.I("nuxui", "ui.Label %s '%s' Measure end width=%d, height=%d, txtW", me.Info().ID, me.Text(), frame.Width, frame.Height, txtW)

}

func (me *label) measureIconSize(icon nux.Widget, width, height nux.MeasureDimen, innerWidth, innerHeight, txtH float32) (hPx, vPx, hMWt, vMWt, hWt, vWt, hPt, vPt, ptWidth, ptHeight float32, hasAutoWeightSize bool, measuredFlags uint8) {
	if icon != nil {
		if s, ok := icon.(nux.Size); ok {
			cf := s.Frame()
			if m := s.Margin(); m != nil {
				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					if m.Left.Value() != 0 {
						switch m.Left.Mode() {
						case nux.Pixel:
							l := m.Left.Value()
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Ems:
							l := txtH * m.Left.Value()
							cf.Margin.Left = util.Roundi32(l)
							hPx += l
							measuredFlags |= flagMeasuredMarginLeft
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								l := m.Left.Value() / 100 * innerWidth
								cf.Margin.Left = util.Roundi32(l)
								hPx += l
								measuredFlags |= flagMeasuredMarginLeft
							case nux.Auto:
								hPt += m.Left.Value()
							}
						case nux.Weight:
							hMWt += m.Left.Value()
						}
					} else {
						measuredFlags |= flagMeasuredMarginLeft
					}
				}

				if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
					if m.Right.Value() != 0 {
						switch m.Right.Mode() {
						case nux.Pixel:
							r := m.Right.Value()
							cf.Margin.Right = util.Roundi32(r)
							hPx += r
							measuredFlags |= flagMeasuredMarginRight
						case nux.Ems:
							r := txtH * m.Right.Value()
							cf.Margin.Right = util.Roundi32(r)
							hPx += r
							measuredFlags |= flagMeasuredMarginRight
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								r := m.Right.Value() / 100 * innerWidth
								cf.Margin.Right = util.Roundi32(r)
								hPx += r
								measuredFlags |= flagMeasuredMarginRight
							case nux.Auto:
								hPt += m.Right.Value()
							}
						case nux.Weight:
							hMWt += m.Right.Value()
						}
					} else {
						measuredFlags |= flagMeasuredMarginRight
					}
				}

				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					if m.Top.Value() != 0 {
						switch m.Top.Mode() {
						case nux.Pixel:
							t := m.Top.Value()
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Ems:
							t := txtH * m.Top.Value()
							cf.Margin.Top = util.Roundi32(t)
							vPx += t
							measuredFlags |= flagMeasuredMarginTop
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								t := m.Top.Value() / 100 * innerHeight
								cf.Margin.Top = util.Roundi32(t)
								vPx += t
								measuredFlags |= flagMeasuredMarginTop
							case nux.Auto:
								vPt += m.Top.Value()
							}
						case nux.Weight:
							vMWt += m.Top.Value()
						}
					} else {
						measuredFlags |= flagMeasuredMarginTop
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					if m.Bottom.Value() != 0 {
						switch m.Bottom.Mode() {
						case nux.Pixel:
							b := m.Bottom.Value()
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Ems:
							b := txtH * m.Bottom.Value()
							cf.Margin.Bottom = util.Roundi32(b)
							vPx += b
							measuredFlags |= flagMeasuredMarginBottom
						case nux.Percent:
							switch width.Mode() {
							case nux.Pixel:
								b := m.Bottom.Value() / 100 * innerHeight
								cf.Margin.Bottom = util.Roundi32(b)
								vPx += b
								measuredFlags |= flagMeasuredMarginBottom
							case nux.Auto:
								vPt += m.Bottom.Value()
							}
						case nux.Weight:
							vMWt += m.Bottom.Value()
						}
					} else {
						measuredFlags |= flagMeasuredMarginBottom
					}
				}
			} else {
				measuredFlags |= flagMeasuredMarginComplete
			}

			canMeasure := true
			if measuredFlags&flagMeasuredWidth != flagMeasuredWidth {
				switch s.Width().Mode() {
				case nux.Pixel:
					w := s.Width().Value()
					cf.Width = util.Roundi32(w)
					setRatioHeightIfNeed(s, s.Frame(), w, nux.Pixel)
					measuredFlags |= flagMeasuredWidth
				case nux.Ems:
					w := txtH * s.Width().Value()
					cf.Width = util.Roundi32(w)
					setRatioHeightIfNeed(s, s.Frame(), w, nux.Pixel)
					measuredFlags |= flagMeasuredWidth
				case nux.Percent:
					w := s.Width().Value() / 100 * innerWidth
					cf.Width = util.Roundi32(w)
					setRatioHeightIfNeed(s, s.Frame(), w, width.Mode())
					if width.Mode() == nux.Pixel {
						measuredFlags |= flagMeasuredWidth
					}
				case nux.Weight:
					hWt += s.Width().Value()
					canMeasure = false
					if width.Mode() != nux.Pixel {
						hasAutoWeightSize = true
					}
				case nux.Auto:
					// measure later
					canMeasure = false
				}
			}

			if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
				switch s.Height().Mode() {
				case nux.Pixel:
					h := s.Height().Value()
					cf.Height = util.Roundi32(h)
					setRatioWidthIfNeed(s, s.Frame(), h, nux.Pixel)
					measuredFlags |= flagMeasuredHeight
				case nux.Ems:
					h := txtH * s.Height().Value()
					cf.Height = util.Roundi32(h)
					setRatioWidthIfNeed(s, s.Frame(), h, nux.Pixel)
					measuredFlags |= flagMeasuredHeight
				case nux.Percent:
					h := s.Height().Value() / 100 * innerHeight
					cf.Height = util.Roundi32(h)
					setRatioWidthIfNeed(s, s.Frame(), h, height.Mode())
					if height.Mode() == nux.Pixel {
						measuredFlags |= flagMeasuredHeight
					}
				case nux.Weight:
					vWt += s.Height().Value()
					canMeasure = false
					if height.Mode() != nux.Pixel {
						hasAutoWeightSize = true
					}
				case nux.Auto:
					// measure later
					canMeasure = false
				}
			}

			if canMeasure && measuredFlags&flagMeasured != flagMeasured {
				if m, ok := icon.(nux.Measure); ok {
					m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))

					if nux.MeasureDimen(cf.Width).Mode() != nux.Pixel ||
						nux.MeasureDimen(cf.Height).Mode() != nux.Pixel {
						log.Fatal("nuxui", "label %s the child %s(%T) measured not completed", me.Info().ID, icon.Info().ID, icon)
					}

					if s.Width().Mode() == nux.Ratio {
						oldwidth := cf.Width
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Height)*s.Width().Value()), nux.Pixel))
						if oldwidth != cf.Width {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					if s.Height().Mode() == nux.Ratio {
						oldheight := cf.Height
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Width)/s.Height().Value()), nux.Pixel))
						if oldheight != cf.Height {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					if width.Mode() != nux.Pixel && s.Width().Mode() == nux.Percent {
						ptWidth = float32(cf.Width) / (s.Width().Value() / 100.0)
					}
					if height.Mode() != nux.Pixel && s.Height().Mode() == nux.Percent {
						ptHeight = float32(cf.Height) / (s.Height().Value() / 100.0)
					}

					measuredFlags |= flagMeasured | flagMeasuredWidth | flagMeasuredHeight
				}
			}

			if measuredFlags&flagMeasuredWidth == flagMeasuredWidth {
				hPx += float32(cf.Width)
			}
			if measuredFlags&flagMeasuredHeight == flagMeasuredHeight {
				vPx += float32(cf.Height)
			}
		} else {
			measuredFlags |= flagMeasuredComplete
		}
	} else {
		measuredFlags |= flagMeasuredComplete
	}
	return
}

func (me *label) measureIconWeightSize(icon nux.Widget, width, height nux.MeasureDimen, hPxRemain, vPxRemain, hWt, vWt float32, measuredFlags uint8) (w, h float32) {
	if icon != nil {
		if s, ok := icon.(nux.Size); ok {
			cf := s.Frame()

			if measuredFlags&flagMeasuredWidth != flagMeasuredWidth {
				switch s.Width().Mode() {
				case nux.Weight:
					w := s.Width().Value() / hWt * hPxRemain
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(w), width.Mode()))
					setRatioHeightIfNeed(s, s.Frame(), w, nux.Pixel)
				case nux.Auto:
					cf.Width = int32(nux.MeasureSpec(util.Roundi32(hPxRemain), width.Mode()))
					setRatioHeightIfNeed(s, s.Frame(), hPxRemain, nux.Pixel)
				}
			}

			if measuredFlags&flagMeasuredHeight != flagMeasuredHeight {
				switch s.Height().Mode() {
				case nux.Weight:
					h := s.Height().Value() / vWt * vPxRemain
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(h), width.Mode()))
					setRatioWidthIfNeed(s, s.Frame(), h, nux.Pixel)
				case nux.Auto:
					cf.Height = int32(nux.MeasureSpec(util.Roundi32(hPxRemain), height.Mode()))
					setRatioWidthIfNeed(s, s.Frame(), hPxRemain, nux.Pixel)
				}
			}

			if measuredFlags&flagMeasured != flagMeasured {
				if m, ok := icon.(nux.Measure); ok {
					m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))

					if nux.MeasureDimen(cf.Width).Mode() != nux.Pixel ||
						nux.MeasureDimen(cf.Height).Mode() != nux.Pixel {
						log.Fatal("nuxui", "label %s the child %s(%T) measured not completed", me.Info().ID, icon.Info().ID, icon)
					}

					if s.Width().Mode() == nux.Ratio {
						oldwidth := cf.Width
						cf.Width = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Height)*s.Width().Value()), nux.Pixel))
						if oldwidth != cf.Width {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					if s.Height().Mode() == nux.Ratio {
						oldheight := cf.Height
						cf.Height = int32(nux.MeasureSpec(util.Roundi32(float32(cf.Width)/s.Height().Value()), nux.Pixel))
						if oldheight != cf.Height {
							log.W("nuxui", "Size auto and ratio make measure twice, please optimization layout")
							m.Measure(nux.MeasureDimen(cf.Width), nux.MeasureDimen(cf.Height))
						}
					}

					w = float32(cf.Width)
					h = float32(cf.Height)
				}
			}
		}
	}
	return
}

func (me *label) measureIconMarginSize(icon nux.Widget, width, height nux.MeasureDimen, innerWidth, innerHeight, hPxRemain, vPxRemain, hMWt, vMWt float32, measuredFlags uint8) (measuredFlagsOut uint8) {
	if icon != nil {
		if s, ok := icon.(nux.Size); ok {

			if m := s.Margin(); m != nil {
				if measuredFlags&flagMeasuredMarginLeft != flagMeasuredMarginLeft {
					if m.Left.Value() != 0 {
						switch m.Left.Mode() {
						case nux.Percent:
							s.Frame().Margin.Left = util.Roundi32(m.Left.Value() / innerWidth)
						case nux.Weight:
							s.Frame().Margin.Left = util.Roundi32(hPxRemain * m.Left.Value() / hMWt)
						}
					}
				}

				if measuredFlags&flagMeasuredMarginRight != flagMeasuredMarginRight {
					if m.Right.Value() != 0 {
						switch m.Right.Mode() {
						case nux.Percent:
							s.Frame().Margin.Right = util.Roundi32(m.Right.Value() / innerWidth)
						case nux.Weight:
							s.Frame().Margin.Right = util.Roundi32(hPxRemain * m.Right.Value() / hMWt)
						}
					}
				}

				if measuredFlags&flagMeasuredMarginTop != flagMeasuredMarginTop {
					if m.Top.Value() != 0 {
						switch m.Top.Mode() {
						case nux.Percent:
							s.Frame().Margin.Top = util.Roundi32(m.Top.Value() / innerHeight)
						case nux.Weight:
							s.Frame().Margin.Top = util.Roundi32(vPxRemain * m.Top.Value() / vMWt)
						}
					}
				}

				if measuredFlags&flagMeasuredMarginBottom != flagMeasuredMarginBottom {
					if m.Bottom.Value() != 0 {
						switch m.Bottom.Mode() {
						case nux.Percent:
							s.Frame().Margin.Bottom = util.Roundi32(m.Bottom.Value() / innerHeight)
						case nux.Weight:
							s.Frame().Margin.Bottom = util.Roundi32(vPxRemain * m.Bottom.Value() / vMWt)
						}
					}
				}
			}
		}
	}
	return measuredFlags
}

// Responsible for determining the position of the widget align, margin...
func (me *label) Layout(x, y, width, height int32) {
	// log.I("nuxui", "label layout x=%d ,y=%d, width=%d, height=%d ", x, y, width, height)
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	innerWidth := float32(width - frame.Padding.Left - frame.Padding.Right)
	innerHeight := float32(height - frame.Padding.Top - frame.Padding.Bottom)

	me.textOffsetX = frame.Padding.Left
	me.textOffsetY = frame.Padding.Top

	var w1, w2, w3, h1, h2, h3 int32
	var sLeft, sTop, sRight, sBottom nux.Size
	var sok bool

	if me.iconLeft != nil {
		if sLeft, sok = me.iconLeft.(nux.Size); sok {
			f := sLeft.Frame()
			w2 += f.Margin.Left + f.Width + f.Margin.Right
			h1 += f.Margin.Top + f.Height + f.Margin.Bottom
		}
	}

	w2 += me.textWidth

	if me.iconTop != nil {
		if sTop, sok = me.iconTop.(nux.Size); sok {
			f := sTop.Frame()
			w1 += f.Margin.Left + f.Width + f.Margin.Right
			h2 += f.Margin.Top + f.Height + f.Margin.Bottom
		}
	}

	h2 += me.textHeight

	if me.iconRight != nil {
		if sRight, sok = me.iconRight.(nux.Size); sok {
			f := sRight.Frame()
			w2 += f.Margin.Left + f.Width + f.Margin.Right
			h3 += f.Margin.Top + f.Height + f.Margin.Bottom
		}
	}

	if me.iconBottom != nil {
		if sBottom, sok = me.iconBottom.(nux.Size); sok {
			f := sBottom.Frame()
			w3 += f.Margin.Left + f.Width + f.Margin.Right
			h2 += f.Margin.Top + f.Height + f.Margin.Bottom
		}
	}

	// log.I("nuxui", "label layout w2=%d, h2=%d", w2, h2)

	switch me.align.Horizontal {
	case Left:
		me.textOffsetX = frame.Padding.Left
		if sLeft != nil {
			sLeft.Frame().X = frame.X + frame.Padding.Left + sLeft.Frame().Margin.Left
			me.textOffsetX += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
		}
		if sTop != nil {
			sTop.Frame().X = frame.X + frame.Padding.Left + sTop.Frame().Margin.Left
		}
		if sRight != nil {
			sRight.Frame().X = frame.X + frame.Padding.Left + +me.textWidth + sRight.Frame().Margin.Left
			if sLeft != nil {
				sRight.Frame().X += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
			}
		}
		if sBottom != nil {
			sBottom.Frame().X = frame.X + frame.Padding.Left + sBottom.Frame().Margin.Left
		}
	case Center:
		me.textOffsetX = frame.Padding.Left + util.Roundi32((innerWidth-float32(w2))/2)
		if sLeft != nil {
			sLeft.Frame().X = frame.X + frame.Padding.Left + util.Roundi32((innerWidth-float32(w2))/2) + sLeft.Frame().Margin.Left
			me.textOffsetX += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
		}
		if sTop != nil {
			sTop.Frame().X = frame.X + frame.Padding.Left + util.Roundi32((innerWidth-float32(w1))/2) + sTop.Frame().Margin.Left
		}
		if sRight != nil {
			sRight.Frame().X = frame.X + frame.Padding.Left + +util.Roundi32((innerWidth-float32(w2))/2) + me.textWidth + sRight.Frame().Margin.Left
			if sLeft != nil {
				sRight.Frame().X += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
			}
		}
		if sBottom != nil {
			sBottom.Frame().X = frame.X + frame.Padding.Left + util.Roundi32((innerWidth-float32(w3))/2) + sBottom.Frame().Margin.Left
		}
	case Right:
		me.textOffsetX = frame.Padding.Left + util.Roundi32(innerWidth-float32(w2))
		if sLeft != nil {
			sLeft.Frame().X = frame.X + frame.Padding.Left + util.Roundi32(innerWidth-float32(w2)) + sLeft.Frame().Margin.Left
			me.textOffsetX += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
		}
		if sTop != nil {
			sTop.Frame().X = frame.X + frame.Padding.Left + util.Roundi32(innerWidth-float32(w1)) + sTop.Frame().Margin.Left
		}
		if sRight != nil {
			sRight.Frame().X = frame.X + frame.Padding.Left + util.Roundi32(innerWidth-float32(w2)) + me.textWidth + sRight.Frame().Margin.Left
			if sLeft != nil {
				sRight.Frame().X += sLeft.Frame().Margin.Left + sLeft.Frame().Width + sLeft.Frame().Margin.Right
			}
		}
		if sBottom != nil {
			sBottom.Frame().X = frame.X + frame.Padding.Left + util.Roundi32(innerWidth-float32(w3))
		}
	}

	switch me.align.Vertical {
	case Top:
		me.textOffsetY = frame.Padding.Top
		if sLeft != nil {
			sLeft.Frame().Y = frame.Y + frame.Padding.Top + sLeft.Frame().Margin.Top
		}
		if sTop != nil {
			sTop.Frame().Y = frame.Y + frame.Padding.Top + sTop.Frame().Margin.Top
			me.textOffsetY += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
		}
		if sRight != nil {
			sRight.Frame().Y = frame.Y + frame.Padding.Top + sRight.Frame().Margin.Top
		}
		if sBottom != nil {
			sBottom.Frame().Y = frame.Y + frame.Padding.Top + me.textHeight + sBottom.Frame().Margin.Top
			if sTop != nil {
				sBottom.Frame().Y += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
			}
		}
	case Center:
		me.textOffsetY = frame.Padding.Top + util.Roundi32((innerHeight-float32(h2))/2)
		if sLeft != nil {
			sLeft.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32((innerHeight-float32(h1))/2) + sLeft.Frame().Margin.Top
		}
		if sTop != nil {
			sTop.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32((innerHeight-float32(h2))/2) + sTop.Frame().Margin.Top
			me.textOffsetY += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
		}
		if sRight != nil {
			sRight.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32((innerHeight-float32(h3))/2) + sRight.Frame().Margin.Top
		}
		if sBottom != nil {
			sBottom.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32((innerHeight-float32(h2))/2) + me.textHeight + sBottom.Frame().Margin.Top
			if sTop != nil {
				sBottom.Frame().Y += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
			}
		}
	case Bottom:
		me.textOffsetY = frame.Padding.Top + util.Roundi32(innerHeight-float32(h2))
		if sLeft != nil {
			sLeft.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32(innerHeight-float32(h1)) + sLeft.Frame().Margin.Top
		}
		if sTop != nil {
			sTop.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32(innerHeight-float32(h2)) + sTop.Frame().Margin.Top
			me.textOffsetY += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
		}
		if sRight != nil {
			sRight.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32(innerHeight-float32(h3)) + sRight.Frame().Margin.Top
		}
		if sBottom != nil {
			sBottom.Frame().Y = frame.Y + frame.Padding.Top + util.Roundi32(innerHeight-float32(h3)) + me.textHeight + sBottom.Frame().Margin.Top
			if sTop != nil {
				sBottom.Frame().Y += sTop.Frame().Margin.Top + sTop.Frame().Width + sTop.Frame().Margin.Bottom
			}
		}
	}

	if me.iconLeft != nil {
		if layout, ok := me.iconLeft.(nux.Layout); ok {
			layout.Layout(sLeft.Frame().X, sLeft.Frame().Y, sLeft.Frame().Width, sLeft.Frame().Height)
		}
	}
	if me.iconTop != nil {
		if layout, ok := me.iconTop.(nux.Layout); ok {
			layout.Layout(sTop.Frame().X, sTop.Frame().Y, sTop.Frame().Width, sTop.Frame().Height)
		}
	}
	if me.iconRight != nil {
		if layout, ok := me.iconRight.(nux.Layout); ok {
			layout.Layout(sRight.Frame().X, sRight.Frame().Y, sRight.Frame().Width, sRight.Frame().Height)
		}
	}
	if me.iconBottom != nil {
		if layout, ok := me.iconBottom.(nux.Layout); ok {
			layout.Layout(sBottom.Frame().X, sBottom.Frame().Y, sBottom.Frame().Width, sBottom.Frame().Height)
		}
	}
}

func (me *label) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	if me.iconLeft != nil {
		if draw, ok := me.iconLeft.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if me.iconRight != nil {
		if draw, ok := me.iconRight.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if me.iconTop != nil {
		if draw, ok := me.iconTop.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	if me.iconBottom != nil {
		if draw, ok := me.iconBottom.(nux.Draw); ok {
			draw.Draw(canvas)
		}
	}

	// draw text
	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X+me.textOffsetX), float32(frame.Y+me.textOffsetY))

	if me.text != "" {
		me.paint.SetColor(me.textColor)
		me.fontLayout.DrawText(canvas, me.font, me.paint, me.text, frame.Width, frame.Height)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}

func checkPercents(pts ...float32) {
	for _, v := range pts {
		if v < 0 || v >= 100 {
			log.Fatal("nuxui", "percent value shouold between 0~100")
		}
	}
}
