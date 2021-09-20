// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Image interface {
	nux.Widget
	nux.Size
	Visual
	nux.Creating
	nux.Created
	nux.Layout
	nux.Measure
	nux.Draw

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

func NewImage() Image {
	me := &image{
		scaleX:    1.0,
		scaleY:    1.0,
		offsetX:   0,
		offsetY:   0,
		scaleType: ScaleType_Matrix,
	}
	me.WidgetSize.Owner = me
	me.WidgetVisual.Owner = me
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)
	return me
}

type image struct {
	nux.WidgetBase
	nux.WidgetSize
	WidgetVisual

	scaleType   ScaleType
	src         string
	srcDrawable ImageDrawable
	scaleX      float32
	scaleY      float32
	offsetX     float32
	offsetY     float32
}

func (me *image) Creating(attr nux.Attr) {
	if attr == nil {
		attr = nux.Attr{}
	}

	me.WidgetBase.Creating(attr)
	me.WidgetSize.Creating(attr)
	me.WidgetVisual.Creating(attr)

	me.src = attr.GetString("src", "")
	me.scaleType = convertScaleTypeFromString(attr.GetString("scaleType", "matrix"))
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

func (me *image) Created(content nux.Widget) {
	if me.src != "" {
		me.srcDrawable = NewImageDrawable(me.src)
	}
}

// TODO if not have Layout, then use default layout to set frame
func (me *image) Layout(dx, dy, left, top, right, bottom int32) {
	ms := me.MeasuredSize()

	// setFrame
	ms.Position.Left = left
	ms.Position.Top = top
	ms.Position.Right = right
	ms.Position.Bottom = bottom
	ms.Position.X = dx
	ms.Position.Y = dy

	var imgW, imgH float32
	if me.srcDrawable != nil {
		imgW = float32(me.srcDrawable.Width())
		imgH = float32(me.srcDrawable.Height())
	}
	innerW := ms.Width - ms.Padding.Left - ms.Padding.Right
	innerH := ms.Height - ms.Padding.Top - ms.Padding.Bottom

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
		r2 := float32(innerW) / float32(innerH)
		if r2 > r {
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
}

func (me *image) Measure(width, height int32) {
	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
		ms := me.MeasuredSize()
		if nux.MeasureSpecMode(width) == nux.Auto && me.srcDrawable != nil {
			ms.Width = nux.MeasureSpec(me.srcDrawable.Width()+ms.Padding.Left+ms.Padding.Right, nux.Pixel)
		} else {
			ms.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto && me.srcDrawable != nil {
			ms.Height = nux.MeasureSpec(me.srcDrawable.Height()+ms.Padding.Top+ms.Padding.Bottom, nux.Pixel)
		} else {
			ms.Height = height
		}
	}

	if me.HasPadding() {
		ms := me.MeasuredSize()

		switch me.PaddingLeft().Mode() {
		case nux.Pixel:
			ms.Padding.Left = util.Roundi32(me.PaddingLeft().Value())
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			ms.Padding.Top = util.Roundi32(me.PaddingTop().Value())
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			ms.Padding.Right = util.Roundi32(me.PaddingRight().Value())
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			ms.Padding.Bottom = util.Roundi32(me.PaddingBottom().Value())
		}
	}
}

func (me *image) onSizeChanged(widget nux.Widget) {

}

func (me *image) onVisualChanged(widget nux.Widget) {
}

func (me *image) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	canvas.Save()
	ms := me.MeasuredSize()
	canvas.Translate(ms.Padding.Left, ms.Padding.Top)
	canvas.ClipRect(0, 0,
		ms.Width-ms.Padding.Left-ms.Padding.Right,
		ms.Height-ms.Padding.Top-ms.Padding.Bottom)

	if me.srcDrawable != nil {
		canvas.TranslateF(me.offsetX, me.offsetY)
		canvas.ScaleF(me.scaleX, me.scaleY)
		me.srcDrawable.Draw(canvas)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}

func (me *image) Src() string {
	return me.src
}

func (me *image) SetSrc(src string) {
	if me.src == src {
		return
	}

	me.src = src

	if me.src != "" {
		me.srcDrawable = NewImageDrawable(me.src)
	}

	nux.RequestRedraw(me)
}

func (me *image) ScaleType() ScaleType {
	return me.scaleType
}

func (me *image) SetScaleType(scaleType ScaleType) {
	me.scaleType = scaleType
}
