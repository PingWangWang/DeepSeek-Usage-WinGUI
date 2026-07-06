# DeepSeek Usage+ — Windows 版用量分析仪表盘

> 为 DeepSeek API 用户打造的实时数据分析工具。一个 EXE 文件，开箱即用，支持图表分析、自动刷新、Webhook 推送和 Key 明细导入。

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](RELEASES)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078d4.svg)](README.md)

## ✨ 核心能力

- **📊 实时数据看板** — 费用、Token、缓存、均价一目了然
- **📈 5 张交互式图表** — 请求趋势、Token 构成、缓存命中率、模型分布等
- **🔄 自动刷新** — 30 秒～1 小时可选，后台定时轮询
- **🔔 Webhook 推送** — 支持钉钉、飞书、企业微信，定时发送用量报告
- **🔑 Key 明细导入** — 从 DeepSeek 导出 ZIP，自动解析 CSV，按 Key/模型统计
- **🌙 深色主题** — 跟随系统自动切换，减少夜间眼睛疲劳
- **📦 单文件分发** — Windows EXE，无需安装依赖，双击即运行
- **🔐 完全本地** — Token 仅保存在本地 AppData，不上传任何数据

## 🎯 适合谁用

- 👨‍💻 **DeepSeek API 用户** — 需要实时监控 API 用量和费用
- 💼 **多 Key 管理者** — 按 Key/模型细分统计，优化成本分配
- 🤖 **LLM 应用开发者** — 观察缓存命中率，调整 Prompt 策略，降低使用成本
- 📊 **数据驱动决策者** — 按天/按周/按月切换统计周期，追踪成本趋势

## 📥 安装 & 使用

### 快速开始（推荐）

