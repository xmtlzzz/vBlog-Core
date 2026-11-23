<template>
  <BlogNav />
  <main id="main-content" class="page-enter">
  <header class="page-header">
    <h1>标签 Tags</h1>
    <p>按标签浏览文章。</p>
  </header>

  <section class="tags-cloud">
    <p v-if="loadingTags" class="tags-status" role="status">正在加载标签…</p>
    <div v-else-if="tagsError" class="tags-status" role="alert">
      <span>标签加载失败，请稍后重试。</span>
      <button type="button" class="retry-btn" @click="fetchTags">重新加载</button>
    </div>
    <p v-else-if="tags.length === 0" class="tags-status" role="status">暂无标签</p>
    <button
      v-else
      v-for="tag in tags"
      :key="tag.id || tag.name"
      type="button"
      :aria-pressed="activeTag === tag.name"
      :class="['tag-chip', { active: activeTag === tag.name }]"
      @click="selectTag(tag.name)"
    >
      {{ tag.name }}
      <span class="count">{{ tag.post_count || 0 }}</span>
    </button>
  </section>

  <section class="tag-posts" v-if="activeTag">
    <div class="tag-posts-header">
      <span class="tag-posts-title">{{ activeTag }}</span>
      <span class="tag-posts-count" v-if="!loadingPosts">{{ filteredPosts.length }} 篇</span>
    </div>
    <div v-if="loadingPosts" class="empty-state" role="status">
      <p>正在加载文章…</p>
    </div>
    <template v-else>
      <PostCard v-for="post in filteredPosts" :key="post.id" :post="post" />
      <div v-if="filteredPosts.length === 0" class="empty-state">
        <p>该标签下暂无文章</p>
      </div>
    </template>
  </section>

  </main>
  <CustomWidgets />
  <BlogFooter />
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api/request'
import { updateMetadata } from '../utils/metadata'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'
import PostCard from '../shared/PostCard.vue'

const route = useRoute()
const router = useRouter()

const tags = ref([])
const loadingTags = ref(true)
const tagsError = ref(false)
const activeTag = ref('')
const filteredPosts = ref([])
const loadingPosts = ref(false)

function updateTagMeta() {
  const tagTitle = activeTag.value ? `标签: ${activeTag.value}` : '标签'
  updateMetadata({}, { title: tagTitle, description: '按主题标签浏览所有文章。', path: '/tags' })
}

async function fetchTags() {
  loadingTags.value = true
  tagsError.value = false
  try {
    const res = await api.get('/tags')
    tags.value = Array.isArray(res) ? res : (res?.data || [])
  } catch {
    tagsError.value = true
    tags.value = []
  } finally {
    loadingTags.value = false
  }
}

async function fetchPosts() {
  if (!activeTag.value) {
    filteredPosts.value = []
    return
  }
  loadingPosts.value = true
  try {
    const res = await api.get('/posts', {
      params: { per_page: 100, tag: activeTag.value, status: 'published' }
    })
    filteredPosts.value = Array.isArray(res) ? res : (res?.data || [])
  } catch {
    filteredPosts.value = []
  } finally {
    loadingPosts.value = false
  }
}

function selectTag(tagName) {
  if (activeTag.value === tagName) {
    activeTag.value = ''
    filteredPosts.value = []
    router.replace({ query: {} })
  } else {
    activeTag.value = tagName
    router.replace({ query: { tag: tagName } })
    fetchPosts()
  }
  updateTagMeta()
}

watch(
  () => route.query.tag,
  (newTag) => {
    if (newTag && newTag !== activeTag.value) {
      activeTag.value = String(newTag)
      fetchPosts()
    } else if (!newTag && activeTag.value) {
      activeTag.value = ''
      filteredPosts.value = []
    }
    updateTagMeta()
  }
)

onMounted(async () => {
  if (route.query.tag) {
    activeTag.value = String(route.query.tag)
  }
  updateTagMeta()
  await fetchTags()
  if (activeTag.value) {
    fetchPosts()
  }
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
.tags-status {
  color: var(--muted);
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}
.retry-btn {
  padding: 4px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  color: var(--fg);
  font-size: 12px;
  cursor: pointer;
}
.retry-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
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
