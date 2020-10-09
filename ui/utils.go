// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

func setRatioHeight(cs nux.Size, cms *nux.MeasuredSize, mode nux.Mode) {
	if cs.Height().Mode() == nux.Ratio {
		cms.Height = nux.MeasureSpec(util.Roundi32(float32(cms.Width)/cs.Height().Value()), mode)
	}
}

func setRatioWidth(cs nux.Size, cms *nux.MeasuredSize, mode nux.Mode) {
	if cs.Width().Mode() == nux.Ratio {
		cms.Width = nux.MeasureSpec(util.Roundi32(float32(cms.Height)*cs.Width().Value()), mode)
	}
}
