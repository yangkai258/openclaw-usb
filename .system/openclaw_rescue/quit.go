package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	user32      = syscall.NewLazyDLL("user32.dll")
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	messageBoxW = user32.NewProc("MessageBoxW")
)

const (
	MB_OK           = 0x00000000
	MB_ICONINFO     = 0x00000040
	MB_SYSTEMMODAL  = 0x00001000
	MB_YESNO        = 0x00000004
)

func winMsgBox(title, msg string, flags uintptr) int {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(msg)
	ret, _, _ := messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), flags)
	return int(ret)
}

func getExePath() string {
	var m = make([]uint16, syscall.MAX_PATH)
	n, _, _ := kernel32.NewProc("GetModuleFileNameW").Call(
		0, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)))
	return syscall.UTF16ToString(m[:n])
}

func killOpenClawProcesses() error {
	procs := []string{"node.exe"}
	for _, name := range procs {
		cmd := exec.Command("taskkill", "/F", "/IM", name)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
	}
	return nil
}

func main() {
	ret := winMsgBox("OpenClaw",
		"确定要完全退出 OpenClaw 吗？\n这将关闭所有相关进程。",
		MB_OK|MB_YESNO|MB_ICONINFO|MB_SYSTEMMODAL)

	if ret != 6 { // IDYES = 6
		return
	}

	killOpenClawProcesses()
	winMsgBox("OpenClaw", "已完全退出。", MB_OK|MB_ICONINFO|MB_SYSTEMMODAL)
	fmt.Println("All processes terminated.")
}