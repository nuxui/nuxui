// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

#include <errno.h>
#include <fcntl.h>
#include <stdint.h>
#include <string.h>
#include <sched.h>
#include <sys/syscall.h>
#include <jni.h>
#include <dlfcn.h>
#include <pthread.h>
#include <unistd.h>
#include <sys/types.h>
#include <android/log.h>
#include "_cgo_export.h"

#ifdef __cplusplus
extern "C" {
#endif

#define LOG_VERB(...) __android_log_print(ANDROID_LOG_VERBOSE, "nuxui", __VA_ARGS__)
#define LOG_ERR(...) __android_log_print(ANDROID_LOG_ERROR, "nuxui", __VA_ARGS__)
#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "nuxui", __VA_ARGS__)

static jfloat sDensity = 1;
static pid_t mainThreadId;

// Call the Go main.main
void nux_call_main_main(){
	uintptr_t mainPC = (uintptr_t)dlsym(RTLD_DEFAULT, "main.main");
	if (!mainPC) {
		LOG_FATAL("missing main.main");
	}
    LOG_VERB("nux_call_main_main 0 ");
	go_nux_callMain(mainPC);
}

static jclass nux_find_class(JNIEnv *env, const char *class_name) {
    jclass clazz = (*env)->FindClass(env, class_name);
    if (clazz == NULL) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find %s", class_name);
        return NULL;
    }
    return clazz;
}

static jmethodID nux_find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
    jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
    if (m == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }
    return m;
}

static jmethodID nux_find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
    jmethodID m = (*env)->GetStaticMethodID(env, clazz, name, sig);
    if (m == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }
    return m;
}

JNIEnv* jnienv;

struct{
    jclass clazz;
    jobject thiz;
    jmethodID instance;
    jmethodID backToUI;
} NuxApplication;

struct{
    jclass clazz;
    // jobject thiz;
    jmethodID drawText;
    jmethodID createStaticLayout;
    jmethodID createBitmap;
} NuxActivity;

struct {
    jclass clazz;
    jmethodID lockCanvas;
    jmethodID unlockCanvasAndPost;
} SurfaceHolder;

struct {
    jclass clazz;
    jmethodID init;
    jmethodID setColor;
    jmethodID getColor;
    jmethodID setTextSize;
    jmethodID setStyle;
    jmethodID setAntiAlias;
} Paint;

struct {
    jclass clazz;
    jfieldID FILL;
} Style;

struct {
    JNIEnv* env;
    jclass clazz;
    jmethodID save;
    jmethodID restore;
    jmethodID translate;
    jmethodID scale;
    jmethodID rotate;
    jmethodID clipRect;
    jmethodID drawColor;
    jmethodID drawRect;
    jmethodID drawRoundRect;
    jmethodID drawBitmap;
} Canvas;

struct {
    jclass clazz;
    jmethodID init;
    jmethodID create;
    jmethodID getWidth;
    jmethodID getHeight;
} StaticLayout;

struct {
    jclass clazz;
    jmethodID getWidth;
    jmethodID getHeight;
    jmethodID recycle;
} Bitmap;

struct {
    jclass clazz;
    jmethodID init;
} Rect;

struct {
    jclass clazz;
    jmethodID init;
} RectF;

