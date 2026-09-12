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
          <div class="field-hint">同时用作首页 Hero 副标题</div>
        </el-form-item>
        <el-form-item label="首页标语">
          <el-input
            v-model="settings.hero_title"
            type="textarea"
            :rows="2"
            placeholder="写代码的人，&#10;也写点别的。"
          />
          <div class="field-hint">首页打字机大标题，支持换行；留空使用默认文案</div>
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

    <!-- About -->
    <div class="settings-section slide-up" style="animation-delay: 150ms">
      <h2 class="section-title">关于页</h2>
      <el-form label-position="top">
        <el-form-item label="技术栈">
          <el-input
            v-model="settings.about_tech"
            type="textarea"
            :rows="5"
            placeholder="Vue|前端框架|vue&#10;Go|后端语言|🐹&#10;Cloudflare Workers|部署|⛅"
          />
          <div class="field-hint">每行一条：名称|角色|图标（图标可填 emoji 或 vue，留空显示首字符）；整项留空使用内置默认卡片</div>
        </el-form-item>
      </el-form>
    </div>

    <!-- QR badge -->
    <div class="settings-section slide-up" style="animation-delay: 250ms">
      <h2 class="section-title">页脚二维码徽章</h2>
      <div class="toggle-list">
        <div class="toggle-item">
          <div class="toggle-info">
            <div class="toggle-label">启用 GitHub 徽章</div>
            <div class="toggle-desc">页脚展示 every-qrcode 生成的动态二维码（点击可在 3D 模型与二维码间变形；需先填写作者 GitHub）</div>
          </div>
          <el-switch v-model="qrBadgeEnabled" />
        </div>
      </div>
      <el-form label-position="top" style="margin-top: 16px">
        <div class="form-grid">
          <el-form-item label="模型形态">
            <el-select v-model="settings.qr_model" style="width: 100%">
              <el-option label="樱花树 Tree" value="tree" />
              <el-option label="地形 Terrain" value="terrain" />
            </el-select>
          </el-form-item>
          <el-form-item label="氛围特效">
            <el-select v-model="settings.qr_effect" style="width: 100%">
              <el-option label="跟随默认" value="auto" />
              <el-option label="平静 Calm" value="calm" />
              <el-option label="落雪 Snow" value="snow" />
              <el-option label="下雨 Rain" value="rain" />
              <el-option label="刮风 Wind" value="wind" />
            </el-select>
          </el-form-item>
          <el-form-item label="配色预设">
            <el-select v-model="settings.qr_palette" style="width: 100%">
              <el-option label="樱花 Sakura（默认）" value="sakura" />
              <el-option label="森林 Forest" value="forest" />
              <el-option label="海洋 Ocean" value="ocean" />
              <el-option label="落日 Sunset" value="sunset" />
            </el-select>
            <div class="field-hint">仅作用于樱花树模型；地形模型使用自带配色</div>
          </el-form-item>
          <el-form-item label="徽章尺寸（px）">
            <el-input-number v-model="qrSizeNum" :min="28" :max="120" :step="4" />
            <div class="field-hint">默认 40；点击变形为二维码时自动放大约 2.6 倍</div>
          </el-form-item>
        </div>
        <el-form-item label="项目二维码导航（模块页）">
          <el-input
            v-model="settings.qr_links"
            type="textarea"
            :rows="4"
            placeholder="My App|https://app.example.com&#10;Demo|https://demo.example.com"
          />
          <div class="field-hint">每行一条「名称|URL」，每个地址确定性生成专属 3D 徽章，展示在「模块」页顶部；留空隐藏该区块</div>
        </el-form-item>
        <el-form-item label="徽章背景色">
          <el-color-picker v-model="qrBackground" />
          <span class="field-hint" style="margin-left: 8px">默认白色；二维码为彩色模块，需浅色背景保证扫码对比度</span>
        </el-form-item>
      </el-form>
    </div>

    <!-- gRPC -->
    <div class="settings-section slide-up" style="animation-delay: 300ms">
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
    <div class="settings-section slide-up" style="animation-delay: 400ms">
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

const qrBadgeEnabled = computed({
  get: () => settings.value.qr_badge !== 'false',
  set: (v) => { settings.value.qr_badge = String(v) }
})

const qrBackground = computed({
  get: () => settings.value.qr_bg || '#ffffff',
  set: (v) => { settings.value.qr_bg = v || '' }
})

const qrSizeNum = computed({
  get: () => {
    const n = parseInt(settings.value.qr_size)
    return Number.isFinite(n) ? Math.min(120, Math.max(28, n)) : 40
  },
  set: (v) => { settings.value.qr_size = String(v) }
})

async function fetchSettings() {
  const res = await api.get('/settings').catch(() => ({}))
  settings.value = Array.isArray(res) ? Object.fromEntries(res.map(s => [s.key, s.value])) : (res || {})
  // qr 配置缺省值（仅前端表单展示，保存时才写入）
  if (settings.value.qr_model === undefined) settings.value.qr_model = 'tree'
  if (settings.value.qr_effect === undefined) settings.value.qr_effect = 'auto'
  if (settings.value.qr_palette === undefined) settings.value.qr_palette = 'sakura'
  if (settings.value.qr_size === undefined) settings.value.qr_size = '40'
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
