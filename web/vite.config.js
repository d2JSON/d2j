import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import EnvironmentPlugin from 'vite-plugin-environment'

export default defineConfig({
  server: {
    host: true,
    port: 3000,
    hmr: {
      port: 3001,
    },
  },
  build: {
    outDir: '../build',
    rollupOptions: {
      cache: false,
    },
  },
  plugins: [
    vue(),
    EnvironmentPlugin({
      API_URL: 'http://localhost:8080/api',
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
