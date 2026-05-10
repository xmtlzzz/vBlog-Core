<template>
  <div class="change-card">
    <div class="change-body">
      <div class="change-title">{{ change.title }}</div>
      <div class="change-meta">
        <span class="change-type">{{ change.type }}</span>
        <span class="change-time">{{ change.timestamp }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({ change: Object })
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
