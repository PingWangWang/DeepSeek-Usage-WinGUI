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

interface DataItem {
  name: string
  value: number
  percent?: number
}

interface Props {
  chartId: string
  title: string
  data: DataItem[] | null
  type?: 'pie' | 'doughnut'
}

const props = withDefaults(defineProps<Props>(), {
  data: () => [],
  type: 'pie',
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

const initChart = () => {
  if (!props.data || props.data.length === 0) {
    return
  }

  const option: EChartsOption = {
    title: {
      text: props.title,
      left: 'center',
      textStyle: {
        color: appStore.isDark ? '#e0e0e0' : '#333',
      },
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: appStore.isDark ? 'rgba(0,0,0,0.7)' : 'rgba(255,255,255,0.9)',
      textStyle: {
        color: appStore.isDark ? '#e0e0e0' : '#333',
      },
      formatter: '{b}: {c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      textStyle: {
        color: appStore.isDark ? '#e0e0e0' : '#333',
      },
    },
    series: [
      {
        name: props.title,
        type: props.type === 'doughnut' ? 'pie' : 'pie',
        radius: props.type === 'doughnut' ? ['40%', '70%'] : '70%',
        data: props.data.map((item, index) => ({
          value: item.value,
          name: item.name,
          itemStyle: {
            color: getChartColor(appStore.isDark, index),
          },
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)',
          },
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
