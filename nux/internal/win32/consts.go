// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package win32

import (
	"syscall"
)

const (
	AlphaShift = 24
	RedShift   = 16
	GreenShift = 8
	BlueShift  = 0
)

/*
* Window field offsets for GetWindowLong()
 */
const (
	GWL_WNDPROC    = -4
	GWL_HINSTANCE  = -6
	GWL_HWNDPARENT = -8
	GWL_STYLE      = -16
	GWL_EXSTYLE    = -20
	GWL_USERDATA   = -21
	GWL_ID         = -12
)

const (
	WM_SETFOCUS        = 0x0007
	WM_KILLFOCUS       = 0x0008
	WM_ENABLE          = 0x000A
	WM_SETREDRAW       = 0x000B
	WM_SETTEXT         = 0x000C
	WM_GETTEXT         = 0x000D
	WM_GETTEXTLENGTH   = 0x000E
	WM_PAINT           = 0x000F
	WM_CLOSE           = 0x0010
	WM_QUERYENDSESSION = 0x0011
	WM_QUERYOPEN       = 0x0013
	WM_ENDSESSION      = 0x0016
	WM_QUIT            = 0x0012
	WM_ERASEBKGND      = 0x0014
	WM_SYSCOLORCHANGE  = 0x0015
	WM_SHOWWINDOW      = 0x0018
	WM_WININICHANGE    = 0x001A
	WM_SETTINGCHANGE   = WM_WININICHANGE

	WM_CUT               = 0x0300
	WM_COPY              = 0x0301
	WM_PASTE             = 0x0302
	WM_CLEAR             = 0x0303
	WM_UNDO              = 0x0304
	WM_RENDERFORMAT      = 0x0305
	WM_RENDERALLFORMATS  = 0x0306
	WM_DESTROYCLIPBOARD  = 0x0307
	WM_DRAWCLIPBOARD     = 0x0308
	WM_PAINTCLIPBOARD    = 0x0309
	WM_VSCROLLCLIPBOARD  = 0x030A
	WM_SIZECLIPBOARD     = 0x030B
	WM_ASKCBFORMATNAME   = 0x030C
	WM_CHANGECBCHAIN     = 0x030D
	WM_HSCROLLCLIPBOARD  = 0x030E
	WM_QUERYNEWPALETTE   = 0x030F
	WM_PALETTEISCHANGING = 0x0310
	WM_PALETTECHANGED    = 0x0311
	WM_HOTKEY            = 0x0312

	// Window Messages
	WM_NULL    = 0x0000
	WM_CREATE  = 0x0001
	WM_DESTROY = 0x0002
	WM_MOVE    = 0x0003
	WM_SIZE    = 0x0005
)

const (
	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
)

const (
	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002
)

/*
* Extended Window Styles
 */
const (
	WS_EX_DLGMODALFRAME       = 0x00000001
	WS_EX_NOPARENTNOTIFY      = 0x00000004
	WS_EX_TOPMOST             = 0x00000008
	WS_EX_ACCEPTFILES         = 0x00000010
	WS_EX_TRANSPARENT         = 0x00000020
	WS_EX_MDICHILD            = 0x00000040
	WS_EX_TOOLWINDOW          = 0x00000080
	WS_EX_WINDOWEDGE          = 0x00000100
	WS_EX_CLIENTEDGE          = 0x00000200
	WS_EX_CONTEXTHELP         = 0x00000400
	WS_EX_RIGHT               = 0x00001000
	WS_EX_LEFT                = 0x00000000
	WS_EX_RTLREADING          = 0x00002000
	WS_EX_LTRREADING          = 0x00000000
	WS_EX_LEFTSCROLLBAR       = 0x00004000
	WS_EX_RIGHTSCROLLBAR      = 0x00000000
	WS_EX_CONTROLPARENT       = 0x00010000
	WS_EX_STATICEDGE          = 0x00020000
	WS_EX_APPWINDOW           = 0x00040000
	WS_EX_OVERLAPPEDWINDOW    = (WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE)
	WS_EX_PALETTEWINDOW       = (WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST)
	WS_EX_LAYERED             = 0x00080000
	WS_EX_NOINHERITLAYOUT     = 0x00100000 // Disable inheritence of mirroring by children
	WS_EX_NOREDIRECTIONBITMAP = 0x00200000
	WS_EX_LAYOUTRTL           = 0x00400000 // Right to left mirroring
	WS_EX_COMPOSITED          = 0x02000000
	WS_EX_NOACTIVATE          = 0x08000000
)

const (
	VK_SHIFT   = 16
	VK_CONTROL = 17
	VK_MENU    = 18
	VK_LWIN    = 0x5B
	VK_RWIN    = 0x5C
)

const (
	MK_LBUTTON = 0x0001
	MK_MBUTTON = 0x0010
	MK_RBUTTON = 0x0002
)

const (
	COLOR_BTNFACE = 15
)

