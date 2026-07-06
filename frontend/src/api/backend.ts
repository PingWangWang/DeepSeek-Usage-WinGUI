// API 调用管理
class BackendAPI {
  async getDashboard(period: string) {
    try {
      return await window.go.main.App.GetDashboard(period)
    } catch (error) {
      console.error('Failed to get dashboard:', error)
      throw error
    }
  }

  async sendSubscription(subID: string) {
    try {
      return await window.go.main.App.SendSubscription(subID)
    } catch (error) {
      console.error('Failed to send subscription:', error)
      throw error
    }
  }

  async setAutoRefresh(interval: number) {
    try {
      return await window.go.main.App.SetAutoRefresh(interval)
    } catch (error) {
      console.error('Failed to set auto refresh:', error)
      throw error
    }
  }

  async checkUpdate() {
    try {
      return await window.go.main.App.CheckUpdate()
    } catch (error) {
      console.error('Failed to check update:', error)
      throw error
    }
  }

  async importKeyDetail(zipPath: string) {
    try {
      return await window.go.main.App.ImportKeyDetail(zipPath)
    } catch (error) {
      console.error('Failed to import key detail:', error)
      throw error
    }
  }

  async getVersion() {
    try {
      return await window.go.main.App.GetVersion()
    } catch (error) {
      console.error('Failed to get version:', error)
      throw error
    }
  }

  async updateToken(token: string) {
    try {
      return await window.go.main.App.UpdateToken(token)
    } catch (error) {
      console.error('Failed to update token:', error)
      throw error
    }
  }
}

export const api = new BackendAPI()
