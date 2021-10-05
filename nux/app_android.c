// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android

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
#include <android/native_activity.h>
#include <android/log.h>
#include <android/configuration.h>
#include <android/looper.h>
#include <android/native_activity.h>
#include "_cgo_export.h"


#define LOG_VERB(...) __android_log_print(ANDROID_LOG_VERBOSE, "nux", __VA_ARGS__)
#define LOG_ERR(...) __android_log_print(ANDROID_LOG_ERROR, "nux", __VA_ARGS__)
#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "nux", __VA_ARGS__)

struct android_app {
    ANativeActivity* activity;
    AConfiguration* config;
    ALooper* looper;
    AInputQueue* inputQueue;
    ANativeWindow* window;

    int msgread;
    int msgwrite;
    pthread_t thread;
    pthread_mutex_t mutex;
    pthread_cond_t cond;

	void* savedState;
    size_t savedStateSize;

    AInputQueue* pendingInputQueue;
    ANativeWindow* pendingWindow;
    ARect pendingContentRect;

	int running;
	int destroyed;
};

enum {
    APP_CMD_INPUT_CHANGED = 1,
    APP_CMD_DESTROY = 2,
};

enum {
    LOOPER_ID_CMD = 1,
    LOOPER_ID_INPUT = 2,
    LOOPER_ID_USER = 3,
};

