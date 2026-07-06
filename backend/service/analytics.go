package service

import (
	"context"
	"sort"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/api"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// AnalyticsService 数据聚合服务
type AnalyticsService struct {
	api *api.DeepSeekClient
	cfg *config.Config
}

// NewAnalyticsService 创建新的分析服务
func NewAnalyticsService(apiClient *api.DeepSeekClient, cfg *config.Config) *AnalyticsService {
	return &AnalyticsService{
		api: apiClient,
		cfg: cfg,
	}
}

// DashboardData 仪表盘数据
type DashboardData struct {
	Period          string                  `json:"period"`
	SummaryMetrics  SummaryMetrics          `json:"summary_metrics"`
	DailyData       []DailyRecord           `json:"daily_data"`
	ModelBreakdown  []ModelStats            `json:"model_breakdown"`
	CacheMetrics    CacheMetrics            `json:"cache_metrics"`
	UserBalance     float64                 `json:"user_balance"`
	EstimatedTokens int                     `json:"estimated_tokens"`
}

// SummaryMetrics 汇总指标
type SummaryMetrics struct {
	DayCost         float64 `json:"day_cost"`         // 当日费用
	MonthCost       float64 `json:"month_cost"`       // 月度费用
	AvgPrice        float64 `json:"avg_price"`        // 均价（每百万 Token）
	TotalTokens     int     `json:"total_tokens"`     // 总 Token 用量
	CacheHitRate    float64 `json:"cache_hit_rate"`   // 缓存命中率
}

// DailyRecord 日常数据
type DailyRecord struct {
	Date      string `json:"date"`
	Requests  int    `json:"requests"`
	Tokens    int    `json:"tokens"`
	CacheHit  int    `json:"cache_hit"`
	CacheMiss int    `json:"cache_miss"`
	Cost      float64 `json:"cost"`
}

// ModelStats 模型统计
type ModelStats struct {
	Model      string  `json:"model"`
	Requests   int     `json:"requests"`
	Tokens     int     `json:"tokens"`
	CacheHit   int     `json:"cache_hit"`
	CacheMiss  int     `json:"cache_miss"`
	Cost       float64 `json:"cost"`
	Percentage float64 `json:"percentage"`
}

// CacheMetrics 缓存指标
type CacheMetrics struct {
	TotalCacheHit  int     `json:"total_cache_hit"`
	TotalCacheMiss int     `json:"total_cache_miss"`
	HitRate        float64 `json:"hit_rate"`
}

// GetDashboard 获取仪表盘数据
func (s *AnalyticsService) GetDashboard(ctx context.Context, period string) (*DashboardData, error) {
	// 获取三个 API 的数据
	amount, err := s.api.GetUsageAmount(ctx, period)
	if err != nil {
		util.Logger.Error("Failed to get amount data", "error", err)
		return nil, err
	}

	cost, err := s.api.GetUsageCost(ctx, period)
	if err != nil {
		util.Logger.Error("Failed to get cost data", "error", err)
		return nil, err
	}

	summary, err := s.api.GetUserSummary(ctx)
	if err != nil {
		util.Logger.Error("Failed to get summary data", "error", err)
		return nil, err
	}

	// 合并和聚合数据
	dashboard := s.aggregateData(amount, cost, summary, period)
	util.Logger.Info("Dashboard data aggregated", "period", period)

	return dashboard, nil
}

// aggregateData 合并数据
func (s *AnalyticsService) aggregateData(
	amount *api.AmountResponse,
	cost *api.CostResponse,
	summary *api.SummaryResponse,
	period string,
) *DashboardData {
	data := &DashboardData{
		Period:          period,
		DailyData:       []DailyRecord{},
		ModelBreakdown:  []ModelStats{},
		UserBalance:     summary.Data.BizData.Balance,
		EstimatedTokens: s.estimateAvailableTokens(summary.Data.BizData.Balance),
	}

	// 聚合日常数据
	dailyMap := make(map[string]*DailyRecord)
	if amount != nil && len(amount.Data.BizData.Data) > 0 {
		for _, day := range amount.Data.BizData.Data {
			record := &DailyRecord{
				Date:      day.Date,
				Requests:  day.Times,
				Tokens:    day.Tokens.CacheHit + day.Tokens.CacheMiss + day.Tokens.Completion,
				CacheHit:  day.Tokens.CacheHit,
				CacheMiss: day.Tokens.CacheMiss,
			}
			dailyMap[day.Date] = record
		}
	}

	// 合并费用数据
	if cost != nil && len(cost.Data.BizData.Data) > 0 {
		for _, day := range cost.Data.BizData.Data {
			if record, ok := dailyMap[day.Date]; ok {
				record.Cost = day.TotalCost
			}
		}
	}

	// 转为列表并排序
	for _, record := range dailyMap {
		data.DailyData = append(data.DailyData, *record)
	}
	sort.Slice(data.DailyData, func(i, j int) bool {
		return data.DailyData[i].Date < data.DailyData[j].Date
	})

	// 聚合模型数据
	modelMap := make(map[string]*ModelStats)
	if amount != nil && len(amount.Data.BizData.Data) > 0 {
		for _, day := range amount.Data.BizData.Data {
			for _, model := range day.Models {
				if _, ok := modelMap[model.Model]; !ok {
					modelMap[model.Model] = &ModelStats{Model: model.Model}
				}
				ms := modelMap[model.Model]
				ms.Requests += model.Times
				ms.Tokens += model.Tokens.CacheHit + model.Tokens.CacheMiss + model.Tokens.Completion
				ms.CacheHit += model.Tokens.CacheHit
				ms.CacheMiss += model.Tokens.CacheMiss
			}
		}
	}

	// 合并模型费用
	if cost != nil && len(cost.Data.BizData.Data) > 0 {
		for _, day := range cost.Data.BizData.Data {
			for _, model := range day.Models {
				if ms, ok := modelMap[model.Model]; ok {
					ms.Cost += model.Cost
				}
			}
		}
	}

	// 转为列表并计算占比
	totalTokens := 0
	for _, ms := range modelMap {
		totalTokens += ms.Tokens
	}
	for _, ms := range modelMap {
		if totalTokens > 0 {
			ms.Percentage = float64(ms.Tokens) / float64(totalTokens)
		}
		data.ModelBreakdown = append(data.ModelBreakdown, *ms)
	}
	sort.Slice(data.ModelBreakdown, func(i, j int) bool {
		return data.ModelBreakdown[i].Tokens > data.ModelBreakdown[j].Tokens
	})

	// 计算汇总指标
	data.SummaryMetrics = s.calculateSummary(data)

	// 计算缓存指标
	totalHit := 0
	totalMiss := 0
	for _, daily := range data.DailyData {
		totalHit += daily.CacheHit
		totalMiss += daily.CacheMiss
	}
	data.CacheMetrics = CacheMetrics{
		TotalCacheHit:  totalHit,
		TotalCacheMiss: totalMiss,
		HitRate:        float64(totalHit) / float64(totalHit+totalMiss),
	}

	return data
}

// calculateSummary 计算汇总指标
func (s *AnalyticsService) calculateSummary(data *DashboardData) SummaryMetrics {
	metrics := SummaryMetrics{}

	// 当日费用和月度费用
	today := time.Now().Format("2006-01-02")
	monthCost := 0.0
	for _, daily := range data.DailyData {
		monthCost += daily.Cost
		if daily.Date == today {
			metrics.DayCost = daily.Cost
		}
	}
	metrics.MonthCost = monthCost

	// 总 Token 和缓存命中率
	totalTokens := 0
	totalHit := 0
	totalMiss := 0
	for _, daily := range data.DailyData {
		totalTokens += daily.Tokens
		totalHit += daily.CacheHit
		totalMiss += daily.CacheMiss
	}
	metrics.TotalTokens = totalTokens
	if totalHit+totalMiss > 0 {
		metrics.CacheHitRate = float64(totalHit) / float64(totalHit+totalMiss)
	}

	// 均价（每百万 Token）
	if totalTokens > 0 {
		metrics.AvgPrice = (monthCost / float64(totalTokens)) * 1e6
	}

	return metrics
}

// estimateAvailableTokens 根据余额和均价估算可用 Token
func (s *AnalyticsService) estimateAvailableTokens(balance float64) int {
	if balance <= 0 {
		return 0
	}
	// 这里使用默认均价估算，实际应使用计算出的均价
	avgPrice := 0.05 // 假设每百万 Token 0.05 元
	availableTokens := (balance / avgPrice) * 1e6
	return int(availableTokens)
}

