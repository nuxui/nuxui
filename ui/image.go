// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
	"nuxui.org/nuxui/util"
)

type Image interface {
	nux.Widget
	nux.Size
	nux.Layout
	nux.Measure
	nux.Draw
	nux.Stateable
	Visual

	Src() string
	SetSrc(src string)
	ScaleType() ScaleType
	SetScaleType(scaleType ScaleType)
}

type ScaleType int32

const (
	ScaleType_Matrix ScaleType = iota
	ScaleType_Center
	ScaleType_CenterCrop
	ScaleType_CenterInside
	ScaleType_FitXY
	ScaleType_FitStart
	ScaleType_FitCenter
	ScaleType_FitEnd
)

type Repeat int32

func NewImage(attr nux.Attr) Image {
	if attr == nil {
		attr = nux.Attr{}
	}

	me := &image{
		scaleX:    1.0,
		scaleY:    1.0,
		offsetX:   0,
		offsetY:   0,
		scaleType: convertScaleTypeFromString(attr.GetString("scaleType", "fitCenter")),
		state:     nux.State_Default,
	}
	me.WidgetBase = nux.NewWidgetBase(attr)
	me.WidgetSize = nux.NewWidgetSize(attr)
	me.WidgetVisual = NewWidgetVisual(me, attr)
	me.srcDrawable = nux.InflateDrawable(attr.Get("src", nil))
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	return me
}

type image struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	state       uint32
	scaleType   ScaleType
	scaleX      float32
	scaleY      float32
	offsetX     float32
	offsetY     float32
	src         string
	srcDrawable ImageDrawable
}

func convertScaleTypeFromString(scaleType string) ScaleType {
	switch scaleType {
	case "matrix":
		return ScaleType_Matrix
	case "center":
		return ScaleType_Center
	case "centerCrop":
		return ScaleType_CenterCrop
	case "centerInside":
		return ScaleType_CenterInside
	case "fitXY":
		return ScaleType_FitXY
	case "fitStart":
		return ScaleType_FitStart
	case "fitCenter":
		return ScaleType_FitCenter
	case "fitEnd":
		return ScaleType_FitEnd
	}

	log.Fatal("nux", fmt.Sprintf("unknow scale type %s, only support 'matrix', 'center', 'centerCrop', 'centerInside', 'fitXY', 'fitStart', 'fitEnd'", scaleType))
	return ScaleType_Center
}

func (me *image) AddState(state uint32) {
	s := me.state
	s |= state
	me.state = s

	if me.background != nil {
		me.background.AddState(state)
	}

	if me.srcDrawable != nil {
		me.srcDrawable.AddState(state)
	}

	if me.foreground != nil {
		me.background.AddState(state)
	}
}

func (me *image) DelState(state uint32) {
	s := me.state
	s &= ^state
	me.state = s

	if me.background != nil {
		me.background.DelState(state)
	}

	if me.srcDrawable != nil {
		me.srcDrawable.DelState(state)
	}

	if me.foreground != nil {
		me.background.DelState(state)
	}
}

func (me *image) State() uint32 {
	return me.state
}

func (me *image) Mount() {
}

func (me *image) Measure(width, height nux.MeasureDimen) {
	frame := me.Frame()

	hPPx, hPPt, vPPx, vPPt, paddingMeasuredFlag := measurePadding(width, height, me.Padding(), frame, -1, 0)
	if hPPt >= 100.0 || vPPt >= 100.0 {
		log.Fatal("nuxui", "padding percent size should at 0% ~ 100%")
	}

	var dw, dh int32
	if me.srcDrawable != nil {
		dw, dh = me.srcDrawable.Size()
	}
	if width.Mode() == nux.Pixel {
		frame.Width = width.Value()
	} else {
		frame.Width = util.Roundi32((float32(dw) + hPPx) / (1 - hPPt/100.0))
		width = nux.MeasureSpec(frame.Width, nux.Pixel)
	}
	if height.Mode() == nux.Pixel {
		frame.Height = height.Value()
	} else {
		frame.Height = util.Roundi32((float32(dh) + vPPx) / (1 - vPPt/100.0))
		height = nux.MeasureSpec(frame.Height, nux.Pixel)
	}

	if paddingMeasuredFlag&flagMeasuredPaddingComplete != flagMeasuredPaddingComplete {
		measurePadding(width, height, me.Padding(), frame, -1, paddingMeasuredFlag)
	}
	return
}

