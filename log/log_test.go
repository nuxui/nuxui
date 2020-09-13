// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	V("log", "TestLog")
	V("log", "TestLog %d", 1)
	V("log", "TestLog %d %s", 1, "hello")
	Fatal("log", "TestLog %d %s %T", 1, "hello", TestLog)
}
