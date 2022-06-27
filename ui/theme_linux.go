// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

/*
	Default theme for ui widgets
*/

func button_theme(kind nux.ThemeKind) nux.Attr {
	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textColor":  "#303030",
			"textSize":   "14",
			"selectable": false,
			"clickable":  true,
			"padding":    nux.Attr{"left": "6px", "top": "3px", "right": "6px", "bottom": "3px"},
			"background": nux.Attr{
				"type": "nuxui.org/nuxui/ui.ShapeDrawable",
				"states": []nux.Attr{
					{"state": "default", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#ffffff",
						"cornerRadius": "6px",
						"shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "1px", "blur": "3px"},
					}},
					{"state": "pressed", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#f0f0f0",
						"cornerRadius": "6px",
						"shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "1px", "blur": "3px"},
					}},
				},
			},
		}
	case nux.ThemeDark:
		return nux.Attr{
			"textColor":  "#ffffff",
			"textSize":   "14",
			"selectable": false,
			"clickable":  true,
			"padding":    nux.Attr{"left": "6px", "top": "3px", "right": "6px", "bottom": "3px"},
			"background": nux.Attr{
				"type": "nuxui.org/nuxui/ui.ShapeDrawable",
				"states": []nux.Attr{
					{"state": "default", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#303030",
						"cornerRadius": "6px",
						"shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "1px", "blur": "3px"},
					}},
					{"state": "pressed", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#000000",
						"cornerRadius": "6px",
						"shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "1px", "blur": "3px"},
					}},
				},
			},
		}
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}

func text_theme(kind nux.ThemeKind) nux.Attr {
	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textColor": "#303030",
			"textSize":  "14",
		}
	case nux.ThemeDark:
		return nux.Attr{
			"textColor": "#ffffff",
			"textSize":  "14",
		}
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}

func editor_theme(kind nux.ThemeKind) nux.Attr {
	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textColor": "#303030",
			"font":      nux.Attr{"size": 14},
			"padding":   nux.Attr{"left": "6px", "top": "3px", "right": "6px", "bottom": "3px"},
			"cursor":    "",
			"background": nux.Attr{
				"type": "nuxui.org/nuxui/ui.ShapeDrawable",
				"states": []nux.Attr{
					{"state": "default", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#ffffff",
						"stroke":       "#00000022",
						"cornerRadius": "2px",
						// "shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "-1px", "blur": "3px"},
					}},
				},
			},
		}
	case nux.ThemeDark:
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}
