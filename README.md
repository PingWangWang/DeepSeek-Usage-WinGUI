# DeepSeek Usage+ — Windows Desktop Analytics Dashboard

> Real-time usage analytics for DeepSeek API users. Single EXE file, zero configuration, with interactive charts, auto-refresh, Webhook notifications, and Key-level cost breakdown.

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](RELEASES)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078d4.svg)](README.md)
[![Tech Stack](https://img.shields.io/badge/Stack-Go%20%2B%20Vue3%20%2B%20Wails-orange.svg)](README.md)

**中文文档** → [README_ZH.md](README_ZH.md)

## ✨ Features

- **📊 Real-time Dashboard** — Daily cost, monthly cost, unit price, total tokens, cache hit rate at a glance
- **📈 Interactive Charts** — 5 ECharts visualizations: request trends, token composition, cache rate, model distribution
- **🔄 Auto-refresh** — 30 seconds to 1 hour intervals, background polling
- **🔔 Webhook Notifications** — Support for DingTalk, Feishu, WeCom; scheduled report delivery
- **🔑 Key Detail Import** — Parse DeepSeek export ZIP files, aggregate by Key/Model, calculate per-unit costs
- **🌙 Dark Theme** — Auto-detect system theme, reduces eye strain
- **📦 Single-file Distribution** — Windows EXE, no dependencies, double-click to run
- **🔐 100% Local** — API tokens stored only in local AppData, zero cloud upload

## 🎯 Who Should Use This

- 👨‍💻 **DeepSeek API Users** — Monitor API usage and costs in real-time
- 💼 **Multi-Key Managers** — Track costs by API Key, optimize budget allocation
- 🤖 **LLM Application Developers** — Analyze cache hit rates, optimize prompt design, reduce costs
- 📊 **Data-Driven Teams** — Daily/weekly/monthly cost trends, make informed decisions

## 📥 Installation & Quick Start

### Download (Recommended)

1. **Download** → Get the latest `DeepSeek-Usage.exe` from [Releases](https://github.com/PingWangWang/DeepSeek-Usage-WinGUI/releases)
2. **Run** → Double-click the EXE to launch
3. **Configure** → Go to Settings tab, paste your DeepSeek API Token (from [platform.deepseek.com](https://platform.deepseek.com))
4. **View** → Return to Dashboard, click "Refresh" to load data

### Build from Source

**Requirements**:
- Go 1.21+
- Node.js 18+
- HTTP Proxy (e.g., Clash Verge)

**Build Steps**:
```bash
git clone https://github.com/PingWangWang/DeepSeek-Usage-WinGUI.git
cd DeepSeek-Usage-WinGUI

# Execute build script (requires proxy for GitHub access)
.\scripts\build.ps1 -ProxyPort 7897

# Or manually:
cd frontend && npm install && npm run build && cd ..
go mod tidy
go build -o DeepSeek-Usage.exe .
```

## 🧭 Usage Examples

### View Dashboard

Default view shows current month's data:

- **Summary Cards** — Daily cost, monthly cost, unit price, total tokens, cache hit rate
- **Charts** — Request trends, token usage, cache rate, model distribution
- **Toolbar** — Month picker, chart toggles, auto-refresh settings

### Import Key Details

1. Visit [DeepSeek Platform](https://platform.deepseek.com/usage), click "Export" to download ZIP
2. In the app, select "Import Key Details", choose the ZIP file
3. App automatically extracts, parses, and aggregates the data
4. View per-key cost, token usage, and cache metrics

### Configure Webhook Notifications

Add subscriptions in Settings:

1. **Delivery Method** — DingTalk/Feishu/WeCom (Webhook URL)
2. **Content Options** — Cost summary, token composition, cache rate, etc.
3. **Frequency** — Custom interval, daily, weekly, or monthly
4. **Key Filtering** — Generate reports for specific API keys

## 📁 Project Structure

```
DeepSeek-Usage-WinGUI/
├── main.go                          # Wails app entry point
├── backend/                         # Go backend
│   ├── app.go                      # Application core
│   ├── api/                        # DeepSeek API client
│   ├── service/                    # Business logic (aggregation, webhooks, imports)
│   ├── task/                       # Scheduled tasks, version checking
│   └── util/                       # Utility functions
├── frontend/                        # Vue 3 frontend
│   ├── src/
│   │   ├── components/             # UI components (dashboard, settings, charts)
│   │   ├── stores/                 # Pinia state management
│   │   └── utils/                  # ECharts, API utilities
│   └── package.json
├── scripts/
│   ├── build.ps1                   # PowerShell build script
│   └── build.sh                    # Shell build script
├── docs/                           # Documentation
│   ├── DEVELOPMENT.md              # Developer guide
│   └── BUILD.md                    # Build instructions
└── README.md / README_ZH.md        # This file
```

## 🛠️ Development

### Start Development Server

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Start dev server with hot reload
wails dev
```

### Build for Production

```bash
# Compile to single Windows EXE
wails build -webview2 Embed -o DeepSeek-Usage.exe

# Or use the build script
.\scripts\build.ps1
```

### Architecture

**Backend (Go)**:
- `API Client` — Call DeepSeek endpoints (usage, cost, user info)
- `Analytics Service` — Aggregate by date, model, and API key
- `Webhook Service` — Support for DingTalk, Feishu, WeCom
- `Export Service` — ZIP extraction, CSV parsing, data aggregation

**Frontend (Vue 3)**:
- `Dashboard` — 5 ECharts visualizations + summary cards
- `Settings` — Token config, auto-refresh, version check
- `Pinia Store` — Application state and data caching

## ❓ FAQ

**Q: Where is my API token stored?**  
A: Stored locally in `%APPDATA%/DeepSeek-Usage/config.json`. Never uploaded to cloud.

**Q: How do I update to a newer version?**  
A: Click "Check Update" in Settings. Download the new EXE from GitHub Releases and replace the old one.

**Q: What if Webhook notifications fail?**  
A: Verify the Webhook URL is correct and your network can access external URLs. Check logs in `%APPDATA%/DeepSeek-Usage/logs/`.

**Q: Does this support macOS/Linux?**  
A: Currently Windows EXE only. The codebase uses Wails (cross-platform), but macOS/Linux builds haven't been tested yet.

## 🚀 Roadmap

- [ ] macOS/Linux builds
- [ ] Auto-update feature (download + replace)
- [ ] Screenshot upload to image hosting (ImgBB)
- [ ] SQLite backend for offline mode
- [ ] Additional Webhook services (Slack, Discord, etc.)

## 🤝 Contributing

Issues and PRs welcome!

### Report a Bug
- Describe reproduction steps
- Attach logs from `%APPDATA%/DeepSeek-Usage/logs/`
- Include Windows and Go versions

### Submit a PR
- Fork and create a feature branch
- Follow existing code style
- Update relevant docs
- Describe your changes in the PR

## 📊 Stats

| Metric | Value |
|--------|-------|
| Total Lines of Code | ~6,000 |
| External Dependencies | Minimal (std lib + Wails only) |
| EXE Size | ~70-100 MB |
| Startup Time | <2 seconds |
| First Query | 1-2 seconds |

## 📝 License

MIT License — Free to use, modify, and distribute. Attribution required.

---

## 🔗 Quick Links

- **DeepSeek Home** — https://www.deepseek.com
- **API Documentation** — https://platform.deepseek.com/api-docs
- **Usage Dashboard** — https://platform.deepseek.com/usage
- **GitHub Repository** — https://github.com/PingWangWang/DeepSeek-Usage-WinGUI

## Acknowledgments

Thanks to DeepSeek for the open API and to the community for feedback and ideas.

---

**Issues & Feedback** → [GitHub Issues](https://github.com/PingWangWang/DeepSeek-Usage-WinGUI/issues)

**Stay Updated** → ⭐ Star this repository
