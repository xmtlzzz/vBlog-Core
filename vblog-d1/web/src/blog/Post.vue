<template>
  <div
    class="reading-progress-bar"
    :style="{ width: `${scrollProgress}%` }"
    aria-hidden="true"
  ></div>
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
          <span v-if="readingStats.words">{{ readingStats.words.toLocaleString() }} 字</span>
          <span>约 {{ post.read_time || readingStats.minutes }} min</span>
          <span>{{ (post.views || 0).toLocaleString() }} views</span>
          <button
            type="button"
            class="meta-action-btn"
            @click="openShareModal"
            title="分享文章或生成海报卡片"
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="18" cy="5" r="3"></circle>
              <circle cx="6" cy="12" r="3"></circle>
              <circle cx="18" cy="19" r="3"></circle>
              <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
              <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
            </svg>
            <span>分享 / 海报</span>
          </button>
        </div>
        <h1 class="article-title">{{ post.title }}</h1>
      </header>

      <div class="article-author">
        <div class="author-avatar">{{ authorName[0] }}</div>
        <div>
          <div class="author-name">{{ authorName }}</div>
          <div v-if="settings.author_bio" class="author-role">{{ settings.author_bio }}</div>
        </div>
      </div>

      <div class="article-body" @click="handleBodyClick">
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

  <!-- Floating Action Buttons -->
  <div class="floating-actions" v-if="post">
    <button
      v-if="tocItems.length"
      type="button"
      class="floating-btn toc-trigger-btn"
      @click="showMobileToc = true"
      title="文章目录"
      aria-label="文章目录"
    >
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="3" y1="6" x2="21" y2="6"></line>
        <line x1="3" y1="12" x2="15" y2="12"></line>
        <line x1="3" y1="18" x2="18" y2="18"></line>
      </svg>
    </button>
    <Transition name="fade">
      <button
        v-show="showTop"
        type="button"
        class="floating-btn back-to-top"
        @click="scrollToTop"
        title="回到顶部"
        aria-label="回到顶部"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M18 15l-6-6-6 6"/>
        </svg>
      </button>
    </Transition>
  </div>

  <!-- Mobile TOC Drawer -->
  <Teleport to="body">
    <Transition name="drawer-fade">
      <div v-if="showMobileToc" class="mobile-toc-overlay" @click.self="showMobileToc = false">
        <div class="mobile-toc-sheet">
          <div class="mobile-toc-header">
            <span class="mobile-toc-title">目录 Contents</span>
            <button class="mobile-toc-close" @click="showMobileToc = false" aria-label="关闭" type="button">✕</button>
          </div>
          <div class="mobile-toc-content" @click="onMobileCatalogClick">
            <MdCatalog
              :editorId="editorId"
              :theme="editorTheme"
              :scrollElement="scrollElement"
              :offsetTop="80"
            />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- Image Lightbox Modal -->
  <Teleport to="body">
    <Transition name="lightbox-fade">
      <div
        v-if="previewImage"
        class="lightbox-overlay"
        @click="previewImage = null"
      >
        <div class="lightbox-content">
          <img :src="previewImage" :alt="previewImageAlt" class="lightbox-img" />
          <div class="lightbox-caption" v-if="previewImageAlt">{{ previewImageAlt }}</div>
        </div>
        <button class="lightbox-close" aria-label="关闭" @click="previewImage = null" type="button">✕</button>
      </div>
    </Transition>
  </Teleport>

  <!-- Copy Link Toast -->
  <Teleport to="body">
    <Transition name="toast-fade">
      <div v-if="copied" class="copy-toast">已复制文章链接到剪贴板</div>
    </Transition>
  </Teleport>

  <!-- Share & Poster Modal -->
  <Teleport to="body">
    <Transition name="lightbox-fade">
      <div
        v-if="showShareModal"
        class="poster-overlay"
        @click.self="showShareModal = false"
      >
        <div class="poster-modal-card">
          <div class="poster-modal-header">
            <h3>分享文章与海报</h3>
            <button class="poster-close-btn" @click="showShareModal = false" type="button" aria-label="关闭">✕</button>
          </div>
          <div class="poster-preview-area">
            <canvas ref="posterCanvasRef" class="poster-canvas"></canvas>
            <div v-if="generatingPoster" class="poster-loading">
              <div class="poster-spinner"></div>
              <span>海报生成中...</span>
            </div>
          </div>
          <div class="poster-modal-actions">
            <button type="button" class="poster-btn poster-btn-primary" @click="downloadPoster">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              </svg>
              <span>保存海报</span>
            </button>
            <button type="button" class="poster-btn" @click="copyPostUrl">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <span>{{ copied ? '已复制' : '复制链接' }}</span>
            </button>
            <button v-if="canNativeShare" type="button" class="poster-btn" @click="nativeShare">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="18" cy="5" r="3"></circle>
                <circle cx="6" cy="12" r="3"></circle>
                <circle cx="18" cy="19" r="3"></circle>
                <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line>
                <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line>
              </svg>
              <span>系统分享</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { createEveryQRCodeIdentity, createQRSvgPath } from '@every-qrcode/core'
