namespace DeepSeekUsageUI.Services;

/// <summary>
/// 导出服务接口
/// </summary>
public interface IExportService
{
    /// <summary>
    /// 导出为 CSV
    /// </summary>
    Task<string> ExportToCsvAsync(string period, CancellationToken ct = default);

    /// <summary>
    /// 导出为 JSON
    /// </summary>
    Task<string> ExportToJsonAsync(string period, CancellationToken ct = default);

    /// <summary>
    /// 导出为 ZIP（包含 CSV + 截图）
    /// </summary>
    Task<byte[]> ExportToZipAsync(string period, CancellationToken ct = default);
}

/// <summary>
/// 截图导出服务接口
/// </summary>
public interface IScreenshotExporter
{
    /// <summary>
    /// 截图当前 UI
    /// </summary>
    Task<byte[]> CaptureUIAsync();

    /// <summary>
    /// 保存截图到文件
    /// </summary>
    Task SaveScreenshotAsync(string filePath, CancellationToken ct = default);
}

/// <summary>
/// 订阅推送服务接口
/// </summary>
public interface ISubscriptionService
{
    /// <summary>
    /// 发送订阅通知（钉钉/邮件/Webhook）
    /// </summary>
    Task SendSubscriptionAsync(SubscriptionConfig config, string content, byte[]? screenshot = null, CancellationToken ct = default);

    /// <summary>
    /// 验证订阅配置（测试连接）
    /// </summary>
    Task<bool> ValidateSubscriptionAsync(SubscriptionConfig config, CancellationToken ct = default);
}

public class SubscriptionConfig { /* 详见 Models/Configuration.cs */ }
