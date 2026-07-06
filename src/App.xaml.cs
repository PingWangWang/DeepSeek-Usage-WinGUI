using System.Windows;

namespace DeepSeekUsageUI;

public partial class App : Application
{
    public App()
    {
        InitializeComponent();
    }

    protected override void OnStartup(StartupEventArgs e)
    {
        base.OnStartup(e);
        // TODO: 注册依赖注入容器
        // var services = new ServiceCollection();
        // services.AddSingleton<IDeepSeekApiClient, DeepSeekApiClient>();
        // services.AddSingleton<IConfigurationManager, ConfigurationManager>();
        // var provider = services.BuildServiceProvider();
    }
}
