// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

func Roundf32(x float32) float32 {
	if x == 0 || x != x {
		return 0
	} else if x > 0 {
		return float32(int32(x + 0.5))
	} else if x < 0 {
		return float32(int32(x - 0.5))
	}
	return x
}

func Roundi32(x float32) int32 {
	if x == 0 || x != x {
		return 0
	} else if x > 0 {
		return int32(x + 0.5)
	} else if x < 0 {
		return int32(x - 0.5)
	}
	return int32(x)
}
