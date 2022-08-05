// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

#include <jni.h>
#include <stdlib.h>
#include <unistd.h>
#include <dlfcn.h>
#include <pthread.h>
#include <sys/types.h>
#include <android/log.h>
#include "_cgo_export.h"

#ifdef __cplusplus
extern "C" {
#endif

#define LOG_VERB(...)  __android_log_print(ANDROID_LOG_VERBOSE, "nuxui", __VA_ARGS__)
#define LOG_ERR(...)   __android_log_print(ANDROID_LOG_ERROR, "nuxui", __VA_ARGS__)
#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "nuxui", __VA_ARGS__)

static JNIEnv* nuxJNIEnv;


// Call the Go main.main
void nux_call_main_main(){
	uintptr_t mainPC = (uintptr_t)dlsym(RTLD_DEFAULT, "main.main");
	if (!mainPC) {
		LOG_FATAL("missing main.main");
	}
	go_nux_callMain(mainPC);
}

jobject nux_newLocalRef(jobject ref){
    return (*nuxJNIEnv)->NewLocalRef(nuxJNIEnv, ref);
}

jobject nux_newGlobalRef(jobject ref){
    return (*nuxJNIEnv)->NewGlobalRef(nuxJNIEnv, ref);
}

void nux_deleteLocalRef(jobject localRef){
    (*nuxJNIEnv)->DeleteLocalRef(nuxJNIEnv, localRef);
}

void nux_deleteGlobalRef(jobject globalRef){
    (*nuxJNIEnv)->DeleteGlobalRef(nuxJNIEnv, globalRef);
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

static jfieldID nux_find_static_fieldid(JNIEnv *env, jclass clazz, const char *name, const char *sig){
    jfieldID f = (*env)->GetStaticFieldID(env, clazz, name, sig);
    if (f == 0) {
        (*env)->ExceptionClear(env);
        LOG_FATAL("cannot find field %s %s", name, sig);
        return 0;
    }
    return f;
}

struct{
    jclass clazz;
    jobject instance;
    jmethodID backToUI;
} NuxApplication;

struct{
    jclass clazz;
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
    jclass clsApplication = nux_find_class(env, "org/nuxui/app/NuxApplication");
    NuxApplication.clazz = nux_newGlobalRef(clsApplication);
    NuxApplication.backToUI = nux_find_method(env, NuxApplication.clazz, "backToUI", "()V");
    NuxApplication.instance = (*env)->GetStaticObjectField(env, NuxApplication.clazz, nux_find_static_fieldid(env, NuxApplication.clazz, "mInstance", "Lorg/nuxui/app/NuxApplication;"));
    nux_deleteLocalRef(clsApplication);
    
    jclass clsActivity = nux_find_class(env, "org/nuxui/app/NuxActivity");
    NuxActivity.clazz = nux_newGlobalRef(clsActivity);
    NuxActivity.drawText = nux_find_static_method(env, NuxActivity.clazz, "drawText", "(Landroid/graphics/Canvas;Ljava/lang/String;ILandroid/text/TextPaint;)V");
    NuxActivity.createStaticLayout = nux_find_static_method(env, NuxActivity.clazz, "createStaticLayout", "(Ljava/lang/String;ILandroid/text/TextPaint;)Landroid/text/StaticLayout;");
    NuxActivity.createBitmap = nux_find_static_method(env, NuxActivity.clazz, "createBitmap", "(Ljava/lang/String;)Landroid/graphics/Bitmap;");
    nux_deleteLocalRef(clsActivity);

    jclass clsSurfaceHolder = nux_find_class(env, "android/view/SurfaceHolder");
    SurfaceHolder.clazz = nux_newGlobalRef(clsSurfaceHolder);
    SurfaceHolder.lockCanvas = nux_find_method(env, SurfaceHolder.clazz, "lockCanvas", "()Landroid/graphics/Canvas;");
    SurfaceHolder.unlockCanvasAndPost = nux_find_method(env, SurfaceHolder.clazz, "unlockCanvasAndPost", "(Landroid/graphics/Canvas;)V");
    nux_deleteLocalRef(clsSurfaceHolder);

    jclass clsPaint = nux_find_class(env, "android/text/TextPaint");
    Paint.clazz = nux_newGlobalRef(clsPaint);
    Paint.init = nux_find_method(env, Paint.clazz, "<init>", "()V");
    Paint.setColor = nux_find_method(env, Paint.clazz, "setColor", "(I)V");
    Paint.getColor = nux_find_method(env, Paint.clazz, "getColor", "()I");
    Paint.setTextSize = nux_find_method(env, Paint.clazz, "setTextSize", "(F)V");
    Paint.setStyle = nux_find_method(env, Paint.clazz, "setStyle", "(Landroid/graphics/Paint$Style;)V");
    Paint.setAntiAlias = nux_find_method(env, Paint.clazz, "setAntiAlias", "(Z)V");
    nux_deleteLocalRef(clsPaint);

    jclass clsStyle = nux_find_class(env, "android/graphics/Paint$Style");
    Style.clazz = nux_newGlobalRef(clsStyle);
    Style.FILL = (*env)->GetStaticFieldID(env, Style.clazz, "FILL", "Landroid/graphics/Paint$Style;");
    nux_deleteLocalRef(clsStyle);

    jclass clsCanvas = nux_find_class(env, "android/graphics/Canvas");
    Canvas.clazz = nux_newGlobalRef(clsCanvas);
    Canvas.save = nux_find_method(env, Canvas.clazz, "save", "()I");
    Canvas.restore = nux_find_method(env, Canvas.clazz, "restore", "()V");
    Canvas.translate = nux_find_method(env, Canvas.clazz, "translate", "(FF)V");
    Canvas.scale = nux_find_method(env, Canvas.clazz, "scale", "(FF)V");
    Canvas.rotate = nux_find_method(env, Canvas.clazz, "rotate", "(F)V");
    Canvas.clipRect = nux_find_method(env, Canvas.clazz, "clipRect", "(FFFF)Z");
    Canvas.drawColor = nux_find_method(env, Canvas.clazz, "drawColor", "(I)V");
    Canvas.drawRect = nux_find_method(env, Canvas.clazz, "drawRect", "(FFFFLandroid/graphics/Paint;)V");
    Canvas.drawRoundRect = nux_find_method(env, Canvas.clazz, "drawRoundRect", "(FFFFFFLandroid/graphics/Paint;)V");
    Canvas.drawBitmap = nux_find_method(env, Canvas.clazz, "drawBitmap", "(Landroid/graphics/Bitmap;Landroid/graphics/Rect;Landroid/graphics/RectF;Landroid/graphics/Paint;)V");
    nux_deleteLocalRef(clsCanvas);

    jclass clsStaticLayout = nux_find_class(env, "android/text/StaticLayout");
    StaticLayout.clazz = nux_newGlobalRef(clsStaticLayout);
    StaticLayout.init = nux_find_method(env, StaticLayout.clazz, "<init>", "(Ljava/lang/CharSequence;Landroid/text/TextPaint;ILandroid/text/Layout$Alignment;FFZ)V");
    StaticLayout.getWidth = nux_find_method(env, StaticLayout.clazz, "getWidth", "()I");
    StaticLayout.getHeight = nux_find_method(env, StaticLayout.clazz, "getHeight", "()I");
    nux_deleteLocalRef(clsStaticLayout);

    jclass clsBitmap = nux_find_class(env, "android/graphics/Bitmap");
    Bitmap.clazz = nux_newGlobalRef(clsBitmap);
    Bitmap.getWidth = nux_find_method(env, Bitmap.clazz, "getWidth", "()I");
    Bitmap.getHeight = nux_find_method(env, Bitmap.clazz, "getHeight", "()I");
    Bitmap.recycle = nux_find_method(env, Bitmap.clazz, "recycle", "()V");
    nux_deleteLocalRef(clsBitmap);

    jclass clsRect = nux_find_class(env, "android/graphics/Rect");
    Rect.clazz = nux_newGlobalRef(clsRect);
    Rect.init = nux_find_method(env, Rect.clazz, "<init>", "(IIII)V");
    nux_deleteLocalRef(clsRect);

    jclass clsRectF = nux_find_class(env, "android/graphics/RectF");
    RectF.clazz = nux_newGlobalRef(clsRectF);
    RectF.init = nux_find_method(env, RectF.clazz, "<init>", "(FFFF)V");
    nux_deleteLocalRef(clsRectF);
}

void nux_BackToUI(){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, NuxApplication.instance, NuxApplication.backToUI);
}

jobject nux_new_Rect(jint left, jint top, jint right, jint bottom){
    jobject rect = (*nuxJNIEnv)->NewObject(nuxJNIEnv, Rect.clazz, Rect.init, left, top, right, bottom);
    jobject g = nux_newGlobalRef(rect);
    nux_deleteLocalRef(rect);
    return g;
}

jobject nux_new_RectF(jfloat left, jfloat top, jfloat right, jfloat bottom){
    jobject rectf = (*nuxJNIEnv)->NewObject(nuxJNIEnv, RectF.clazz, RectF.init, left, top, right, bottom);
    jobject g = nux_newGlobalRef(rectf);
    nux_deleteLocalRef(rectf);
    return g;
}

jobject nux_surfaceHolder_lockCanvas(jobject surfaceHolder){
    jobject canvas = (*nuxJNIEnv)->CallObjectMethod(nuxJNIEnv, surfaceHolder, SurfaceHolder.lockCanvas);
    jobject g = nux_newGlobalRef(canvas);
    nux_deleteLocalRef(canvas);
    return g;
}

void nux_surfaceHolder_unlockCanvas(jobject surfaceHolder, jobject canvas){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, surfaceHolder, SurfaceHolder.unlockCanvasAndPost, canvas);
}

