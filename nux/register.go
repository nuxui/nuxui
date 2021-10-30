// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"reflect"

	"github.com/nuxui/nuxui/log"
)

var _widgetList = make(map[string]Creator)
var _mixinList = make(map[string]MixinCreator)

// func CreateWidget(widgetType interface{}) Widget {
// 	// TODO:: check widgetType is struct and (*Widget)nil
// 	// return reflect.New(reflect.TypeOf(widgetType)).Elem().Interface()
// 	i := reflect.New(reflect.TypeOf(widgetType).Elem()).Interface()
// 	w, ok := i.(Widget)
// 	if !ok {
// 		log.Fatal("nuxui", "Create widget from widget type %T faild", widgetType)
// 	}
// 	return w
// }

// TODO mutex, or check if is in init func?
// TODO What should I do if I have already registered, but I wrote it directly in the template, and the package is not included in the code, and golang cannot find it? Displayed use nux.Use(...), automatically generated for later use
// eg: ui.RegisterWidget((*MyWidget)nil)
func RegisterWidget(widget interface{}, creator Creator) {
	// TODO Type check, judgment is not (**MyWidget)nil
	// if widget != nil {
	// 	panic("regist widget faild, widget argument should like (*MyWidget)nil")
	// }

	// if w, ok := widget.(Widget); !ok {
	// 	log.Fatal("nuxui", "%T %T is not a Widget, register faild.", w, widget)
	// }

	t := reflect.TypeOf(widget).Elem()
	widgetName := t.PkgPath() + "." + t.Name()
	if _, ok := _widgetList[widgetName]; ok {
		log.Fatal("nuxui", "Widget %s is already registed", widgetName)
	} else {
		_widgetList[widgetName] = creator
	}
}

func FindRegistedWidgetCreatorByType(widget interface{}) Creator {
	log.I("nuxui", "FindRegistedWidgetCreatorByType %T", widget)
	if widgetName, ok := widget.(string); ok {
		return FindRegistedWidgetCreatorByName(widgetName)
	}

	t := reflect.TypeOf(widget).Elem()
	widgetName := t.PkgPath() + "." + t.Name()
	return FindRegistedWidgetCreatorByName(widgetName)
}

func FindRegistedWidgetCreatorByName(name string) Creator {
	if c, ok := _widgetList[name]; ok {
		return c
	}
	log.Fatal("nuxui", "widget '%s' can not find, make sure it was registed", name)
	return nil
}

/////////////////////// mixin //////////////////////////

func CreateMixin(mixinType interface{}) Mixin {
	i := reflect.New(reflect.TypeOf(mixinType).Elem()).Interface()
	m, ok := i.(Mixin)
	if !ok {
		log.Fatal("nuxui", "Create mixin from mixin type %T faild", mixinType)
	}
	return m
}

func RegisterMixin(mixin interface{}, creator MixinCreator) {
	// if _, ok := mixin.(Mixin); !ok {
	// 	log.Fatal("nuxui", "%T is not a Mixin, register faild.", mixin)
	// }

	t := reflect.TypeOf(mixin).Elem()
	mixinName := t.PkgPath() + "." + t.Name()
	if _, ok := _mixinList[mixinName]; ok {
		log.Fatal("nuxui", "Mixin %s is already registed", mixinName)
	} else {
		_mixinList[mixinName] = creator
	}
}

func FindRegistedMixinCreatorByType(mixin interface{}) MixinCreator {
	if mixinName, ok := mixin.(string); ok {
		return FindRegistedMixinCreatorByName(mixinName)
	}

	t := reflect.TypeOf(mixin).Elem()
	mixinName := t.PkgPath() + "." + t.Name()
	return FindRegistedMixinCreatorByName(mixinName)
}

func FindRegistedMixinCreatorByName(name string) MixinCreator {
	if c, ok := _mixinList[name]; ok {
		return c
	}
	log.Fatal("nuxui", "Mixin '%s' can not find, make sure it was registed", name)
	return nil
}

func IsView(widget Widget) bool {
	_, ret := widget.(viewfuncs)

	return ret
}

// only find
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

// func FindMixin(widget Widget, id string) Mixin {
// 	for _, mixin := range Mixins(widget) {
// 		if id == mixin.ID()  {
// 			return mixin
// 		}
// 	}

// 	log.Fatal("nuxui", "the mixin id '%s' was not found in widget %T\n", id, widget)
// 	return nil
// }
