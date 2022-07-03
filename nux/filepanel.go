// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Response int

const (
	ResponseCancel Response = 0
	ResponseOK     Response = 1
)

// ----------- ViewFilePanel -------------------
func ViewFilePanel() *viewFilePanel {
	return &viewFilePanel{}
}

type viewFilePanel struct {
	title           string
	directory       string
	activeFileNames []string
}

func (me *viewFilePanel) SetTitle(title string) *viewFilePanel {
	me.title = title
	return me
}

func (me *viewFilePanel) SetDirectory(directory string) *viewFilePanel {
	me.directory = directory
	return me
}

func (me *viewFilePanel) SetActiveFileNames(activeFileNames []string) *viewFilePanel {
	me.activeFileNames = activeFileNames
	return me
}

func (me *viewFilePanel) Show() {
	showViewFilePanel(me)
}

// ------------- PickFilePanel -----------------
func PickFilePanel() *pickFilePanel {
	return &pickFilePanel{}
}

type pickFilePanel struct {
	title                   string
	directory               string
	filters                 []string
	allowedContentTypes     []string
	canChooseFiles          bool
	canChooseDirectories    bool
	allowsMultipleSelection bool
}

func (me *pickFilePanel) SetDirectory(directory string) *pickFilePanel {
	me.directory = directory
	return me
}

func (me *pickFilePanel) SetFiltersByExtension(filters []string) *pickFilePanel {
	me.filters = filters
	return me
}

func (me *pickFilePanel) AllowsMultipleSelection() *pickFilePanel {
	me.allowsMultipleSelection = true
	return me
}

func (me *pickFilePanel) AllowChooseFiles() *pickFilePanel {
	me.canChooseFiles = true
	return me
}

func (me *pickFilePanel) AllowChooseDirectories() *pickFilePanel {
	me.canChooseDirectories = true
	return me
}

func (me *pickFilePanel) ShowModal(callback func(ok bool, results []string)) {
	ok, paths := showPickFilePanel(me)
	if callback != nil {
		callback(ok, paths)
	}
}

// ---------------- SaveFilePanel -------------------
func SaveFilePanel() *saveFilePanel {
	return &saveFilePanel{}
}

type saveFilePanel struct {
	directory string
	saveName  string
}

func (me *saveFilePanel) SetDirectory(directory string) *saveFilePanel {
	me.directory = directory
	return me
}

func (me *saveFilePanel) SetSaveName(saveName string) *saveFilePanel {
	me.saveName = saveName
	return me
}

func (me *saveFilePanel) ShowModal(callback func(ok bool, saveName string)) {
	ok, saveName := showSaveFilePanel(me)
	if callback != nil {
		callback(ok, saveName)
	}
}
