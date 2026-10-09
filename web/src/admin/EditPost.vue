<template>
  <div class="edit-post-page">
    <!-- Top bar -->
    <div class="editor-topbar">
      <div class="topbar-left">
        <button class="back-btn" @click="router.push('/admin/posts')">← 返回文章列表</button>
        <span v-if="draftSavedTime" class="draft-indicator" :title="'草稿已保存在本地浏览器中'">
          <span class="draft-dot"></span>已自动保存草稿 {{ draftSavedTime }}
        </span>
      </div>
      <div class="editor-actions">
        <el-button :loading="saving" @click="handleSave('draft')">保存草稿</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave('published')">发布</el-button>
      </div>
    </div>

    <!-- Metadata: WordPress-style 2-column -->
    <div class="meta-grid">
      <div class="meta-item meta-full">
        <label>标题 Title</label>
        <el-input v-model="form.title" placeholder="文章标题" size="large" />
      </div>
      <div class="meta-item">
        <label>状态 Status</label>
        <el-select v-model="form.status" style="width: 100%">
          <el-option label="草稿 Draft" value="draft" />
          <el-option label="已发布 Published" value="published" />
          <el-option label="已归档 Archived" value="archived" />
        </el-select>
      </div>
      <div class="meta-item">
        <label>标签 Tags</label>
        <el-select v-model="form.tagNames" multiple filterable allow-create placeholder="选择或输入标签" style="width: 100%">
          <el-option v-for="t in allTags" :key="t.id" :label="t.name" :value="t.name" />
        </el-select>
      </div>
      <div class="meta-item meta-full">
        <label>摘要 Excerpt</label>
        <el-input v-model="form.excerpt" type="textarea" :rows="2" placeholder="文章摘要（可选，留空自动生成）" />
      </div>
    </div>

    <!-- Markdown Editor -->
    <div class="editor-wrap">
      <MdEditor
        v-model="form.content"
        :theme="editorTheme"
        language="zh-CN"
        class="vblog-md-editor"
        :preview="true"
        :htmlPreview="true"
        @onUploadImg="onUploadImg"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import api from '../api/request'
import { useThemeStore } from '../stores/theme'

const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()

const isEdit = computed(() => !!route.params.id)
const postId = computed(() => route.params.id)
const editorTheme = computed(() => themeStore.theme === 'dark' ? 'dark' : 'light')

const saving = ref(false)
const saved = ref(false)
const allTags = ref([])
const form = reactive({ title: '', content: '', excerpt: '', status: 'published', tagNames: [] })

// ── 本地草稿防丢（Auto-save draft）──
const draftKey = computed(() => 'vblog_draft_' + (postId.value || 'new'))
const draftSavedTime = ref('')
let autoSaveTimer = null
let isRestoring = false

const isDirty = computed(() => form.title.trim() || form.content.trim() || form.excerpt.trim())

