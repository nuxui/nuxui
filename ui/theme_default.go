// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux"
)

const (
	svg_check_default = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M832 928 192 928c-52.928 0-96-43.072-96-96L96 192c0-52.928 43.072-96 96-96l640 0c52.928 0 96 43.072 96 96l0 640C928 884.928 884.928 928 832 928zM192 160C174.368 160 160 174.368 160 192l0 640c0 17.664 14.368 32 32 32l640 0c17.664 0 32-14.336 32-32L864 192c0-17.632-14.336-32-32-32L192 160z" fill="#e6e6e6"></path></svg>`
	svg_check_checked = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M832 96 192 96C139.072 96 96 139.072 96 192l0 640c0 52.928 43.072 96 96 96l640 0c52.928 0 96-43.072 96-96L928 192C928 139.072 884.928 96 832 96zM727.232 438.432l-256.224 259.008c-0.064 0.064-0.192 0.096-0.256 0.192-0.096 0.064-0.096 0.192-0.192 0.256-2.048 1.984-4.576 3.2-6.944 4.544-1.184 0.672-2.144 1.696-3.392 2.176-3.84 1.536-7.904 2.336-11.968 2.336-4.096 0-8.224-0.8-12.096-2.4-1.28-0.544-2.304-1.632-3.52-2.304-2.368-1.344-4.832-2.528-6.88-4.544-0.064-0.064-0.096-0.192-0.16-0.256-0.064-0.096-0.192-0.096-0.256-0.192l-126.016-129.504c-12.32-12.672-12.032-32.928 0.64-45.248 12.672-12.288 32.896-12.064 45.248 0.64l103.264 106.112 233.28-235.808c12.416-12.576 32.704-12.672 45.248-0.256C739.52 405.632 739.648 425.888 727.232 438.432z" fill="#1296db"></path></svg>`
	svg_radio_default = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M512 960C264.96 960 64 759.04 64 512S264.96 64 512 64s448 200.96 448 448S759.04 960 512 960zM512 128C300.256 128 128 300.256 128 512c0 211.744 172.256 384 384 384 211.744 0 384-172.256 384-384C896 300.256 723.744 128 512 128z" fill="#e6e6e6"></path></svg>`
	svg_radio_checked = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M510.54537 62.365396c-247.564375 0-448.16002 200.594621-448.16002 448.16002s200.595644 448.161043 448.16002 448.161043S958.70539 758.088768 958.70539 510.524392 758.109746 62.365396 510.54537 62.365396L510.54537 62.365396zM510.54537 921.410484c-226.449475 0-410.885068-184.219675-410.885068-410.886091 0-226.449475 184.218652-410.885068 410.885068-410.885068 226.449475 0 410.885068 184.219675 410.885068 410.885068C921.430438 736.974891 736.994845 921.410484 510.54537 921.410484L510.54537 921.410484zM510.54537 921.410484" fill="#1296db"></path><path d="M510.54537 165.786861c-190.470029 0-344.737532 154.267503-344.737532 344.737532 0 190.463889 154.267503 344.737532 344.737532 344.737532 190.463889 0 344.737532-154.273642 344.737532-344.737532C855.282902 320.054363 701.00926 165.786861 510.54537 165.786861L510.54537 165.786861zM510.54537 165.786861" fill="#1296db"></path></svg>`
	svg_switch_on     = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M309.794 819.848h404.239c169.953 0 307.726-137.774 307.726-307.727s-137.774-307.726-307.726-307.726h-404.239c-169.952 0-307.726 137.773-307.726 307.726s137.774 307.727 307.726 307.727zM709.794 233.5c154.23 0 279.268 125.031 279.268 279.268s-125.039 279.268-279.268 279.268c-154.24 0-279.269-125.033-279.269-279.268 0-154.238 125.030-279.268 279.269-279.268z" fill="#1296fa"></path></svg>`
	svg_switch_off    = `<svg viewBox="0 0 1024 1024" width="32" height="32"><path d="M309.794 819.848h404.239c169.953 0 307.726-137.774 307.726-307.727s-137.774-307.726-307.726-307.726h-404.239c-169.952 0-307.726 137.773-307.726 307.726s137.774 307.727 307.726 307.727zM309.445 233.5c154.23 0 279.268 125.031 279.268 279.268s-125.039 279.268-279.268 279.268c-154.24 0-279.269-125.033-279.269-279.268 0-154.238 125.030-279.268 279.269-279.268z" fill="#e6e6e6"></path></svg>`
)

