<template>
  <BlogNav />
  <div class="page-enter">
  <header class="page-header">
    <h1>标签 Tags</h1>
    <p>按标签浏览文章。</p>
  </header>

  <section class="tags-cloud">
    <button
      :class="['tag-chip', { active: !activeTag }]"
      @click="activeTag = ''; filteredPosts = []"
    >全部</button>
    <button
      v-for="tag in tags"
      :key="tag.id || tag.name"
      :class="['tag-chip', { active: activeTag === tag.name }]"
      @click="activeTag = tag.name; fetchPosts()"
    >
      {{ tag.name }}
      <span class="count">{{ tag.post_count || 0 }}</span>
    </button>
  </section>

  <section class="tag-posts" v-if="activeTag">
    <div class="tag-posts-header">
      <span class="tag-posts-title">{{ activeTag }}</span>
      <span class="tag-posts-count">{{ filteredPosts.length }} 篇</span>
    </div>
    <PostCard v-for="post in filteredPosts" :key="post.id" :post="post" />
    <div v-if="filteredPosts.length === 0" class="empty-state">
      <p>该标签下暂无文章</p>
    </div>
  </section>

  </div>
  <BlogFooter />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import PostCard from '../shared/PostCard.vue'

const tags = ref([])
const activeTag = ref('')
const filteredPosts = ref([])

async function fetchPosts() {
  if (!activeTag.value) {
    filteredPosts.value = []
    return
  }
  const res = await api.get('/posts', {
    params: { per_page: 100, tag: activeTag.value, status: 'published' }
  }).catch(() => ({ posts: [] }))
  filteredPosts.value = res.posts || []
}

onMounted(async () => {
  const res = await api.get('/tags').catch(() => [])
  tags.value = Array.isArray(res) ? res : (res.tags || [])
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
.tags-cloud {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 48px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.tag-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 100px;
  border: 1px solid var(--border);
  background: var(--surface);
  font-size: 15px;
  font-weight: 500;
  color: var(--fg);
  cursor: pointer;
  transition: all 0.2s;
  font-family: var(--font-sans);
}
.tag-chip:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}
.tag-chip.active {
  background: var(--fg);
  color: var(--bg);
  border-color: var(--fg);
}
.tag-chip .count {
  font-family: var(--font-mono);
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  color: var(--muted);
  background: var(--tag-bg, #f0f0f0);
  padding: 2px 8px;
  border-radius: 100px;
}
.tag-chip:hover .count {
  background: rgba(37, 99, 235, 0.1);
  color: var(--accent);
}
.tag-chip.active .count {
  background: rgba(255, 255, 255, 0.15);
  color: inherit;
}
.tag-posts {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px 80px;
}
.tag-posts-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}
.tag-posts-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.01em;
}
.tag-posts-count {
  font-family: var(--font-mono);
  font-size: 13px;
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
  .tags-cloud {
    padding: 0 16px 32px;
    gap: 8px;
  }
  .tag-chip {
    padding: 8px 16px;
    font-size: 14px;
  }
  .tag-posts {
    padding: 0 16px 48px;
  }
}
</style>
