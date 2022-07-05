// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"nuxui.org/nuxui/nux/internal/darwin"
)

func showViewFileDialog(dialog *viewFileDialog) {
	if dialog.directory != "" {
		darwin.SharedNSWorkspace().OpenURL(dialog.directory)
	}

	if len(dialog.activeFileNames) > 0 {
		urls := []darwin.NSURL{}
		for _, name := range dialog.activeFileNames {
			url := darwin.NSURLfileURLWithPath(name)
			if !url.IsNil() {
				urls = append(urls, url)
			}
		}

		darwin.SharedNSWorkspace().ActivateFileViewerSelectingURLs(urls)
	}
}

func showPickFileDialog(dialog *pickFileDialog) (ok bool, paths []string) {
	p := darwin.SharedNSOpenPanel()
	p.SetCanChooseFiles(dialog.allowsChooseFiles)
	p.SetCanChooseDirectories(dialog.allowsChooseFolders)
	p.SetCanCreateDirectories(dialog.allowsCreateFolders)
	p.SetAllowsMultipleSelection(dialog.allowsMultipleSelection)

	if dialog.directory != "" {
		p.SetDirectoryURL(darwin.NSURLfileURLWithPath(dialog.directory))
	}

	// add filters
	if len(dialog.filters) > 0 {
		types := []darwin.UTType{}
		for _, ext := range dialog.filters {
			if ext != "" {
				uttype := darwin.UTTypeWithFilenameExtension(ext)
				if !uttype.IsNil() {
					types = append(types, uttype)
				}
			}
		}
		p.SetAllowedContentTypes(types)
	}

	ok = p.RunModal() == darwin.NSModalResponseOK
	for _, url := range p.URLs() {
		paths = append(paths, url.Path())
	}
	return
}

func showSaveFileDialog(dialog *saveFileDialog) (ok bool, saveName string) {
	p := darwin.SharedNSSavePanel()
	if dialog.directory != "" {
		p.SetDirectoryURL(darwin.NSURLfileURLWithPath(dialog.directory))
	}

	if dialog.saveName != "" {
		p.SetNameFieldStringValue(dialog.saveName)
	}

	return p.RunModal() == darwin.NSModalResponseOK, p.URL().Path()
}
