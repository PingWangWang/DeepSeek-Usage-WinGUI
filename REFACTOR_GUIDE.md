# WPF 重构指南

## 概述

将现有 Greasy​monkey 脚本（JavaScript 浏览器增强）重构为独立 WPF（Windows Presentation Foundation）桌面应用。

**原始实现**：GoJS 脚本 (`assert/DeepSeek-Usage.user.js`)
- 运行环境：浏览器（Tampermonkey/Greasy​monkey）
- 功能范围：DeepSeek API 用量页增强
- 核心特性：图表、订阅推送、Key 明细分析

**目标实现**：WPF .NET 8 应用
- 运行环境：Windows 桌面（独立应用）
- 功能范围：完整的 API 用量分析工具
- 技术栈：MVVM + Community Toolkit

---

## 目录结构

```
src/
├── App.xaml                          # 应用全局配置
├── App.xaml.cs                       # 应用启动逻辑
├── DeepSeekUsageUI.csproj            # 项目文件
├── DeepSeekUsageUI.csproj.filters    # VS 筛选器
│
├── Models/                           # 数据模型
│   ├── ApiModels.cs                  # API 响应数据结构
│   ├── UIModels.cs                   # UI 绑定数据模型
│   └── Configuration.cs              # 配置数据结构
│
├── Views/                            # XAML 视图
│   ├── MainWindow.xaml               # 主窗口
│   ├── MainWindow.xaml.cs
│   ├── Dialogs/
│   │   ├── SettingsDialog.xaml       # 设置对话框
│   │   ├── SubscriptionDialog.xaml   # 订阅设置
│   │   └── KeyFilterDialog.xaml      # Key 筛选
│   └── Controls/
│       ├── ChartPanel.xaml           # 图表面板
│       ├── StatsSummary.xaml         # 统计摘要
│       └── KeyDetailTable.xaml       # Key 明细表
│
├── ViewModels/                       # MVVM 视图模型
│   ├── MainWindowViewModel.cs        # 主窗口 ViewModel
│   ├── StatsSummaryViewModel.cs      # 统计摘要 ViewModel
│   ├── KeyDetailViewModel.cs         # Key 明细 ViewModel
│   └── SettingsViewModel.cs          # 设置 ViewModel
│
├── Services/                         # 业务逻辑
│   ├── API/
│   │   ├── DeepSeekApiClient.cs      # API 客户端
│   │   ├── AuthTokenProvider.cs      # Token 管理
│   │   └── ApiModels.cs              # API 数据模型
│   │
│   ├── Storage/
│   │   ├── LocalStorageService.cs    # 本地存储（配置/缓存）
│   │   ├── CacheManager.cs           # 缓存管理
│   │   └── ConfigurationManager.cs   # 配置管理
│   │
│   └── Export/
│       ├── ExportService.cs          # 导出服务（CSV/JSON/ZIP）
│       ├── ScreenshotExporter.cs     # 截图导出
│       └── SubscriptionService.cs    # 订阅推送（钉钉/邮件等）
│
├── Resources/                        # 资源
│   ├── Brushes.xaml                  # 颜色/画刷
│   ├── Styles.xaml                   # 全局样式
│   └── Icons.xaml                    # 图标资源
│
├── Utilities/                        # 工具类
│   ├── FormatHelper.cs               # 格式化工具（数字/货币/百分比）
│   ├── LoggerHelper.cs               # 日志记录
│   └── ValidationHelper.cs           # 数据验证
│
└── Assets/                           # 图片/图标资源
    ├── app.ico
    └── ...
```

---

## 模块映射关系

### API 数据获取

**脚本中**：
```javascript
async function loadData(period, signal) {
  const [summaryJson, amountJson, costJson] = await Promise.all([
    fetchJson("/api/v0/users/get_user_summary", signal),
    fetchJson(`/api/v0/usage/amount?${query}`, signal),
    fetchJson(`/api/v0/usage/cost?${query}`, signal),
  ]);
  // ...
}
```

**WPF 中** → `Services/API/DeepSeekApiClient.cs`:
- 三个独立方法：`GetUserSummaryAsync()`、`GetUsageAmountAsync()`、`GetUsageCostAsync()`
- `AuthTokenProvider` 负责 Token 获取和刷新

---

### 状态管理

