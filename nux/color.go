// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"strconv"
	"strings"

	"nuxui.org/nuxui/log"
)

type Color uint32

const (
	Transparent Color = 0x00000000
	White       Color = 0xFFFFFFFF
	Black       Color = 0xFF000000
	Red         Color = 0xFFFF0000
	Green       Color = 0xFF00FF00
	Blue        Color = 0xFF0000FF
	Purple      Color = 0xFFFF00FF
	Yellow      Color = 0xFFFFFF00
)

func ParseColor(color string, defaultValue Color) Color {
	if color == "" {
		return defaultValue
	}

	if strlen(color) >= 1 && color[0] == '#' {
		color = strings.Replace(color, "#", "", 1)
	} else if strlen(color) >= 2 && color[0] == '0' && color[1] == 'x' {
		color = strings.Replace(color, "0x", "", 1)
	}

	i, e := strconv.ParseUint(color, 16, 32)
	if e != nil {
		log.E("color", "%s", e)
		return defaultValue
	}

	return Color(i)
}

func (me Color) ARGBf() (a, r, g, b float32) {
	a = float32((me>>24)&0xff) / 255
	r = float32((me>>16)&0xff) / 255
	g = float32((me>>8)&0xff) / 255
	b = float32((me)&0xff) / 255
	return
}

func (me Color) Equal(color Color) bool {
	return uint32(me) == uint32(color)
}