jint nux_canvas_save(jobject canvas){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, canvas, Canvas.save);
}

void nux_canvas_restore(jobject canvas){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.restore);
}

void nux_canvas_translate(jobject canvas, jfloat x, jfloat y){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.translate, x, y);
}

void nux_canvas_scale(jobject canvas, jfloat x, jfloat y){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.scale, x, y);
}

void nux_canvas_rotate(jobject canvas, jfloat degrees){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.rotate, degrees);
}

void nux_canvas_clipRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom){
    jboolean ret = (*nuxJNIEnv)->CallBooleanMethod(nuxJNIEnv, canvas, Canvas.clipRect, left, top, right, bottom);
}

void nux_canvas_drawColor(jobject canvas, uint32_t color){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.drawColor, color);
}

void nux_canvas_drawRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.drawRect, left, top, right, bottom, paint);
}

void nux_canvas_drawRoundRect(jobject canvas, jfloat left, jfloat top, jfloat right, jfloat bottom, jfloat rx, jfloat ry, jobject paint){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.drawRoundRect, left, top, right, bottom, rx, ry, paint);
}

void nux_canvas_drawBitmap(jobject canvas, jobject bitmap, jfloat left, jfloat top, jfloat right, jfloat bottom, jobject paint){
    jobject rectf = nux_new_RectF(left, top, right, bottom);
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, canvas, Canvas.drawBitmap, bitmap, NULL, rectf, paint);
    nux_deleteGlobalRef(rectf);
}

