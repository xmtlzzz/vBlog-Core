<template>
  <BlogNav />
  <div :class="['post-layout', 'fade-in', { 'has-toc': tocItems.length }]" v-if="post">
    <!-- 宽屏对称占位：与右侧目录等宽，确保文章主体绝对居中 -->
    <div class="toc-spacer" v-if="tocItems.length" aria-hidden="true"></div>
    <article class="article">
      <router-link to="/" class="back-link">← 返回首页</router-link>

      <header class="article-header">
        <div class="article-meta">
          <span
            v-for="tag in (post.tags || [])"
            :key="tag.id || tag.name || tag"
            class="tag"
          >{{ tag.name || tag }}</span>
          <span>{{ formatDate(post.created_at) }}</span>
          <span>{{ post.read_time || 0 }} min</span>
          <span>{{ (post.views || 0).toLocaleString() }} views</span>
        </div>
        <h1 class="article-title">{{ post.title }}</h1>
        <p class="article-deck" v-if="post.excerpt">{{ plainExcerpt(post.excerpt) }}</p>
      </header>

      <div class="article-author">
        <div class="author-avatar">{{ authorName[0] }}</div>
        <div>
          <div class="author-name">{{ authorName }}</div>
          <div v-if="settings.author_bio" class="author-role">{{ settings.author_bio }}</div>
        </div>
      </div>

      <div class="article-body">
        <MdPreview
          :editorId="editorId"
          :modelValue="articleMarkdown(post.content || '')"
          :theme="editorTheme"
          language="zh-CN"
          :previewTheme="'github'"
          :codeTheme="'atom'"
          :showCodeRowNumber="true"
          @onGetCatalog="onGetCatalog"
        />
      </div>

      <footer class="article-footer">
        <div class="footer-tags">
          <span
            v-for="tag in (post.tags || [])"
            :key="tag.id || tag.name || tag"
            class="tag"
          >{{ tag.name || tag }}</span>
        </div>
      </footer>

      <nav class="post-nav" v-if="prevPost || nextPost">
        <router-link v-if="prevPost" :to="'/post/' + prevPost.id" class="post-nav-item prev">
          <div class="post-nav-label">← 上一篇 Previous</div>
          <div class="post-nav-title">{{ prevPost.title }}</div>
        </router-link>
        <router-link v-if="nextPost" :to="'/post/' + nextPost.id" class="post-nav-item next">
          <div class="post-nav-label">下一篇 Next →</div>
          <div class="post-nav-title">{{ nextPost.title }}</div>
        </router-link>
      </nav>

      <CommentSection :post-id="route.params.id" />
    </article>

    <aside class="toc" v-if="tocItems.length">
      <div class="toc-title">目录 Contents</div>
      <MdCatalog
        :editorId="editorId"
        :theme="editorTheme"
        :scrollElement="scrollElement"
        :offsetTop="80"
      />
    </aside>
  </div>

  <article class="article not-found" v-else-if="loaded">
    <h2>文章不存在</h2>
    <p>该文章可能已被删除或链接无效。</p>
    <router-link to="/" class="back-link">← 返回首页</router-link>
  </article>

  <CustomWidgets />
  <BlogFooter />

  <!-- Back to top -->
  <Transition name="fade">
    <button v-show="showTop" class="back-to-top" @click="scrollToTop" title="回到顶部">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 15l-6-6-6 6"/>
      </svg>
    </button>
  </Transition>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api/request'
import { formatDate } from '../utils/format'
import { plainExcerpt, articleMarkdown } from '../utils/markdown'
import { updateMetadata } from '../utils/metadata'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'
import CommentSection from '../shared/CommentSection.vue'
import { useThemeStore } from '../stores/theme'
import { MdPreview, MdCatalog } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

const route = useRoute()
const themeStore = useThemeStore()
const post = ref(null)
const settings = ref({})
const authorName = computed(() => post.value?.author || settings.value.author_name || settings.value.about_name || settings.value.site_title || 'vBlog')
let requestId = 0
const loaded = ref(false)
const tocItems = ref([])
const prevPost = ref(null)
const nextPost = ref(null)
const showTop = ref(false)
const editorId = 'vblog-post-content'
const scrollElement = ref(typeof document !== 'undefined' ? document.documentElement : null)
const editorTheme = computed(() => themeStore.theme === 'dark' ? 'dark' : 'light')

function onScroll() {
  showTop.value = window.scrollY > 400
}
function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function onGetCatalog(list) {
  tocItems.value = list || []
}

async function fetchAdjacentPosts() {
  try {
    const res = await api.get('/posts', { params: { page: 1, per_page: 50, status: 'published' } })
    const allPosts = res.data || []
    const currentId = Number(route.params.id)
    const idx = allPosts.findIndex(p => p.id === currentId)
    if (idx >= 0 && idx < allPosts.length - 1) prevPost.value = allPosts[idx + 1]
    if (idx > 0) nextPost.value = allPosts[idx - 1]
  } catch {
    // silently fail
  }
}

