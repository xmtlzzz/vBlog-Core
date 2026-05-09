<template>
  <router-link
    :to="`/post/${post.id}`"
    :class="['post-card', { pinned: post.is_pinned }]"
  >
    <div class="post-card-body">
      <div v-if="post.is_pinned" class="pin-badge">置顶 Pinned</div>
      <div class="post-meta">
        <el-tag
          v-for="tag in (post.tags || [])"
          :key="tag"
          size="small"
          type="info"
          effect="plain"
          class="meta-tag"
        >{{ tag }}</el-tag>
        <span class="meta-date">{{ formatDate(post.created_at) }}</span>
      </div>
      <div class="post-title">{{ post.title }}</div>
      <div class="post-excerpt">{{ post.excerpt }}</div>
      <div class="post-stats">
        <span class="read-time">{{ post.read_time || 0 }} min</span>
        <span class="views">{{ (post.views || 0).toLocaleString() }} views</span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
defineProps({
  post: { type: Object, required: true }
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}
</script>

<style scoped>
.post-card {
  display: block;
  text-decoration: none;
  color: inherit;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
  transition: opacity 0.15s;
}
.post-card:hover {
  opacity: 0.7;
}
.post-card:last-child {
  border-bottom: none;
}
.post-card.pinned {
  background: var(--accent-soft);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 8px;
  border-bottom: none;
}
.pin-badge {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--accent);
  margin-bottom: 4px;
}
.post-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
  font-size: 13px;
  color: var(--muted);
  flex-wrap: wrap;
}
.meta-tag {
  font-size: 12px;
}
.meta-date {
  font-size: 13px;
  color: var(--muted);
}
.post-title {
  font-family: var(--font-display);
  font-size: 19px;
  font-weight: 600;
  letter-spacing: -0.01em;
  line-height: 1.35;
  color: var(--fg);
  margin-bottom: 6px;
}
.post-excerpt {
  font-size: 14px;
  color: var(--muted);
  line-height: 1.55;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 8px;
}
.post-stats {
  display: flex;
  gap: 16px;
}
.post-stats .read-time,
.post-stats .views {
  font-size: 12px;
  color: var(--muted);
  font-family: var(--font-mono);
  font-variant-numeric: tabular-nums;
}

@media (max-width: 640px) {
  .post-card {
    padding: 16px 0;
  }
  .post-card.pinned {
    padding: 16px;
  }
  .post-title {
    font-size: 17px;
  }
}
</style>
