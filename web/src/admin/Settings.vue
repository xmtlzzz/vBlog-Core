<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <div class="header-actions">
        <el-button @click="resetSettings">重置默认</el-button>
        <el-button type="primary" :loading="saving" @click="saveSettings">保存设置</el-button>
      </div>
    </div>

    <!-- General -->
    <div class="settings-section">
      <h2 class="section-title">通用设置</h2>
      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="站点标题">
            <el-input v-model="settings.site_title" placeholder="vBlog" />
          </el-form-item>
          <el-form-item label="副标题">
            <el-input v-model="settings.subtitle" placeholder="副标题" />
          </el-form-item>
        </div>
        <el-form-item label="站点描述">
          <el-input v-model="settings.description" type="textarea" :rows="2" placeholder="站点描述" />
        </el-form-item>
        <div class="form-grid">
          <el-form-item label="语言">
            <el-select v-model="settings.language" style="width: 100%">
              <el-option label="中文" value="zh-CN" />
              <el-option label="English" value="en" />
            </el-select>
          </el-form-item>
          <el-form-item label="每页文章数">
            <el-input-number v-model="perPageNum" :min="1" :max="50" />
          </el-form-item>
        </div>
      </el-form>
    </div>

    <!-- Author -->
    <div class="settings-section">
      <h2 class="section-title">作者信息</h2>
      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="姓名">
            <el-input v-model="settings.author_name" placeholder="作者名" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="settings.author_email" placeholder="email@example.com" />
          </el-form-item>
        </div>
        <el-form-item label="个人简介">
          <el-input v-model="settings.author_bio" type="textarea" :rows="2" placeholder="个人简介" />
        </el-form-item>
        <el-form-item label="GitHub">
          <el-input v-model="settings.author_github" placeholder="https://github.com/username" />
        </el-form-item>
      </el-form>
    </div>

    <!-- Features -->
    <div class="settings-section">
      <h2 class="section-title">功能开关</h2>
      <div class="toggle-list">
        <div class="toggle-item">
          <div class="toggle-info">
            <div class="toggle-label">评论功能</div>
            <div class="toggle-desc">允许访客在文章下方发表评论</div>
          </div>
          <el-switch v-model="enableComments" />
        </div>
        <div class="toggle-item">
          <div class="toggle-info">
            <div class="toggle-label">RSS 订阅</div>
            <div class="toggle-desc">生成 RSS feed 供读者订阅</div>
          </div>
          <el-switch v-model="enableRss" />
        </div>
        <div class="toggle-item">
          <div class="toggle-info">
            <div class="toggle-label">阅读计数</div>
            <div class="toggle-desc">记录并展示每篇文章的阅读量</div>
          </div>
          <el-switch v-model="enableViewCounter" />
        </div>
      </div>
    </div>

    <!-- Database (read-only) -->
    <div class="settings-section">
      <h2 class="section-title">数据库信息</h2>
      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="主机">
            <el-input :model-value="settings.db_host" disabled />
          </el-form-item>
          <el-form-item label="端口">
            <el-input :model-value="settings.db_port" disabled />
          </el-form-item>
        </div>
        <div class="form-grid">
          <el-form-item label="数据库名">
            <el-input :model-value="settings.db_name" disabled />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input :model-value="settings.db_user" disabled />
          </el-form-item>
        </div>
        <el-form-item label="密码">
          <el-input :model-value="settings.db_password" type="password" disabled show-password />
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const settings = ref({})
const saving = ref(false)

const perPageNum = computed({
  get: () => parseInt(settings.value.posts_per_page) || 5,
  set: (v) => { settings.value.posts_per_page = String(v) }
})

const enableComments = computed({
  get: () => settings.value.enable_comments !== 'false',
  set: (v) => { settings.value.enable_comments = String(v) }
})

const enableRss = computed({
  get: () => settings.value.enable_rss !== 'false',
  set: (v) => { settings.value.enable_rss = String(v) }
})

const enableViewCounter = computed({
  get: () => settings.value.enable_view_counter !== 'false',
  set: (v) => { settings.value.enable_view_counter = String(v) }
})

async function fetchSettings() {
  const res = await api.get('/settings').catch(() => ({}))
  settings.value = Array.isArray(res) ? Object.fromEntries(res.map(s => [s.key, s.value])) : (res || {})
}

async function saveSettings() {
  saving.value = true
  try {
    await api.put('/settings', settings.value)
    ElMessage.success('设置已保存')
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function resetSettings() {
  try {
    await ElMessageBox.confirm('确定重置所有设置为默认值？', '确认重置', { type: 'warning' })
    await api.post('/settings/reset')
    ElMessage.success('已重置为默认设置')
    fetchSettings()
  } catch {}
}

onMounted(fetchSettings)
</script>

<style scoped>
.settings-page {
  max-width: 800px;
}
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.page-title {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 600;
  color: var(--fg);
}
.header-actions {
  display: flex;
  gap: 8px;
}
.settings-section {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
}
.section-title {
  font-family: var(--font-display);
  font-size: 16px;
  font-weight: 600;
  color: var(--fg);
  margin-bottom: 16px;
}
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 20px;
}
.toggle-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.toggle-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.toggle-item:last-child {
  border-bottom: none;
}
.toggle-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--fg);
}
.toggle-desc {
  font-size: 12px;
  color: var(--muted);
  margin-top: 2px;
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
