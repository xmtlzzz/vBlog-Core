<template>
  <div class="dashboard">
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Total Posts</div>
        <div class="stat-value">{{ stats.total_posts }}</div>
        <div class="stat-change positive">已发布文章</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Total Views</div>
        <div class="stat-value">{{ stats.total_views.toLocaleString() }}</div>
        <div class="stat-change positive">累计阅读量</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Comments</div>
        <div class="stat-value">{{ stats.total_comments }}</div>
        <div class="stat-change neutral">全部评论</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Tags</div>
        <div class="stat-value">{{ stats.total_tags }}</div>
        <div class="stat-change neutral">标签数量</div>
      </div>
    </div>

    <div class="section">
      <h2 class="section-title">最近文章</h2>
      <el-table :data="posts" stripe style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column label="标签" width="180">
          <template #default="{ row }">
            <el-tag
              v-for="tag in (row.tags || [])"
              :key="tag.id || tag.name"
              size="small"
              style="margin-right: 4px"
            >{{ tag.name }}</el-tag>
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
            <el-button size="small" text type="primary" @click="$router.push('/admin/posts')">编辑</el-button>
            <el-button size="small" text type="danger" @click="deletePost(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const stats = ref({ total_posts: 0, total_views: 0, total_comments: 0, total_tags: 0 })
const posts = ref([])

function statusLabel(s) {
  return { published: '已发布', draft: '草稿', archived: '已归档' }[s] || s
}

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
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
  max-width: 1200px;
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}
.stat-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
}
.stat-label {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 8px;
}
.stat-value {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  color: var(--fg);
  letter-spacing: -0.02em;
}
.stat-change {
  font-size: 12px;
  margin-top: 4px;
}
.stat-change.positive { color: var(--success); }
.stat-change.neutral { color: var(--muted); }
.section-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--fg);
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
