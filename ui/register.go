// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

func init() {
	nux.RegisterWidget((*Row)(nil), func(attr nux.Attr) nux.Widget { return NewRow(attr) })
	nux.RegisterWidget((*Column)(nil), func(attr nux.Attr) nux.Widget { return NewColumn(attr) })
	nux.RegisterWidget((*Layer)(nil), func(attr nux.Attr) nux.Widget { return NewLayer(attr) })
	nux.RegisterWidget((*Scroll)(nil), func(attr nux.Attr) nux.Widget { return NewScroll(attr) })
	nux.RegisterWidget((*Text)(nil), func(attr nux.Attr) nux.Widget { return NewText(attr) })
	nux.RegisterWidget((*Label)(nil), func(attr nux.Attr) nux.Widget { return NewLabel(attr) })
	nux.RegisterWidget((*Button)(nil), func(attr nux.Attr) nux.Widget { return NewButton(attr) })
	nux.RegisterWidget((*Image)(nil), func(attr nux.Attr) nux.Widget { return NewImage(attr) })
	nux.RegisterWidget((*Editor)(nil), func(attr nux.Attr) nux.Widget { return NewEditor(attr) })
	nux.RegisterWidget((*Opt)(nil), func(attr nux.Attr) nux.Widget { return NewOpt(attr) })
	nux.RegisterWidget((*Radio)(nil), func(attr nux.Attr) nux.Widget { return NewRadio(attr) })
	nux.RegisterWidget((*Options)(nil), func(attr nux.Attr) nux.Widget { return NewOptions(attr) })

	nux.RegisterDrawable((*ColorDrawable)(nil), func(attr nux.Attr) nux.Drawable { return NewColorDrawable(attr) })
	nux.RegisterDrawable((*ImageDrawable)(nil), func(attr nux.Attr) nux.Drawable { return NewImageDrawable(attr) })
	nux.RegisterDrawable((*ShapeDrawable)(nil), func(attr nux.Attr) nux.Drawable { return NewShapeDrawable(attr) })
}