const (
	CW_USEDEFAULT = 0x80000000 - 0x100000000

	HWND_MESSAGE = syscall.Handle(^uintptr(2)) // -3

	SWP_NOSIZE = 0x0001
)

const (
	BI_RGB         = 0
	DIB_RGB_COLORS = 0

	AC_SRC_OVER  = 0x00
	AC_SRC_ALPHA = 0x01

	WHEEL_DELTA = 120
)

const (
	DCB_RESET      = 0x0001
	DCB_ACCUMULATE = 0x0002
	DCB_DIRTY      = DCB_ACCUMULATE
	DCB_SET        = (DCB_RESET | DCB_ACCUMULATE)
	DCB_ENABLE     = 0x0004
	DCB_DISABLE    = 0x0008
)

/*
* RedrawWindow() flags
 */
const (
	RDW_INVALIDATE      = 0x0001
	RDW_INTERNALPAINT   = 0x0002
	RDW_ERASE           = 0x0004
	RDW_VALIDATE        = 0x0008
	RDW_NOINTERNALPAINT = 0x0010
	RDW_NOERASE         = 0x0020
	RDW_NOCHILDREN      = 0x0040
	RDW_ALLCHILDREN     = 0x0080
	RDW_UPDATENOW       = 0x0100
	RDW_ERASENOW        = 0x0200
	RDW_FRAME           = 0x0400
	RDW_NOFRAME         = 0x0800
)

/*
 * Standard Icon IDs
 */
const (
	IDI_APPLICATION = 32512
	IDI_HAND        = 32513
	IDI_QUESTION    = 32514
	IDI_EXCLAMATION = 32515
	IDI_ASTERISK    = 32516
	IDI_WINLOGO     = 32517
	IDI_SHIELD      = 32518
	IDI_WARNING     = IDI_EXCLAMATION
	IDI_ERROR       = IDI_HAND
	IDI_INFORMATION = IDI_ASTERISK
)

// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.10240.0/um/WinUser.h
const (
	IDC_ARROW       = 32512
	IDC_IBEAM       = 32513
	IDC_WAIT        = 32514
	IDC_CROSS       = 32515
	IDC_UPARROW     = 32516
	IDC_SIZE        = 32640
	IDC_ICON        = 32641
	IDC_SIZENWSE    = 32642
	IDC_SIZENESW    = 32643
	IDC_SIZEWE      = 32644
	IDC_SIZENS      = 32645
	IDC_SIZEALL     = 32646
	IDC_NO          = 32648
	IDC_HAND        = 32649
	IDC_APPSTARTING = 32650
	IDC_HELP        = 32651
)

/*
 * ShowWindow() Commands
 */
const (
	SW_HIDE            = 0
	SW_SHOWNORMAL      = 1
	SW_NORMAL          = 1
	SW_SHOWMINIMIZED   = 2
	SW_SHOWMAXIMIZED   = 3
	SW_MAXIMIZE        = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11
	SW_MAX             = 11
)

/*
 * WM_KEYUP/DOWN/CHAR HIWORD(lParam) flags
 */
const (
	KF_EXTENDED = 0x0100
	KF_DLGMODE  = 0x0800
	KF_MENUMODE = 0x1000
	KF_ALTDOWN  = 0x2000
	KF_REPEAT   = 0x4000
	KF_UP       = 0x8000
)

/*
 * Virtual Keys, Standard Set
 */
const (
	VK_LBUTTON  = 0x01
	VK_RBUTTON  = 0x02
	VK_CANCEL   = 0x03
	VK_MBUTTON  = 0x04 /* NOT contiguous with L & RBUTTON */
	VK_XBUTTON1 = 0x05 /* NOT contiguous with L & RBUTTON */
	VK_XBUTTON2 = 0x06 /* NOT contiguous with L & RBUTTON */
)

const (
	WM_USER                   = 0x0400
	WM_NOTIFY                 = 0x004E
	WM_INPUTLANGCHANGEREQUEST = 0x0050
	WM_INPUTLANGCHANGE        = 0x0051
	WM_TCARD                  = 0x0052
	WM_HELP                   = 0x0053
	WM_USERCHANGED            = 0x0054
	WM_NOTIFYFORMAT           = 0x0055
	NFR_ANSI                  = 1
	NFR_UNICODE               = 2
	NF_QUERY                  = 3
	NF_REQUERY                = 4
	WM_CONTEXTMENU            = 0x007B
	WM_STYLECHANGING          = 0x007C
	WM_STYLECHANGED           = 0x007D
	WM_DISPLAYCHANGE          = 0x007E
	WM_GETICON                = 0x007F
	WM_SETICON                = 0x0080
)

const (
	WM_CTLCOLORMSGBOX    = 0x0132
	WM_CTLCOLOREDIT      = 0x0133
	WM_CTLCOLORLISTBOX   = 0x0134
	WM_CTLCOLORBTN       = 0x0135
	WM_CTLCOLORDLG       = 0x0136
	WM_CTLCOLORSCROLLBAR = 0x0137
	WM_CTLCOLORSTATIC    = 0x0138
	MN_GETHMENU          = 0x01E1
)

