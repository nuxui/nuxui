// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

// Application app
type Application interface {
	Manifest() Manifest
	MainWindow() Window
	SendEvent(Event)
	sendEventAndWaitDone(Event)
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
	log.V("nux", "Run...")
	run()
}

func (me *application) loop() {
	// application created and started , so here is started?
	// me.Created(struct{}{})

	log.V("nux", "run go loop...")
	var event Event
	var wait bool
	for {
		select {
		case event = <-me.event:
			wait = false
		case event = <-me.eventWaitDone:
			wait = true
		}

		// log.V("nux", "select a event wait = %s", wait)

		switch event.Type() {
		case Type_WindowEvent:
			e := event.(WindowEvent)
			switch e.Action() {
			case Action_WindowCreated:
				log.V("nux", "Action_WindowCreated")
				if f, ok := e.Window().(AnyCreated); ok {
					f.Created(App().Manifest().Main())
				}
			case Action_WindowMeasured:
				log.V("nux", "Action_WindowMeasured width=%d, height=%d", e.Window().ContentWidth(), e.Window().ContentHeight())
				if f, ok := e.Window().(Measure); ok {
					f.Measure(e.Window().ContentWidth(), e.Window().ContentHeight())
				}

				if f, ok := e.Window().(Layout); ok {
					f.Layout(0, 0, 0, 0, e.Window().ContentWidth(), e.Window().ContentHeight())
				}

				if f, ok := e.Window().(Draw); ok {
					if canvas, err := e.Window().LockCanvas(); err == nil {
						f.Draw(canvas)
						e.Window().UnlockCanvas()
					}
				}
			case Action_WindowDraw:
				log.V("nux", "Action_WindowDraw")
				if f, ok := e.Window().(Draw); ok {
					if canvas, err := e.Window().LockCanvas(); err == nil {
						f.Draw(canvas)
						e.Window().UnlockCanvas()
					}
				}
			}
		// App().MainWindow()
		// App().Manifest().Root()
		case Type_InputEvent:
			log.V("nux", "Type_InputEvent")
			if f, ok := App().MainWindow().(Draw); ok {
				if canvas, err := App().MainWindow().LockCanvas(); err == nil {
					f.Draw(canvas)
					App().MainWindow().UnlockCanvas()
				}
			}
		case Type_AppExit:
			goto end
		}

		if wait {
			// log.V("nux", "wait event done")
			me.eventDone <- struct{}{}
			wait = false
		}
	}
end:
	if wait {
		me.eventDone <- struct{}{}
		wait = false
	}
	log.V("nux", "end of nux loop")
}

func RequestLayout(widget Widget) {
}

func RequestRedraw(widget Widget) {
}

func StartTextInput() {
	startTextInput()
}
