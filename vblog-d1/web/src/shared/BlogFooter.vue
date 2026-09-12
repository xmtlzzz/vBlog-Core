<template>
  <footer ref="footerRef" class="blog-footer">
    <span>© 2026 vBlog Core · 用代码写作，用文字思考</span>
    <div class="footer-links">
      <template v-if="settings.author_github">
        <span
          v-if="badgeEnabled && badgeVisible"
          class="qr-badge"
          :class="{ zoomed: qrZoomed }"
          :style="{ background: qrBackground }"
          title="GitHub · 点击变形为二维码"
        >
          <EveryQrBadge
            :url="settings.author_github"
            :model="qrModel"
            :effect="qrEffect"
            @viewchange="qrZoomed = $event"
          />
        </span>
        <a :href="settings.author_github" target="_blank" rel="noopener">GitHub</a>
      </template>
      <a href="/api/rss" target="_blank">RSS</a>
    </div>
  </footer>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted, onUnmounted } from 'vue'
import api from '../api/request'

const EveryQrBadge = defineAsyncComponent(() => import('./EveryQrBadge.vue'))

const settings = ref({})
const footerRef = ref(null)
const badgeVisible = ref(false)
const qrZoomed = ref(false)
let observer = null

const badgeEnabled = computed(() => settings.value.qr_badge !== 'false')
const qrModel = computed(() => (settings.value.qr_model === 'terrain' ? 'terrain' : 'tree'))
const qrEffect = computed(() => {
  const v = settings.value.qr_effect
  return ['calm', 'snow', 'rain', 'wind'].includes(v) ? v : ''
})
const qrBackground = computed(() =>
  /^#[0-9a-fA-F]{3,8}$/.test(settings.value.qr_bg || '') ? settings.value.qr_bg : '#ffffff'
)

onMounted(async () => {
  try {
    settings.value = await api.get('/settings')
  } catch {}
  // 徽章组件含 3D 渲染器，页脚进入视口时才动态加载
  observer = new IntersectionObserver((entries) => {
    if (entries.some((e) => e.isIntersecting)) {
      observer?.disconnect()
      observer = null
      badgeVisible.value = true
    }
  })
  if (footerRef.value) observer.observe(footerRef.value)
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
.blog-footer {
  max-width: 1080px;
  margin: 0 auto;
  padding: 24px 24px 48px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--muted);
}
.blog-footer a {
  color: var(--muted);
  text-decoration: none;
  transition: color 0.15s;
}
.blog-footer a:hover {
  color: var(--fg);
}
.footer-links {
  display: flex;
  gap: 16px;
  align-items: center;
}
.qr-badge {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border);
  cursor: pointer;
  transition: width 0.25s ease, height 0.25s ease, border-color 0.15s;
  z-index: 5;
}
.qr-badge:hover {
  border-color: var(--fg);
}
/* 二维码视图放大，保证可扫描；回到模型视图时恢复 */
.qr-badge.zoomed {
  width: 104px;
  height: 104px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

@media (max-width: 640px) {
  .blog-footer {
    padding: 24px 16px 32px;
    flex-direction: column;
    gap: 12px;
  }
}
</style>