const (
	WM_MOUSEFIRST    = 0x0200
	WM_MOUSEMOVE     = 0x0200
	WM_LBUTTONDOWN   = 0x0201
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONDOWN   = 0x0204
	WM_RBUTTONUP     = 0x0205
	WM_RBUTTONDBLCLK = 0x0206
	WM_MBUTTONDOWN   = 0x0207
	WM_MBUTTONUP     = 0x0208
	WM_MBUTTONDBLCLK = 0x0209
	WM_MOUSEWHEEL    = 0x020A
	WM_XBUTTONDOWN   = 0x020B
	WM_XBUTTONUP     = 0x020C
	WM_XBUTTONDBLCLK = 0x020D
	WM_MOUSEHWHEEL   = 0x020E
)

const (
	WM_KEYFIRST    = 0x0100
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_CHAR        = 0x0102
	WM_DEADCHAR    = 0x0103
	WM_SYSKEYDOWN  = 0x0104
	WM_SYSKEYUP    = 0x0105
	WM_SYSCHAR     = 0x0106
	WM_SYSDEADCHAR = 0x0107
	WM_UNICHAR     = 0x0109
	UNICODE_NOCHAR = 0xFFFF
)

const (
	WM_IME_STARTCOMPOSITION = 0x010D
	WM_IME_ENDCOMPOSITION   = 0x010E
	WM_IME_COMPOSITION      = 0x010F
	WM_IME_KEYLAST          = 0x010F
)

const (
	WM_DEVMODECHANGE = 0x001B
	WM_ACTIVATEAPP   = 0x001C
	WM_FONTCHANGE    = 0x001D
	WM_TIMECHANGE    = 0x001E
	WM_CANCELMODE    = 0x001F
	WM_SETCURSOR     = 0x0020
	WM_MOUSEACTIVATE = 0x0021
	WM_CHILDACTIVATE = 0x0022
	WM_QUEUESYNC     = 0x0023
	WM_GETMINMAXINFO = 0x0024
)

const (
	WM_MOUSEHOVER      = 0x02A1
	WM_MOUSELEAVE      = 0x02A3
	WM_NCMOUSEHOVER    = 0x02A0
	WM_NCMOUSELEAVE    = 0x02A2
	WM_NCCREATE        = 0x0081
	WM_NCDESTROY       = 0x0082
	WM_NCCALCSIZE      = 0x0083
	WM_NCHITTEST       = 0x0084
	WM_NCPAINT         = 0x0085
	WM_NCACTIVATE      = 0x0086
	WM_GETDLGCODE      = 0x0087
	WM_SYNCPAINT       = 0x0088
	WM_NCMOUSEMOVE     = 0x00A0
	WM_NCLBUTTONDOWN   = 0x00A1
	WM_NCLBUTTONUP     = 0x00A2
	WM_NCLBUTTONDBLCLK = 0x00A3
	WM_NCRBUTTONDOWN   = 0x00A4
	WM_NCRBUTTONUP     = 0x00A5
	WM_NCRBUTTONDBLCLK = 0x00A6
	WM_NCMBUTTONDOWN   = 0x00A7
	WM_NCMBUTTONUP     = 0x00A8
	WM_NCMBUTTONDBLCLK = 0x00A9
	WM_NCXBUTTONDOWN   = 0x00AB
	WM_NCXBUTTONUP     = 0x00AC
	WM_NCXBUTTONDBLCLK = 0x00AD
)

const (
	WM_NEXTMENU       = 0x0213
	WM_SIZING         = 0x0214
	WM_CAPTURECHANGED = 0x0215
	WM_MOVING         = 0x0216
)

// Status
const (
	Ok                        GpStatus = 0
	GenericError              GpStatus = 1
	InvalidParameter          GpStatus = 2
	OutOfMemory               GpStatus = 3
	ObjectBusy                GpStatus = 4
	InsufficientBuffer        GpStatus = 5
	NotImplemented            GpStatus = 6
	Win32Error                GpStatus = 7
	WrongState                GpStatus = 8
	Aborted                   GpStatus = 9
	FileNotFound              GpStatus = 10
	ValueOverflow             GpStatus = 11
	AccessDenied              GpStatus = 12
	UnknownImageFormat        GpStatus = 13
	FontFamilyNotFound        GpStatus = 14
	FontStyleNotFound         GpStatus = 15
	NotTrueTypeFont           GpStatus = 16
	UnsupportedGdiplusVersion GpStatus = 17
	GdiplusNotInitialized     GpStatus = 18
	PropertyNotFound          GpStatus = 19
	PropertyNotSupported      GpStatus = 20
	ProfileNotFound           GpStatus = 21
)

const (
	CombineModeReplace    GpCombineMode = 0
	CombineModeIntersect  GpCombineMode = 1
	CombineModeUnion      GpCombineMode = 2
	CombineModeXor        GpCombineMode = 3
	CombineModeExclude    GpCombineMode = 4
	CombineModeComplement GpCombineMode = 5
)

