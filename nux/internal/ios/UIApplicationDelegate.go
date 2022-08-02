// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ios

// https://developer.apple.com/documentation/uikit/uiapplicationdelegate?language=objc
type UIApplicationDelegate interface {
	// Initializing the App
	WillFinishLaunchingWithOptions(launchOptions map[string]uintptr) bool
	DidFinishLaunchingWithOptions(launchOptions map[string]uintptr) bool
}
