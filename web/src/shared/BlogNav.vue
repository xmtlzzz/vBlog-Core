<template>
  <nav class="top">
    <div class="inner">
      <router-link class="nav-brand" to="/">
        <span class="dot"></span>vBlog
      </router-link>
      <div class="nav-links">
        <router-link to="/">首页</router-link>
        <router-link to="/archives">归档</router-link>
        <router-link to="/tags">标签</router-link>
        <router-link to="/about">关于</router-link>
      </div>
      <div class="nav-right">
        <router-link to="/admin" class="admin-btn">后台</router-link>
        <el-button circle size="small" class="theme-toggle" @click="themeStore.toggle()" aria-label="切换主题">
          <span v-if="themeStore.theme === 'dark'">☀</span>
          <span v-else>☾</span>
        </el-button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useThemeStore } from '../stores/theme'
const themeStore = useThemeStore()
</script>

<style scoped>
.top {
  position: sticky;
  top: 0;
  z-index: 100;
  background: rgba(250, 250, 250, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
  transition: background 0.3s ease, border-color 0.3s ease;
}
[data-theme="dark"] .top {
  background: rgba(10, 10, 10, 0.85);
}
.inner {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.nav-brand {
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--fg);
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 8px;
}
.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
}
.nav-links a {
  font-size: 14px;
  color: var(--muted);
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
}
.nav-links a:hover,
.nav-links :deep(.router-link-active) {
  color: var(--fg);
  background: rgba(0, 0, 0, 0.04);
}
[data-theme="dark"] .nav-links a:hover,
[data-theme="dark"] .nav-links :deep(.router-link-active) {
  background: rgba(255, 255, 255, 0.06);
}
.nav-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.admin-btn {
  font-size: 13px;
  color: var(--muted);
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid var(--border);
  transition: all 0.15s;
}
.admin-btn:hover {
  color: var(--fg);
  border-color: var(--fg);
}
.theme-toggle {
  width: 36px;
  height: 36px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--fg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  font-size: 16px;
}
.theme-toggle:hover {
  background: rgba(0, 0, 0, 0.04);
}
[data-theme="dark"] .theme-toggle:hover {
  background: rgba(255, 255, 255, 0.06);
}

@media (max-width: 640px) {
  .nav-links {
    gap: 0;
  }
  .nav-links a {
    padding: 6px 8px;
    font-size: 13px;
  }
  .admin-btn {
    padding: 4px 8px;
    font-size: 12px;
  }
}
</style>
