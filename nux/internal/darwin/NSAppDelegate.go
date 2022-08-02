// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package darwin

type NSApplicationTerminateReply int

const (
	NSTerminateCancel = 0
	NSTerminateNow    = 1
	NSTerminateLater  = 2
)

// https://developer.apple.com/documentation/appkit/nsapplicationdelegate?language=objc
type NSApplicationDelegate interface {
	// Launching Applications
	ApplicationWillFinishLaunching(notification NSNotification)
	ApplicationDidFinishLaunching(notification NSNotification)

	// Managing Active Status
	ApplicationWillBecomeActive(notification NSNotification)
	ApplicationDidBecomeActive(notification NSNotification)
	ApplicationWillResignActive(notification NSNotification)
	ApplicationDidResignActive(notification NSNotification)

	// Terminating Applications
	ApplicationShouldTerminate(sender NSApplication) NSApplicationTerminateReply
	ApplicationShouldTerminateAfterLastWindowClosed(sender NSApplication) bool
	ApplicationWillTerminate(notification NSNotification)

	// Hiding Applications
	ApplicationWillHide(notification NSNotification)
	ApplicationDidHide(notification NSNotification)
	ApplicationWillUnhide(notification NSNotification)
	ApplicationDidUnhide(notification NSNotification)

	// Managing Windows
}
