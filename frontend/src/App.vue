<template>
  <div class="app" :style="{ colorScheme: appStore.theme }">
    <Header ref="headerRef" />
    <main class="main-content">
      <!-- 等待初始化完成 -->
      <div v-if="!appInitialized" class="loading">
        <p>📦 初始化中...</p>
      </div>
      <template v-else>
        <Dashboard v-if="headerRef?.currentView === 'dashboard'" />
        <Settings v-else-if="headerRef?.currentView === 'settings'" />
      </template>
    </main>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Header from './components/Header.vue'
import Dashboard from './components/Dashboard.vue'
import Settings from './components/Settings.vue'
import Footer from './components/Footer.vue'
import { useAppStore } from './stores/useAppStore'

const appStore = useAppStore()
const headerRef = ref<InstanceType<typeof Header> | null>(null)
const appInitialized = ref(false)

onMounted(async () => {
  await appStore.init()
  appInitialized.value = true

  // 仅在没有保存主题时才根据系统偏好设置主题
  const savedTheme = localStorage.getItem('app_theme')
  if (!savedTheme && window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    appStore.theme = 'dark'
  }
})
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
  transition: background-color 0.3s ease;
}

.app[style*="dark"] {
  background-color: #1a1a1a;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  font-size: 16px;
  color: #666;
}

:dark .loading {
  color: #999;
}
</style>

