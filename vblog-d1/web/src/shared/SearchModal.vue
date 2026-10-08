<template>
  <Teleport to="body">
    <Transition name="search-modal">
      <div v-if="isOpen" class="search-modal-backdrop" @click="close">
        <div class="search-modal-panel" @click.stop>
          <div class="search-modal-header">
            <svg class="search-modal-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
            </svg>
            <input
              ref="inputRef"
              v-model="query"
              class="search-modal-input"
              placeholder="搜索文章标题、内容或标签…"
              aria-label="搜索文章"
              @input="onInput"
              @keydown.down.prevent="selectNext"
              @keydown.up.prevent="selectPrev"
              @keydown.enter.prevent="selectCurrent"
              @keydown.escape="close"
            />
            <button class="search-modal-close" @click="close" aria-label="关闭搜索" title="关闭 (Esc)">
              ESC
            </button>
          </div>

          <div class="search-modal-body" ref="listRef">
            <div v-if="loading" class="search-modal-status">
              正在搜索…
            </div>
            <div v-else-if="results.length > 0" class="search-modal-list" role="listbox">
              <div
                v-for="(item, idx) in results"
                :key="item.id"
                role="option"
                :aria-selected="idx === selectedIndex"
                :class="['search-modal-item', { 'is-active': idx === selectedIndex }]"
                @mouseenter="selectedIndex = idx"
                @click="goToPost(item.id)"
              >
                <div class="search-item-header">
                  <span class="search-item-title">{{ item.title }}</span>
                  <span v-if="item.tags && item.tags.length" class="search-item-tag">
                    {{ item.tags[0].name || item.tags[0] }}
                  </span>
                </div>
                <div v-if="item.excerpt" class="search-item-excerpt">
                  {{ cleanExcerpt(item.excerpt) }}
                </div>
              </div>
            </div>
            <div v-else-if="query.trim()" class="search-modal-status">
              未找到与 “{{ query }}” 相关的文章
            </div>
            <div v-else class="search-modal-status search-modal-hint">
              输入关键词，在全站文章中快速查找
            </div>
          </div>

          <div class="search-modal-footer">
            <div class="search-modal-keys">
              <span><kbd>↑</kbd><kbd>↓</kbd> 切换</span>
              <span><kbd>↵</kbd> 选择</span>
              <span><kbd>ESC</kbd> 退出</span>
            </div>
            <span class="search-modal-count" v-if="results.length > 0">
              找到 {{ results.length }} 篇结果
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api/request'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'close'])

const router = useRouter()
const isOpen = ref(false)
const query = ref('')
const results = ref([])
const loading = ref(false)
const selectedIndex = ref(0)
const inputRef = ref(null)
const listRef = ref(null)
let searchTimer = null

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val) {
    query.value = ''
    results.value = []
    selectedIndex.value = 0
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
})

function close() {
  emit('update:modelValue', false)
  emit('close')
}

function cleanExcerpt(text) {
  if (!text) return ''
  return text.replace(/^[#\s*`>_-]+/gm, '').replace(/\[([^\]]+)\]\([^)]+\)/g, '$1').trim()
}

function onInput() {
  clearTimeout(searchTimer)
  const q = query.value.trim()
  if (!q) {
    results.value = []
    loading.value = false
    return
  }
  loading.value = true
  searchTimer = setTimeout(async () => {
    try {
      const res = await api.get('/posts', {
        params: { page: 1, per_page: 20, search: q, status: 'published' }
      })
      results.value = res.data || []
      selectedIndex.value = 0
    } catch {
      results.value = []
    } finally {
      loading.value = false
    }
  }, 220)
}

function selectNext() {
  if (results.value.length === 0) return
  selectedIndex.value = (selectedIndex.value + 1) % results.value.length
  scrollToActive()
}

function selectPrev() {
  if (results.value.length === 0) return
  selectedIndex.value = (selectedIndex.value - 1 + results.value.length) % results.value.length
  scrollToActive()
}

function selectCurrent() {
  if (results.value.length > 0 && results.value[selectedIndex.value]) {
    goToPost(results.value[selectedIndex.value].id)
  }
}

function scrollToActive() {
  nextTick(() => {
    const list = listRef.value
    if (!list) return
    const activeItem = list.querySelector('.search-modal-item.is-active')
    if (activeItem) {
      activeItem.scrollIntoView({ block: 'nearest' })
    }
  })
}

function goToPost(id) {
  close()
  router.push(`/post/${id}`)
}

function onGlobalKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault()
    if (!isOpen.value) {
      emit('update:modelValue', true)
    } else {
      close()
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', onGlobalKeydown)
})

onUnmounted(() => {
  clearTimeout(searchTimer)
  window.removeEventListener('keydown', onGlobalKeydown)
})
</script>

<style scoped>
.search-modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 2000;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 80px 16px 24px;
}

.search-modal-panel {
  width: 100%;
  max-width: 620px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg, 12px);
  box-shadow: 0 20px 48px -12px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.search-modal-header {
  display: flex;
  align-items: center;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  gap: 12px;
}

.search-modal-icon {
  color: var(--muted);
  flex-shrink: 0;
}

.search-modal-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  font-size: 16px;
  color: var(--fg);
  font-family: var(--font-sans);
}

.search-modal-input::placeholder {
  color: var(--muted);
}

.search-modal-close {
  background: var(--surface-hover, var(--border));
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--muted);
  font-size: 11px;
  font-family: var(--font-mono);
  padding: 2px 6px;
  cursor: pointer;
  transition: all 0.15s;
}

.search-modal-close:hover {
  color: var(--fg);
  border-color: var(--fg);
}

.search-modal-body {
  max-height: 420px;
  overflow-y: auto;
  overscroll-behavior: contain;
  padding: 8px;
}

.search-modal-status {
  padding: 36px 16px;
  text-align: center;
  color: var(--muted);
  font-size: 14px;
}

.search-modal-hint {
  font-size: 13px;
}

.search-modal-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-modal-item {
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.search-modal-item.is-active {
  background: var(--accent-soft, rgba(37, 99, 235, 0.08));
}

.search-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.search-item-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--fg);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.search-modal-item.is-active .search-item-title {
  color: var(--accent);
}

.search-item-tag {
  font-size: 11px;
  font-weight: 500;
  color: var(--tag-fg);
  background: var(--tag-bg);
  padding: 2px 6px;
  border-radius: 4px;
  flex-shrink: 0;
}

.search-item-excerpt {
  font-size: 12px;
  color: var(--muted);
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.search-modal-footer {
  padding: 10px 18px;
  border-top: 1px solid var(--border);
  background: var(--surface);
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--muted);
}

.search-modal-keys {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-modal-keys kbd {
  font-family: var(--font-mono);
  background: var(--border);
  color: var(--muted);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 11px;
  margin-right: 2px;
}

.search-modal-count {
  font-size: 11px;
  font-family: var(--font-mono);
}

/* Animations */
.search-modal-enter-active,
.search-modal-leave-active {
  transition: opacity 0.18s ease;
}

.search-modal-enter-from,
.search-modal-leave-to {
  opacity: 0;
}

.search-modal-enter-active .search-modal-panel,
.search-modal-leave-active .search-modal-panel {
  transition: transform 0.18s ease, opacity 0.18s ease;
}

.search-modal-enter-from .search-modal-panel {
  transform: translateY(-16px) scale(0.98);
  opacity: 0;
}

.search-modal-leave-to .search-modal-panel {
  transform: translateY(-8px) scale(0.99);
  opacity: 0;
}
</style>
