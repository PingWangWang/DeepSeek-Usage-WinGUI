# 开发指南

## 本地开发环境搭建

### 1. 前置条件

- Go 1.21+：[下载](https://golang.org/dl/)
- Node.js 18+：[下载](https://nodejs.org/)
- Git

### 2. 安装依赖

```bash
# 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 克隆项目
git clone https://github.com/PingWangWang/DeepSeek-Usage-WinGUI.git
cd DeepSeek-Usage-WinGUI

# 安装 Go 依赖
go mod download

# 安装前端依赖
cd frontend
npm install
cd ..
```

### 3. 开发模式

```bash
# 启动开发服务器（带热加载）
wails dev

# 应用会在 http://localhost:5173 打开
# 修改前端代码自动刷新
# 修改后端代码需要手动重启（或 Ctrl+R）
```

### 4. 生产构建

```bash
# 构建为 EXE（输出到 build/bin/）
wails build -webview2 Embed

# 或使用脚本
scripts\build.bat       # Windows
./scripts/build.sh      # macOS/Linux
```

## 项目结构详解

### 后端 (backend/)

#### `app.go`
应用主体，暴露给 Wails 前端的 API 方法：
- `GetDashboard(period)` — 获取仪表盘数据
- `SendSubscription(subID)` — 发送订阅报告
- `CheckUpdate()` — 检查版本更新
- `ImportKeyDetail(zipPath)` — 导入 Key 明细

#### `api/deepseek.go`
DeepSeek API 客户端：
- `GetUsageAmount(period)` — 获取用量数据
- `GetUsageCost(period)` — 获取费用数据
- `GetUserSummary()` — 获取用户信息

#### `service/analytics.go`
数据聚合服务：
- 合并三个 API 的返回值
- 计算日常数据、模型统计
- 计算汇总指标（当日费用、均价、缓存命中率等）

#### `service/subscription.go`
订阅推送服务：
- 构建报告（Markdown/截图）
- 发送到钉钉/飞书/企业微信
- 定时检查和触发

#### `task/scheduler.go`
定时任务调度器：
- 管理自动刷新、订阅检查等后台任务
- 基于时间间隔触发回调

#### `task/updater.go`
版本更新检查：
- 查询 GitHub Release API
- 下载新版本 EXE
- 自动替换升级

#### `config/config.go`
配置管理：
- 从 `%APPDATA%/DeepSeek-Usage/config.json` 加载
- 管理 Token、UI 开关、订阅配置等

### 前端 (frontend/)

#### `components/`
Vue 3 组件：
- `Header.vue` — 顶部导航
- `Dashboard.vue` — 仪表盘容器
- `Toolbar.vue` — 工具栏（月份选择、开关等）
- `SummaryCards.vue` — 汇总卡片
- `Card.vue` — 单个卡片组件
- `Charts.vue` — 图表区域（ECharts）

#### `stores/`
Pinia 状态管理：
- `useAppStore.ts` — 应用全局状态（Token、主题、UI 开关）
- `useDataStore.ts` — 仪表盘数据状态

#### `styles/`
全局样式：
- `main.css` — 重置、主题、滚动条等

### 配置文件

#### `go.mod`
Go 依赖管理（Wails、resty、echarts 库等）

#### `wails.json`
Wails 框架配置：
- 前后端目录
- 编译目标（Windows)
- 窗口大小、标题等

#### `frontend/package.json`
npm 依赖：
- vue@3
- pinia（状态管理）
- echarts（图表）
- axios（HTTP 客户端）

#### `frontend/tsconfig.json`
TypeScript 配置

#### `frontend/vite.config.ts`
Vite 构建配置

## 常见开发任务

### 添加新的 API 端点

1. 在后端 `app.go` 中添加方法：
```go
func (a *App) MyNewMethod(param string) (result interface{}, error error) {
    // 实现逻辑
    return result, nil
}
```

2. 前端自动生成 TypeScript binding，可直接调用：
```typescript
const result = await window.go.main.App.MyNewMethod('param')
```

### 修改配置存储格式

1. 编辑 `backend/config/config.go` 中的 `Config` 结构体
2. 现有配置会自动 migrate（JSON Unmarshal 处理不存在的字段）

### 添加新的图表

1. 在 `frontend/src/components/Charts.vue` 中引入 ECharts
2. 定义图表选项
3. 监听数据变化，调用 `echarts.setOption()`

### 实现订阅推送（Webhook）

1. 编辑 `backend/service/subscription.go` 中的 `sendToWebhook()` 方法
2. 构建 Markdown 或 JSON 报告体
3. 调用 `http.Post()` 发送到 Webhook URL

## 测试

### 手动测试

1. 填写有效的 DeepSeek API Token
2. 点击"刷新"观察数据是否正确加载和显示

### 日志

- 日志输出到 `%APPDATA%/DeepSeek-Usage/logs/app-YYYY-MM-DD.log`
- 控制台也会输出日志

## 故障排查

### 应用启动失败

- 查看日志文件（见上）
- 确认 WebView2 运行时已安装
- 检查网络连接

### API 调用失败

- 确认 Token 有效
- 检查网络代理设置
- 查看日志中的 API 错误信息

### 前端组件无法加载

- 运行 `npm install` 确保依赖完整
- 清除浏览器缓存
- 重启开发服务器

## 发布新版本

1. 修改 `main.go` 中的 AppVersion（或标签 git tag）
2. 运行 `scripts/build.bat` 编译 EXE
3. 创建 GitHub Release，上传 EXE
4. 用户应用会自动检测到更新

## 性能优化

- 【后端】使用缓存减少 API 调用
- 【前端】使用 Vue 3 的 `<Suspense>` 处理异步加载
- 【图表】使用 ECharts 的 incremental 模式处理大数据集