void nux_canvas_drawText(jobject canvas, char* text, jint width, jobject paint){
    jstring str = (*nuxJNIEnv)->NewStringUTF(nuxJNIEnv, text);
    (*nuxJNIEnv)->CallStaticVoidMethod(nuxJNIEnv, NuxActivity.clazz, NuxActivity.drawText, canvas, str, width, paint);
    nux_deleteLocalRef(str);
}

// --------------- StaticLayout ----------------------

jobject nux_new_StaticLayout(jstring text, jint width, jobject paint){
    jobject layout = (*nuxJNIEnv)->CallStaticObjectMethod(nuxJNIEnv, NuxActivity.clazz, NuxActivity.createStaticLayout, text, width, paint);
    jobject g = nux_newGlobalRef(layout);
    nux_deleteLocalRef(layout);
    return g;
}

jint nux_staticLayout_getWidth(jobject staticLayout){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, staticLayout, StaticLayout.getWidth);
}

jint nux_staticLayout_getHeight(jobject staticLayout){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, staticLayout, StaticLayout.getHeight);
}

// --------------- Paint ----------------------

jobject nux_new_Paint(){
    jobject paint = (*nuxJNIEnv)->NewObject(nuxJNIEnv, Paint.clazz, Paint.init);
    jobject g = nux_newGlobalRef(paint);
    nux_deleteLocalRef(paint);
    return g;
}

