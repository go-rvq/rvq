import {expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import ImageSource from '../TipTap/components/ImageSource.vue'
import {mountTemplate} from '@/lib/__tests__/testutils'

const png = () => new File([new Uint8Array([137, 80, 78, 71])], 'a.png', {type: 'image/png'})

it('an image taken: sent by upload, its URL the value', async () => {
  const upload = vi.fn(async () => '/system/a.png')
  const w = mountTemplate(ImageSource, {modelValue: '', upload})
  await (w.vm as any).take(png())
  await flushPromises()
  expect(upload).toHaveBeenCalled()
  expect(w.emitted('update:modelValue')).toEqual([['/system/a.png']])
  w.unmount()
})

it('no upload: a data URL; not an image: nothing', async () => {
  const w = mountTemplate(ImageSource, {modelValue: ''})
  await (w.vm as any).take(new File(['x'], 'a.txt', {type: 'text/plain'}))
  await (w.vm as any).take(png())
  await flushPromises()
  const got = w.emitted('update:modelValue') as string[][]
  expect(got).toHaveLength(1)
  expect(got[0][0]).toMatch(/^data:image\/png;base64,/)
  w.unmount()
})
