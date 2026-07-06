package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// ReportData 报告数据
type ReportData struct {
	Period              string
	DayCost             float64
	MonthlyCost         float64
	AvgPrice            float64
	TotalTokens         int
	CacheHitRate        float64
	EstimatedTokens     int
	UserBalance         float64
	TopKeys             []KeyStat
	ModelBreakdown      []ModelStat
	SummaryMarkdown     string
}

// KeyStat Key 统计
type KeyStat struct {
	Key     string
	Cost    float64
	Tokens  int
	Percent float64
}

// ModelStat 模型统计
type ModelStat struct {
	Model  string
	Tokens int
	Cost   float64
}

// webhook 请求体定义

// DingTalkWebhookPayload 钉钉 Webhook 请求体
type DingTalkWebhookPayload struct {
	MsgType string                 `json:"msgtype"`
	Text    DingTalkText           `json:"text,omitempty"`
	Markdown DingTalkMarkdown      `json:"markdown,omitempty"`
	At      DingTalkAt             `json:"at,omitempty"`
}

type DingTalkText struct {
	Content string `json:"content"`
}

type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkAt struct {
	AtAll bool `json:"isAtAll"`
}

// FeiShuWebhookPayload 飞书 Webhook 请求体
type FeiShuWebhookPayload struct {
	MsgType string      `json:"msg_type"`
	Content interface{} `json:"content"`
}

type FeiShuTextContent struct {
	Text string `json:"text"`
}

type FeiShuMarkdownContent struct {
	MD string `json:"md"`
}

// WeChatWebhookPayload 企业微信 Webhook 请求体
type WeChatWebhookPayload struct {
	MsgType string      `json:"msgtype"`
	Text    WeChatText  `json:"text,omitempty"`
	Markdown WeChatMarkdown `json:"markdown,omitempty"`
}

type WeChatText struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type WeChatMarkdown struct {
	Content string `json:"content"`
}

// WebhookService Webhook 推送服务
type WebhookService struct {
	cfg *config.Config
	httpClient *http.Client
}