void nux_init_classes(JNIEnv *env){
    jnienv = env;

    jclass clsApplication = nux_find_class(env, "org/nuxui/app/NuxApplication");
    NuxApplication.clazz = (*env)->NewGlobalRef(env, clsApplication);
    NuxApplication.instance = nux_find_static_method(env, NuxApplication.clazz, "instance", "()Lorg/nuxui/app/NuxApplication;");
    NuxApplication.backToUI = nux_find_method(env, NuxApplication.clazz, "backToUI", "()V");
    
    jclass clsActivity = nux_find_class(env, "org/nuxui/app/NuxActivity");
    NuxActivity.clazz = (*env)->NewGlobalRef(env, clsActivity);
    NuxActivity.drawText = nux_find_static_method(env, NuxActivity.clazz, "drawText", "(Landroid/graphics/Canvas;Ljava/lang/String;ILandroid/text/TextPaint;)V");
    NuxActivity.createStaticLayout = nux_find_static_method(env, NuxActivity.clazz, "createStaticLayout", "(Ljava/lang/String;ILandroid/text/TextPaint;)Landroid/text/StaticLayout;");
    NuxActivity.createBitmap = nux_find_static_method(env, NuxActivity.clazz, "createBitmap", "(Ljava/lang/String;)Landroid/graphics/Bitmap;");

    jclass clsSurfaceHolder = nux_find_class(env, "android/view/SurfaceHolder");
    SurfaceHolder.clazz = (*env)->NewGlobalRef(env, clsSurfaceHolder);
    SurfaceHolder.lockCanvas = nux_find_method(env, SurfaceHolder.clazz, "lockCanvas", "()Landroid/graphics/Canvas;");
    SurfaceHolder.unlockCanvasAndPost = nux_find_method(env, SurfaceHolder.clazz, "unlockCanvasAndPost", "(Landroid/graphics/Canvas;)V");

    jclass clsPaint = nux_find_class(env, "android/text/TextPaint");
    Paint.clazz = (*env)->NewGlobalRef(env, clsPaint);
    Paint.init = nux_find_method(env, Paint.clazz, "<init>", "()V");
    Paint.setColor = nux_find_method(env, Paint.clazz, "setColor", "(I)V");
    Paint.getColor = nux_find_method(env, Paint.clazz, "getColor", "()I");
    Paint.setTextSize = nux_find_method(env, Paint.clazz, "setTextSize", "(F)V");
    Paint.setStyle = nux_find_method(env, Paint.clazz, "setStyle", "(Landroid/graphics/Paint$Style;)V");
    Paint.setAntiAlias = nux_find_method(env, Paint.clazz, "setAntiAlias", "(Z)V");

    jclass clsStyle = nux_find_class(env, "android/graphics/Paint$Style");
    Style.clazz = (*env)->NewGlobalRef(env, clsStyle);
    Style.FILL = (*env)->GetStaticFieldID(env, Style.clazz, "FILL", "Landroid/graphics/Paint$Style;");

    jclass clsCanvas = nux_find_class(env, "android/graphics/Canvas");
    Canvas.clazz = (*env)->NewGlobalRef(env, clsCanvas);
    Canvas.save = nux_find_method(env, Canvas.clazz, "save", "()I");
    Canvas.restore = nux_find_method(env, Canvas.clazz, "restore", "()V");
    Canvas.translate = nux_find_method(env, Canvas.clazz, "translate", "(FF)V");
    Canvas.scale = nux_find_method(env, Canvas.clazz, "scale", "(FF)V");
    Canvas.rotate = nux_find_method(env, Canvas.clazz, "rotate", "(F)V");
    Canvas.clipRect = nux_find_method(env, Canvas.clazz, "clipRect", "(FFFF)Z");
    // Canvas.drawColor = nux_find_method(env, Canvas.clazz, "drawColor", "(I)V");
    Canvas.drawRect = nux_find_method(env, Canvas.clazz, "drawRect", "(FFFFLandroid/graphics/Paint;)V");
    Canvas.drawRoundRect = nux_find_method(env, Canvas.clazz, "drawRoundRect", "(FFFFFFLandroid/graphics/Paint;)V");
    Canvas.drawBitmap = nux_find_method(env, Canvas.clazz, "drawBitmap", "(Landroid/graphics/Bitmap;Landroid/graphics/Rect;Landroid/graphics/RectF;Landroid/graphics/Paint;)V");

    jclass clsStaticLayout = nux_find_class(env, "android/text/StaticLayout");
    StaticLayout.clazz = (*env)->NewGlobalRef(env, clsStaticLayout);
    StaticLayout.init = nux_find_method(env, StaticLayout.clazz, "<init>", "(Ljava/lang/CharSequence;Landroid/text/TextPaint;ILandroid/text/Layout$Alignment;FFZ)V");
    StaticLayout.getWidth = nux_find_method(env, StaticLayout.clazz, "getWidth", "()I");
    StaticLayout.getHeight = nux_find_method(env, StaticLayout.clazz, "getHeight", "()I");

    jclass clsBitmap = nux_find_class(env, "android/graphics/Bitmap");
    Bitmap.clazz = (*env)->NewGlobalRef(env, clsBitmap);
    Bitmap.getWidth = nux_find_method(env, Bitmap.clazz, "getWidth", "()I");
    Bitmap.getHeight = nux_find_method(env, Bitmap.clazz, "getHeight", "()I");
    Bitmap.recycle = nux_find_method(env, Bitmap.clazz, "recycle", "()V");

    jclass clsRect = nux_find_class(env, "android/graphics/Rect");
    Rect.clazz = (*env)->NewGlobalRef(env, clsRect);
    Rect.init = nux_find_method(env, Rect.clazz, "<init>", "(IIII)V");

    jclass clsRectF = nux_find_class(env, "android/graphics/RectF");
    RectF.clazz = (*env)->NewGlobalRef(env, clsRectF);
    RectF.init = nux_find_method(env, RectF.clazz, "<init>", "(FFFF)V");
}

