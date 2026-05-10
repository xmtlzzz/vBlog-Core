<template>
  <div class="settings-overlay" @click.self="$emit('close')">
    <div class="settings-panel">
      <h2>连接设置</h2>
      <label>
        服务器地址
        <input v-model="addr" placeholder="localhost:50051" />
      </label>
      <label>
        API Key
        <input v-model="apiKey" type="password" placeholder="输入 API Key" />
      </label>
      <div class="settings-actions">
        <button @click="$emit('close')">取消</button>
        <button class="primary" @click="connect">连接</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['close', 'connect'])
const addr = ref(localStorage.getItem('vblog_addr') || 'localhost:50051')
const apiKey = ref(localStorage.getItem('vblog_key') || '')

function connect() {
  localStorage.setItem('vblog_addr', addr.value)
  localStorage.setItem('vblog_key', apiKey.value)
  emit('connect', { addr: addr.value, apiKey: apiKey.value })
}
</script>

<style scoped>
.settings-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); display: flex; align-items: center; justify-content: center; z-index: 100; }
.settings-panel { background: white; padding: 24px; border-radius: 12px; width: 320px; }
.settings-panel h2 { margin: 0 0 16px; }
.settings-panel label { display: block; margin-bottom: 12px; font-size: 13px; }
.settings-panel input { display: block; width: 100%; margin-top: 4px; padding: 8px; border: 1px solid #e5e5e5; border-radius: 6px; box-sizing: border-box; }
.settings-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 16px; }
.settings-actions button { padding: 8px 16px; border-radius: 6px; border: 1px solid #e5e5e5; cursor: pointer; }
.settings-actions .primary { background: #2563eb; color: white; border-color: #2563eb; }
</style>
