// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Label can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Label does not do this and uses the new AlignedText

import (
	"math"
	"time"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Label interface {
	nux.Widget
	nux.Size
	nux.Stateable
	Visual

	Text() string
	SetText(text string)
}

type label struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	text               string
	textSize           float32
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

	state uint32
}

func NewLabel(attr nux.Attr) Label {
	me := &label{
		text:               attr.GetString("text", ""),
		textSize:           attr.GetFloat32("textSize", 12),
		textColor:          attr.GetColor("textColor", nux.White),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		align:              NewAlign(attr.GetAttr("align", nux.Attr{"horizontal": "center", "vertical": "center"})),
		paint:              nux.NewPaint(),
	}

	if icon := attr.GetAttr("icon", nil); icon != nil {
		if iconLeft := icon.GetAttr("left", nil); iconLeft != nil {
			me.iconLeft = nux.InflateLayoutAttr(nil, iconLeft)
		}
		if iconTop := icon.GetAttr("top", nil); iconTop != nil {
			me.iconTop = nux.InflateLayoutAttr(nil, iconTop)
		}
		if iconRight := icon.GetAttr("right", nil); iconRight != nil {
			me.iconRight = nux.InflateLayoutAttr(nil, iconRight)
		}
		if iconBottom := icon.GetAttr("bottom", nil); iconBottom != nil {
			me.iconBottom = nux.InflateLayoutAttr(nil, iconBottom)
		}
	}

	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)

	return me
}

func (me *label) Mount() {
	log.I("nuxui", "label Mount")
	nux.OnTapDown(me.Info().Self, me.onTapDown)
	nux.OnTapUp(me.Info().Self, me.onTapUp)
	nux.OnTapCancel(me.Info().Self, me.onTapUp)
	nux.OnTap(me.Info().Self, me.onTap)
}

func (me *label) Eject() {
	log.I("nuxui", "label Eject")
}

func (me *label) onTapDown(detail nux.GestureDetail) {
	log.V("nuxui", "label onTapDown")
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
		if me.iconLeft != nil {
			if s, ok := me.iconLeft.(nux.Stateable); ok {
				s.AddState(nux.State_Pressed)
				changed = true
			}
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) onTapUp(detail nux.GestureDetail) {
	log.V("nuxui", "label onTapUp")
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
		if me.iconLeft != nil {
			if s, ok := me.iconLeft.(nux.Stateable); ok {
				s.DelState(nux.State_Pressed)
				changed = true
			}
		}
	}
	if changed {
		nux.RequestRedraw(me)
	}
}

func (me *label) onTap(detail nux.GestureDetail) {
	log.V("nuxui", "label onTap")
}

func (me *label) AddState(state uint32) {
	s := me.state
	s |= state
	me.state = s
	// me.applyState()
}

func (me *label) DelState(state uint32) {
	s := me.state
	s &= ^state
	me.state = s
	// me.applyState()
}

func (me *label) State() uint32 {
	return me.state
}

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

