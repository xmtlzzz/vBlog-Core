<template>
  <BlogNav />
  <div class="page-enter">
    <div class="friends-container">
      <header class="page-header">
        <h1 class="page-title">友情链接 Friends</h1>
        <p class="page-desc">海内存知己，天涯若比邻。欢迎与技术爱好者、生活记录者互换友链。</p>
      </header>

      <!-- Friends List Grid -->
      <section class="friends-section" v-if="friendsList.length">
        <div class="friends-grid">
          <a
            v-for="(item, idx) in friendsList"
            :key="idx"
            :href="item.link"
            target="_blank"
            rel="noopener noreferrer"
            class="friend-card"
          >
            <div class="friend-avatar-wrap">
              <img
                v-if="item.avatar && !avatarErrors[idx]"
                :src="item.avatar"
                :alt="item.name"
                class="friend-avatar"
                @error="avatarErrors[idx] = true"
                loading="lazy"
              />
              <div v-else class="friend-avatar-fallback">
                {{ (item.name || 'F').charAt(0).toUpperCase() }}
              </div>
            </div>
            <div class="friend-info">
              <div class="friend-header">
                <span class="friend-name">{{ item.name }}</span>
                <svg class="external-icon" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                  <polyline points="15 3 21 3 21 9"></polyline>
                  <line x1="10" y1="14" x2="21" y2="3"></line>
                </svg>
              </div>
              <p class="friend-desc" :title="item.desc">{{ item.desc || '暂无简介' }}</p>
            </div>
          </a>
        </div>
      </section>

      <!-- Empty State -->
      <div v-else-if="loaded" class="empty-state">
        <p>暂无友情链接，欢迎按照下方格式申请交换！</p>
      </div>

      <!-- Apply / Site Info Section -->
      <section class="apply-section">
        <div class="apply-card">
          <div class="apply-header">
            <h2>申请互换友链</h2>
            <button class="copy-site-info-btn" @click="copySiteInfo">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <span>{{ copiedInfo ? '已复制本站信息' : '复制本站信息' }}</span>
            </button>
          </div>
          <p class="apply-hint">
            在申请友链前，请先将本站添加到您的博客中。申请可通过邮件发送至
            <a v-if="settings.author_email" :href="'mailto:' + settings.author_email">{{ settings.author_email }}</a>
            <span v-else>站长邮箱</span>
            或在关于页留言。
          </p>
          <div class="site-info-box">
            <div class="info-row">
              <span class="info-label">名称：</span>
              <span class="info-value">{{ settings.site_title || 'vBlog' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">简介：</span>
              <span class="info-value">{{ settings.description || settings.subtitle || '写代码的人，也写点别的。' }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">链接：</span>
              <span class="info-value">{{ siteUrl }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">头像：</span>
              <span class="info-value">{{ siteUrl }}/avatar.png</span>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
  <CustomWidgets />
  <BlogFooter />
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import CustomWidgets from '../shared/CustomWidgets.vue'

const settings = ref({})
const loaded = ref(false)
const copiedInfo = ref(false)
const avatarErrors = reactive({})

const siteUrl = computed(() => {
  if (typeof window !== 'undefined') {
    return window.location.origin
  }
  return 'https://vblog.xmtlz.dev'
})

// 解析 settings.friends (每行：名称|链接|头像|简介)
const friendsList = computed(() => {
  const raw = (settings.value.friends || '').trim()
  if (!raw) return []
  return raw
    .split('\n')
    .map(line => line.trim())
    .filter(Boolean)
    .map(line => {
      const parts = line.split('|').map(s => s.trim())
      return {
        name: parts[0] || '',
        link: parts[1] || '#',
        avatar: parts[2] || '',
        desc: parts[3] || ''
      }
    })
    .filter(item => item.name && item.link && item.link !== '#')
})

function copySiteInfo() {
  const name = settings.value.site_title || 'vBlog'
  const desc = settings.value.description || settings.value.subtitle || '写代码的人，也写点别的。'
  const text = `博客名称：${name}\n博客网址：${siteUrl.value}\n博客简介：${desc}\n头像链接：${siteUrl.value}/avatar.png`
  try {
    navigator.clipboard.writeText(text)
    copiedInfo.value = true
    setTimeout(() => { copiedInfo.value = false }, 2500)
  } catch {
    // clipboard error fallback
  }
}

onMounted(async () => {
  try {
    const res = await api.get('/settings').catch(() => ({}))
    settings.value = Array.isArray(res) ? Object.fromEntries(res.map(s => [s.key, s.value])) : (res || {})
  } finally {
    loaded.value = true
  }
})
</script>

<style scoped>
.friends-container {
  max-width: 860px;
  margin: 0 auto;
  padding: 56px 24px 80px;
}

.page-header {
  margin-bottom: 40px;
}
.page-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 10px;
  color: var(--fg);
}
.page-desc {
  font-size: 15px;
  color: var(--muted);
  margin: 0;
  line-height: 1.6;
}

/* 友链卡片网格 */
.friends-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  margin-bottom: 48px;
}

.friend-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px;
  border-radius: var(--radius-lg);
  background: var(--surface);
  border: 1px solid var(--border);
  text-decoration: none;
  color: inherit;
  transition: all 0.22s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.friend-card:hover {
  transform: translateY(-2px);
  border-color: var(--accent);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
}

.friend-avatar-wrap {
  flex-shrink: 0;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--tag-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--border);
}

.friend-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.friend-avatar-fallback {
  font-size: 20px;
  font-weight: 700;
  color: var(--accent);
}

.friend-info {
  flex: 1;
  min-width: 0;
}

.friend-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 4px;
}

.friend-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--fg);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.external-icon {
  flex-shrink: 0;
  color: var(--muted);
  opacity: 0.6;
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.friend-card:hover .external-icon {
  opacity: 1;
  color: var(--accent);
  transform: translate(1px, -1px);
}

.friend-desc {
  font-size: 13px;
  color: var(--muted);
  margin: 0;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-state {
  text-align: center;
  padding: 48px 20px;
  color: var(--muted);
  font-size: 14px;
  background: var(--surface);
  border: 1px dashed var(--border);
  border-radius: var(--radius-lg);
  margin-bottom: 48px;
}

/* 申请卡片 */
.apply-section {
  margin-top: 32px;
}

.apply-card {
  padding: 24px;
  border-radius: var(--radius-lg);
  background: var(--surface);
  border: 1px solid var(--border);
}

.apply-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 12px;
}

.apply-header h2 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--fg);
}

.copy-site-info-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  padding: 6px 14px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--fg);
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-site-info-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}

.apply-hint {
  font-size: 13px;
  color: var(--muted);
  margin: 0 0 16px;
  line-height: 1.6;
}

.apply-hint a {
  color: var(--accent);
  text-decoration: none;
}

.apply-hint a:hover {
  text-decoration: underline;
}

.site-info-box {
  background: var(--tag-bg);
  border-radius: var(--radius);
  padding: 14px 18px;
  font-size: 13px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
  word-break: break-all;
}

.info-label {
  font-weight: 600;
  color: var(--fg);
  flex-shrink: 0;
}

.info-value {
  color: var(--muted);
}
</style>