watch(
  () => [form.title, form.content, form.excerpt, form.tagNames],
  () => {
    if (saved.value || isRestoring) return
    if (!form.title.trim() && !form.content.trim()) return
    clearTimeout(autoSaveTimer)
    autoSaveTimer = setTimeout(() => {
      try {
        const data = {
          title: form.title,
          content: form.content,
          excerpt: form.excerpt,
          tagNames: [...form.tagNames],
          status: form.status,
          savedAt: Date.now()
        }
        localStorage.setItem(draftKey.value, JSON.stringify(data))
        const d = new Date()
        draftSavedTime.value = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:${String(d.getSeconds()).padStart(2, '0')}`
      } catch {
        // quota exceeded or storage disabled
      }
    }, 1500)
  },
  { deep: true }
)

function checkLocalDraft() {
  try {
    const raw = localStorage.getItem(draftKey.value)
    if (!raw) return
    const draft = JSON.parse(raw)
    if (!draft || (!draft.title && !draft.content)) return

    // 如果与当前内容一致则不提示
    const isDifferent = draft.title !== form.title || draft.content !== form.content || draft.excerpt !== form.excerpt
    if (!isDifferent) return

    const timeStr = draft.savedAt ? new Date(draft.savedAt).toLocaleTimeString() : '之前'
    ElMessageBox.confirm(`检测到本地浏览器保存有未发布的草稿（${timeStr}），是否恢复？`, '恢复草稿', {
      confirmButtonText: '恢复草稿',
      cancelButtonText: '放弃',
      type: 'info'
    }).then(() => {
      isRestoring = true
      form.title = draft.title || ''
      form.content = draft.content || ''
      form.excerpt = draft.excerpt || ''
      if (Array.isArray(draft.tagNames)) form.tagNames = draft.tagNames
      if (draft.status) form.status = draft.status
      ElMessage.success('已恢复本地草稿')
      setTimeout(() => { isRestoring = false }, 500)
    }).catch(() => {
      // 用户选择放弃，不覆盖
    })
  } catch {
    // parse error
  }
}

onBeforeRouteLeave(async () => {
  if (saved.value || !isDirty.value) return true
  try {
    await ElMessageBox.confirm('当前内容未保存，确定离开？', '未保存的更改', {
      confirmButtonText: '离开',
      cancelButtonText: '留下',
      type: 'warning'
    })
    return true
  } catch {
    return false
  }
})

function onBeforeUnload(e) {
  if (!isDirty.value) return
  e.preventDefault()
  e.returnValue = ''
}
window.addEventListener('beforeunload', onBeforeUnload)
onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', onBeforeUnload)
  clearTimeout(autoSaveTimer)
})

async function fetchTags() {
  const res = await api.get('/tags').catch(() => [])
  allTags.value = Array.isArray(res) ? res : (res.data || [])
}

async function fetchPost() {
  if (!isEdit.value) return
  try {
    const res = await api.get(`/posts/${postId.value}`)
    const post = res.data || res
    form.title = post.title || ''
    form.content = post.content || ''
    form.excerpt = post.excerpt || ''
    form.status = post.status || 'published'
    form.tagNames = (post.tags || []).map(t => t.name)
    // 加载完已有文章后检查本地草稿
    checkLocalDraft()
  } catch {
    ElMessage.error('文章加载失败')
    router.push('/admin/posts')
  }
}

async function handleSave(statusOverride) {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }
  const status = statusOverride || form.status || 'published'
  saving.value = true
  try {
    const payload = {
      title: form.title,
      content: form.content,
      excerpt: form.excerpt,
      status,
      tags: form.tagNames.map(name => ({ name }))
    }
    form.status = status
    if (isEdit.value) {
      await api.put(`/posts/${postId.value}`, payload)
      saved.value = true
      localStorage.removeItem(draftKey.value)
      draftSavedTime.value = ''
      ElMessage.success(status === 'published' ? '文章已发布更新' : '草稿已保存')
    } else {
      await api.post('/posts', payload)
      saved.value = true
      localStorage.removeItem(draftKey.value)
      draftSavedTime.value = ''
      ElMessage.success(status === 'published' ? '文章已创建并发布' : '草稿已保存')
      router.push('/admin/posts')
    }
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

// ── 图片上传前端 WebP 压缩转换 ──
async function compressImageToWebP(file) {
  // SVG, GIF（动图）或非图片文件直接返回原文件
  if (!file.type.startsWith('image/') || file.type === 'image/svg+xml' || file.type === 'image/gif') {
    return file
  }
  return new Promise((resolve) => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      URL.revokeObjectURL(url)
      let { width, height } = img
      const maxDim = 2560 // 限制最大边长不超过 2560px
      if (width > maxDim || height > maxDim) {
        if (width > height) {
          height = Math.round((height * maxDim) / width)
          width = maxDim
        } else {
          width = Math.round((width * maxDim) / height)
          height = maxDim
        }
      }
      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height
      const ctx = canvas.getContext('2d')
      ctx.drawImage(img, 0, 0, width, height)
      canvas.toBlob((blob) => {
        if (blob && blob.size < file.size) {
          const newName = file.name.replace(/\.[^.]+$/, '') + '.webp'
          const webpFile = new File([blob], newName, { type: 'image/webp' })
          resolve(webpFile)
        } else {
          resolve(file)
        }
      }, 'image/webp', 0.85)
    }
    img.onerror = () => {
      URL.revokeObjectURL(url)
      resolve(file)
    }
    img.src = url
  })
}

async function onUploadImg(files, callback) {
  const urls = []
  for (const file of files) {
    const processedFile = await compressImageToWebP(file)
    const formData = new FormData()
    formData.append('file', processedFile)
    try {
      const res = await api.post('/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      urls.push(res.url)
    } catch {
      // handled by interceptor
    }
  }
  callback(urls)
}

onMounted(() => {
  fetchTags()
  // Check for uploaded markdown content from Posts page
  const uploadedContent = sessionStorage.getItem('md-upload-content')
  const uploadedTitle = sessionStorage.getItem('md-upload-title')
  if (uploadedContent) {
    form.content = uploadedContent
    if (uploadedTitle) form.title = uploadedTitle
    sessionStorage.removeItem('md-upload-content')
    sessionStorage.removeItem('md-upload-title')
  } else if (!isEdit.value) {
    checkLocalDraft()
  } else {
    fetchPost()
  }
})
</script>

<style scoped>
.edit-post-page {
  width: 100%;
  max-width: 1200px;
}
.editor-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  gap: 16px;
}
.topbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.back-btn {
  background: none;
  border: none;
  color: var(--muted);
  font-size: 13px;
  cursor: pointer;
  padding: 4px 0;
  transition: color 0.15s;
}
.back-btn:hover {
  color: var(--fg);
}
.draft-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--muted);
  background: var(--tag-bg);
  padding: 3px 10px;
  border-radius: 999px;
}
.draft-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--success, #16a34a);
}
.editor-actions {
  display: flex;
  gap: 8px;
}
.meta-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 20px;
}
.meta-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.meta-item label {
  font-size: 12px;
  font-weight: 600;
  color: var(--muted);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.meta-full {
  grid-column: 1 / -1;
}
.editor-wrap {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  height: calc(100vh - 240px);
  min-height: 520px;
  position: relative;
  isolation: isolate;
  z-index: 1;
}
.vblog-md-editor {
  height: 100% !important;
}
.vblog-md-editor :deep(.md-editor-code-head) {
  z-index: 2 !important;
}
</style>
