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
		cms.Height = nux.MeasureSpec(util.Roundi32(width/cs.Height().Value()), mode)
	}
}

func setRatioWidth(cs nux.Size, cms *nux.MeasuredSize, height float32, mode nux.Mode) {
	if cs.Width().Mode() == nux.Ratio {
		cms.Width = nux.MeasureSpec(util.Roundi32(height*cs.Width().Value()), mode)
	}
}
