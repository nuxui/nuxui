// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

#include <jni.h>
#include <dlfcn.h>
#include <errno.h>
#include <fcntl.h>
#include <stdint.h>
#include <string.h>
#include <poll.h>
#include <pthread.h>
#include <sched.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/syscall.h>
#include <android/native_activity.h>
#include <android/log.h>
#include <android/configuration.h>
#include <android/looper.h>
#include <android/native_activity.h>
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
void call_main_main(){
	uintptr_t mainPC = dlsym(RTLD_DEFAULT, "main.main");
	if (!mainPC) {
		LOG_FATAL("missing main.main");
	}
	go_callMain(mainPC);
}

static jclass find_class(JNIEnv *env, const char *class_name) {
    jclass clazz = (*env)->FindClass(env, class_name);
    if (clazz == NULL) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find %s", class_name);
        return NULL;
    }
    return clazz;
}

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
    jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
    if (m == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find method %s %s", name, sig);
        return 0;
    }
    return m;
}

static jmethodID find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
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
    jmethodID drawText;
    jmethodID createStaticLayout;
    jmethodID createImage;
    jmethodID backToUI;
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

void initClasses(JNIEnv *env, jobject activity){
    jnienv = env;

    jclass clsActivity = find_class(env, "org/nuxui/app/NuxActivity");
    NuxActivity.thiz = (*env)->NewGlobalRef(env, activity);
    NuxActivity.clazz = (*env)->NewGlobalRef(env, clsActivity);
    NuxActivity.drawText = find_static_method(env, NuxActivity.clazz, "drawText", "(Landroid/graphics/Canvas;Ljava/lang/String;ILandroid/text/TextPaint;)V");
    NuxActivity.createStaticLayout = find_static_method(env, NuxActivity.clazz, "createStaticLayout", "(Ljava/lang/String;ILandroid/text/TextPaint;)Landroid/text/StaticLayout;");
    NuxActivity.createImage = find_static_method(env, NuxActivity.clazz, "createImage", "(Ljava/lang/String;)Landroid/graphics/Bitmap;");
    NuxActivity.backToUI = find_method(env, NuxActivity.clazz, "backToUI", "()V");
    LOG_VERB("backToUI %p ", NuxActivity.backToUI);

    jclass clsSurfaceHolder = find_class(env, "android/view/SurfaceHolder");
    SurfaceHolder.clazz = (*env)->NewGlobalRef(env, clsSurfaceHolder);
    SurfaceHolder.lockCanvas = find_method(env, SurfaceHolder.clazz, "lockCanvas", "()Landroid/graphics/Canvas;");
    SurfaceHolder.unlockCanvasAndPost = find_method(env, SurfaceHolder.clazz, "unlockCanvasAndPost", "(Landroid/graphics/Canvas;)V");

    jclass clsPaint = find_class(env, "android/text/TextPaint");
    Paint.clazz = (*env)->NewGlobalRef(env, clsPaint);
    Paint.init = find_method(env, Paint.clazz, "<init>", "()V");
    Paint.setColor = find_method(env, Paint.clazz, "setColor", "(I)V");
    Paint.setTextSize = find_method(env, Paint.clazz, "setTextSize", "(F)V");
    Paint.setStyle = find_method(env, Paint.clazz, "setStyle", "(Landroid/graphics/Paint$Style;)V");
    Paint.setAntiAlias = find_method(env, Paint.clazz, "setAntiAlias", "(Z)V");

    jclass clsStyle = find_class(env, "android/graphics/Paint$Style");
    Style.clazz = (*env)->NewGlobalRef(env, clsStyle);
    Style.FILL = (*env)->GetStaticFieldID(env, Style.clazz, "FILL", "Landroid/graphics/Paint$Style;");

    jclass clsCanvas = find_class(env, "android/graphics/Canvas");
    Canvas.clazz = (*env)->NewGlobalRef(env, clsCanvas);
    Canvas.save = find_method(env, Canvas.clazz, "save", "()I");
    Canvas.restore = find_method(env, Canvas.clazz, "restore", "()V");
    Canvas.translate = find_method(env, Canvas.clazz, "translate", "(FF)V");
    Canvas.scale = find_method(env, Canvas.clazz, "scale", "(FF)V");
    Canvas.rotate = find_method(env, Canvas.clazz, "rotate", "(F)V");
    Canvas.clipRect = find_method(env, Canvas.clazz, "clipRect", "(FFFF)Z");
    Canvas.drawColor = find_method(env, Canvas.clazz, "drawColor", "(I)V");
    Canvas.drawBitmap = find_method(env, Canvas.clazz, "drawBitmap", "(Landroid/graphics/Bitmap;Landroid/graphics/Rect;Landroid/graphics/RectF;Landroid/graphics/Paint;)V");

    jclass clsStaticLayout = find_class(env, "android/text/StaticLayout");
    StaticLayout.clazz = (*env)->NewGlobalRef(env, clsStaticLayout);
    StaticLayout.init = find_method(env, StaticLayout.clazz, "<init>", "(Ljava/lang/CharSequence;Landroid/text/TextPaint;ILandroid/text/Layout$Alignment;FFZ)V");
    StaticLayout.getWidth = find_method(env, StaticLayout.clazz, "getWidth", "()I");
    StaticLayout.getHeight = find_method(env, StaticLayout.clazz, "getHeight", "()I");

    jclass clsBitmap = find_class(env, "android/graphics/Bitmap");
    Bitmap.clazz = (*env)->NewGlobalRef(env, clsBitmap);
    Bitmap.getWidth = find_method(env, Bitmap.clazz, "getWidth", "()I");
    Bitmap.getHeight = find_method(env, Bitmap.clazz, "getHeight", "()I");
    Bitmap.recycle = find_method(env, Bitmap.clazz, "recycle", "()V");

    jclass clsRect = find_class(env, "android/graphics/Rect");
    Rect.clazz = (*env)->NewGlobalRef(env, clsRect);
    Rect.init = find_method(env, Rect.clazz, "<init>", "(IIII)V");

    jclass clsRectF = find_class(env, "android/graphics/RectF");
    RectF.clazz = (*env)->NewGlobalRef(env, clsRectF);
    RectF.init = find_method(env, RectF.clazz, "<init>", "(FFFF)V");
}