import api from '../api/request'
import { formatDate } from '../utils/format'
import { articleMarkdown } from '../utils/markdown'
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
const scrollProgress = ref(0)
const previewImage = ref(null)
const previewImageAlt = ref('')
const showMobileToc = ref(false)
const copied = ref(false)
let copyTimer = null

// ── 分享与海报状态 ──
const showShareModal = ref(false)
const generatingPoster = ref(false)
const posterCanvasRef = ref(null)
const canNativeShare = computed(() => typeof navigator !== 'undefined' && !!navigator.share)

const editorId = 'vblog-post-content'
const scrollElement = ref(typeof document !== 'undefined' ? document.documentElement : null)
const editorTheme = computed(() => themeStore.theme === 'dark' ? 'dark' : 'light')

function calcReadingStats(raw) {
  if (!raw) return { words: 0, minutes: 1 }
  const clean = raw
    .replace(/```[\s\S]*?```/g, '')
    .replace(/!\[.*?\]\(.*?\)/g, '')
    .replace(/\[(.*?)\]\(.*?\)/g, '$1')
    .replace(/<[^>]+>/g, '')
  const cjkMatches = clean.match(/[\u4e00-\u9fa5]/g) || []
  const nonCjk = clean.replace(/[\u4e00-\u9fa5]/g, ' ')
  const wordMatches = nonCjk.match(/[a-zA-Z0-9_\u00C0-\u024F]+/g) || []
  const total = cjkMatches.length + wordMatches.length
  const minutes = Math.max(1, Math.ceil(total / 350))
  return { words: total, minutes }
}

const readingStats = computed(() => calcReadingStats(post.value?.content || ''))