void nux_paint_setTextSize(jobject paint, jfloat textSize){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, paint, Paint.setTextSize, textSize);
}

void nux_paint_setColor(jobject paint, uint32_t color){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, paint, Paint.setColor, color);
}

jint nux_paint_getColor(jobject paint){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, paint, Paint.getColor);
}

void nux_paint_setStyle(jobject paint, jint style){
    // LOG_ERR("Paint.setStyle=%p objPaint=%p Style.FILL=%p", Paint.setStyle, objPaint, Style.FILL);
    // (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, paint, Paint.setStyle, Style.FILL);
}

void nux_paint_setAntiAlias(jobject paint, jboolean aa){
    (*nuxJNIEnv)->CallVoidMethod(nuxJNIEnv, paint, Paint.setAntiAlias, aa);
}

void nux_paint_measureText(jobject paint, char* text, jint width, jint *outWidth, jint* outHeight){
    jstring str = (*nuxJNIEnv)->NewStringUTF(nuxJNIEnv, text);
    jobject layout = nux_new_StaticLayout(str, width, paint);
    *outWidth  = nux_staticLayout_getWidth(layout);
    *outHeight = nux_staticLayout_getHeight(layout);
    nux_deleteLocalRef(str);
    nux_deleteGlobalRef(layout);
}

// --------------- Bitmap ----------------------

jobject nux_createBitmap(char* fileName){
    jstring str = (*nuxJNIEnv)->NewStringUTF(nuxJNIEnv, fileName);
    jobject bitmap = (*nuxJNIEnv)->CallStaticObjectMethod(nuxJNIEnv, NuxActivity.clazz, NuxActivity.createBitmap, str);
    jobject g = nux_newGlobalRef(bitmap);
    nux_deleteLocalRef(str);
    nux_deleteLocalRef(bitmap);
    return g;
}

jint nux_bitmap_getWidth(jobject bitmap){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, bitmap, Bitmap.getWidth);
}

jint nux_bitmap_getHeight(jobject bitmap){
    return (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, bitmap, Bitmap.getHeight);
}

void nux_bitmap_recycle(jobject bitmap){
    (*nuxJNIEnv)->CallIntMethod(nuxJNIEnv, bitmap, Bitmap.recycle);
}

jobject nux_NuxApplication_instance(){
    return NuxApplication.instance;
}

// --------------- NuxActivity lifecycle ----------------------

void native_NuxActivity_onCreate(JNIEnv *env, jobject activity, jbyteArray native_saved_state) {
    LOG_VERB("native_NuxActivity_onCreate thiz=%p, gettid = %d", activity, gettid());
    go_NuxActivity_onCreate(activity);
}

void native_NuxActivity_onStart(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onStart begin ");
    go_NuxActivity_onStart(activity);
}

void native_NuxActivity_onRestart(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onRestart begin ");
    go_NuxActivity_onRestart(activity);
}

void native_NuxActivity_onResume(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onResume begin ");
    go_NuxActivity_onResume(activity);
}

void native_NuxActivity_onPause(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onPause begin ");
    go_NuxActivity_onPause(activity);
}

void native_NuxActivity_onStop(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onStop begin ");
    go_NuxActivity_onStop(activity);
}

void native_NuxActivity_onDestroy(JNIEnv *env, jobject activity) {
    LOG_VERB("native_NuxActivity_onDestroy begin ");
    go_NuxActivity_onDestroy(activity);
}

void native_NuxActivity_surfaceCreated(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceCreated activity=%ud,%p, surfaceHolder=%ud,%p", activity, activity, surfaceHolder, surfaceHolder);
    go_NuxActivity_surfaceCreated(activity, surfaceHolder);
}

void native_NuxActivity_surfaceChanged(JNIEnv *env, jobject activity, jobject surfaceHolder,
                                                         jint format, jint width, jint height) {
    LOG_VERB("native_NuxActivity_surfaceChanged begin ");
    go_NuxActivity_surfaceChanged(activity, surfaceHolder, format, width, height);
}

