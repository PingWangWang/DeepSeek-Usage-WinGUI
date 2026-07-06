<template>
  <header class="header">
    <div class="header-content">
      <h1>DeepSeek Usage+</h1>
      <div class="header-actions">
        <button
          @click="currentView = 'dashboard'"
          :class="['nav-btn', { active: currentView === 'dashboard' }]"
        >
          📊 仪表盘
        </button>
        <button
          @click="currentView = 'settings'"
          :class="['nav-btn', { active: currentView === 'settings' }]"
        >
          ⚙️ 设置
        </button>
        <button @click="toggleTheme" class="btn-icon">{{ appStore.isDark ? '☀️' : '🌙' }}</button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAppStore } from '../stores/useAppStore'

const appStore = useAppStore()
const currentView = ref('dashboard')

defineExpose({ currentView })

const toggleTheme = () => {
  appStore.theme = appStore.theme === 'light' ? 'dark' : 'light'
  document.documentElement.style.colorScheme = appStore.theme
}
</script>

<style scoped>
.header {
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
  padding: 15px 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

:dark .header {
  background-color: #1a1a1a;
  border-bottom-color: #333;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
}

h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

:dark h1 {
  color: #e0e0e0;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.nav-btn {
  padding: 8px 16px;
  background-color: transparent;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
}

:dark .nav-btn {
  border-color: #444;
  color: #999;
}

.nav-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.nav-btn.active {
  background-color: #1890ff;
  border-color: #1890ff;
  color: white;
}

.btn-icon {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.btn-icon:hover {
  background-color: #f0f0f0;
}

:dark .btn-icon:hover {
  background-color: #2a2a2a;
}
</style>

