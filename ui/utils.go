// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

func setRatioHeight(cs nux.Size, cms *nux.Frame, width float32, mode nux.Mode) {
	if cs.Height().Mode() == nux.Ratio {
		if cms.Width == 0 {
			cms.Height = 0
		} else {
			cms.Height = nux.MeasureSpec(util.Roundi32(width/cs.Height().Value()), mode)
		}
	}
}

func setRatioWidth(cs nux.Size, cms *nux.Frame, height float32, mode nux.Mode) {
	if cs.Width().Mode() == nux.Ratio {
		if cms.Height == 0 {
			cms.Width = 0
		} else {
			cms.Width = nux.MeasureSpec(util.Roundi32(height*cs.Width().Value()), mode)
		}
	}
}

func setNewWidth(frame *nux.Frame, originWidth, newWidth int32) {
	switch nux.MeasureSpecMode(originWidth) {
	case nux.Pixel:
		frame.Width = originWidth
	case nux.Unlimit:
		frame.Width = newWidth
	case nux.Auto:
		if nux.MeasureSpecValue(newWidth) > nux.MeasureSpecValue(originWidth) {
			frame.Width = nux.MeasureSpec(nux.MeasureSpecValue(originWidth), nux.Pixel)
		} else {
			frame.Width = newWidth
		}
	}

}

func setNewHeight(frame *nux.Frame, originHeight, newHeight int32) {
	switch nux.MeasureSpecMode(originHeight) {
	case nux.Pixel:
		frame.Height = originHeight
	case nux.Unlimit:
		frame.Height = newHeight
	case nux.Auto:
		if nux.MeasureSpecValue(newHeight) > nux.MeasureSpecValue(originHeight) {
			frame.Height = nux.MeasureSpec(nux.MeasureSpecValue(originHeight), nux.Pixel)
		} else {
			frame.Height = newHeight
		}
	}
}
