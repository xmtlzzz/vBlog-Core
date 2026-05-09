<template>
  <BlogNav />
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

  <!-- Filter bar -->
  <section class="filter-bar fade-in" style="animation-delay: 350ms">
    <button
      :class="['filter-btn', { active: !activeTag }]"
      @click="setTag('')"
    >全部 All</button>
    <button
      v-for="tag in tags"
      :key="tag.id || tag.name"
      :class="['filter-btn', { active: activeTag === tag.name }]"
      @click="setTag(tag.name)"
    >{{ tag.name }}</button>
  </section>

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
      layout="prev, pager, next"
      :total="total"
      :page-size="perPage"
      :current-page="page"
      @current-change="handlePageChange"
    />
  </section>

  <BlogFooter />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import PostCard from '../shared/PostCard.vue'

const stats = ref({ total_posts: 0, total_views: 0, total_tags: 0 })
const tags = ref([])
const posts = ref([])
const activeTag = ref('')
const page = ref(1)
const perPage = 5
const total = ref(0)

const statItems = computed(() => [
  { value: stats.value.total_posts, label: '篇文章 Articles' },
  { value: stats.value.total_views.toLocaleString(), label: '次阅读 Views' },
  { value: stats.value.total_tags, label: '个标签 Tags' }
])

async function fetchPosts() {
  const res = await api.get('/posts', {
    params: { page: page.value, per_page: perPage, tag: activeTag.value, status: 'published' }
  })
  posts.value = res.posts || []
  total.value = res.total || 0
}

function setTag(tag) {
  activeTag.value = tag
  page.value = 1
  fetchPosts()
}

function handlePageChange(p) {
  page.value = p
  fetchPosts()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  const [statsRes, tagsRes] = await Promise.all([
    api.get('/dashboard/stats').catch(() => ({ total_posts: 0, total_views: 0, total_tags: 0 })),
    api.get('/tags').catch(() => [])
  ])
  stats.value = {
    total_posts: statsRes.total_posts || 0,
    total_views: statsRes.total_views || 0,
    total_tags: statsRes.total_tags || 0
  }
  tags.value = Array.isArray(tagsRes) ? tagsRes : (tagsRes.tags || [])
  await fetchPosts()
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
.filter-bar {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 24px;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.filter-btn {
  font-size: 13px;
  padding: 6px 14px;
  border-radius: 100px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--muted);
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: var(--font-sans);
}
.filter-btn:hover {
  border-color: var(--fg);
  color: var(--fg);
  transform: translateY(-1px);
}
.filter-btn:active {
  transform: scale(0.96);
}
.filter-btn.active {
  background: var(--fg);
  color: var(--bg);
  border-color: var(--fg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
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
  .filter-bar {
    padding: 0 16px 16px;
  }
  .post-list {
    padding: 0 16px 20px;
  }
  .pagination-wrap {
    padding: 0 16px 48px;
  }
}
</style>