func (me *label) Measure(width, height int32) {
	var vPPt float32 // horizontal padding percent
	var vPPx float32 // horizontal padding pixel
	var hPPt float32 //
	var hPPx float32 //
	// var vPx float32  // vertical total pixel
	// var hPx float32

	frame := me.Frame()

	me.paint.SetTextSize(me.textSize)
	txtW, txtH := me.paint.MeasureText(me.text, float32(nux.MeasureSpecValue(width)), float32(nux.MeasureSpecValue(height)))
	me.textWidth = int32(math.Ceil(float64(txtW)))
	me.textHeight = int32(math.Ceil(float64(txtH)))

	// 1. Calculate its own padding size
	if p := me.Padding(); p != nil {
		if p.Left.Value() != 0 {
			switch p.Left.Mode() {
			case nux.Pixel:
				l := p.Left.Value()
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
			case nux.Ems:
				l := txtH * p.Left.Value()
				frame.Padding.Left = util.Roundi32(l)
				hPPx += l
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					l := p.Left.Value() / 100 * float32(nux.MeasureSpecValue(width))
					frame.Padding.Left = util.Roundi32(l)
					hPPx += l
				case nux.Auto:
					hPPt += p.Left.Value()
				}
			}
		}

		if p.Right.Value() != 0 {
			switch p.Right.Mode() {
			case nux.Pixel:
				r := p.Right.Value()
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
			case nux.Ems:
				r := txtH * p.Right.Value()
				frame.Padding.Right = util.Roundi32(r)
				hPPx += r
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					r := p.Right.Value() / 100 * float32(nux.MeasureSpecValue(width))
					frame.Padding.Right = util.Roundi32(r)
					hPPx += r
				case nux.Auto:
					hPPt += p.Right.Value()
				}
			}
		}

		if p.Top.Value() != 0 {
			switch p.Top.Mode() {
			case nux.Pixel:
				t := p.Top.Value()
				frame.Padding.Top = util.Roundi32(t)
				vPPx += t
			case nux.Ems:
				t := txtH * p.Top.Value()
				frame.Padding.Top = util.Roundi32(t)
				hPPx += t
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					t := p.Top.Value() / 100 * float32(nux.MeasureSpecValue(height))
					frame.Padding.Top = util.Roundi32(t)
					vPPx += t
				case nux.Auto:
					vPPt += p.Top.Value()
				}
			}
		}

		if p.Bottom.Value() != 0 {
			switch p.Bottom.Mode() {
			case nux.Pixel:
				b := p.Bottom.Value()
				frame.Padding.Bottom = util.Roundi32(b)
				vPPx += b
			case nux.Ems:
				b := txtH * p.Bottom.Value()
				frame.Padding.Bottom = util.Roundi32(b)
				hPPx += b
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					b := p.Bottom.Value() / 100 * float32(nux.MeasureSpecValue(height))
					frame.Padding.Bottom = util.Roundi32(b)
					vPPx += b
				case nux.Auto:
					vPPt += p.Bottom.Value()
				}
			}
		}
	}

	innerWidth := float32(nux.MeasureSpecValue(width))*(1.0-hPPt/100.0) - hPPx
	innerHeight := float32(nux.MeasureSpecValue(height))*(1.0-vPPt/100.0) - vPPx

	hWtl, vWtl, hl, vl, hPtl, vPtl, ml := me.measureIconSize(me.iconLeft, width, height, innerWidth, innerHeight, txtH)
	hWtt, vWtt, ht, vt, hPtt, vPtt, mt := me.measureIconSize(me.iconTop, width, height, innerWidth, innerHeight, txtH)
	hWtr, vWtr, hr, vr, hPtr, vPtr, mr := me.measureIconSize(me.iconRight, width, height, innerWidth, innerHeight, txtH)
	hWtb, vWtb, hb, vb, hPtb, vPtb, mb := me.measureIconSize(me.iconBottom, width, height, innerWidth, innerHeight, txtH)

	hPx := hl + hr + txtW
	if ht > hPx {
		hPx = ht
	}
	if hb > hPx {
		hPx = hb
	}
	vPx := vt + vb + txtH
	if vl > vPx {
		vPx = vl
	}
	if vr > vPx {
		vPx = vr
	}

	if !ml {
		me.measureIconWeightSize(me.iconLeft, width, height, innerWidth-hl-txtW-hr, innerHeight-vl, hPx, vPx, hPtl+hPtr, vPtl, hWtl+hWtr, vWtl)
	}
	if !mt {
		me.measureIconWeightSize(me.iconTop, width, height, innerWidth-ht, innerHeight-vt-txtH-vb, hPx, vPx, hPtt, vPtt+vPtb, hWtt, vWtt+vWtb)
	}
	if !mr {
		me.measureIconWeightSize(me.iconRight, width, height, innerWidth-hl-txtW-hr, innerHeight-vr, hPx, vPx, hPtl+hPtr, vPtr, hWtl+hWtr, vWtr)
	}
	if !mb {
		me.measureIconWeightSize(me.iconBottom, width, height, innerWidth-hb, innerHeight-vt-txtH-vb, hPx, vPx, hPtb, vPtt+vPtb, hWtb, vWtt+vWtb)
	}

	if txtW < innerWidth-hPx {
		txtW = innerWidth - hPx
	}

	if nux.MeasureSpecMode(width) == nux.Pixel {
		frame.Width = nux.MeasureSpecValue(width)
	} else {
		// w := (innerWidth + hPPx) / (1.0 - hPPt/100.0)
		w := hPx + hPPx
		frame.Width = int32(math.Ceil(float64(w)))
	}

	if nux.MeasureSpecMode(height) == nux.Pixel {
		frame.Height = nux.MeasureSpecValue(height)
	} else {
		// h := (txtH + vPPx) / (1.0 - vPPt/100.0)
		h := vPx + vPPx
		frame.Height = int32(math.Ceil(float64(h)))
	}

	log.I("nuxui", "label measure width=%d, height=%d, txtW=%d, txtH=%d", frame.Width, frame.Height, me.textWidth, me.textHeight)
}

