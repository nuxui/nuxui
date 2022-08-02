// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && ios

package nux

import ()

func showViewFileDialog(dialog *viewFileDialog) {
}

func showPickFileDialog(dialog *pickFileDialog) (ok bool, paths []string) {
	return
}

func showSaveFileDialog(dialog *saveFileDialog) (ok bool, saveName string) {
	return
}
