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
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/element-plus')) {
            return 'element-plus'
          }
          if (id.includes('node_modules/@element-plus/icons-vue')) {
            return 'element-icons'
          }
          if (id.includes('node_modules/md-editor-v3') || id.includes('node_modules/markdown-it')) {
            return 'markdown'
          }
          if (id.includes('node_modules/highlight.js')) {
            return 'highlight'
          }
        }
      }
    }
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
