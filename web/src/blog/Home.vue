<template>
  <BlogNav />
  <div class="page-enter">
  <!-- Hero -->
  <header class="hero fade-in">
    <h1>写代码的人，<br/>也写点别的。</h1>
    <p>一个关于系统设计、工程实践与极客生活的博客。用 Markdown 写作，为 vibe coder 而建。</p>
  </header>

  <!-- Stats bar -->
  <section class="stats-bar fade-in" style="animation-delay: 150ms">
    <div class="stat-item" v-for="(s, i) in statItems" :key="i" :style="{ animationDelay: (200 + i * 80) + 'ms' }">
      <span class="stat-value">{{ s.value }}</span>
      <span class="stat-label">{{ s.label }}</span>
    </div>
  </section>

  <!-- Search overlay (Ctrl+F) -->
  <Transition name="search-slide">
    <section v-if="showSearch" class="search-overlay">
      <div class="search-box">
        <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
        </svg>
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          class="search-input"
          placeholder="搜索文章..."
          @input="debounceSearch"
          @keyup.escape="closeSearch"
        />
        <button v-if="searchQuery" class="search-clear" @click="closeSearch">✕</button>
      </div>
    </section>
  </Transition>

  <!-- Post list -->
  <section class="post-list">
    <TransitionGroup name="post-list" tag="div">
      <PostCard v-for="(post, i) in posts" :key="post.id" :post="post" :style="{ animationDelay: (i * 60) + 'ms' }" class="fade-in" />
    </TransitionGroup>
    <div v-if="posts.length === 0" class="empty-state fade-in">
      <p>暂无文章</p>
    </div>
  </section>

  <!-- Pagination -->
  <section v-if="total > perPage" class="pagination-wrap fade-in">
    <el-pagination
      layout="prev, pager, next, jumper"
      :total="total"
      :page-size="perPage"
      :current-page="page"
      @current-change="handlePageChange"
    />
  </section>

  </div>
  <CustomWidgets />
  <BlogFooter />
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'
import PostCard from '../shared/PostCard.vue'

const stats = ref({ total_posts: 0, total_views: 0, total_tags: 0 })
const posts = ref([])
const searchQuery = ref('')
const showSearch = ref(false)
const searchInputRef = ref(null)
const page = ref(1)
const perPage = 5
const total = ref(0)
let searchTimer = null

const statItems = computed(() => [
  { value: stats.value.total_posts, label: '篇文章 Articles' },
  { value: stats.value.total_views.toLocaleString(), label: '次阅读 Views' },
  { value: stats.value.total_tags, label: '个标签 Tags' }
])

async function fetchPosts() {
  const params = { page: page.value, per_page: perPage, status: 'published' }
  if (searchQuery.value) params.search = searchQuery.value
  const res = await api.get('/posts', { params })
  posts.value = res.data || []
  total.value = res.total || 0
}

function debounceSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { page.value = 1; fetchPosts() }, 300)
}

function openSearch() {
  showSearch.value = true
  nextTick(() => searchInputRef.value?.focus())
}

function closeSearch() {
  showSearch.value = false
  if (searchQuery.value) {
    searchQuery.value = ''
    page.value = 1
    fetchPosts()
  }
}

function onKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
    e.preventDefault()
    openSearch()
  }
}

function handlePageChange(p) {
  page.value = p
  fetchPosts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  window.addEventListener('keydown', onKeydown)
  const statsRes = await api.get('/dashboard/stats').catch(() => ({ total_posts: 0, total_views: 0, total_tags: 0 }))
  stats.value = {
    total_posts: statsRes.total_posts || 0,
    total_views: statsRes.total_views || 0,
    total_tags: statsRes.total_tags || 0
  }
  await fetchPosts()
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
.hero {
  max-width: 1080px;
  margin: 0 auto;
  padding: 72px 24px 48px;
}
.hero h1 {
  font-family: var(--font-display);
  font-size: clamp(32px, 5vw, 48px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.1;
  color: var(--fg);
  margin-bottom: 12px;
}
.hero p {
  font-size: 17px;
  color: var(--muted);
  max-width: 520px;
  line-height: 1.6;
}
.stats-bar {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 40px;
  display: flex;
  gap: 32px;
}
.stat-item {
  display: flex;
  align-items: baseline;
  gap: 6px;
  animation: fadeIn 0.4s ease both;
}
.stat-value {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--fg);
}
.stat-label {
  font-size: 13px;
  color: var(--muted);
}
/* Search overlay */
.search-overlay {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 28px;
}
.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  max-width: 520px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 0 16px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.search-box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}
.search-icon {
  flex-shrink: 0;
  color: var(--muted);
  display: flex;
  align-items: center;
}
.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 15px;
  color: var(--fg);
  padding: 12px 0;
  font-family: var(--font-sans);
}
.search-input::placeholder {
  color: var(--muted);
}
.search-clear {
  flex-shrink: 0;
  background: none;
  border: none;
  color: var(--muted);
  font-size: 14px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  line-height: 1;
  transition: all 0.15s;
}
.search-clear:hover {
  color: var(--fg);
  background: var(--border);
}
.search-slide-enter-active,
.search-slide-leave-active {
  transition: all 0.2s ease;
}
.search-slide-enter-from,
.search-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
/* TransitionGroup for post list */
.post-list-enter-active {
  animation: fadeIn 0.3s ease both;
}
.post-list-leave-active {
  animation: fadeOut 0.2s ease both;
}
@keyframes fadeOut {
  to { opacity: 0; transform: translateY(-8px); }
}
.post-list {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 20px;
}
.empty-state {
  text-align: center;
  padding: 64px 24px;
  color: var(--muted);
}
.pagination-wrap {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 80px;
  display: flex;
  justify-content: center;
}

@media (max-width: 640px) {
  .hero {
    padding: 48px 16px 32px;
  }
  .stats-bar {
    padding: 0 16px 24px;
    gap: 20px;
  }
  .stat-value {
    font-size: 20px;
  }
  .search-overlay {
    padding: 0 16px 20px;
  }
  .search-box {
    max-width: 100%;
  }
  .post-list {
    padding: 0 16px 20px;
  }
  .pagination-wrap {
    padding: 0 16px 48px;
  }
}
</style>
