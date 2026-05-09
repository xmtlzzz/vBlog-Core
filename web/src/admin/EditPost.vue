<template>
  <div class="edit-post-page">
    <!-- Top bar -->
    <div class="editor-topbar">
      <button class="back-btn" @click="router.push('/admin/posts')">← 返回文章列表</button>
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
        style="height: 600px"
        :preview="true"
        :htmlPreview="true"
        @onUploadImg="onUploadImg"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
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
const allTags = ref([])
const form = reactive({ title: '', content: '', excerpt: '', status: 'draft', tagNames: [] })

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
    form.status = post.status || 'draft'
    form.tagNames = (post.tags || []).map(t => t.name)
  } catch {
    ElMessage.error('文章加载失败')
    router.push('/admin/posts')
  }
}

async function handleSave(status) {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: form.title,
      content: form.content,
      excerpt: form.excerpt,
      status,
      tags: form.tagNames.map(name => ({ name }))
    }
    if (isEdit.value) {
      await api.put(`/posts/${postId.value}`, payload)
      ElMessage.success('文章已更新')
    } else {
      await api.post('/posts', payload)
      ElMessage.success('文章已创建')
      router.push('/admin/posts')
    }
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function onUploadImg(files, callback) {
  const urls = []
  for (const file of files) {
    const formData = new FormData()
    formData.append('file', file)
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
}
</style>
