<template>
  <div class="chart-container">
    <div :id="chartId" style="width: 100%; height: 400px"></div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
import type { EChartsOption } from 'echarts'
import { chartManager, getChartColor } from '@/utils/charts'
import { useAppStore } from '@/stores/useAppStore'

interface Props {
  chartId: string
  title: string
  data: Array<{ date: string; value: number }> | null
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
})

const appStore = useAppStore()

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  chartManager.dispose(props.chartId)
})

watch(() => props.data, () => {
  initChart()
})

watch(() => appStore.theme, () => {
  initChart()
})

// 初始化图表
const initChart = () => {
  if (!props.data || props.data.length === 0) {
    return
  }

  const dates = props.data.map((item) => item.date)
  const values = props.data.map((item) => item.value)

  const option: EChartsOption = {
    title: {
      text: props.title,
      left: 'center',
      textStyle: {
        color: appStore.isDark ? '#e0e0e0' : '#333',
      },
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: appStore.isDark ? 'rgba(0,0,0,0.7)' : 'rgba(255,255,255,0.9)',
      textStyle: {
        color: appStore.isDark ? '#e0e0e0' : '#333',
      },
    },
    grid: {
      left: '10%',
      right: '10%',
      bottom: '10%',
      top: '15%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        color: appStore.isDark ? '#999' : '#666',
      },
      axisLine: {
        lineStyle: {
          color: appStore.isDark ? '#333' : '#ddd',
        },
      },
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: appStore.isDark ? '#999' : '#666',
      },
      splitLine: {
        lineStyle: {
          color: appStore.isDark ? '#2a2a2a' : '#f0f0f0',
        },
      },
    },
    series: [
      {
        data: values,
        type: 'line',
        smooth: true,
        itemStyle: {
          color: getChartColor(appStore.isDark, 0),
        },
        areaStyle: {
          color: getChartColor(appStore.isDark, 0) + '40',
        },
      },
    ],
  }

  chartManager.setOption(props.chartId, option)
}
</script>

<style scoped>
.chart-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-bottom: 15px;
}

:dark .chart-container {
  background-color: #2a2a2a;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}
</style>
