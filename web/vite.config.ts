import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

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
        target: process.env.DOCKHELM_API || 'http://127.0.0.1:8080',
        changeOrigin: true,
        ws: false,
      },
    },
  },
})
