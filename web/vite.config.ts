import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:23817',
    },
  },
  build: {
    outDir: '../server/static',
    emptyOutDir: true,
  },
  test: {
    environment: 'jsdom',
    // globals 开启后 @testing-library/react 自动在 afterEach 清理 DOM
    globals: true,
  },
})
