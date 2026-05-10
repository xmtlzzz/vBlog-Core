<template>
  <div class="stats-bar" v-if="stats">
    <div class="stat">
      <span class="stat-val">{{ stats.pvToday || 0 }}</span>
      <span class="stat-label">PV Today</span>
      <span class="stat-delta" v-if="pvDelta !== null">{{ pvDelta > 0 ? '+' : '' }}{{ pvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.uvToday || 0 }}</span>
      <span class="stat-label">UV Today</span>
      <span class="stat-delta" v-if="uvDelta !== null">{{ uvDelta > 0 ? '+' : '' }}{{ uvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.totalViews || 0 }}</span>
      <span class="stat-label">Total Views</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.totalPosts || 0 }}</span>
      <span class="stat-label">Posts</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ stats: Object })

const pvDelta = computed(() => {
  if (!props.stats?.pvYesterday) return null
  return Math.round(((props.stats.pvToday - props.stats.pvYesterday) / props.stats.pvYesterday) * 100)
})
const uvDelta = computed(() => {
  if (!props.stats?.uvYesterday) return null
  return Math.round(((props.stats.uvToday - props.stats.uvYesterday) / props.stats.uvYesterday) * 100)
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
