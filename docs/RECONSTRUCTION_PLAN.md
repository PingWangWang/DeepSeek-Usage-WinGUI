# DeepSeek Usage+ → Go Windows 桌面程序重构方案

## 需求理解

**现状**：Tampermonkey 用户脚本（6269 行 JS），支持用量分析、订阅推送、Key 明细导入等功能

**目标**：
- 保留全部功能
- 转换为单文件 Windows EXE
- 从 DeepSeek API 获取数据
- 个人使用场景
- 包含完整工程化：目录结构、版本管理、打包、自动更新

---

## 项目现状分析

### 现有功能模块

1. **数据层**（API 调用 → 数据归一化）
   - 3 个 API 端点：amount、cost、summary
   - Token 认证（localStorage 提取）
   - 响应解包（biz_data）、字段映射

2. **业务逻辑层**（数据聚合 → 指标计算）
   - 日/周/月聚合
   - 均价计算、可用 Token 预估
   - Key/模型分组统计

3. **图表层**（ECharts 配置）
   - 7 张图表定义
   - 主题适配、响应式布局

4. **UI 层**（HTML 渲染 → 交互）
   - 仪表盘卡片、表格、图表、模式框
   - 用户交互（月份切换、开关、筛选）

5. **数据持久化**（localStorage）
   - 配置（图表开关、刷新间隔、订阅配置）
   - 缓存（Key 明细、日常数据）

6. **高级功能**
   - Key 明细导入（ZIP → CSV 解析）
   - 订阅推送（Webhook / 剪贴板 / 截图）
   - 自动刷新定时器
   - 订阅定时检查

---

## 方案一：Wails + Vue 3 + Electron 风格（推荐）

**一句话概括**：使用 Wails（Go + Vue3 前后端分离）构建，编译为单文件 EXE，保持 Web UI 灵活性，天然支持跨更新。

### 核心思路

利用 Wails 框架的优势：
- **Go 后端**处理数据获取、API 调用、定时任务、系统集成
- **Vue 3 前端**负责 UI 渲染、ECharts 图表、用户交互
- **单文件 EXE**通过 Wails 编译时嵌入所有资源
- **自动更新**通过 Go 后端检查新版本、下载替换

### 详细设计

#### 1. 目录结构

```
DeepSeek-Usage-WinGUI/
├── go.mod / go.sum               # Go 依赖管理
├── main.go                        # Wails 应用入口
├── wails.json                     # Wails 配置
├── version.go                     # 版本信息（自动生成）
│
├── backend/                       # Go 后端代码
│   ├── app.go                     # 应用主体、Wails 绑定
│   ├── config/
│   │   ├── config.go              # 配置加载/保存
│   │   └── paths.go               # 路径管理（AppData）
│   ├── api/
│   │   ├── deepseek.go            # DeepSeek API 客户端
│   │   └── types.go               # API 响应结构体
│   ├── service/
│   │   ├── analytics.go           # 数据聚合/计算
│   │   ├── export.go              # ZIP 导入/CSV 解析
│   │   ├── subscription.go        # 订阅逻辑、Webhook 推送
│   │   └── cache.go               # 缓存管理
│   ├── task/
│   │   ├── scheduler.go           # 定时任务调度
│   │   └── updater.go             # 版本检查/自动更新
│   └── util/
│       ├── logger.go              # 日志
│       └── format.go              # 数值格式化
│
├── frontend/                      # Vue 3 前端代码
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── components/
│   │   │   ├── Dashboard.vue      # 仪表盘容器
│   │   │   ├── SummaryCards.vue   # 汇总卡片
│   │   │   ├── Charts/
│   │   │   │   ├── RequestChart.vue
│   │   │   │   ├── TokenChart.vue
│   │   │   │   ├── CacheRateChart.vue
│   │   │   │   ├── CompositionChart.vue
│   │   │   │   ├── ModelChart.vue
│   │   │   │   └── KeyCostChart.vue
│   │   │   ├── Tables/
│   │   │   │   ├── ModelTable.vue
│   │   │   │   └── KeyTable.vue
│   │   │   ├── Subscription/
│   │   │   │   ├── SubscriptionManager.vue
│   │   │   │   ├── SubscriptionForm.vue
│   │   │   │   └── SubscriptionPreview.vue
│   │   │   ├── Settings.vue       # 设置面板
│   │   │   └── Toolbar.vue        # 工具栏
│   │   ├── stores/
│   │   │   ├── useAppStore.ts     # 全局状态（Pinia）
│   │   │   └── useDataStore.ts    # 数据状态
│   │   ├── api/
│   │   │   └── backend.ts         # 调用后端的 Wails binding
│   │   ├── hooks/
│   │   │   ├── useCharts.ts       # ECharts 初始化/自适应
│   │   │   └── useAutoRefresh.ts  # 自动刷新逻辑
│   │   └── styles/
│   │       ├── main.css
│   │       └── theme.css          # 明/暗主题
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts             # Vite 构建配置
│
├── assets/                        # 图标、资源
│   ├── icon.ico
│   └── icon.png
│
├── .github/workflows/
│   └── release.yml                # GitHub Actions 自动构建 Release
│
├── scripts/
│   ├── build.sh / build.bat       # 构建脚本
│   ├── pack.sh / pack.bat         # 打包脚本
│   └── version.sh / version.bat   # 版本号更新脚本
│
└── docs/
    ├── DEVELOPMENT.md             # 开发指南
    ├── BUILD.md                   # 构建说明
    └── CHANGELOG.md               # 更新日志
```