func (me *label) measureIconSize(icon nux.Widget, width, height int32, innerWidth, innerHeight, txtH float32) (hWt, vWt, hPxUsed, vPxUsed, hPt, vPt float32, measured bool) {
	if icon != nil {
		if s, ok := icon.(nux.Size); ok {
			if m := s.Margin(); m != nil {
				if m.Left.Value() != 0 {
					switch m.Left.Mode() {
					case nux.Pixel:
						l := m.Left.Value()
						s.Frame().Margin.Left = util.Roundi32(l)
						hPxUsed += l
					case nux.Ems:
						l := txtH * m.Left.Value()
						s.Frame().Margin.Left = util.Roundi32(l)
						hPxUsed += l
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							l := m.Left.Value() / 100 * innerWidth
							s.Frame().Margin.Left = util.Roundi32(l)
							hPxUsed += l
						case nux.Auto:
							hPt += m.Left.Value()
						}
					case nux.Weight:
						hWt += m.Left.Value()
					}
				}

				if m.Right.Value() != 0 {
					switch m.Right.Mode() {
					case nux.Pixel:
						r := m.Right.Value()
						s.Frame().Margin.Right = util.Roundi32(r)
						hPxUsed += r
					case nux.Ems:
						r := txtH * m.Right.Value()
						s.Frame().Margin.Right = util.Roundi32(r)
						hPxUsed += r
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							r := m.Right.Value() / 100 * innerWidth
							s.Frame().Margin.Right = util.Roundi32(r)
							hPxUsed += r
						case nux.Auto:
							hPt += m.Right.Value()
						}
					case nux.Weight:
						hWt += m.Right.Value()
					}
				}

				if m.Top.Value() != 0 {
					switch m.Top.Mode() {
					case nux.Pixel:
						t := m.Top.Value()
						s.Frame().Margin.Top = util.Roundi32(t)
						vPxUsed += t
					case nux.Ems:
						t := txtH * m.Top.Value()
						s.Frame().Margin.Top = util.Roundi32(t)
						vPxUsed += t
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							t := m.Top.Value() / 100 * innerHeight
							s.Frame().Margin.Top = util.Roundi32(t)
							vPxUsed += t
						case nux.Auto:
							vPt += m.Top.Value()
						}
					case nux.Weight:
						vWt += m.Top.Value()
					}
				}

				if m.Bottom.Value() != 0 {
					switch m.Bottom.Mode() {
					case nux.Pixel:
						b := m.Bottom.Value()
						s.Frame().Margin.Bottom = util.Roundi32(b)
						vPxUsed += b
					case nux.Ems:
						b := txtH * m.Bottom.Value()
						s.Frame().Margin.Bottom = util.Roundi32(b)
						vPxUsed += b
					case nux.Percent:
						switch nux.MeasureSpecMode(width) {
						case nux.Pixel:
							b := m.Bottom.Value() / 100 * innerHeight
							s.Frame().Margin.Bottom = util.Roundi32(b)
							vPxUsed += b
						case nux.Auto:
							vPt += m.Bottom.Value()
						}
					case nux.Weight:
						vWt += m.Bottom.Value()
					}
				}
			}

			canMeasure := true
			wm := s.Width().Mode()
			switch wm {
			case nux.Pixel:
				w := s.Width().Value()
				s.Frame().Width = util.Roundi32(w)
				setRatioHeight(s, s.Frame(), w, nux.Pixel)
				hPxUsed += w
			case nux.Ems:
				w := txtH * s.Width().Value()
				s.Frame().Width = util.Roundi32(w)
				setRatioHeight(s, s.Frame(), w, nux.Pixel)
				hPxUsed += w
			case nux.Percent:
				switch nux.MeasureSpecMode(width) {
				case nux.Pixel:
					w := s.Width().Value() / 100 * innerWidth
					s.Frame().Width = util.Roundi32(w)
					setRatioHeight(s, s.Frame(), w, nux.Pixel)
					hPxUsed += w
				case nux.Auto:
					hPt += s.Width().Value()
					canMeasure = false
				}
			case nux.Weight:
				hWt += s.Width().Value()
				canMeasure = false
			case nux.Auto:
				// measure later
				canMeasure = false
			}

			hm := s.Height().Mode()
			switch hm {
			case nux.Pixel:
				h := s.Height().Value()
				s.Frame().Height = util.Roundi32(h)
				setRatioWidth(s, s.Frame(), h, nux.Pixel)
				vPxUsed += h
			case nux.Ems:
				h := txtH * s.Height().Value()
				s.Frame().Height = util.Roundi32(h)
				setRatioWidth(s, s.Frame(), h, nux.Pixel)
				vPxUsed += h
			case nux.Percent:
				switch nux.MeasureSpecMode(height) {
				case nux.Pixel:
					h := s.Height().Value() / 100 * innerHeight
					s.Frame().Height = util.Roundi32(h)
					setRatioWidth(s, s.Frame(), h, nux.Pixel)
					vPxUsed += h
				case nux.Auto:
					hPt += s.Height().Value()
					canMeasure = false
				}
			case nux.Weight:
				vWt += s.Height().Value()
				canMeasure = false
			case nux.Auto:
				// measure later
				canMeasure = false
			}

			if canMeasure {
				if m, ok := icon.(nux.Measure); ok {
					m.Measure(s.Frame().Width, s.Frame().Height)
					measured = true
				}
			}
		}
	}
	return
}

