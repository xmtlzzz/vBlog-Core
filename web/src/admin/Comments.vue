<template>
  <div class="comments-page">
    <div class="page-header">
      <h1 class="page-title">评论管理</h1>
    </div>

    <div class="filter-row">
      <el-input v-model="filters.search" placeholder="搜索评论..." clearable style="width: 240px" @input="debounceFetch" />
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 140px" @change="fetchComments">
        <el-option label="全部" value="" />
        <el-option label="待审核" value="pending" />
        <el-option label="已通过" value="approved" />
        <el-option label="垃圾" value="spam" />
      </el-select>
    </div>

    <div class="comment-list fade-in stagger">
      <div v-for="c in comments" :key="c.id" class="comment-card card-hover">
        <div class="comment-avatar">{{ (c.author_name || '?')[0].toUpperCase() }}</div>
        <div class="comment-body">
          <div class="comment-header">
            <span class="comment-author">{{ c.author_name }}</span>
            <span class="comment-email">{{ c.author_email }}</span>
            <span :class="['status-badge', c.status]">{{ statusLabel(c.status) }}</span>
          </div>
          <p class="comment-text">{{ c.body }}</p>
          <div class="comment-footer">
            <span class="comment-post">文章 #{{ c.post_id }}</span>
            <span class="comment-time">{{ formatDate(c.created_at) }}</span>
          </div>
          <div class="comment-actions">
            <el-button v-if="c.status === 'pending'" size="small" type="success" text @click="approveComment(c)">通过</el-button>
            <el-button v-if="c.status === 'pending'" size="small" type="warning" text @click="spamComment(c)">垃圾</el-button>
            <el-button size="small" type="danger" text @click="deleteComment(c.id)">删除</el-button>
          </div>
        </div>
      </div>
      <div v-if="comments.length === 0" class="empty-state">暂无评论</div>
    </div>

    <div class="pagination-wrap" v-if="total > perPage">
      <el-pagination
        layout="prev, pager, next"
        :total="total"
        :page-size="perPage"
        :current-page="page"
        @current-change="p => { page = p; fetchComments() }"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'
import { formatDate } from '../utils/format'

const comments = ref([])
const total = ref(0)
const page = ref(1)
const perPage = 10
let debounceTimer = null

const filters = reactive({ search: '', status: '' })

function statusLabel(s) {
  return { pending: '待审核', approved: '已通过', spam: '垃圾' }[s] || s
}


function debounceFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; fetchComments() }, 300)
}

async function fetchComments() {
  const params = { page: page.value, per_page: perPage }
  if (filters.search) params.search = filters.search
  if (filters.status) params.status = filters.status
  const res = await api.get('/comments', { params }).catch(() => ({ data: [], total: 0 }))
  comments.value = res.data || []
  total.value = res.total || 0
}

async function approveComment(c) {
  await api.patch(`/comments/${c.id}/approve`)
  c.status = 'approved'
  ElMessage.success('已通过')
}

async function spamComment(c) {
  await api.patch(`/comments/${c.id}/spam`)
  c.status = 'spam'
  ElMessage.success('已标记为垃圾')
}

async function deleteComment(id) {
  try {
    await ElMessageBox.confirm('确定删除此评论？', '确认删除', { type: 'warning' })
    await api.delete(`/comments/${id}`)
    ElMessage.success('已删除')
    fetchComments()
  } catch {}
}

onMounted(fetchComments)
</script>

<style scoped>
.comments-page {
  width: 100%;
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
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.comment-card {
  display: flex;
  gap: 16px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px 20px;
}
.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-soft);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}
.comment-body {
  flex: 1;
  min-width: 0;
}
.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}
.comment-author {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg);
}
.comment-email {
  font-size: 12px;
  color: var(--muted);
}
.comment-text {
  font-size: 14px;
  color: var(--fg);
  line-height: 1.6;
  margin-bottom: 8px;
}
.comment-footer {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 8px;
}
.comment-actions {
  display: flex;
  gap: 8px;
}
.status-badge {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
}
.status-badge.pending {
  background: var(--warning-soft);
  color: var(--warning);
}
.status-badge.approved {
  background: var(--success-soft);
  color: var(--success);
}
.status-badge.spam {
  background: var(--error-soft);
  color: var(--error);
}
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
.empty-state {
  text-align: center;
  padding: 48px;
  color: var(--muted);
}
</style>
