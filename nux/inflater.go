// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"os"
	"path/filepath"
	"strings"

	"nuxui.org/nuxui/log"
)

func InflateStyles(styles ...string) Attr {
	attrs := make([]Attr, len(styles))
	for i, style := range styles {
		attrs[i] = InflateStyle(style)
	}
	return MergeAttrs(attrs...)
}

func InflateStyle(style string) Attr {
	attr := ParseAttr(style)
	if s := attr.GetAttr("style", nil); s != nil {
		return s
	}

	log.Fatal("nuxui", "not found 'style' node")
	return nil
}

func InflateLayout(parent Widget, layout string, styles Attr) Widget {
	attr := ParseAttr(layout)
	layoutAttr := attr.GetAttr("layout", nil)
	if layoutAttr == nil {
		log.Fatal("nuxui", "not found 'layout' node")
		return parent
	}

	return InflateLayoutAttr(parent, layoutAttr, styles)
}

func InflateLayoutAttr(parent Widget, layout Attr, styles Attr) Widget {
	widget := inflateLayoutAttr(nil, layout, styles)

	if parent == nil {
		return widget
	}

	if c, ok := parent.(Component); ok {
		c.SetContent(widget)
	} else if p, ok := parent.(Parent); ok {
		p.AddChild(widget)
	} else if isView(parent) {
		log.Fatal("nuxui", "%T is not a parent widget", parent)
		return nil
	}

	return parent
}

func inflateLayoutAttr(parent Widget, layoutAttr Attr, styleAttr Attr) Widget {
	var widgetTheme, widgetStyle Attr
	if themeNames := layoutAttr.GetStringArray("theme", nil); themeNames != nil {
		widgetTheme = AppTheme().GetStyle(themeNames...)
	}

	if styleNames := layoutAttr.GetStringArray("style", nil); styleNames != nil {
		for _, styleName := range styleNames {
			if style := styleAttr.GetAttr(styleName, nil); style != nil {
				widgetStyle = MergeAttrs(widgetStyle, style)
			}
		}
	}

	mergeAttr := MergeAttrs(widgetTheme, widgetStyle, layoutAttr)
	typeFullName := mergeAttr.GetString("type", "")
	if typeFullName == "" {
		log.Fatal("nuxui", `must specified "type" node in layout attr`)
		return nil
	}

	widgetCreator := FindTypeCreator(typeFullName)

	if widget, ok := widgetCreator(mergeAttr).(Widget); ok {
		widget.Info().Self = widget

		if childrenNode, ok := layoutAttr["children"]; ok {
			if children, ok := childrenNode.([]any); ok {
				if p, ok := widget.(Parent); ok {
					for _, child := range children {
						if childAttr, ok := child.(Attr); ok {
							childWidget := inflateLayoutAttr(widget, childAttr, styleAttr)
							p.AddChild(childWidget)
						}
					}
				} else {
					log.Fatal("nuxui", "%T is not a WidgetParent but has 'children' node", p)
				}
			} else {
				log.Fatal("nuxui", "'children' node must be widget array")
			}
		}

		if content, ok := layoutAttr["content"]; ok {
			if compt, ok := widget.(Component); ok {
				if childAttr, ok := content.(Attr); ok {
					childWidget := inflateLayoutAttr(widget, childAttr, styleAttr)
					compt.SetContent(childWidget)
				}
			}
		}
		return widget
	} else {
		log.Fatal("nuxui", `the type "%s" is not a Widget`, typeFullName)
	}

	return nil
}

// return: can be nil
func InflateDrawable(drawable any) Drawable {
	if drawable == nil {
		return nil
	}

	switch t := drawable.(type) {
	case string:
		if strings.HasPrefix(t, "#") {
			return InflateDrawable(Attr{"type": "nuxui.org/nuxui/ui.ColorDrawable", "color": t})
		} else if strings.Count(t, "/") > 1 {
			return InflateDrawable(Attr{"type": "nuxui.org/nuxui/ui.ImageDrawable", "src": t})
		} else {
			if path, err := filepath.Abs(t); err == nil {
				if fileInfo, err := os.Stat(path); err == nil && !fileInfo.IsDir() {
					return InflateDrawable(Attr{"type": "nuxui.org/nuxui/ui.ImageDrawable", "src": path})
				}
			}
		}
	case Attr:
		return InflateDrawableAttr(t)
	}
	return nil
}

// return: can be nil
func InflateDrawableAttr(attr Attr) Drawable {
	if attr == nil {
		return nil
	}

	typeFullName := attr.GetString("type", "")
	if typeFullName == "" {
		log.Fatal("nuxui", `must specified "type" node in drawable attr %s`, attr)
		return nil
	}

	drawableCreator := FindTypeCreator(typeFullName)
	if d, ok := drawableCreator(attr).(Drawable); ok {
		return d
	}
	log.Fatal("nuxui", `the type "%s" is not a Drawable`, typeFullName)
	return nil
}
