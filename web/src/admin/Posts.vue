<template>
  <div class="posts-page">
    <div class="page-header">
      <h1 class="page-title">全部文章 All Posts</h1>
      <div class="header-actions">
        <el-button @click="triggerUpload">Markdown 上传</el-button>
        <el-button type="primary" @click="router.push('/admin/posts/new')">新文章 New Post</el-button>
        <input ref="fileInput" type="file" accept=".md,.markdown" style="display:none" @change="handleUpload" />
      </div>
    </div>

    <div class="filter-row">
      <el-input v-model="filters.search" placeholder="搜索标题或摘要 Search..." clearable style="width: 240px" @input="debounceFetch" />
      <el-select v-model="filters.status" placeholder="全部状态 All Status" clearable style="width: 160px" @change="fetchPosts">
        <el-option label="已发布 Published" value="published" />
        <el-option label="草稿 Draft" value="draft" />
        <el-option label="已归档 Archived" value="archived" />
      </el-select>
      <el-select v-model="filters.tag" placeholder="全部标签 All Tags" clearable style="width: 160px" @change="fetchPosts">
        <el-option v-for="t in allTags" :key="t.id" :label="t.name" :value="t.name" />
      </el-select>
    </div>

    <div class="slide-up">
    <el-table :data="posts" stripe>
      <el-table-column prop="title" label="标题 Title" min-width="200" show-overflow-tooltip />
      <el-table-column label="标签 Tags" width="180">
        <template #default="{ row }">
          <span v-for="tag in (row.tags || [])" :key="tag.id || tag.name" class="tag-pill">{{ tag.name }}</span>
          <span v-if="!row.tags?.length" style="color: var(--muted); font-size: 13px">-</span>
        </template>
      </el-table-column>
      <el-table-column label="状态 Status" width="100">
        <template #default="{ row }">
          <span :class="['status-badge', 'status-' + row.status]">{{ statusLabel(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="views" label="阅读 Views" width="100" />
      <el-table-column label="日期 Date" width="120">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作 Actions" width="200" fixed="right">
        <template #default="{ row }">
          <div class="action-btn-group">
            <button class="action-btn" @click="viewPost(row)">查看 View</button>
            <button class="action-btn" @click="router.push(`/admin/posts/${row.id}/edit`)">编辑</button>
            <button class="action-btn action-danger" @click="deletePost(row.id)">删除</button>
          </div>
        </template>
      </el-table-column>
    </el-table>
    </div>

    <div class="pagination-wrap" v-if="total > perPage">
      <span class="pagination-summary">共 {{ total }} 篇文章 · 每页
        <el-select v-model="perPage" size="small" style="width: 80px; margin: 0 4px" @change="onPerPageChange">
          <el-option :value="10" label="10" /><el-option :value="20" label="20" /><el-option :value="50" label="50" />
        </el-select>
        条
      </span>
      <el-pagination
        layout="prev, pager, next, jumper"
        :total="total"
        :page-size="perPage"
        :current-page="page"
        @current-change="p => { page = p; fetchPosts() }"
      />
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'
import { formatDate } from '../utils/format'

const router = useRouter()
const fileInput = ref(null)

const posts = ref([])
const allTags = ref([])
const total = ref(0)
const page = ref(1)
const perPage = ref(10)
let debounceTimer = null

const filters = reactive({ search: '', status: '', tag: '' })

function statusLabel(s) {
  return { published: '已发布', draft: '草稿', archived: '已归档' }[s] || s
}

function viewPost(row) {
  window.open(`/post/${row.id}`, '_blank')
}

function onPerPageChange() {
  page.value = 1
  fetchPosts()
}

function debounceFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; fetchPosts() }, 300)
}

async function fetchPosts() {
  const params = { page: page.value, per_page: perPage.value }
  if (filters.search) params.search = filters.search
  if (filters.status) params.status = filters.status
  if (filters.tag) params.tag = filters.tag
  const res = await api.get('/posts', { params }).catch(() => ({ data: [], total: 0 }))
  posts.value = res.data || res.posts || []
  total.value = res.total || 0
}

async function fetchTags() {
  const res = await api.get('/tags').catch(() => ({ data: [] }))
  allTags.value = Array.isArray(res) ? res : (res.data || [])
}

async function deletePost(id) {
  try {
    await ElMessageBox.confirm('文章将移入回收站，可在回收站中恢复或彻底删除。', '移入回收站', { type: 'warning', confirmButtonText: '移入回收站', cancelButtonText: '取消' })
    await api.delete(`/posts/${id}`)
    ElMessage.success('已移入回收站')
    fetchPosts()
  } catch {}
}

function triggerUpload() {
  fileInput.value.click()
}

function handleUpload(e) {
  const file = e.target.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    sessionStorage.setItem('md-upload-content', reader.result)
    // Extract title from first heading if present
    const match = reader.result.match(/^#\s+(.+)$/m)
    if (match) sessionStorage.setItem('md-upload-title', match[1])
    router.push('/admin/posts/new')
  }
  reader.readAsText(file)
  e.target.value = ''
}

onMounted(() => {
  fetchTags()
  fetchPosts()
})
</script>

<style scoped>
.posts-page {
  width: 100%;
}
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.page-title {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 600;
  color: var(--fg);
}
.filter-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.pagination-wrap {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
}
.pagination-summary {
  font-size: 13px;
  color: var(--muted);
}
.status-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: var(--radius);
  display: inline-block;
}
.status-published { background: var(--success-soft); color: var(--success); }
.status-draft { background: var(--warning-soft); color: var(--warning); }
.status-archived { background: var(--error-soft); color: var(--muted); }
.tag-pill {
  font-size: 11px;
  padding: 2px 7px;
  border-radius: var(--radius);
  background: var(--tag-bg);
  color: var(--tag-fg);
  margin-right: 4px;
}
.action-btn-group {
  display: flex;
  align-items: center;
  gap: 4px;
}
.action-btn {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--fg);
  cursor: pointer;
  white-space: nowrap;
}
.action-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}
.action-danger {
  color: var(--error);
  border-color: var(--error-soft);
}
.action-danger:hover {
  color: var(--error);
  border-color: var(--error);
  background: var(--error-soft);
}
</style>
