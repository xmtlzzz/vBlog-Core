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
    <div class="settings-section slide-up">
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
    <div class="settings-section slide-up" style="animation-delay: 100ms">
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

    <!-- gRPC -->
    <div class="settings-section slide-up" style="animation-delay: 200ms">
      <h2 class="section-title">gRPC 监控</h2>
      <el-form label-position="top">
        <el-form-item label="API Key">
          <div class="api-key-row">
            <el-input v-model="settings.grpc_api_key" :type="showKey ? 'text' : 'password'" placeholder="桌面客户端连接密钥" />
            <el-button @click="showKey = !showKey">{{ showKey ? '隐藏' : '显示' }}</el-button>
            <el-button @click="generateKey">生成</el-button>
          </div>
          <div class="field-hint">桌面监控客户端连接时需要此密钥，留空则不启用认证</div>
        </el-form-item>
        <el-form-item label="gRPC 端口">
          <el-input v-model="settings.grpc_port" placeholder="50051" />
        </el-form-item>
      </el-form>
    </div>

    <!-- Features -->
    <div class="settings-section slide-up" style="animation-delay: 300ms">
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

  </div>

  <Transition name="fade">
    <button v-show="showTop" class="back-to-top" @click="scrollToTop" title="回到顶部">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 15l-6-6-6 6"/>
      </svg>
    </button>
  </Transition>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api/request'

const settings = ref({})
const saving = ref(false)
const showKey = ref(false)

function generateKey() {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let key = 'vblog_'
  for (let i = 0; i < 32; i++) {
    key += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  settings.value.grpc_api_key = key
}

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

const showTop = ref(false)
function onScroll() { showTop.value = window.scrollY > 400 }
function scrollToTop() { window.scrollTo({ top: 0, behavior: 'smooth' }) }

onMounted(() => {
  fetchSettings()
  window.addEventListener('scroll', onScroll, { passive: true })
})
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<style scoped>
.settings-page {
  width: 100%;
  max-width: 960px;
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
  border-radius: var(--radius);
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
.api-key-row {
  display: flex;
  gap: 8px;
}
.api-key-row .el-input {
  flex: 1;
}
.field-hint {
  font-size: 12px;
  color: var(--muted);
  margin-top: 4px;
}

.back-to-top {
  position: fixed;
  bottom: 32px;
  right: 32px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--fg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  transition: all 0.2s ease;
  z-index: 50;
}
.back-to-top:hover {
  border-color: var(--accent);
  color: var(--accent);
  transform: translateY(-2px);
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

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
