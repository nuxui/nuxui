// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

func TestDimen(t *testing.T) {
	defer log.Close()
	log.V("test", "%s, %s, %s", nux.SDimen("10.02px"), nux.ADimen(10.02, nux.Pixel), nux.Dimen(10))
	log.V("test", "%s, %s, %s", nux.SDimen("-10.02px"), nux.ADimen(-10.02, nux.Pixel), nux.Dimen(-10))
	log.V("test", "%s, %s", nux.SDimen("10.02%"), nux.ADimen(10.02, nux.Percent))
	log.V("test", "%s, %s", nux.SDimen("-10.02%"), nux.ADimen(-10.02, nux.Percent))
	log.V("test", "%s, %s", nux.SDimen("10.02wt"), nux.ADimen(10.02, nux.Weight))
	// log.V("test", "%s, %s", nux.SDimen("10.02wt"), nux.ADimen(-10.02, nux.Weight))
	log.V("test", "%s, %s", nux.SDimen("3.25:4"), nux.ADimen(3.25/4, nux.Ratio))
	// log.V("test", "%s, %s", nux.SDimen("3.25:4"), nux.ADimen(-3.25/4, nux.Ratio))
}
