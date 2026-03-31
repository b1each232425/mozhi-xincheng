import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5220, // 换成 5174 或其他没用的端口
  },
  resolve: {
   alias: {
  '@': fileURLToPath(new URL('./src', import.meta.url))
}
  }
})