#### 2. 后端架构（Go）

**API 客户端** (`backend/api/deepseek.go`)
```go
type DeepSeekClient struct {
    baseURL string
    token   string
    client  *http.Client
}

// 三个核心方法
func (c *DeepSeekClient) GetUsageAmount(period string) (*AmountResponse, error)
func (c *DeepSeekClient) GetUsageCost(period string) (*CostResponse, error)
func (c *DeepSeekClient) GetUserSummary() (*SummaryResponse, error)
```

**业务逻辑层** (`backend/service/analytics.go`)
```go
type AnalyticsService struct {
    api    *DeepSeekClient
    cache  *CacheManager
}

// 数据聚合
func (s *AnalyticsService) GetDashboard(period string) (*DashboardData, error)
func (s *AnalyticsService) CalculateMetrics(data *RawData) *Metrics
func (s *AnalyticsService) GroupByModel(data *RawData) map[string]*ModelStats
```

**订阅推送** (`backend/service/subscription.go`)
```go
type SubscriptionService struct {
    config Config
}

func (s *SubscriptionService) SendReport(sub *Subscription, content string) error
func (s *SubscriptionService) SendToWebhook(url string, message map[string]interface{}) error
```

**定时任务** (`backend/task/scheduler.go`)
```go
type Scheduler struct {
    jobs map[string]*ScheduledJob
}

func (s *Scheduler) RegisterAutoRefresh(interval time.Duration, callback func()) error
func (s *Scheduler) RegisterSubscriptionCheck(subscription *Subscription, callback func()) error
func (s *Scheduler) Stop()
```

**自动更新** (`backend/task/updater.go`)
```go
type Updater struct {
    currentVersion string
    repoURL        string
}

func (u *Updater) CheckUpdate() (*UpdateInfo, error)
func (u *Updater) DownloadAndUpdate(info *UpdateInfo) error
```

**Wails 绑定** (`backend/app.go`)
```go
type App struct {
    analytics      *AnalyticsService
    subscription   *SubscriptionService
    scheduler      *Scheduler
    updater        *Updater
    config         *Config
}

// 暴露给前端的方法（自动生成 TypeScript 类型）
func (a *App) GetDashboard(period string) (*DashboardData, error)
func (a *App) SendSubscription(subID string) error
func (a *App) SetAutoRefresh(interval int) error
func (a *App) CheckUpdate() (*UpdateInfo, error)
func (a *App) ImportKeyDetail(zipPath string) error
```

#### 3. 前端架构（Vue 3 + TypeScript）

**状态管理** (`frontend/src/stores/useAppStore.ts`)
- 配置状态（图表开关、刷新间隔）
- 订阅管理状态
- 主题状态

**数据获取** (`frontend/src/api/backend.ts`)
```typescript
// 自动生成的 Wails binding
export const backend = {
  getDashboard: (period: string) => Promise<DashboardData>,
  sendSubscription: (subID: string) => Promise<void>,
  // ...
};
```

**组件设计**
- Dashboard 容器：管理整体数据和状态
- 功能组件：卡片、表格、图表独立封装
- Toolbar：月份选择、开关按钮、自动刷新
- 设置面板：配置编辑、导入导出

**图表初始化** (`frontend/src/hooks/useCharts.ts`)
- 使用组合式 API 管理图表生命周期
- 响应式图表尺寸（ResizeObserver）
- 主题自适应（对应明/暗模式）

