<template>
  <div class="admin-layout">
    <aside class="sidebar" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand"><span class="dot"></span>vBlog <span class="badge">Admin</span></div>
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div class="nav-section-title">概览 Overview</div>
          <router-link to="/admin" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">◎</span> 仪表盘 Dashboard
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">内容 Content</div>
          <router-link to="/admin/posts" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">≡</span> 文章 Posts
          </router-link>
          <router-link to="/admin/tags" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">◉</span> 标签 Tags
          </router-link>
          <router-link to="/admin/comments" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">❝</span> 评论 Comments
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">系统 System</div>
          <router-link to="/admin/custom" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">⬡</span> 组件 Custom
          </router-link>
          <router-link to="/admin/settings" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">⚙</span> 设置 Settings
          </router-link>
        </div>
      </nav>
      <div class="sidebar-footer">
        <div class="sidebar-user">
          <div class="user-avatar">A</div>
          <div class="user-info">
            <div class="user-name">Admin</div>
            <div class="user-role">超级管理员</div>
          </div>
          <button class="logout-btn" @click="handleLogout" title="退出登录">⏻</button>
        </div>
      </div>
    </aside>

    <div class="main-area">
      <header class="topbar">
        <div class="topbar-left">
          <button class="menu-toggle" @click="sidebarOpen = !sidebarOpen">☰</button>
          <template v-for="(crumb, i) in breadcrumbs" :key="crumb.path">
            <router-link v-if="crumb.path" :to="crumb.path" class="breadcrumb-link">{{ crumb.label }}</router-link>
            <span v-if="crumb.path" class="breadcrumb-sep">/</span>
            <span v-if="i === breadcrumbs.length - 1" class="breadcrumb-current">{{ crumb.label }}</span>
          </template>
        </div>
        <div class="topbar-right">
          <router-link to="/" class="view-blog-link">← 查看博客 View Blog</router-link>
          <button class="theme-toggle" @click="themeStore.toggle()">
            {{ themeStore.theme === 'dark' ? '☀️' : '🌙' }}
          </button>
        </div>
      </header>
      <main class="content">
        <router-view />
      </main>
    </div>

    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useThemeStore } from '../stores/theme'
import { useAuthStore } from '../stores/auth'
const themeStore = useThemeStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(false)

function handleLogout() {
  authStore.logout()
  router.push('/admin/login')
}

const pageMap = {
  '/admin': '仪表盘 Dashboard',
  '/admin/posts': '文章 Posts',
  '/admin/tags': '标签 Tags',
  '/admin/comments': '评论 Comments',
  '/admin/custom': '组件 Custom',
  '/admin/settings': '设置 Settings'
}

const breadcrumbs = computed(() => {
  const crumbs = []
  // Find parent page
  const basePath = '/' + route.path.split('/').slice(1, 3).join('/')
  const parentLabel = pageMap[basePath]
  if (parentLabel && basePath !== '/admin') {
    crumbs.push({ label: parentLabel, path: basePath })
  }
  // Current page (use meta.breadcrumb if available, otherwise pageMap)
  const currentLabel = route.meta?.breadcrumb || pageMap[route.path]
  if (currentLabel) {
    crumbs.push({ label: currentLabel, path: '' })
  }
  // Fallback: dashboard
  if (crumbs.length === 0) {
    crumbs.push({ label: '仪表盘 Dashboard', path: '' })
  }
  return crumbs
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg);
  color: var(--fg);
}
.sidebar {
  width: 220px;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
  transition: transform 0.3s;
}
.sidebar-brand {
  padding: 20px 20px 16px;
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--fg);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 8px;
}
.sidebar-brand .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent);
}
.sidebar-brand .badge {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: var(--radius);
  background: var(--accent-soft);
  color: var(--accent);
  margin-left: auto;
}
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}
.nav-section {
  margin-bottom: 8px;
}
.nav-section-title {
  padding: 8px 20px 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--muted);
}
.sidebar-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  font-size: 13px;
  color: var(--muted);
  text-decoration: none;
  transition: all 0.15s;
  border-radius: var(--radius);
}
.sidebar-link:hover {
  color: var(--fg);
  background: var(--accent-soft);
}
.sidebar-link.router-link-active {
  color: var(--fg);
  background: var(--surface);
  box-shadow: 0 1px 2px rgba(0,0,0,0.06);
  font-weight: 500;
}
.nav-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}
.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border);
}
.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
}
.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}
.user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--fg);
}
.user-role {
  font-size: 11px;
  color: var(--muted);
}
.logout-btn {
  margin-left: auto;
  background: none;
  border: none;
  color: var(--muted);
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.15s;
}
.logout-btn:hover {
  color: var(--error);
  background: var(--error-soft);
}
.main-area {
  flex: 1;
  margin-left: 220px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
.topbar {
  height: 56px;
  border-bottom: 1px solid var(--border);
  background: var(--nav-bg);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
}
.topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.menu-toggle {
  display: none;
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--fg);
  padding: 4px;
}
.view-blog-link {
  font-size: 13px;
  color: var(--muted);
  text-decoration: none;
  padding: 4px 10px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  transition: all 0.15s;
}
.view-blog-link:hover {
  color: var(--fg);
  border-color: var(--fg);
}
.breadcrumb-link {
  font-size: 14px;
  color: var(--muted);
  text-decoration: none;
  transition: color 0.15s;
}
.breadcrumb-link:hover {
  color: var(--accent);
}
.breadcrumb-sep {
  color: var(--border);
  font-size: 14px;
  margin: 0 4px;
}
.breadcrumb-current {
  font-size: 14px;
  font-weight: 500;
  color: var(--fg);
}
.topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.theme-toggle {
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 6px 10px;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.15s;
}
.theme-toggle:hover {
  border-color: var(--fg);
}
/* new-post-btn removed — not in topbar */
.content {
  flex: 1;
  padding: 24px;
  overflow-x: hidden;
}
.sidebar-overlay {
  display: none;
}

@media (max-width: 900px) {
  .sidebar {
    transform: translateX(-100%);
  }
  .sidebar.open {
    transform: translateX(0);
  }
  .main-area {
    margin-left: 0;
  }
  .menu-toggle {
    display: block;
  }
  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.3);
    z-index: 99;
  }
}
</style>