void deleteGlobalRef(jobject globalRef){
    (*jnienv)->DeleteGlobalRef(jnienv, globalRef);
}

void deleteLocalRef(jobject localRef){
    (*jnienv)->DeleteLocalRef(jnienv, localRef);
}

void backToUI(){
    (*jnienv)->CallVoidMethod(jnienv, NuxActivity.thiz, NuxActivity.backToUI);
}

int isMainThread(){
    return mainThreadId == gettid();
}

jobject new_Rect(jint left, jint top, jint right, jint bottom){
    return (*jnienv)->NewObject(jnienv, Rect.clazz, Rect.init, left, top, right, bottom);
}

jobject new_RectF(jfloat left, jfloat top, jfloat right, jfloat bottom){
    return (*jnienv)->NewObject(jnienv, RectF.clazz, RectF.init, left, top, right, bottom);
}

jobject surfaceHolder_lockCanvas(jobject surfaceHolder){
    return (*jnienv)->CallObjectMethod(jnienv, surfaceHolder, SurfaceHolder.lockCanvas);
}

void surfaceHolder_unlockCanvas(jobject surfaceHolder, jobject canvas){
    (*jnienv)->CallVoidMethod(jnienv, surfaceHolder, SurfaceHolder.unlockCanvasAndPost, canvas);
}

jint canvas_save(jobject canvas){
    return (*jnienv)->CallIntMethod(jnienv, canvas, Canvas.save);
}

void canvas_restore(jobject canvas){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.restore);
}

void canvas_translate(jobject canvas, jfloat x, jfloat y){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.translate, x, y);
}

void canvas_scale(jobject canvas, jfloat x, jfloat y){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.scale, x, y);
}

void canvas_rotate(jobject canvas, jfloat degrees){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.rotate, degrees);
}