function onScroll() {
  const total = document.documentElement.scrollHeight - window.innerHeight
  scrollProgress.value = total > 0 ? Math.min(100, Math.max(0, (window.scrollY / total) * 100)) : 0
  showTop.value = window.scrollY > 400
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function copyPostUrl() {
  try {
    navigator.clipboard.writeText(window.location.href)
    copied.value = true
    clearTimeout(copyTimer)
    copyTimer = setTimeout(() => {
      copied.value = false
    }, 2200)
  } catch {
    // clipboard api unavailable
  }
}

async function openShareModal() {
  showShareModal.value = true
  await nextTick()
  generatePoster()
}

function nativeShare() {
  if (navigator.share) {
    navigator.share({
      title: post.value?.title || 'vBlog',
      text: post.value?.excerpt || post.value?.title || '',
      url: window.location.href
    }).catch(() => {})
  }
}

function downloadPoster() {
  const canvas = posterCanvasRef.value
  if (!canvas) return
  const link = document.createElement('a')
  link.download = `${post.value?.title || 'post'}-poster.png`
  link.href = canvas.toDataURL('image/png')
  link.click()
}

function fillRoundRect(ctx, x, y, width, height, radius, fillStyle) {
  ctx.save()
  ctx.fillStyle = fillStyle
  if (ctx.roundRect) {
    ctx.beginPath()
    ctx.roundRect(x, y, width, height, radius)
    ctx.fill()
  } else {
    ctx.fillRect(x, y, width, height)
  }
  ctx.restore()
}

function strokeRoundRect(ctx, x, y, width, height, radius, strokeStyle, lineWidth = 1) {
  ctx.save()
  ctx.strokeStyle = strokeStyle
  ctx.lineWidth = lineWidth
  if (ctx.roundRect) {
    ctx.beginPath()
    ctx.roundRect(x, y, width, height, radius)
    ctx.stroke()
  } else {
    ctx.strokeRect(x, y, width, height)
  }
  ctx.restore()
}

function wrapPosterText(ctx, text, x, y, maxWidth, lineHeight, maxLines = 4) {
  const chars = String(text || '').trim().split('')
  let line = ''
  let currentY = y
  let linesCount = 0

  for (let n = 0; n < chars.length; n++) {
    const testLine = line + chars[n]
    const metrics = ctx.measureText(testLine)
    if (metrics.width > maxWidth && n > 0) {
      linesCount++
      if (linesCount >= maxLines) {
        ctx.fillText(line + '...', x, currentY)
        return currentY + lineHeight
      }
      ctx.fillText(line, x, currentY)
      line = chars[n]
      currentY += lineHeight
    } else {
      line = testLine
    }
  }
  if (line) {
    ctx.fillText(line, x, currentY)
    currentY += lineHeight
  }
  return currentY
}

async function generatePoster() {
  const canvas = posterCanvasRef.value
  if (!canvas || !post.value) return
  generatingPoster.value = true

  try {
    const width = 640
    const height = 860
    canvas.width = width * 2
    canvas.height = height * 2

    const ctx = canvas.getContext('2d')
    ctx.scale(2, 2)

    const isDark = themeStore.theme === 'dark'
    const bgGradient = ctx.createLinearGradient(0, 0, width, height)
    if (isDark) {
      bgGradient.addColorStop(0, '#1c1f26')
      bgGradient.addColorStop(1, '#0e1014')
    } else {
      bgGradient.addColorStop(0, '#ffffff')
      bgGradient.addColorStop(1, '#f1f5f9')
    }

    fillRoundRect(ctx, 0, 0, width, height, 20, bgGradient)
    strokeRoundRect(ctx, 0, 0, width, height, 20, isDark ? '#2d3342' : '#e2e8f0', 1.5)

    const siteTitle = settings.value.site_title || 'vBlog'
    ctx.fillStyle = isDark ? '#60a5fa' : '#2563eb'
    ctx.beginPath()
    ctx.arc(42, 46, 6, 0, Math.PI * 2)
    ctx.fill()

    ctx.fillStyle = isDark ? '#f3f4f6' : '#1e293b'
    ctx.font = 'bold 18px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif'
    ctx.fillText(siteTitle, 56, 52)

    const tag = post.value.tags?.[0]?.name || post.value.tags?.[0] || 'Article'
    const tagText = '#' + tag
    ctx.font = '13px -apple-system, BlinkMacSystemFont, sans-serif'
    const tagWidth = ctx.measureText(tagText).width + 20
    fillRoundRect(ctx, width - 36 - tagWidth, 34, tagWidth, 26, 6, isDark ? 'rgba(59, 130, 246, 0.15)' : 'rgba(37, 99, 235, 0.08)')
    ctx.fillStyle = isDark ? '#93c5fd' : '#2563eb'
    ctx.fillText(tagText, width - 36 - tagWidth + 10, 51)

    ctx.strokeStyle = isDark ? '#262b37' : '#e2e8f0'
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.moveTo(36, 76)
    ctx.lineTo(width - 36, 76)
    ctx.stroke()

    ctx.fillStyle = isDark ? '#ffffff' : '#0f172a'
    ctx.font = 'bold 28px -apple-system, BlinkMacSystemFont, "SF Pro Display", sans-serif'
    let currentY = 120
    currentY = wrapPosterText(ctx, post.value.title, 36, currentY, width - 72, 38, 3)

    currentY += 12
    ctx.fillStyle = isDark ? '#9ca3af' : '#64748b'
    ctx.font = '13px -apple-system, BlinkMacSystemFont, sans-serif'
    const dateStr = formatDate(post.value.created_at)
    const readMin = post.value.read_time || readingStats.value.minutes || 1
    const words = readingStats.value.words || 0
    const metaStr = `${dateStr}  ·  约 ${readMin} 分钟阅读  ·  ${words.toLocaleString()} 字`
    ctx.fillText(metaStr, 36, currentY)

    currentY += 28
    const rawExcerpt = (post.value.excerpt || post.value.content || '')
      .replace(/```[\s\S]*?```/g, '')
      .replace(/!\[.*?\]\(.*?\)/g, '')
      .replace(/[#*`>~_-]/g, '')
      .trim()
    const boxHeight = 240
    fillRoundRect(ctx, 36, currentY, width - 72, boxHeight, 12, isDark ? '#14171e' : '#f8fafc')
    strokeRoundRect(ctx, 36, currentY, width - 72, boxHeight, 12, isDark ? '#232834' : '#e5e7eb', 1)
    fillRoundRect(ctx, 36, currentY, 4, boxHeight, 2, isDark ? '#3b82f6' : '#2563eb')

    ctx.fillStyle = isDark ? '#d1d5db' : '#334155'
    ctx.font = '15px -apple-system, BlinkMacSystemFont, sans-serif'
    wrapPosterText(ctx, rawExcerpt, 54, currentY + 36, width - 108, 28, 6)

    const footerY = height - 120
    ctx.strokeStyle = isDark ? '#262b37' : '#e2e8f0'
    ctx.beginPath()
    ctx.moveTo(36, footerY - 20)
    ctx.lineTo(width - 36, footerY - 20)
    ctx.stroke()

    ctx.fillStyle = isDark ? '#ffffff' : '#0f172a'
    ctx.font = 'bold 16px -apple-system, BlinkMacSystemFont, sans-serif'
    ctx.fillText(authorName.value, 36, footerY + 16)

    ctx.fillStyle = isDark ? '#9ca3af' : '#64748b'
    ctx.font = '13px -apple-system, BlinkMacSystemFont, sans-serif'
    const bioText = settings.value.description || settings.value.subtitle || '长按或扫描二维码阅读全文'
    ctx.fillText(bioText.slice(0, 24), 36, footerY + 40)

    ctx.fillStyle = isDark ? '#60a5fa' : '#2563eb'
    ctx.font = '12px "JetBrains Mono", monospace'
    const domain = typeof window !== 'undefined' ? window.location.hostname : 'vblog.xmtlz.dev'
    ctx.fillText(domain, 36, footerY + 62)

    const qrSize = 96
    const qrX = width - 36 - qrSize
    const qrY = footerY - 5

    try {
      const url = window.location.href
      const identity = await createEveryQRCodeIdentity(url, { identityScope: 'url' }).catch(() => null)
      if (identity) {
        const qr = createQRSvgPath(identity.qr)
        if (qr && qr.path) {
          fillRoundRect(ctx, qrX - 8, qrY - 8, qrSize + 16, qrSize + 16, 10, '#ffffff')
          strokeRoundRect(ctx, qrX - 8, qrY - 8, qrSize + 16, qrSize + 16, 10, '#e2e8f0', 1)
          ctx.save()
          ctx.translate(qrX, qrY)
          const scale = qrSize / qr.size
          ctx.scale(scale, scale)
          ctx.fillStyle = '#0f172a'
          ctx.fill(new Path2D(qr.path))
          ctx.restore()
        }
      }
    } catch {
      fillRoundRect(ctx, qrX, qrY, qrSize, qrSize, 8, isDark ? '#1e2430' : '#e2e8f0')
      ctx.fillStyle = isDark ? '#9ca3af' : '#475569'
      ctx.font = '12px sans-serif'
      ctx.fillText('vBlog', qrX + 28, qrY + 52)
    }
  } finally {
    generatingPoster.value = false
  }
}

function handleBodyClick(e) {
  const img = e.target.closest('img')
  if (img && img.src && !img.closest('.author-avatar')) {
    previewImage.value = img.src
    previewImageAlt.value = img.alt || ''
  }
}

function onMobileCatalogClick(e) {
  if (e.target.closest('.md-editor-catalog-link') || e.target.closest('span')) {
    setTimeout(() => {
      showMobileToc.value = false
    }, 250)
  }
}

function onKeydown(e) {
  if (e.key === 'Escape') {
    if (previewImage.value) previewImage.value = null
    if (showMobileToc.value) showMobileToc.value = false
    if (showShareModal.value) showShareModal.value = false
  }
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
  window.addEventListener('keydown', onKeydown)
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
  window.removeEventListener('keydown', onKeydown)
  clearTimeout(copyTimer)
})
</script>

<style scoped>
/* Reading progress bar */
.reading-progress-bar {
  position: fixed;
  top: 0;
  left: 0;
  height: 3px;
  background: var(--accent);
  z-index: 999;
  transition: width 0.08s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none;
  box-shadow: 0 0 8px var(--accent);
}

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
.meta-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: 1px solid var(--border);
  color: var(--muted);
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
}
.meta-action-btn:hover {
  color: var(--fg);
  border-color: var(--fg);
}
.article-title {
  font-family: var(--font-display);
  font-size: clamp(28px, 4vw, 36px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.15;
  color: var(--fg);
  margin-bottom: 0;
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
.article-body :deep(img) {
  max-width: 100%;
  height: auto;
  cursor: zoom-in;
  border-radius: 6px;
  transition: opacity 0.2s ease;
}
.article-body :deep(img:hover) {
  opacity: 0.94;
}
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

/* Floating Actions */
.floating-actions {
  position: fixed;
  bottom: 32px;
  right: 32px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  z-index: 50;
}
.floating-btn {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--fg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.floating-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0,0,0,0.14);
}
.toc-trigger-btn {
  display: none;
}
@media (max-width: 1099px) {
  .toc-trigger-btn {
    display: flex;
  }
}

/* Mobile TOC Drawer */
.mobile-toc-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(4px);
  z-index: 2000;
  display: flex;
  align-items: flex-end;
}
.mobile-toc-sheet {
  width: 100%;
  max-height: 75vh;
  background: var(--surface);
  border-top-left-radius: 18px;
  border-top-right-radius: 18px;
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.25);
  animation: slideUpSheet 0.24s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes slideUpSheet {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}
.mobile-toc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}
.mobile-toc-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--fg);
}
.mobile-toc-close {
  background: transparent;
  border: none;
  color: var(--muted);
  font-size: 18px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}
