// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

void nux_deleteGlobalRef(jobject globalRef);
void nux_deleteLocalRef(jobject localRef);
*/
import "C"

type (
	JNIEnv *C.JNIEnv
)

type (
	JObject       C.jobject
	Activity      C.jobject
	Application   C.jobject
	View          C.jobject
	SurfaceHolder C.jobject
	Canvas        C.jobject
	Paint         C.jobject
	Path          C.jobject
	Bitmap        C.jobject
	StaticLayout  C.jobject
)

func DeleteGlobalRef(globalRef JObject) {
	C.nux_deleteGlobalRef(C.jobject(globalRef))
}

func DeleteLocalRef(localRef JObject) {
	C.nux_deleteLocalRef(C.jobject(localRef))
}
