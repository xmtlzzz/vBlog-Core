<template>
  <div class="app">
    <header class="app-header">
      <h1>vBlog 监控</h1>
      <button @click="showSettings = true">设置</button>
    </header>
    <StatsBar :stats="stats" />
    <div class="change-feed">
      <div class="feed-title">最近变动</div>
      <ChangeCard v-for="c in changes.slice(0, 3)" :key="c.id" :change="c" />
      <div v-if="!changes.length" class="empty">暂无变动</div>
    </div>
    <TrendPanel :points="trends" :granularity="granularity" @change="loadTrends" />
    <Settings v-if="showSettings" @close="showSettings = false" @connect="handleConnect" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import StatsBar from './components/StatsBar.vue'
import ChangeCard from './components/ChangeCard.vue'
import TrendPanel from './components/TrendPanel.vue'
import Settings from './components/Settings.vue'

const stats = ref({})
const changes = ref([])
const trends = ref([])
const granularity = ref('day')
const showSettings = ref(false)
let refreshInterval = null

function getApp() {
  return window.go?.main?.App
}

async function handleConnect({ addr, apiKey }) {
  const app = getApp()
  if (!app) { alert('应用未就绪'); return }
  try {
    await app.Connect(addr, apiKey)
    showSettings.value = false
    await refresh()
    try {
      await app.WatchChanges(apiKey, 0)
    } catch (e) {
      console.error('WatchChanges error:', e)
    }
    refreshInterval = setInterval(refresh, 10000)
  } catch (e) {
    alert('连接失败: ' + e)
  }
}

async function refresh() {
  const app = getApp()
  if (!app) return
  try {
    const s = await app.GetLatestStats()
    if (s) stats.value = s
    await loadTrends(granularity.value)
  } catch (e) {
    console.error('refresh error:', e)
  }
}

async function loadTrends(g) {
  const app = getApp()
  if (!app) return
  granularity.value = g
  const resp = await app.GetTrends(g, 14)
  trends.value = resp?.points || []
}

onMounted(() => {
  if (window.runtime?.EventsOn) {
    window.runtime.EventsOn('change', (event) => {
      changes.value.unshift(event)
      if (changes.value.length > 50) changes.value.pop()
    })
  }
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>

<style>
* { margin: 0; box-sizing: border-box; }
body { font-family: system-ui, sans-serif; background: #fafafa; color: #171717; }
.app { max-width: 420px; margin: 0 auto; padding: 16px; }
.app-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.app-header h1 { font-size: 18px; }
.app-header button { padding: 6px 12px; border: 1px solid #e5e5e5; border-radius: 6px; background: white; cursor: pointer; }
.feed-title { font-size: 14px; font-weight: 600; margin-bottom: 8px; }
.change-feed { display: flex; flex-direction: column; gap: 8px; margin: 16px 0; }
.empty { text-align: center; color: #737373; padding: 24px; }
</style>