// Call the Go main.main
void call_main_main(){
	uintptr_t mainPC = (uintptr_t)dlsym(RTLD_DEFAULT, "main.main");
	if (!mainPC) {
		LOG_FATAL("missing main.main");
	}
	callMain(mainPC);
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

// Fix per Google for bug https://code.google.com/p/android/issues/detail?id=41755
void process_input(struct android_app* app) {
    AInputEvent* event = NULL;
    while (AInputQueue_getEvent(app->inputQueue, &event) >= 0) {
		LOG_VERB("new event %d", AInputEvent_getType(event));

		// https://developer.android.com/ndk/reference/group/input#ainputqueue_predispatchevent
        if (AInputQueue_preDispatchEvent(app->inputQueue, event) != 0) { 
			LOG_VERB("AInputQueue_preDispatchEvent %d", AInputEvent_getType(event));
            continue;
        }
        int32_t handled = (int32_t)onInputEvent(event);
        AInputQueue_finishEvent(app->inputQueue, event, handled);
    }
}

void android_app_pre_exec_cmd(struct android_app* android_app, int8_t cmd) {
    switch (cmd) {
        case APP_CMD_INPUT_CHANGED:
            LOG_VERB("APP_CMD_INPUT_CHANGED\n");
            pthread_mutex_lock(&android_app->mutex);
            if (android_app->inputQueue != NULL) {
                AInputQueue_detachLooper(android_app->inputQueue);
            }
            android_app->inputQueue = android_app->pendingInputQueue;
            if (android_app->inputQueue != NULL) {
                LOG_VERB("Attaching input queue to looper");
                AInputQueue_attachLooper(android_app->inputQueue, android_app->looper, LOOPER_ID_INPUT, NULL, NULL);
            }
            pthread_cond_broadcast(&android_app->cond);
            pthread_mutex_unlock(&android_app->mutex);
            break;
    }
}

void android_app_loop(struct android_app *android_app){
	int events;
	int ident;
	int8_t cmd;
	while(1){
		events = 0;
		ident = 0;
		cmd = 0;

		ident = ALooper_pollAll(-1, NULL, &events, NULL);
		switch (ident) {
		case LOOPER_ID_CMD:
			if ( read(android_app->msgread, &cmd, sizeof(cmd)) == sizeof(cmd) ) {
                if (cmd == APP_CMD_DESTROY) {
                    goto end;
                }
				android_app_pre_exec_cmd(android_app, cmd);
			} else {
				LOG_ERR("No data on command pipe!");
			}
			break;
		case LOOPER_ID_INPUT:
			process_input(android_app);
			break;
		case LOOPER_ID_USER:
			break;
		}
	}
end:
    LOG_VERB("exited, end of android looper");
}

void android_app_destroy(struct android_app *android_app){
    LOG_VERB("android_app_destroied");
    // ANativeActivity_finish(android_app->activity);

    close(android_app->msgread);
    close(android_app->msgwrite);
    pthread_cond_destroy(&android_app->cond);
    pthread_mutex_destroy(&android_app->mutex);
    free(android_app);
}

void* android_app_entry(void* param) {
    struct android_app* android_app = (struct android_app*)param;
	
	// https://developer.android.com/ndk/reference/group/configuration
    android_app->config = AConfiguration_new();
    AConfiguration_fromAssetManager(android_app->config, android_app->activity->assetManager);

    ALooper* looper = ALooper_prepare(ALOOPER_PREPARE_ALLOW_NON_CALLBACKS);
    ALooper_addFd(looper, android_app->msgread, LOOPER_ID_CMD, ALOOPER_EVENT_INPUT, NULL, NULL);
    android_app->looper = looper;

    pthread_mutex_lock(&android_app->mutex);
    android_app->running = 1;
    pthread_cond_broadcast(&android_app->cond);
    pthread_mutex_unlock(&android_app->mutex);

    nativeLoopPrepared();

    android_app_loop(android_app);

    ALooper_removeFd(looper, android_app->msgread);
    ALooper_release(looper);

    pthread_mutex_lock(&android_app->mutex);
    android_app->running = 0;
    pthread_cond_broadcast(&android_app->cond);
    pthread_mutex_unlock(&android_app->mutex);



    // android_app_destroy(android_app);
    return NULL;
}

struct android_app* android_app_create(ANativeActivity* activity, void* savedState, size_t savedStateSize) {
    call_main_main();

    initWindow(activity);
    
    // window.creating called must before pthread_create
    // init window at here

    struct android_app* android_app = (struct android_app*)malloc(sizeof(struct android_app));
    memset(android_app, 0, sizeof(struct android_app));
    android_app->activity = activity;

	// createApp( (uintptr_t)android_app);

    pthread_mutex_init(&android_app->mutex, NULL);
    pthread_cond_init(&android_app->cond, NULL);

    if (savedState != NULL) {
        android_app->savedState = malloc(savedStateSize);
        android_app->savedStateSize = savedStateSize;
        memcpy(android_app->savedState, savedState, savedStateSize);
    }

    int msgpipe[2];
    if (pipe(msgpipe)) {
        LOG_FATAL("could not create pipe: %s", strerror(errno));
        return NULL;
    }
    android_app->msgread = msgpipe[0];
    android_app->msgwrite = msgpipe[1];

    pthread_attr_t attr; 
    pthread_attr_init(&attr);
    pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
    pthread_create(&android_app->thread, &attr, android_app_entry, android_app);

    // Wait for thread to start.
    pthread_mutex_lock(&android_app->mutex);
    while (!android_app->running) {
        pthread_cond_wait(&android_app->cond, &android_app->mutex);
    }
    pthread_mutex_unlock(&android_app->mutex);

    return android_app;
}

void android_app_write_cmd(struct android_app* android_app, int8_t cmd) {
    if (write(android_app->msgwrite, &cmd, sizeof(cmd)) != sizeof(cmd)) {
        LOG_ERR("Failure writing android_app cmd: %s\n", strerror(errno));
    }
}

void android_app_set_input(struct android_app* android_app, AInputQueue* inputQueue) {
    pthread_mutex_lock(&android_app->mutex);
    android_app->pendingInputQueue = inputQueue;
    android_app_write_cmd(android_app, APP_CMD_INPUT_CHANGED);
    while (android_app->inputQueue != android_app->pendingInputQueue) {
        pthread_cond_wait(&android_app->cond, &android_app->mutex);
    }
    pthread_mutex_unlock(&android_app->mutex);
}

void onInputQueueCreated(ANativeActivity *activity, AInputQueue *queue){
    LOG_VERB("onInputQueueCreated");
	android_app_set_input((struct android_app*)activity->instance, queue);
}

void onInputQueueDestroyed(ANativeActivity *activity, AInputQueue *queue){
    LOG_VERB("onInputQueueDestroyed");
	android_app_set_input((struct android_app*)activity->instance, NULL);
}

void onDestroyA(ANativeActivity *activity){
    LOG_VERB("onDestroyA");
    onDestroy(activity);

    struct android_app* android_app = (struct android_app*)activity->instance;
    pthread_mutex_lock(&android_app->mutex);
    android_app_write_cmd(android_app, APP_CMD_DESTROY);
    while (android_app->running != 0) {
        pthread_cond_wait(&android_app->cond, &android_app->mutex);
    }
    pthread_mutex_unlock(&android_app->mutex);

    close(android_app->msgread);
    close(android_app->msgwrite);
    pthread_cond_destroy(&android_app->cond);
    pthread_mutex_destroy(&android_app->mutex);
    free(android_app);

    LOG_VERB("onDestroyA end");
}

void ANativeActivity_onCreate(ANativeActivity *activity, void* savedState, size_t savedStateSize) {
    LOG_VERB("ANativeActivity_onCreate");
	// These functions match the methods on Activity, described at
	// https://developer.android.com/ndk/reference/struct/a-native-activity-callbacks#struct_a_native_activity_callbacks
	//
	// Note that onNativeWindowResized is not called on resize. Avoid it.
	// https://code.google.com/p/android/issues/detail?id=180645

    activity->callbacks->onStart = onStart;
    activity->callbacks->onResume = onResume;
    activity->callbacks->onPause = onPause;
    activity->callbacks->onStop = onStop;
	activity->callbacks->onDestroy = onDestroyA;
    activity->callbacks->onLowMemory = onLowMemory;
    activity->callbacks->onInputQueueCreated = onInputQueueCreated;
    activity->callbacks->onInputQueueDestroyed = onInputQueueDestroyed;
    activity->callbacks->onSaveInstanceState = onSaveInstanceState;
    activity->callbacks->onConfigurationChanged = onConfigurationChanged;
    activity->callbacks->onContentRectChanged = onContentRectChanged;
    activity->callbacks->onNativeWindowCreated = onNativeWindowCreated;
	activity->callbacks->onNativeWindowResized = onNativeWindowResized;
	activity->callbacks->onNativeWindowRedrawNeeded = onNativeWindowRedrawNeeded;
    activity->callbacks->onWindowFocusChanged = onWindowFocusChanged;
    activity->callbacks->onNativeWindowDestroyed = onNativeWindowDestroyed;

	activity->instance = android_app_create(activity, savedState, savedStateSize);
}