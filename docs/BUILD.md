# 构建指南

## 环境要求

- Go 1.21+
- Node.js 18+
- Wails 2.8+
- Windows 10 或更新版本（用于编译 Windows EXE）

## 编译步骤

### 方法 1：使用 Wails CLI

```bash
# 开发模式（带热加载）
wails dev

# 生产构建（单文件 EXE）
wails build -webview2 Embed
```

输出文件：`build/bin/DeepSeek-Usage.exe`

### 方法 2：使用构建脚本

#### Windows

```batch
# 运行批处理脚本
scripts\build.bat
```

#### macOS/Linux

```bash
# 运行 Shell 脚本
./scripts/build.sh
```

## 编译参数说明

### Wails 构建标志

```bash
wails build [options]

# 常用选项
-webview2 Embed              # 嵌入 WebView2 运行时（推荐）
-webview2 Existing           # 依赖系统已安装的 WebView2
-clean                       # 清理前后端构建缓存
-o <name>                    # 输出 EXE 名称
-tags <tags>                 # 构建标签（如 dev, prod）
```

### Go 构建参数

```bash
go build -ldflags "-X main.AppVersion=1.0.0" .
```

通过 `-ldflags` 注入版本号到二进制文件。

## 输出文件大小

| 配置 | 大小 | 说明 |
|------|------|------|
| WebView2 Embed | 50-100MB | 包含 WebView2 运行时（推荐，首次运行快） |
| WebView2 Existing | 30-40MB | 依赖系统已安装运行时（用户需自行安装） |

## 发布流程

### GitHub Actions（自动）

1. 创建版本标签：
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. GitHub Actions 自动运行 `.github/workflows/release.yml`，生成 Release

### 手动发布

1. 编译 EXE：
```bash
scripts\build.bat
```

2. 创建 GitHub Release，上传 `DeepSeek-Usage.exe`

3. 用户应用启动时自动检查新版本

## 版本号管理

### 版本来源优先级

1. Git tag（推荐）：`git describe --tags --always`
2. 环境变量：`$VERSION`
3. 硬编码：`main.go` 中的 `AppVersion`

### 语义化版本

遵循 semver：`major.minor.patch`

- **major**：重大功能或破坏性变更
- **minor**：新功能，向后兼容
- **patch**：Bug 修复

示例：
- `v1.0.0` — 首个正式版本
- `v1.1.0` — 增加新功能
- `v1.0.1` — Bug 修复

## 优化编译

### 减少 EXE 体积

```bash
# 移除调试符号和符号表
go build -ldflags="-s -w" .

# 使用 UPX 压缩（可选）
upx -9 DeepSeek-Usage.exe
```

### 加快编译速度

```bash
# 并行编译
go build -p 4 .

# 缓存 Go 模块
go mod vendor
go build -mod vendor .
```

## 交叉编译

目前仅支持 Windows AMD64。未来可扩展到：

```bash
# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o DeepSeek-Usage-mac .

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o DeepSeek-Usage-linux .
```

## 故障排查

### 前端资源无法嵌入

- 确保 `frontend/dist/` 已生成
- 检查 `wails.json` 中的 `bin` 字段指向正确目录

### WebView2 相关错误

- 确认目标系统已安装 WebView2 运行时
- 使用 `-webview2 Embed` 时会自动包含运行时

### 编译失败

- 清理缓存：`wails build -clean`
- 检查 Go 版本：`go version` 应为 1.21+
- 检查依赖：`go mod tidy`

## CI/CD 配置

### GitHub Actions 示例

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - uses: actions/setup-node@v3
      - run: scripts\build.bat
      - uses: softprops/action-gh-release@v1
        with:
          files: DeepSeek-Usage.exe
```
