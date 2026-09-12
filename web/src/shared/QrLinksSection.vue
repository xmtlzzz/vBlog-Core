<template>
  <section v-if="projects.length" class="qr-links fade-in">
    <div class="qr-links-toolbar" role="group" aria-label="氛围特效">
      <button
        v-for="opt in EFFECTS"
        :key="opt.value"
        type="button"
        class="qr-effect-btn"
        :class="{ active: effect === opt.value }"
        @click="activeEffect = opt.value"
      >{{ opt.label }}</button>
    </div>
    <div class="qr-links-grid">
      <div v-for="(p, i) in projects" :key="p.url" class="qr-link-card">
        <span class="qr-link-badge" :title="p.url + ' · 点击变形为二维码'">
          <EveryQrBadge :url="p.url" :model="model" :effect="effect" :palette="paletteOf(p, i)" />
        </span>
        <a class="qr-link-name" :href="p.url" target="_blank" rel="noopener">{{ p.name }}</a>
        <a class="qr-link-host" :href="p.url" target="_blank" rel="noopener">{{ hostOf(p.url) }}</a>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import api from '../api/request'
import { parseQrLinks, QR_PALETTE_NAMES } from '../utils/qrLinks'

const EveryQrBadge = defineAsyncComponent(() => import('./EveryQrBadge.vue'))

const settings = ref({})

const projects = computed(() => parseQrLinks(settings.value.qr_links))

// 与后台 QR 徽章共用同一组外观配置
const model = computed(() => (settings.value.qr_model === 'terrain' ? 'terrain' : 'tree'))

const EFFECTS = [
  { value: '', label: '默认' },
  { value: 'calm', label: '平静' },
  { value: 'wind', label: '刮风' },
  { value: 'rain', label: '下雨' },
  { value: 'snow', label: '落雪' }
]

// null = 跟随后台配置；点击按钮后仅本地切换，不写回设置
const configuredEffect = computed(() => {
  const v = settings.value.qr_effect
  return ['calm', 'snow', 'rain', 'wind'].includes(v) ? v : ''
})
const activeEffect = ref(null)
const effect = computed(() => (activeEffect.value === null ? configuredEffect.value : activeEffect.value))

// 配色分发：行内第三段手动指定优先，否则按行轮换预设（相邻卡片不同色）
function paletteOf(p, index) {
  return p.palette || QR_PALETTE_NAMES[index % QR_PALETTE_NAMES.length]
}

function hostOf(url) {
  try {
    return new URL(url).host.replace(/^www\./, '')
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
/* 与模块页内容同容器（1080px 居中），作为页面内容直接呈现 */
.qr-links {
  max-width: 1080px;
  margin: 0 auto 8px;
  padding: 0 24px;
}
.qr-links-toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}
.qr-effect-btn {
  font-size: 12px;
  color: var(--muted);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 4px 12px;
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s;
}
.qr-effect-btn:hover {
  color: var(--fg);
  border-color: var(--fg);
}
.qr-effect-btn.active {
  color: var(--accent);
  border-color: var(--accent);
  background: var(--accent-soft);
}
.qr-links-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
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
/* 不设内框：徽章撑满卡片宽度，树可以长满整个画布 */
.qr-link-badge {
  display: block;
  width: 100%;
  max-width: 260px;
  aspect-ratio: 1;
  border-radius: 10px;
  overflow: hidden;
  background: #fff;
  cursor: pointer;
  margin: 0 auto 8px;
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
/* 可点击的域名胶囊：新标签直达项目 */
.qr-link-host {
  font-size: 12px;
  color: var(--muted);
  font-family: var(--font-mono);
  text-decoration: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 3px 10px;
  margin-top: 6px;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.15s, border-color 0.15s;
}
.qr-link-host:hover {
  color: var(--accent);
  border-color: var(--accent);
}
</style>
