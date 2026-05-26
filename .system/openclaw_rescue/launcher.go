package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	MB_OK             = 0x00000000
	MB_ICONERROR      = 0x00000010
	MB_ICONWARNING    = 0x00000030
	MB_SYSTEMMODAL    = 0x00001000
	CREATE_NO_WINDOW  = 0x08000000
)

func winMsgBox(title, msg string, flags uintptr) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	msgPtr, _ := syscall.UTF16PtrFromString(msg)
	messageBoxW.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), flags)
}

func winErrorBox(err error) {
	if err == nil {
		return
	}
	winMsgBox("OpenClaw 启动器 - 错误", err.Error(), MB_OK|MB_ICONERROR|MB_SYSTEMMODAL)
	writeErrorLog(err.Error())
	os.Exit(1)
}

func getExePath() string {
	var m = make([]uint16, syscall.MAX_PATH)
	n, _, _ := kernel32.NewProc("GetModuleFileNameW").Call(
		0, uintptr(unsafe.Pointer(&m[0])), uintptr(len(m)))
	return syscall.UTF16ToString(m[:n])
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readConfig(path string) (map[string]interface{}, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	var cfg map[string]interface{}
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, false
	}
	providers, ok := cfg["models"].(map[string]interface{})
	if !ok {
		return nil, false
	}
	pList, ok := providers["providers"].(map[string]interface{})
	if !ok {
		return nil, false
	}
	for _, v := range pList {
		p, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if apiKey, ok := p["apiKey"].(string); ok && apiKey != "" && apiKey != "***" {
			if apiKey == "USER_GETS_OWN_KEY" || len(apiKey) < 10 {
				continue
			}
			return cfg, true
		}
	}
	return nil, false
}

func logsDir() string {
	exePath := getExePath()
	return filepath.Join(filepath.Dir(exePath), ".system", "logs")
}

func writeStartupLog(format string, args ...interface{}) {
	dir := logsDir()
	os.MkdirAll(dir, 0755)
	f, _ := os.OpenFile(filepath.Join(dir, "startup.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if f == nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}

func writeErrorLog(msg string) {
	dir := logsDir()
	os.MkdirAll(dir, 0755)
	f, _ := os.OpenFile(filepath.Join(dir, "error.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if f == nil {
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func startNode(script string, args ...string) error {
	nodeExe := filepath.Join(filepath.Dir(getExePath()), "node", "node.exe")
	cmdArgs := []string{script}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(nodeExe, cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}
	cmd.Dir = filepath.Dir(getExePath())
	return cmd.Start()
}

func waitForPort(port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		cmd := exec.Command("cmd", "/c", "netstat -an")
		out, err := cmd.Output()
		if err == nil && strings.Contains(string(out), port) && strings.Contains(string(out), "LISTENING") {
			return nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return fmt.Errorf("timeout")
}

func openBrowser(url string) {
	cmd := exec.Command("cmd", "/c", "start", "", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

func main() {
	exePath := getExePath()
	rootDir := filepath.Dir(exePath)

	nodeExe := filepath.Join(rootDir, "node", "node.exe")
	setupJs := filepath.Join(rootDir, "setup.js")
	openclawMjs := filepath.Join(rootDir, "openclaw", "openclaw.mjs")
	configPath := filepath.Join(rootDir, ".openclaw", "openclaw.json")

	nodeOk := fileExists(nodeExe)
	setupOk := fileExists(setupJs)
	openclawOk := fileExists(openclawMjs)
	_, hasKey := readConfig(configPath)

	writeStartupLog("Launcher start. node=%v setup=%v openclaw=%v hasKey=%v",
		nodeOk, setupOk, openclawOk, hasKey)

	if !nodeOk {
		winErrorBox(fmt.Errorf("未找到 Node.js 运行环境"))
	}

	// Case 1: no valid API key configured → show activation page
	if !hasKey {
		if !setupOk {
			winErrorBox(fmt.Errorf("未找到激活页面 setup.js"))
		}
		writeStartupLog("No valid API key, launching activation page")
		err := startNode(setupJs)
		if err != nil {
			writeErrorLog(fmt.Sprintf("startNode setup.js error: %v", err))
			winErrorBox(fmt.Errorf("启动激活页面失败: %v", err))
		}
		err = waitForPort(":8080", 10*time.Second)
		if err != nil {
			writeErrorLog(fmt.Sprintf("waitForPort :8080 error: %v", err))
			winErrorBox(fmt.Errorf("激活页面启动超时，请检查 .system\\logs\\error.log"))
		}
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://127.0.0.1:8080")
		return
	}

	// Case 2: valid API key → start OpenClaw directly
	if openclawOk {
		writeStartupLog("Starting OpenClaw normal mode")
		err := startNode(openclawMjs, "gateway", "start")
		if err != nil {
			writeErrorLog(fmt.Sprintf("startNode openclaw error: %v", err))
			winErrorBox(fmt.Errorf("启动 OpenClaw 失败: %v", err))
		}
		err = waitForPort(":8080", 15*time.Second)
		if err != nil {
			writeErrorLog(fmt.Sprintf("waitForPort :8080 error: %v", err))
			winErrorBox(fmt.Errorf("服务启动超时，请检查 .system\\logs\\error.log"))
		}
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://127.0.0.1:8080")
		writeStartupLog("OpenClaw started successfully")
		return
	}

	// Case 3: no openclaw.mjs but has setup.js → fall back to activation
	if setupOk {
		winMsgBox("OpenClaw 启动器", "启动异常，正在启动急救模式...",
			MB_OK|MB_ICONWARNING|MB_SYSTEMMODAL)
		err := startNode(setupJs)
		if err != nil {
			winErrorBox(fmt.Errorf("启动急救模式失败: %v", err))
		}
		waitForPort(":8080", 10*time.Second)
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://127.0.0.1:8080")
		return
	}

	winErrorBox(fmt.Errorf("未找到任何可用入口文件"))
}