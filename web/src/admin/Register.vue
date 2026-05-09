<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">vBlog Admin</h1>
      <p class="login-subtitle">注册新账号</p>
      <el-form @submit.prevent="handleRegister">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.email" placeholder="邮箱（可选）" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password @keyup.enter="handleRegister" />
        </el-form-item>
        <el-button type="primary" size="large" :loading="loading" style="width: 100%" @click="handleRegister">
          注册
        </el-button>
      </el-form>
      <div class="login-footer">
        <router-link to="/admin/login">已有账号？去登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import api from '../api/request'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const form = reactive({ username: '', password: '', email: '' })

async function handleRegister() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await api.post('/auth/register', { username: form.username, password: form.password, email: form.email })
    if (res.access_token) {
      authStore.setToken(res.access_token)
      ElMessage.success('注册成功')
      await router.push('/admin')
    } else {
      ElMessage.error('注册失败：未收到令牌')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg);
}
.login-card {
  width: 360px;
  padding: 40px 32px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}
.login-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  text-align: center;
  color: var(--fg);
  margin-bottom: 4px;
  letter-spacing: -0.02em;
}
.login-subtitle {
  text-align: center;
  font-size: 14px;
  color: var(--muted);
  margin-bottom: 24px;
}
.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
}
.login-footer a {
  color: var(--accent);
  text-decoration: none;
}
.login-footer a:hover {
  text-decoration: underline;
}
</style>
