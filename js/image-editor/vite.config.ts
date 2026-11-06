import {fileURLToPath, URL} from 'node:url'
import {resolve} from 'node:path'
import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// vx-image-editor: a UMD registering its components in the admin's Vue app
// (window.__goplaidVueComponentRegisters) — Vue and Vuetify are the page's.
export default defineConfig({
  build: {
    outDir: 'dist',
    emptyOutDir: true,
    lib: {
      entry: resolve(__dirname, 'src/lib/main.ts'),
      formats: ['umd'],
      name: 'rvqImageEditor',
      fileName: () => 'image-editor.umd.js',
    },
    copyPublicDir: false,
    cssCodeSplit: false,
    rollupOptions: {
      external: ['vue'],
      output: {
        assetFileNames: () => 'image-editor.css',
        globals: {vue: 'Vue'},
      },
    },
  },
  plugins: [vue()],
  resolve: {
    alias: {'@': fileURLToPath(new URL('./src', import.meta.url))},
  },
  test: {
    environment: 'jsdom',
  },
} as any)
