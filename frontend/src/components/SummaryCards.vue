<template>
  <div class="summary-cards">
    <Card title="当日费用" :value="formatCNY(metrics.day_cost)" />
    <Card title="月度费用" :value="formatCNY(metrics.month_cost)" />
    <Card title="均价" :value="formatPrice(metrics.avg_price)" />
    <Card title="总用量" :value="formatNumber(metrics.total_tokens)" unit="Token" />
    <Card title="缓存命中率" :value="formatPercent(metrics.cache_hit_rate)" />
  </div>
</template>

<script setup lang="ts">
import type { SummaryMetrics } from '../stores/useDataStore'
import Card from './Card.vue'

defineProps<{
  metrics: SummaryMetrics
}>()

const formatCNY = (value: number) => `¥${value.toFixed(2)}`
const formatPrice = (value: number) => `¥${value.toFixed(2)}/M`
const formatNumber = (value: number) => value.toLocaleString()
const formatPercent = (value: number) => `${(value * 100).toFixed(2)}%`
</script>

<style scoped>
.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
}
</style>