void canvas_clipRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom){
    jboolean ret = (*jnienv)->CallBooleanMethod(jnienv, canvas, Canvas.clipRect, left, top, right, bottom);
}

void canvas_drawColor(jobject canvas, uint32_t color){
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawColor, color);
}

void canvas_drawBitmap(jobject canvas, jobject bitmap, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint){
    jobject rect = new_Rect(0,0,0,0);
    jobject rectf = new_RectF(left, top, right, bottom);
    (*jnienv)->CallVoidMethod(jnienv, canvas, Canvas.drawBitmap, bitmap, NULL, rectf, paint);
    (*jnienv)->DeleteLocalRef(jnienv, rect);
    (*jnienv)->DeleteLocalRef(jnienv, rectf);
}

void canvas_drawText(jobject canvas, char* text, jint width, jobject paint){
    jstring str = (*jnienv)->NewStringUTF(jnienv, text);
    (*jnienv)->CallStaticVoidMethod(jnienv, NuxActivity.clazz, NuxActivity.drawText, canvas, str, width, paint);
    (*jnienv)->DeleteLocalRef(jnienv, str);
}

// --------------- StaticLayout ----------------------

jobject new_StaticLayout(jstring text, jint width, jobject paint){
    return (*jnienv)->CallStaticObjectMethod(jnienv, NuxActivity.clazz, NuxActivity.createStaticLayout, text, width, paint);
}

jint staticLayout_getWidth(jobject staticLayout){
    return (*jnienv)->CallIntMethod(jnienv, staticLayout, StaticLayout.getWidth);
}

jint staticLayout_getHeight(jobject staticLayout){
    return (*jnienv)->CallIntMethod(jnienv, staticLayout, StaticLayout.getHeight);
}

// --------------- Paint ----------------------

jobject new_Paint(){
    return (*jnienv)->NewObject(jnienv, Paint.clazz, Paint.init);
}

void paint_setTextSize(jobject paint, jfloat textSize){
    textSize *= sDensity;
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setTextSize, textSize);
}

void paint_setColor(jobject paint, uint32_t color){
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setColor, color);
}

void paint_setStyle(jobject paint, jint style){
    // LOG_ERR("Paint.setStyle=%p objPaint=%p Style.FILL=%p", Paint.setStyle, objPaint, Style.FILL);
    // (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setStyle, Style.FILL);
}

void paint_setAntiAlias(jobject paint, jboolean aa){
    (*jnienv)->CallVoidMethod(jnienv, paint, Paint.setAntiAlias, aa);
}

void paint_measureText(jobject paint, char* text, jint width, jint *outWidth, jint* outHeight){
    jstring str = (*jnienv)->NewStringUTF(jnienv, text);
    jobject layout = new_StaticLayout(str, width, paint);
    *outWidth = staticLayout_getWidth(layout);
    *outHeight = staticLayout_getHeight(layout);
    (*jnienv)->DeleteLocalRef(jnienv, str);
    (*jnienv)->DeleteLocalRef(jnienv, layout);
}

// --------------- Bitmap ----------------------

jobject createImage(char* fileName){
    jstring str = (*jnienv)->NewStringUTF(jnienv, fileName);
    jobject bitmap = (*jnienv)->CallStaticObjectMethod(jnienv, NuxActivity.clazz, NuxActivity.createImage, str);
    (*jnienv)->DeleteLocalRef(jnienv, str);
    return (*jnienv)->NewGlobalRef(jnienv, bitmap);
}

jint bitmap_getWidth(jobject bitmap){
    return (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.getWidth);
}

jint bitmap_getHeight(jobject bitmap){
    return (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.getHeight);
}

void bitmap_recycle(jobject bitmap){
    (*jnienv)->CallIntMethod(jnienv, bitmap, Bitmap.recycle);
    (*jnienv)->DeleteGlobalRef(jnienv, bitmap);
}

// --------------- NuxActivity lifecycle ----------------------

void NuxActivity_onCreateNative(JNIEnv *env, jobject activity, jbyteArray native_saved_state, jfloat density) {
    sDensity = density;
    mainThreadId = gettid();
    LOG_VERB("NuxActivity_onCreateNative thiz=%p, gettid = %d", activity, gettid());
    initClasses(env, activity);
    call_main_main();
}

