<template>
  <BlogNav />
  <div class="post-layout fade-in" v-if="post">
    <nav class="toc" v-if="tocItems.length">
      <div class="toc-title">目录 Contents</div>
      <a v-for="item in tocItems" :key="item.id" :href="'#' + item.id"
         class="toc-link" :class="{ active: activeToc === item.id }">
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

  <BlogFooter />
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CommentSection from '../shared/CommentSection.vue'

const route = useRoute()
const post = ref(null)
const loaded = ref(false)
const tocItems = ref([])
const activeToc = ref('')
const prevPost = ref(null)
const nextPost = ref(null)

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function markdownToHtml(md) {
  if (!md) return ''
  let html = md
    // Code blocks
    .replace(/```(\w*)\n([\s\S]*?)```/g, '<pre><code class="lang-$1">$2</code></pre>')
    // Inline code
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    // Headings
    .replace(/^### (.+)$/gm, '<h3>$1</h3>')
    .replace(/^## (.+)$/gm, (match, text) => {
      const id = text.trim().toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '')
      return `<h2 id="${id}">${text}</h2>`
    })
    .replace(/^# (.+)$/gm, '<h1>$1</h1>')
    // Bold and italic
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.+?)\*/g, '<em>$1</em>')
    // Links
    .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>')
    // Images
    .replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1" />')
    // Blockquotes
    .replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>')
    // Unordered lists
    .replace(/^- (.+)$/gm, '<li>$1</li>')
    // Horizontal rules
    .replace(/^---$/gm, '<hr />')
    // Paragraphs
    .replace(/\n\n/g, '</p><p>')
    // Line breaks
    .replace(/\n/g, '<br />')

  // Wrap loose <li> in <ul>
  html = html.replace(/(<li>.*?<\/li>)+/gs, '<ul>$&</ul>')
  // Wrap in paragraph if not starting with block element
  if (!/^<(h[1-6]|pre|ul|blockquote|hr)/.test(html)) {
    html = '<p>' + html + '</p>'
  }
  return html
}

function buildToc(html) {
  const items = []
  const regex = /<h2 id="([^"]+)">([^<]+)<\/h2>/g
  let match
  while ((match = regex.exec(html)) !== null) {
    items.push({ id: match[1], text: match[2] })
  }
  return items
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
  return post.value ? markdownToHtml(post.value.content) : ''
})

onMounted(async () => {
  try {
    const res = await api.get(`/posts/${route.params.id}`)
    post.value = res
    await nextTick()
    tocItems.value = buildToc(renderedContent.value)
    fetchAdjacentPosts()
  } catch {
    post.value = null
  } finally {
    loaded.value = true
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
}
.article-body :deep(h3) {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--fg);
  margin-top: 32px;
  margin-bottom: 12px;
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
.article-body :deep(code) {
  font-family: var(--font-mono);
  font-size: 13px;
  background: var(--tag-bg, #f0f0f0);
  border: 1px solid var(--code-border, var(--border));
  padding: 2px 6px;
  border-radius: 4px;
}
.article-body :deep(pre) {
  background: var(--tag-bg, #f0f0f0);
  border: 1px solid var(--code-border, var(--border));
  border-radius: var(--radius);
  padding: 16px 20px;
  margin-bottom: 20px;
  overflow-x: auto;
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
}
.article-body :deep(pre code) {
  background: none;
  border: none;
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

@media (max-width: 640px) {
  .post-layout {
    padding: 40px 16px 48px;
  }
  .post-nav {
    grid-template-columns: 1fr;
  }
}
</style>
