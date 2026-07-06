<template>
  <div class="dashboard">
    <Toolbar />

    <!-- 加载状态 -->
    <div v-if="dataStore.loading" class="loading">
      <p>📊 加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-if="dataStore.error" class="error">
      <p>❌ 加载失败: {{ dataStore.error }}</p>
      <button @click="refresh" class="btn">重试</button>
    </div>

    <!-- 数据展示 -->
    <template v-if="dataStore.data">
      <!-- 汇总卡片 -->
      <SummaryCards :metrics="dataStore.data.summary_metrics" />

      <!-- 图表区块 -->
      <div v-if="appStore.sectionVisible.requests" class="section">
        <LineChart
          chart-id="requests-chart"
          title="API 请求趋势"
          :data="requestsData"
        />
      </div>

      <div v-if="appStore.sectionVisible.tokens" class="section">
        <BarChart
          chart-id="tokens-chart"
          title="Tokens 用量构成"
          x-axis-key="date"
          :data="dataStore.data.daily_data"
          :series="tokenSeries"
          stack
        />
      </div>

      <div v-if="appStore.sectionVisible.cacheRate" class="section">
        <LineChart
          chart-id="cache-chart"
          title="缓存命中率"
          :data="cacheData"
        />
      </div>

      <div v-if="appStore.sectionVisible.composition" class="section">
        <BarChart
          chart-id="composition-chart"
          title="Token 构成（月度）"
          x-axis-key="name"
          :data="compositionData"
          :series="[{ key: 'value', name: 'Tokens' }]"
        />
      </div>

      <div v-if="appStore.sectionVisible.models" class="section">
        <PieChart
          chart-id="models-chart"
          title="模型分布"
          :data="modelsData"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useAppStore } from '@/stores/useAppStore'
import { useDataStore } from '@/stores/useDataStore'
import { api } from '@/api/backend'
import Toolbar from './Toolbar.vue'
import SummaryCards from './SummaryCards.vue'
import LineChart from './LineChart.vue'
import BarChart from './BarChart.vue'
import PieChart from './PieChart.vue'

const appStore = useAppStore()
const dataStore = useDataStore()

onMounted(() => {
  appStore.init()
  refresh()
})

const refresh = async () => {
  try {
    await dataStore.fetchDashboard(dataStore.selectedPeriod)
  } catch (error) {
    console.error('Failed to refresh:', error)
  }
}

// 计算图表数据
const requestsData = computed(() => {
  if (!dataStore.data) return []
  return dataStore.data.daily_data.map((item) => ({
    date: item.date,
    value: item.requests,
  }))
})

const tokenSeries = [
  { key: 'completion', name: '输出' },
  { key: 'cache_miss', name: '未缓存' },
  { key: 'cache_hit', name: '缓存命中' },
]

const cacheData = computed(() => {
  if (!dataStore.data) return []
  return dataStore.data.daily_data.map((item) => ({
    date: item.date,
    value: item.cache_hit + item.cache_miss > 0
      ? (item.cache_hit / (item.cache_hit + item.cache_miss)) * 100
      : 0,
  }))
})

const compositionData = computed(() => {
  if (!dataStore.data) return []
  const total = dataStore.data.cache_metrics.total_cache_hit +
                dataStore.data.cache_metrics.total_cache_miss
  return [
    { name: '缓存命中', value: dataStore.data.cache_metrics.total_cache_hit },
    { name: '未缓存', value: dataStore.data.cache_metrics.total_cache_miss },
  ]
})

const modelsData = computed(() => {
  if (!dataStore.data) return []
  return dataStore.data.model_breakdown.slice(0, 8).map((item) => ({
    name: item.model,
    value: item.tokens,
  }))
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.loading,
.error {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  text-align: center;
  color: #666;
}

:dark .loading,
:dark .error {
  background-color: #2a2a2a;
  color: #999;
}

.error {
  color: #d32f2f;
}

.error .btn {
  margin-top: 15px;
}

.section {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

