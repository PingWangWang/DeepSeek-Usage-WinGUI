// Wails 运行时类型定义
declare global {
  namespace Window {
    const go: {
      main: {
        App: {
          GetDashboard(period: string): Promise<DashboardData>
          SendSubscription(subID: string): Promise<void>
          SetAutoRefresh(interval: number): Promise<void>
          CheckUpdate(): Promise<UpdateInfo>
          ImportKeyDetail(zipPath: string): Promise<KeyDetailResult>
          GetVersion(): Promise<string>
          UpdateToken(token: string): Promise<void>
        }
      }
    }
  }
}

export {}
