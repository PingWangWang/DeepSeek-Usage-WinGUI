<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <select v-model="selectedPeriod" class="period-select">
        <option value="current_month">本月</option>
        <option value="last_month">上月</option>
        <option value="last_3_months">最近3个月</option>
        <option value="last_12_months">最近12个月</option>
      </select>
      <button @click="refresh" class="btn">🔄 刷新</button>
    </div>

    <div class="toolbar-right">
      <label v-for="(visible, key) in sectionVisible" :key="key" class="toggle-label">
        <input type="checkbox" v-model="sectionVisible[key]" />
        {{ formatLabel(key) }}
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '../stores/useAppStore'
import { useDataStore } from '../stores/useDataStore'

const appStore = useAppStore()
const dataStore = useDataStore()

const sectionVisible = computed(() => appStore.sectionVisible)
const selectedPeriod = computed({
  get: () => dataStore.selectedPeriod,
  set: (value) => dataStore.setPeriod(value),
})

const refresh = () => {
  dataStore.fetchDashboard(selectedPeriod.value)
}

const formatLabel = (key: string) => {
  const labels: Record<string, string> = {
    requests: '请求',
    tokens: 'Tokens',
    cacheRate: '缓存',
    composition: '构成',
    models: '模型',
  }
  return labels[key] || key
}
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  gap: 15px;
  align-items: center;
}

.period-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.btn {
  padding: 8px 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn:hover {
  background-color: #0050b3;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
}

.toggle-label input {
  cursor: pointer;
}
</style>
