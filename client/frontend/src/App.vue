<template>
  <div class="app">
    <header class="app-header">
      <h1>vBlog Monitor</h1>
      <button @click="showSettings = true">Settings</button>
    </header>
    <StatsBar :stats="stats" />
    <div class="change-feed">
      <ChangeCard v-for="c in changes" :key="c.id" :change="c" />
      <div v-if="!changes.length" class="empty">No changes yet</div>
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

async function handleConnect({ addr, apiKey }) {
  try {
    await window.go.main.App.Connect(addr, apiKey)
    showSettings.value = false
    await refresh()
    await window.go.main.App.WatchChanges(apiKey, 0)
    refreshInterval = setInterval(refresh, 30000)
  } catch (e) {
    alert('Connection failed: ' + e)
  }
}

async function refresh() {
  try {
    stats.value = await window.go.main.App.GetLatestStats()
    await loadTrends(granularity.value)
  } catch {}
}

async function loadTrends(g) {
  granularity.value = g
  const resp = await window.go.main.App.GetTrends(g, 14)
  trends.value = resp?.points || []
}

onMounted(() => {
  if (window.runtime) {
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
.change-feed { display: flex; flex-direction: column; gap: 8px; margin: 16px 0; }
.empty { text-align: center; color: #737373; padding: 24px; }
</style>