void native_NuxActivity_surfaceRedrawNeeded(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceRedrawNeeded begin ");
    go_NuxActivity_surfaceRedrawNeeded(activity, surfaceHolder);
}

void native_NuxActivity_surfaceDestroyed(JNIEnv *env, jobject activity, jobject surfaceHolder) {
    LOG_VERB("native_NuxActivity_surfaceDestroyed begin ");
    go_NuxActivity_surfaceDestroyed(activity, surfaceHolder);
}

jboolean native_NuxActivity_onTouch(JNIEnv *env, jobject thiz, jint device_id, jint pointer_id, jint action, jfloat x, jfloat y) {
    if (go_NuxActivity_onTouch(device_id, pointer_id, action, x, y) > 0) {
        return (jboolean)1;
    }
    return (jboolean)0;
}

void native_NuxApplication_onBackToUI(JNIEnv *env, jobject thiz) {
    LOG_VERB("native_NuxApplication_onBackToUI ");
    go_nux_backToUI();
}

void native_NuxApplication_onConfigurationChanged(JNIEnv *env, jobject application, jobject newConfig){
    LOG_VERB("native_NuxApplication_onConfigurationChanged ");
    go_NuxApplication_onConfigurationChanged(application, newConfig);
}

void native_NuxApplication_onCreate(JNIEnv *env, jobject application, jfloat density){
    LOG_VERB("native_NuxApplication_onCreate ");
    nuxJNIEnv = env;
    // put here can make sure the init func and main func run only once
    nux_init_classes(env);
    // any go code run first time will triger go init func 
    nux_call_main_main();
    
    go_NuxApplication_onCreate(application);
}

void native_NuxApplication_onLowMemory(JNIEnv *env, jobject application){
    LOG_VERB("native_NuxApplication_onLowMemory ");
    go_NuxApplication_onLowMemory(application);
}

void native_NuxApplication_onTerminate(JNIEnv *env, jobject application){
    LOG_VERB("native_NuxApplication_onTerminate ");
    go_NuxApplication_onTerminate(application);
}

void native_NuxApplication_onTrimMemory(JNIEnv *env, jobject application, jint level){
    LOG_VERB("native_NuxApplication_onTrimMemory ");
    go_NuxApplication_onTrimMemory(application, level);
}

JNIEXPORT jint JNI_OnLoad(JavaVM* vm, void* reserved) {
    JNIEnv* env;
    if ( (*vm)->GetEnv(vm, (void**)(&env), JNI_VERSION_1_6) != JNI_OK ) {
        return JNI_ERR;
    }

    // NuxApplication methods
    jclass clsNuxApplication = (*env)->FindClass(env, "org/nuxui/app/NuxApplication");
    if (clsNuxApplication == NULL) return JNI_ERR;

    static const JNINativeMethod methodsNuxApplication[] = {
        {"native_NuxApplication_onConfigurationChanged", "(Landroid/content/res/Configuration;)V", &native_NuxApplication_onConfigurationChanged},
        {"native_NuxApplication_onCreate",               "(F)V",                                   &native_NuxApplication_onCreate},
        {"native_NuxApplication_onLowMemory",            "()V",                                    &native_NuxApplication_onLowMemory},
        {"native_NuxApplication_onTerminate",            "()V",                                    &native_NuxApplication_onTerminate},
        {"native_NuxApplication_onTrimMemory",           "(I)V",                                   &native_NuxApplication_onTrimMemory},
        {"native_NuxApplication_onBackToUI",             "()V",                                    &native_NuxApplication_onBackToUI},
    };

    int rc = (*env)->RegisterNatives(env, clsNuxApplication, methodsNuxApplication, sizeof(methodsNuxApplication)/sizeof(JNINativeMethod));
    if (rc != JNI_OK) return rc;

    // NuxActivity methods
    jclass clsNuxActivity = (*env)->FindClass(env, "org/nuxui/app/NuxActivity");
    if (clsNuxActivity == NULL) return JNI_ERR;

    static const JNINativeMethod methodsNuxActivity[] = {
        {"native_NuxActivity_onCreate",              "([B)V",                                &native_NuxActivity_onCreate},
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


