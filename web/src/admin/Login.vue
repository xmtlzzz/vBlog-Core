<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">vBlog Admin</h1>
      <div class="login-tabs">
        <span :class="{ active: mode === 'login' }" @click="mode = 'login'">登录</span>
        <span :class="{ active: mode === 'register' }" @click="mode = 'register'">注册</span>
      </div>
      <el-form @submit.prevent="mode === 'login' ? handleLogin() : handleRegister()">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" size="large" />
        </el-form-item>
        <el-form-item v-if="mode === 'register'">
          <el-input v-model="form.email" placeholder="邮箱（可选）" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password @keyup.enter="mode === 'login' ? handleLogin() : handleRegister()" />
        </el-form-item>
        <el-button type="primary" size="large" :loading="loading" style="width: 100%" @click="mode === 'login' ? handleLogin() : handleRegister()">
          {{ mode === 'login' ? '登录' : '注册' }}
        </el-button>
      </el-form>
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
const mode = ref('login')
const form = reactive({ username: '', password: '', email: '' })

async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await api.post('/auth/login', { username: form.username, password: form.password })
    authStore.setToken(res.access_token)
    ElMessage.success('登录成功')
    router.push('/admin')
  } catch {
    // error handled by interceptor
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await api.post('/auth/register', { username: form.username, password: form.password, email: form.email })
    authStore.setToken(res.access_token)
    ElMessage.success('注册成功')
    router.push('/admin')
  } catch {
    // error handled by interceptor
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
  border-radius: 12px;
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
.login-tabs {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 24px;
}
.login-tabs span {
  font-size: 14px;
  color: var(--muted);
  cursor: pointer;
  padding-bottom: 4px;
  border-bottom: 2px solid transparent;
}
.login-tabs span.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}
</style>
