// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import "nuxui.org/nuxui/nux"

func init() {
	nux.RegisterType((*Row)(nil), func(attr nux.Attr) any { return NewRow(attr) })
	nux.RegisterType((*Column)(nil), func(attr nux.Attr) any { return NewColumn(attr) })
	nux.RegisterType((*Layer)(nil), func(attr nux.Attr) any { return NewLayer(attr) })
	nux.RegisterType((*Scroll)(nil), func(attr nux.Attr) any { return NewScroll(attr) })
	nux.RegisterType((*Text)(nil), func(attr nux.Attr) any { return NewText(attr) })
	nux.RegisterType((*Label)(nil), func(attr nux.Attr) any { return NewLabel(attr) })
	nux.RegisterType((*Button)(nil), func(attr nux.Attr) any { return NewButton(attr) })
	nux.RegisterType((*Image)(nil), func(attr nux.Attr) any { return NewImage(attr) })
	nux.RegisterType((*Editor)(nil), func(attr nux.Attr) any { return NewEditor(attr) })
	nux.RegisterType((*Check)(nil), func(attr nux.Attr) any { return NewCheck(attr) })
	nux.RegisterType((*Radio)(nil), func(attr nux.Attr) any { return NewRadio(attr) })
	nux.RegisterType((*Options)(nil), func(attr nux.Attr) any { return NewOptions(attr) })
	// nux.RegisterType((*SeekBar)(nil), func(attr nux.Attr) any { return NewSeekBar(attr) })
	nux.RegisterType((*Switch)(nil), func(attr nux.Attr) any { return NewSwitch(attr) })

	nux.RegisterType((*ColorDrawable)(nil), func(attr nux.Attr) any { return NewColorDrawable(attr) })
	nux.RegisterType((*ImageDrawable)(nil), func(attr nux.Attr) any { return NewImageDrawable(attr) })
	nux.RegisterType((*ShapeDrawable)(nil), func(attr nux.Attr) any { return NewShapeDrawable(attr) })
}
