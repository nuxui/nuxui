// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/darwin"

)

func TestLog(t *testing.T) {
	darwin.NewNSWindow(800, 600)
	log.I("nuxui", "h")
}
