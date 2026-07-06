import * as echarts from 'echarts'
import type { ECharts, EChartsOption } from 'echarts'

// ECharts 实例管理
class ChartManager {
  private charts: Map<string, ECharts> = new Map()

  // 获取或创建图表实例
  getChart(elementId: string): ECharts {
    if (!this.charts.has(elementId)) {
      const dom = document.getElementById(elementId)
      if (!dom) {
        throw new Error(`Element ${elementId} not found`)
      }
      const chart = echarts.init(dom, { useDirtyRect: true })
      this.charts.set(elementId, chart)

      // 监听窗口尺寸变化
      window.addEventListener('resize', () => this.resize(elementId))
    }
    return this.charts.get(elementId)!
  }

  // 设置图表选项
  setOption(elementId: string, option: EChartsOption) {
    const chart = this.getChart(elementId)
    chart.setOption(option)
  }

  // 调整大小
  resize(elementId: string) {
    const chart = this.charts.get(elementId)
    if (chart) {
      chart.resize()
    }
  }

  // 销毁图表
  dispose(elementId: string) {
    const chart = this.charts.get(elementId)
    if (chart) {
      chart.dispose()
      this.charts.delete(elementId)
    }
  }

  // 销毁所有图表
  disposeAll() {
    this.charts.forEach((chart) => {
      chart.dispose()
    })
    this.charts.clear()
  }
}

export const chartManager = new ChartManager()

// 主题配置
export const darkTheme = {
  color: [
    '#5470c6',
    '#ee6666',
    '#9ac900',
    '#fac858',
    '#ee6590',
    '#73c0de',
    '#3ba272',
    '#fc8452',
  ],
  backgroundColor: '#1a1a1a',
  textStyle: {
    color: '#e0e0e0',
  },
}

export const lightTheme = {
  color: [
    '#5470c6',
    '#ee6666',
    '#91cc75',
    '#fac858',
    '#ee6590',
    '#73c0de',
    '#3ba272',
    '#fc8452',
  ],
  backgroundColor: '#ffffff',
  textStyle: {
    color: '#333',
  },
}

// 图表颜色辅助函数
export function getChartColor(isDark: boolean, index: number) {
  const theme = isDark ? darkTheme : lightTheme
  return theme.color[index % theme.color.length]
}
