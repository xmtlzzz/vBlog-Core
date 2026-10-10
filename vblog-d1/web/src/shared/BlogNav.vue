<template>
  <nav class="top">
    <div class="inner">
      <router-link class="nav-brand" to="/">
        <img class="mascot" src="/nav-mascot.png" alt="vBlog" />vBlog
      </router-link>
      <div id="blog-navigation" :class="['nav-links', { 'is-open': menuOpen }]">
        <router-link to="/" @click="menuOpen = false">首页 Home</router-link>
        <router-link to="/archives" @click="menuOpen = false">归档 Archives</router-link>
        <router-link to="/modules" @click="menuOpen = false">模块 Modules</router-link>
        <router-link to="/tags" @click="menuOpen = false">标签 Tags</router-link>
        <a :href="`https://vmaker.xmtlz.dev?theme=${themeStore.theme}`" target="_blank" rel="noopener" class="nav-external" @click="menuOpen = false">作品 Projects</a>
        <router-link to="/friends" @click="menuOpen = false">友链 Friends</router-link>
        <router-link to="/about" @click="menuOpen = false">关于 About</router-link>
      </div>
      <div class="nav-right">
        <button class="nav-search-btn" @click="handleSearch" aria-label="搜索文章" title="搜索文章 (⌘K / /)">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
          </svg>
        </button>
        <router-link to="/admin" class="admin-btn">后台<span class="admin-en"> Admin</span></router-link>
        <ThemeToggleAnimated
          :is-dark="themeStore.theme === 'dark'"
          size="sm"
          @toggle="themeStore.toggle()"
        />
        <button
          class="menu-toggle"
          :aria-expanded="menuOpen"
          aria-controls="blog-navigation"
          :aria-label="menuOpen ? '收起导航菜单' : '展开导航菜单'"
          @click="menuOpen = !menuOpen"
        >
          <svg v-if="!menuOpen" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="4" y1="6" x2="20" y2="6"/>
            <line x1="4" y1="12" x2="20" y2="12"/>
            <line x1="4" y1="18" x2="20" y2="18"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </div>
  </nav>

  <!-- 移动端抽屉遮罩 -->
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="menuOpen" class="mobile-nav-scrim" @click="menuOpen = false" />
    </Transition>
  </Teleport>

  <!-- 全局原地 Command Palette 弹窗 -->
  <SearchModal v-model="searchModalOpen" />
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from '../stores/theme'
import SearchModal from './SearchModal.vue'
import ThemeToggleAnimated from './ThemeToggleAnimated.vue'

const route = useRoute()
const themeStore = useThemeStore()
const menuOpen = ref(false)
const searchModalOpen = ref(false)

watch(() => route.fullPath, () => { menuOpen.value = false })

function onGlobalKeydown(event) {
  if (event.key === 'Escape') menuOpen.value = false
  const isCmdK = (event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k'
  const isSlash = event.key === '/' && !event.ctrlKey && !event.metaKey && !event.altKey
  if (isCmdK || isSlash) {
    const active = document.activeElement
    const tag = active?.tagName?.toLowerCase()
    const isInput = tag === 'input' || tag === 'textarea' || active?.isContentEditable
    if (!isInput) {
      event.preventDefault()
      handleSearch()
    }
  }
}

function onOpenSearchEvent() {
  searchModalOpen.value = true
}

onMounted(() => {
  window.addEventListener('keydown', onGlobalKeydown)
  window.addEventListener('vblog-open-search', onOpenSearchEvent)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onGlobalKeydown)
  window.removeEventListener('vblog-open-search', onOpenSearchEvent)
})

function handleSearch() {
  menuOpen.value = false
  searchModalOpen.value = true
}
</script>

<style scoped>
.top {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--nav-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border);
  transition: background 260ms ease, border-color 260ms ease;
}
.inner {
  min-width: 0;
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
.mascot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}
.nav-links {
  min-width: 0;
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
  background: var(--card-hover);
}
[data-theme="dark"] .nav-links a:hover,
[data-theme="dark"] .nav-links :deep(.router-link-active) {
  background: var(--card-hover);
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
.nav-search-btn {
  width: 36px;
  height: 36px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
  border-radius: 6px;
}
.nav-search-btn:hover {
  color: var(--fg);
  border-color: var(--accent);
  background: var(--card-hover);
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
  background: var(--card-hover);
}
.menu-toggle {
  display: none;
  width: 36px;
  height: 36px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm, 6px);
  background: var(--surface);
  color: var(--muted);
  cursor: pointer;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}
.menu-toggle:hover {
  color: var(--fg);
  border-color: var(--accent);
  background: var(--card-hover);
}
button:focus-visible, a:focus-visible { outline: 2px solid var(--accent); outline-offset: 3px; }

@media (max-width: 800px) {
  /* 两行布局：品牌 + 右侧操作一行，导航链接第二行横向滑动，避免折行溢出 */
  .inner {
    height: auto;
    flex-wrap: wrap;
    padding: 10px 12px;
  }
  .nav-brand {
    font-size: 16px;
  }
  .nav-right {
    gap: 6px;
  }
  .admin-btn {
    padding: 5px 8px;
    font-size: 12px;
  }
  .nav-search-btn {
    width: 32px;
    height: 32px;
  }
  .theme-toggle {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }
  .nav-links {
    display: none;
    order: 3;
    flex: 1 1 100%;
    flex-wrap: wrap;
    justify-content: flex-start;
    gap: 2px;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
    padding-bottom: 6px;
    margin-top: 2px;
  }
  .nav-links.is-open { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); margin-top: 12px; padding-top: 12px; border-top: 1px solid var(--border); }
  .menu-toggle { display: flex; }
  .nav-links::-webkit-scrollbar {
    display: none;
  }
  .nav-links a {
    white-space: nowrap;
    flex-shrink: 0;
    padding: 6px 10px;
    font-size: 13px;
  }
}

@media (max-width: 380px) {
  .admin-en {
    display: none;
  }
}

@media (max-width: 360px) {
  .inner {
    padding: 8px 10px;
  }
  .nav-brand {
    font-size: 15px;
    gap: 6px;
  }
  .nav-right {
    gap: 4px;
  }
  .admin-btn {
    padding: 4px 6px;
    font-size: 11px;
  }
}

.mobile-nav-scrim {
  position: fixed;
  inset: 0;
  z-index: 90;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}
</style>