**脚本中**：全局 `state` 对象 + localStorage
```javascript
const state = {
  selectedPeriod: "",
  sectionVisible: loadSectionVisible(),
  keyDetailData: null,
  autoRefreshInterval: loadAutoRefreshInterval(),
  // ...
};
```

**WPF 中**：
- `Services/Storage/ConfigurationManager.cs` — 读写 AppData（替代 localStorage）
- `ViewModels/*` — MVVM 对象属性（INotifyPropertyChanged）绑定到 XAML

---

### UI 渲染

**脚本中**：纯 HTML + ECharts
```javascript
function renderSummary(data) {
  // HTML 字符串 + DOM 操作
}
const chart = echarts.init(container);
chart.setOption(option);
```

**WPF 中**：
- **Views/** — XAML 声明式布局
- **ViewModels/** — 数据源 + 命令处理
- **图表库选择**：
  - 推荐 **OxyPlot** 或 **Syncfusion Charts**（WPF native）
  - 或保留 Web 嵌入（WebView2 + ECharts）— 见 [WebView2 折衷](#webview2折衷)

---

### 持久化存储

**脚本中**：`localStorage.setItem()` / `localStorage.getItem()`
```javascript
function saveKeyDetailData() {
  const payload = { data, unitPrices, updateTime };
  localStorage.setItem("dsapi_plus_key_detail", JSON.stringify(payload));
}
```

**WPF 中** → `Services/Storage/LocalStorageService.cs`:
- 存储位置：`%APPDATA%\DeepSeekUsageUI\config.json`（Windows 标准做法）
- 实现：JSON 序列化 + 文件 I/O

---

### 配置管理

**脚本中**：单独的 load/save 函数对
```javascript
function loadGroupByModel() { /* ... */ }
function saveGroupByModel() { /* ... */ }
```

**WPF 中** → `Models/Configuration.cs` + `Services/Storage/ConfigurationManager.cs`:
```csharp
public class AppSettings
{
    public bool GroupByModel { get; set; }
    public int AutoRefreshInterval { get; set; }
    public KeyFilterOptions KeyFilter { get; set; }
    public List<SubscriptionConfig> Subscriptions { get; set; }
    // ...
}
```

---

## 核心类设计指南

### 1. Models（数据模型）

**ApiModels.cs**：与 API 响应 1:1 对应
```csharp
public record UserSummaryResponse(
    long CurrentToken,
    decimal TotalUsage,
    decimal MonthlyUsage,
    List<MonthlyCost> MonthlyCosts,
    List<Wallet> NormalWallets
);

public record UsageAmount(
    string Model,
    int Request,
    int Response,
    int PromptCacheHit,
    int PromptCacheMiss
);
```

**UIModels.cs**：绑定到 UI 的视图模型（与 ViewModel 的区别）
```csharp
public class SummaryUIModel : ObservableObject
{
    [ObservableProperty]
    string currentTokenDisplay;
    
    [ObservableProperty]
    List<ChartSeriesData> monthlyTrendData;
}
```

**Configuration.cs**：应用配置
```csharp
public class AppSettings
{
    public string ApiToken { get; set; }
    public bool AutoRefresh { get; set; }
    public int RefreshIntervalMs { get; set; }
    // ...
}
```

---

### 2. ViewModels（视图模型）

基类：使用 `ObservableObject` from `CommunityToolkit.Mvvm`

**MainWindowViewModel.cs**：
```csharp
public partial class MainWindowViewModel : ObservableObject
{
    private readonly IDeepSeekApiClient _apiClient;
    private readonly IConfigurationManager _configManager;
    
    [ObservableProperty]
    string selectedPeriod = GetCurrentPeriod();
    
    [ObservableProperty]
    SummaryUIModel summaryData;
    
    [ObservableProperty]
    bool isLoading;
    
    [RelayCommand]
    async Task RefreshData()
    {
        try
        {
            IsLoading = true;
            var period = SelectedPeriod;
            var (summary, usage, cost) = await _apiClient.LoadDataAsync(period);
            // 数据聚合 + UI 模型转换
            SummaryData = MapToUIModel(summary, usage, cost);
        }
        finally { IsLoading = false; }
    }
    
    [RelayCommand]
    void OpenSettings() { /* ... */ }
}
```

---

### 3. Services（业务逻辑）

**DeepSeekApiClient.cs**：
```csharp
public interface IDeepSeekApiClient
{
    Task<UserSummaryResponse> GetUserSummaryAsync(CancellationToken ct = default);
    Task<UsageAmountResponse> GetUsageAmountAsync(string period, CancellationToken ct = default);
    Task<UsageCostResponse> GetUsageCostAsync(string period, CancellationToken ct = default);
}

