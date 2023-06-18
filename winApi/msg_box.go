package winApi

import (
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	messageBox              = user32.NewProc("MessageBoxW")
	mbYesNo         uintptr = 0x00000004
	mbIconQuestion  uintptr = 0x00000020
	mbDefButton2    uintptr = 0x00000100
	MSG_BOX_YES_BTN int     = 1
	MSG_BOX_NO_BTN  int     = 0
)

func MessageBox(title string, text string) int {
	lpText, _ := syscall.UTF16PtrFromString(text)
	lpCaption, _ := syscall.UTF16PtrFromString(title)
	ret, _, _ := messageBox.Call(
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		mbYesNo|mbIconQuestion|mbDefButton2)
	if int(ret) == 6 { // Yes button
		return MSG_BOX_YES_BTN
	} else {
		return MSG_BOX_NO_BTN
	}
	return MSG_BOX_NO_BTN
}
