namespace DeepSeekUsageUI.Services;

/// <summary>
/// DeepSeek API 客户端接口
/// </summary>
public interface IDeepSeekApiClient
{
    /// <summary>
    /// 获取用户概览（余额、总用量等）
    /// </summary>
    Task<UserSummaryResponse> GetUserSummaryAsync(CancellationToken ct = default);

    /// <summary>
    /// 获取指定月份用量（按模型/Key 分布）
    /// </summary>
    Task<UsageAmountResponse> GetUsageAmountAsync(string period, CancellationToken ct = default);

    /// <summary>
    /// 获取指定月份成本（按模型/Key 分布）
    /// </summary>
    Task<UsageCostResponse> GetUsageCostAsync(string period, CancellationToken ct = default);

    /// <summary>
    /// 一次性加载完整月份数据
    /// </summary>
    Task<(UserSummaryResponse, UsageAmountResponse, UsageCostResponse)> LoadFullMonthDataAsync(
        string period, CancellationToken ct = default);
}

/// <summary>
/// Token 管理接口
/// </summary>
public interface IAuthTokenProvider
{
    /// <summary>
    /// 获取有效的 API Token
    /// </summary>
    Task<string> GetValidTokenAsync(CancellationToken ct = default);

    /// <summary>
    /// 设置 Token（手动输入）
    /// </summary>
    void SetToken(string token);

    /// <summary>
    /// 清除已缓存的 Token
    /// </summary>
    void ClearToken();

    /// <summary>
    /// 获取 Token 来源（用于调试）
    /// </summary>
    string GetTokenSource();
}

public record UserSummaryResponse(
    long CurrentToken,
    decimal TotalUsage,
    decimal MonthlyUsage,
    List<MonthlyCost> MonthlyCosts,
    List<Wallet> NormalWallets,
    List<Wallet> BonusWallets
);

public record MonthlyCost(string Currency, decimal Amount);
public record Wallet(string Type, decimal Balance, string Currency);

public record UsageAmountResponse(
    string Period,
    List<ModelUsageItem> Models,
    List<DailyUsageItem> Days,
    AggregateUsage Aggregate
);

public record ModelUsageItem(
    string Model,
    int Request,
    int ResponseToken,
    int PromptCacheMissToken,
    int PromptCacheHitToken
);

public record DailyUsageItem(
    string Date,
    List<ModelUsageItem> Models,
    AggregateUsage Aggregate
);

public record AggregateUsage(
    int Request,
    int Response,
    int PromptMiss,
    int PromptHit,
    int Tokens
);

public record UsageCostResponse(
    string Period,
    List<CurrencyBlock> Currencies
);

public record CurrencyBlock(
    string Currency,
    decimal TotalAmount,
    List<ModelCostItem> ModelCosts,
    List<DailyCostItem> Days
);

public record ModelCostItem(string Model, decimal Amount);
public record DailyCostItem(string Date, decimal Amount);
