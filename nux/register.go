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

var _widgetList = make(map[string]Creator)
var _mixinList = make(map[string]MixinCreator)

func CreateWidget(widgetType interface{}) Widget {
	// TODO:: check widgetType is struct and (*Widget)nil
	// return reflect.New(reflect.TypeOf(widgetType)).Elem().Interface()
	i := reflect.New(reflect.TypeOf(widgetType).Elem()).Interface()
	w, ok := i.(Widget)
	if !ok {
		log.Fatal("nux", fmt.Sprintf("Create widget from widget type %T faild", widgetType))
	}
	return w
}

// TODO mutex, or check if is in init func?
// TODO What should I do if I have already registered, but I wrote it directly in the template, and the package is not included in the code, and golang cannot find it? Displayed use nux.Use(...), automatically generated for later use
// eg: ui.RegisterWidget((*MyWidget)nil)
func RegisterWidget(widget interface{}, creator Creator) {
	// TODO Type check, judgment is not (**MyWidget)nil
	// if widget != nil {
	// 	panic("regist widget faild, widget argument should like (*MyWidget)nil")
	// }

	// if w, ok := widget.(Widget); !ok {
	// 	log.Fatal("nux", fmt.Sprintf("%T %T is not a Widget, register faild.", w, widget))
	// }

	t := reflect.TypeOf(widget).Elem()
	widgetName := t.PkgPath() + "." + t.Name()
	if _, ok := _widgetList[widgetName]; ok {
		log.Fatal("nux", fmt.Sprintf("Widget %s is already registed", widgetName))
	} else {
		_widgetList[widgetName] = creator
	}
}

func FindRegistedWidgetCreatorByType(widget interface{}) Creator {
	log.I("nux", "FindRegistedWidgetCreatorByType %T", widget)
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
	log.Fatal("nux", fmt.Sprintf("widget '%s' can not find, make sure it was registed", name))
	return nil
}

/////////////////////// mixin //////////////////////////

func CreateMixin(mixinType interface{}) Mixin {
	i := reflect.New(reflect.TypeOf(mixinType).Elem()).Interface()
	m, ok := i.(Mixin)
	if !ok {
		log.Fatal("nux", fmt.Sprintf("Create mixin from mixin type %T faild", mixinType))
	}
	return m
}

