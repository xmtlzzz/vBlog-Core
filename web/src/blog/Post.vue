<template>
  <BlogNav />
  <div class="post-layout fade-in" v-if="post">
    <nav class="toc" v-if="tocItems.length">
      <div class="toc-title">目录 Contents</div>
      <a v-for="item in tocItems" :key="item.id" :href="'#' + item.id"
         class="toc-link" :class="{ active: activeToc === item.id, 'toc-h3': item.level === 3 }"
         @click="scrollToHeading($event, item.id)">
        {{ item.text }}
      </a>
    </nav>

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
      <p class="article-deck" v-if="post.excerpt">{{ post.excerpt }}</p>
    </header>

    <div class="article-author">
      <div class="author-avatar">{{ post.author?.[0] || 'V' }}</div>
      <div>
        <div class="author-name">{{ post.author || 'vBlog Admin' }}</div>
        <div class="author-role">全栈工程师 / 极客博主</div>
      </div>
    </div>

    <div class="article-body" v-html="renderedContent"></div>

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
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api/request'
import { formatDate } from '../utils/format'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'
import CommentSection from '../shared/CommentSection.vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark-dimmed.css'

const route = useRoute()
const post = ref(null)
const loaded = ref(false)
const tocItems = ref([])
const activeToc = ref('')
const prevPost = ref(null)
const nextPost = ref(null)
const showTop = ref(false)
let observer = null

function onScroll() {
  showTop.value = window.scrollY > 400
}
function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function scrollToHeading(e, id) {
  e.preventDefault()
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
    activeToc.value = id
    history.replaceState(null, '', `#${id}`)
  }
}

function generateSlug(text, index) {
  const clean = text
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')
    .replace(/[^\w\u4e00-\u9fa5\-]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
  return clean || `heading-${index}`
}

// 使用 markdown-it 渲染，html:false（默认值）转义内嵌 HTML，防止存储型 XSS；开启 linkify 方便裸链接
const md = new MarkdownIt({ html: false, linkify: true })

// 为 h2/h3 注入锚点 id（供右侧目录定位），支持中文字符与重复标题去重
md.core.ruler.push('vblog_heading_ids', (state) => {
  const tokens = state.tokens
  const headingCounts = {}
  let headingIdx = 0
  for (let i = 0; i < tokens.length; i++) {
    const t = tokens[i]
    if (t.type === 'heading_open' && (t.tag === 'h2' || t.tag === 'h3')) {
      const inline = tokens[i + 1]
      const text = inline && inline.type === 'inline' ? inline.content.trim() : ''
      headingIdx++
      let slug = generateSlug(text, headingIdx)
      if (headingCounts[slug] !== undefined) {
        headingCounts[slug]++
        slug = `${slug}-${headingCounts[slug]}`
      } else {
        headingCounts[slug] = 0
      }
      t.attrSet('id', slug)
    }
  }
})

// 代码块语法高亮与复制按钮封装
md.renderer.rules.fence = function (tokens, idx) {
  const token = tokens[idx]
  const info = token.info ? token.info.trim() : ''
  const lang = info ? info.split(/\s+/)[0] : ''
  const code = token.content

  let highlighted = md.utils.escapeHtml(code)
  if (lang && hljs.getLanguage(lang)) {
    try {
      highlighted = hljs.highlight(code, { language: lang, ignoreIllegals: true }).value
    } catch {}
  } else {
    try {
      highlighted = hljs.highlightAuto(code).value
    } catch {}
  }

  const langLabel = (lang || 'code').toLowerCase()
  const encoded = encodeURIComponent(code)

  return `<div class="code-block-wrapper"><div class="code-header"><span class="code-lang">${langLabel}</span><button class="copy-code-btn" type="button" data-code="${encoded}">复制</button></div><pre><code class="hljs ${lang ? 'language-' + lang : ''}">${highlighted}</code></pre></div>\n`
}

function buildToc(html) {
  const items = []
  const regex = /<(h[23]) id="([^"]*)">([\s\S]*?)<\/\1>/g
  let match
  while ((match = regex.exec(html)) !== null) {
    items.push({
      level: match[1] === 'h3' ? 3 : 2,
      id: match[2],
      text: match[3].replace(/<[^>]+>/g, '').trim()
    })
  }
  return items
}

function initScrollSpy() {
  if (observer) {
    observer.disconnect()
    observer = null
  }
  const headings = document.querySelectorAll('.article-body h2[id], .article-body h3[id]')
  if (!headings.length) return

  observer = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (entry.isIntersecting) {
        activeToc.value = entry.target.id
        break
      }
    }
  }, {
    rootMargin: '-80px 0px -60% 0px',
    threshold: 0.1
  })

  headings.forEach(h => observer.observe(h))
}

function handleArticleClick(e) {
  const btn = e.target.closest('.copy-code-btn')
  if (!btn) return
  const rawCode = decodeURIComponent(btn.getAttribute('data-code') || '')
  if (!rawCode) return

  navigator.clipboard.writeText(rawCode).then(() => {
    btn.textContent = '已复制!'
    btn.classList.add('copied')
    setTimeout(() => {
      btn.textContent = '复制'
      btn.classList.remove('copied')
    }, 2000)
  }).catch(() => {
    btn.textContent = '复制失败'
    setTimeout(() => {
      btn.textContent = '复制'
    }, 2000)
  })
}

