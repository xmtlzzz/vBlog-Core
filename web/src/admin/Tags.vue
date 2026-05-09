<template>
  <div class="tags-page">
    <div class="page-header">
      <h1 class="page-title">标签管理</h1>
      <el-button type="primary" @click="openCreate">新标签 New Tag</el-button>
    </div>

    <div class="tag-grid fade-in stagger">
      <div v-for="tag in tags" :key="tag.id" class="tag-card card-hover">
        <div class="tag-card-header">
          <span class="tag-name">{{ tag.name }}</span>
          <el-tag size="small" type="info">{{ tag.post_count || 0 }} 篇</el-tag>
        </div>
        <p class="tag-desc">{{ tag.description || '暂无描述' }}</p>
        <div class="tag-meta">
          <span>ID: {{ tag.id }}</span>
          <span>{{ formatDate(tag.created_at) }}</span>
        </div>
        <div class="tag-actions">
          <el-button size="small" text type="primary" @click="openEdit(tag)">编辑</el-button>
          <el-button size="small" text type="danger" @click="deleteTag(tag.id)">删除</el-button>
        </div>
      </div>
      <div v-if="tags.length === 0" class="empty-state">暂无标签</div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑标签' : '新建标签'" width="420px" destroy-on-close>
      <el-form label-position="top">
        <el-form-item label="标签名">
          <el-input v-model="form.name" placeholder="标签名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="标签描述（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveTag">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const tags = ref([])
const dialogVisible = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = reactive({ name: '', description: '' })

function formatDate(d) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('zh-CN')
}

async function fetchTags() {
  const res = await api.get('/tags').catch(() => ({ data: [] }))
  tags.value = Array.isArray(res) ? res : (res.data || [])
}

function openCreate() {
  editingId.value = null
  Object.assign(form, { name: '', description: '' })
  dialogVisible.value = true
}

function openEdit(tag) {
  editingId.value = tag.id
  Object.assign(form, { name: tag.name, description: tag.description || '' })
  dialogVisible.value = true
}

async function saveTag() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入标签名')
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await api.put(`/tags/${editingId.value}`, { name: form.name, description: form.description })
      ElMessage.success('标签已更新')
    } else {
      await api.post('/tags', { name: form.name, description: form.description })
      ElMessage.success('标签已创建')
    }
    dialogVisible.value = false
    fetchTags()
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function deleteTag(id) {
  try {
    await ElMessageBox.confirm('确定删除此标签？', '确认删除', { type: 'warning' })
    await api.delete(`/tags/${id}`)
    ElMessage.success('已删除')
    fetchTags()
  } catch {}
}

onMounted(fetchTags)
</script>

<style scoped>
.tags-page {
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
.tag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}
.tag-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
  transition: border-color 0.15s;
}
.tag-card:hover {
  border-color: var(--accent);
}
.tag-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.tag-name {
  font-weight: 600;
  font-size: 16px;
  color: var(--fg);
}
.tag-desc {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 12px;
  line-height: 1.5;
}
.tag-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--muted);
  margin-bottom: 12px;
}
.tag-actions {
  display: flex;
  gap: 8px;
}
.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 48px;
  color: var(--muted);
}
</style>
