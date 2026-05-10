<template>
  <div class="trend-panel">
    <div class="trend-header">
      <h3>Trends</h3>
      <div class="trend-tabs">
        <button v-for="g in ['day','week','month']" :key="g"
          :class="{ active: granularity === g }"
          @click="$emit('change', g)">{{ g }}</button>
      </div>
    </div>
    <div class="trend-chart" v-if="points?.length">
      <div v-for="(p, i) in points" :key="i" class="trend-bar-group">
        <div class="trend-bar" :style="{ height: barHeight(p.pv) + '%' }"></div>
        <span class="trend-label">{{ p.label?.slice(5) || '' }}</span>
      </div>
    </div>
    <div v-else class="trend-empty">No trend data</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ points: Array, granularity: String })

const maxPv = computed(() => Math.max(...(props.points?.map(p => p.pv) || [1])))
const barHeight = (pv) => Math.round((pv / maxPv.value) * 100)
</script>

<style scoped>
.trend-panel { padding: 16px; background: var(--surface, #fff); border: 1px solid var(--border, #e5e5e5); border-radius: 8px; }
.trend-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.trend-header h3 { margin: 0; font-size: 14px; }
.trend-tabs button { background: none; border: 1px solid #e5e5e5; padding: 4px 10px; border-radius: 4px; cursor: pointer; font-size: 12px; }
.trend-tabs button.active { background: #2563eb; color: white; border-color: #2563eb; }
.trend-chart { display: flex; gap: 4px; height: 100px; align-items: flex-end; }
.trend-bar-group { flex: 1; display: flex; flex-direction: column; align-items: center; }
.trend-bar { width: 100%; background: #2563eb; border-radius: 2px 2px 0 0; min-height: 2px; }
.trend-label { font-size: 10px; color: #737373; margin-top: 4px; }
.trend-empty { text-align: center; color: #737373; padding: 20px; }
</style>
