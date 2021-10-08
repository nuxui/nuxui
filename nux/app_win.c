// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

﻿//go:build windows

#include "_cgo_export.h"
#include <windows.h>
#include <windowsx.h>
#include <imm.h>
#include <stdio.h>

HINSTANCE theInstance;
const char nuxWindowClass[] = "nuxWindowClass";
LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam);

int win32_main()
{
    WNDCLASSEX wc;
    HWND hwnd;
    MSG msg;

    wc.cbSize        = sizeof(WNDCLASSEX);
    wc.style         = 0;
    wc.lpfnWndProc   = WndProc;
    wc.cbClsExtra    = 0;
    wc.cbWndExtra    = 0;
    wc.hInstance     = theInstance;
    wc.hIcon         = LoadIcon(NULL, IDI_APPLICATION);
    wc.hCursor       = LoadCursor(NULL, IDC_ARROW);
    wc.hbrBackground = (HBRUSH)(COLOR_WINDOW+1);
    wc.lpszMenuName  = NULL;
    wc.lpszClassName = nuxWindowClass;
    wc.hIconSm       = LoadIcon(NULL, IDI_APPLICATION);

    if(!RegisterClassEx(&wc))
    {
        MessageBox(NULL, "Window Registration Failed!", "Error!",
            MB_ICONEXCLAMATION | MB_OK);
        return 0;
    }

    hwnd = CreateWindowEx(
        WS_EX_CLIENTEDGE,
        nuxWindowClass,
        "The title of my window",
        WS_OVERLAPPEDWINDOW,
        CW_USEDEFAULT, 
        CW_USEDEFAULT, 
        800, 
        600,
        NULL, 
        NULL, 
        theInstance, 
        NULL);

    if(hwnd == NULL)
    {
        MessageBox(NULL, "Window Creation Failed!", "Error!",
            MB_ICONEXCLAMATION | MB_OK);
        return 0;
    }

    ShowWindow(hwnd, SW_SHOWDEFAULT);
    UpdateWindow(hwnd);

    while(GetMessage(&msg, NULL, 0, 0) > 0)
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }
    return msg.wParam;
}

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{
    switch(msg)
    {
        /*mouse event begin*/
        case WM_NCHITTEST:
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_MOUSEMOVE:
        break;
        case WM_LBUTTONDOWN:
        case WM_LBUTTONUP:
        case WM_MBUTTONDOWN:
        case WM_MBUTTONUP:
        case WM_RBUTTONDOWN:
        case WM_RBUTTONUP:
        case WM_XBUTTONDOWN:
        case WM_XBUTTONUP:
        {
            // TODO:: did mouse event contain shift/ctrl key? 
            printf("mouse event msg = 0x%04x, wx=%d, wy=%d, x=%d, y=%d\n", 
                msg, GET_KEYSTATE_WPARAM(wParam), GET_WHEEL_DELTA_WPARAM(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam));
            go_mouseEvent(hwnd, msg, GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam));
            break;
        }
        case WM_MOUSEWHEEL:
        {
            int aMouseInfo[3];
            int lines = 1;
            if ( SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &aMouseInfo, 0) > 0 ){
                lines = aMouseInfo[0]; // get lines number of scroll each time
            }
            go_scrollEvent(hwnd, 0, (GET_WHEEL_DELTA_WPARAM(wParam) * lines / (double) WHEEL_DELTA) );
            break;
        }
        case WM_MOUSEHWHEEL:
        {
            int aMouseInfo[3];
            int lines = 1;
            if ( SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &aMouseInfo, 0) > 0 ){
                lines = aMouseInfo[0]; // get lines number of scroll each time
            }
            go_scrollEvent(hwnd, -(GET_WHEEL_DELTA_WPARAM(wParam) * lines / (double) WHEEL_DELTA), (double)0 );
            break;
        }
        case WM_KEYDOWN:
        case WM_KEYUP:
        case WM_SYSKEYDOWN:
        case WM_SYSKEYUP:
        {
            // printf("WM_KEYUP, wParam=0x%X, lParam=0x%X, &lParam=0x%X, lParam=0x%X, high=0x%X, low=0x%X\n", msg, wParam, lParam, &lParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            go_keyEvent(hwnd, msg, (UINT32)LOWORD(wParam), 0, 0, NULL);
            break;

        }
        case WM_CHAR:
            printf("WM_CHAR, wParam=%d, lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_SYSCHAR:
            printf("WM_SYSCHAR, wParam=%d, lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_DEADCHAR:
            printf("WM_DEADCHAR, wParam=%d, lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_SYSDEADCHAR:
            printf("WM_SYSDEADCHAR, wParam=%d, lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_UNICHAR:
            printf("WM_UNICHAR, wParam=%d, char=%c lParam=%d\n", msg, wParam, wParam, lParam);
            break;
        case WM_IME_STARTCOMPOSITION:
            printf("WM_IME_STARTCOMPOSITION, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_ENDCOMPOSITION:
            printf("WM_IME_ENDCOMPOSITION, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_COMPOSITION:   // = WM_IME_KEYLAST
        {
            HIMC hIMC;
            DWORD dwSize;
            HGLOBAL hstr;
            LPWSTR lpstr;
            printf("WM_IME_COMPOSITION wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            // if (lParam & GCS_RESULTSTR) 
            // {
                go_typingEvent(hwnd);
            // }
            break;
        }
        case WM_IME_CHAR:
        {
            printf("WM_IME_CHAR wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        }
        case WM_IME_SETCONTEXT:
            printf("WM_IME_SETCONTEXT, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_NOTIFY:
            printf("WM_IME_NOTIFY, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_CONTROL:
            printf("WM_IME_CONTROL, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_COMPOSITIONFULL:
            printf("WM_IME_COMPOSITIONFULL, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_SELECT:
            printf("WM_IME_SELECT, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_REQUEST:
            printf("WM_IME_REQUEST, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_KEYDOWN:
            printf("WM_IME_KEYDOWN, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_IME_KEYUP:
            printf("WM_IME_KEYUP, wParam=%d lParam=%d, high=%d, low=%d\n", msg, wParam, lParam, (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));
            break;
        case WM_CREATE:
        case WM_PAINT:
        case WM_SIZE:
            windowAction(hwnd, msg);
            break;

        // nonclient area of a window
        case WM_NCLBUTTONDBLCLK:
        case WM_NCLBUTTONDOWN:
        case WM_NCLBUTTONUP:
        case WM_NCMBUTTONDBLCLK:
        case WM_NCMBUTTONDOWN:
        case WM_NCMBUTTONUP:
        case WM_NCMOUSEHOVER:
        case WM_NCMOUSELEAVE:
        case WM_NCMOUSEMOVE:
        case WM_NCRBUTTONDBLCLK:
        case WM_NCRBUTTONDOWN:
        case WM_NCRBUTTONUP:
        case WM_NCXBUTTONDBLCLK:
        case WM_NCXBUTTONDOWN:
        case WM_NCXBUTTONUP:
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_CAPTURECHANGED:
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_LBUTTONDBLCLK:
        case WM_MBUTTONDBLCLK:
        case WM_RBUTTONDBLCLK:
        case WM_XBUTTONDBLCLK:
            printf("mouse double click msg=0x%04x\n", msg);
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_MOUSEACTIVATE:
        case WM_MOUSEHOVER:
        case WM_MOUSELEAVE:
            printf("mouse active hover leave msg=0x%04x\n", msg);
            return DefWindowProc(hwnd, msg, wParam, lParam);
        /*mouse event end*/
        case WM_CLOSE:
            DestroyWindow(hwnd);
            break;
        case WM_DESTROY:
            PostQuitMessage(0);
            break;
        default:
            // printf("######## event msg = 0x%04x\n", msg);
            return DefWindowProc(hwnd, msg, wParam, lParam);
    }
    // return 0;
    return DefWindowProc(hwnd, msg, wParam, lParam);
}

void invalidate(HWND hwnd){
    SendMessage(hwnd, WM_PAINT, 0, 0);
}