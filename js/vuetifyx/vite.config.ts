// Plugins
import Components from 'unplugin-vue-components/vite'
import {Vuetify3Resolver} from 'unplugin-vue-components/resolvers'
import Vue from '@vitejs/plugin-vue'
import Vuetify, {transformAssetUrls} from 'vite-plugin-vuetify'
import ViteFonts from 'unplugin-fonts/vite'
import {resolve} from 'path'
import vueJsx from '@vitejs/plugin-vue-jsx'

import svgLoader from 'vite-svg-loader';

// Utilities
import {defineConfig, loadEnv} from 'vite'
import {fileURLToPath, URL} from 'node:url'
import {existsSync} from 'node:fs'

// The diff browser (<vx-diff-browser>) is the gad IDE's (web/ide-vuetify,
// src/diff): from the source of a gad checkout — GAD_DIR, or the one beside
// this workspace —, as js/gadide has the IDE.
const gadDir = resolve(process.env.GAD_DIR || resolve(__dirname, '../../../../../../../gade/src/github.com/gad-lang/gad'))
const gadDiff = resolve(gadDir, 'web/ide-vuetify/src/diff/index.ts')
if (!existsSync(gadDiff)) throw new Error(`vuetifyx: no ${gadDiff} (set GAD_DIR to a gad checkout with its submodules)`)

// https://vitejs.dev/config/
export default ({mode}) => {
  process.env = {...process.env, ...loadEnv(mode, process.cwd())};
  return defineConfig({
    build: {
      // minify: false,
      outDir: 'dist',
      emptyOutDir: true,
      lib: {
        entry: resolve(__dirname, 'src/lib/main.ts'),
        formats: ['umd'],
        name: 'vuetifyxjs'
      },
      copyPublicDir: false,
      rollupOptions: {
        external: ['vue', 'vuetify'],
        output: {
          assetFileNames: (assetInfo) => {
            return 'vuetifyxjs.css'
          },
          globals: {
            vue: 'Vue',
            vuetify: 'Vuetify'
          }
        }
      },
      cssCodeSplit: false,
    },

    publicDir: './src/demo/public',

    plugins: [
      svgLoader(),
      Vue({
        template: { transformAssetUrls }
      }),
      vueJsx(),
      // https://github.com/vuetifyjs/vuetify-loader/tree/master/packages/vite-plugin#readme
      Vuetify({
        autoImport: { labs: true },
        styles: {
          configFile: 'src/styles/settings.scss',
        },
      }),
      Components({
        dts: true,
        dirs: ['src/demo/components', 'src/lib'],
        resolvers: [Vuetify3Resolver()],
        include: [/\.vue$/]

      }),
      ViteFonts({
        google: {
          families: [{
            name: 'Roboto',
            styles: 'wght@100;300;400;500;700;900'
          }]
        }
      })//,
      //tsconfigPaths()
    ],
    define: { 'process.env': {} },
    css: {
      preprocessorOptions: {
        scss: {
          // additionalData: `@use "@/styles/settings.scss" as settings;`, // Adjust path as needed
        },
      },
    },
    resolve: {
      alias: {
        '@gad-lang/ide-vuetify/diff': gadDiff,
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
      // one of each — the gad source's imports resolved here, not in its
      // node_modules: CodeMirror's state fields need it, Vuetify's provides
      dedupe: ['vue', 'vuetify', '@codemirror/state', '@codemirror/view', '@codemirror/language',
        '@codemirror/autocomplete', '@codemirror/lint', '@codemirror/commands', '@codemirror/merge',
        '@codemirror/theme-one-dark', '@codemirror/lang-css', '@codemirror/lang-go', '@codemirror/lang-html',
        '@codemirror/lang-javascript', '@codemirror/lang-json',
        'dockview-vue', 'dockview-core', '@gad-lang/codemirror-gad'],
      extensions: [
        '.js',
        '.json',
        '.jsx',
        '.mjs',
        '.ts',
        '.tsx',
        '.vue'
      ]
    },
    server: {
      port: parseInt(process.env.VITE_PORT || '3000'),
      host: process.env.VITE_HOST === 'true',
      allowedHosts: true
    }
  })
}
