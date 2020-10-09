// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

const (
	permission_granted = 0 + iota
	permission_denied
	permission_not_allowed
)

func RequestPermissions(permission []string, callback func(map[string]int)) {

}

func CheckPermissions(permission []string) map[string]int {
	return nil
}