/*
* Parameter for SystemParametersInfo.
 */
const (
	SPI_GETBEEP                     = 0x0001
	SPI_SETBEEP                     = 0x0002
	SPI_GETMOUSE                    = 0x0003
	SPI_SETMOUSE                    = 0x0004
	SPI_GETBORDER                   = 0x0005
	SPI_SETBORDER                   = 0x0006
	SPI_GETKEYBOARDSPEED            = 0x000A
	SPI_SETKEYBOARDSPEED            = 0x000B
	SPI_LANGDRIVER                  = 0x000C
	SPI_ICONHORIZONTALSPACING       = 0x000D
	SPI_GETSCREENSAVETIMEOUT        = 0x000E
	SPI_SETSCREENSAVETIMEOUT        = 0x000F
	SPI_GETSCREENSAVEACTIVE         = 0x0010
	SPI_SETSCREENSAVEACTIVE         = 0x0011
	SPI_GETGRIDGRANULARITY          = 0x0012
	SPI_SETGRIDGRANULARITY          = 0x0013
	SPI_SETDESKWALLPAPER            = 0x0014
	SPI_SETDESKPATTERN              = 0x0015
	SPI_GETKEYBOARDDELAY            = 0x0016
	SPI_SETKEYBOARDDELAY            = 0x0017
	SPI_ICONVERTICALSPACING         = 0x0018
	SPI_GETICONTITLEWRAP            = 0x0019
	SPI_SETICONTITLEWRAP            = 0x001A
	SPI_GETMENUDROPALIGNMENT        = 0x001B
	SPI_SETMENUDROPALIGNMENT        = 0x001C
	SPI_SETDOUBLECLKWIDTH           = 0x001D
	SPI_SETDOUBLECLKHEIGHT          = 0x001E
	SPI_GETICONTITLELOGFONT         = 0x001F
	SPI_SETDOUBLECLICKTIME          = 0x0020
	SPI_SETMOUSEBUTTONSWAP          = 0x0021
	SPI_SETICONTITLELOGFONT         = 0x0022
	SPI_GETFASTTASKSWITCH           = 0x0023
	SPI_SETFASTTASKSWITCH           = 0x0024
	SPI_SETDRAGFULLWINDOWS          = 0x0025
	SPI_GETDRAGFULLWINDOWS          = 0x0026
	SPI_GETNONCLIENTMETRICS         = 0x0029
	SPI_SETNONCLIENTMETRICS         = 0x002A
	SPI_GETMINIMIZEDMETRICS         = 0x002B
	SPI_SETMINIMIZEDMETRICS         = 0x002C
	SPI_GETICONMETRICS              = 0x002D
	SPI_SETICONMETRICS              = 0x002E
	SPI_SETWORKAREA                 = 0x002F
	SPI_GETWORKAREA                 = 0x0030
	SPI_SETPENWINDOWS               = 0x0031
	SPI_GETHIGHCONTRAST             = 0x0042
	SPI_SETHIGHCONTRAST             = 0x0043
	SPI_GETKEYBOARDPREF             = 0x0044
	SPI_SETKEYBOARDPREF             = 0x0045
	SPI_GETSCREENREADER             = 0x0046
	SPI_SETSCREENREADER             = 0x0047
	SPI_GETANIMATION                = 0x0048
	SPI_SETANIMATION                = 0x0049
	SPI_GETFONTSMOOTHING            = 0x004A
	SPI_SETFONTSMOOTHING            = 0x004B
	SPI_SETDRAGWIDTH                = 0x004C
	SPI_SETDRAGHEIGHT               = 0x004D
	SPI_SETHANDHELD                 = 0x004E
	SPI_GETLOWPOWERTIMEOUT          = 0x004F
	SPI_GETPOWEROFFTIMEOUT          = 0x0050
	SPI_SETLOWPOWERTIMEOUT          = 0x0051
	SPI_SETPOWEROFFTIMEOUT          = 0x0052
	SPI_GETLOWPOWERACTIVE           = 0x0053
	SPI_GETPOWEROFFACTIVE           = 0x0054
	SPI_SETLOWPOWERACTIVE           = 0x0055
	SPI_SETPOWEROFFACTIVE           = 0x0056
	SPI_SETCURSORS                  = 0x0057
	SPI_SETICONS                    = 0x0058
	SPI_GETDEFAULTINPUTLANG         = 0x0059
	SPI_SETDEFAULTINPUTLANG         = 0x005A
	SPI_SETLANGTOGGLE               = 0x005B
	SPI_GETWINDOWSEXTENSION         = 0x005C
	SPI_SETMOUSETRAILS              = 0x005D
	SPI_GETMOUSETRAILS              = 0x005E
	SPI_SETSCREENSAVERRUNNING       = 0x0061
	SPI_SCREENSAVERRUNNING          = SPI_SETSCREENSAVERRUNNING
	SPI_GETFILTERKEYS               = 0x0032
	SPI_SETFILTERKEYS               = 0x0033
	SPI_GETTOGGLEKEYS               = 0x0034
	SPI_SETTOGGLEKEYS               = 0x0035
	SPI_GETMOUSEKEYS                = 0x0036
	SPI_SETMOUSEKEYS                = 0x0037
	SPI_GETSHOWSOUNDS               = 0x0038
	SPI_SETSHOWSOUNDS               = 0x0039
	SPI_GETSTICKYKEYS               = 0x003A
	SPI_SETSTICKYKEYS               = 0x003B
	SPI_GETACCESSTIMEOUT            = 0x003C
	SPI_SETACCESSTIMEOUT            = 0x003D
	SPI_GETSERIALKEYS               = 0x003E
	SPI_SETSERIALKEYS               = 0x003F
	SPI_GETSOUNDSENTRY              = 0x0040
	SPI_SETSOUNDSENTRY              = 0x0041
	SPI_GETSNAPTODEFBUTTON          = 0x005F
	SPI_SETSNAPTODEFBUTTON          = 0x0060
	SPI_GETMOUSEHOVERWIDTH          = 0x0062
	SPI_SETMOUSEHOVERWIDTH          = 0x0063
	SPI_GETMOUSEHOVERHEIGHT         = 0x0064
	SPI_SETMOUSEHOVERHEIGHT         = 0x0065
	SPI_GETMOUSEHOVERTIME           = 0x0066
	SPI_SETMOUSEHOVERTIME           = 0x0067
	SPI_GETWHEELSCROLLLINES         = 0x0068
	SPI_SETWHEELSCROLLLINES         = 0x0069
	SPI_GETMENUSHOWDELAY            = 0x006A
	SPI_SETMENUSHOWDELAY            = 0x006B
	SPI_GETWHEELSCROLLCHARS         = 0x006C
	SPI_SETWHEELSCROLLCHARS         = 0x006D
	SPI_GETSHOWIMEUI                = 0x006E
	SPI_SETSHOWIMEUI                = 0x006F
	SPI_GETMOUSESPEED               = 0x0070
	SPI_SETMOUSESPEED               = 0x0071
	SPI_GETSCREENSAVERRUNNING       = 0x0072
	SPI_GETDESKWALLPAPER            = 0x0073
	SPI_GETAUDIODESCRIPTION         = 0x0074
	SPI_SETAUDIODESCRIPTION         = 0x0075
	SPI_GETSCREENSAVESECURE         = 0x0076
	SPI_SETSCREENSAVESECURE         = 0x0077
	SPI_GETHUNGAPPTIMEOUT           = 0x0078
	SPI_SETHUNGAPPTIMEOUT           = 0x0079
	SPI_GETWAITTOKILLTIMEOUT        = 0x007A
	SPI_SETWAITTOKILLTIMEOUT        = 0x007B
	SPI_GETWAITTOKILLSERVICETIMEOUT = 0x007C
	SPI_SETWAITTOKILLSERVICETIMEOUT = 0x007D
	SPI_GETMOUSEDOCKTHRESHOLD       = 0x007E
	SPI_SETMOUSEDOCKTHRESHOLD       = 0x007F
	SPI_GETPENDOCKTHRESHOLD         = 0x0080
	SPI_SETPENDOCKTHRESHOLD         = 0x0081
	SPI_GETWINARRANGING             = 0x0082
	SPI_SETWINARRANGING             = 0x0083
	SPI_GETMOUSEDRAGOUTTHRESHOLD    = 0x0084
	SPI_SETMOUSEDRAGOUTTHRESHOLD    = 0x0085
	SPI_GETPENDRAGOUTTHRESHOLD      = 0x0086
	SPI_SETPENDRAGOUTTHRESHOLD      = 0x0087
	SPI_GETMOUSESIDEMOVETHRESHOLD   = 0x0088
	SPI_SETMOUSESIDEMOVETHRESHOLD   = 0x0089
	SPI_GETPENSIDEMOVETHRESHOLD     = 0x008A
	SPI_SETPENSIDEMOVETHRESHOLD     = 0x008B
	SPI_GETDRAGFROMMAXIMIZE         = 0x008C
	SPI_SETDRAGFROMMAXIMIZE         = 0x008D
	SPI_GETSNAPSIZING               = 0x008E
	SPI_SETSNAPSIZING               = 0x008F
	SPI_GETDOCKMOVING               = 0x0090
	SPI_SETDOCKMOVING               = 0x0091
)

