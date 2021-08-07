import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import config from './src/constants/config'
import { resolve } from 'path'
import svgLoader from 'vite-svg-loader'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), svgLoader()],
  resolve: {
    alias: {
      '@/': resolve(process.cwd(), 'src') + '/',
    },
  },
  server: {
    proxy: {
      '/api': {
        target: `${config.apiHost}${config.apiPrefix}`,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})
