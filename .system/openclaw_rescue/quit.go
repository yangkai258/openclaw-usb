package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
	"time"
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
	MB_ICONWARNING  = 0x00000030
)

func winMsgBox(title, msg string, flags uintptr) int {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(msg)
	ret, _, _ := messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), flags)
	return int(ret)
}

func killProcess(name string) error {
	cmd := exec.Command("taskkill", "/F", "/IM", name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func isProcessRunning(name string) bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+name)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, _ := cmd.Output()
	return strings.Contains(string(out), name)
}

func main() {
	// Step 1: Confirmation
	ret := winMsgBox("OpenClaw",
		"确定要完全退出 OpenClaw 吗？\n这将关闭所有相关进程。",
		MB_OK|MB_YESNO|MB_ICONINFO|MB_SYSTEMMODAL)

	if ret != 6 { // IDYES = 6
		return
	}

	// Step 2: Kill all OpenClaw related processes
	procs := []string{"node.exe", "openclaw.exe"}
	for _, name := range procs {
		killProcess(name)
	}

	// Step 3: Wait a moment for processes to fully terminate
	time.Sleep(800 * time.Millisecond)

	// Step 4: Check if any node.exe is still running
	stillRunning := isProcessRunning("node.exe")
	if stillRunning {
		winMsgBox("OpenClaw",
			"注意：仍有进程在运行。\n请手动关闭后再尝试拔U盘。",
			MB_OK|MB_ICONWARNING|MB_SYSTEMMODAL)
	} else {
		// Step 5: Safe to eject
		winMsgBox("OpenClaw",
			"✅ 所有进程已关闭。\n\n现在可以安全拔出U盘了。",
			MB_OK|MB_ICONINFO|MB_SYSTEMMODAL)
	}

	fmt.Println("All processes terminated.")
}