func RegisterMixin(mixin interface{}, creator MixinCreator) {
	// if _, ok := mixin.(Mixin); !ok {
	// 	log.Fatal("nux", fmt.Sprintf("%T is not a Mixin, register faild.", mixin))
	// }

	t := reflect.TypeOf(mixin).Elem()
	mixinName := t.PkgPath() + "." + t.Name()
	if _, ok := _mixinList[mixinName]; ok {
		log.Fatal("nux", fmt.Sprintf("Mixin %s is already registed", mixinName))
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
	log.Fatal("nux", fmt.Sprintf("Mixin '%s' can not find, make sure it was registed", name))
	return nil
}

func RenderWidget(widget Widget) Widget {
	if reflect.TypeOf(widget).Kind() != reflect.Ptr {
		log.Fatal("nux", fmt.Sprintf("Widget %T should a pointer. eg: &MyWidget{}", widget)) // TODO:: tips
	}

	if IsComponent(widget) {
		var content Widget

		if fn, ok := widget.(Render); ok { // RenderRender
			executeCreating(widget, Attr{})
			content = RenderRender(fn.Render()) // TODO:: &Header{}. template()
		} else if fn, ok := widget.(Template); ok {
			attrs := ParseAttr(fn.Template())
			executeCreating(widget, attrs)
			content = RenderTemplate(widget, attrs)
		} else {
			log.Fatal("nux", fmt.Sprintf("%T is not a component, the code should not run here", widget))
		}

		if c, ok := widget.(Created); ok {
			c.Created(content)
		}

		return NewComponent(widget, content)
	}

	return widget
}

func RenderRender(widget Widget) Widget {
	w := RenderWidget(widget)
	if p, ok := w.(Parent); ok {
		for _, child := range p.Children() {
			if compt, ok := child.(Component); ok {
				child = compt.Content()
			}

			r := RenderRender(child)
			p.ReplaceChild(child, r)
		}
	}
	return w
}

func RenderTemplate(component Widget, attrs Attr) Widget {
	layout := attrs.GetAttr("layout", make(Attr, 0))
	return renderLayout(component, layout)
}

func renderLayout(component Widget, layout Attr) (ret Widget) {
	if widgetNode := layout.GetString("widget", ""); widgetNode != "" {
		// TODO:: copyAttribute(layout)
		// w0 := reflect.New(reflect.TypeOf(widgetType).Elem()).Elem().Interface()
		// w := reflect.New(widgetType).Elem().Interface()
		// w := info.Creator(layout)
		// var w interface{} = &w0
		// w := CreateWidget(widgetType)
		w := FindRegistedWidgetCreatorByName(widgetNode)()

		mixins := layout.GetStringArray("mixins", []string{})
		for _, name := range mixins {
			mixin := FindRegistedMixinCreatorByName(name)()
			if h, ok := mixin.(ComptHelper); ok {
				h.AssignComponent(component)
			}
			AddMixins(w, mixin)
		}

		// TODO:: The execution order of the life cycle? From parent to child, because the measure draw layout is like this.
		executeCreating(w, layout)

		ret = RenderWidget(w)

		if childrenNode, ok := layout["children"]; ok {
			if children, ok := childrenNode.([]interface{}); ok {
				if p, ok := ret.(Parent); ok {
					for _, childNode := range children {
						if child, ok := childNode.(Attr); ok {
							childWidget := renderLayout(component, child)
							p.AddChild(childWidget)
						}
					}
				} else {
					log.Fatal("nux", fmt.Sprintf("%T is not a WidgetGroup", p))
				}
			} else {
				log.Fatal("nux", "unknow type for Children.")
			}
		}

		executeCreated(ret)

	} else {
		log.Fatal("nux", `must specified "widget"`)
	}
	return ret
}

func IsWidget(widget Widget) bool {
	_, hasLayout := widget.(Layout)
	_, hasDraw := widget.(Draw)
	_, hasTemplate := widget.(Template)
	_, hasRender := widget.(Render)

	return hasLayout || !hasDraw || hasTemplate || hasRender
}

func IsComponent(widget Widget) bool {
	_, hasLayout := widget.(Layout)
	_, hasDraw := widget.(Draw)
	_, hasTemplate := widget.(Template)
	_, hasRender := widget.(Render)

	return !hasLayout && !hasDraw && (hasTemplate || hasRender)
}

func IsView(widget Widget) bool {
	_, hasLayout := widget.(Layout)
	_, hasDraw := widget.(Draw)
	_, hasTemplate := widget.(Template)
	_, hasRender := widget.(Render)

	return !hasTemplate && !hasRender && (hasLayout || hasDraw)
}

// TODO measure root widget, should move to app package
func MeasureWidget(widget Widget) {
	if m, ok := widget.(Measure); ok {
		// TODO:: dynamic obtain size
		if s, ok := widget.(Size); ok {
			ms := s.MeasuredSize()
			m.Measure(MeasureSpec(ms.Width, Pixel), MeasureSpec(ms.Height, Pixel))
		}
	}
}

func LayoutWidget(widget Widget) {
	if layout, ok := widget.(Layout); ok {
		if s, ok := widget.(Size); ok {
			layout.Layout(0, 0, 0, 0, s.MeasuredSize().Width, s.MeasuredSize().Height)
		}
	}
}

func DrawWidget(canvas Canvas, widget Widget) {
	if s, ok := widget.(Size); ok {
		if d, ok := widget.(Draw); ok {
			ms := s.MeasuredSize()
			canvas.Save()
			canvas.Translate(ms.Position.Left, ms.Position.Top)
			canvas.ClipRect(0, 0, ms.Width, ms.Height)
			d.Draw(canvas)
			canvas.Restore()
		}
	}
}

func DrawWidget2(canvas Canvas, widget Widget) {
	drawWidgetById0(canvas, widget, "title")
}

func drawWidgetById0(canvas Canvas, widget Widget, id string) {
	if strings.Compare("title", widget.ID()) == 0 {
		if s, ok := widget.(Size); ok {
			if d, ok := widget.(Draw); ok {
				ms := s.MeasuredSize()
				canvas.Save()
				canvas.Translate(ms.Position.Left, ms.Position.Top)
				canvas.ClipRect(0, 0, ms.Width, ms.Height)
				if strings.Compare("title", widget.ID()) == 0 {
					d.Draw(canvas)
				}
				canvas.Restore()
			}
		}
	} else {
		if p, ok := widget.(Parent); ok {

			if s, ok := widget.(Size); ok {
				ms := s.MeasuredSize()
				canvas.Save()
				canvas.Translate(ms.Position.Left, ms.Position.Top)

				for _, child := range p.Children() {
					if compt, ok := child.(Component); ok {
						child = compt.Content()
					}

					drawWidgetById0(canvas, child, id)
				}

				canvas.Restore()
			}

		}
	}
}

// func GetComponentRootWidget(widget Widget) Widget {
// 	return nil
// }

// only find
func Find(widget Widget, id string) Widget {
	if id == "" {
		log.Fatal("nux", fmt.Sprintf("the widget %T id must be specified", widget))
	}

	w := find(widget, id)
	if w == nil {
		log.Fatal("nux", fmt.Sprintf("the id '%s' was not found in widget %T\n", id, widget))
	}

	return w
}

func find(widget Widget, id string) Widget {
	if strings.Compare(id, widget.ID()) == 0 {
		return widget
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
// 		if strings.Compare(id, mixin.ID()) == 0 {
// 			return mixin
// 		}
// 	}

// 	log.Fatal("nux", fmt.Sprintf("the mixin id '%s' was not found in widget %T\n", id, widget))
// 	return nil
// }
