package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

// 应用版本号（构建时通过 ldflags 注入）
var (
	AppVersion = "1.0.0"
	AppTitle   = "DeepSeek Usage+"
)

func main() {
	// 创建应用实例
	app := backend.NewApp(AppVersion)

	// 创建应用选项
	opts := &options.App{
		Title:  AppTitle,
		Width:  1400,
		Height: 900,
		Assets: assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			app.Shutdown(ctx)
		},
		Bind: []interface{}{
			app,
		},
	}

	// 运行应用
	err := wails.Run(opts)
	if err != nil {
		panic(err)
	}
}