async function fetchAdjacentPosts() {
  try {
    const res = await api.get('/posts', { params: { page: 1, per_page: 100, status: 'published' } })
    const allPosts = res.data || []
    const currentId = Number(route.params.id)
    const idx = allPosts.findIndex(p => p.id === currentId)
    if (idx > 0) prevPost.value = allPosts[idx - 1]
    if (idx >= 0 && idx < allPosts.length - 1) nextPost.value = allPosts[idx + 1]
  } catch {
    // silently fail
  }
}

const renderedContent = computed(() => {
  return post.value ? md.render(post.value.content || '') : ''
})

async function loadPost(id) {
  post.value = null
  prevPost.value = null
  nextPost.value = null
  tocItems.value = []
  try {
    const res = await api.get(`/posts/${id}`)
    post.value = res
    await nextTick()
    tocItems.value = buildToc(renderedContent.value)
    fetchAdjacentPosts()
    nextTick(() => {
      initScrollSpy()
      if (window.location.hash) {
        const targetId = decodeURIComponent(window.location.hash.slice(1))
        const targetEl = document.getElementById(targetId)
        if (targetEl) {
          targetEl.scrollIntoView({ behavior: 'smooth' })
          activeToc.value = targetId
        }
      }
    })
  } catch {
    post.value = null
  } finally {
    loaded.value = true
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  window.addEventListener('click', handleArticleClick)
  loadPost(route.params.id)
})

watch(() => route.params.id, (newId) => {
  if (newId) {
    window.scrollTo({ top: 0 })
    loadPost(newId)
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('click', handleArticleClick)
  if (observer) {
    observer.disconnect()
    observer = null
  }
})
</script>

<style scoped>
.post-layout {
  display: flex;
  justify-content: center;
  gap: 48px;
  max-width: 1080px;
  margin: 0 auto;
  padding: 64px 24px 80px;
}
.article {
  max-width: 720px;
  min-width: 0;
  flex: 1;
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
.toc {
  display: none;
  width: 200px;
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
.toc-link {
  display: block;
  font-size: 13px;
  color: var(--muted);
  text-decoration: none;
  padding: 6px 0 6px 12px;
  border-left: 2px solid transparent;
  transition: all 0.15s;
  line-height: 1.4;
}
.toc-link:hover,
.toc-link.active {
  color: var(--accent);
  border-left-color: var(--accent);
}
.toc-link.toc-h3 {
  padding-left: 20px;
  font-size: 12px;
}
@media (min-width: 1100px) {
  .toc {
    display: block;
  }
}
.article-body {
  font-size: 16px;
  line-height: 1.75;
  color: var(--fg);
}
.article-body :deep(h2) {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--fg);
  margin-top: 48px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
  scroll-margin-top: 72px;
}
.article-body :deep(h3) {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--fg);
  margin-top: 32px;
  margin-bottom: 12px;
  scroll-margin-top: 72px;
}
.article-body :deep(p) {
  margin-bottom: 16px;
}
.article-body :deep(ul),
.article-body :deep(ol) {
  margin-bottom: 16px;
  padding-left: 24px;
}
.article-body :deep(li) {
  margin-bottom: 8px;
}
.article-body :deep(a) {
  color: var(--accent);
  text-decoration: underline;
  text-underline-offset: 2px;
}
.article-body :deep(a:hover) {
  opacity: 0.8;
}
.article-body :deep(:not(pre) > code) {
  font-family: var(--font-mono);
  font-size: 13px;
  background: var(--tag-bg, #f0f0f0);
  border: 1px solid var(--code-border, var(--border));
  padding: 2px 6px;
  border-radius: 4px;
}
.article-body :deep(.code-block-wrapper) {
  position: relative;
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border);
  background: #22272e;
}
.article-body :deep(.code-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: rgba(0, 0, 0, 0.25);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.article-body :deep(.code-lang) {
  font-size: 12px;
  font-family: var(--font-mono);
  color: #768390;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 500;
}
.article-body :deep(.copy-code-btn) {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #adbac7;
  font-size: 11px;
  padding: 3px 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
  line-height: 1.4;
}
.article-body :deep(.copy-code-btn:hover) {
  background: rgba(255, 255, 255, 0.18);
  color: #fff;
}
.article-body :deep(.copy-code-btn.copied) {
  background: #347d39;
  border-color: #46954a;
  color: #fff;
}
.article-body :deep(pre) {
  margin-bottom: 0;
  border-radius: 0;
  border: none;
  background: #22272e;
  padding: 16px 20px;
  overflow-x: auto;
  font-family: var(--font-mono);
  font-size: 13.5px;
  line-height: 1.6;
}
.article-body :deep(pre code.hljs) {
  background: transparent;
  padding: 0;
}
.article-body :deep(blockquote) {
  border-left: 3px solid var(--accent);
  padding: 12px 20px;
  margin: 24px 0;
  color: var(--fg);
  background: var(--accent-soft);
  border-radius: 0 var(--radius) var(--radius) 0;
}
.article-body :deep(hr) {
  border: none;
  border-top: 1px solid var(--border);
  margin: 32px 0;
}
.article-body :deep(img) {
  max-width: 100%;
  border-radius: var(--radius);
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
