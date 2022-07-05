// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package com

type HRESULT uint32

// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.10240.0/shared/guiddef.h
type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}
