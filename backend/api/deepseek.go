package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

const (
	APIBaseURL = "https://api.deepseek.com"
	APITimeout = 30 * time.Second
)

// DeepSeekClient DeepSeek API 客户端
type DeepSeekClient struct {
	baseURL string
	token   string
	client  *resty.Client
}

// NewDeepSeekClient 创建新的 API 客户端
func NewDeepSeekClient(token string) *DeepSeekClient {
	client := resty.New().
		SetBaseURL(APIBaseURL).
		SetTimeout(APITimeout).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))

	return &DeepSeekClient{
		baseURL: APIBaseURL,
		token:   token,
		client:  client,
	}
}

// AmountResponse API 响应：用量数据
type AmountResponse struct {
	Code int `json:"code"`
	Data struct {
		BizData struct {
			Data []struct {
				Date  string `json:"date"`
				Times int    `json:"times"`
				Tokens struct {
					CacheHit   int `json:"cache_hit"`
					CacheMiss  int `json:"cache_miss"`
					Completion int `json:"completion"`
				} `json:"tokens"`
				Models []struct {
					Model  string `json:"model"`
					Times  int    `json:"times"`
					Tokens struct {
						CacheHit   int `json:"cache_hit"`
						CacheMiss  int `json:"cache_miss"`
						Completion int `json:"completion"`
					} `json:"tokens"`
				} `json:"models"`
			} `json:"data"`
		} `json:"biz_data"`
	} `json:"data"`
	Message string `json:"message"`
}

// CostResponse API 响应：费用数据
type CostResponse struct {
	Code int `json:"code"`
	Data struct {
		BizData struct {
			Data []struct {
				Date      string  `json:"date"`
				TotalCost float64 `json:"total_cost"`
				Models    []struct {
					Model     string  `json:"model"`
					Cost      float64 `json:"cost"`
					CacheHit  float64 `json:"cache_hit"`
					CacheMiss float64 `json:"cache_miss"`
					Completion float64 `json:"completion"`
				} `json:"models"`
			} `json:"data"`
		} `json:"biz_data"`
	} `json:"data"`
	Message string `json:"message"`
}

// SummaryResponse API 响应：用户信息
type SummaryResponse struct {
	Code int `json:"code"`
	Data struct {
		BizData struct {
			UserID        string  `json:"user_id"`
			Balance       float64 `json:"balance"`
			TotalTokens   int     `json:"total_tokens"`
			UsedTokens    int     `json:"used_tokens"`
			AvailableTokens int   `json:"available_tokens"`
		} `json:"biz_data"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetUsageAmount 获取用量数据
func (c *DeepSeekClient) GetUsageAmount(ctx context.Context, period string) (*AmountResponse, error) {
	var resp AmountResponse
	r, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("period", period).
		SetResult(&resp).
		Get("/usage/amount")

	if err != nil {
		util.Logger.Error("Failed to get usage amount", "error", err)
		return nil, err
	}

	if r.StatusCode() != 200 {
		util.Logger.Warn("API returned non-200 status", "code", r.StatusCode())
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("api error: %s (code %d)", resp.Message, resp.Code)
	}

	util.Logger.Debug("Got usage amount", "period", period, "days", len(resp.Data.BizData.Data))
	return &resp, nil
}

// GetUsageCost 获取费用数据
func (c *DeepSeekClient) GetUsageCost(ctx context.Context, period string) (*CostResponse, error) {
	var resp CostResponse
	r, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("period", period).
		SetResult(&resp).
		Get("/usage/cost")

	if err != nil {
		util.Logger.Error("Failed to get usage cost", "error", err)
		return nil, err
	}

	if r.StatusCode() != 200 {
		util.Logger.Warn("API returned non-200 status", "code", r.StatusCode())
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("api error: %s (code %d)", resp.Message, resp.Code)
	}

	util.Logger.Debug("Got usage cost", "period", period, "days", len(resp.Data.BizData.Data))
	return &resp, nil
}

// GetUserSummary 获取用户信息
func (c *DeepSeekClient) GetUserSummary(ctx context.Context) (*SummaryResponse, error) {
	var resp SummaryResponse
	r, err := c.client.R().
		SetContext(ctx).
		SetResult(&resp).
		Get("/users/get_user_summary")

	if err != nil {
		util.Logger.Error("Failed to get user summary", "error", err)
		return nil, err
	}

	if r.StatusCode() != 200 {
		util.Logger.Warn("API returned non-200 status", "code", r.StatusCode())
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("api error: %s (code %d)", resp.Message, resp.Code)
	}

	util.Logger.Debug("Got user summary", "balance", resp.Data.BizData.Balance)
	return &resp, nil
}

// SetToken 更新 API 令牌
func (c *DeepSeekClient) SetToken(token string) {
	c.token = token
	c.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}