void nux_deleteGlobalRef(jobject globalRef){
    (*jnienv)->DeleteGlobalRef(jnienv, globalRef);
}

void nux_deleteLocalRef(jobject localRef){
    (*jnienv)->DeleteLocalRef(jnienv, localRef);
}

void nux_BackToUI(){
    (*jnienv)->CallVoidMethod(jnienv, NuxApplication.thiz, NuxApplication.backToUI);
}

jobject nux_new_Rect(jint left, jint top, jint right, jint bottom){
    return (*jnienv)->NewObject(jnienv, Rect.clazz, Rect.init, left, top, right, bottom);
}

jobject nux_new_RectF(jfloat left, jfloat top, jfloat right, jfloat bottom){
    return (*jnienv)->NewObject(jnienv, RectF.clazz, RectF.init, left, top, right, bottom);
}

jobject nux_surfaceHolder_lockCanvas(jobject surfaceHolder){
    return (*jnienv)->CallObjectMethod(jnienv, surfaceHolder, SurfaceHolder.lockCanvas);
}

void nux_surfaceHolder_unlockCanvas(jobject surfaceHolder, jobject canvas){
    (*jnienv)->CallVoidMethod(jnienv, surfaceHolder, SurfaceHolder.unlockCanvasAndPost, canvas);
}

jint nux_canvas_save(jobject canvas){
    return (*jnienv)->CallIntMethod(jnienv, canvas, Canvas.save);
}

void nux_canvas_restore(jobject canvas){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.restore);
}

void nux_canvas_translate(jobject canvas, jfloat x, jfloat y){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.translate, x, y);
}

void nux_canvas_scale(jobject canvas, jfloat x, jfloat y){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.scale, x, y);
}

void nux_canvas_rotate(jobject canvas, jfloat degrees){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.rotate, degrees);
}

void nux_canvas_clipRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom){
    jboolean ret = (*jnienv)->CallBooleanMethod(jnienv, canvas, Canvas.clipRect, left, top, right, bottom);
}

// void nux_canvas_drawColor(jobject canvas, uint32_t color){
//     (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawColor, color);
// }

void nux_canvas_drawRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawRect, left, top, right, bottom, paint);
}

void nux_canvas_drawRoundRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jfloat rx, jfloat ry, jobject paint){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawRoundRect, left, top, right, bottom, rx, ry, paint);
}

void nux_canvas_drawBitmap(jobject canvas, jobject bitmap, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint){
    // jobject rect = nux_new_Rect(0,0,0,0);
    jobject rectf = nux_new_RectF(left, top, right, bottom);
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawBitmap, bitmap, NULL, rectf, paint);
    // (*jnienv)->DeleteLocalRef(jnienv, rect);
    (*jnienv)->DeleteLocalRef(jnienv, rectf);
}

void nux_canvas_drawText(jobject canvas, char* text, jint width, jobject paint){
    jstring str = (*jnienv)->NewStringUTF(jnienv, text);
    (*jnienv)->CallStaticVoidMethod(jnienv, NuxActivity.clazz, NuxActivity.drawText, canvas, str, width, paint);
    (*jnienv)->DeleteLocalRef(jnienv, str);
}

// --------------- StaticLayout ----------------------

jobject nux_new_StaticLayout(jstring text, jint width, jobject paint){
    return (*jnienv)->CallStaticObjectMethod(jnienv, NuxActivity.clazz, NuxActivity.createStaticLayout, text, width, paint);
}

jint nux_staticLayout_getWidth(jobject staticLayout){
    return (*jnienv)->CallIntMethod(jnienv, staticLayout, StaticLayout.getWidth);
}

