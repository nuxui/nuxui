// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

func InflateLayout(ctx Context, parent Widget, layout string) Widget {
	attrs := ParseAttr(layout)
	layoutAttr := attrs.GetAttr("layout", nil)
	if layoutAttr == nil {
		log.E("nuxui", "not found 'layout' node")
		return parent
	}

	return Inflate(ctx, parent, layoutAttr)
}

func Inflate(ctx Context, parent Widget, attr Attr) Widget {
	widget := inflate(ctx, nil, attr)

	if parent == nil {
		return widget
	}

	if c, ok := parent.(Component); ok {
		c.SetContent(widget)
	} else if p, ok := parent.(Parent); ok {
		p.AddChild(parent)
	} else if IsView(parent) {
		log.Fatal("nuxui", "%T is not a parent widget", parent)
		return nil
	}

	return parent
}

func inflate(ctx Context, parent Widget, attr Attr) Widget {
	widgetName := attr.GetString("widget", "")
	if widgetName == "" {
		log.Fatal("nuxui", `must specified "widget"`)
		return nil
	}

	widgetCreator := FindRegistedWidgetCreatorByName(widgetName)
	widget := widgetCreator(ctx, attr)

	if childrenNode, ok := attr["children"]; ok {
		if children, ok := childrenNode.([]interface{}); ok {
			if p, ok := widget.(Parent); ok {
				for _, child := range children {
					if childAttr, ok := child.(Attr); ok {
						childWidget := inflate(ctx, widget, childAttr)
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
