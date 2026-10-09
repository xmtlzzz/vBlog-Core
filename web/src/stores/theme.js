import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const urlTheme = typeof window !== 'undefined' ? new URLSearchParams(window.location.search).get('theme') : null
  const initialTheme = (urlTheme === 'dark' || urlTheme === 'light')
    ? urlTheme
    : (localStorage.getItem('vblog-theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
  const theme = ref(initialTheme)

  function applyTheme(t) {
    document.documentElement.setAttribute('data-theme', t)
    const isDark = t === 'dark'
    document.documentElement.classList.toggle('dark', isDark)
    document.documentElement.classList.toggle('theme-dark', isDark)
    document.documentElement.classList.toggle('theme-light', !isDark)
  }

  function toggle() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    applyTheme(theme.value)
    localStorage.setItem('vblog-theme', theme.value)
  }

  function init() { applyTheme(theme.value) }

  return { theme, toggle, init }
})
