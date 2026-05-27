# OpenClaw U盘便携版

把 OpenClaw 装进U盘，插到任意Windows电脑即插即用。

## 目录结构

```
D:\OpenClaw_Portable\                    ← 根目录（可见文件）
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
        ├── platform.json              ← ⭐ 平台地址配置（换域名只改这里）
        ├── version.json                ← 版本信息
        ├── openclaw_normal\             ← 日常版配置
        │   └── config_default.json
        ├── openclaw_rescue\            ← 救援版工具源码
        │   ├── config_default.json
        │   ├── launcher.go
        │   ├── main.py
        │   ├── quit.go
        │   └── 激活账号.go
        └── logs\                        ← 日志目录
            ├── startup.log
            └── error.log
```

## 核心配置文件

### `.system/platform.json`

平台地址的唯一真相来源（换域名只改这里）：

```json
{
  "platform": {
    "baseUrl": "http://3295b30e.r8.cpolar.cn",
    "apiPath": "/v1"
  }
}
```

所有用到平台地址的地方都会读取这个文件：
- `setup.js` 激活页面
- `openclaw.json` 中的 baseUrl

### `.system/version.json`

版本信息：

```json
{
  "version": "1.0.0",
  "build": "2026-05-27",
  "description": "OpenClaw U盘便携版"
}
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
  → 检测是否还有残留 → 弹窗提示是否可以拔盘
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

用户在平台注册账号 → 充值获得 Token 额度 → 在后台获取 API Key → 填入激活页面

平台地址在 `.system/platform.json` 中配置。

## 商业模式

| 成本 | 售价 | 利润 |
|------|------|------|
| ~1元/百万Token | 3-5元/百万Token | 2-4元/百万Token |

用户买U盘 → 注册平台账号 → 充值 → 使用 OpenClaw

## GitHub

https://github.com/yangkai258/openclaw-usb