jint nux_staticLayout_getHeight(jobject staticLayout){
    return (*jnienv)->CallIntMethod(jnienv, staticLayout, StaticLayout.getHeight);
}

// --------------- Paint ----------------------

jobject nux_new_Paint(){
    return (*jnienv)->NewObject(jnienv, Paint.clazz, Paint.init);
}

void nux_paint_setTextSize(jobject paint, jfloat textSize){
    textSize *= sDensity;
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setTextSize, textSize);
}

void nux_paint_setColor(jobject paint, uint32_t color){
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setColor, color);
}

jint nux_paint_getColor(jobject paint){
    return (*jnienv)->CallIntMethod(jnienv, paint, Paint.getColor);
}

void nux_paint_setStyle(jobject paint, jint style){
    // LOG_ERR("Paint.setStyle=%p objPaint=%p Style.FILL=%p", Paint.setStyle, objPaint, Style.FILL);
    // (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setStyle, Style.FILL);
}

void nux_paint_setAntiAlias(jobject paint, jboolean aa){
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setAntiAlias, aa);
}

void nux_paint_measureText(jobject paint, char* text, jint width, jint *outWidth, jint* outHeight){
    jstring str = (*jnienv)->NewStringUTF(jnienv, text);
    jobject layout = nux_new_StaticLayout(str, width, paint);
    *outWidth = nux_staticLayout_getWidth(layout);
    *outHeight = nux_staticLayout_getHeight(layout);
    (*jnienv)->DeleteLocalRef(jnienv, str);
    (*jnienv)->DeleteLocalRef(jnienv, layout);
}

// --------------- Bitmap ----------------------

jobject nux_createBitmap(char* fileName){
    jstring str = (*jnienv)->NewStringUTF(jnienv, fileName);
    jobject bitmap = (*jnienv)->CallStaticObjectMethod(jnienv, NuxActivity.clazz, NuxActivity.createBitmap, str);
    (*jnienv)->DeleteLocalRef(jnienv, str);
    return (*jnienv)->NewGlobalRef(jnienv, bitmap);
}

jint nux_bitmap_getWidth(jobject bitmap){
    return (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.getWidth);
}

jint nux_bitmap_getHeight(jobject bitmap){
    return (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.getHeight);
}

void nux_bitmap_recycle(jobject bitmap){
    (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.recycle);
    (*jnienv)->DeleteGlobalRef(jnienv, bitmap);
}

jobject nux_NuxApplication_instance(){
    jobject instance = (*jnienv)->CallStaticObjectMethod(jnienv, NuxApplication.clazz, NuxApplication.instance);
    return (*jnienv)->NewGlobalRef(jnienv, instance);
}

// --------------- NuxActivity lifecycle ----------------------

void native_NuxActivity_onCreate(JNIEnv *env, jobject activity, jbyteArray native_saved_state) {
    mainThreadId = gettid();
    LOG_VERB("native_NuxActivity_onCreate thiz=%p, gettid = %d", activity, gettid());
    // NuxActivity.thiz = (*env)->NewGlobalRef(env, activity);  // TODO:: delete this ref when destroy

    LOG_VERB("native_NuxActivity_onCreate gettid 1 = %d", gettid());
    // nux_call_main_main();
    LOG_VERB("native_NuxActivity_onCreate gettid 2 = %d", gettid());
}

void native_NuxActivity_onStart(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onStart begin ");
}

void native_NuxActivity_onRestart(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onRestart begin ");
}

void native_NuxActivity_onResume(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onResume begin ");
}

void native_NuxActivity_onPause(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onPause begin ");
}

void native_NuxActivity_onStop(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onStop begin ");
}

void native_NuxActivity_onDestroy(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onDestroy begin ");
}

void native_NuxActivity_surfaceRedrawNeeded(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceRedrawNeeded begin ");
    go_nux_surfaceRedrawNeeded(env, activity, surfaceHolder);
}

void native_NuxActivity_surfaceCreated(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceCreated activity=%ud,%p, surfaceHolder=%ud,%p", activity, activity, surfaceHolder, surfaceHolder);
    go_nux_surfaceCreated(env, activity, surfaceHolder);
}

