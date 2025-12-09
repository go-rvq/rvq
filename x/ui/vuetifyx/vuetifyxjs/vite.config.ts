// Plugins
import Components from 'unplugin-vue-components/vite'
import { Vuetify3Resolver } from 'unplugin-vue-components/resolvers'
import Vue from '@vitejs/plugin-vue'
import Vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import ViteFonts from 'unplugin-fonts/vite'
import { resolve } from 'path'
import vueJsx from '@vitejs/plugin-vue-jsx'

import svgLoader from 'vite-svg-loader';

// Utilities
import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default ({mode}) => {
  process.env = {...process.env, ...loadEnv(mode, process.cwd())};
  return defineConfig({
    build: {
      // minify: false,
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
        '@': fileURLToPath(new URL('./src', import.meta.url))
      },
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
