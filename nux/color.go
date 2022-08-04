// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"strconv"
	"strings"
)

// rgba color
type Color uint32

const (
	Transparent Color = 0x00000000
	White       Color = 0xFFFFFFFF
	Black       Color = 0x000000FF
	Red         Color = 0xFF0000FF
	Green       Color = 0x00FF00FF
	Blue        Color = 0x0000FFFF
	Purple      Color = 0xFF00FFFF
	Yellow      Color = 0xFFFF00FF
)

func ParseColor(color string, defaultValue Color) (Color, error) {
	c := strings.TrimSpace(color)

	if c == "" {
		return defaultValue, fmt.Errorf("empty string for parse color, use default value instead")
	} else if c[0] == '#' {
		c = c[1:]
	} else if strlen(c) >= 2 && c[0] == '0' && (c[1] == 'x' || c[1] == 'X') {
		c = c[2:]
	} else {
		return defaultValue, fmt.Errorf("color should start with #, use default value instead")
	}

	if strlen(c) == 6 {
		c += "FF"
	}
	if strlen(c) != 8 {
		return defaultValue, fmt.Errorf("cannot convert %s to Color, use default value instead", color)
	}

	i, e := strconv.ParseUint(c, 16, 32)
	if e != nil {
		return defaultValue, fmt.Errorf("parse color error: %s", e.Error())
	}

	return Color(i), nil
}

func (me Color) RGBAf() (r, g, b, a float32) {
	r = float32((me>>24)&0xff) / 255
	g = float32((me>>16)&0xff) / 255
	b = float32((me>>8)&0xff) / 255
	a = float32((me)&0xff) / 255
	return
}

func (me Color) ARGB() uint32 {
	return uint32(((me & 0xff) << 24) | (me >> 8))
}

func (me Color) Equal(color Color) bool {
	return uint32(me) == uint32(color)
}

func FromARGB(color uint32) Color {
	return Color(((color & 0xff) << 24) | (color >> 8))
}
