<template>
  <div class="trash-page">
    <div class="page-header">
      <h1 class="page-title">回收站</h1>
    </div>

    <div v-if="posts.length" class="trash-list">
      <div v-for="post in posts" :key="post.id" class="trash-item">
        <div class="trash-info">
          <div class="trash-title">{{ post.title }}</div>
          <div class="trash-meta">
            <span>{{ post.status }}</span>
            <span>删除于 {{ post.deleted_at }}</span>
          </div>
        </div>
        <div class="trash-actions">
          <el-button size="small" type="primary" text @click="restore(post.id)">恢复</el-button>
          <el-button size="small" type="danger" text @click="permanentDelete(post.id)">彻底删除</el-button>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">回收站为空</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const posts = ref([])

async function fetchTrash() {
  const res = await api.get('/posts/trash').catch(() => ({ data: [] }))
  posts.value = Array.isArray(res) ? res : (res.data || [])
}

async function restore(id) {
  await api.post(`/posts/${id}/restore`)
  ElMessage.success('已恢复')
  fetchTrash()
}

async function permanentDelete(id) {
  try {
    await ElMessageBox.confirm('彻底删除后无法恢复，确定？', '确认删除', { type: 'warning' })
    await api.delete(`/posts/${id}/permanent`)
    ElMessage.success('已彻底删除')
    fetchTrash()
  } catch {}
}

onMounted(fetchTrash)
</script>

<style scoped>
.trash-page { width: 100%; max-width: 960px; }
.page-header { display: flex; align-items: center; margin-bottom: 24px; }
.page-title { font-family: var(--font-display); font-size: 22px; font-weight: 600; color: var(--fg); }
.trash-list { display: flex; flex-direction: column; gap: 8px; }
.trash-item {
  display: flex; align-items: center; justify-content: space-between;
  padding: 16px; background: var(--surface); border: 1px solid var(--border);
  border-radius: var(--radius);
}
.trash-title { font-weight: 500; font-size: 14px; }
.trash-meta { display: flex; gap: 12px; font-size: 12px; color: var(--muted); margin-top: 4px; }
.trash-actions { display: flex; gap: 8px; }
.empty-state { text-align: center; padding: 48px; color: var(--muted); }
</style>
