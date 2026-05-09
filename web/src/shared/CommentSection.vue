<template>
  <section class="comment-section" v-if="enabled">
    <h3 class="comment-title">评论 Comments <span class="comment-count">({{ comments.length }})</span></h3>

    <!-- Comment list -->
    <div class="comment-list">
      <div v-for="c in comments" :key="c.id" class="comment-item">
        <div class="comment-avatar">{{ (c.author_name || '?')[0].toUpperCase() }}</div>
        <div class="comment-content">
          <div class="comment-header">
            <span class="comment-name">{{ c.author_name }}</span>
            <span class="comment-time">{{ formatTime(c.created_at) }}</span>
          </div>
          <p class="comment-body">{{ c.body }}</p>
        </div>
      </div>
      <div v-if="comments.length === 0" class="comment-empty">暂无评论，来发表第一条吧</div>
    </div>

    <!-- Comment form -->
    <form class="comment-form" @submit.prevent="submitComment">
      <div class="form-row">
        <input v-model="form.author_name" placeholder="昵称 *" required class="form-input" />
        <input v-model="form.author_email" placeholder="邮箱（可选）" type="email" class="form-input" />
      </div>
      <textarea v-model="form.body" placeholder="写下你的评论..." required class="form-textarea" rows="4"></textarea>
      <div class="form-footer">
        <span class="form-hint">评论将在审核后显示</span>
        <button type="submit" class="submit-btn" :disabled="submitting">
          {{ submitting ? '提交中...' : '发表评论' }}
        </button>
      </div>
    </form>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api/request'

const props = defineProps({ postId: [Number, String] })

const enabled = ref(false)
const comments = ref([])
const submitting = ref(false)
const form = ref({ author_name: '', author_email: '', body: '' })

function formatTime(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now - d
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + ' 分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + ' 小时前'
  return d.toLocaleDateString('zh-CN')
}

async function fetchComments() {
  try {
    const res = await api.get(`/posts/${props.postId}/comments`)
    comments.value = res.data || []
  } catch { comments.value = [] }
}

async function submitComment() {
  if (!form.value.author_name.trim() || !form.value.body.trim()) return
  submitting.value = true
  try {
    await api.post(`/posts/${props.postId}/comments`, form.value)
    ElMessage.success('评论已提交，等待审核')
    form.value = { author_name: form.value.author_name, author_email: form.value.author_email, body: '' }
  } catch {
    ElMessage.error('提交失败')
  } finally { submitting.value = false }
}

onMounted(async () => {
  try {
    const settings = await api.get('/settings')
    enabled.value = settings.enable_comments === 'true'
    if (enabled.value) fetchComments()
  } catch { enabled.value = false }
})
</script>

<style scoped>
.comment-section {
  margin-top: 48px;
  padding-top: 32px;
  border-top: 1px solid var(--border);
}
.comment-title {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 24px;
}
.comment-count {
  font-size: 14px;
  font-weight: 400;
  color: var(--muted);
}
.comment-list {
  margin-bottom: 32px;
}
.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--border);
}
.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--accent-soft);
  color: var(--accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}
.comment-content {
  flex: 1;
  min-width: 0;
}
.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}
.comment-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--fg);
}
.comment-time {
  font-size: 12px;
  color: var(--muted);
}
.comment-body {
  font-size: 15px;
  line-height: 1.6;
  color: var(--fg);
}
.comment-empty {
  text-align: center;
  padding: 32px 0;
  color: var(--muted);
  font-size: 14px;
}
.comment-form {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 20px;
}
.form-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}
.form-input {
  flex: 1;
  font-size: 14px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--fg);
  outline: none;
  font-family: var(--font-sans);
  transition: border-color 0.15s;
}
.form-input:focus {
  border-color: var(--accent);
}
.form-textarea {
  width: 100%;
  font-size: 14px;
  padding: 10px 12px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--fg);
  outline: none;
  font-family: var(--font-sans);
  resize: vertical;
  transition: border-color 0.15s;
  box-sizing: border-box;
}
.form-textarea:focus {
  border-color: var(--accent);
}
.form-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
}
.form-hint {
  font-size: 12px;
  color: var(--muted);
}
.submit-btn {
  font-size: 13px;
  font-weight: 500;
  padding: 8px 20px;
  border-radius: var(--radius);
  border: none;
  background: var(--accent);
  color: white;
  cursor: pointer;
  transition: opacity 0.15s;
  font-family: var(--font-sans);
}
.submit-btn:hover { opacity: 0.9; }
.submit-btn:active { transform: scale(0.97); }
.submit-btn:disabled { opacity: 0.5; cursor: not-allowed; }
</style>
