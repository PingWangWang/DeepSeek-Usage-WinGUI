package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// Config 应用全局配置
type Config struct {
	// API 认证
	Token string `json:"token"`

	// UI 状态
	SectionVisible          map[string]bool `json:"section_visible"`           // 图表显示/隐藏开关
	NativeContentVisible    bool            `json:"native_content_visible"`    // 原生内容显示
	KeyTableVisible         bool            `json:"key_table_visible"`         // Key 明细表格显示
	KeyDetailDailyVisible   bool            `json:"key_detail_daily_visible"`  // 每日明细显示
	KeyDetailChartVisible   bool            `json:"key_detail_chart_visible"`  // Key 费用分布图显示
	GroupByModel            bool            `json:"group_by_model"`            // 按模型分组

	// 刷新控制
	AutoRefreshInterval int `json:"auto_refresh_interval"` // 自动刷新间隔（毫秒）

	// Key 筛选
	KeyFilter KeyFilter `json:"key_filter"`

	// 订阅配置
	Subscriptions []Subscription `json:"subscriptions"`

	// 缓存数据
	LastPanelData    interface{} `json:"last_panel_data"`
	KeyDetailData    interface{} `json:"key_detail_data"`
	KeyDetailDailyData interface{} `json:"key_detail_daily_data"`
	KeyDetailUpdateTime string `json:"key_detail_update_time"`
	KeyUnitPrices    map[string]interface{} `json:"key_unit_prices"`

	// 订阅最后发送时间
	SubscriptionLastSent map[string]string `json:"subscription_last_sent"`
}

// KeyFilter Key 筛选配置
type KeyFilter struct {
	Mode string   `json:"mode"` // "all" 或 "custom"
	Keys []string `json:"keys"` // 选中的 Key 列表
}

// Subscription 订阅配置
type Subscription struct {
	ID         string                 `json:"id"`          // 订阅 ID（UUID）
	Name       string                 `json:"name"`        // 订阅名称
	Enabled    bool                   `json:"enabled"`     // 是否启用
	Type       string                 `json:"type"`        // "webhook" / "clipboard" / "preview"
	WebhookURL string                 `json:"webhook_url"` // Webhook URL（type=webhook 时）
	ImgBBKey   string                 `json:"imgbb_key"`   // ImgBB API Key（截图模式）
	Format     string                 `json:"format"`      // "markdown" / "image"
	Frequency  string                 `json:"frequency"`   // "interval" / "daily" / "weekly" / "monthly"
	Interval   int                    `json:"interval"`    // 间隔（分钟，frequency=interval 时）
	DayOfWeek  int                    `json:"day_of_week"` // 周几（frequency=weekly 时，0=Sunday）
	DayOfMonth int                    `json:"day_of_month"`// 每月几号（frequency=monthly 时）
	Time       string                 `json:"time"`        // 发送时间 HH:MM（daily/weekly/monthly）
	Content    SubscriptionContent    `json:"content"`     // 内容选项
	KeyFilter  KeyFilter              `json:"key_filter"`  // Key 筛选
}

// SubscriptionContent 订阅内容选项
type SubscriptionContent struct {
	IncludeSummary bool `json:"include_summary"`
	IncludeTokens  bool `json:"include_tokens"`
	IncludeCache   bool `json:"include_cache"`
	IncludeComposition bool `json:"include_composition"`
	IncludeKeyDetail bool `json:"include_key_detail"`
	TopNKeys       int  `json:"top_n_keys"` // 前 N 个 Key（0=全部）
}

var (
	configDir  string
	configPath string
)

// init 初始化配置路径
func init() {
	configDir = getConfigDir()
	configPath = filepath.Join(configDir, "config.json")
	os.MkdirAll(configDir, 0755)
}

// getConfigDir 获取配置目录（%APPDATA%/DeepSeek-Usage）
func getConfigDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		// Fallback: 使用用户主目录
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".deepseek-usage")
	}
	return filepath.Join(appData, "DeepSeek-Usage")
}

// LoadConfig 从文件加载配置，不存在则返回默认值
func LoadConfig() *Config {
	cfg := &Config{
		SectionVisible: map[string]bool{
			"requests":    true,
			"tokens":      true,
			"cacheRate":   true,
			"composition": true,
			"models":      true,
		},
		NativeContentVisible:   true,
		KeyTableVisible:        false,
		KeyDetailDailyVisible:  false,
		KeyDetailChartVisible:  true,
		GroupByModel:           false,
		AutoRefreshInterval:    0, // 默认关闭
		KeyFilter:              KeyFilter{Mode: "all"},
		Subscriptions:          []Subscription{},
		SubscriptionLastSent:   make(map[string]string),
		KeyUnitPrices:          make(map[string]interface{}),
	}

	// 尝试从文件加载
	data, err := os.ReadFile(configPath)
	if err != nil {
		util.Logger.Warn("Config file not found, using defaults", "path", configPath)
		return cfg
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		util.Logger.Error("Failed to parse config", "error", err)
		return cfg
	}

	return cfg
}

// Save 保存配置到文件
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		util.Logger.Error("Failed to marshal config", "error", err)
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		util.Logger.Error("Failed to write config", "error", err)
		return err
	}

	util.Logger.Info("Config saved", "path", configPath)
	return nil
}

// GetConfigDir 返回配置目录
func GetConfigDir() string {
	return configDir
}

// GetConfigPath 返回配置文件路径
func GetConfigPath() string {
	return configPath
}
