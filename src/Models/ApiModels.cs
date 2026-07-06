namespace DeepSeekUsageUI.Models;

/// <summary>
/// API 响应模型 — 直接映射 DeepSeek API 返回数据
/// </summary>

public record UserSummaryResponse(
    long CurrentToken,
    decimal TotalUsage,
    decimal MonthlyUsage,
    List<MonthlyCost> MonthlyCosts,
    List<Wallet> NormalWallets,
    List<Wallet> BonusWallets
);

public record MonthlyCost(
    string Currency,
    decimal Amount
);

public record Wallet(
    string Type,
    decimal Balance,
    string Currency
);

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

public record ModelCostItem(
    string Model,
    decimal Amount
);

public record DailyCostItem(
    string Date,
    decimal Amount
);