type FontStyle int32

const (
	FontStyleRegular    FontStyle = 0
	FontStyleBold       FontStyle = 1
	FontStyleItalic     FontStyle = 2
	FontStyleBoldItalic FontStyle = 3
	FontStyleUnderline  FontStyle = 4
	FontStyleStrikeout  FontStyle = 8
)

/*
 * GetSystemMetrics() codes
 */
type SMCode int

const (
	SM_CXSCREEN                    SMCode = 0
	SM_CYSCREEN                    SMCode = 1
	SM_CXVSCROLL                   SMCode = 2
	SM_CYHSCROLL                   SMCode = 3
	SM_CYCAPTION                   SMCode = 4
	SM_CXBORDER                    SMCode = 5
	SM_CYBORDER                    SMCode = 6
	SM_CXDLGFRAME                  SMCode = 7
	SM_CYDLGFRAME                  SMCode = 8
	SM_CYVTHUMB                    SMCode = 9
	SM_CXHTHUMB                    SMCode = 10
	SM_CXICON                      SMCode = 11
	SM_CYICON                      SMCode = 12
	SM_CXCURSOR                    SMCode = 13
	SM_CYCURSOR                    SMCode = 14
	SM_CYMENU                      SMCode = 15
	SM_CXFULLSCREEN                SMCode = 16
	SM_CYFULLSCREEN                SMCode = 17
	SM_CYKANJIWINDOW               SMCode = 18
	SM_MOUSEPRESENT                SMCode = 19
	SM_CYVSCROLL                   SMCode = 20
	SM_CXHSCROLL                   SMCode = 21
	SM_DEBUG                       SMCode = 22
	SM_SWAPBUTTON                  SMCode = 23
	SM_RESERVED1                   SMCode = 24
	SM_RESERVED2                   SMCode = 25
	SM_RESERVED3                   SMCode = 26
	SM_RESERVED4                   SMCode = 27
	SM_CXMIN                       SMCode = 28
	SM_CYMIN                       SMCode = 29
	SM_CXSIZE                      SMCode = 30
	SM_CYSIZE                      SMCode = 31
	SM_CXFRAME                     SMCode = 32
	SM_CYFRAME                     SMCode = 33
	SM_CXMINTRACK                  SMCode = 34
	SM_CYMINTRACK                  SMCode = 35
	SM_CXDOUBLECLK                 SMCode = 36
	SM_CYDOUBLECLK                 SMCode = 37
	SM_CXICONSPACING               SMCode = 38
	SM_CYICONSPACING               SMCode = 39
	SM_MENUDROPALIGNMENT           SMCode = 40
	SM_PENWINDOWS                  SMCode = 41
	SM_DBCSENABLED                 SMCode = 42
	SM_CMOUSEBUTTONS               SMCode = 43
	SM_CXFIXEDFRAME                SMCode = SM_CXDLGFRAME /* ;win40 name change */
	SM_CYFIXEDFRAME                SMCode = SM_CYDLGFRAME /* ;win40 name change */
	SM_CXSIZEFRAME                 SMCode = SM_CXFRAME    /* ;win40 name change */
	SM_CYSIZEFRAME                 SMCode = SM_CYFRAME    /* ;win40 name change */
	SM_SECURE                      SMCode = 44
	SM_CXEDGE                      SMCode = 45
	SM_CYEDGE                      SMCode = 46
	SM_CXMINSPACING                SMCode = 47
	SM_CYMINSPACING                SMCode = 48
	SM_CXSMICON                    SMCode = 49
	SM_CYSMICON                    SMCode = 50
	SM_CYSMCAPTION                 SMCode = 51
	SM_CXSMSIZE                    SMCode = 52
	SM_CYSMSIZE                    SMCode = 53
	SM_CXMENUSIZE                  SMCode = 54
	SM_CYMENUSIZE                  SMCode = 55
	SM_ARRANGE                     SMCode = 56
	SM_CXMINIMIZED                 SMCode = 57
	SM_CYMINIMIZED                 SMCode = 58
	SM_CXMAXTRACK                  SMCode = 59
	SM_CYMAXTRACK                  SMCode = 60
	SM_CXMAXIMIZED                 SMCode = 61
	SM_CYMAXIMIZED                 SMCode = 62
	SM_NETWORK                     SMCode = 63
	SM_CLEANBOOT                   SMCode = 67
	SM_CXDRAG                      SMCode = 68
	SM_CYDRAG                      SMCode = 69
	SM_SHOWSOUNDS                  SMCode = 70
	SM_CXMENUCHECK                 SMCode = 71 /* Use instead of GetMenuCheckMarkDimensions()! */
	SM_CYMENUCHECK                 SMCode = 72
	SM_SLOWMACHINE                 SMCode = 73
	SM_MIDEASTENABLED              SMCode = 74
	SM_MOUSEWHEELPRESENT           SMCode = 75
	SM_XVIRTUALSCREEN              SMCode = 76
	SM_YVIRTUALSCREEN              SMCode = 77
	SM_CXVIRTUALSCREEN             SMCode = 78
	SM_CYVIRTUALSCREEN             SMCode = 79
	SM_CMONITORS                   SMCode = 80
	SM_SAMEDISPLAYFORMAT           SMCode = 81
	SM_IMMENABLED                  SMCode = 82
	SM_CXFOCUSBORDER               SMCode = 83
	SM_CYFOCUSBORDER               SMCode = 84
	SM_TABLETPC                    SMCode = 86
	SM_MEDIACENTER                 SMCode = 87
	SM_STARTER                     SMCode = 88
	SM_SERVERR2                    SMCode = 89
	SM_MOUSEHORIZONTALWHEELPRESENT SMCode = 91
	SM_CXPADDEDBORDER              SMCode = 92
	SM_DIGITIZER                   SMCode = 94
	SM_MAXIMUMTOUCHES              SMCode = 95
	SM_REMOTESESSION               SMCode = 0x1000
	SM_SHUTTINGDOWN                SMCode = 0x2000
	SM_REMOTECONTROL               SMCode = 0x2001
	SM_CARETBLINKINGENABLED        SMCode = 0x2002
	SM_CONVERTIBLESLATEMODE        SMCode = 0x2003
	SM_SYSTEMDOCKED                SMCode = 0x2004
)

