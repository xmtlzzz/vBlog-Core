import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('vblog-token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let authNoticeTimer = null
function notifyAuthExpired(msg) {
  if (!authNoticeTimer) {
    ElMessage.closeAll()
    ElMessage.warning(msg)
    authNoticeTimer = setTimeout(() => {
      authNoticeTimer = null
    }, 2500)
  }
}

api.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response?.status === 401) {
      const hadToken = !!localStorage.getItem('vblog-token')
      localStorage.removeItem('vblog-token')

      const msg = err.response?.data?.error === 'token expired'
        ? '登录已过期，请重新登录'
        : (hadToken ? '登录状态已失效，请重新登录' : '请先登录')

      sessionStorage.setItem('vblog_auth_expired_msg', msg)
      notifyAuthExpired(msg)

      const currentPath = router.currentRoute.value?.path || ''
      if (currentPath.startsWith('/admin') && currentPath !== '/admin/login') {
        router.push({
          path: '/admin/login',
          query: { redirect: router.currentRoute.value?.fullPath || '/admin' }
        }).catch(() => {})
      }
      return Promise.reject(err)
    }

    ElMessage.error(err.response?.data?.error || '请求失败')
    return Promise.reject(err)
  }
)

export default api
