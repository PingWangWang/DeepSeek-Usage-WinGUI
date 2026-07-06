namespace DeepSeekUsageUI.Models;

/// <summary>
/// 应用配置模型 — 对应 AppData 中的配置文件
/// </summary>
public class AppSettings
{
    /// <summary>
    /// DeepSeek API Token（可选）
    /// </summary>
    public string? ApiToken { get; set; }

    /// <summary>
    /// 自动刷新间隔（毫秒）
    /// </summary>
    public int AutoRefreshIntervalMs { get; set; } = 0;

    /// <summary>
    /// 按模型分组显示
    /// </summary>
    public bool GroupByModel { get; set; } = false;

    /// <summary>
    /// 是否显示原生内容区块
    /// </summary>
    public bool ShowNativeContent { get; set; } = true;

    /// <summary>
    /// Key 筛选配置
    /// </summary>
    public KeyFilterConfig? KeyFilter { get; set; }

    /// <summary>
    /// 订阅推送配置列表
    /// </summary>
    public List<SubscriptionConfig> Subscriptions { get; set; } = [];

    /// <summary>
    /// 图表可见性配置
    /// </summary>
    public SectionVisibilityConfig SectionVisibility { get; set; } = new();

    /// <summary>
    /// 最后选择的月份 (yyyy-MM)
    /// </summary>
    public string? LastSelectedPeriod { get; set; }

    /// <summary>
    /// 最后加载成功的数据
    /// </summary>
    public KeyDetailData? CachedKeyDetailData { get; set; }
}

public class KeyFilterConfig
{
    public string Mode { get; set; } = "all"; // "all" | "include" | "exclude"
    public List<string> Keys { get; set; } = [];
}

public class SubscriptionConfig
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string Name { get; set; } = "";
    public string TargetType { get; set; } = ""; // "dingtalk" | "email" | "webhook"
    public string TargetUrl { get; set; } = "";
    public bool Enabled { get; set; } = true;
    public string ScheduleCron { get; set; } = "0 0 * * *";
    public ContentOptions ContentOptions { get; set; } = new();
}

public class ContentOptions
{
    public bool IncludeSummary { get; set; } = true;
    public bool IncludeTodayDetail { get; set; } = false;
    public bool IncludeMonthDetail { get; set; } = false;
    public bool IncludeChart { get; set; } = false;
    public bool SendAsScreenshot { get; set; } = false;
}

public class SectionVisibilityConfig
{
    public bool Requests { get; set; } = false;
    public bool Tokens { get; set; } = false;
    public bool CacheRate { get; set; } = false;
    public bool Composition { get; set; } = false;
    public bool Models { get; set; } = false;
}

/// <summary>
/// Key 明细聚合数据（从导出接口获取）
/// </summary>
public class KeyDetailData
{
    public List<KeyAggregate> Keys { get; set; } = [];
    public Dictionary<string, ModelUnitPrices> UnitPrices { get; set; } = [];
    public DateTime UpdateTime { get; set; }
    public KeyDailyData? DailyData { get; set; }
}

public class KeyAggregate
{
    public string Key { get; set; } = "";
    public int Request { get; set; }
    public int Response { get; set; }
    public int PromptMiss { get; set; }
    public int PromptHit { get; set; }
    public int Tokens => Response + PromptMiss + PromptHit;
    public double CacheHitRate
    {
        get
        {
            int total = PromptMiss + PromptHit;
            return total > 0 ? (double)PromptHit / total : 0;
        }
    }
}

public class ModelUnitPrices
{
    public decimal PromptMissPrice { get; set; }
    public decimal PromptHitPrice { get; set; }
    public decimal ResponsePrice { get; set; }
}

public class KeyDailyData
{
    public List<string> Dates { get; set; } = [];
    public List<DailySeries> Series { get; set; } = [];
    public List<DailySeries>? Requests { get; set; }
    public List<DailySeries>? Tokens { get; set; }
    public List<DailySeries>? Miss { get; set; }
    public List<DailySeries>? Hit { get; set; }
}

public class DailySeries
{
    public string Name { get; set; } = "";
    public List<int> Data { get; set; } = [];
}
