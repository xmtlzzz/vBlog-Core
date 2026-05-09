<template>
  <BlogNav />
  <div class="page-enter">
  <header class="page-header">
    <h1>归档 Archives</h1>
    <p>按时间线浏览所有文章。</p>
    <div class="archive-count">共 {{ allPosts.length }} 篇文章</div>
  </header>

  <main class="archive-list">
    <div v-for="group in yearGroups" :key="group.year" class="year-group">
      <div class="year-label">
        {{ group.year }} <span class="year-count">— {{ group.posts.length }} 篇</span>
      </div>
      <div class="timeline">
        <router-link
          v-for="post in group.posts"
          :key="post.id"
          :to="`/post/${post.id}`"
          class="archive-item"
        >
          <div class="timeline-dot"></div>
          <span class="archive-date">{{ formatDay(post.created_at) }}</span>
          <div class="archive-content">
            <div class="archive-title">{{ post.title }}</div>
            <div class="archive-excerpt">{{ post.excerpt }}</div>
            <div class="archive-tags">
              <span
                v-for="tag in (post.tags || [])"
                :key="tag"
                class="archive-tag"
              >{{ tag }}</span>
            </div>
          </div>
          <div class="archive-stats">
            <span>{{ (post.views || 0).toLocaleString() }} views</span>
            <span>{{ post.read_time || 0 }} min</span>
          </div>
        </router-link>
      </div>
    </div>

    <div v-if="allPosts.length === 0" class="empty-state">
      <p>暂无文章</p>
    </div>
  </main>

  </div>
  <BlogFooter />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'

const allPosts = ref([])

const yearGroups = computed(() => {
  const groups = {}
  for (const post of allPosts.value) {
    const year = new Date(post.created_at).getFullYear()
    if (!groups[year]) groups[year] = []
    groups[year].push(post)
  }
  return Object.keys(groups)
    .sort((a, b) => b - a)
    .map(year => ({ year: Number(year), posts: groups[year] }))
})

function formatDay(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${mm}-${dd}`
}

onMounted(async () => {
  const res = await api.get('/posts', { params: { per_page: 100, status: 'published' } }).catch(() => ({ posts: [] }))
  allPosts.value = res.posts || []
})
</script>

<style scoped>
.page-header {
  max-width: 1080px;
  margin: 0 auto;
  padding: 56px 24px 32px;
}
.page-header h1 {
  font-family: var(--font-display);
  font-size: clamp(28px, 4vw, 40px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin-bottom: 8px;
}
.page-header p {
  font-size: 16px;
  color: var(--muted);
}
.archive-count {
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
  font-size: 13px;
  color: var(--muted);
  margin-top: 12px;
}
.archive-list {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 80px;
}
.year-group {
  margin-bottom: 48px;
}
.year-label {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--fg);
  margin-bottom: 20px;
  display: flex;
  align-items: baseline;
  gap: 10px;
}
.year-count {
  font-family: var(--font-mono);
  font-size: 13px;
  font-variant-numeric: tabular-nums;
  color: var(--muted);
  font-weight: 400;
}
.timeline {
  border-left: 2px solid var(--border);
  padding-left: 24px;
  margin-left: 4px;
}
.archive-item {
  display: grid;
  grid-template-columns: 80px 1fr auto;
  gap: 16px;
  align-items: start;
  padding: 16px 0;
  text-decoration: none;
  color: inherit;
  transition: opacity 0.15s;
  position: relative;
}
.archive-item:hover {
  opacity: 0.7;
}
.timeline-dot {
  position: absolute;
  left: -29px;
  top: 20px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--border);
  border: 2px solid var(--bg);
}
.archive-item:first-child .timeline-dot {
  background: var(--accent);
}
.archive-date {
  font-family: var(--font-mono);
  font-size: 13px;
  font-variant-numeric: tabular-nums;
  color: var(--muted);
  padding-top: 2px;
}
.archive-title {
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 600;
  letter-spacing: -0.01em;
  line-height: 1.35;
  color: var(--fg);
  margin-bottom: 4px;
}
.archive-excerpt {
  font-size: 13px;
  color: var(--muted);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.archive-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  padding-top: 4px;
}
.archive-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  background: var(--tag-bg, #f0f0f0);
  color: var(--tag-fg, #525252);
  font-weight: 500;
}
.archive-stats {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  padding-top: 2px;
  white-space: nowrap;
}
.archive-stats span {
  font-family: var(--font-mono);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  color: var(--muted);
}
.empty-state {
  text-align: center;
  padding: 64px 24px;
  color: var(--muted);
}

@media (max-width: 640px) {
  .page-header {
    padding: 40px 16px 24px;
  }
  .archive-list {
    padding: 0 16px 48px;
  }
  .archive-item {
    grid-template-columns: 60px 1fr;
  }
  .archive-stats {
    display: none;
  }
}
</style>
