import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 应用状态
  const version = ref('1.0.0')
  const token = ref('')
  const autoRefreshInterval = ref(0)
  const theme = ref<'light' | 'dark'>('light')

  // UI 状态
  const sectionVisible = ref({
    requests: true,
    tokens: true,
    cacheRate: true,
    composition: true,
    models: true,
  })

  // 初始化：从后端恢复配置
  const init = async () => {
    try {
      // 检查window.go是否可用
      if (!window.go?.main?.App?.GetConfig) {
        console.warn('Wails binding not ready, using localStorage fallback')
        loadFromLocalStorage()
        return
      }

      const config = await window.go.main.App.GetConfig()
      if (config) {
        token.value = config.token || ''
        autoRefreshInterval.value = config.auto_refresh_interval || 0
        if (config.section_visible) {
          sectionVisible.value = config.section_visible
        }
      }

      // 从localStorage恢复主题（后端不管理主题）
      const savedTheme = localStorage.getItem('app_theme') as 'light' | 'dark' | null
      if (savedTheme) {
        theme.value = savedTheme
      }

      console.log('App config loaded from backend')
    } catch (error) {
      console.error('Failed to load app config from backend:', error)
      loadFromLocalStorage()
    }
  }

  // localStorage备用加载
  const loadFromLocalStorage = () => {
    const savedToken = localStorage.getItem('app_token')
    const savedInterval = localStorage.getItem('app_auto_refresh_interval')
    const savedTheme = localStorage.getItem('app_theme') as 'light' | 'dark' | null
    const savedSectionVisible = localStorage.getItem('app_section_visible')

    if (savedToken) token.value = savedToken
    if (savedInterval) autoRefreshInterval.value = parseInt(savedInterval, 10)
    if (savedTheme) theme.value = savedTheme
    if (savedSectionVisible) sectionVisible.value = JSON.parse(savedSectionVisible)
  }

  // 计算属性
  const isDark = computed(() => theme.value === 'dark')

  // 监听配置变化，自动持久化
  watch(() => token.value, async (newToken) => {
    localStorage.setItem('app_token', newToken)
    try {
      if (window.go?.main?.App?.SaveConfig) {
        await window.go.main.App.SaveConfig({
          token: newToken,
          auto_refresh_interval: autoRefreshInterval.value,
          section_visible: sectionVisible.value,
        })
      }
    } catch (error) {
      console.error('Failed to save token to backend:', error)
    }
  })

  watch(() => autoRefreshInterval.value, async (newInterval) => {
    localStorage.setItem('app_auto_refresh_interval', newInterval.toString())
    try {
      if (window.go?.main?.App?.SaveConfig) {
        await window.go.main.App.SaveConfig({
          token: token.value,
          auto_refresh_interval: newInterval,
          section_visible: sectionVisible.value,
        })
      }
    } catch (error) {
      console.error('Failed to save auto refresh interval to backend:', error)
    }
  })

  watch(() => theme.value, (newTheme) => {
    localStorage.setItem('app_theme', newTheme)
  })

  watch(() => sectionVisible.value, async (newVisible) => {
    localStorage.setItem('app_section_visible', JSON.stringify(newVisible))
    try {
      if (window.go?.main?.App?.SaveConfig) {
        await window.go.main.App.SaveConfig({
          token: token.value,
          auto_refresh_interval: autoRefreshInterval.value,
          section_visible: newVisible,
        })
      }
    } catch (error) {
      console.error('Failed to save section visibility to backend:', error)
    }
  }, { deep: true })

  return {
    version,
    token,
    autoRefreshInterval,
    theme,
    sectionVisible,
    isDark,
    init,
  }
})
