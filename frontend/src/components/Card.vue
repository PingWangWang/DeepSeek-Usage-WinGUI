<template>
  <div class="card">
    <div class="card-title">{{ title }}</div>
    <div class="card-value">{{ value }}</div>
    <div v-if="unit" class="card-unit">{{ unit }}</div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  value: string
  unit?: string
}>()
</script>

<style scoped>
.card {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.card-title {
  font-size: 14px;
  color: #999;
  margin-bottom: 10px;
}

.card-value {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.card-unit {
  font-size: 12px;
  color: #999;
}
</style>