.mobile-toc-content {
  padding: 16px 20px 32px;
  overflow-y: auto;
  max-height: calc(75vh - 60px);
}
.mobile-toc-content :deep(.md-editor-catalog-indicator) {
  background-color: var(--accent);
}
.mobile-toc-content :deep(.md-editor-catalog-active > span) {
  color: var(--accent);
  font-weight: 600;
}
.drawer-fade-enter-active, .drawer-fade-leave-active {
  transition: opacity 0.2s ease;
}
.drawer-fade-enter-from, .drawer-fade-leave-to {
  opacity: 0;
}

/* Lightbox Modal */
.lightbox-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.88);
  backdrop-filter: blur(8px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: zoom-out;
  padding: 24px;
  box-sizing: border-box;
}
.lightbox-content {
  max-width: 90vw;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.lightbox-img {
  max-width: 100%;
  max-height: 85vh;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 16px 48px rgba(0,0,0,0.45);
}
.lightbox-caption {
  margin-top: 12px;
  color: #e5e5e5;
  font-size: 13px;
  text-align: center;
}
.lightbox-close {
  position: absolute;
  top: 24px;
  right: 24px;
  width: 38px;
  height: 38px;
  border-radius: 50%;
  background: rgba(255,255,255,0.15);
  color: #fff;
  border: none;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s;
}
.lightbox-close:hover {
  background: rgba(255,255,255,0.3);
}
.lightbox-fade-enter-active, .lightbox-fade-leave-active {
  transition: opacity 0.2s ease;
}
.lightbox-fade-enter-from, .lightbox-fade-leave-to {
  opacity: 0;
}

/* Copy Toast */
.copy-toast {
  position: fixed;
  bottom: 84px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--fg);
  color: var(--bg);
  padding: 8px 18px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  z-index: 1000;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
  pointer-events: none;
}
.toast-fade-enter-active, .toast-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.toast-fade-enter-from, .toast-fade-leave-to {
  opacity: 0;
  transform: translate(-50%, 8px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ── Poster Modal ── */
.poster-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.poster-modal-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  max-width: 400px;
  width: 100%;
  max-height: 92vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.35);
  animation: posterPop 0.22s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes posterPop {
  from { opacity: 0; transform: scale(0.96) translateY(8px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.poster-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
}

.poster-modal-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: var(--fg);
}

.poster-close-btn {
  background: none;
  border: none;
  font-size: 16px;
  color: var(--muted);
  cursor: pointer;
  padding: 4px;
  border-radius: var(--radius);
  transition: color 0.15s ease;
}

.poster-close-btn:hover {
  color: var(--fg);
}

.poster-preview-area {
  padding: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow-y: auto;
  max-height: calc(90vh - 120px);
  background: var(--bg);
}

.poster-canvas {
  width: 100%;
  max-width: 320px;
  height: auto;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  display: block;
}

.poster-loading {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: var(--surface);
  color: var(--muted);
  font-size: 13px;
}

.poster-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.poster-modal-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  background: var(--surface);
}

.poster-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  padding: 8px 12px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--fg);
  cursor: pointer;
  transition: all 0.18s ease;
}

.poster-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.poster-btn-primary {
  background: var(--accent);
  color: #ffffff;
  border-color: var(--accent);
}

.poster-btn-primary:hover {
  filter: brightness(1.1);
  color: #ffffff;
}

@media (max-width: 640px) {
  .post-layout {
    padding: 40px 16px 48px;
  }
  .post-nav {
    grid-template-columns: 1fr;
  }
  .floating-actions {
    bottom: 24px;
    right: 20px;
  }
}
</style>
