import { defineConfig } from 'vite'
import { resolve } from 'node:path'
import Vue from '@vitejs/plugin-vue'
import VueJsx from '@vitejs/plugin-vue-jsx'

// The bundle of the portals of a site (src/site.ts): rvq-portal.umd.js, Vue
// the page's global (window.Vue, with its compiler).
export default defineConfig({
  plugins: [Vue(), VueJsx()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    lib: {
      entry: resolve(__dirname, 'src/site.ts'),
      formats: ['umd'],
      name: 'RvqPortal',
      fileName: () => 'rvq-portal.umd.js'
    },
    rollupOptions: {
      external: ['vue'],
      output: { globals: { vue: 'Vue' } }
    }
  }
})
