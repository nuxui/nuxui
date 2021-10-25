// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"github.com/nuxui/nuxui/log"
)

// Application app
type Application interface {
	Manifest() Manifest
	MainWindow() Window
	// SendEvent(Event)
	// SendEventAndWait(Event)
	RequestLayout(Widget)
	RequestRedraw(Widget)
	Terminate()
}

// App get the nux app instance
func App() Application {
	return app()
}

// Init init application
func Init(manifest string) {
	attr := ParseAttr(manifest)
	theApp.manifest = NewManifest(attr.GetAttr("manifest", Attr{}))
	theApp.window = newWindow(Attr{})
}

// Run run application
func Run() {
	log.V("nuxui", "Run...")
	run()
}

func (me *application) handleEvent(event Event) {
	switch event.Type() {
	case Type_WindowEvent:
		switch event.Action() {
		case Action_WindowCreated:
			log.V("nuxui", "Action_WindowCreated")
			if f, ok := event.Window().(OnCreate); ok {
				f.OnCreate()
			}
			log.V("nuxui", "Action_WindowCreated end")
		case Action_WindowMeasured:
			log.V("nuxui", "Action_WindowMeasured width=%d, height=%d", event.Window().ContentWidth(), event.Window().ContentHeight())
			if f, ok := event.Window().(Measure); ok {
				f.Measure(event.Window().ContentWidth(), event.Window().ContentHeight())
			}

			if f, ok := event.Window().(Layout); ok {
				f.Layout(0, 0, 0, 0, event.Window().ContentWidth(), event.Window().ContentHeight())
			}
		case Action_WindowDraw:
			// log.V("nuxui", "Action_WindowDraw")
			if f, ok := event.Window().(Draw); ok {
				if canvas, err := event.Window().LockCanvas(); err == nil {
					f.Draw(canvas)
					event.Window().UnlockCanvas()
				}
			}
		}
	// App().MainWindow()
	// App().Manifest().Root()
	case Type_PointerEvent:
		event.Window().handlePointerEvent(event.(PointerEvent))
	case Type_ScrollEvent:
		event.Window().handleScrollEvent(event.(ScrollEvent))
	case Type_KeyEvent:
		event.Window().handleKeyEvent(event.(KeyEvent))
	case Type_TypeEvent:
		event.Window().handleTypeEvent(event.(TypeEvent))
	case Type_BackToUI:
		if f, ok := event.Data().(func()); ok {
			f()
		} else {
			log.Fatal("nuxui", "can not accept type %T for event Type_BackToUI, only supported type func()", event.Data())
		}
	case Type_AppExit:
		// goto end
	}

}

func (me *application) RequestLayout(widget Widget) {
	// TODO:: schedule layout
	if w := GetWidgetWindow(widget); w != nil {
		e := &event{
			etype:  Type_WindowEvent,
			action: Action_WindowMeasured,
			window: w,
		}
		me.handleEvent(e)
	}
}

func RequestLayout(widget Widget) {
	theApp.RequestLayout(widget)
}

func RequestRedraw(widget Widget) {
	theApp.RequestRedraw(widget)
}

func StartTextInput() {
	startTextInput()
}

func StopTextInput() {
	stopTextInput()
}

func SetTextInputRect(x, y, w, h float32) {
	setTextInputRect(x, y, w, h)
}

func RequestFocus(widget Widget) {
	if w := GetWidgetWindow(widget); w != nil {
		w.requestFocus(widget)
	}
}

var decorWindowList map[Widget]Window = map[Widget]Window{}

// GetWidgetWindow get the window of widget  return can be nil if widget not attach to window
func GetWidgetWindow(widget Widget) Window {

	return theApp.window
}

// if widget.Parent() == nil {
// 	if util.GetTypeName(widget) == "github.com/nuxui/nuxui/ui.layer" {
// 		if w, ok := decorWindowList[widget]; ok {
// 			return w
// 		}
// 	}

// 	// log.Fatal("nuxui", "can not run here %T", widget)
// 	// widget not attached to window yet
// 	return nil
// }
// return GetWidgetWindow(widget.Parent())

func RunOnUI(callback func()) {
	runOnUI(callback)
}
