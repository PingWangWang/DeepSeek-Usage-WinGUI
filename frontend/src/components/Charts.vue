<template>
  <div class="charts">
    <div class="chart-placeholder">
      <p>📊 图表组件（ECharts）</p>
      <p>在此处渲染 7 张交互式图表</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardData } from '../stores/useDataStore'

defineProps<{
  data: DashboardData
}>()
</script>

<style scoped>
.charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 15px;
}

.chart-placeholder {
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
  color: #999;
}

p {
  margin: 10px 0;
}
</style>
