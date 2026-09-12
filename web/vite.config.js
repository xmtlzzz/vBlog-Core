import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          // 页脚 GitHub 二维码徽章（every-qrcode Web Component）
          isCustomElement: (tag) => tag === 'every-qr-code'
        }
      }
    })
  ],
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
