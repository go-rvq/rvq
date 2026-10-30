// bun test: a block of code made the code viewer's (its language, label,
// copy words, host) before Prism highlights it.
import { expect, test } from 'bun:test'
import { Window } from 'happy-dom'
import { labelOf, languageOf, prepare } from './code.js'

test('the language of a code element, and its name', () => {
  expect(languageOf({ className: 'language-Go' })).toBe('go')
  expect(languageOf({ className: 'x language-c++ y' })).toBe('c++')
  expect(languageOf({ className: '' })).toBe('')
  expect(labelOf('yml')).toBe('YAML')
  expect(labelOf('rust')).toBe('rust')
})

test('a block prepared: in a host, its pre the viewer\'s, an unknown language none', () => {
  const window = new Window()
  globalThis.document = window.document
  const doc = window.document
  doc.body.innerHTML = '<div><pre><code class="language-go">x := 1\n</code></pre><pre><code class="language-rust">fn x(){}</code></pre></div>'
  const Prism = { languages: { go: {} } }
  const words = { copy: 'Copiar', copied: 'Copiado!', copyError: 'Ctrl+C' }
  const [go, rust] = doc.querySelectorAll('code')
  const pre = prepare(go, Prism, words)
  expect(pre.parentElement.className).toBe('mdkit-code-host')
  expect(pre.classList.contains('mdkit-code')).toBe(true)
  expect(pre.classList.contains('line-numbers')).toBe(true)
  expect(pre.getAttribute('data-language')).toBe('Go')
  expect(pre.getAttribute('data-prismjs-copy')).toBe('Copiar')
  expect(go.className).toBe('language-go')
  expect(go.textContent).toBe('x := 1')
  prepare(rust, Prism, words)
  expect(rust.className).toBe('language-none')
  expect(rust.parentElement.getAttribute('data-language')).toBe('rust')
  delete globalThis.document
})
