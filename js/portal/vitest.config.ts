import { defineConfig } from 'vitest/config'
import Vue from '@vitejs/plugin-vue'
import VueJsx from '@vitejs/plugin-vue-jsx'

export default defineConfig({
  plugins: [Vue(), VueJsx()],
  resolve: {
    // the templates of a content are compiled in the browser: Vue with its
    // compiler
    alias: { vue: 'vue/dist/vue.esm-bundler.js' }
  },
  test: {
    environment: 'jsdom'
  }
})
