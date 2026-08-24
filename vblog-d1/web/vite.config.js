import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // wrangler dev 默认监听 8787（Go 版才是 8080）
      '/api': 'http://127.0.0.1:8787'
    }
  }
})
