// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go:build (linux && !android)

package linux

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <strings.h>
#include <string.h>
#include <locale.h>
#include <sys/syscall.h>
#include <stdint.h>
*/
import "C"

type Category int

const (
	LC_CTYPE          Category = 0  //C.LC_CTYPE
	LC_NUMERIC        Category = 1  //C.LC_NUMERIC
	LC_TIME           Category = 2  //C.LC_TIME
	LC_COLLATE        Category = 3  //C.LC_COLLATE
	LC_MONETARY       Category = 4  //C.LC_MONETARY
	LC_MESSAGES       Category = 5  //C.LC_MESSAGES
	LC_ALL            Category = 6  //C.LC_ALL
	LC_PAPER          Category = 7  //C.LC_PAPER
	LC_NAME           Category = 8  //C.LC_NAME
	LC_ADDRESS        Category = 9  //C.LC_ADDRESS
	LC_TELEPHONE      Category = 10 //C.LC_TELEPHONE
	LC_MEASUREMENT    Category = 11 //C.LC_MEASUREMENT
	LC_IDENTIFICATION Category = 12 //C.LC_IDENTIFICATION
)
