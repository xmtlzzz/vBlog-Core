<template>
  <router-view />
</template>
<script setup>
import { useThemeStore } from './stores/theme'
const theme = useThemeStore()
theme.init()
</script>
