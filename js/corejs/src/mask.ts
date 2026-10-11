// v-mask="'(###) ###-####'": the input under the element writes its value by
// the mask as it is typed — # a digit, A a letter, * either, anything else
// itself; what does not fit its place is left out (the schemaform's
// [mask=…], the same as the site's). Before the input's own listeners: the
// v-model of a field reads the value written.

import type { Directive } from 'vue'

/** masked is v written by mask. */
export function masked(mask: string, v: string): string {
  const raw = (v || '').replace(/[^0-9\p{L}]/gu, '')
  let out = ''
  let i = 0
  for (let k = 0; k < mask.length && i < raw.length; k++) {
    const m = mask[k]
    const c = raw[i]
    if (m === '#' || m === 'A' || m === '*') {
      const fits = m === '#' ? /\d/.test(c) : m === 'A' ? /\p{L}/u.test(c) : true
      i++
      if (fits) out += c
      else k--
      continue
    }
    out += m
    if (c === m) i++
  }
  return out
}

const handlers = new WeakMap<HTMLElement, (ev: Event) => void>()

export const maskDirective: Directive<HTMLElement, string> = {
  mounted(el, binding) {
    const h = (ev: Event) => {
      const t = ev.target as HTMLInputElement
      const mask = binding.value
      if (!mask || !t || !('value' in t)) return
      const v = masked(mask, t.value)
      if (v !== t.value) t.value = v
    }
    handlers.set(el, h)
    el.addEventListener('input', h, true)
  },
  beforeUnmount(el) {
    const h = handlers.get(el)
    if (h) el.removeEventListener('input', h, true)
  }
}
