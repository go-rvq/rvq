// Builds the CodeMirror bundle (bundle.js → rvq's js/vendor/codemirror),
// every importer given one copy of CodeMirror's core and of Lezer's — the
// workspace keeps more than one version, and two copies of
// @codemirror/state make CodeMirror refuse the extensions of the other.
//   bun src/codemirror/build.js
import { dirname, resolve } from 'node:path'

const here = dirname(new URL(import.meta.url).pathname)
const root = resolve(here, '../..')
const one = /^@(codemirror\/(state|view|language)|lezer\/(common|highlight|lr))$/

const out = await Bun.build({
  entrypoints: [resolve(here, 'bundle.js')],
  target: 'browser',
  format: 'iife',
  minify: true,
  plugins: [{
    name: 'one-codemirror',
    setup(build) {
      build.onResolve({ filter: one }, (args) => ({ path: Bun.resolveSync(args.path, root) }))
    },
  }],
})
if (!out.success) {
  console.error(out.logs.join('\n'))
  process.exit(1)
}
await Bun.write(resolve(root, '../vendor/codemirror/6/codemirror.min.js'), out.outputs[0])
console.log('codemirror.min.js', (out.outputs[0].size / 1024).toFixed(0), 'KB')
