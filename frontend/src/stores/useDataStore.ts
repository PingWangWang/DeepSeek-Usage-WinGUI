import { defineStore } from 'pinia'
import { ref } from 'vue'

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
      // TODO: 调用后端 API
      // const response = await window.go.main.App.GetDashboard(period)
      // data.value = response
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error'
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
