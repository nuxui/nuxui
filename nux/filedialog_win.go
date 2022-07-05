// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package nux

import (
	// "fmt"
	"nuxui.org/nuxui/log"
	"nuxui.org/nuxui/nux/internal/win32"
	"nuxui.org/nuxui/nux/internal/win32/com"
	"path/filepath"
	"syscall"
	"unsafe"
)

func showViewFileDialog(dialog *viewFileDialog) {
	if dialog.directory != "" {
		_, err := win32.ShellExecute(0, "open", dialog.directory, "", "", win32.SW_SHOWNORMAL)
		if err != nil {
			log.E("nuxui", "error call ShellExecute %s", err.Error())
		}
	}

	if len(dialog.activeFileNames) > 0 {
		all := map[string][]*win32.ITEMIDLIST{}

		for _, name := range dialog.activeFileNames {
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

func showPickFileDialog(dialog *pickFileDialog) (ok bool, paths []string) {
	if dialog.allowsChooseFolders{
		return showPickFolderPanel(dialog)
	}

	return showPickJustFileDialog(dialog)
}

func showPickFolderPanel(dialog *pickFileDialog) (ok bool, paths []string) {
	_, err := com.CoInitializeEx(nil, com.COINIT_APARTMENTTHREADED|com.COINIT_DISABLE_OLE1DDE)
	if err != nil {
		log.E("nuxui", "error call CoInitializeEx %s", err.Error())
		return false, nil
	}
	defer com.CoUninitialize()

	var fileOpenDialog *com.IFileOpenDialog
	_, err = com.CoCreateInstance(
		&com.CLSID_FileOpenDialog,
		nil,
		com.CLSCTX_ALL,
		&com.IID_IFileOpenDialog,
		unsafe.Pointer(&fileOpenDialog))

	if err != nil {
		log.E("nuxui", "error call CoCreateInstance %s", err.Error())
		return false, nil
	}
	defer fileOpenDialog.Release()

	options, err := fileOpenDialog.GetOptions()
	if err != nil {
		log.E("nuxui", "error call GetOptions %s", err.Error())
		return false, nil
	}

	if dialog.allowsChooseFiles {

	}

	if dialog.allowsChooseFolders {
		options |= com.FOS_PICKFOLDERS
	}

	if dialog.allowsMultipleSelection {
		options |= com.FOS_ALLOWMULTISELECT
	}

	if dialog.allowsCreateFolders {
		// explorer default allowed
	}

	fileOpenDialog.SetOptions(options)

	ret, err := fileOpenDialog.Show(theApp.MainWindow().native().hwnd)
	if err != nil {
		log.E("nuxui", "error call Show %s", err.Error())
		return false, nil
	}
	if ret != com.S_OK {
		return false, nil
	}

	if dialog.allowsMultipleSelection {
		results, err := fileOpenDialog.GetResults()
		if err != nil {
			log.E("nuxui", "error call GetResults %s", err.Error())
			return false, nil
		}
		defer results.Release()

		count, err := results.GetCount()
		if err != nil {
			log.E("nuxui", "error call GetCount %s", err.Error())
			return false, nil
		}

		var i uint32
		for i = 0; i != count; i++ {
			item, err := results.GetItemAt(i)
			if err != nil {
				log.E("nuxui", "error call GetItemAt %s", err.Error())
				return false, nil
			}
			defer item.Release()
			name, err := item.GetDisplayName(com.SIGDN_FILESYSPATH)
			if err != nil {
				log.E("nuxui", "error call GetDisplayName %s", err.Error())
				return false, nil
			}
			paths = append(paths, name)
		}
	} else {
		item, err := fileOpenDialog.GetResult()
		if err != nil {
			log.E("nuxui", "error call GetResult %s", err.Error())
			return false, nil
		}
		defer item.Release()
		name, err := item.GetDisplayName(com.SIGDN_FILESYSPATH)
		if err != nil {
			log.E("nuxui", "error call GetDisplayName %s", err.Error())
			return false, nil
		}
		paths = append(paths, name)
	}

	return true, paths
}

func showPickJustFileDialog(dialog *pickFileDialog) (ok bool, paths []string) {
	filePath := make([]uint16, 102400) // TODO:: out of range MAX_PATH ?
	var ofn win32.OPENFILENAME
	ofn.StructSize = uint32(unsafe.Sizeof(ofn))
	ofn.Owner = theApp.MainWindow().native().hwnd
	ofn.File = &filePath[0]
	ofn.MaxFile = uint32(len(filePath))

	if dialog.directory != "" {
		ofn.InitialDir, _ = syscall.UTF16PtrFromString(dialog.directory)
	}

	ofn.Flags |= win32.OFN_EXPLORER | win32.OFN_DONTADDTORECENT
	// ofn.Flags |= win32.OFN_FORCESHOWHIDDEN
	if dialog.allowsChooseFiles {
	}

	if dialog.allowsChooseFolders {
	}

	if dialog.allowsMultipleSelection {
		ofn.Flags |= win32.OFN_ALLOWMULTISELECT
	}

	if dialog.allowsCreateFolders {
		// explorer default allowed
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

func showPickFileDialog_deprecated(dialog *pickFileDialog) (ok bool, paths []string) {

	var info win32.BROWSEINFO
	info.Owner = theApp.MainWindow().native().hwnd
	info.Flags = win32.BIF_RETURNONLYFSDIRS

	_, err := win32.SHBrowseForFolder(&info)
	if err != nil {
		log.E("nuxui", "call SHBrowseForFolder error %s", err.Error())
	}

	return false, []string{}
}

func showSaveFileDialog(dialog *saveFileDialog) (ok bool, saveName string) {
	name, _ := syscall.UTF16FromString(dialog.saveName)
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

	if dialog.directory != "" {
		ofn.InitialDir, _ = syscall.UTF16PtrFromString(dialog.directory)
	}

	ok = win32.GetSaveFileName(&ofn)
	saveName = syscall.UTF16ToString(filePath)
	return
}
