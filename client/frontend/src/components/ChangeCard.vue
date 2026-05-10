<template>
  <div class="change-card">
    <div class="change-body">
      <div class="change-title">{{ change.title }}</div>
      <div class="change-meta">
        <span class="change-type">{{ typeLabel(change.type) }}</span>
        <span class="change-time">{{ formatTime(change.timestamp) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({ change: Object })

const typeMap = {
  new_post: '新文章',
  new_comment: '新评论',
  view_milestone: '阅读里程碑',
  pv_milestone: '访问里程碑',
  tag_added: '新标签',
}
function typeLabel(type) {
  return typeMap[type] || type
}
function formatTime(ts) {
  if (!ts) return ''
  const d = new Date(ts)
  if (isNaN(d.getTime())) return ts
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<style scoped>
.change-card {
  padding: 12px;
  border: 1px solid var(--border, #e5e5e5);
  border-radius: 8px;
  background: var(--surface, #fff);
}
.change-title { font-weight: 500; font-size: 14px; }
.change-meta { display: flex; gap: 8px; margin-top: 4px; font-size: 12px; color: #737373; }
.change-type { background: #eff6ff; color: #2563eb; padding: 1px 6px; border-radius: 4px; }
</style>
