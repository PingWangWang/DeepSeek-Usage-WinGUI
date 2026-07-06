package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// Updater 版本更新检查
type Updater struct {
	currentVersion string
	repoURL        string
	httpClient     *http.Client
}

// UpdateInfo 更新信息
type UpdateInfo struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	HasUpdate      bool   `json:"has_update"`
	DownloadURL    string `json:"download_url"`
	ChangeLog      string `json:"change_log"`
	PublishedAt    string `json:"published_at"`
}

// GitHubRelease GitHub Release API 响应
type GitHubRelease struct {
	TagName         string `json:"tag_name"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	PublishedAt     string `json:"published_at"`
	Assets          []GitHubAsset `json:"assets"`
	Prerelease      bool   `json:"prerelease"`
}

// GitHubAsset GitHub Release 资源
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
	ContentType        string `json:"content_type"`
}

// NewUpdater 创建新的更新检查器
func NewUpdater(version string) *Updater {
	return &Updater{
		currentVersion: version,
		repoURL:        "https://api.github.com/repos/PingWangWang/DeepSeek-Usage-WinGUI/releases/latest",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckUpdate 检查是否有更新
func (u *Updater) CheckUpdate(ctx context.Context) (*UpdateInfo, error) {
	info := &UpdateInfo{
		CurrentVersion: u.currentVersion,
		LatestVersion:  u.currentVersion,
		HasUpdate:      false,
	}

	// 从 GitHub API 获取最新版本
	release, err := u.fetchLatestRelease(ctx)
	if err != nil {
		util.Logger.Warn("Failed to check update", "error", err)
		return info, nil // 返回当前版本信息，但不报错
	}

	if release == nil {
		return info, nil
	}

	info.LatestVersion = release.TagName
	info.ChangeLog = release.Body
	info.PublishedAt = release.PublishedAt

	// 查找 Windows EXE 下载链接
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, "DeepSeek-Usage") && strings.HasSuffix(asset.Name, ".exe") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	// 比较版本号
	if u.isNewer(release.TagName, u.currentVersion) {
		info.HasUpdate = true
	}

	util.Logger.Info("Update check completed", "current", u.currentVersion, "latest", info.LatestVersion, "has_update", info.HasUpdate)
	return info, nil
}

// fetchLatestRelease 从 GitHub API 获取最新版本
func (u *Updater) fetchLatestRelease(ctx context.Context) (*GitHubRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.repoURL, nil)
	if err != nil {
		return nil, err
	}

	// 添加用户代理（GitHub API 要求）
	req.Header.Set("User-Agent", "DeepSeek-Usage-WinGUI")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("parse json failed: %w", err)
	}

	return &release, nil
}

// isNewer 比较版本号，判断 newVersion 是否比 currentVersion 更新
// 使用简单的字符串比较（v1.2.0 vs v1.1.9）
func (u *Updater) isNewer(newVersion, currentVersion string) bool {
	// 移除 'v' 前缀
	newVer := strings.TrimPrefix(newVersion, "v")
	curVer := strings.TrimPrefix(currentVersion, "v")

	// 简单的版本比较：按字符比较
	// 更精确的做法是解析为 semver，但对于简单情况字符比较足够
	return strings.Compare(newVer, curVer) > 0
}

// DownloadAndUpdate 下载并更新（暂未实现）
func (u *Updater) DownloadAndUpdate(ctx context.Context, info *UpdateInfo) error {
	if !info.HasUpdate {
		return fmt.Errorf("no update available")
	}

	if info.DownloadURL == "" {
		return fmt.Errorf("download url not found")
	}

	// TODO: 实现自动更新逻辑
	// 1. 下载新版本 EXE 到临时目录
	// 2. 验证 SHA256 校验和
	// 3. 备份当前版本
	// 4. 替换 EXE（或创建 updater.exe 作为 wrapper）
	// 5. 重启应用

	util.Logger.Info("Update started", "target", info.LatestVersion)
	return fmt.Errorf("not implemented yet")
}
