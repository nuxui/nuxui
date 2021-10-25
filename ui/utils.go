// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

func setRatioHeight(cs nux.Size, cms *nux.MeasuredSize, width float32, mode nux.Mode) {
	if cs.Height().Mode() == nux.Ratio {
		if cms.Width == 0 {
			cms.Height = 0
		} else {
			cms.Height = nux.MeasureSpec(util.Roundi32(width/cs.Height().Value()), mode)
		}
	}
}

func setRatioWidth(cs nux.Size, cms *nux.MeasuredSize, height float32, mode nux.Mode) {
	if cs.Width().Mode() == nux.Ratio {
		if cms.Height == 0 {
			cms.Width = 0
		} else {
			cms.Width = nux.MeasureSpec(util.Roundi32(height*cs.Width().Value()), mode)
		}
	}
}

func setNewWidth(ms *nux.MeasuredSize, originWidth, newWidth int32) {
	switch nux.MeasureSpecMode(originWidth) {
	case nux.Pixel:
		ms.Width = originWidth
	case nux.Unlimit:
		ms.Width = newWidth
	case nux.Auto:
		if nux.MeasureSpecValue(newWidth) > nux.MeasureSpecValue(originWidth) {
			ms.Width = nux.MeasureSpec(nux.MeasureSpecValue(originWidth), nux.Pixel)
		} else {
			ms.Width = newWidth
		}
	}

}

func setNewHeight(ms *nux.MeasuredSize, originHeight, newHeight int32) {
	switch nux.MeasureSpecMode(originHeight) {
	case nux.Pixel:
		ms.Height = originHeight
	case nux.Unlimit:
		ms.Height = newHeight
	case nux.Auto:
		if nux.MeasureSpecValue(newHeight) > nux.MeasureSpecValue(originHeight) {
			ms.Height = nux.MeasureSpec(nux.MeasureSpecValue(originHeight), nux.Pixel)
		} else {
			ms.Height = newHeight
		}
	}
}

func getAttr(attrs ...nux.Attr) nux.Attr {
	if len(attrs) == 0 {
		return nux.Attr{}
	}
	return attrs[0]
}
