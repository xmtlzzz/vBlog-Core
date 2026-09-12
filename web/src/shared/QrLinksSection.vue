<template>
  <section v-if="projects.length" class="qr-links fade-in">
    <div class="qr-links-grid">
      <div v-for="p in projects" :key="p.url" class="qr-link-card">
        <span class="qr-link-badge" :title="p.url + ' · 点击变形为二维码'">
          <EveryQrBadge :url="p.url" :model="model" :effect="effect" :palette="palette" />
        </span>
        <a class="qr-link-name" :href="p.url" target="_blank" rel="noopener">{{ p.name }}</a>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import api from '../api/request'
import { parseQrLinks } from '../utils/qrLinks'

const EveryQrBadge = defineAsyncComponent(() => import('./EveryQrBadge.vue'))

const settings = ref({})

const projects = computed(() => parseQrLinks(settings.value.qr_links))

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

onMounted(async () => {
  try {
    settings.value = await api.get('/settings')
  } catch {}
})
</script>

<style scoped>
/* 作为「功能模块」页内容的一部分直接呈现，无独立标题 */
.qr-links {
  margin-bottom: 8px;
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
  text-align: center;
  word-break: break-all;
}
.qr-link-name:hover {
  color: var(--accent);
}
</style>
