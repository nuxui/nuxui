// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build !ios

#include "_cgo_export.h"
#include <pthread.h>
#include <stdio.h>

#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>
#import <OpenGL/gl3.h>

int32_t window_getWidth(uintptr_t window){
    return (int32_t)([(NSWindow*)window frame].size.width);
}

int32_t window_getHeight(uintptr_t window){
    return (int32_t)([(NSWindow*)window frame].size.height);
}
int32_t window_getContentWidth(uintptr_t window){
    return (int32_t)([[(NSWindow*)window contentView] bounds].size.width);
}

int32_t window_getContentHeight(uintptr_t window){
    return (int32_t)([[(NSWindow*)window contentView] bounds].size.height);
}

uint8_t* renew_buffer(int32_t width, int32_t height){
    size_t pitch = width * sizeof(uint32_t);
    uint8_t *buffer = malloc(pitch * height);
    return buffer;
}

size_t get_pitch(int32_t width){
    size_t pitch = width * sizeof(uint32_t);
    return pitch;
}

void flush_buffer(uintptr_t window, void* buf){
    NSWindow* w = (NSWindow*)window;
    size_t width = w.contentView.bounds.size.width;
    size_t height = w.contentView.bounds.size.height;
    NSLog(@"flush_buffer w=%d , h=%d", width, height);
    uint8_t *buffer = (uint8_t *)buf;
    size_t pitch = width * sizeof(uint32_t);

    NSBitmapImageRep *rep = [[[NSBitmapImageRep alloc] initWithBitmapDataPlanes:&buffer
        pixelsWide:width pixelsHigh:height
        bitsPerSample:8 samplesPerPixel:4 hasAlpha:YES isPlanar:NO
        colorSpaceName:NSDeviceRGBColorSpace
        bytesPerRow:pitch bitsPerPixel:sizeof(uint32_t) * 8] autorelease];
    NSImage *image = [[[NSImage alloc] initWithSize:NSMakeSize(width, height)] autorelease];
    [image addRepresentation:rep];
    w.contentView.wantsLayer = YES;
    w.contentView.layer.contents = image;
}