// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	// "nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/linux"
	"nuxui.org/nuxui/nux/internal/linux/gtk3"
)

func showViewFileDialog(dialog *viewFileDialog) {
	if dialog.directory != "" {
		linux.System("xdg-open " + dialog.directory)
	}
}

// https://stackoverflow.com/questions/45153305/gtk-filechooserdialog-select-files-and-folders-vala
func showPickFileDialog(dialog *pickFileDialog) (ok bool, paths []string) {
	if !gtk3.InitCheck() {
		return false, nil
	}

	action := gtk3.FILE_CHOOSER_ACTION_OPEN
	if dialog.allowsChooseFolders {
		action = gtk3.FILE_CHOOSER_ACTION_SELECT_FOLDER
	}

	chooser := gtk3.FileChooserDialogNew("", 0, action)

	if dialog.directory != "" {
		chooser.SetCurrentFolder(dialog.directory)
	}

	if !dialog.allowsChooseFolders {
		for k, v := range dialog.filters {
			filter := gtk3.FileFilterNew()
			filter.SetName(k)
			for _, f := range v {
				filter.AddPattern("*." + f)
			}
			chooser.AddFilter(filter)
		}
	}

	chooser.SetSelectMultiple(dialog.allowsMultipleSelection)
	chooser.SetCreateFolders(true)

	defer chooser.Close()
	return chooser.Run() == gtk3.RESPONSE_ACCEPT, chooser.GetFilenames()
}

func showSaveFileDialog(dialog *saveFileDialog) (ok bool, saveName string) {
	if !gtk3.InitCheck() {
		return false, ""
	}
	chooser := gtk3.FileChooserDialogNew("", 0, gtk3.FILE_CHOOSER_ACTION_SAVE)

	if dialog.directory != "" {
		chooser.SetCurrentFolder(dialog.directory)
	}

	if dialog.saveName != "" {
		chooser.SetCurrentName(dialog.saveName)
	}
	chooser.SetCreateFolders(true)
	chooser.SetDoOverwriteConfirmation(true)

	defer chooser.Close()
	return chooser.Run() == gtk3.RESPONSE_ACCEPT, chooser.GetFilename()
}
