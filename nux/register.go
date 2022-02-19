// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nuxui/nuxui/log"
)

var _widgetList = make(map[string]WidgetCreator)
var _drawableList = make(map[string]DrawableCreator)

func RegisterWidget(widget any, creator WidgetCreator) {
	if debug_register {
		t := fmt.Sprintf("%T", widget)
		ss := strings.Split(t, ".")
		sname := ss[len(ss)-1]
		if sname[0] < 'A' || sname[0] > 'Z' {
			log.Fatal("nuxui", "register widget %s faild, widget is not exported", t)
		}
		if strings.Count(t, "*") != 1 {
			log.Fatal("nuxui", "register widget %s faild, widget argument should like `(*MyWidget)nil`", t)
		}
	}

	t := reflect.TypeOf(widget).Elem()
	widgetName := t.PkgPath() + "." + t.Name()
	if _, ok := _widgetList[widgetName]; ok {
		log.Fatal("nuxui", "Widget %s is already registed", widgetName)
	} else {
		_widgetList[widgetName] = creator
	}
}

func FindRegistedWidgetCreator(widget any) WidgetCreator {
	if widgetName, ok := widget.(string); ok {
		return findRegistedWidgetCreatorByName(widgetName)
	}

	t := reflect.TypeOf(widget).Elem()
	widgetName := t.PkgPath() + "." + t.Name()
	return findRegistedWidgetCreatorByName(widgetName)
}

func findRegistedWidgetCreatorByName(name string) WidgetCreator {
	if c, ok := _widgetList[name]; ok {
		return c
	}
	log.Fatal("nuxui", "widget '%s' can not find, make sure it was registed", name)
	return nil
}

func Find(widget Widget, id string) Widget {
	if id == "" {
		log.Fatal("nuxui", "the widget %T id must be specified", widget)
		return nil
	}

	w := find(widget, id)
	if w == nil {
		log.Fatal("nuxui", "the id '%s' was not found in widget %T\n", id, widget)
		return nil
	}

	return w
}

func find(widget Widget, id string) Widget {
	if id == widget.Info().ID {
		return widget
	}

	if c, ok := widget.(Component); ok {
		if ret := find(c.Content(), id); ret != nil {
			return ret
		}
	}

	if p, ok := widget.(Parent); ok {
		for _, child := range p.Children() {
			if _, ok := child.(Component); ok {
				continue
			}

			if ret := find(child, id); ret != nil {
				return ret
			}
		}
	}

	return nil
}

func RegisterDrawable(drawable any, creator DrawableCreator) {
	if debug_register {
		t := fmt.Sprintf("%T", drawable)
		ss := strings.Split(t, ".")
		sname := ss[len(ss)-1]
		if sname[0] < 'A' || sname[0] > 'Z' {
			log.Fatal("nuxui", "register drawable %s faild, drawable is not exported", t)
		}
		if strings.Count(t, "*") != 1 {
			log.Fatal("nuxui", "register drawable %s faild, drawable argument should like `(*MyDrawable)nil`", t)
		}
	}

	t := reflect.TypeOf(drawable).Elem()
	drawableName := t.PkgPath() + "." + t.Name()
	if _, ok := _drawableList[drawableName]; ok {
		log.Fatal("nuxui", "Drawable %s is already registed", drawableName)
	} else {
		_drawableList[drawableName] = creator
	}
}

func FindRegistedDrawableCreator(drawable any) DrawableCreator {
	if drawableName, ok := drawable.(string); ok {
		return findRegistedDrawableCreatorByName(drawableName)
	}

	t := reflect.TypeOf(drawable).Elem()
	drawableName := t.PkgPath() + "." + t.Name()
	return findRegistedDrawableCreatorByName(drawableName)
}

func findRegistedDrawableCreatorByName(name string) DrawableCreator {
	if c, ok := _drawableList[name]; ok {
		return c
	}
	log.Fatal("nuxui", "drawable '%s' can not find, make sure it was registed", name)
	return nil
}

func RegisterTheme(themeFunc func(string) Attr) {

}
