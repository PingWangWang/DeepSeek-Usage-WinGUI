package service

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// KeyDetailData Key 级别的详细数据
type KeyDetailData struct {
	Key            string                 `json:"key"`
	Model          string                 `json:"model"`
	Requests       int                    `json:"requests"`
	Tokens         int                    `json:"tokens"`
	PromptMiss     int                    `json:"prompt_miss"`     // 未缓存 Token
	PromptHit      int                    `json:"prompt_hit"`      // 缓存命中 Token
	Completion     int                    `json:"completion"`      // 输出 Token
	Cost           float64                `json:"cost"`
	CacheMissCost  float64                `json:"cache_miss_cost"` // 未缓存成本
	CacheHitCost   float64                `json:"cache_hit_cost"`  // 缓存成本
	CompletionCost float64                `json:"completion_cost"` // 输出成本
	DailyCosts     map[string]float64     `json:"daily_costs"`     // 按日期的费用
}

// KeyDetailResult 导入结果
type KeyDetailResult struct {
	Data            []*KeyDetailData        `json:"data"`
	UnitPrices      map[string]UnitPrice    `json:"unit_prices"`     // 按模型的单价
	UpdateTime      string                  `json:"update_time"`
	DailyData       map[string]*DailyDetail `json:"daily_data"`      // 按日期统计
	KeyCount        int                     `json:"key_count"`
	ModelCount      int                     `json:"model_count"`
	TotalCost       float64                 `json:"total_cost"`
}

// UnitPrice 模型单价
type UnitPrice struct {
	Model       string  `json:"model"`
	PromptMiss  float64 `json:"prompt_miss"`  // 未缓存 Prompt 的单价
	PromptHit   float64 `json:"prompt_hit"`   // 缓存命中 Prompt 的单价
	Completion  float64 `json:"completion"`   // Completion 的单价
}

// DailyDetail 每日明细
type DailyDetail struct {
	Date      string                 `json:"date"`
	TotalCost float64                `json:"total_cost"`
	ByCost    map[string]float64     `json:"by_key"`    // 按 Key 的成本
}

// ExportService 导出/导入服务
type ExportService struct {
	cfg *config.Config
}

// NewExportService 创建导出服务
func NewExportService(cfg *config.Config) *ExportService {
	return &ExportService{cfg: cfg}
}

// ImportKeyDetail 从 ZIP 文件导入 Key 明细
func (e *ExportService) ImportKeyDetail(zipPath string) (*KeyDetailResult, error) {
	if zipPath == "" {
		return nil, fmt.Errorf("zip path is required")
	}

	// 检查文件是否存在
	if _, err := os.Stat(zipPath); err != nil {
		util.Logger.Error("ZIP file not found", "path", zipPath, "error", err)
		return nil, fmt.Errorf("file not found: %s", zipPath)
	}

	// 打开 ZIP 文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		util.Logger.Error("Failed to open ZIP", "error", err)
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	result := &KeyDetailResult{
		Data:       []*KeyDetailData{},
		UnitPrices: make(map[string]UnitPrice),
		UpdateTime: time.Now().Format(time.RFC3339),
		DailyData:  make(map[string]*DailyDetail),
	}

	// 解析所有 CSV 文件（支持多个月份的导出）
	for _, file := range reader.File {
		if !strings.HasPrefix(file.Name, "amount-") || !strings.HasSuffix(file.Name, ".csv") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			util.Logger.Warn("Failed to open CSV in ZIP", "file", file.Name, "error", err)
			continue
		}

		if err := e.parseCSV(rc, result); err != nil {
			rc.Close()
			util.Logger.Warn("Failed to parse CSV", "file", file.Name, "error", err)
			continue
		}
		rc.Close()
	}

	// 后处理：计算聚合数据
	e.aggregateData(result)

	util.Logger.Info("Key detail imported", "keys", result.KeyCount, "models", result.ModelCount, "cost", result.TotalCost)
	return result, nil
}

// parseCSV 解析单个 CSV 文件
func (e *ExportService) parseCSV(reader io.Reader, result *KeyDetailResult) error {
	csvReader := csv.NewReader(reader)

	// 跳过首行（表头）
	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// 识别列索引
	colIdx := e.parseCSVHeader(header)
	if !colIdx.isValid() {
		return fmt.Errorf("invalid CSV format")
	}

	// 读取数据行
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			util.Logger.Warn("Error reading CSV row", "error", err)
			continue
		}

		if len(record) <= colIdx.Date {
			continue
		}

		// 解析行数据
		item := e.parseCSVRow(record, colIdx)
		if item == nil {
			continue
		}

		result.Data = append(result.Data, item)
	}

	return nil
}

// CSVColumnIndex CSV 列索引
type CSVColumnIndex struct {
	Date           int
	Key            int
	Model          int
	PromptMiss     int
	PromptHit      int
	Completion     int
	CacheMissCost  int
	CacheHitCost   int
	CompletionCost int
}

func (c *CSVColumnIndex) isValid() bool {
	return c.Date >= 0 && c.Key >= 0 && c.Model >= 0
}

