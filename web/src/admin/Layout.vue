<template>
  <div class="admin-layout">
    <aside class="sidebar" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand"><span class="dot"></span>vBlog <span class="badge">Admin</span></div>
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div class="nav-section-title">概览 Overview</div>
          <router-link to="/admin" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
            </span>
            仪表盘 Dashboard
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">内容 Content</div>
          <router-link to="/admin/posts" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
            </span>
            文章 Posts
          </router-link>
          <router-link to="/admin/tags" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>
            </span>
            标签 Tags
          </router-link>
          <router-link to="/admin/comments" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>
            </span>
            评论 Comments
          </router-link>
          <router-link to="/admin/trash" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
            </span>
            回收站 Trash
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">系统 System</div>
          <router-link to="/admin/custom" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
            </span>
            组件 Custom
          </router-link>
          <router-link to="/admin/settings" class="sidebar-link" @click="sidebarOpen = false">
            <span class="nav-icon">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
            </span>
            设置 Settings
          </router-link>
        </div>
      </nav>
      <div class="sidebar-footer">
        <div class="sidebar-user">
          <div class="user-avatar">{{ (currentUser.username || 'A')[0].toUpperCase() }}</div>
          <div class="user-info">
            <div class="user-name">{{ currentUser.username }}</div>
            <div class="user-role">超级管理员</div>
          </div>
          <button class="logout-btn" @click="handleLogout" title="退出登录" aria-label="退出登录">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18.36 6.64a9 9 0 1 1-12.73 0"/><line x1="12" y1="2" x2="12" y2="12"/></svg>
          </button>
        </div>
      </div>
    </aside>

    <div class="main-area">
      <header class="topbar">
        <div class="topbar-left">
          <button class="menu-toggle" @click="sidebarOpen = !sidebarOpen" aria-label="切换侧栏">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          </button>
          <template v-for="(crumb, i) in breadcrumbs" :key="crumb.path">
            <router-link v-if="crumb.path" :to="crumb.path" class="breadcrumb-link">{{ crumb.label }}</router-link>
            <span v-if="crumb.path" class="breadcrumb-sep">/</span>
            <span v-if="i === breadcrumbs.length - 1" class="breadcrumb-current">{{ crumb.label }}</span>
          </template>
        </div>
        <div class="topbar-right">
          <router-link to="/" class="view-blog-link">← 查看博客 View Blog</router-link>
          <ThemeToggleAnimated
            :is-dark="themeStore.theme === 'dark'"
            size="sm"
            @toggle="themeStore.toggle()"
          />
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useThemeStore } from '../stores/theme'
import { useAuthStore } from '../stores/auth'
import ThemeToggleAnimated from '../shared/ThemeToggleAnimated.vue'
import api from '../api/request'

const themeStore = useThemeStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const sidebarOpen = ref(false)
const currentUser = ref({ username: 'Admin' })

onMounted(async () => {
  try {
    const me = await api.get('/auth/me')
    if (me?.username) {
      currentUser.value.username = me.username
      authStore.setUser(me)
    }
  } catch {
    // 401 is handled automatically by request interceptor
  }
})

function handleLogout() {
  authStore.logout()
  router.push('/admin/login')
}

const pageMap = {
  '/admin': '仪表盘 Dashboard',
  '/admin/posts': '文章 Posts',
  '/admin/tags': '标签 Tags',
  '/admin/comments': '评论 Comments',
  '/admin/trash': '回收站 Trash',
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
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
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
