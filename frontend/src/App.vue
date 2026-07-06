<template>
  <div class="app" :style="{ colorScheme: appStore.theme }">
    <Header ref="headerRef" />
    <main class="main-content">
      <Dashboard v-if="headerRef?.currentView === 'dashboard'" />
      <Settings v-else-if="headerRef?.currentView === 'settings'" />
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

onMounted(() => {
  appStore.init()
  // 设置系统主题
  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
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
</style>