#### 4. 数据持久化

**配置文件**（`%APPDATA%/DeepSeek-Usage/config.json`）
```json
{
  "token": "sk-...",
  "sectionVisible": {"requests": true, "tokens": true, ...},
  "autoRefreshInterval": 300000,
  "subscriptions": [...],
  "keyFilter": {"mode": "all", "keys": [...]},
  ...
}
```

**缓存**（`%APPDATA%/DeepSeek-Usage/cache/`）
- 日常数据缓存（避免频繁 API 调用）
- Key 明细 ZIP 缓存

#### 5. 版本管理与更新

**版本号**：遵循 semver（v1.0.0）
- 每次构建时从 git tag 自动提取
- 内嵌到 EXE 的 version.go

**更新检查**
- 启动时后台检查 GitHub Release
- 用户可手动触发检查
- 下载新版本到临时目录，退出时替换（或使用 updater.exe 包装）

**发布流程**
- git tag → GitHub Actions 自动编译 Windows EXE
- Release 页自动生成 changelog
- 用户应用内点击"更新"下载安装

#### 6. 打包为单文件 EXE

**Wails 原生支持**
```bash
wails build -webview2 Embed -o DeepSeek-Usage.exe
```
- 所有前端资源嵌入到 EXE（无需 HTML/JS 外置文件）
- WebView2 runtime 可选打包或依赖系统已安装

**大小优化**
- 前端资源压缩（Vite 打包）
- Go 二进制精简（`-ldflags "-s -w"`）
- 预计 EXE 大小：50-100MB（包含 WebView2）

### 优点

✅ **完整的工程化支持**：版本管理、自动更新、CI/CD
✅ **UI 灵活性**：保留 Web 开发的便利（CSS、Vue 组件化）
✅ **性能**：Go 后端高效，前端 Vue 3 响应式
✅ **跨平台基础**：Wails 框架支持 macOS/Linux 扩展（未来升级）
✅ **单文件 EXE**：零依赖启动
✅ **原生集成**：可调用 Windows API（系统托盘、快捷键、系统通知等）

### 缺点

⚠️ **学习成本**：需要熟悉 Go、Vue 3、Wails 框架
⚠️ **EXE 体积**：50-100MB（较大）
⚠️ **首次启动**：初始化前端框架有延迟（通常 < 1s）
⚠️ **第三方依赖**：Go 依赖库需要管理版本

---

## 方案二：纯 Go + 原生 Windows UI（Fyne）

**一句话概括**：使用 Fyne GUI 框架，完全原生 Go 开发，体积小，性能高，但 UI 定制化程度较低。

### 核心思路

- Go 代码 100%，无前端框架
- Fyne 提供跨平台 GUI（Windows、macOS、Linux）
- 编译后 EXE 体积小（~20-30MB）
- 开发速度快，依赖简单

### 详细设计

**目录结构**简化：
```
DeepSeek-Usage-WinGUI/
├── main.go                       # 应用入口 + Fyne UI 初始化
├── backend/                      # 与方案一相同
│   ├── app.go
│   ├── api/
│   ├── service/
│   └── ...
└── ui/
    ├── dashboard.go              # 仪表盘 UI
    ├── charts.go                 # ECharts 集成？或 Fyne Charts
    ├── components.go             # 可重用 UI 组件
    └── theme.go                  # 主题
```

**UI 实现**：
- Fyne 内置组件（Container、VBox、HBox、Table）
- 第三方图表库（fyne-x、或嵌入 ECharts via HTML Canvas）
- 样式配置（颜色、字体、间距）

### 优点

✅ **体积小**：EXE 仅 20-30MB
✅ **性能高**：纯 Go，无虚拟机开销
✅ **开发快**：Fyne API 简洁
✅ **真正跨平台**：Windows/macOS/Linux 代码共用

### 缺点

⚠️ **UI 灵活性不足**：难以实现复杂 ECharts
⚠️ **图表能力弱**：Fyne 内置图表有限，第三方库社区生态较差
⚠️ **自定义成本高**：样式定制受限
⚠️ **前端开发体验差**：不如 Web 开发便利

---

## 方案三：Go + WebView2（轻量级 Wails 替代）

**一句话概括**：自己实现 Go + WebView2 集成（无依赖 Wails 框架），使用原生 HTML/CSS/JS，体积最小。

### 核心思路

- 用 `github.com/jchv/go-webview2` 直接嵌入 WebView2
- 前端用原生 JavaScript（不依赖 Vue、React，减少打包体积）
- 自己实现前后端通信的 bridge
- 完全掌控编译过程

