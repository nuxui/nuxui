// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nuxui/nuxui/log"
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

	if len(color) >= 1 && color[0] == '#' {
		color = strings.Replace(color, "#", "", 1)
	} else if len(color) >= 2 && color[0] == '0' && color[1] == 'x' {
		color = strings.Replace(color, "0x", "", 1)
	}

	i, e := strconv.ParseUint(color, 16, 32)
	if e != nil {
		log.E("color", fmt.Sprintln(e))
		return defaultValue
	}

	return Color(i)
}
