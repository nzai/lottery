import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // dev 模式下前端 5173 把 /api 转发给 Go 服务
      // （端口 23817 是用户指定的非默认端口，见 AGENTS.md）
      '/api': 'http://localhost:23817',
    },
  },
  build: {
    // 构建产物直接输出到 server/static——这是 go:embed 的嵌入源目录，
    // 部署时前端页面随 Go 二进制一起打包（单文件部署，见 AGENTS.md）
    outDir: '../server/static',
    emptyOutDir: true,
  },
  test: {
    environment: 'jsdom',
    // globals 开启后 @testing-library/react 自动在 afterEach 清理 DOM
    globals: true,
  },
})
