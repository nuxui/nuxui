// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import (
	"nuxui.org/nuxui/nux/internal/win32/com"
	"syscall"
	"unsafe"
)

var (
	modeshell32                     = syscall.NewLazyDLL("shell32.dll")
	procSHBrowseForFolderW          = modeshell32.NewProc("SHBrowseForFolderW")
	procShellExecuteW               = modeshell32.NewProc("ShellExecuteW")
	procILCreateFromPathW           = modeshell32.NewProc("ILCreateFromPathW")
	procILFree                      = modeshell32.NewProc("ILFree")
	procSHOpenFolderAndSelectItems  = modeshell32.NewProc("SHOpenFolderAndSelectItems")
	procSHCreateItemFromParsingName = modeshell32.NewProc("SHCreateItemFromParsingName")
)

// https://docs.microsoft.com/en-us/windows/win32/api/shlobj_core/nf-shlobj_core-shbrowseforfolderw
func SHBrowseForFolder(info *BROWSEINFO) (uintptr, error) {
	ret, _, err := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(info)))
	if ret != 0 {
		err = nil
	}
	return ret, err
}

// https://docs.microsoft.com/zh-cn/windows/win32/api/shellapi/nf-shellapi-shellexecutew
// showCmd https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-showwindow
func ShellExecute(hwnd uintptr, operation, file, parameters, directory string, showCmd int32) (uintptr, error) {
	sOperation, _ := syscall.UTF16PtrFromString(operation)
	sFile, _ := syscall.UTF16PtrFromString(file)
	sParam, _ := syscall.UTF16PtrFromString(parameters)
	sDir, _ := syscall.UTF16PtrFromString(directory)
	ret, _, err := procShellExecuteW.Call(
		hwnd,
		uintptr(unsafe.Pointer(sOperation)),
		uintptr(unsafe.Pointer(sFile)),
		uintptr(unsafe.Pointer(sParam)),
		uintptr(unsafe.Pointer(sDir)),
		uintptr(showCmd))
	if ret != 0 {
		err = nil
	}
	return ret, err
}

func ILCreateFromPath(path string) (*ITEMIDLIST, error) {
	sPath, _ := syscall.UTF16PtrFromString(path)
	ret, _, err := procILCreateFromPathW.Call(uintptr(unsafe.Pointer(sPath)))
	if ret != 0 {
		err = nil
		return (*ITEMIDLIST)(unsafe.Pointer(ret)), nil
	}
	return nil, err
}

func ILFree(pointer *ITEMIDLIST) error {
	ret, _, err := procILFree.Call(uintptr(unsafe.Pointer(pointer)))
	if ret != 0 {
		err = nil
	}
	return err
}

func SHOpenFolderAndSelectItems(dir *ITEMIDLIST, count uint32, childArray []*ITEMIDLIST, flags uint32) error {
	ret, _, err := procSHOpenFolderAndSelectItems.Call(
		uintptr(unsafe.Pointer(dir)),
		uintptr(count),
		uintptr(unsafe.Pointer(&childArray[0])),
		uintptr(flags))
	if ret == 0 {
		err = nil
	}
	return err
}

// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-shcreateitemfromparsingname
func SHCreateItemFromParsingName(name string, pbc *com.IBindCtx, iid *com.GUID, object unsafe.Pointer) error {
	pszPath, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	ret, _, err := procSHCreateItemFromParsingName.Call(
		uintptr(unsafe.Pointer(pszPath)),
		uintptr(unsafe.Pointer(pbc)),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(object)))

	if ret == 0 {
		return nil
	}
	return err
}
