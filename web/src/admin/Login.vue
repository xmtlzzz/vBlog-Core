<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">vBlog Admin</h1>
      <p class="login-subtitle">登录管理后台</p>
      <el-form @submit.prevent="handleLogin">
        <el-form-item>
          <el-input
            v-model="form.username"
            placeholder="用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          style="width: 100%"
          @click="handleLogin"
        >
          登录
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
const form = reactive({ username: '', password: '' })

async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await api.post('/auth/login', {
      username: form.username,
      password: form.password
    })
    authStore.setToken(res.access_token)
    ElMessage.success('登录成功')
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
.login-subtitle {
  text-align: center;
  font-size: 14px;
  color: var(--muted);
  margin-bottom: 32px;
}
</style>