const (
	HTERROR       = (-2)
	HTTRANSPARENT = (-1)
	HTNOWHERE     = 0
	HTCLIENT      = 1
	HTCAPTION     = 2
	HTSYSMENU     = 3
	HTGROWBOX     = 4
	HTSIZE        = HTGROWBOX
	HTMENU        = 5
	HTHSCROLL     = 6
	HTVSCROLL     = 7
	HTMINBUTTON   = 8
	HTMAXBUTTON   = 9
	HTLEFT        = 10
	HTRIGHT       = 11
	HTTOP         = 12
	HTTOPLEFT     = 13
	HTTOPRIGHT    = 14
	HTBOTTOM      = 15
	HTBOTTOMLEFT  = 16
	HTBOTTOMRIGHT = 17
	HTBORDER      = 18
	HTREDUCE      = HTMINBUTTON
	HTZOOM        = HTMAXBUTTON
	HTSIZEFIRST   = HTLEFT
	HTSIZELAST    = HTBOTTOMRIGHT
	HTOBJECT      = 19
	HTCLOSE       = 20
	HTHELP        = 21
)

// https://github.com/tpn/winsdk-10/blob/master/Include/10.0.14393.0/um/commdlg.h
const (
	OFN_READONLY             = 0x00000001
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_HIDEREADONLY         = 0x00000004
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_SHOWHELP             = 0x00000010
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_NOVALIDATE           = 0x00000100
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_CREATEPROMPT         = 0x00002000
	OFN_SHAREAWARE           = 0x00004000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOLONGNAMES          = 0x00040000 // force no long names for 4.x modules
	OFN_EXPLORER             = 0x00080000 // new look commdlg
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_LONGNAMES            = 0x00200000 // force long names for 3.x modules
	OFN_ENABLEINCLUDENOTIFY  = 0x00400000 // send include message to callback
	OFN_ENABLESIZING         = 0x00800000
	OFN_DONTADDTORECENT      = 0x02000000
	OFN_FORCESHOWHIDDEN      = 0x10000000 // Show All files including System and hidden files
	OFN_EX_NOPLACESBAR       = 0x00000001
	OFN_SHAREFALLTHROUGH     = 2
	OFN_SHARENOWARN          = 1
	OFN_SHAREWARN            = 0
)