// parseCSVHeader 解析 CSV 表头，返回列索引
func (e *ExportService) parseCSVHeader(header []string) CSVColumnIndex {
	idx := CSVColumnIndex{
		Date:           -1,
		Key:            -1,
		Model:          -1,
		PromptMiss:     -1,
		PromptHit:      -1,
		Completion:     -1,
		CacheMissCost:  -1,
		CacheHitCost:   -1,
		CompletionCost: -1,
	}

	for i, col := range header {
		col = strings.TrimSpace(strings.ToLower(col))
		switch {
		case strings.Contains(col, "date"):
			idx.Date = i
		case strings.Contains(col, "key"):
			idx.Key = i
		case strings.Contains(col, "model"):
			idx.Model = i
		case strings.Contains(col, "prompt") && strings.Contains(col, "miss"):
			idx.PromptMiss = i
		case strings.Contains(col, "prompt") && strings.Contains(col, "hit"):
			idx.PromptHit = i
		case strings.Contains(col, "completion") && !strings.Contains(col, "cost"):
			idx.Completion = i
		case strings.Contains(col, "cache_miss_cost") || (strings.Contains(col, "miss") && strings.Contains(col, "cost")):
			idx.CacheMissCost = i
		case strings.Contains(col, "cache_hit_cost") || (strings.Contains(col, "hit") && strings.Contains(col, "cost")):
			idx.CacheHitCost = i
		case strings.Contains(col, "completion") && strings.Contains(col, "cost"):
			idx.CompletionCost = i
		}
	}

	return idx
}

// parseCSVRow 解析 CSV 行
func (e *ExportService) parseCSVRow(record []string, colIdx CSVColumnIndex) *KeyDetailData {
	if colIdx.Date < 0 || colIdx.Key < 0 || len(record) <= colIdx.Model {
		return nil
	}

	item := &KeyDetailData{
		Key:        strings.TrimSpace(record[colIdx.Key]),
		Model:      strings.TrimSpace(record[colIdx.Model]),
		DailyCosts: make(map[string]float64),
	}

	if colIdx.Date >= 0 && colIdx.Date < len(record) {
		item.DailyCosts[strings.TrimSpace(record[colIdx.Date])] = 0
	}

	// 解析数值字段
	if colIdx.PromptMiss >= 0 && colIdx.PromptMiss < len(record) {
		item.PromptMiss = parseInt(record[colIdx.PromptMiss])
	}
	if colIdx.PromptHit >= 0 && colIdx.PromptHit < len(record) {
		item.PromptHit = parseInt(record[colIdx.PromptHit])
	}
	if colIdx.Completion >= 0 && colIdx.Completion < len(record) {
		item.Completion = parseInt(record[colIdx.Completion])
	}

	item.Tokens = item.PromptMiss + item.PromptHit + item.Completion

	// 解析费用
	if colIdx.CacheMissCost >= 0 && colIdx.CacheMissCost < len(record) {
		item.CacheMissCost = parseFloat(record[colIdx.CacheMissCost])
	}
	if colIdx.CacheHitCost >= 0 && colIdx.CacheHitCost < len(record) {
		item.CacheHitCost = parseFloat(record[colIdx.CacheHitCost])
	}
	if colIdx.CompletionCost >= 0 && colIdx.CompletionCost < len(record) {
		item.CompletionCost = parseFloat(record[colIdx.CompletionCost])
	}

	item.Cost = item.CacheMissCost + item.CacheHitCost + item.CompletionCost

	return item
}

// aggregateData 聚合数据
func (e *ExportService) aggregateData(result *KeyDetailResult) {
	keyMap := make(map[string]*KeyDetailData)
	modelSet := make(map[string]bool)
	dailyMap := make(map[string]*DailyDetail)

	for _, item := range result.Data {
		// 按 (Key, Model) 聚合
		key := item.Key + "|" + item.Model
		if existing, ok := keyMap[key]; ok {
			existing.PromptMiss += item.PromptMiss
			existing.PromptHit += item.PromptHit
			existing.Completion += item.Completion
			existing.Tokens += item.Tokens
			existing.Cost += item.Cost
			existing.CacheMissCost += item.CacheMissCost
			existing.CacheHitCost += item.CacheHitCost
			existing.CompletionCost += item.CompletionCost
		} else {
			keyMap[key] = item
		}

		modelSet[item.Model] = true

		// 按日期聚合
		for date, cost := range item.DailyCosts {
			if daily, ok := dailyMap[date]; ok {
				daily.TotalCost += cost
				daily.ByCost[item.Key] += cost
			} else {
				dailyMap[date] = &DailyDetail{
					Date:      date,
					TotalCost: cost,
					ByCost:    map[string]float64{item.Key: cost},
				}
			}
		}

		result.TotalCost += item.Cost
	}

	// 替换结果数据
	result.Data = make([]*KeyDetailData, 0, len(keyMap))
	keySet := make(map[string]bool)
	for _, item := range keyMap {
		result.Data = append(result.Data, item)
		keySet[item.Key] = true
	}

	// 计算单价（按模型平均）
	modelCosts := make(map[string]struct{ tokens int; cost float64 })
	for _, item := range result.Data {
		if m, ok := modelCosts[item.Model]; ok {
			m.tokens += item.Tokens
			m.cost += item.Cost
			modelCosts[item.Model] = m
		} else {
			modelCosts[item.Model] = struct{ tokens int; cost float64 }{item.Tokens, item.Cost}
		}
	}

	for model, mc := range modelCosts {
		if mc.tokens > 0 {
			result.UnitPrices[model] = UnitPrice{
				Model:      model,
				PromptMiss: mc.cost / float64(mc.tokens),
				Completion: mc.cost / float64(mc.tokens),
			}
		}
	}

	result.DailyData = dailyMap
	result.KeyCount = len(keySet)
	result.ModelCount = len(modelSet)
}

// 辅助函数
func parseInt(s string) int {
	s = strings.TrimSpace(s)
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}
