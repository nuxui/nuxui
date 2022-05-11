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

func InflateLayout(parent Widget, layout string) Widget {
	attrs := ParseAttr(layout)
	layoutAttr := attrs.GetAttr("layout", nil)
	if layoutAttr == nil {
		log.E("nuxui", "not found 'layout' node")
		return parent
	}

	return InflateLayoutAttr(parent, layoutAttr)
}

func InflateLayoutAttr(parent Widget, attr Attr) Widget {
	widget := inflateLayoutAttr(nil, attr)

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

func inflateLayoutAttr(parent Widget, attr Attr) Widget {
	widgetName := attr.GetString("widget", "")
	if widgetName == "" {
		log.Fatal("nuxui", `must specified "widget"`)
		return nil
	}

	widgetCreator := FindRegistedWidgetCreator(widgetName)

	var widget Widget
	if theme := attr.GetAttr("theme", nil); theme != nil {
		if themeAttr := appTheme.GetAttr(
			widgetName,
			theme.GetString("name", ""),
			theme.GetString("kind", ""),
			theme.GetString("style", "")); themeAttr != nil {
			widget = widgetCreator(MergeAttrs(themeAttr, attr))
		}
	}

	if widget == nil {
		widget = widgetCreator(attr)
	}

	widget.Info().Self = widget

	if childrenNode, ok := attr["children"]; ok {
		if children, ok := childrenNode.([]any); ok {
			if p, ok := widget.(Parent); ok {
				for _, child := range children {
					if childAttr, ok := child.(Attr); ok {
						childWidget := inflateLayoutAttr(widget, childAttr)
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
	return widget
}

// return: can be nil
func InflateDrawable(drawable any) Drawable {
	if drawable == nil {
		return nil
	}

	switch t := drawable.(type) {
	case string:
		if strings.HasPrefix(t, "#") {
			return InflateDrawable(Attr{"drawable": "nuxui.org/nuxui/ui.ColorDrawable", "color": t})
		} else if strings.HasPrefix(t, "assets/") || strings.HasPrefix(t, "http://") {
			return InflateDrawable(Attr{"drawable": "nuxui.org/nuxui/ui.ImageDrawable", "src": t})
		} else {
			if path, err := filepath.Abs(t); err == nil {
				if fileInfo, err := os.Stat(path); err == nil && !fileInfo.IsDir() {
					return InflateDrawable(Attr{"drawable": "nuxui.org/nuxui/ui.ImageDrawable", "src": path})
				}
			}
		}
	case Attr:
		return InflateDrawableAttr(t)
	}
	return nil
}

func InflateDrawableAttr(attr Attr) Drawable {
	drawableName := attr.GetString("drawable", "")
	if drawableName == "" {
		log.Fatal("nuxui", `must specified "drawable"`)
		return nil
	}

	drawableCreator := FindRegistedDrawableCreator(drawableName)
	return drawableCreator(attr)
}
