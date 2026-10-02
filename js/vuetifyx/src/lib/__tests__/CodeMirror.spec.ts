import {expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'

import CodeMirror from '../CodeMirror'
import {mountTemplate} from '@/lib/__tests__/testutils'

// The Gad languages load their grammar (@gad-lang/codemirror-gad); a name the
// registry has not is plain text.
it.each([
  ['gad', {}],
  ['gadx', {}],
  ['gadt', {}],
  ['gadt', {templateStart: '{', templateEnd: '}'}],
])('CodeMirror %s %o loads its grammar', async (language, delims) => {
  const wrapper = mountTemplate(CodeMirror, {language, modelValue: '', ...delims})
  // the grammar is imported asynchronously
  await vi.waitFor(() => expect((wrapper.vm as any).resolvedLang.length).toBeGreaterThan(0), {timeout: 5000})
})

it('CodeMirror without a grammar is plain text', async () => {
  const wrapper = mountTemplate(CodeMirror, {language: 'nope', modelValue: ''})
  await flushPromises()
  expect((wrapper.vm as any).resolvedLang.length).toBe(0)
})
