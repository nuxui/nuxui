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

#define printMsg(msgstr) printf("%s, msg=0x%X, wParam=0x%X, lParam=0x%X, W_high=0x%X, W_low=0x%X, L_high=0x%X, L_low=0x%X\n", (msgstr), msg, wParam, lParam, (SHORT)HIWORD(wParam), (SHORT)LOWORD(wParam), (SHORT)HIWORD(lParam), (SHORT)LOWORD(lParam));

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
{
    switch(msg)
    {
        /*mouse event begin*/
        case WM_NCHITTEST:
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_MOUSEMOVE:
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
            printMsg("WM_MOUSEBUTTON");
            go_mouseEvent(hwnd, msg, HIWORD(wParam), GET_X_LPARAM(lParam), GET_Y_LPARAM(lParam));
            break;
        }
        case WM_MOUSEWHEEL:
        {
            POINT p;
            p.x = GET_X_LPARAM(lParam);
            p.y = GET_Y_LPARAM(lParam);
            ScreenToClient(hwnd, &p);

            int aMouseInfo[3];
            int lines = 1;
            if ( SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &aMouseInfo, 0) > 0 ){
                lines = aMouseInfo[0]; // get lines number of scroll each time
            }
            go_scrollEvent(hwnd, (float)p.x, (float)p.y, 0, (GET_WHEEL_DELTA_WPARAM(wParam) * lines / (float) WHEEL_DELTA) );
            break;
        }
        case WM_MOUSEHWHEEL:
        {
            POINT p;
            p.x = GET_X_LPARAM(lParam);
            p.y = GET_Y_LPARAM(lParam);
            ScreenToClient(hwnd, &p);

            int aMouseInfo[3];
            int lines = 1;
            if ( SystemParametersInfoW(SPI_GETWHEELSCROLLLINES, 0, &aMouseInfo, 0) > 0 ){
                lines = aMouseInfo[0]; // get lines number of scroll each time
            }
            go_scrollEvent(hwnd, (float)p.x, (float)p.y, -(GET_WHEEL_DELTA_WPARAM(wParam) * lines / (float) WHEEL_DELTA), 0 );
            break;
        }
        case WM_KEYDOWN:
        case WM_KEYUP:
        case WM_SYSKEYDOWN:
        case WM_SYSKEYUP:
        {
            // printMsg("WM_KEY");
            go_keyEvent(hwnd, msg, (UINT32)LOWORD(wParam), 0, 0, NULL);
            break;

        }
        case WM_CHAR:
        {
            // printMsg("WM_CHAR");
            if((SHORT)HIWORD(lParam) != 0){  // 0 is IME char, ignore
                UINT32 keycode = (UINT32)LOWORD(wParam);
                if ( (keycode >= 0x20 && keycode <= 0x7E) || keycode == 0x09 || keycode == 0x0A || keycode == 0x0B || keycode == 0x0D ){
                    go_typeEvent(hwnd, msg, wParam, lParam);
                }
            }
            break;
        }
        case WM_SYSCHAR:
            printMsg("WM_SYSCHAR");
            break;
        case WM_DEADCHAR:
            printMsg("WM_DEADCHAR");
            break;
        case WM_SYSDEADCHAR:
            printMsg("WM_SYSDEADCHAR");
            break;
        case WM_UNICHAR:
            printMsg("WM_UNICHAR");
            break;
        case WM_IME_STARTCOMPOSITION:
            printMsg("WM_IME_STARTCOMPOSITION");
            break;
        case WM_IME_ENDCOMPOSITION:
            printMsg("WM_IME_ENDCOMPOSITION");
            break;
        case WM_IME_COMPOSITION:   // = WM_IME_KEYLAST
        {
            // printMsg("WM_IME_COMPOSITION");
            go_typeEvent(hwnd, msg, wParam, lParam);
            break;
        }
        case WM_IME_CHAR:
        {
            printMsg("WM_IME_CHAR");
            break;
        }
        case WM_IME_SETCONTEXT:
            printMsg("WM_IME_SETCONTEXT");
            break;
        case WM_IME_NOTIFY:
            printMsg("WM_IME_NOTIFY");
            break;
        case WM_IME_CONTROL:
            printMsg("WM_IME_CONTROL");
            break;
        case WM_IME_COMPOSITIONFULL:
            printMsg("WM_IME_COMPOSITIONFULL");
            break;
        case WM_IME_SELECT:
            printMsg("WM_IME_SELECT");
            break;
        case WM_IME_REQUEST:
            printMsg("WM_IME_REQUEST");
            break;
        case WM_IME_KEYDOWN:
            printMsg("WM_IME_KEYDOWN");
            break;
        case WM_IME_KEYUP:
            printMsg("WM_IME_KEYUP");
            break;
        case WM_CREATE:
            printMsg("WM_CREATE");
            go_windowAction(hwnd, msg);break;
        case WM_PAINT:
            printMsg("WM_PAINT");
            go_windowAction(hwnd, msg);break;
        case WM_SIZE:
            printMsg("WM_SIZE");
            go_windowAction(hwnd, msg);break;
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
            printMsg("WM_XBUTTONDBLCLK");
            return DefWindowProc(hwnd, msg, wParam, lParam);
        case WM_MOUSEACTIVATE:
        case WM_MOUSEHOVER:
        case WM_MOUSELEAVE:
            printMsg("WM_MOUSELEAVE");
            return DefWindowProc(hwnd, msg, wParam, lParam);
        /*mouse event end*/
        case WM_USER:
            go_backToUI();
            break;
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

void setTextInputRect(HWND hwnd, LONG x, LONG y, LONG w, LONG h)
{
    HIMC hIMC = ImmGetContext(hwnd);
    if (hIMC) {
        COMPOSITIONFORM comp;
        comp.dwStyle = CFS_POINT;
        comp.ptCurrentPos.x = x;
        comp.ptCurrentPos.y = y;
        ImmSetCompositionWindow(hIMC, &comp);
        ImmReleaseContext(hwnd, hIMC);
    }
}