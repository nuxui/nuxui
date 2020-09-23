// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

// Application app
type Application interface {
	Manifest() Manifest
	MainWindow() Window
	// KeyWindow() Window
	// Windows() []Window
	// FindWindow() Window
	SendEvent(Event)
	sendEventAndWaitDone(Event)
	Terminate()
}

// App get the nux app instance
func App() Application {
	return app()
}

// Init init application
func Init(manifest string) {
	attr := ParseAttr(manifest)
	if c, ok := App().(Creating); ok {
		c.Creating(attr)
	}
}

// Run run application
func Run() {
	log.V("nuxui", "Run...")
	run()
}

func (me *application) loop() {
	// application created and started , so here is started?
	// me.Created(struct{}{})

	log.V("nuxui", "run go loop...")
	var event Event
	var wait bool
	for {
		select {
		case event = <-me.event:
			wait = false
		case event = <-me.eventWaitDone:
			wait = true
		}
		// log.V("nuxui", "select a event wait = %s", wait)

		switch event.Type() {
		case Type_WindowEvent:
			switch event.Action() {
			case Action_WindowCreated:
				log.V("nuxui", "Action_WindowCreated")
				if f, ok := event.Window().(AnyCreated); ok {
					f.Created(App().Manifest().Main())
				}
			case Action_WindowMeasured:
				log.V("nuxui", "Action_WindowMeasured width=%d, height=%d", event.Window().ContentWidth(), event.Window().ContentHeight())
				if f, ok := event.Window().(Measure); ok {
					f.Measure(event.Window().ContentWidth(), event.Window().ContentHeight())
				}

				if f, ok := event.Window().(Layout); ok {
					f.Layout(0, 0, 0, 0, event.Window().ContentWidth(), event.Window().ContentHeight())
				}

				if f, ok := event.Window().(Draw); ok {
					if canvas, err := event.Window().LockCanvas(); err == nil {
						f.Draw(canvas)
						event.Window().UnlockCanvas()
					}
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
			event.Window().handlePointerEvent(event)
		case Type_KeyEvent:
			event.Window().handleKeyEvent(event)
		case Type_TypingEvent:
			event.Window().handleTypingEvent(event)
		case Type_BackToUI:
			if f, ok := event.Data().(func()); ok {
				f()
			} else {
				log.Fatal("nuxui", "can not accept type %T for event Type_BackToUI, only supported type func()", event.Data())
			}
		case Type_AppExit:
			goto end
		}

		if wait {
			// log.V("nuxui", "wait event done")
			wait = false
			me.eventDone <- struct{}{}
		}
	}
end:
	if wait {
		wait = false
		me.eventDone <- struct{}{}
	}
	log.V("nuxui", "end of nux loop")
}

func RequestLayout(widget Widget) {
	e := &event{
		etype:  Type_WindowEvent,
		action: Action_WindowMeasured,
		window: GetWidgetWindow(widget),
	}
	go func() {
		App().sendEventAndWaitDone(e)
	}()
}

func RequestRedraw(widget Widget) {
	e := &event{
		etype:  Type_WindowEvent,
		action: Action_WindowMeasured,
		window: GetWidgetWindow(widget),
	}
	go func() {
		App().sendEventAndWaitDone(e)
	}()
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
	GetWidgetWindow(widget).requestFocus(widget)
}

func GetWidgetWindow(widget Widget) Window {
	if widget.Parent() == nil {
		if decor, ok := widget.(Decor); ok {
			return decor.Window()
		} else {
			log.Fatal("nuxui", "can not run here %T", widget)
			return nil
		}
	}
	return GetWidgetWindow(widget.Parent())
}
