# OpenClaw U盘便携版

把 OpenClaw 装进U盘，插到任意Windows电脑即插即用。

## 目录结构

```
D:\OpenClaw_Portable\                    ← 根目录（5个可见文件）
│
├── OpenClaw.exe                         ← ⭐ 主程序（Go双模启动器）
├── 激活账号.exe                          ← 激活账号工具
├── 急救与重置工具.exe                    ← 恢复出厂设置
├── 彻底退出软件.exe                      ← 安全退出+提示拔盘
├── 快速入门指南.pdf                      ← 使用说明
├── setup.js                             ← 激活页面服务
│
└── (全部隐藏)
    ├── node\                            ← Node.js 便携版
    │   └── node.exe
    ├── openclaw\                        ← OpenClaw 主程序
    │   └── openclaw.mjs
    ├── .openclaw\                      ← 用户配置
    │   └── openclaw.json
    └── .system\
        ├── openclaw_normal\             ← 日常版配置
        │   └── config_default.json
        ├── openclaw_rescue\             ← 救援版工具
        │   ├── config_default.json
        │   ├── launcher.go             ← 救援模式启动器源码
        │   ├── main.py                  ← 急救脚本
        │   ├── quit.go                  ← 退出工具源码
        │   └── 激活账号.go               ← 激活工具源码
        └── logs\                        ← 日志目录
            ├── startup.log
            └── error.log
```

## 设计理念

### 双模式机制

| 模式 | 条件 | 行为 |
|------|------|------|
| **日常版** | 有有效 API Key | 直接启动 OpenClaw |
| **救援版** | 无有效 API Key 或启动失败 | 弹出激活页面 |

### 激活流程

```
插入U盘 → 运行 OpenClaw.exe
  ├── 检测到有效 API Key → 直接启动 ✅
  └── 未检测到 → 启动 setup.js (端口8080)
        → 打开浏览器 http://127.0.0.1:8080
        → 用户输入 API Key → 验证 → 保存配置 → 启动
```

### 安全退出流程

```
运行"彻底退出软件.exe"
  → 弹出确认对话框
  → 确定 → kill 所有 node.exe 进程
  → 弹出"可以安全拔U盘了" → 用户拔盘
```

### 出厂重置流程

```
运行"急救与重置工具.exe"
  → 弹出确认对话框
  → 确定 → 从 config_default.json 恢复配置
  → 重置后需重新激活
```

## 开发说明

### 工具链

| 文件 | 源码 | 语言 | 说明 |
|------|------|------|------|
| OpenClaw.exe | launcher.go | Go | 双模启动器（需编译） |
| 激活账号.exe | 激活账号.go | Go | 单独打开激活页面（需编译） |
| 彻底退出软件.exe | quit.go | Go | 安全退出工具（需编译） |
| 急救与重置工具.exe | main.py | Python | 恢复出厂设置（需Python） |

### 编译方法

```bash
# Go 工具（需要 Go 1.21+）
cd .system/openclaw_rescue
go build -ldflags="-s -w" -o ..\..\OpenClaw.exe launcher.go
go build -ldflags="-s -w" -o ..\..\激活账号.exe 激活账号.go
go build -ldflags="-s -w" -o ..\..\彻底退出软件.exe quit.go

# Python 工具（需要 Python 3.x）
cd .system/openclaw_rescue
pyinstaller --onefile --console --name "急救与重置工具.exe" main.py
```

### API Key 来源

平台地址（中转站）：`http://3295b30e.r8.cpolar.cn`

用户需要：
1. 打开平台注册账号
2. 充值获得 Token 额度
3. 在后台获取 API Key
4. 填入激活页面

## 文件说明

### setup.js

激活页面服务器（Node.js 内置 http.server）：
- `GET /` → 激活页面 HTML
- `POST /save-config` → 保存 API Key 到配置
- `GET /start` → 触发 OpenClaw 启动

### config_default.json

出厂默认配置，launcher.go 和 main.py 都依赖它做重置。

## 商业模式

| 成本 | 售价 | 利润 |
|------|------|------|
| ~1元/百万Token | 3-5元/百万Token | 2-4元/百万Token |

用户买U盘 → 注册平台账号 → 充值 → 使用 OpenClaw

## GitHub

https://github.com/yangkai258/openclaw-usb