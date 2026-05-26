package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32      = syscall.NewLazyDLL("user32.dll")
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	messageBoxW = user32.NewProc("MessageBoxW")
)

const (
	MB_OK           = 0x00000000
	MB_ICONINFO      = 0x00000040
	MB_ICONERROR     = 0x00000010
	MB_SYSTEMMODAL   = 0x00001000
	CREATE_NO_WINDOW = 0x08000000
)

func winMsgBox(title, msg string, flags uintptr) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(msg)
	messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), flags)
}

func getExePath() (string, error) {
	var m = make([]uint16, syscall.MAX_PATH)
	n, _, err := kernel32.NewProc("GetModuleFileNameW").Call(
		0, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)))
	if n == 0 {
		return "", err
	}
	return syscall.UTF16ToString(m[:n]), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	exePath, err := getExePath()
	if err != nil {
		winMsgBox("错误", fmt.Sprintf("无法获取程序路径: %v", err), MB_OK|MB_ICONERROR|MB_SYSTEMMODAL)
		return
	}

	rootDir := filepath.Dir(exePath)
	rescueScript := filepath.Join(rootDir, ".system", "openclaw_rescue", "main.py")

	if !fileExists(rescueScript) {
		winMsgBox("错误", fmt.Sprintf("未找到救援脚本: %s", rescueScript), MB_OK|MB_ICONERROR|MB_SYSTEMMODAL)
		return
	}

	// Try python first, fall back to python.exe in .system\env
	pythonExe := filepath.Join(rootDir, ".system", "env", "python.exe")
	if !fileExists(pythonExe) {
		winMsgBox("错误", "未找到 Python 运行环境", MB_OK|MB_ICONERROR|MB_SYSTEMMODAL)
		return
	}

	cmd := exec.Command(pythonExe, rescueScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}
	err = cmd.Start()
	if err != nil {
		winMsgBox("错误", fmt.Sprintf("启动救援脚本失败: %v", err), MB_OK|MB_ICONERROR|MB_SYSTEMMODAL)
		return
	}

	// Wait for the rescue script to finish (it shows its own MessageBox)
	cmd.Process.Wait()
	time.Sleep(500 * time.Millisecond)
}