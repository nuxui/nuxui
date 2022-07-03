// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	// "fmt"
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
	"path/filepath"
	"syscall"
	"unsafe"
)

func showViewFilePanel(panel *viewFilePanel) {
	if panel.directory != "" {
		_, err := win32.ShellExecute(0, "open", panel.directory, "", "", win32.SW_SHOWNORMAL)
		if err != nil {
			log.E("nuxui", "error call ShellExecute %s", err.Error())
		}
	}

	if len(panel.activeFileNames) > 0 {
		all := map[string][]*win32.ITEMIDLIST{}

		for _, name := range panel.activeFileNames {
			list, ok := all[filepath.Dir(name)]
			if !ok {
				list = []*win32.ITEMIDLIST{}
			}
			item, err := win32.ILCreateFromPath(name)
			if err != nil {
				log.E("nuxui", "error call ILCreateFromPath(%s) %s", name, err.Error())
			} else {
				list = append(list, item)
			}
			all[filepath.Dir(name)] = list
		}

		// free
		defer func(clearall map[string][]*win32.ITEMIDLIST) {
			for _, v := range clearall {
				for _, item := range v {
					win32.ILFree(item)
				}
			}
		}(all)

		for k, v := range all {
			if v != nil && len(v) > 0 {
				dir, err := win32.ILCreateFromPath(k)
				if err != nil {
					log.E("nuxui", "error call ILCreateFromPath(%s) %s", k, err.Error())
				} else {
					win32.SHOpenFolderAndSelectItems(dir, uint32(len(v)), v, 0)
					win32.ILFree(dir)
				}
			}
		}
	}
}

func showPickFilePanel(panel *pickFilePanel) (ok bool, paths []string) {
	filePath := make([]uint16, 102400) // TODO:: out of range MAX_PATH ?
	var ofn win32.OPENFILENAME
	ofn.StructSize = uint32(unsafe.Sizeof(ofn))
	ofn.Owner = theApp.MainWindow().native().hwnd
	ofn.File = &filePath[0]
	ofn.MaxFile = uint32(len(filePath))

	if panel.directory != "" {
		ofn.InitialDir, _ = syscall.UTF16PtrFromString(panel.directory)
	}

	ofn.Flags |= win32.OFN_EXPLORER
	// ofn.Flags |= win32.OFN_DONTADDTORECENT
	// ofn.Flags |= win32.OFN_FORCESHOWHIDDEN
	if panel.canChooseFiles {
	}

	if panel.canChooseDirectories {
	}

	if panel.allowsMultipleSelection {
		ofn.Flags |= win32.OFN_ALLOWMULTISELECT
	}

	ok = win32.GetOpenFileName(&ofn)

	// for _, c := range filePath {
	// 	fmt.Printf("%d",c)
	// }
	// fmt.Println()

	if ok {
		path := syscall.UTF16ToString(filePath)
		paths = append(paths, path)

		split := path
		for i := len(split); filePath[i+1] != 0; i += len(split) + 1 {
			split = syscall.UTF16ToString(filePath[i+1:])
			paths = append(paths, path+"\\"+split)
		}
	}

	return ok, paths
}

func showPickFilePanel2(panel *pickFilePanel) (ok bool, paths []string) {

	var info win32.BROWSEINFO
	info.Owner = theApp.MainWindow().native().hwnd
	info.Flags = win32.BIF_RETURNONLYFSDIRS

	_, err := win32.SHBrowseForFolder(&info)
	if err != nil {
		log.E("nuxui", "call SHBrowseForFolder error %s", err.Error())
	}

	return false, []string{}
}

func showSaveFilePanel(panel *saveFilePanel) (ok bool, saveName string) {
	name, _ := syscall.UTF16FromString(panel.saveName)
	filePath := make([]uint16, len(name)+10240)

	for i, c := range name {
		filePath[i] = c
	}

	var ofn win32.OPENFILENAME
	ofn.StructSize = uint32(unsafe.Sizeof(ofn))
	ofn.Owner = theApp.MainWindow().native().hwnd
	ofn.File = &filePath[0]
	ofn.MaxFile = uint32(len(filePath))
	ofn.Flags = win32.OFN_OVERWRITEPROMPT

	if panel.directory != "" {
		ofn.InitialDir, _ = syscall.UTF16PtrFromString(panel.directory)
	}

	ok = win32.GetSaveFileName(&ofn)
	saveName = syscall.UTF16ToString(filePath)
	return
}
