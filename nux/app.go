// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"math"
	"runtime"

	"nuxui.org/nuxui/log"
)

type Application interface {
	MainWindow() Window
	Terminate()
	Run()
}

type application struct {
	native           *nativeApp
	mainThreadID     uint64
	mainWindow       Window
	createMainWindow func()
	invalidateSignal chan *Rect
	measureSignal    chan Widget
	windowPrepared   chan struct{}
}

const (
	invalidateSignalSize = 50
)

var theApp = &application{
	native:           createNativeApp(),
	invalidateSignal: make(chan *Rect, invalidateSignalSize),
	measureSignal:    make(chan Widget, invalidateSignalSize),
	windowPrepared:   make(chan struct{}, 1),
}

func init() {
	runtime.LockOSThread()
	theApp.mainThreadID = currentThreadID()
	timerLoopInstance.init()
}

func App() Application {
	return theApp
}

func Run(createMainWindow func()) {
	if createMainWindow == nil {
		log.Fatal("nuxui", "createMainWindow can not be nil")
	}
	if theApp.createMainWindow != nil {
		log.Fatal("nuxui", "do not call Run again")
	}

	theApp.createMainWindow = createMainWindow

	// in android, the thread id of init func is different with main func
	// but it is still run in order, has not effect
	if runtime.GOOS == "android" {
		theApp.mainThreadID = currentThreadID()
	}

	if tid := currentThreadID(); tid != theApp.mainThreadID {
		log.Fatal("nuxui", "main called on thread %d, but init run on %d", tid, theApp.mainThreadID)
	}

	defer runtime.UnlockOSThread()

	go refreshLoop()

	theApp.Run()
}

func (me *application) Run() {
	me.native.run()
}

// not nil
func (me *application) MainWindow() Window {
	if me.mainWindow == nil {
		log.Fatal("nuxui", "main window not created yet")
	}
	return me.mainWindow
}

func (me *application) Terminate() {
	me.native.terminate()
}

func (me *application) onDidFinishLaunch() {
}

func createNativeApp() *nativeApp {
	return createNativeApp_()
}

func RunOnUI(callback func()) {
	runOnUI(callback)
}

func refreshLoop() {
	var i, l int
	var rect *Rect
	var dirtRect *Rect = &Rect{}

	<-theApp.windowPrepared

	for {
		select {
		case rect = <-theApp.invalidateSignal:
			// time.Sleep(16 * time.Millisecond)
			dirtRect.X = math.MaxInt32
			dirtRect.Y = math.MaxInt32
			dirtRect.Width = 0
			dirtRect.Height = 0
			dirtRect.Expand(rect)
			l = len(theApp.invalidateSignal)
			for i = 0; i != l; i++ {
				rect = <-theApp.invalidateSignal
				dirtRect.Expand(rect)
			}
			invalidateRectAsync(dirtRect)
		case widget := <-theApp.measureSignal:
			// time.Sleep(16 * time.Millisecond)
			l = len(theApp.measureSignal)
			for i = 0; i != l; i++ {
				// TODO:: find outermost widget
				<-theApp.measureSignal
			}
			requestLayoutAsync(widget)
		}

		// time.Sleep(16 * time.Millisecond)
	}
}

func RequestLayout(widget Widget) {
	go func() {
		theApp.measureSignal <- widget
	}()
}

func requestLayoutAsync(widget Widget) {
	runOnUI(func() {
		theApp.mainWindow.resize()
	})
}

func RequestRedraw(widget Widget) {
	if s, ok := widget.(Size); ok {
		f := s.Frame()
		InvalidateRect(f.X, f.Y, f.Width, f.Height)
	}
}

func InvalidateRect(x, y, width, height int32) {
	var rect *Rect
	var dirtRect *Rect = &Rect{x, y, width, height}

	if l := len(theApp.invalidateSignal); l >= invalidateSignalSize {
		for i := 0; i != l-1; i++ {
			rect = <-theApp.invalidateSignal
			dirtRect.Expand(rect)
		}
	}
	theApp.invalidateSignal <- dirtRect
}

func invalidateRectAsync(dirtRect *Rect) {
	invalidateRectAsync_(dirtRect)
}

func StartTextInput() {
	startTextInput()
}

func StopTextInput() {
	stopTextInput()
}

func SetTextInputRect(x, y, width, height float32) {
	DebugCheckMainThread()
	// setTextInputRect(x, y, width, height)
}

func RequestFocus(widget Widget) {
	DebugCheckMainThread()
	theApp.MainWindow().requestFocus(widget)
}

func IsMainThread() bool {
	return theApp.mainThreadID == currentThreadID()
}

func currentThreadID() uint64 {
	return currentThreadID_()
}

func ScreenSize() (width, height int32) {
	return screenSize()
}
