import {configDefaults, defineConfig, mergeConfig} from 'vitest/config'
import viteConfig from './vite.config'

// vite.config exports a function of the env ({mode}): it is called with the
// test's env, and its config merged with the test's.
export default defineConfig((env) =>
  mergeConfig(
    viteConfig(env),
    defineConfig({
      test: {
        server: {
          deps: {
            // @gad-lang/codemirror-gad ships ES modules whose imports have no
            // extension (./complete): a bundler resolves them, Node does not —
            // so it goes through vite here too, as in the build
            inline: ['vuetify', '@gad-lang/codemirror-gad'],
          },
        },
        environment: 'jsdom',
        exclude: [...configDefaults.exclude, 'e2e/*'],
      },
    }),
  ),
)
