package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

const CREATE_NO_WINDOW = 0x08000000

func getExePath() string {
	var m = make([]uint16, syscall.MAX_PATH)
	n, _, _ := kernel32.NewProc("GetModuleFileNameW").Call(
		0, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)))
	return syscall.UTF16ToString(m[:n])
}

func winErrorBox(err error) {
	if err == nil {
		return
	}
	user32.NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(err.Error()))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("错误"))),
		0x00000010|0x00001000|0x00000000)
	os.Exit(1)
}

func isPortInUse(port string) bool {
	conn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

func main() {
	exePath := getExePath()
	root := filepath.Dir(exePath)
	setupJs := filepath.Join(root, "setup.js")
	nodeExe := filepath.Join(root, "node", "node.exe")

	if !fileExists(setupJs) {
		winErrorBox(fmt.Errorf("未找到 setup.js"))
	}
	if !fileExists(nodeExe) {
		winErrorBox(fmt.Errorf("未找到 node.exe"))
	}

	// If port 8080 is already in use, server is already running — just open browser
	if isPortInUse("8080") {
		exec.Command("cmd", "/c", "start", "", "http://127.0.0.1:8080").Run()
		return
	}

	cmd := exec.Command(nodeExe, setupJs)
	cmd.Dir = root
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}
	err := cmd.Start()
	if err != nil {
		winErrorBox(fmt.Errorf("启动激活页面失败: %v", err))
	}
	// setup.js 启动后会根据配置自动决定是否打开浏览器
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}