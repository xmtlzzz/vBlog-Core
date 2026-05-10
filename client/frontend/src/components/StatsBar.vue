<template>
  <div class="stats-bar" v-if="stats">
    <div class="stat">
      <span class="stat-val">{{ stats.pv_today || 0 }}</span>
      <span class="stat-label">今日 PV</span>
      <span class="stat-delta" v-if="pvDelta !== null">{{ pvDelta > 0 ? '+' : '' }}{{ pvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.uv_today || 0 }}</span>
      <span class="stat-label">今日 UV</span>
      <span class="stat-delta" v-if="uvDelta !== null">{{ uvDelta > 0 ? '+' : '' }}{{ uvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.total_views || 0 }}</span>
      <span class="stat-label">总阅读量</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.total_posts || 0 }}</span>
      <span class="stat-label">文章数</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ stats: Object })

const pvDelta = computed(() => {
  if (!props.stats?.pv_yesterday) return null
  return Math.round(((props.stats.pv_today - props.stats.pv_yesterday) / props.stats.pv_yesterday) * 100)
})
const uvDelta = computed(() => {
  if (!props.stats?.uv_yesterday) return null
  return Math.round(((props.stats.uv_today - props.stats.uv_yesterday) / props.stats.uv_yesterday) * 100)
})
</script>

<style scoped>
.stats-bar {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  padding: 16px;
  background: var(--surface, #fff);
  border: 1px solid var(--border, #e5e5e5);
  border-radius: 8px;
}
.stat { text-align: center; }
.stat-val { display: block; font-size: 24px; font-weight: 700; }
.stat-label { display: block; font-size: 11px; color: #737373; text-transform: uppercase; }
.stat-delta { display: block; font-size: 12px; margin-top: 2px; }
.stat-delta { color: #16a34a; }
</style>
