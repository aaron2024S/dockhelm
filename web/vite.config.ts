import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 开发时 /api 转发到本地跑的后端。
// 端口顺序与后端保持一致：DOCKHELM_PORT → PORT → 5923（可用 DOCKHELM_API 直接给完整地址）。
const apiPort = process.env.DOCKHELM_PORT || process.env.PORT || '5923'
const apiTarget = process.env.DOCKHELM_API || `http://127.0.0.1:${apiPort}`

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    // 关掉 sourcemap：产物要内嵌进二进制，越小越好
    sourcemap: false,
    chunkSizeWarningLimit: 1200,
  },
  server: {
    port: 5273,
    proxy: {
      // 本地开发时把 /api 转到 go run 起来的后端
      '/api': {
        target: apiTarget,
        changeOrigin: true,
        ws: false,
      },
    },
  },
})
