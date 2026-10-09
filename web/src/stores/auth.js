import { defineStore } from 'pinia'
import { ref } from 'vue'
import { isTokenExpired } from '../utils/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('vblog-token') || '')
  const user = ref(null)

  function setToken(t) {
    token.value = t
    localStorage.setItem('vblog-token', t)
  }

  function setUser(u) {
    user.value = u
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('vblog-token')
  }

  const isLoggedIn = () => !!token.value && !isTokenExpired(token.value)

  return { token, user, setToken, setUser, logout, isLoggedIn }
})
