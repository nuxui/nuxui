// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import (
	"nuxui.org/nuxui/log"
	"os"
)

type Response int

const (
	ResponseCancel Response = 0
	ResponseOK     Response = 1
)

// ----------- ViewFileDialog -------------------
func ViewFileDialog() *viewFileDialog {
	return &viewFileDialog{}
}

type viewFileDialog struct {
	directory       string
	activeFileNames []string
}

func (me *viewFileDialog) SetDirectory(directory string) *viewFileDialog {
	// TODO:: check is dir
	me.directory = directory
	return me
}

func (me *viewFileDialog) SetActiveFileNames(activeFileNames []string) *viewFileDialog {
	// TODO:: check is in dir ?
	me.activeFileNames = activeFileNames
	return me
}

func (me *viewFileDialog) Show() {
	if !_fileCheck(me.directory, true) {
		return
	}

	showViewFileDialog(me)
}

func _fileCheck(fileName string, isdir bool) bool {
	if fileName == "" {
		return true
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.E("nuxui", "%s", err.Error())
		return false
	}

	if isdir && !fileInfo.IsDir() {
		log.E("nuxui", "%s is not a directory", fileName)
		return false
	}

	if !isdir && fileInfo.IsDir() {
		log.E("nuxui", "%s is not a file", fileName)
		return false
	}
	return true
}

// ------------- PickFileDialog -----------------
func PickFileDialog() *pickFileDialog {
	return &pickFileDialog{}
}

// pick single or multi files, folders
// both pick files and folders can not implement at windows and linux
type pickFileDialog struct {
	directory               string
	filters                 map[string][]string
	allowsChooseFiles       bool
	allowsChooseFolders     bool
	allowsMultipleSelection bool
	allowsCreateFolders     bool
}

func (me *pickFileDialog) SetDirectory(directory string) *pickFileDialog {
	me.directory = directory
	return me
}

func (me *pickFileDialog) SetExtensionFilters(filters map[string][]string) *pickFileDialog {
	me.filters = filters
	return me
}

func (me *pickFileDialog) AllowsMultipleSelection() *pickFileDialog {
	me.allowsMultipleSelection = true
	return me
}

func (me *pickFileDialog) AllowsChooseFiles() *pickFileDialog {
	me.allowsChooseFiles = true
	return me
}

func (me *pickFileDialog) AllowsChooseFolders() *pickFileDialog {
	me.allowsChooseFolders = true
	return me
}

func (me *pickFileDialog) AllowsCreateFolders() *pickFileDialog {
	me.allowsCreateFolders = true
	return me
}

func (me *pickFileDialog) ShowModal(callback func(ok bool, results []string)) {
	if !_fileCheck(me.directory, true) {
		return
	}

	ok, paths := showPickFileDialog(me)
	if callback != nil {
		callback(ok, paths)
	}
}

// ---------------- SaveFileDialog -------------------
func SaveFileDialog() *saveFileDialog {
	return &saveFileDialog{}
}

type saveFileDialog struct {
	directory string
	saveName  string
}

func (me *saveFileDialog) SetDirectory(directory string) *saveFileDialog {
	// TODO:: check is dir
	me.directory = directory
	return me
}

func (me *saveFileDialog) SetSaveName(saveName string) *saveFileDialog {
	me.saveName = saveName
	return me
}

func (me *saveFileDialog) ShowModal(callback func(ok bool, saveName string)) {
	if !_fileCheck(me.directory, true) {
		return
	}

	ok, saveName := showSaveFileDialog(me)
	if callback != nil {
		callback(ok, saveName)
	}
}