1. **下载** → 从 [Release](https://github.com/PingWangWang/DeepSeek-Usage-WinGUI/releases) 页面下载最新 `DeepSeek-Usage.exe`
2. **运行** → 双击 EXE 文件启动应用
3. **配置** → 进入"设置"页面，输入 DeepSeek API Token（从 [platform.deepseek.com](https://platform.deepseek.com) 复制）
4. **查看** → 返回"仪表盘"，点击"刷新"获取最新数据

### 从源码编译

**前置条件**：
- Go 1.21+
- Node.js 18+
- Clash 代理（或任何 HTTP 代理）

**编译步骤**：
```bash
git clone https://github.com/PingWangWang/DeepSeek-Usage-WinGUI.git
cd DeepSeek-Usage-WinGUI

# 使用 PowerShell 执行构建脚本
.\scripts\build.ps1 -ProxyPort 7897  # 按需调整代理端口

# 或手动步骤
cd frontend && npm install && npm run build && cd ..
go mod tidy
go build -o DeepSeek-Usage.exe .
```

## 🧭 使用示例

### 查看仪表盘

启动应用后，默认显示当月数据：

- **汇总卡片** — 当日费用、月度费用、均价、总 Token、缓存命中率
- **图表** — 请求趋势、Token 用量、缓存率、模型分布等
- **工具栏** — 月份切换、图表显示/隐藏、自动刷新设置

### 导入 Key 明细

1. 进入 [DeepSeek 平台](https://platform.deepseek.com/usage)，点击"导出"下载 ZIP
2. 在应用中选择"导入 Key 明细"，选择 ZIP 文件
3. 应用自动解压、解析、聚合数据
4. 查看 Key 级别的费用、Token、缓存数据

### 配置 Webhook 推送

在"设置"中添加订阅：

1. **接收方式** — 钉钉/飞书/企业微信（Webhook URL）
2. **内容选择** — 费用摘要、Token 构成、缓存命中率等
3. **发送频率** — 自定义间隔、每日、每周、每月
4. **Key 筛选** — 可针对特定 API Key 生成报告

## 📁 项目结构

```
DeepSeek-Usage-WinGUI/
├── main.go                          # Wails 应用入口
├── backend/                         # Go 后端
│   ├── app.go                      # 应用主体
│   ├── api/                        # DeepSeek API 客户端
│   ├── service/                    # 业务逻辑（聚合、推送、导入）
│   ├── task/                       # 定时任务、版本检查
│   └── util/                       # 工具函数
├── frontend/                        # Vue 3 前端
│   ├── src/
│   │   ├── components/             # UI 组件（仪表盘、设置、图表）
│   │   ├── stores/                 # Pinia 状态管理
│   │   └── utils/                  # ECharts、API 工具
│   └── package.json
├── scripts/
│   ├── build.ps1                   # PowerShell 构建脚本
│   └── build.sh                    # Shell 构建脚本
├── docs/                           # 文档
│   ├── DEVELOPMENT.md              # 开发指南
│   └── BUILD.md                    # 构建说明
└── README.md / README_ZH.md        # 本文件

```

## 🛠️ 开发

### 启动开发服务

```bash
# 需要 Wails 安装
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 启动开发服务（自动热加载）
wails dev
```

### 构建生产版本

```bash
# 编译为单文件 Windows EXE
wails build -webview2 Embed -o DeepSeek-Usage.exe

# 或使用脚本
.\scripts\build.ps1
```

### 代码结构

**后端 (Go)**：
- `API 客户端` — 调用 DeepSeek 三个端点（用量、费用、用户信息）
- `数据聚合` — 日常、模型、Key 维度的分组统计
- `订阅推送` — Webhook、剪贴板、面板预览三种模式
- `导入导出` — ZIP 解压、CSV 解析、数据聚合

**前端 (Vue 3)**：
- `仪表盘` — 5 张 ECharts 图表 + 汇总卡片
- `设置页面` — Token 配置、自动刷新、版本检查
- `Pinia 状态** — 应用状态、数据缓存

## ❓ 常见问题

**Q: Token 会被保存在哪里？**  
A: 存储在 `%APPDATA%/DeepSeek-Usage/config.json`，仅本地保存，不会上传任何数据。

**Q: 如何更新到最新版本？**  
A: 在"设置"中点击"检查更新"，应用会从 GitHub Release 检查新版本，手动下载替换 EXE 即可。

**Q: Webhook 推送失败怎么办？**  
A: 检查 Webhook URL 是否正确，网络是否能访问外网，日志文件位于 `%APPDATA%/DeepSeek-Usage/logs/` 中。

**Q: 支持 macOS/Linux 吗？**  
A: 目前仅提供 Windows EXE。代码使用 Wails 框架，理论上可以跨平台编译，但尚未测试。

## 🚀 后续计划

- [ ] macOS/Linux 编译支持
- [ ] 自动更新功能（下载 + 替换）
- [ ] 截图上传到图床（ImgBB）
- [ ] 数据库存储（SQLite）支持离线模式
- [ ] 更多 Webhook 服务（Slack、Discord 等）

## 🤝 参与贡献

欢迎提交 Issue 和 PR！

### 报告 Bug
- 描述重现步骤
- 附加日志文件（`%APPDATA%/DeepSeek-Usage/logs/`）
- 说明操作系统版本和 Go 版本

### 提交 PR
- Fork 项目，创建特性分支
- 代码遵循现有风格
- 更新相关文档
- 提交 PR 并描述改动

## 📊 项目数据

| 指标 | 数值 |
|------|------|
| 代码行数 | ~6,000 行 |
| 依赖数 | 极少（仅用标准库 + Wails） |
| EXE 大小 | ~70-100 MB |
| 启动时间 | <2 秒 |
| 首次查询 | 1-2 秒 |

## 📝 许可证

MIT License — 自由使用、修改、分发，需保留版权声明。

---

## 🔗 相关链接

- **DeepSeek 官网** — https://www.deepseek.com
- **API 文档** — https://platform.deepseek.com/api-docs
- **用量查看** — https://platform.deepseek.com/usage
- **GitHub 仓库** — https://github.com/PingWangWang/DeepSeek-Usage-WinGUI

## 致谢

感谢 DeepSeek 提供开放的 API，感谢社区用户的反馈和建议。

---

**问题 & 反馈** → [GitHub Issues](https://github.com/PingWangWang/DeepSeek-Usage-WinGUI/issues)

**关注最新版本** → ⭐ Star 本项目
