// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"runtime"
	"time"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/darwin"
)

type nativeWindow struct {
	ptr darwin.NSWindow
}

func newNativeWindow(attr Attr) *nativeWindow {
	darwin.SetWindowEventHandler(nativeWindowEventHandler)

	width, height := measureWindowSize(attr.GetDimen("width", "50%"), attr.GetDimen("height", "50%"))
	me := &nativeWindow{
		ptr: darwin.NewNSWindow(width, height),
	}
	me.SetTitle(attr.GetString("title", ""))

	runtime.SetFinalizer(me, freeWindow)
	return me
}

func freeWindow(me *nativeWindow) {
	darwin.NSObject_release(uintptr(me.ptr))
}

func (me *nativeWindow) Center() {
	me.ptr.Center()
}

func (me *nativeWindow) Show() {
	me.ptr.MakeKeyAndOrderFront()
}

func (me *nativeWindow) ContentSize() (width, height int32) {
	return me.ptr.ContentSize()
}

func (me *nativeWindow) Title() string {
	return me.ptr.Title()
}

func (me *nativeWindow) SetTitle(title string) {
	me.ptr.SetTitle(title)
}

func (me *nativeWindow) lockCanvas() Canvas {
	return newCanvas(darwin.CGCurrentContext())
}

func (me *nativeWindow) unlockCanvas(canvas Canvas) {
	// canvas.Flush()
}

func (me *nativeWindow) draw(canvas Canvas, decor Widget) {
	if decor != nil {
		if f, ok := decor.(Draw); ok {
			_, h := me.ContentSize()
			canvas.Save()
			canvas.Translate(0, float32(h))
			canvas.Scale(1, -1)
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
		}
	}
}

func nativeWindowEventHandler(event any) bool {
	switch e := event.(type) {
	case darwin.NSEvent:
		return handleNSEvent(e)
	case *darwin.TypingEvent:
		return handleTypingEvent(e)
	case *darwin.WindowEvent:
		switch e.Type {
		case darwin.Event_WindowDidResize:
			theApp.window.resize()
		case darwin.Event_WindowDrawRect:
			theApp.window.draw()
		}
	}

	return false
}

func handleNSEvent(event darwin.NSEvent) bool {
	switch event.Type() {
	case darwin.NSEventTypeMouseEntered, darwin.NSEventTypeMouseExited,
		darwin.NSEventTypeLeftMouseDown, darwin.NSEventTypeLeftMouseUp,
		darwin.NSEventTypeRightMouseDown, darwin.NSEventTypeRightMouseUp,
		darwin.NSEventTypeMouseMoved, darwin.NSEventTypeLeftMouseDragged,
		darwin.NSEventTypeRightMouseDragged, darwin.NSEventTypeOtherMouseDown,
		darwin.NSEventTypeOtherMouseUp, darwin.NSEventTypeOtherMouseDragged:
		return handlePointerEvent(event)
	case darwin.NSEventTypePressure:
		// TODO pressure
	case darwin.NSEventTypeKeyDown, darwin.NSEventTypeKeyUp, darwin.NSEventTypeFlagsChanged:
		return handleKeyEvent(event)
	case darwin.NSEventTypeScrollWheel:
		return handleScrollEvent2(event)
	case darwin.NSEventTypeAppKitDefined:
		log.V("nuxui", "darwin.NSEventTypeAppKitDefined subtype=%d", event.Subtype())
	case darwin.NSEventTypeSystemDefined:
		log.V("nuxui", "darwin.NSEventTypeSystemDefined subtype=%d", event.Subtype())
	case darwin.NSEventTypeApplicationDefined:
		log.V("nuxui", "darwin.NSEventTypeApplicationDefined subtype=%d", event.Subtype())
	case darwin.NSEventTypePeriodic:
		log.V("nuxui", "darwin.NSEventTypePeriodic subtype=%d", event.Subtype())
		// if event.HasPreciseScrollingDeltas() {
		// }
	default:
		// x, y := event.LocationInWindow()
		// log.V("nuxui", "sendEvent type=%d, x=%f, y=%f, window=%d", event.Type(), x, y, event.Window())
	}
	// x, y := event.LocationInWindow()
	// log.V("nuxui", "sendEvent type=%d, x=%f, y=%f", event.Type(), x, y)
	return false
}

func handleTypingEvent(tevent *darwin.TypingEvent) bool {
	// 0 = Action_Input, 1 = Action_Preedit,
	act := Action_Input
	if tevent.Action == 1 {
		act = Action_Preedit
	}

	e := &typingEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_TypingEvent,
			action: act,
		},
		text:     tevent.Text,
		location: tevent.Location,
	}

	return App().MainWindow().handleTypingEvent(e)
}

func handleScrollEvent2(nsevent darwin.NSEvent) bool {
	x, y := nsevent.LocationInWindow()
	e := &scrollEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_ScrollEvent,
			action: Action_Scroll,
		},
		x:       x,
		y:       y,
		scrollX: nsevent.ScrollingDeltaX(),
		scrollY: nsevent.ScrollingDeltaY(),
	}

	if nsevent.HasPreciseScrollingDeltas() {
		e.scrollX *= 0.1
		e.scrollY *= 0.1
	}

	return App().MainWindow().handleScrollEvent(e)
}