var (
	img_check_default nux.Image
	img_check_checked nux.Image
	img_radio_default nux.Image
	img_radio_checked nux.Image
	img_switch_on     nux.Image
	img_switch_off    nux.Image
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
					{"state": "hovered | pressed", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#f0f0f0",
						"cornerRadius": "6px",
						"shadow":       nux.Attr{"color": "#00000088", "x": 0, "y": "1px", "blur": "3px"},
					}},
					{"state": "hovered", "shape": nux.Attr{
						"shape":        "rect",
						"solid":        "#f6f6f6",
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
		// TODO::
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}

func check_theme(kind nux.ThemeKind) nux.Attr {
	if img_check_default == nil {
		img_check_default = nux.ReadImageSVG(svg_check_default)
		img_check_checked = nux.ReadImageSVG(svg_check_checked)
	}

	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textSize":  14,
			"textColor": "#ffffff",
			"margin":    nux.Attr{"left": "16px"},
			"padding":   nux.Attr{"left": "0px", "top": "8px", "right": "0px", "bottom": "8px"},
			"icon": nux.Attr{
				"left": nux.Attr{
					"type":   "nuxui.org/nuxui/ui.Image",
					"width":  "1.2em",
					"height": "1.2em",
					"margin": nux.Attr{"right": "6px"},
					"src": nux.Attr{
						"type": "nuxui.org/nuxui/ui.ImageDrawable",
						"states": []nux.Attr{
							{"state": "default", "src": img_check_default},
							{"state": "checked", "src": img_check_checked},
						},
					},
				},
			},
		}
	case nux.ThemeDark:
		// TODO::
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}

func radio_theme(kind nux.ThemeKind) nux.Attr {
	if img_radio_default == nil {
		img_radio_default = nux.ReadImageSVG(svg_radio_default)
		img_radio_checked = nux.ReadImageSVG(svg_radio_checked)
	}

	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textSize":  14,
			"textColor": "#ffffff",
			"margin":    nux.Attr{"left": "16px"},
			"padding":   nux.Attr{"left": "0px", "top": "8px", "right": "0px", "bottom": "8px"},
			"icon": nux.Attr{
				"left": nux.Attr{
					"type":   "nuxui.org/nuxui/ui.Image",
					"width":  "1.2em",
					"height": "1.2em",
					"margin": nux.Attr{"right": "6px"},
					"src": nux.Attr{
						"type": "nuxui.org/nuxui/ui.ImageDrawable",
						"states": []nux.Attr{
							{"state": "default", "src": img_radio_default},
							{"state": "checked", "src": img_radio_checked},
						},
					},
				},
			},
		}
	case nux.ThemeDark:
		// TODO::
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}

func switch_theme(kind nux.ThemeKind) nux.Attr {
	if img_switch_on == nil {
		img_switch_on = nux.ReadImageSVG(svg_switch_on)
		img_switch_off = nux.ReadImageSVG(svg_switch_off)
	}

	switch kind {
	case nux.ThemeLight:
		return nux.Attr{
			"textSize":  14,
			"textColor": "#ffffff",
			"margin":    nux.Attr{"left": "16px"},
			"padding":   nux.Attr{"left": "0px", "top": "8px", "right": "0px", "bottom": "8px"},
			"icon": nux.Attr{
				"left": nux.Attr{
					"type":   "nuxui.org/nuxui/ui.Image",
					"width":  "1.2em",
					"height": "1.2em",
					"margin": nux.Attr{"right": "6px"},
					"src": nux.Attr{
						"type": "nuxui.org/nuxui/ui.ImageDrawable",
						"states": []nux.Attr{
							{"state": "default", "src": img_switch_off},
							{"state": "checked", "src": img_switch_on},
						},
					},
				},
			},
		}
	case nux.ThemeDark:
		// TODO::
	}
	log.E("nuxui", "unknow theme kind")
	return nux.Attr{}
}
