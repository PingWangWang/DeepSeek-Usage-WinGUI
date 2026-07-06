#!/bin/bash
# DeepSeek Usage+ 构建脚本

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")

echo "🔨 Building DeepSeek Usage+ v${VERSION}..."

# 1. 构建前端
echo "📦 Building frontend..."
cd "$PROJECT_ROOT/frontend"
npm install
npm run build

# 2. 构建 Go 应用
echo "🏗️  Building Go application..."
cd "$PROJECT_ROOT"

# 注入版本号
LDFLAGS="-X main.AppVersion=${VERSION}"

# Windows AMD64
GOOS=windows GOARCH=amd64 go build \
  -ldflags "$LDFLAGS" \
  -o "DeepSeek-Usage.exe" \
  .

echo "✅ Build complete: DeepSeek-Usage.exe (v${VERSION})"
