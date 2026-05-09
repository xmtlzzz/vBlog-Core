<template>
  <div class="admin-layout">
    <aside class="sidebar" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand">vBlog Admin</div>
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div class="nav-section-title">概览</div>
          <router-link to="/admin" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">📊</span> 仪表盘
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">内容</div>
          <router-link to="/admin/posts" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">📝</span> 文章管理
          </router-link>
          <router-link to="/admin/tags" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">🏷️</span> 标签管理
          </router-link>
          <router-link to="/admin/comments" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">💬</span> 评论管理
          </router-link>
        </div>
        <div class="nav-section">
          <div class="nav-section-title">系统</div>
          <router-link to="/admin/custom" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">🧩</span> 组件定制
          </router-link>
          <router-link to="/admin/settings" class="nav-link" @click="sidebarOpen = false">
            <span class="nav-icon">⚙️</span> 系统设置
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
        </div>
      </div>
    </aside>

    <div class="main-area">
      <header class="topbar">
        <div class="topbar-left">
          <button class="menu-toggle" @click="sidebarOpen = !sidebarOpen">☰</button>
          <router-link to="/" class="back-link">← 首页</router-link>
          <span class="breadcrumb-sep">/</span>
          <span class="breadcrumb-current">{{ currentPage }}</span>
        </div>
        <div class="topbar-right">
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
import { useRoute } from 'vue-router'
import { useThemeStore } from '../stores/theme'

const themeStore = useThemeStore()
const route = useRoute()
const sidebarOpen = ref(false)

const pageMap = {
  '/admin': '仪表盘',
  '/admin/posts': '文章管理',
  '/admin/tags': '标签管理',
  '/admin/comments': '评论管理',
  '/admin/custom': '组件定制',
  '/admin/settings': '系统设置'
}

const currentPage = computed(() => pageMap[route.path] || '仪表盘')
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
  background: var(--surface);
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
.nav-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  font-size: 14px;
  color: var(--muted);
  text-decoration: none;
  transition: all 0.15s;
  border-left: 3px solid transparent;
}
.nav-link:hover {
  color: var(--fg);
  background: var(--accent-soft);
}
.nav-link.router-link-exact-active {
  color: var(--accent);
  background: var(--accent-soft);
  border-left-color: var(--accent);
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
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
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
  background: var(--surface);
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
.back-link {
  font-size: 14px;
  color: var(--muted);
  text-decoration: none;
  transition: color 0.15s;
}
.back-link:hover {
  color: var(--accent);
}
.breadcrumb-sep {
  color: var(--border);
  font-size: 14px;
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
.content {
  flex: 1;
  padding: 24px;
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
