# OpenClaw U盘便携版

> 插上U盘就能用的OpenClaw AI助手，内置Token中转站支持

## 文件说明

| 文件 | 说明 |
|------|------|
| `使用说明.txt` | 使用指南 |
| `注册账号.txt` | 注册步骤 |
| `启动OpenClaw.bat` | 启动服务 |
| `关闭OpenClaw.bat` | 关闭服务 |
| `setup.js` | 激活引导页面 |
| `openclaw_portable_marker.txt` | 盘符检测标记 |

## 目录结构

```
OpenClaw_Portable/
├── node/                    # Node.js 便携版
├── openclaw/                # OpenClaw 本体
├── .openclaw/               # 配置文件
├── setup.js                 # 激活引导
├── 使用说明.txt             # 使用指南
├── 注册账号.txt             # 注册步骤
├── 启动OpenClaw.bat         # 启动脚本
├── 关闭OpenClaw.bat         # 关闭脚本
└── openclaw_portable_marker.txt  # 盘符标记
```

## 开发说明

### 构建步骤

1. 下载 Node.js 便携版到 `node/` 目录
2. 下载 OpenClaw 到 `openclaw/` 目录
3. 配置 `.openclaw/openclaw.json`（API地址等）

### 更新setup.js

激活引导页逻辑在 `setup.js`，包含：
- 注册引导
- API Key验证
- 启动OpenClaw

## 部署平台要求

- One API 已部署（http://3295b30e.r8.cpolar.cn/）
- cpolar 内网穿透已配置

## License

MIT