var lastMouseEvent map[PointerButton]PointerEvent = map[PointerButton]PointerEvent{}

func handlePointerEvent(nsevent darwin.NSEvent) bool {
	x, y := nsevent.LocationInWindow()
	etype := nsevent.Type()

	e := &pointerEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_PointerEvent,
			action: Action_None,
		},
		pointer: 0,
		button:  ButtonNone,
		kind:    Kind_Mouse,
		x:       x,
		y:       y,
	}

	switch etype {
	case darwin.NSEventTypeMouseMoved:
		e.action = Action_Hover
		e.button = ButtonNone
		e.pointer = 0
	case darwin.NSEventTypeLeftMouseDown:
		e.action = Action_Down
		e.button = ButtonPrimary
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case darwin.NSEventTypeLeftMouseUp:
		e.action = Action_Up
		e.button = ButtonPrimary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	case darwin.NSEventTypeLeftMouseDragged:
		e.action = Action_Drag
		e.button = ButtonPrimary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case darwin.NSEventTypeRightMouseDown:
		e.action = Action_Down
		e.button = ButtonSecondary
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case darwin.NSEventTypeRightMouseUp:
		e.action = Action_Up
		e.button = ButtonSecondary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	case darwin.NSEventTypeRightMouseDragged:
		e.action = Action_Drag
		e.button = ButtonSecondary
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	case darwin.NSEventTypeOtherMouseDown:
		e.action = Action_Down
		buttonNumber := nsevent.ButtonNumber()
		switch buttonNumber {
		case 2:
			e.button = ButtonMiddle
		case 3:
			e.button = ButtonX1
		case 4:
			e.button = ButtonX2
		default:
			e.button = PointerButton(buttonNumber)
		}
		e.pointer = time.Now().UnixNano()
		lastMouseEvent[e.button] = e
	case darwin.NSEventTypeOtherMouseUp:
		e.action = Action_Up
		buttonNumber := nsevent.ButtonNumber()
		switch buttonNumber {
		case 2:
			e.button = ButtonMiddle
		case 3:
			e.button = ButtonX1
		case 4:
			e.button = ButtonX2
		default:
			e.button = PointerButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
			delete(lastMouseEvent, e.button)
		}
	case darwin.NSEventTypeOtherMouseDragged:
		e.action = Action_Drag
		buttonNumber := nsevent.ButtonNumber()
		switch buttonNumber {
		case 2:
			e.button = ButtonMiddle
		case 3:
			e.button = ButtonX1
		case 4:
			e.button = ButtonX2
		default:
			e.button = PointerButton(buttonNumber)
		}
		if v, ok := lastMouseEvent[e.button]; ok {
			e.pointer = v.Pointer()
		}
	}

	return App().MainWindow().handlePointerEvent(e)
}

var lastModifierKeyEvent map[KeyCode]bool = map[KeyCode]bool{}

func handleKeyEvent(nsevent darwin.NSEvent) bool {
	keyCode := KeyCode(nsevent.KeyCode())

	if keyCode == 0x36 { // kVK_Command
		keyCode = Key_Command
	}
	e := &keyEvent{
		event: event{
			window: theApp.MainWindow(),
			time:   time.Now(),
			etype:  Type_KeyEvent,
			action: Action_None,
		},
		keyCode:       keyCode,
		modifierFlags: convertModifierFlags(nsevent.ModifierFlags()),
	}

	switch nsevent.Type() {
	case darwin.NSEventTypeKeyDown:
		e.action = Action_Down
		e.repeat = nsevent.IsARepeat()
		e.keyRune = nsevent.Characters()
	case darwin.NSEventTypeKeyUp:
		e.action = Action_Up
		e.repeat = nsevent.IsARepeat()
		e.keyRune = nsevent.Characters()
	case darwin.NSEventTypeFlagsChanged:
		if down, ok := lastModifierKeyEvent[e.keyCode]; ok && down {
			lastModifierKeyEvent[e.keyCode] = false
			e.action = Action_Up
		} else {
			lastModifierKeyEvent[e.keyCode] = true
			e.action = Action_Down
		}
	}

	return App().MainWindow().handleKeyEvent(e)
}

func convertModifierFlags(flags darwin.NSEventModifierFlags) uint32 {
	var mods uint32 = 0
	if flags&darwin.NSEventModifierFlagShift == darwin.NSEventModifierFlagShift {
		mods |= Mod_Shift
	}
	if flags&darwin.NSEventModifierFlagControl == darwin.NSEventModifierFlagControl {
		mods |= Mod_Control
	}
	if flags&darwin.NSEventModifierFlagOption == darwin.NSEventModifierFlagOption {
		mods |= Mod_Alt
	}
	if flags&darwin.NSEventModifierFlagCommand == darwin.NSEventModifierFlagCommand {
		mods |= Mod_Super
	}
	if flags&darwin.NSEventModifierFlagCapsLock == darwin.NSEventModifierFlagCapsLock {
		mods |= Mod_CapsLock
	}
	return mods
}
