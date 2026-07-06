<template>
  <footer class="footer">
    <p>DeepSeek Usage+ {{ version }} | © 2024 PingWangWang</p>
  </footer>
</template>

<script setup lang="ts">
import { useAppStore } from '../stores/useAppStore'

const appStore = useAppStore()
const version = appStore.version
</script>

<style scoped>
.footer {
  background-color: #f5f5f5;
  border-top: 1px solid #e0e0e0;
  padding: 15px;
  text-align: center;
  font-size: 12px;
  color: #999;
}

p {
  margin: 0;
}
</style>
