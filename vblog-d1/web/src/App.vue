<template>
  <router-view />
</template>
<script setup>
import { onMounted } from 'vue'
import { useThemeStore } from './stores/theme'
import api from './api/request'
const theme = useThemeStore()
theme.init()

// 浏览器标签页标题跟随后台「站点标题」（index.html 里是静态兜底值）
onMounted(async () => {
  try {
    const s = await api.get('/settings')
    if (s && s.site_title) document.title = s.site_title
  } catch {}
})
</script>
