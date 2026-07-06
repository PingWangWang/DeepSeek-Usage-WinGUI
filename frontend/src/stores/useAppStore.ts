import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

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

  // 初始化
  const init = () => {
    // 从后端获取配置
    console.log('App store initialized')
  }

  // 计算属性
  const isDark = computed(() => theme.value === 'dark')

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