func (me *label) measureIconWeightSize(icon nux.Widget, width, height int32, hPxRemain, vPxRemain, hPx, vPx, hPt, vPt, hWt, vWt float32) (measured bool) {

	return
	// TODO::
	if icon != nil {
		if s, ok := icon.(nux.Size); ok {
			switch s.Width().Mode() {
			case nux.Weight:
				s.Frame().Width = util.Roundi32(hPxRemain * s.Width().Value() / hWt)
			case nux.Auto:
			}

			switch s.Height().Mode() {
			case nux.Weight:
				s.Frame().Height = util.Roundi32(vPxRemain * s.Height().Value() / vWt)
			case nux.Auto:
			}

			if m := s.Margin(); m != nil {
				if m.Left.Value() != 0 {
					switch m.Left.Mode() {
					case nux.Weight:
						s.Frame().Margin.Left = util.Roundi32(hPxRemain * m.Left.Value() / hWt)
					}
				}

				if m.Right.Value() != 0 {
					switch m.Right.Mode() {
					case nux.Weight:
						s.Frame().Margin.Right = util.Roundi32(hPxRemain * m.Right.Value() / hWt)
					}
				}

				if m.Top.Value() != 0 {
					switch m.Top.Mode() {
					case nux.Weight:
						s.Frame().Margin.Top = util.Roundi32(vPxRemain * m.Top.Value() / vWt)
					}
				}

				if m.Bottom.Value() != 0 {
					switch m.Bottom.Mode() {
					case nux.Weight:
						s.Frame().Margin.Bottom = util.Roundi32(vPxRemain * m.Bottom.Value() / vWt)
					}
				}
			}

			if m, ok := icon.(nux.Measure); ok {
				m.Measure(s.Frame().Width, s.Frame().Height)
				measured = true
			}

		}
	}
	return
}

// Responsible for determining the position of the widget align, margin...
func (me *label) Layout(x, y, width, height int32) {
	log.I("nuxui", "label layout x=%d ,y=%d, width=%d, height=%d ", x, y, width, height)
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	var innerHeight float32 = float32(height)
	var innerWidth float32 = float32(width)

	innerHeight -= float32(frame.Padding.Top + frame.Padding.Bottom)
	innerWidth -= float32(frame.Padding.Left + frame.Padding.Right)

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

	log.I("nuxui", "label layout w2=%d, h2=%d", w2, h2)

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

	// draw text
	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X+me.textOffsetX), float32(frame.Y+me.textOffsetY))

	if me.text != "" {
		me.paint.SetTextSize(me.textSize)
		me.paint.SetColor(me.textColor)
		canvas.DrawText(me.text, float32(frame.Width), float32(frame.Height), me.paint)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
