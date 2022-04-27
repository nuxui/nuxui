// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	"github.com/nuxui/nuxui/log"
)

func TestLog(t *testing.T) {
	log.V("log", "TestLog")
	log.V("log", "TestLog %d", 1)
	log.V("log", "TestLog %d %s", 1, "hello")
}
