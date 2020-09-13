// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "github.com/nuxui/nuxui/nux"

func init() {
	nux.RegisterWidget((*Row)(nil), func() nux.Widget { return NewRow() })
	nux.RegisterWidget((*Column)(nil), func() nux.Widget { return NewColumn() })
	// RegisterWidget((*Layer)(nil), func() Widget { return NewLayer() })
	nux.RegisterWidget((*Text)(nil), func() nux.Widget { return NewText() })
	// RegisterWidget((*Image)(nil), func() Widget { return NewImage() })
	// RegisterWidget((*Editor)(nil), func() Widget { return NewEditor() })
	// RegisterWidget((*Pager)(nil), func() Widget { return NewPager() })

	// RegisterMixin((*GestureMixin)(nil), func() Mixin { return NewGestureMixin() })
}
