<template>
  <div class="dashboard">
    <div class="stats-grid fade-in stagger">
      <div class="stat-card card-hover">
        <div class="stat-label">◎ 文章总数 Total Posts</div>
        <div class="stat-card-value">{{ stats.total_posts }}</div>
        <div class="stat-change">已发布文章</div>
      </div>
      <div class="stat-card card-hover">
        <div class="stat-label">≡ 总阅读量 Total Views</div>
        <div class="stat-card-value">{{ stats.total_views.toLocaleString() }}</div>
        <div class="stat-change">累计阅读量</div>
      </div>
      <div class="stat-card card-hover">
        <div class="stat-label">❝ 评论数 Comments</div>
        <div class="stat-card-value">{{ stats.total_comments }}</div>
        <div class="stat-change">全部评论</div>
      </div>
      <div class="stat-card card-hover">
        <div class="stat-label">◉ 标签数 Tags</div>
        <div class="stat-card-value">{{ stats.total_tags }}</div>
        <div class="stat-change">标签数量</div>
      </div>
    </div>

    <div class="section slide-up">
      <div class="table-header">
        <span class="table-title">文章管理 Post Management</span>
        <div class="table-actions">
          <el-input v-model="search" placeholder="搜索 Search..." style="width: 180px" clearable />
          <el-select v-model="statusFilter" placeholder="全部状态 All Status" style="width: 160px" clearable>
            <el-option label="全部状态 All Status" value="" />
            <el-option label="已发布 Published" value="published" />
            <el-option label="草稿 Draft" value="draft" />
          </el-select>
        </div>
      </div>
      <el-table :data="filteredPosts" stripe style="width: 100%">
        <el-table-column prop="title" label="标题 Title" min-width="200" show-overflow-tooltip />
        <el-table-column label="标签 Tags" width="180">
          <template #default="{ row }">
            <span
              v-for="tag in (row.tags || [])"
              :key="tag.id || tag.name"
              class="tag-pill"
            >{{ tag.name }}</span>
            <span v-if="!row.tags?.length" style="color: var(--muted); font-size: 13px">-</span>
          </template>
        </el-table-column>
        <el-table-column label="状态 Status" width="100">
          <template #default="{ row }">
            <span :class="['status-badge', 'status-' + row.status]">{{ statusLabel(row.status) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="阅读 Views" width="80" />
        <el-table-column label="日期 Date" width="120">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作 Actions" width="140" fixed="right">
          <template #default="{ row }">
            <div class="action-btn-group">
              <button class="action-btn" @click="viewPost(row)">查看 View</button>
              <button class="action-btn" @click="$router.push('/admin/posts')">编辑</button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const stats = ref({ total_posts: 0, total_views: 0, total_comments: 0, total_tags: 0 })
const posts = ref([])
const search = ref('')
const statusFilter = ref('')

const filteredPosts = computed(() => {
  return posts.value.filter(p => {
    const matchSearch = !search.value || p.title?.toLowerCase().includes(search.value.toLowerCase())
    const matchStatus = !statusFilter.value || p.status === statusFilter.value
    return matchSearch && matchStatus
  })
})

function statusLabel(s) {
  return { published: '已发布', draft: '草稿', archived: '已归档' }[s] || s
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}

function viewPost(row) {
  window.open(`/post/${row.id}`, '_blank')
}

async function deletePost(id) {
  try {
    await ElMessageBox.confirm('确定删除这篇文章？', '确认删除', { type: 'warning' })
    await api.delete(`/posts/${id}`)
    posts.value = posts.value.filter(p => p.id !== id)
    stats.value.total_posts--
    ElMessage.success('已删除')
  } catch {}
}

onMounted(async () => {
  const [statsRes, postsRes] = await Promise.all([
    api.get('/dashboard/stats').catch(() => ({})),
    api.get('/posts', { params: { page: 1, per_page: 10 } }).catch(() => ({ data: [] }))
  ])
  stats.value = {
    total_posts: statsRes.total_posts || 0,
    total_views: statsRes.total_views || 0,
    total_comments: statsRes.total_comments || 0,
    total_tags: statsRes.total_tags || 0
  }
  posts.value = postsRes.data || postsRes.posts || []
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}
.stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
}
.stat-card-label {
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 8px;
}
.stat-label {
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 8px;
}
.stat-card-value {
  font-size: 28px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  font-family: var(--font-mono);
  color: var(--fg);
}
.stat-change {
  font-size: 12px;
  margin-top: 4px;
  font-family: var(--font-mono);
  color: var(--success);
}
.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.table-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  color: var(--fg);
}
.table-actions {
  display: flex;
  gap: 8px;
  align-items: center;
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
  transition: all 0.15s ease;
}
.action-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}
.action-btn:active {
  transform: scale(0.96);
}
</style>
