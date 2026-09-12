<template>
  <section v-if="projects.length" class="qr-links fade-in">
    <h2 class="qr-links-title">项目导航 · 扫码访问</h2>
    <div class="qr-links-grid">
      <div v-for="p in projects" :key="p.url" class="qr-link-card">
        <span class="qr-link-badge" :title="p.url + ' · 点击变形为二维码'">
          <EveryQrBadge :url="p.url" :model="model" :effect="effect" :palette="palette" />
        </span>
        <a class="qr-link-name" :href="p.url" target="_blank" rel="noopener">{{ p.name }}</a>
        <span class="qr-link-host">{{ hostOf(p.url) }}</span>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import api from '../api/request'

const EveryQrBadge = defineAsyncComponent(() => import('./EveryQrBadge.vue'))

const settings = ref({})

const projects = computed(() => {
  const raw = (settings.value.qr_links || '').trim()
  if (!raw) return []
  const list = raw
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [name, url] = line.split('|').map((s) => s.trim())
      return { name: name || url || '', url: url || name || '' }
    })
    .filter((p) => /^https?:\/\//.test(p.url))
  return list
})

// 与后台 QR 徽章共用同一组外观配置
const model = computed(() => (settings.value.qr_model === 'terrain' ? 'terrain' : 'tree'))
const effect = computed(() => {
  const v = settings.value.qr_effect
  return ['calm', 'snow', 'rain', 'wind'].includes(v) ? v : ''
})
const palette = computed(() => {
  const v = settings.value.qr_palette
  return ['sakura', 'forest', 'ocean', 'sunset'].includes(v) ? v : 'sakura'
})

function hostOf(url) {
  try {
    return new URL(url).host
  } catch {
    return url
  }
}

onMounted(async () => {
  try {
    settings.value = await api.get('/settings')
  } catch {}
})
</script>

<style scoped>
.qr-links {
  margin-bottom: 28px;
}
.qr-links-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 16px;
  letter-spacing: -0.01em;
}
.qr-links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 14px;
}
.qr-link-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  transition: border-color 0.15s;
}
.qr-link-card:hover {
  border-color: var(--accent);
}
.qr-link-badge {
  display: block;
  width: 112px;
  height: 112px;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid var(--border);
  background: #fff;
  cursor: pointer;
  margin-bottom: 6px;
}
.qr-link-badge every-qr-code,
.qr-link-badge canvas,
.qr-link-badge button {
  display: block;
  width: 100%;
  height: 100%;
}
.qr-link-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--fg);
  text-decoration: none;
}
.qr-link-name:hover {
  color: var(--accent);
}
.qr-link-host {
  font-size: 12px;
  color: var(--muted);
  font-family: var(--font-mono);
}
</style>
