#!/usr/bin/env pwsh
# DeepSeek Usage+ 构建脚本 (PowerShell)

param(
    [string]$Output = "DeepSeek-Usage.exe",
    [int]$ProxyPort = 7897
)

$ErrorActionPreference = "Stop"

# 获取脚本所在目录，然后进入项目根目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# 进入项目目录
Push-Location $ProjectRoot

try {
    # 配置
    $VERSION = "1.0.0"

    # 颜色函数
    function Write-Success {
        Write-Host "✅ $args" -ForegroundColor Green
    }

    function Write-ErrorMsg {
        Write-Host "❌ $args" -ForegroundColor Red
    }

    function Write-Info {
        Write-Host "ℹ️  $args" -ForegroundColor Cyan
    }

    function Write-Step {
        Write-Host "📦 $args" -ForegroundColor Yellow
    }

    # 设置代理
    Write-Info "设置 Clash 代理 (端口 $ProxyPort)..."
    $env:HTTP_PROXY = "http://127.0.0.1:$ProxyPort"
    $env:HTTPS_PROXY = "http://127.0.0.1:$ProxyPort"
    $env:ALL_PROXY = "http://127.0.0.1:$ProxyPort"

    # 验证代理
    Write-Info "验证代理连接..."
    $proxyTest = Test-NetConnection -ComputerName github.com -Port 443 -InformationLevel Quiet -WarningAction SilentlyContinue
    if ($proxyTest) {
        Write-Success "代理连接正常"
    } else {
        Write-ErrorMsg "代理连接失败！请检查 Clash 是否运行"
        exit 1
    }

    # 1. 清理缓存
    Write-Step "清理 Go 缓存..."
    try {
        go clean -modcache
        Write-Success "缓存已清理"
    }
    catch {
        Write-ErrorMsg "清理缓存失败: $_"
    }

    # 2. 下载 Go 依赖
    Write-Step "下载 Go 依赖..."
    try {
        go mod tidy
        Write-Success "依赖下载成功"
    }
    catch {
        Write-ErrorMsg "依赖下载失败: $_"
        Write-Info "重试..."
        Start-Sleep -Seconds 3
        go mod tidy
    }

    # 3. 构建前端
    Write-Step "构建前端..."
    try {
        Push-Location frontend

        # 检查 node_modules
        if (-not (Test-Path node_modules)) {
            Write-Info "安装 npm 依赖..."
            npm install
        }

        npm run build
        Write-Success "前端构建成功"

        Pop-Location
    }
    catch {
        Write-ErrorMsg "前端构建失败: $_"
        Pop-Location
        exit 1
    }

    # 4. 检查/安装 Wails CLI
    Write-Step "检查 Wails CLI..."
    try {
        $wailsCheck = wails --version 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Wails not installed"
        }
        Write-Success "Wails CLI 已安装"
    }
    catch {
        Write-Info "安装 Wails CLI..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
        Write-Success "Wails CLI 安装成功"
    }

    # 5. 构建 Go 应用
    Write-Step "编译 Go 应用..."
    try {
        wails build -webview2 Embed -o $Output

        # 从 build/bin 目录复制最终输出到项目根目录
        $BuildOutput = Join-Path $ProjectRoot "build" "bin" $Output
        if (Test-Path $BuildOutput) {
            Copy-Item -Path $BuildOutput -Destination (Join-Path $ProjectRoot $Output) -Force
            $FileSize = (Get-Item (Join-Path $ProjectRoot $Output)).Length / 1MB
            Write-Success "编译成功！"
            Write-Host ""
            Write-Host "📦 输出文件: $Output" -ForegroundColor Green
            Write-Host "📊 大小: $([math]::Round($FileSize, 2))MB" -ForegroundColor Green
            Write-Host "🔖 版本: $VERSION" -ForegroundColor Green
            Write-Host ""
            Write-Host "可以直接运行: .\$Output" -ForegroundColor Cyan
        }
        else {
            throw "Wails 构建失败或输出文件未创建"
        }
    }
    catch {
        Write-ErrorMsg "编译失败: $_"
        exit 1
    }

    Write-Info "构建完成！"
}
finally {
    Pop-Location
}

