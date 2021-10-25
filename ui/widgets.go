// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

func init() {
	nux.RegisterWidget((*Row)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewRow(ctx, attr...) })
	nux.RegisterWidget((*Column)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewColumn(ctx, attr...) })
	nux.RegisterWidget((*Layer)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewLayer(ctx, attr...) })
	nux.RegisterWidget((*Scroll)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewScroll(ctx, attr...) })
	nux.RegisterWidget((*Text)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewText(ctx, attr...) })
	nux.RegisterWidget((*Image)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewImage(ctx, attr...) })
	nux.RegisterWidget((*Editor)(nil), func(ctx nux.Context, attr ...nux.Attr) nux.Widget { return NewEditor(ctx, attr...) })
	// RegisterWidget((*Pager)(nil), func() Widget { return NewPager() })

	// RegisterMixin((*GestureMixin)(nil), func() Mixin { return NewGestureMixin() })
}
