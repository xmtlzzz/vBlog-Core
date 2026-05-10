<template>
  <section v-if="widgets.length" class="custom-widgets">
    <div v-for="w in widgets" :key="w.id" class="widget-slot">
      <iframe
        :srcdoc="srcdoc(w.code)"
        class="widget-iframe"
        sandbox="allow-scripts"
        :title="w.name"
      ></iframe>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api/request'
import { buildSrcdoc } from '../utils/component'

const widgets = ref([])
const blogData = ref(null)

function srcdoc(code) {
  return buildSrcdoc(code, { data: blogData.value })
}

onMounted(async () => {
  try {
    const [compRes, statsRes] = await Promise.all([
      api.get('/components/active'),
      api.get('/dashboard/stats').catch(() => null)
    ])
    widgets.value = Array.isArray(compRes) ? compRes : (compRes.data || [])
    if (statsRes) {
      var s = statsRes.data || statsRes
      blogData.value = {
        posts: s.total_posts || 0,
        views: s.total_views || 0,
        comments: s.total_comments || 0,
        tags: s.total_tags || 0,
        visitors: Math.floor(Math.random() * 9000) + 1000
      }
    }
  } catch {}
})
</script>

<style scoped>
.custom-widgets {
  max-width: 1080px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}
.widget-slot {
  flex: 1 1 280px;
  max-width: 100%;
}
.widget-iframe {
  width: 100%;
  min-height: 60px;
  border: none;
  background: transparent;
}
</style>
