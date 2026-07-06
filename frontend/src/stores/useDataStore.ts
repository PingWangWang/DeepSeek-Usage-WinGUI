import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as App from '@/../../wailsjs/go/backend/App'

export interface DashboardData {
  period: string
  summary_metrics: SummaryMetrics
  daily_data: DailyRecord[]
  model_breakdown: ModelStats[]
  cache_metrics: CacheMetrics
  user_balance: number
  estimated_tokens: number
}

export interface SummaryMetrics {
  day_cost: number
  month_cost: number
  avg_price: number
  total_tokens: number
  cache_hit_rate: number
}

export interface DailyRecord {
  date: string
  requests: number
  tokens: number
  cache_hit: number
  cache_miss: number
  cost: number
}

export interface ModelStats {
  model: string
  requests: number
  tokens: number
  cache_hit: number
  cache_miss: number
  cost: number
  percentage: number
}

export interface CacheMetrics {
  total_cache_hit: number
  total_cache_miss: number
  hit_rate: number
}

export const useDataStore = defineStore('data', () => {
  const loading = ref(false)
  const error = ref<string | null>(null)
  const data = ref<DashboardData | null>(null)
  const selectedPeriod = ref('current_month')

  // 获取仪表盘数据
  const fetchDashboard = async (period: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await App.GetDashboard(period)
      if (!response) {
        throw new Error('从后端获取数据失败：服务返回空值')
      }
      data.value = response
      error.value = null
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Unknown error'
      error.value = errorMsg
      console.error('Failed to fetch dashboard:', err)
      data.value = null
    } finally {
      loading.value = false
    }
  }

  // 切换周期
  const setPeriod = (period: string) => {
    selectedPeriod.value = period
    fetchDashboard(period)
  }

  return {
    loading,
    error,
    data,
    selectedPeriod,
    fetchDashboard,
    setPeriod,
  }
})
