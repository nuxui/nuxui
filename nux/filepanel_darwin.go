// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
)

type nativeFilePicker darwin.NSSavePanel

func showViewFilePanel(panel *viewFilePanel) {
	if panel.directory != "" {
		darwin.SharedNSWorkspace().OpenURL(panel.directory)
	}

	if len(panel.activeFileNames) > 0 {
		urls := []darwin.NSURL{}
		for _, name := range panel.activeFileNames {
			url := darwin.NSURLfileURLWithPath(name)
			if !url.IsNil() {
				urls = append(urls, url)
			}
		}

		darwin.SharedNSWorkspace().ActivateFileViewerSelectingURLs(urls)
	}
}

func showPickFilePanel(panel *pickFilePanel) (ok bool, paths []string) {
	p := darwin.SharedNSOpenPanel()
	p.SetCanChooseFiles(panel.canChooseFiles)
	p.SetCanChooseDirectories(panel.canChooseDirectories)
	p.SetAllowsMultipleSelection(panel.allowsMultipleSelection)

	if panel.directory != "" {
		p.SetDirectoryURL(darwin.NSURLfileURLWithPath(panel.directory))
	}
	types := []darwin.UTType{}
	for _, ext := range panel.filters {
		if ext != "" {
			uttype := darwin.UTTypeWithFilenameExtension(ext)
			if !uttype.IsNil() {
				types = append(types, uttype)
			}
		}
	}
	p.SetAllowedContentTypes(types)

	ok = p.RunModal() == darwin.NSModalResponseOK
	for _, url := range p.URLs() {
		paths = append(paths, url.Path())
	}
	return
}

func showSaveFilePanel(panel *saveFilePanel) (ok bool, saveName string) {
	p := darwin.SharedNSSavePanel()
	if panel.directory != "" {
		p.SetDirectoryURL(darwin.NSURLfileURLWithPath(panel.directory))
	}

	if panel.saveName != "" {
		p.SetNameFieldStringValue(panel.saveName)
	}

	return p.RunModal() == darwin.NSModalResponseOK, p.URL().Path()
}