void native_NuxActivity_surfaceChanged(JNIEnv *env, jobject activity, jobject surfaceHolder,
                                                         jint format, jint width, jint height) {
    LOG_VERB("native_NuxActivity_surfaceChanged begin ");
    go_nux_surfaceChanged(env, activity, surfaceHolder, format, width, height);
}

void native_NuxActivity_surfaceDestroyed(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceDestroyed begin ");
}

jboolean native_NuxActivity_onTouch(JNIEnv *env, jobject thiz, jint device_id, jint pointer_id, jint action, jfloat x, jfloat y) {
    go_nux_onPointerEvent(device_id, pointer_id, action, x, y);
}

void native_NuxApplication_onBackToUI(JNIEnv *env, jobject thiz) {
    LOG_VERB("native_NuxApplication_onBackToUI ");
    go_nux_backToUI();
}

void native_NuxApplication_onCreate(JNIEnv *env, jobject thiz, jfloat density){
    sDensity = density;
    LOG_VERB("native_NuxApplication_onCreate ");

    // put here can make sure the init func and main func run only once
    nux_init_classes(env);

}

JNIEXPORT jint JNI_OnLoad(JavaVM* vm, void* reserved) {
    LOG_VERB("JNI_OnLoad gettid = %d", gettid());

    JNIEnv* env;
    if ( (*vm)->GetEnv(vm, (void**)(&env), JNI_VERSION_1_6) != JNI_OK ) {
        return JNI_ERR;
    }

    // NuxApplication methods
    jclass clsNuxApplication = (*env)->FindClass(env, "org/nuxui/app/NuxApplication");
    if (clsNuxApplication == NULL) return JNI_ERR;

    static const JNINativeMethod methodsNuxApplication[] = {
        {"native_NuxApplication_onCreate",           "(F)V",                                 &native_NuxApplication_onCreate},
        {"native_NuxApplication_onBackToUI",         "()V",                                  &native_NuxApplication_onBackToUI},
    };

    int rc = (*env)->RegisterNatives(env, clsNuxApplication, methodsNuxApplication, sizeof(methodsNuxApplication)/sizeof(JNINativeMethod));
    if (rc != JNI_OK) return rc;

    // NuxActivity methods
    jclass clsNuxActivity = (*env)->FindClass(env, "org/nuxui/app/NuxActivity");
    if (clsNuxActivity == NULL) return JNI_ERR;

    static const JNINativeMethod methodsNuxActivity[] = {
        {"native_NuxActivity_onCreate",              "([B)V",                               &native_NuxActivity_onCreate},
        {"native_NuxActivity_onStart",               "()V",                                  &native_NuxActivity_onStart},
        {"native_NuxActivity_onRestart",             "()V",                                  &native_NuxActivity_onRestart},
        {"native_NuxActivity_onResume",              "()V",                                  &native_NuxActivity_onResume},
        {"native_NuxActivity_onPause",               "()V",                                  &native_NuxActivity_onPause},
        {"native_NuxActivity_onStop",                "()V",                                  &native_NuxActivity_onStop},
        {"native_NuxActivity_onDestroy",             "()V",                                  &native_NuxActivity_onDestroy},
        {"native_NuxActivity_surfaceRedrawNeeded",   "(Landroid/view/SurfaceHolder;)V",      &native_NuxActivity_surfaceRedrawNeeded},
        {"native_NuxActivity_surfaceCreated",        "(Landroid/view/SurfaceHolder;)V",      &native_NuxActivity_surfaceCreated},
        {"native_NuxActivity_surfaceChanged",        "(Landroid/view/SurfaceHolder;III)V",   &native_NuxActivity_surfaceChanged},
        {"native_NuxActivity_surfaceDestroyed",      "(Landroid/view/SurfaceHolder;)V",      &native_NuxActivity_surfaceDestroyed},
        {"native_NuxActivity_onTouch",               "(IIIFF)Z",                             &native_NuxActivity_onTouch},
    };
    LOG_VERB("JNI_OnLoad gettid 0 = %d", gettid());

    rc = (*env)->RegisterNatives(env, clsNuxActivity, methodsNuxActivity, sizeof(methodsNuxActivity)/sizeof(JNINativeMethod));
    if (rc != JNI_OK) return rc;

    return JNI_VERSION_1_6;
}


#ifdef __cplusplus
}
#endif


