<template>
  <iframe ref="frame" :srcdoc="source" sandbox="allow-scripts" @load="syncTheme"></iframe>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { buildSrcdoc } from '../utils/component'
import { useThemeStore } from '../stores/theme'

const props = defineProps({ code: String, data: Object, padding: { type: String, default: '12px' }, emptyMessage: String })
const themeStore = useThemeStore()
const initialTheme = themeStore.theme
const frame = ref(null)
const source = computed(() => {
  const code = props.code || (props.emptyMessage ? `function render(root) { root.textContent = ${JSON.stringify(props.emptyMessage)} }` : '')
  return buildSrcdoc(code, { data: props.data, padding: props.padding, theme: initialTheme })
})
function syncTheme() {
  frame.value?.contentWindow?.postMessage({ type: 'vblog-theme', theme: themeStore.theme }, '*')
}
watch(() => themeStore.theme, syncTheme)
</script>
