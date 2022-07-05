// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gtk3

type GtkFileChooserConfirmation int

const (
	FILE_CHOOSER_CONFIRMATION_CONFIRM GtkFileChooserConfirmation = iota
	FILE_CHOOSER_CONFIRMATION_ACCEPT_FILENAME
	FILE_CHOOSER_CONFIRMATION_SELECT_AGAIN
)

type GtkFileChooserAction int

const (
	FILE_CHOOSER_ACTION_OPEN GtkFileChooserAction = iota
	FILE_CHOOSER_ACTION_SAVE
	FILE_CHOOSER_ACTION_SELECT_FOLDER
	FILE_CHOOSER_ACTION_CREATE_FOLDER
)

type GtkResponseType int

const (
	RESPONSE_NONE         GtkResponseType = -1
	RESPONSE_REJECT       GtkResponseType = -2
	RESPONSE_ACCEPT       GtkResponseType = -3
	RESPONSE_DELETE_EVENT GtkResponseType = -4
	RESPONSE_OK           GtkResponseType = -5
	RESPONSE_CANCEL       GtkResponseType = -6
	RESPONSE_CLOSE        GtkResponseType = -7
	RESPONSE_YES          GtkResponseType = -8
	RESPONSE_NO           GtkResponseType = -9
	RESPONSE_APPLY        GtkResponseType = -10
	RESPONSE_HELP         GtkResponseType = -11
)