### 优点

✅ **体积最小**：EXE ~40-50MB（无 Wails 框架开销）
✅ **完全自主**：不依赖第三方框架更新
✅ **性能最优**：直接调用 WebView2 API

### 缺点

⚠️ **开发复杂**：需要自己实现前后端通信、资源嵌入等细节
⚠️ **维护成本高**：框架选型风险高，社区生态弱
⚠️ **文档不足**：相比 Wails 学习成本更高

---

## 方案对比

| 维度 | 方案一（Wails + Vue） | 方案二（Fyne） | 方案三（WebView2） |
|------|------------------------|---------------|-------------------|
| **工作量** | 15-20 人天 | 12-15 人天 | 18-25 人天 |
| **EXE 体积** | 50-100MB | 20-30MB | 40-50MB |
| **首次启动** | ~1-2s | <500ms | <1s |
| **UI 灵活性** | 高（Vue 生态完整） | 低（Fyne 组件有限） | 中（原生 JS） |
| **图表能力** | 优秀（ECharts） | 一般（Fyne Charts） | 优秀（ECharts） |
| **开发体验** | 优秀（Web 开发） | 中（纯 Go） | 一般（手动 bridge） |
| **跨平台扩展** | 容易（Wails 支持） | 容易（Fyne 原生支持） | 困难（WebView2 仅 Windows） |
| **更新机制** | 成熟（Wails 内置）| 无（需自己实现） | 无（需自己实现） |
| **学习成本** | 中（Go + Vue + Wails） | 低（纯 Go） | 高（手动集成） |
| **长期维护** | 低风险（生态活跃） | 中风险（社区小） | 高风险（框架风险） |

---

## 建议

**推荐方案一（Wails + Vue 3）**

理由：
1. **工程化最完整**：内置版本管理、自动更新、资源嵌入、CI/CD 支持
2. **功能保留度高**：ECharts 完全兼容，订阅推送、定时任务都能无缝迁移
3. **开发效率最高**：Web 开发模式熟悉，热刷新调试便利
4. **长期易维护**：Vue 和 Wails 都是活跃社区，更新有保障
5. **未来扩展容易**：跨平台、系统集成（托盘、通知等）都有成熟方案
6. **符合现实需求**：单文件 EXE、订阅推送、自动更新，Wails 都原生支持

**方案二的适用场景**：
- 如果 UI 需求简单（只需列表、基本图表）
- 如果对 EXE 体积有严格限制（不能超过 30MB）
- 如果想要最快的启动速度和最好的性能

**方案三不推荐**：
- 开发成本高、维护风险大，得不偿失
- 不如 Wails 成熟稳定

---

## 后续步骤

1. ✅ **确认方案**：评审上述设计，确认使用方案一（Wails）
2. 🔨 **项目初始化**：`wails create` 创建项目骨架，配置 Go 依赖（echarts、jszip、webhook 客户端）
3. 📂 **代码迁移**：
   - 将 JavaScript 逻辑移植到 Go（API 调用、数据聚合、CSV 解析、Webhook 推送）
   - 前端逐组件重写：Vue 3 + TypeScript
4. 🧪 **功能验证**：逐个验证 API 调用、图表渲染、订阅推送、自动更新
5. 📦 **打包发布**：CI/CD 流程配置，自动生成 Release EXE
6. 🚀 **版本管理**：确定 v1.0.0 版本号、开源协议、更新机制

---

## 补充：与原脚本的关键差异

| 功能点 | 原脚本（JS） | 新程序（Go）| 迁移说明 |
|--------|-----------|---------|---------|
| Token 获取 | localStorage 自动提取 | 用户手动输入/保存 | 首启动时弹出配置框 |
| API 调用 | fetch + Tampermonkey GM_xmlhttpRequest | Go http 客户端 | 无感迁移 |
| 定时刷新 | setInterval | time.Ticker | 后台定时任务 |
| 截图上传 | html2canvas + ImgBB | Go 调用系统截图工具或库 | 可选功能 |
| CSV 解析 | 原生 JS 正则 | Go encoding/csv | 性能更好 |
| ZIP 导入 | JSZip (JS) | archive/zip (Go) | 更稳定高效 |
| localStorage | 浏览器存储 | JSON 文件（%APPDATA%） | 程序启动时加载 |
| 主题自适应 | 监听 CSS 变量 | Qt/WinAPI 获取系统主题 | 更精准 |