func (me *image) Layout(x, y, width, height int32) {
	if me.Background() != nil {
		me.Background().SetBounds(x, y, width, height)
	}

	if me.Foreground() != nil {
		me.Foreground().SetBounds(x, y, width, height)
	}

	frame := me.Frame()

	var imgW, imgH float32
	if me.srcDrawable != nil {
		w, h := me.srcDrawable.Size()
		imgW = float32(w)
		imgH = float32(h)
	}
	innerW := frame.Width - frame.Padding.Left - frame.Padding.Right
	innerH := frame.Height - frame.Padding.Top - frame.Padding.Bottom

	if imgW == 0 || imgH == 0 || innerW == 0 || innerH == 0 {
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = 0
		me.offsetY = 0
		return
	}

	switch me.scaleType {
	case ScaleType_Matrix:
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = 0
		me.offsetY = 0
	case ScaleType_Center:
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = (float32(innerW) - imgW) / 2
		me.offsetY = (float32(innerH) - imgH) / 2
	case ScaleType_CenterCrop:
		r := imgW / imgH
		ir := float32(innerW) / float32(innerH)
		if ir > r {
			newH := float32(innerW) / r
			me.scaleX = float32(innerW) / imgW
			me.scaleY = newH / imgH
			me.offsetX = 0
			me.offsetY = (float32(innerH) - newH) / 2
		} else {
			newW := float32(innerH) * r
			me.scaleX = newW / imgW
			me.scaleY = float32(innerH) / imgH
			me.offsetX = (float32(innerW) - newW) / 2
			me.offsetY = 0
		}
	case ScaleType_CenterInside:
		if imgW > float32(innerW) || imgH > float32(innerH) {
			r := imgW / imgH
			r2 := float32(innerW) / float32(innerH)
			if r2 > r {
				newW := float32(innerH) * r
				me.scaleX = newW / imgW
				me.scaleY = float32(innerH) / imgH
				me.offsetX = (float32(innerW) - newW) / 2
				me.offsetY = 0
			} else {
				newH := float32(innerW) / r
				me.scaleX = float32(innerW) / imgW
				me.scaleY = newH / imgH
				me.offsetX = 0
				me.offsetY = (float32(innerH) - newH) / 2
			}
		} else {
			me.scaleX = 1.0
			me.scaleY = 1.0
			me.offsetX = (float32(innerW) - imgW) / 2
			me.offsetY = (float32(innerH) - imgH) / 2
		}

	case ScaleType_FitXY:
		me.scaleX = float32(innerW) / imgW
		me.scaleY = float32(innerH) / imgH
		me.offsetX = 0
		me.offsetY = 0
	case ScaleType_FitCenter, ScaleType_FitStart, ScaleType_FitEnd:
		r := imgW / imgH
		r2 := float32(innerW) / float32(innerH)
		if r2 > r {
			newW := float32(innerH) * r
			me.scaleX = newW / imgW
			me.scaleY = float32(innerH) / imgH

			switch me.scaleType {
			case ScaleType_FitStart:
				me.offsetX = 0
				me.offsetY = 0
			case ScaleType_FitCenter:
				me.offsetX = (float32(innerW) - newW) / 2
				me.offsetY = 0
			case ScaleType_FitEnd:
				me.offsetX = float32(innerW) - newW
				me.offsetY = 0
			}
		} else {
			newH := float32(innerW) / r
			me.scaleX = float32(innerW) / imgW
			me.scaleY = newH / imgH

			switch me.scaleType {
			case ScaleType_FitStart:
				me.offsetX = 0
				me.offsetY = 0
			case ScaleType_FitCenter:
				me.offsetX = 0
				me.offsetY = (float32(innerH) - newH) / 2
			case ScaleType_FitEnd:
				me.offsetX = 0
				me.offsetY = float32(innerH) - newH
			}
		}
	}

	if me.srcDrawable != nil {
		me.srcDrawable.SetBounds(frame.X+frame.Padding.Left+util.Roundi32(me.offsetX), frame.Y+frame.Padding.Top+util.Roundi32(me.offsetY), util.Roundi32(imgW*me.scaleX), util.Roundi32(imgH*me.scaleY))
	}
}

func (me *image) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *image) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	if me.srcDrawable != nil {
		me.srcDrawable.Draw(canvas)
	}

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}

func (me *image) Src() string {
	// return me.srcDrawable.Src()
	return ""
}

func (me *image) SetSrc(src string) {
	if me.src == src {
		return
	}

	me.src = src

	if me.src != "" {
		me.srcDrawable = NewImageDrawableWithResource(me.src)
	}

	nux.RequestLayout(me)
	nux.RequestRedraw(me)
}

func (me *image) ScaleType() ScaleType {
	return me.scaleType
}

func (me *image) SetScaleType(scaleType ScaleType) {
	me.scaleType = scaleType
}
