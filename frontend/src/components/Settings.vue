<template>
  <div class="settings-panel">
    <h2>⚙️ 设置</h2>

    <!-- Token 设置 -->
    <div class="settings-section">
      <h3>API Token</h3>
      <p class="description">输入你的 DeepSeek API Token，用于获取用量数据</p>
      <div class="form-group">
        <label>API Token</label>
        <div class="token-input">
          <input
            v-model="token"
            :type="showToken ? 'text' : 'password'"
            placeholder="sk-..."
            @keyup.enter="saveToken"
          />
          <button @click="toggleShowToken" class="btn-icon">{{ showToken ? '👁️' : '👁️‍🗨️' }}</button>
        </div>
        <small>Token 安全地存储在本地，不会被上传</small>
      </div>
      <button @click="saveToken" class="btn" :disabled="!token">保存 Token</button>
    </div>

    <!-- 自动刷新设置 -->
    <div class="settings-section">
      <h3>自动刷新</h3>
      <div class="form-group">
        <label>刷新间隔</label>
        <select v-model="autoRefreshInterval" @change="updateAutoRefresh">
          <option :value="0">关闭</option>
          <option :value="60000">1 分钟</option>
          <option :value="300000">5 分钟</option>
          <option :value="600000">10 分钟</option>
          <option :value="1800000">30 分钟</option>
          <option :value="3600000">1 小时</option>
        </select>
      </div>
    </div>

    <!-- 主题设置 -->
    <div class="settings-section">
      <h3>主题</h3>
      <div class="form-group">
        <label>外观模式</label>
        <select v-model="theme" @change="updateTheme">
          <option value="light">浅色模式</option>
          <option value="dark">深色模式</option>
        </select>
      </div>
    </div>

    <!-- 版本信息 -->
    <div class="settings-section">
      <h3>关于</h3>
      <div class="about-info">
        <p><strong>应用版本</strong>：{{ appVersion }}</p>
        <p>
          <strong>最新版本</strong>：{{ latestVersion }}
          <span v-if="hasUpdate" class="badge badge-update">有更新</span>
        </p>
        <button v-if="hasUpdate" @click="checkUpdate" class="btn">检查更新</button>
        <button v-else @click="checkUpdate" class="btn">检查更新</button>
      </div>
    </div>

    <!-- 提示信息 -->
    <div v-if="message" :class="['message', `message-${message.type}`]">
      {{ message.text }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAppStore } from '@/stores/useAppStore'
import { useDataStore } from '@/stores/useDataStore'
import { api } from '@/api/backend'

const appStore = useAppStore()

// 用computed动态绑定appStore的值，保持响应式
const token = computed({
  get: () => appStore.token,
  set: (val) => {
    appStore.token = val
  }
})

const autoRefreshInterval = computed({
  get: () => appStore.autoRefreshInterval || 0,
  set: (val) => {
    appStore.autoRefreshInterval = val
  }
})

const theme = computed({
  get: () => appStore.theme,
  set: (val) => {
    appStore.theme = val as 'light' | 'dark'
    document.documentElement.style.colorScheme = val
  }
})

const showToken = ref(false)
const appVersion = ref('1.0.0')
const latestVersion = ref('')
const hasUpdate = ref(false)
const message = ref<{ type: 'success' | 'error'; text: string } | null>(null)

onMounted(async () => {
  try {
    const version = await api.getVersion()
    appVersion.value = version
  } catch (error) {
    console.error('Failed to get version:', error)
  }

  checkUpdate()
})

const toggleShowToken = () => {
  showToken.value = !showToken.value
}

const saveToken = async () => {
  if (!appStore.token) {
    showMessage('error', '请输入 Token')
    return
  }

  try {
    await api.updateToken(appStore.token)
    // 主动更新localStorage（触发store的持久化）
    localStorage.setItem('app_token', appStore.token)
    showMessage('success', 'Token 保存成功')
    // 保存Token后刷新仪表盘数据
    const dataStore = useDataStore()
    await dataStore.fetchDashboard(dataStore.selectedPeriod)
  } catch (error) {
    showMessage('error', 'Token 保存失败')
    console.error('Failed to update token:', error)
  }
}

const updateAutoRefresh = async () => {
  try {
    await api.setAutoRefresh(appStore.autoRefreshInterval)
    showMessage('success', '自动刷新已更新')
  } catch (error) {
    showMessage('error', '更新自动刷新失败')
    console.error('Failed to update auto refresh:', error)
  }
}

const updateTheme = () => {
  document.documentElement.style.colorScheme = appStore.theme
  showMessage('success', `已切换为${appStore.theme === 'dark' ? '深色' : '浅色'}模式`)
}

const checkUpdate = async () => {
  try {
    const info = await api.checkUpdate()
    latestVersion.value = info.latest_version
    hasUpdate.value = info.has_update
    if (hasUpdate.value) {
      showMessage('success', `发现新版本：${info.latest_version}`)
    } else {
      showMessage('success', '已是最新版本')
    }
  } catch (error) {
    showMessage('error', '检查更新失败')
    console.error('Failed to check update:', error)
  }
}

const showMessage = (type: 'success' | 'error', text: string) => {
  message.value = { type, text }
  setTimeout(() => {
    message.value = null
  }, 3000)
}
</script>

<style scoped>
.settings-panel {
  background-color: #fff;
  border-radius: 8px;
  padding: 30px;
  max-width: 600px;
  margin: 0 auto;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

:dark .settings-panel {
  background-color: #2a2a2a;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

h2 {
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
}

:dark h2 {
  color: #e0e0e0;
}

.settings-section {
  margin-bottom: 30px;
  padding-bottom: 30px;
  border-bottom: 1px solid #ddd;
}

:dark .settings-section {
  border-bottom-color: #444;
}

.settings-section:last-child {
  border-bottom: none;
}

h3 {
  margin-bottom: 15px;
  font-size: 16px;
  color: #333;
}

:dark h3 {
  color: #e0e0e0;
}

.description {
  margin-bottom: 15px;
  color: #999;
  font-size: 14px;
}

:dark .description {
  color: #666;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #666;
}

:dark label {
  color: #999;
}

input,
select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

:dark input,
:dark select {
  background-color: #1a1a1a;
  color: #e0e0e0;
  border-color: #444;
}

input:focus,
select:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.token-input {
  display: flex;
  gap: 8px;
}

.token-input input {
  flex: 1;
}

.btn-icon {
  padding: 8px 12px;
  border: 1px solid #ddd;
  background-color: #f5f5f5;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

:dark .btn-icon {
  background-color: #1a1a1a;
  border-color: #444;
}

small {
  display: block;
  margin-top: 8px;
  color: #999;
  font-size: 12px;
}

:dark small {
  color: #666;
}

.about-info {
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 15px;
}

:dark .about-info {
  background-color: #1a1a1a;
}

.about-info p {
  margin-bottom: 8px;
  font-size: 14px;
}

.badge {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  margin-left: 8px;
}

.badge-update {
  background-color: #ff7a45;
  color: white;
}

.message {
  padding: 12px;
  border-radius: 4px;
  margin-top: 20px;
  text-align: center;
  font-size: 14px;
}

.message-success {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
  color: #52c41a;
}

:dark .message-success {
  background-color: #162312;
  border-color: #274916;
}

.message-error {
  background-color: #fff1f0;
  border: 1px solid #ffccc7;
  color: #ff4d4f;
}

:dark .message-error {
  background-color: #2f1410;
  border-color: #58181c;
}
</style>