public class DeepSeekApiClient : IDeepSeekApiClient
{
    private readonly HttpClient _http;
    private readonly IAuthTokenProvider _tokenProvider;
    
    public async Task<UsageAmountResponse> GetUsageAmountAsync(string period, CancellationToken ct)
    {
        var token = await _tokenProvider.GetValidTokenAsync();
        var url = $"/api/v0/usage/amount?{ToQueryString(period)}";
        // ...
    }
}
```

**ConfigurationManager.cs**：
```csharp
public interface IConfigurationManager
{
    AppSettings Load();
    void Save(AppSettings settings);
    T GetSetting<T>(string key, T defaultValue = default);
    void SetSetting<T>(string key, T value);
}

public class ConfigurationManager : IConfigurationManager
{
    private readonly string _configPath = Path.Combine(
        Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData),
        "DeepSeekUsageUI", "config.json");
    
    public AppSettings Load()
    {
        if (!File.Exists(_configPath)) return new AppSettings();
        var json = File.ReadAllText(_configPath);
        return JsonConvert.DeserializeObject<AppSettings>(json) ?? new AppSettings();
    }
    
    public void Save(AppSettings settings)
    {
        Directory.CreateDirectory(Path.GetDirectoryName(_configPath));
        var json = JsonConvert.SerializeObject(settings, Formatting.Indented);
        File.WriteAllText(_configPath, json);
    }
}
```

---

### 4. Views（XAML 视图）

**MVVM 绑定模式**：

```xaml
<Window x:Class="DeepSeekUsageUI.Views.MainWindow"
        xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
        Title="DeepSeek 用量分析" Width="1200" Height="800">
    
    <Window.DataContext>
        <vm:MainWindowViewModel />
    </Window.DataContext>
    
    <Grid>
        <!-- 顶部控制栏 -->
        <StackPanel Orientation="Horizontal" Margin="16">
            <TextBlock Text="选择月份:" Margin="0,0,8,0" VerticalAlignment="Center"/>
            <ComboBox ItemsSource="{Binding AvailablePeriods}"
                      SelectedItem="{Binding SelectedPeriod}"/>
            <Button Command="{Binding RefreshDataCommand}" Margin="8,0,0,0">刷新</Button>
            <Button Command="{Binding OpenSettingsCommand}" Margin="8,0,0,0">设置</Button>
        </StackPanel>
        
        <!-- 统计摘要 -->
        <local:StatsSummary DataContext="{Binding SummaryData}" Margin="16,50,16,16"/>
        
        <!-- 加载指示器 -->
        <ProgressBar IsIndeterminate="{Binding IsLoading}" Height="2"/>
    </Grid>
</Window>
```

---

## 关键功能迁移清单

| 功能 | 脚本文件位置 | WPF 实现位置 | 状态 |
|------|-----------|-----------|------|
| API 数据获取 | `loadData()` | `DeepSeekApiClient.cs` | ⏳ TODO |
| Token 管理 | `getStoredAuthToken()` | `AuthTokenProvider.cs` | ⏳ TODO |
| 数据聚合 | `normalizeAmount()` 等 | `Services/DataAggregation.cs` | ⏳ TODO |
| 本地存储 | `localStorage.*` | `ConfigurationManager.cs` | ⏳ TODO |
| 图表绘制 | ECharts | OxyPlot 或 WebView2 | 📋 决策中 |
| 订阅推送 | `SubscriptionService` 在脚本中 | `Services/Export/SubscriptionService.cs` | ⏳ TODO |
| CSV/JSON 导出 | 脚本的拼接字符串 | `Services/Export/ExportService.cs` | ⏳ TODO |
| 自动刷新 | `setInterval()` | `DispatcherTimer` | ⏳ TODO |
| Key 筛选 | 脚本 filter 逻辑 | `Dialogs/KeyFilterDialog.xaml` | ⏳ TODO |

---

## 开发流程

### Phase 1：基础架构 ✓
- [x] 创建 .csproj 和筛选器文件
- [x] 定义 Models（API 响应 + UI 绑定）
- [ ] 配置 Dependency Injection

### Phase 2：核心服务
- [ ] 实现 `DeepSeekApiClient` — HTTP 请求 + 响应解析
- [ ] 实现 `ConfigurationManager` — 配置持久化
- [ ] 实现 `AuthTokenProvider` — Token 获取和验证

### Phase 3：UI 框架
- [ ] 主窗口 XAML + ViewModel
- [ ] 统计摘要面板
- [ ] Key 明细表格

### Phase 4：数据聚合 & 图表
- [ ] 数据规范化 + 聚合逻辑
- [ ] 选择图表库（OxyPlot / Syncfusion / WebView2）
- [ ] 实现各类图表组件

### Phase 5：高级功能
- [ ] 订阅推送（钉钉/邮件）
- [ ] 导出功能（CSV/JSON/ZIP）
- [ ] 自动刷新机制

### Phase 6：打磨 & 发布
- [ ] 错误处理 + 日志
- [ ] 国际化（可选）
- [ ] 单元测试
- [ ] 构建 installer / MSIX

---

## 技术决策

### WebView2 折衷方案

如果使用 **OxyPlot/Syncfusion** 需要完全重写图表逻辑，可考虑：

**临时方案**：嵌入 WebView2 运行 ECharts
```csharp
public class WebChartRenderer
{
    private readonly WebView2Control _webView;
    
