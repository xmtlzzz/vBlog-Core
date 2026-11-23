import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue({ template: { transformAssetUrls: false } })],
  test: { environment: 'jsdom', css: { include: [/tokens\.css/] } }
})
