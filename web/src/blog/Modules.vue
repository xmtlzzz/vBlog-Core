<template>
  <BlogNav />
  <div class="page-enter">
  <header class="page-header">
    <h1>功能模块 Modules</h1>
    <p>通过组件定制上传的功能模块，每个模块独立运行。</p>
  </header>

  <main class="modules-grid" v-if="modules.length">
    <div v-for="mod in modules" :key="mod.id" class="module-card fade-in">
      <div class="module-card-top">
        <div class="module-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="3" width="18" height="18" rx="2"/>
            <path d="M3 9h18"/>
            <path d="M9 21V9"/>
          </svg>
        </div>
        <span class="module-status">
          <span class="status-dot"></span>
          Running
        </span>
      </div>
      <div class="module-name">{{ mod.name }}</div>
      <div class="module-desc">{{ mod.description || '自定义功能模块' }}</div>
      <div class="module-preview">
        <iframe
          :srcdoc="srcdoc(mod.code)"
          class="module-iframe"
          sandbox="allow-scripts"
          :title="mod.name"
        ></iframe>
      </div>
      <div class="module-footer">
        <span class="module-meta">v{{ mod.version || '1.0' }}</span>
        <span class="module-meta">{{ mod.origin || 'custom' }}</span>
      </div>
    </div>
  </main>

  <div v-else class="empty-state fade-in">
    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="color:var(--muted)">
      <rect x="3" y="3" width="18" height="18" rx="2"/>
      <path d="M3 9h18"/>
      <path d="M9 21V9"/>
    </svg>
    <p>暂无功能模块</p>
    <span>前往后台 → 组件定制上传模块</span>
  </div>

  </div>
  <BlogFooter />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import { buildSrcdoc } from '../utils/component'

const modules = ref([])
const blogData = ref(null)

function srcdoc(code) {
  return buildSrcdoc(code, { data: blogData.value })
}

onMounted(async () => {
  try {
    const [compRes, statsRes] = await Promise.all([
      api.get('/components/active'),
      api.get('/dashboard/stats').catch(() => null)
    ])
    modules.value = Array.isArray(compRes) ? compRes : (compRes.data || [])
    if (statsRes) {
      var s = statsRes.data || statsRes
      blogData.value = {
        posts: s.total_posts || 0,
        views: s.total_views || 0,
        comments: s.total_comments || 0,
        tags: s.total_tags || 0,
        visitors: Math.floor(Math.random() * 9000) + 1000
      }
    }
  } catch {}
})
</script>

<style scoped>
.page-header {
  max-width: 1080px;
  margin: 0 auto;
  padding: 56px 24px 32px;
}
.page-header h1 {
  font-family: var(--font-display);
  font-size: clamp(28px, 4vw, 40px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin-bottom: 8px;
}
.page-header p {
  font-size: 16px;
  color: var(--muted);
}

.modules-grid {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 80px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 16px;
}

.module-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  transition: border-color 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
}
.module-card:hover {
  border-color: var(--accent);
  box-shadow: 0 4px 16px rgba(0,0,0,0.06);
}

.module-card-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.module-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--accent-soft);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
}
.module-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--success);
  font-weight: 500;
}
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--success);
  animation: pulse 2s infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.module-name {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 4px;
}
.module-desc {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 16px;
  line-height: 1.5;
}

.module-preview {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px;
  margin-bottom: 16px;
  flex: 1;
  min-height: 80px;
}
.module-iframe {
  width: 100%;
  min-height: 60px;
  border: none;
  background: transparent;
}

.module-footer {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--muted);
  font-family: var(--font-mono);
}

.empty-state {
  max-width: 1080px;
  margin: 0 auto;
  padding: 80px 24px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.empty-state p {
  font-size: 16px;
  color: var(--fg);
  font-weight: 500;
}
.empty-state span {
  font-size: 13px;
  color: var(--muted);
}

@media (max-width: 640px) {
  .page-header {
    padding: 40px 16px 24px;
  }
  .modules-grid {
    padding: 0 16px 48px;
    grid-template-columns: 1fr;
  }
}
</style>