// NewWebhookService 创建 Webhook 服务
func NewWebhookService(cfg *config.Config) *WebhookService {
	return &WebhookService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendToWebhook 发送到 Webhook
func (w *WebhookService) SendToWebhook(ctx context.Context, webhookURL string, reportType string, content string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// 根据 URL 判断服务类型，自动选择格式
	if err := w.detectAndSend(ctx, webhookURL, reportType, content); err != nil {
		return err
	}

	return nil
}

// detectAndSend 自动检测 Webhook 类型并发送
func (w *WebhookService) detectAndSend(ctx context.Context, webhookURL string, reportType string, content string) error {
	var payload interface{}

	// 根据 URL 域名识别服务类型
	if strings.Contains(webhookURL, "dingtalk") || strings.Contains(webhookURL, "oapi.dingtalk") {
		payload = w.buildDingTalkPayload(reportType, content)
	} else if strings.Contains(webhookURL, "open.feishu") || strings.Contains(webhookURL, "feishu") {
		payload = w.buildFeiShuPayload(reportType, content)
	} else if strings.Contains(webhookURL, "qyapi.weixin") || strings.Contains(webhookURL, "wechat") {
		payload = w.buildWeChatPayload(reportType, content)
	} else {
		// 默认使用通用格式（JSON）
		payload = map[string]string{"message": content}
	}

	return w.postWebhook(ctx, webhookURL, payload)
}

// buildDingTalkPayload 构建钉钉消息体
func (w *WebhookService) buildDingTalkPayload(reportType string, content string) DingTalkWebhookPayload {
	if reportType == "markdown" {
		return DingTalkWebhookPayload{
			MsgType: "markdown",
			Markdown: DingTalkMarkdown{
				Title: "DeepSeek API 用量报告",
				Text:  content,
			},
			At: DingTalkAt{AtAll: false},
		}
	}

	return DingTalkWebhookPayload{
		MsgType: "text",
		Text:    DingTalkText{Content: content},
		At:      DingTalkAt{AtAll: false},
	}
}

// buildFeiShuPayload 构建飞书消息体
func (w *WebhookService) buildFeiShuPayload(reportType string, content string) FeiShuWebhookPayload {
	if reportType == "markdown" {
		return FeiShuWebhookPayload{
			MsgType: "post",
			Content: FeiShuMarkdownContent{MD: content},
		}
	}

	return FeiShuWebhookPayload{
		MsgType: "text",
		Content: FeiShuTextContent{Text: content},
	}
}

// buildWeChatPayload 构建企业微信消息体
func (w *WebhookService) buildWeChatPayload(reportType string, content string) WeChatWebhookPayload {
	if reportType == "markdown" {
		return WeChatWebhookPayload{
			MsgType:  "markdown",
			Markdown: WeChatMarkdown{Content: content},
		}
	}

	return WeChatWebhookPayload{
		MsgType: "text",
		Text:    WeChatText{Content: content},
	}
}

// postWebhook 通用 POST 请求
func (w *WebhookService) postWebhook(ctx context.Context, webhookURL string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		util.Logger.Error("Failed to marshal payload", "error", err)
		return fmt.Errorf("marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(data))
	if err != nil {
		util.Logger.Error("Failed to create request", "error", err)
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		util.Logger.Error("Failed to send webhook", "error", err)
		return fmt.Errorf("send failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		util.Logger.Error("Webhook returned error status", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("webhook error: status %d, body %s", resp.StatusCode, string(body))
	}

	util.Logger.Info("Webhook sent successfully", "status", resp.StatusCode)
	return nil
}

// BuildMarkdownReport 构建 Markdown 格式报告
func BuildMarkdownReport(data *ReportData) string {
	var buf bytes.Buffer

	buf.WriteString("# DeepSeek API 用量报告\n\n")
	buf.WriteString(fmt.Sprintf("📅 时间：%s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// 汇总卡片
	buf.WriteString("## 📊 汇总指标\n\n")
	buf.WriteString(fmt.Sprintf("| 指标 | 数值 |\n"))
	buf.WriteString(fmt.Sprintf("|------|------|\n"))
	buf.WriteString(fmt.Sprintf("| 当日费用 | ¥%.2f |\n", data.DayCost))
	buf.WriteString(fmt.Sprintf("| 月度费用 | ¥%.2f |\n", data.MonthlyCost))
	buf.WriteString(fmt.Sprintf("| 均价 | ¥%.4f/M Token |\n", data.AvgPrice))
	buf.WriteString(fmt.Sprintf("| 总 Token | %d |\n", data.TotalTokens))
	buf.WriteString(fmt.Sprintf("| 缓存命中率 | %.2f%% |\n\n", data.CacheHitRate*100))

	// 模型分布
	if len(data.ModelBreakdown) > 0 {
		buf.WriteString("## 🤖 模型分布\n\n")
		for _, m := range data.ModelBreakdown {
			buf.WriteString(fmt.Sprintf("- **%s**：%d Tokens (¥%.2f)\n", m.Model, m.Tokens, m.Cost))
		}
		buf.WriteString("\n")
	}

	// Top Keys
	if len(data.TopKeys) > 0 {
		buf.WriteString("## 🔑 Top Keys\n\n")
		for i, k := range data.TopKeys {
			if i >= 5 { // 仅显示前 5 个
				break
			}
			buf.WriteString(fmt.Sprintf("%d. %s：¥%.2f (%.1f%%)\n", i+1, k.Key, k.Cost, k.Percent*100))
		}
		buf.WriteString("\n")
	}

	buf.WriteString("---\n")
	buf.WriteString("📌 _本报告由 DeepSeek Usage+ 自动生成_\n")

	return buf.String()
}
