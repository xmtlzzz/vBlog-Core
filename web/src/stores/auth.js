import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('vblog-token') || '')
  function setToken(t) { token.value = t; localStorage.setItem('vblog-token', t) }
  function logout() { token.value = ''; localStorage.removeItem('vblog-token') }
  const isLoggedIn = () => !!token.value
  return { token, setToken, logout, isLoggedIn }
})
