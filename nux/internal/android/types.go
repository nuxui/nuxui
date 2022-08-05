// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

/*
#include <jni.h>
#include <stdlib.h>

jobject nux_newLocalRef(jobject ref);
jobject nux_newGlobalRef(jobject ref);
void    nux_deleteGlobalRef(jobject globalRef);
void    nux_deleteLocalRef(jobject localRef);
*/
import "C"

type (
	JNIEnv *C.JNIEnv
)

type (
	JObject       C.jobject
	JLObject      C.jobject // Local Ref
	JGObject      C.jobject // Global Ref
	Activity      C.jobject
	Application   C.jobject
	View          C.jobject
	SurfaceHolder C.jobject
	Canvas        C.jobject
	Paint         C.jobject
	Path          C.jobject
	Bitmap        C.jobject
	StaticLayout  C.jobject
	Configuration C.jobject
)

func NewLocalRef(ref JLObject) JLObject {
	return JLObject(C.nux_newLocalRef(C.jobject(ref)))
}

func NewGlobalRef(ref JLObject) JGObject {
	return JGObject(C.nux_newGlobalRef(C.jobject(ref)))
}

func DeleteGlobalRef(globalRef JGObject) {
	C.nux_deleteGlobalRef(C.jobject(globalRef))
}

func DeleteLocalRef(localRef JLObject) {
	C.nux_deleteLocalRef(C.jobject(localRef))
}
