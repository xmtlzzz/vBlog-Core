<template>
  <BlogNav />
  <main class="not-found-page page-enter">
    <div class="not-found-container">
      <div class="not-found-code">404</div>
      <h1 class="not-found-title">页面未找到</h1>
      <p class="not-found-desc">您访问的页面不存在，或者已经被移动或删除。</p>
      <div class="not-found-actions">
        <router-link to="/" class="btn-primary">← 返回首页</router-link>
        <router-link to="/archives" class="btn-secondary">浏览归档</router-link>
      </div>
    </div>
  </main>
  <CustomWidgets />
  <BlogFooter />
</template>

<script setup>
import { onMounted } from 'vue'
import { updateMetadata } from '../utils/metadata'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'

onMounted(() => {
  updateMetadata({}, { title: '404 页面未找到', description: '您访问的页面不存在，或者已经被移动或删除。' })
})
</script>

<style scoped>
.not-found-page {
  min-height: calc(100vh - 160px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem 1.5rem;
  box-sizing: border-box;
}

.not-found-container {
  text-align: center;
  max-width: 480px;
  margin: 0 auto;
}

.not-found-code {
  font-size: clamp(4rem, 12vw, 6rem);
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.04em;
  color: var(--accent);
  opacity: 0.9;
  margin-bottom: 1rem;
}

.not-found-title {
  font-size: 1.6rem;
  font-weight: 700;
  margin: 0 0 0.75rem;
  color: var(--fg);
}

.not-found-desc {
  font-size: 0.95rem;
  line-height: 1.6;
  color: var(--muted);
  margin: 0 0 2rem;
}

.not-found-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 1.4rem;
  border-radius: var(--radius, 8px);
  background: var(--accent);
  color: #fff;
  font-size: 0.9rem;
  font-weight: 500;
  text-decoration: none;
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-secondary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.65rem 1.4rem;
  border-radius: var(--radius, 8px);
  background: var(--surface);
  border: 1px solid var(--border);
  color: var(--fg);
  font-size: 0.9rem;
  font-weight: 500;
  text-decoration: none;
  transition: background 0.2s ease, transform 0.2s ease;
}

.btn-secondary:hover {
  background: var(--card-hover);
  transform: translateY(-1px);
}
</style>