void NuxActivity_onStartNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onStartNative begin ");
}

void NuxActivity_onRestartNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onRestartNative begin ");
}

void NuxActivity_onResumeNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onResumeNative begin ");
}

void NuxActivity_onPauseNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onPauseNative begin ");
}

void NuxActivity_onStopNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onStopNative begin ");
}

void NuxActivity_onDestroyNative(JNIEnv *env, jobject activity) {
    LOG_VERB("NuxActivity_onDestroyNative begin ");
}

void NuxActivity_surfaceRedrawNeededNative(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("NuxActivity_surfaceRedrawNeededNative begin ");
    go_surfaceRedrawNeeded(env, activity, surfaceHolder);
}

void NuxActivity_surfaceCreatedNative(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("NuxActivity_surfaceCreatedNative activity=%ud,%p, surfaceHolder=%ud,%p", activity, activity, surfaceHolder, surfaceHolder);
    go_surfaceCreated(env, activity, surfaceHolder);
}

void NuxActivity_surfaceChangedNative(JNIEnv *env, jobject activity, jobject surfaceHolder,
                                                         jint format, jint width, jint height) {
    LOG_VERB("NuxActivity_surfaceChangedNative begin ");
    go_surfaceChanged(env, activity, surfaceHolder, format, width, height);
}

void NuxActivity_surfaceDestroyedNative(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("NuxActivity_surfaceDestroyedNative begin ");
}

jboolean NuxActivity_onTouchNative(JNIEnv *env, jobject thiz, jint device_id, jint pointer_id, jint action, jfloat x, jfloat y) {
    go_onPointerEvent(device_id, pointer_id, action, x, y);
}

void NuxActivity_onBackToUI(JNIEnv *env, jobject thiz) {
    LOG_VERB("NuxActivity_onBackToUI ");
    go_backToUI();
}

JNIEXPORT jint JNI_OnLoad(JavaVM* vm, void* reserved) {
    JNIEnv* env;
    if ( (*vm)->GetEnv(vm, (void**)(&env), JNI_VERSION_1_6) != JNI_OK ) {
        return JNI_ERR;
    }

    jclass c = (*env)->FindClass(env, "org/nuxui/app/NuxActivity");
    if (c == NULL) return JNI_ERR;

    static const JNINativeMethod methods[] = {
        {"onCreateNative",              "([BF)V",                               &NuxActivity_onCreateNative},
        {"onStartNative",               "()V",                                  &NuxActivity_onStartNative},
        {"onRestartNative",             "()V",                                  &NuxActivity_onRestartNative},
        {"onResumeNative",              "()V",                                  &NuxActivity_onResumeNative},
        {"onPauseNative",               "()V",                                  &NuxActivity_onPauseNative},
        {"onStopNative",                "()V",                                  &NuxActivity_onStopNative},
        {"onDestroyNative",             "()V",                                  &NuxActivity_onDestroyNative},
        {"surfaceRedrawNeededNative",   "(Landroid/view/SurfaceHolder;)V",      &NuxActivity_surfaceRedrawNeededNative},
        {"surfaceCreatedNative",        "(Landroid/view/SurfaceHolder;)V",      &NuxActivity_surfaceCreatedNative},
        {"surfaceChangedNative",        "(Landroid/view/SurfaceHolder;III)V",   &NuxActivity_surfaceChangedNative},
        {"surfaceDestroyedNative",      "(Landroid/view/SurfaceHolder;)V",      &NuxActivity_surfaceDestroyedNative},
        {"onTouchNative",               "(IIIFF)Z",                             &NuxActivity_onTouchNative},
        {"onBackToUI",                  "()V",                                  &NuxActivity_onBackToUI},
    };
    int rc = (*env)->RegisterNatives(env, c, methods, sizeof(methods)/sizeof(JNINativeMethod));
    if (rc != JNI_OK) return rc;

    return JNI_VERSION_1_6;
}


#ifdef __cplusplus
}
#endif


