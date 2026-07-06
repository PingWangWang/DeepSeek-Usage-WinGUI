namespace DeepSeekUsageUI.Services;

/// <summary>
/// 配置管理接口
/// </summary>
public interface IConfigurationManager
{
    /// <summary>
    /// 加载所有配置
    /// </summary>
    AppSettings Load();

    /// <summary>
    /// 保存所有配置
    /// </summary>
    void Save(AppSettings settings);

    /// <summary>
    /// 获取单个配置值
    /// </summary>
    T? GetSetting<T>(string key);

    /// <summary>
    /// 设置单个配置值
    /// </summary>
    void SetSetting<T>(string key, T value);
}

/// <summary>
/// 本地存储接口（缓存、临时数据）
/// </summary>
public interface ILocalStorageService
{
    /// <summary>
    /// 保存对象到存储
    /// </summary>
    Task SaveAsync<T>(string key, T value);

    /// <summary>
    /// 读取对象
    /// </summary>
    Task<T?> LoadAsync<T>(string key);

    /// <summary>
    /// 删除存储项
    /// </summary>
    Task DeleteAsync(string key);

    /// <summary>
    /// 清空所有存储
    /// </summary>
    Task ClearAsync();
}

public class AppSettings { /* 详见 Models/Configuration.cs */ }
