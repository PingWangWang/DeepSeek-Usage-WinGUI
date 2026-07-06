package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// SubscriptionService 订阅推送服务
type SubscriptionService struct {
	cfg *config.Config
}

// NewSubscriptionService 创建新的订阅服务
func NewSubscriptionService(cfg *config.Config) *SubscriptionService {
	return &SubscriptionService{
		cfg: cfg,
	}
}

// SendReport 发送订阅报告
func (s *SubscriptionService) SendReport(ctx context.Context, subID string) error {
	// 查找订阅配置
	var sub *config.Subscription
	for i := range s.cfg.Subscriptions {
		if s.cfg.Subscriptions[i].ID == subID {
			sub = &s.cfg.Subscriptions[i]
			break
		}
	}

	if sub == nil {
		return fmt.Errorf("subscription not found: %s", subID)
	}

	if !sub.Enabled {
		util.Logger.Warn("Subscription is disabled", "id", subID)
		return nil
	}

	// 根据类型发送
	switch sub.Type {
	case "webhook":
		return s.sendToWebhook(ctx, sub)
	case "clipboard":
		return s.sendToClipboard(ctx, sub)
	case "preview":
		return s.sendToPreview(ctx, sub)
	default:
		return fmt.Errorf("unknown subscription type: %s", sub.Type)
	}
}

// sendToWebhook 发送到 Webhook
func (s *SubscriptionService) sendToWebhook(ctx context.Context, sub *config.Subscription) error {
	util.Logger.Info("Sending to webhook", "url", sub.WebhookURL)
	// TODO: 实现 Webhook 推送逻辑
	return nil
}

// sendToClipboard 发送到剪贴板
func (s *SubscriptionService) sendToClipboard(ctx context.Context, sub *config.Subscription) error {
	util.Logger.Info("Sending to clipboard")
	// TODO: 实现剪贴板复制逻辑
	return nil
}

// sendToPreview 预览模式
func (s *SubscriptionService) sendToPreview(ctx context.Context, sub *config.Subscription) error {
	util.Logger.Info("Preview mode", "id", sub.ID)
	// TODO: 实现预览逻辑
	return nil
}

// CheckAndSend 检查是否应该发送报告（定时任务回调）
func (s *SubscriptionService) CheckAndSend(ctx context.Context, sub *config.Subscription) error {
	if !sub.Enabled {
		return nil
	}

	now := time.Now()
	lastSent := s.cfg.SubscriptionLastSent[sub.ID]

	// 根据频率判断是否应该发送
	var shouldSend bool
	switch sub.Frequency {
	case "interval":
		// 检查间隔是否已过
		if lastSent == "" {
			shouldSend = true
		} else {
			lastSentTime, err := time.Parse(time.RFC3339, lastSent)
			if err != nil {
				shouldSend = true
			} else {
				nextSend := lastSentTime.Add(time.Duration(sub.Interval) * time.Minute)
				shouldSend = now.After(nextSend)
			}
		}

	case "daily":
		// 检查今天是否已发送
		if lastSent != "" {
			lastSentTime, _ := time.Parse(time.RFC3339, lastSent)
			if lastSentTime.Day() == now.Day() {
				shouldSend = false
			}
		}

	case "weekly":
		// 检查本周是否已发送
		if lastSent != "" {
			lastSentTime, _ := time.Parse(time.RFC3339, lastSent)
			if lastSentTime.Year() == now.Year() && lastSentTime.YearDay()/7 == now.YearDay()/7 {
				shouldSend = false
			}
		}

	case "monthly":
		// 检查本月是否已发送
		if lastSent != "" {
			lastSentTime, _ := time.Parse(time.RFC3339, lastSent)
			if lastSentTime.Year() == now.Year() && lastSentTime.Month() == now.Month() {
				shouldSend = false
			}
		}
	}

	if !shouldSend {
		return nil
	}

	// 发送报告
	if err := s.SendReport(ctx, sub.ID); err != nil {
		util.Logger.Error("Failed to send subscription", "id", sub.ID, "error", err)
		return err
	}

	// 更新最后发送时间
	s.cfg.SubscriptionLastSent[sub.ID] = now.Format(time.RFC3339)
	s.cfg.Save()

	util.Logger.Info("Subscription sent", "id", sub.ID)
	return nil
}
