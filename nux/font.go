// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"strings"
)

// TODO interface
type Font struct {
	Family string // TODO:: face or family
	Size   float32
	Color  Color
	Weight int
}

func (me *Font) Creating(attr Attr) {
	me.Size = attr.GetFloat32("size", 14.0)
	me.Color = attr.GetColor("color", Black)
	family := strings.Split(attr.GetString("family", "Sans"), ",")
	for i, f := range family {
		family[i] = strings.TrimSpace(f)
		if strings.ContainsRune(family[i], ' ') {
			family[i] = `"` + family[i] + `"`
		}
	}

	me.Family = strings.Join(family, ",")
}

// typedef enum {
// 	PANGO_ELLIPSIZE_NONE,
// 	PANGO_ELLIPSIZE_START,
// 	PANGO_ELLIPSIZE_MIDDLE,
// 	PANGO_ELLIPSIZE_END
//} PangoEllipsizeMode;