    public async Task RenderChartAsync(ChartOption option)
    {
        var html = GenerateEChartsHTML(option);
        await _webView.NavigateToString(html);
    }
    
    private string GenerateEChartsHTML(ChartOption option)
    {
        return $"""
            <!DOCTYPE html>
            <html>
            <head><script src="https://cdn.jsdelivr.net/npm/echarts@5.6.0/dist/echarts.min.js"></script></head>
            <body style="margin:0">
                <div id="chart" style="width:100%;height:100%;"></div>
                <script>
                    var chart = echarts.init(document.getElementById('chart'));
                    chart.setOption({JsonConvert.SerializeObject(option)});
                </script>
            </body>
            </html>
            """;
    }
}
```

**优点**：快速迁移，图表代码无需改动
**缺点**：额外开销，不够原生

**最终推荐**：选择 OxyPlot（轻量级、纯 C# native）

---

## 约定俗成

### 命名规范
- 类名：PascalCase（如 `MainWindowViewModel`）
- 属性：PascalCase（如 `SelectedPeriod`）
- 私有字段：`_camelCase`（如 `_apiClient`）
- XAML 资源键：`camelCase`（如 `primaryBrush`）

### 代码组织
- 每个逻辑单元一个文件
- ViewModel 文件名与 View 对应（如 `MainWindowViewModel` ↔ `MainWindow.xaml`）
- Service 接口定义在同一文件开头

### UI/ViewModel 分离
- ViewModel 不能引用 View
- View 只通过数据绑定与 ViewModel 交互
- Code-behind 仅做必要的 XAML 初始化

---

## 环境设置

### 必要软件
- Visual Studio 2022 或更新版本（Community/Professional）
- .NET 8 SDK
- Windows 10/11

### 项目配置
```bash
# 恢复依赖
dotnet restore

# 构建
dotnet build -c Release

# 运行
dotnet run
```

### 调试快捷方式
```bash
# Hot reload
dotnet watch run

# 单元测试
dotnet test
```

---

## 参考资源

- [WPF 官方文档](https://learn.microsoft.com/en-us/dotnet/desktop/wpf/)
- [MVVM Community Toolkit](https://github.com/CommunityToolkit/dotnet/tree/main/components/MVVM)
- [OxyPlot WPF Documentation](https://oxyplot.readthedocs.io/)
- [Newtonsoft.Json 文档](https://www.newtonsoft.com/json/help/html/Introduction.htm)

---

## 常见问题

### Q: 为什么不继续使用 Greasy​monkey 脚本？
A: 脚本运行在浏览器环境，功能受限（无本地文件访问、无系统通知等）。独立应用可以：
- 无依赖运行（不依赖浏览器）
- 系统级集成（任务栏、通知等）
- 更好的性能 + 离线缓存

### Q: 如何处理 API Token 安全？
A: Token 本地存储在 Windows Credential Manager（`Services/AuthTokenProvider.cs`）而非明文文件。

### Q: 可以在 MacOS/Linux 上运行吗？
A: 使用 .NET 8 理论可行，但 WPF 仅 Windows 原生。跨平台可改用 **MAUI** 或 **Avalonia**。

---

最后修改：2026-07-06
