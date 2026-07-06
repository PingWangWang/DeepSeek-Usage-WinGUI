import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import * as App from '@/../../wailsjs/go/backend/App'

let tokenSaveTimer: ReturnType<typeof setTimeout> | null = null
let intervalSaveTimer: ReturnType<typeof setTimeout> | null = null
let visibilitySaveTimer: ReturnType<typeof setTimeout> | null = null

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
      const config = await App.GetConfig()
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

  // 监听Token变化，防抖保存
  watch(() => token.value, async (newToken) => {
    localStorage.setItem('app_token', newToken)
    if (tokenSaveTimer) clearTimeout(tokenSaveTimer)
    tokenSaveTimer = setTimeout(async () => {
      try {
        await App.SaveConfig({
          token: newToken,
          auto_refresh_interval: autoRefreshInterval.value,
          section_visible: sectionVisible.value,
        })
      } catch (error) {
        console.warn('Failed to save token to backend, using localStorage fallback:', error)
      }
    }, 300)
  })

  // 监听自动刷新间隔变化，防抖保存
  watch(() => autoRefreshInterval.value, async (newInterval) => {
    localStorage.setItem('app_auto_refresh_interval', newInterval.toString())
    if (intervalSaveTimer) clearTimeout(intervalSaveTimer)
    intervalSaveTimer = setTimeout(async () => {
      try {
        await App.SaveConfig({
          token: token.value,
          auto_refresh_interval: newInterval,
          section_visible: sectionVisible.value,
        })
      } catch (error) {
        console.warn('Failed to save auto refresh interval to backend, using localStorage fallback:', error)
      }
    }, 300)
  })

  // 监听主题变化，同步到localStorage和后端
  watch(() => theme.value, async (newTheme) => {
    localStorage.setItem('app_theme', newTheme)
    try {
      await App.SaveConfig({
        token: token.value,
        auto_refresh_interval: autoRefreshInterval.value,
        section_visible: sectionVisible.value,
        theme: newTheme,
      })
    } catch (error) {
      console.warn('Failed to save theme to backend, using localStorage fallback:', error)
    }
  })

  // 监听UI可见性变化，防抖保存
  watch(() => sectionVisible.value, async (newVisible) => {
    localStorage.setItem('app_section_visible', JSON.stringify(newVisible))
    if (visibilitySaveTimer) clearTimeout(visibilitySaveTimer)
    visibilitySaveTimer = setTimeout(async () => {
      try {
        await App.SaveConfig({
          token: token.value,
          auto_refresh_interval: autoRefreshInterval.value,
          section_visible: newVisible,
        })
      } catch (error) {
        console.warn('Failed to save section visibility to backend, using localStorage fallback:', error)
      }
    }, 300)
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
