// API 调用管理 - 使用 Wails 生成的绑定
import * as App from '@/../../wailsjs/go/backend/App'

class BackendAPI {
  async getDashboard(period: string) {
    try {
      return await App.GetDashboard(period)
    } catch (error) {
      console.error('Failed to get dashboard:', error)
      throw error
    }
  }

  async sendSubscription(subID: string) {
    try {
      return await App.SendSubscription(subID)
    } catch (error) {
      console.error('Failed to send subscription:', error)
      throw error
    }
  }

  async setAutoRefresh(interval: number) {
    try {
      return await App.SetAutoRefresh(interval)
    } catch (error) {
      console.error('Failed to set auto refresh:', error)
      throw error
    }
  }

  async checkUpdate() {
    try {
      return await App.CheckUpdate()
    } catch (error) {
      console.error('Failed to check update:', error)
      throw error
    }
  }

  async importKeyDetail(zipPath: string) {
    try {
      return await App.ImportKeyDetail(zipPath)
    } catch (error) {
      console.error('Failed to import key detail:', error)
      throw error
    }
  }

  async getVersion() {
    try {
      return await App.GetVersion()
    } catch (error) {
      console.error('Failed to get version:', error)
      throw error
    }
  }

  async updateToken(token: string) {
    try {
      return await App.UpdateToken(token)
    } catch (error) {
      console.error('Failed to update token:', error)
      throw error
    }
  }

  async getConfig() {
    try {
      return await App.GetConfig()
    } catch (error) {
      console.error('Failed to get config:', error)
      throw error
    }
  }

  async saveConfig(configData: Record<string, any>) {
    try {
      return await App.SaveConfig(configData)
    } catch (error) {
      console.error('Failed to save config:', error)
      throw error
    }
  }
}

export const api = new BackendAPI()
