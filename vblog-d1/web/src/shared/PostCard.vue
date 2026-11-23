<template>
  <article :class="['post-card fade-in', { pinned: post.pinned || post.is_pinned }]">
    <div class="post-card-body">
      <div v-if="post.pinned || post.is_pinned" class="pin-badge">置顶 Pinned</div>
      <div class="post-meta">
        <router-link
          v-for="tag in (post.tags || [])"
          :key="tag.id || tag.name || tag"
          :to="`/tags?tag=${encodeURIComponent(tag.name || tag)}`"
          class="tag"
        >{{ tag.name || tag }}</router-link>
        <span class="meta-date">{{ formatDate(post.created_at) }}</span>
      </div>
      <h2 class="post-title">
        <router-link :to="`/post/${post.id}`" class="post-title-link">
          {{ displayTitle }}
        </router-link>
      </h2>
      <div class="post-excerpt">{{ plainExcerpt(post.excerpt || post.content) }}</div>
    </div>
    <div class="post-stats">
      <span class="read-time">{{ displayReadTime }} min</span>
      <span class="views">{{ (post.views || 0).toLocaleString() }} views</span>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { formatDate } from '../utils/format'
import { cleanTitle, plainExcerpt } from '../utils/markdown'

const props = defineProps({
  post: { type: Object, required: true }
})

const displayTitle = computed(() => cleanTitle(props.post.title))

const displayReadTime = computed(() => {
  if (props.post.read_time) return props.post.read_time
  const text = props.post.content || props.post.excerpt || ''
  const words = text.replace(/\s+/g, '').length
  return Math.max(1, Math.ceil(words / 350))
})
</script>

<style scoped>
.post-card {
  position: relative;
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
.post-card:hover .post-title,
.post-card:hover .post-title-link {
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
  position: relative;
  z-index: 2;
  display: inline-block;
  background: var(--tag-bg);
  color: var(--tag-fg);
  padding: 2px 8px;
  border-radius: var(--radius-sm, 6px);
  font-size: 12px;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.15s ease;
}
.tag:hover {
  background: var(--accent-soft);
  color: var(--accent);
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
  margin: 0 0 6px;
}
.post-title-link {
  color: inherit;
  text-decoration: none;
}
.post-title-link::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
}
.post-title-link:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
  border-radius: 4px;
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
