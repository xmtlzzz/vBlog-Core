import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({ baseURL: '/api' })
api.interceptors.request.use(config => {
  const token = localStorage.getItem('vblog-token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})
api.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('vblog-token')
      window.location.href = '/admin/login'
    }
    ElMessage.error(err.response?.data?.error || '请求失败')
    return Promise.reject(err)
  }
)
export default api
