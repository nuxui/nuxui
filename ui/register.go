// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

func init() {
	nux.RegisterWidget((*Row)(nil), func(attrs ...nux.Attr) nux.Widget { return NewRow(attrs...) })
	nux.RegisterWidget((*Column)(nil), func(attrs ...nux.Attr) nux.Widget { return NewColumn(attrs...) })
	nux.RegisterWidget((*Layer)(nil), func(attrs ...nux.Attr) nux.Widget { return NewLayer(attrs...) })
	nux.RegisterWidget((*Scroll)(nil), func(attrs ...nux.Attr) nux.Widget { return NewScroll(attrs...) })
	nux.RegisterWidget((*Text)(nil), func(attrs ...nux.Attr) nux.Widget { return NewText(attrs...) })
	nux.RegisterWidget((*Button)(nil), func(attrs ...nux.Attr) nux.Widget { return NewButton(attrs...) })
	nux.RegisterWidget((*Image)(nil), func(attrs ...nux.Attr) nux.Widget { return NewImage(attrs...) })
	nux.RegisterWidget((*Editor)(nil), func(attrs ...nux.Attr) nux.Widget { return NewEditor(attrs...) })

	nux.RegisterDrawable((*ColorDrawable)(nil), func(owner nux.Widget, attrs ...nux.Attr) nux.Drawable { return NewColorDrawable(owner, attrs...) })
	nux.RegisterDrawable((*ImageDrawable)(nil), func(owner nux.Widget, attrs ...nux.Attr) nux.Drawable { return NewImageDrawable(owner, attrs...) })
	nux.RegisterDrawable((*ShapeDrawable)(nil), func(owner nux.Widget, attrs ...nux.Attr) nux.Drawable { return NewShapeDrawable(owner, attrs...) })
	nux.RegisterDrawable((*StateDrawable)(nil), func(owner nux.Widget, attrs ...nux.Attr) nux.Drawable { return NewStateDrawable(owner, attrs...) })
}
