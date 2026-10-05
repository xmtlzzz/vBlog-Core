<template>
  <router-view />
</template>
<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useThemeStore } from './stores/theme'
import api from './api/request'
import { updateMetadata } from './utils/metadata'
const theme = useThemeStore()
theme.init()

// 浏览器标签页标题跟随后台「站点标题」（index.html 里是静态兜底值）
const route = useRoute()
const settings = ref({})
const pageTitles = { '/archives': '归档', '/modules': '模块', '/tags': '标签', '/about': '关于' }
function updatePage() {
  // Article pages own their metadata once the article has loaded.
  if (!route.path.startsWith('/post/')) updateMetadata(settings.value, { title: pageTitles[route.path], path: route.path })
}
watch(() => route.path, updatePage, { immediate: true })
api.get('/settings').then(value => { settings.value = value || {}; updatePage() }).catch(() => {})
</script>
