<template>
  <BlogNav />
  <div class="page-enter">
  <div class="about-container">
    <div class="about-hero">
      <div class="about-avatar">{{ initial }}</div>
      <div class="about-intro">
        <h1>{{ settings.site_title || 'vBlog Core' }}</h1>
        <p>{{ settings.author_name || '匿名作者' }}</p>
      </div>
    </div>

    <section class="about-section" v-if="settings.author_bio">
      <h2>关于我</h2>
      <p>{{ settings.author_bio }}</p>
    </section>

    <section class="about-section">
      <h2>技术栈</h2>
      <div class="tech-grid">
        <div v-for="tech in techStack" :key="tech.name" class="tech-card">
          <div class="tech-card-icon">{{ tech.icon }}</div>
          <div class="tech-card-name">{{ tech.name }}</div>
          <div class="tech-card-role">{{ tech.role }}</div>
        </div>
      </div>
    </section>

    <section class="about-section" v-if="settings.author_github || settings.author_email">
      <h2>联系方式</h2>
      <ul class="links-list">
        <li v-if="settings.author_github">
          <span class="link-label">GitHub</span>
          <a :href="settings.author_github" target="_blank">{{ settings.author_github }}</a>
        </li>
        <li v-if="settings.author_email">
          <span class="link-label">Email</span>
          <a :href="`mailto:${settings.author_email}`">{{ settings.author_email }}</a>
        </li>
      </ul>
    </section>
  </div>
  </div>
  <BlogFooter />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'

const settings = ref({})

const initial = computed(() => {
  const name = settings.value.author_name || 'V'
  return name.charAt(0).toUpperCase()
})

const techStack = [
  { name: 'React', icon: '⚛', role: '前端框架' },
  { name: 'Go', icon: '🔵', role: '后端语言' },
  { name: 'PostgreSQL', icon: '🐘', role: '数据库' },
  { name: 'Wails', icon: '🖥', role: '桌面客户端' },
  { name: 'gRPC', icon: '📡', role: 'RPC 通信' },
  { name: 'Markdown', icon: '📝', role: '内容格式' },
]

onMounted(async () => {
  const res = await api.get('/settings').catch(() => ({}))
  settings.value = res || {}
})
</script>

<style scoped>
.about-container {
  max-width: 720px;
  margin: 0 auto;
  padding: 56px 24px 80px;
}
.about-hero {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 48px;
}
.about-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--accent-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 600;
  color: var(--accent);
  flex-shrink: 0;
  border: 2px solid var(--border);
}
.about-intro h1 {
  font-family: var(--font-display);
  font-size: clamp(28px, 4vw, 36px);
  font-weight: 600;
  letter-spacing: -0.03em;
  line-height: 1.15;
  margin-bottom: 6px;
}
.about-intro p {
  font-size: 15px;
  color: var(--muted);
}
.about-section {
  margin-bottom: 40px;
}
.about-section h2 {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 600;
  letter-spacing: -0.01em;
  margin-bottom: 16px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border);
}
.about-section p {
  font-size: 15px;
  color: var(--fg);
  line-height: 1.75;
  margin-bottom: 14px;
}
.tech-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
  margin-top: 16px;
}
.tech-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  transition: all 0.15s;
}
.tech-card:hover {
  border-color: var(--accent);
}
.tech-card-name {
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 4px;
}
.tech-card-role {
  font-size: 12px;
  color: var(--muted);
}
.tech-card-icon {
  font-size: 20px;
  margin-bottom: 8px;
}
.links-list {
  list-style: none;
  padding: 0;
  margin-top: 12px;
}
.links-list li {
  padding: 10px 0;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.links-list li:last-child {
  border-bottom: none;
}
.links-list a {
  font-size: 14px;
  color: var(--accent);
  text-decoration: none;
  font-family: var(--font-mono);
}
.links-list a:hover {
  text-decoration: underline;
}
.link-label {
  font-size: 14px;
  color: var(--fg);
}

@media (max-width: 640px) {
  .about-container {
    padding: 40px 16px 48px;
  }
  .about-hero {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  .tech-grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
