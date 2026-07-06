package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/api"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/service"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/task"
)

// App 应用主体，暴露给 Wails 前端调用的方法
type App struct {
	version         string
	ctx             context.Context
	config          *config.Config
	deepseekClient  *api.DeepSeekClient
	analyticsSvc    *service.AnalyticsService
	subscriptionSvc *service.SubscriptionService
	webhookSvc      *service.WebhookService
	exportSvc       *service.ExportService
	scheduler       *task.Scheduler
	updater         *task.Updater
}

// NewApp 创建新应用实例
func NewApp(version string) *App {
	return &App{
		version: version,
	}
}

// Startup 应用启动时的初始化逻辑
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// 加载配置
	cfg := config.LoadConfig()
	a.config = cfg

	// 初始化 API 客户端
	a.deepseekClient = api.NewDeepSeekClient(cfg.Token)

	// 初始化业务服务
	a.analyticsSvc = service.NewAnalyticsService(a.deepseekClient, cfg)
	a.subscriptionSvc = service.NewSubscriptionService(cfg)
	a.webhookSvc = service.NewWebhookService(cfg)
	a.exportSvc = service.NewExportService(cfg)

	// 初始化定时任务调度器
	a.scheduler = task.NewScheduler()

	// 初始化更新检查器
	a.updater = task.NewUpdater(a.version)

	// 启动后台定时任务（如有配置）
	a.startBackgroundTasks()
}

// Shutdown 应用关闭时的清理逻辑
func (a *App) Shutdown(ctx context.Context) {
	if a.scheduler != nil {
		a.scheduler.Stop()
	}
}

// startBackgroundTasks 启动后台定时任务
func (a *App) startBackgroundTasks() {
	// 启动自动刷新定时器（如果配置了）
	if a.config.AutoRefreshInterval > 0 {
		a.scheduler.RegisterAutoRefresh(
			time.Duration(a.config.AutoRefreshInterval)*time.Millisecond,
			func() {
				// 这里可以调用前端刷新数据的方法，或直接更新缓存
				// 目前先留空，前端可手动触发刷新
			},
		)
	}

	// 启动订阅检查定时器
	for _, sub := range a.config.Subscriptions {
		a.scheduler.RegisterSubscriptionCheck(&sub, func() {
			// 检查订阅是否应该发送报告
		})
	}
}

// GetDashboard 获取仪表盘数据
func (a *App) GetDashboard(period string) (interface{}, error) {
	return a.analyticsSvc.GetDashboard(a.ctx, period)
}

// SendSubscription 发送订阅报告
func (a *App) SendSubscription(subID string) error {
	// 获取订阅配置
	var sub *config.Subscription
	for i := range a.config.Subscriptions {
		if a.config.Subscriptions[i].ID == subID {
			sub = &a.config.Subscriptions[i]
			break
		}
	}

	if sub == nil {
		return fmt.Errorf("subscription not found: %s", subID)
	}

	// 获取最新数据
	data, err := a.analyticsSvc.GetDashboard(a.ctx, "current_month")
	if err != nil {
		return err
	}

	// 构建报告
	report := service.BuildMarkdownReport(&service.ReportData{
		Period:          "current_month",
		DayCost:         data.SummaryMetrics.DayCost,
		MonthlyCost:     data.SummaryMetrics.MonthCost,
		AvgPrice:        data.SummaryMetrics.AvgPrice,
		TotalTokens:     data.SummaryMetrics.TotalTokens,
		CacheHitRate:    data.SummaryMetrics.CacheHitRate,
		EstimatedTokens: data.EstimatedTokens,
		UserBalance:     data.UserBalance,
	})

	// 发送到 Webhook
	if sub.Type == "webhook" {
		return a.webhookSvc.SendToWebhook(a.ctx, sub.WebhookURL, sub.Format, report)
	}

	return nil
}

// SetAutoRefresh 设置自动刷新间隔
func (a *App) SetAutoRefresh(interval int) error {
	a.config.AutoRefreshInterval = interval
	a.config.Save()
	a.startBackgroundTasks()
	return nil
}

// CheckUpdate 检查版本更新
func (a *App) CheckUpdate() (interface{}, error) {
	return a.updater.CheckUpdate(a.ctx)
}

// ImportKeyDetail 导入 Key 明细 ZIP 文件
func (a *App) ImportKeyDetail(zipPath string) (interface{}, error) {
	return a.exportSvc.ImportKeyDetail(zipPath)
}

// GetVersion 获取当前应用版本
func (a *App) GetVersion() string {
	return a.version
}

// UpdateToken 更新 API Token
func (a *App) UpdateToken(token string) error {
	a.config.Token = token
	a.config.Save()
	a.deepseekClient.SetToken(token)
	return nil
}