async function loadPost(id) {
  const currentRequest = ++requestId
  loaded.value = false
  post.value = null
  prevPost.value = null
  nextPost.value = null
  tocItems.value = []
  try {
    const [res, siteSettings] = await Promise.all([api.get(`/posts/${id}`), api.get('/settings').catch(() => ({}))])
    if (currentRequest !== requestId) return
    settings.value = siteSettings || {}
    post.value = res
    updateMetadata(settings.value, { title: res.title, description: res.excerpt || res.content, article: true, path: `/post/${id}` })
    fetchAdjacentPosts()
  } catch {
    if (currentRequest === requestId) {
      post.value = null
      updateMetadata(settings.value, { title: '文章不存在' })
    }
  } finally {
    if (currentRequest === requestId) loaded.value = true
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  loadPost(route.params.id)
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    window.scrollTo({ top: 0 })
    loadPost(newId)
  }
})

onUnmounted(() => {
  ++requestId
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.post-layout {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 40px;
  max-width: 1440px;
  margin: 0 auto;
  padding: 64px 24px 80px;
  box-sizing: border-box;
}
.article {
  overflow-wrap: break-word;
  word-break: normal;
  max-width: 820px;
  min-width: 0;
  flex: 1;
  width: 100%;
}
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--muted);
  text-decoration: none;
  margin-bottom: 32px;
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid var(--border);
  transition: all 0.15s;
}
.back-link:hover {
  color: var(--fg);
  border-color: var(--fg);
}
.article-header {
  margin-bottom: 24px;
}
.article-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  font-size: 13px;
  color: var(--muted);
  flex-wrap: wrap;
}
.tag {
  display: inline-block;
  background: var(--tag-bg);
  color: var(--tag-fg);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}
.article-title {
  font-family: var(--font-display);
  font-size: clamp(28px, 4vw, 36px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.15;
  color: var(--fg);
  margin-bottom: 16px;
}
.article-deck {
  font-size: 17px;
  color: var(--muted);
  line-height: 1.6;
}
.article-author {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid var(--border);
  margin-top: 20px;
  margin-bottom: 32px;
}
.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
}
.author-name {
  font-size: 14px;
  font-weight: 500;
}
.author-role {
  font-size: 12px;
  color: var(--muted);
}
.toc-spacer {
  display: none;
  width: 240px;
  flex-shrink: 0;
  pointer-events: none;
  visibility: hidden;
}
.toc {
  display: none;
  width: 240px;
  flex-shrink: 0;
  position: sticky;
  top: 80px;
  align-self: flex-start;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}
.toc-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}
.toc :deep(.md-editor-catalog-indicator) {
  background-color: var(--accent);
}
.toc :deep(.md-editor-catalog-active > span) {
  color: var(--accent);
  font-weight: 500;
}
.toc :deep(.md-editor-catalog-link span:hover) {
  color: var(--accent);
}
@media (min-width: 1100px) {
  .toc {
    display: block;
  }
}
@media (min-width: 1360px) {
  .toc-spacer {
    display: block;
  }
}

.article-body {
  min-width: 0;
  font-size: 16px;
  line-height: 1.75;
  color: var(--fg);
}
.article-body :deep(img) { max-width: 100%; height: auto; }
.article-body :deep(pre), .article-body :deep(.md-editor-code) { max-width: 100%; overflow-x: auto; }
.article-body :deep(table) { display: block; max-width: 100%; overflow-x: auto; }
.article-body :deep(.md-editor-previewOnly) {
  background: transparent !important;
  border: none;
  padding: 0;
}
.article-body :deep(.md-editor-preview) {
  padding: 0;
  font-family: var(--font-sans);
  color: var(--fg);
}
.article-body :deep(code),
.article-body :deep(pre),
.article-body :deep(.md-editor-code) {
  font-family: var(--font-mono) !important;
}
.article-body :deep(.md-editor-code) {
  border-radius: 8px;
  margin: 20px 0;
  border: 1px solid var(--border);
}

.article-footer {
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}
.footer-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.post-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
}
.post-nav-item {
  text-decoration: none;
  color: var(--fg);
  padding: 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  transition: all 0.15s;
}
.post-nav-item:hover {
  border-color: var(--fg);
}
.post-nav-item.next {
  text-align: right;
}
.post-nav-label {
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 4px;
}
.post-nav-title {
  font-size: 14px;
  font-weight: 500;
}
.not-found {
  text-align: center;
  padding-top: 120px;
}
.not-found h2 {
  font-family: var(--font-display);
  font-size: 24px;
  margin-bottom: 8px;
}
.not-found p {
  color: var(--muted);
  margin-bottom: 24px;
}

/* Back to top button */
.back-to-top {
  position: fixed;
  bottom: 32px;
  right: 32px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--fg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  transition: all 0.2s ease;
  z-index: 50;
}
.back-to-top:hover {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.12);
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .post-layout {
    padding: 40px 16px 48px;
  }
  .post-nav {
    grid-template-columns: 1fr;
  }
}
</style>
