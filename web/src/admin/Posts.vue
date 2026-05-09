<template>
  <div class="posts-page">
    <div class="page-header">
      <h1 class="page-title">文章管理</h1>
      <el-button type="primary" @click="openCreate">新文章 New Post</el-button>
    </div>

    <div class="filter-row">
      <el-input v-model="filters.search" placeholder="搜索文章..." clearable style="width: 240px" @input="debounceFetch" />
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="fetchPosts">
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
        <el-option label="已归档" value="archived" />
      </el-select>
      <el-select v-model="filters.tag" placeholder="标签" clearable style="width: 140px" @change="fetchPosts">
        <el-option v-for="t in allTags" :key="t.id" :label="t.name" :value="t.name" />
      </el-select>
    </div>

    <el-table :data="posts" stripe>
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column label="标签" width="180">
        <template #default="{ row }">
          <el-tag v-for="tag in (row.tags || [])" :key="tag.id || tag.name" size="small" style="margin-right: 4px">{{ tag.name }}</el-tag>
          <span v-if="!row.tags?.length" style="color: var(--muted); font-size: 13px">-</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <span :class="['status-badge', row.status]">{{ statusLabel(row.status) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="views" label="阅读" width="80" />
      <el-table-column label="日期" width="120">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button size="small" text type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" text type="danger" @click="deletePost(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap" v-if="total > perPage">
      <el-pagination
        layout="prev, pager, next"
        :total="total"
        :page-size="perPage"
        :current-page="page"
        @current-change="p => { page = p; fetchPosts() }"
      />
    </div>

    <!-- Create/Edit Modal -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑文章' : '新建文章'"
      width="680px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="文章标题" />
        </el-form-item>
        <el-form-item label="标签">
          <el-select v-model="form.tagNames" multiple filterable allow-create placeholder="选择或输入标签" style="width: 100%">
            <el-option v-for="t in allTags" :key="t.id" :label="t.name" :value="t.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.excerpt" type="textarea" :rows="2" placeholder="文章摘要" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="10" placeholder="Markdown 内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="savePost">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const posts = ref([])
const allTags = ref([])
const total = ref(0)
const page = ref(1)
const perPage = 10
const dialogVisible = ref(false)
const editingId = ref(null)
const saving = ref(false)
let debounceTimer = null

const filters = reactive({ search: '', status: '', tag: '' })
const form = reactive({ title: '', content: '', excerpt: '', status: 'draft', tagNames: [] })

function statusLabel(s) {
  return { published: '已发布', draft: '草稿', archived: '已归档' }[s] || s
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}

function debounceFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; fetchPosts() }, 300)
}

async function fetchPosts() {
  const params = { page: page.value, per_page: perPage }
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

function openCreate() {
  editingId.value = null
  Object.assign(form, { title: '', content: '', excerpt: '', status: 'draft', tagNames: [] })
  dialogVisible.value = true
}

function openEdit(post) {
  editingId.value = post.id
  Object.assign(form, {
    title: post.title || '',
    content: post.content || '',
    excerpt: post.excerpt || '',
    status: post.status || 'draft',
    tagNames: (post.tags || []).map(t => t.name)
  })
  dialogVisible.value = true
}

async function savePost() {
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
      status: form.status,
      tags: form.tagNames.map(name => ({ name }))
    }
    if (editingId.value) {
      await api.put(`/posts/${editingId.value}`, payload)
      ElMessage.success('文章已更新')
    } else {
      await api.post('/posts', payload)
      ElMessage.success('文章已创建')
    }
    dialogVisible.value = false
    fetchPosts()
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function deletePost(id) {
  try {
    await ElMessageBox.confirm('确定删除这篇文章？', '确认删除', { type: 'warning' })
    await api.delete(`/posts/${id}`)
    ElMessage.success('已删除')
    fetchPosts()
  } catch {}
}

onMounted(() => {
  fetchTags()
  fetchPosts()
})
</script>

<style scoped>
.posts-page {
  max-width: 1200px;
}
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
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
  justify-content: center;
  margin-top: 20px;
}
.status-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
}
.status-badge.published {
  background: rgba(22, 163, 74, 0.1);
  color: var(--success);
}
.status-badge.draft {
  background: rgba(245, 158, 11, 0.1);
  color: var(--warning);
}
.status-badge.archived {
  background: rgba(115, 115, 115, 0.1);
  color: var(--muted);
}
</style>