//SHBrowseForFolder flags
const (
	BIF_RETURNONLYFSDIRS    = 0x00000001
	BIF_DONTGOBELOWDOMAIN   = 0x00000002
	BIF_STATUSTEXT          = 0x00000004
	BIF_RETURNFSANCESTORS   = 0x00000008
	BIF_EDITBOX             = 0x00000010
	BIF_VALIDATE            = 0x00000020
	BIF_NEWDIALOGSTYLE      = 0x00000040
	BIF_BROWSEINCLUDEURLS   = 0x00000080
	BIF_USENEWUI            = BIF_EDITBOX | BIF_NEWDIALOGSTYLE
	BIF_UAHINT              = 0x00000100
	BIF_NONEWFOLDERBUTTON   = 0x00000200
	BIF_NOTRANSLATETARGETS  = 0x00000400
	BIF_BROWSEFORCOMPUTER   = 0x00001000
	BIF_BROWSEFORPRINTER    = 0x00002000
	BIF_BROWSEINCLUDEFILES  = 0x00004000
	BIF_SHAREABLE           = 0x00008000
	BIF_BROWSEFILEJUNCTIONS = 0x00010000
)

//MessageBox flags
const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND
	MB_DEFBUTTON1        = 0x00000000
	MB_DEFBUTTON2        = 0x00000100
	MB_DEFBUTTON3        = 0x00000200
	MB_DEFBUTTON4        = 0x00000300
)

// ------------------------- GDI+ -------------------

const LANG_NEUTRAL = 0x00

