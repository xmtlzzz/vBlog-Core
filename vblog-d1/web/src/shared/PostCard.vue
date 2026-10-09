<template>
  <router-link
    :to="`/post/${post.id}`"
    :class="['post-card fade-in', { pinned: post.pinned || post.is_pinned }]"
  >
    <div class="post-card-body">
      <div v-if="post.pinned || post.is_pinned" class="pin-badge">置顶 Pinned</div>
      <div class="post-meta">
        <span
          v-for="tag in (post.tags || [])"
          :key="tag.id || tag.name || tag"
          class="tag"
        >{{ tag.name || tag }}</span>
        <span class="meta-date">{{ formatDate(post.created_at) }}</span>
      </div>
      <div class="post-title">{{ post.title }}</div>
      <div class="post-excerpt">{{ plainExcerpt(post.excerpt || post.content) }}</div>
    </div>
    <div class="post-stats">
      <span class="read-time">{{ post.read_time || 0 }} min</span>
      <span class="views">{{ (post.views || 0).toLocaleString() }} views</span>
    </div>
  </router-link>
</template>

<script setup>
import { formatDate } from '../utils/format'
import { plainExcerpt } from '../utils/markdown'

defineProps({
  post: { type: Object, required: true }
})
</script>

<style scoped>
.post-card {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px;
  align-items: start;
  padding: 20px 0;
  border-bottom: 1px solid var(--border);
  text-decoration: none;
  color: var(--fg);
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.2s ease, border-color 0.2s ease;
}
.post-card:hover {
  transform: translateY(-2px);
  border-color: var(--accent);
}
.post-card-body { min-width: 0; overflow-wrap: anywhere; }
.post-card:hover .post-title {
  color: var(--accent);
}
.post-card:last-child {
  border-bottom: none;
}
.post-card.pinned {
  background: var(--accent-soft);
  border-radius: var(--radius-lg);
  padding: 20px;
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
.tag {
  display: inline-block;
  background: var(--tag-bg);
  color: var(--tag-fg);
  padding: 2px 8px;
  border-radius: var(--radius-sm, 6px);
  font-size: 12px;
  font-weight: 500;
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
}
.post-stats {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--muted);
  white-space: nowrap;
}

@media (max-width: 640px) {
  .post-card {
    padding: 16px 0;
    grid-template-columns: minmax(0, 1fr);
    gap: 10px;
  }
  .post-stats { flex-direction: row; align-items: center; gap: 16px; }
  .post-card.pinned {
    padding: 16px;
  }
  .post-title {
    font-size: 17px;
  }
}
</style>