// Unit
const (
	UnitWorld      = 0 // 0 -- World coordinate (non-physical unit)
	UnitDisplay    = 1 // 1 -- Variable -- for PageTransform only
	UnitPixel      = 2 // 2 -- Each unit is one device pixel.
	UnitPoint      = 3 // 3 -- Each unit is a printer's point, or 1/72 inch.
	UnitInch       = 4 // 4 -- Each unit is 1 inch.
	UnitDocument   = 5 // 5 -- Each unit is 1/300 inch.
	UnitMillimeter = 6 // 6 -- Each unit is 1 millimeter.
)

const (
	AlphaMask = 0xff000000
	RedMask   = 0x00ff0000
	GreenMask = 0x0000ff00
	BlueMask  = 0x000000ff
)

// QualityMode
const (
	QualityModeInvalid = iota - 1
	QualityModeDefault
	QualityModeLow  // Best performance
	QualityModeHigh // Best rendering quality
)

// Alpha Compositing mode
const (
	CompositingModeSourceOver = iota // 0
	CompositingModeSourceCopy        // 1
)

// Alpha Compositing quality
const (
	CompositingQualityInvalid = iota + QualityModeInvalid
	CompositingQualityDefault
	CompositingQualityHighSpeed
	CompositingQualityHighQuality
	CompositingQualityGammaCorrected
	CompositingQualityAssumeLinear
)

// InterpolationMode
const (
	InterpolationModeInvalid = iota + QualityModeInvalid
	InterpolationModeDefault
	InterpolationModeLowQuality
	InterpolationModeHighQuality
	InterpolationModeBilinear
	InterpolationModeBicubic
	InterpolationModeNearestNeighbor
	InterpolationModeHighQualityBilinear
	InterpolationModeHighQualityBicubic
)

type GpSmoothingMode int32

const (
	SmoothingModeInvalid      GpSmoothingMode = QualityModeInvalid
	SmoothingModeDefault      GpSmoothingMode = 0
	SmoothingModeHighSpeed    GpSmoothingMode = 1
	SmoothingModeHighQuality  GpSmoothingMode = 2
	SmoothingModeNone         GpSmoothingMode = 3
	SmoothingModeAntiAlias8x4 GpSmoothingMode = 4
	SmoothingModeAntiAlias    GpSmoothingMode = 4
	SmoothingModeAntiAlias8x8 GpSmoothingMode = 5
)

type GpFlushIntention int32

const (
	FlushIntentionFlush GpFlushIntention = 0
	FlushIntentionSync  GpFlushIntention = 1
)

// Pixel Format Mode
const (
	PixelOffsetModeInvalid = iota + QualityModeInvalid
	PixelOffsetModeDefault
	PixelOffsetModeHighSpeed
	PixelOffsetModeHighQuality
	PixelOffsetModeNone // No pixel offset
	PixelOffsetModeHalf // Offset by -0.5, -0.5 for fast anti-alias perf
)

// Text Rendering Hint
const (
	TextRenderingHintSystemDefault            = iota // Glyph with system default rendering hint
	TextRenderingHintSingleBitPerPixelGridFit        // Glyph bitmap with hinting
	TextRenderingHintSingleBitPerPixel               // Glyph bitmap without hinting
	TextRenderingHintAntiAliasGridFit                // Glyph anti-alias bitmap with hinting
	TextRenderingHintAntiAlias                       // Glyph anti-alias bitmap without hinting
	TextRenderingHintClearTypeGridFit                // Glyph CT bitmap with hinting
)

// BrushType
const (
	BrushTypeSolidColor GpBrushType = iota
	BrushTypeHatchFill
	BrushTypeTextureFill
	BrushTypePathGradient
	BrushTypeLinearGradient
)

// LineCap
const (
	LineCapFlat GpLineCap = iota
	LineCapSquare
	LineCapRound
	LineCapTriangle
	LineCapNoAnchor
	LineCapSquareAnchor
	LineCapRoundAnchor
	LineCapDiamondAnchor
	LineCapArrowAnchor
	LineCapCustom
	LineCapAnchorMask
)

// LineJoin
const (
	LineJoinMiter GpLineJoin = iota
	LineJoinBevel
	LineJoinRound
	LineJoinMiterClipped
)

// DashCap
const (
	DashCapFlat GpDashCap = iota
	DashCapRound
	DashCapTriangle
)

// DashStyle
const (
	DashStyleSolid GpDashStyle = iota
	DashStyleDash
	DashStyleDot
	DashStyleDashDot
	DashStyleDashDotDot
	DashStyleCustom
)

// PenAlignment
const (
	PenAlignmentCenter GpPenAlignment = iota
	PenAlignmentInset
)

// MatrixOrder
const (
	MatrixOrderPrepend GpMatrixOrder = iota
	MatrixOrderAppend
)

// PenType
const (
	PenTypeSolidColor GpPenType = iota
	PenTypeHatchFill
	PenTypeTextureFill
	PenTypePathGradient
	PenTypeLinearGradient
	PenTypeUnknown
)

type GpFillMode int

const (
	FillModeAlternate GpFillMode = iota
	FillModeWinding
)
