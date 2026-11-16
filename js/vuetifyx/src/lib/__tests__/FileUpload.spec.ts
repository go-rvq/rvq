import {afterEach, expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import {createHash} from 'node:crypto'
import {Sha256} from '../upload/sha256'
import {PART_SIZE, hashFile} from '../upload/parts'
import FileUpload from '../upload/FileUpload.vue'
import {mountTemplate} from '@/lib/__tests__/testutils'

const nodeSha = (b: Uint8Array) => createHash('sha256').update(b).digest('hex')

afterEach(() => {
  vi.unstubAllGlobals()
})

it('the sha256 fed in pieces is the one of the content whole', () => {
  for (const n of [0, 1, 55, 56, 63, 64, 65, 1000, 70000]) {
    const data = new Uint8Array(n).map((_, i) => (i * 31 + 7) & 255)
    const h = new Sha256()
    // pieces of odd sizes, crossing the blocks
    for (let at = 0, step = 1; at < n; at += step, step = (step * 3) % 97 + 1) h.update(data.subarray(at, Math.min(n, at + step)))
    expect(h.hex()).toBe(nodeSha(data))
  }
})

it('a file hashed slice by slice', async () => {
  const data = new Uint8Array(PART_SIZE + 12345).map((_, i) => i & 255)
  const seen: number[] = []
  expect(await hashFile(new Blob([data]), (p) => seen.push(p))).toBe(nodeSha(data))
  expect(seen.at(-1)).toBe(1)
})

// a fake server: the events posted, the hashes it has
function fakeServer(has: Record<string, any>) {
  const events: {ev: string; form: Record<string, any>}[] = []
  vi.stubGlobal('fetch', vi.fn(async (url: string, init: any) => {
    const ev = new URL(url, 'http://x').searchParams.get('__execute_event__')!
    const form: Record<string, any> = {}
    for (const [k, v] of (init.body as FormData).entries()) form[k] = v
    events.push({ev, form})
    let data: any = null
    if (ev === 'uploadparts_Exists') data = has[form.hash] ?? null
    if (ev === 'uploadparts_Done') data = {id: form.id, name: form.name, size: Number(form.total), hash: form.hash}
    return {ok: true, json: async () => ({data})}
  }))
  return events
}

const fileOf = (text: string, name = 'a.txt') => new File([text], name)
const hashOf = (text: string) => nodeSha(new TextEncoder().encode(text))

it('a file that exists: refused without sending it (blockExisting), its link opening it', async () => {
  const events = fakeServer({[hashOf('hello')]: {id: '7', title: 'Hello file', href: '/admin/files/7', createdAt: '14/11/2026'}})
  const w = mountTemplate(FileUpload, {url: '/admin/files', blockExisting: true, modelValue: ''})
  await (w.vm as any).add([fileOf('hello')])
  await flushPromises()
  expect(events.map((e) => e.ev)).toEqual(['uploadparts_Exists'])
  const link = w.find('[data-file-upload-existing]')
  expect(link.attributes('href')).toBe('/admin/files/7')
  expect(link.text()).toBe('Hello file')
  await link.trigger('click')
  // the file, and the portal of the message: the target of what opens
  const opened = w.emitted('open')![0][0] as any
  expect(opened).toMatchObject({id: '7', title: 'Hello file', href: '/admin/files/7', createdAt: '14/11/2026'})
  expect(opened.portal).toMatch(/^vx-file-upload-/)
  expect(w.find('[data-file-upload-exists] go-plaid-portal').attributes('portal-name')).toBe(opened.portal)
  // nothing put in the form: the file not sent
  expect(w.emitted('update:modelValue')?.some((e) => e[0])).toBeFalsy()
  w.unmount()
})

it('a file that exists, elsewhere: only its hash; a new one sent in parts', async () => {
  const events = fakeServer({[hashOf('hello')]: {id: '7', title: 'Hello file', href: '/admin/files/7'}})
  const w = mountTemplate(FileUpload, {url: '/admin/files', multiple: true})
  await (w.vm as any).add([fileOf('hello'), fileOf('new content', 'b.txt')])
  await flushPromises()
  const uploaded = w.emitted('uploaded')!.map((e) => e[0])
  expect(uploaded[0]).toMatchObject({hash: hashOf('hello'), name: 'a.txt'})
  expect((uploaded[0] as any).mtime).toBeGreaterThan(0)
  expect((uploaded[1] as any).hash).toBe(hashOf('new content'))
  expect(events.map((e) => e.ev)).toEqual(['uploadparts_Exists', 'uploadparts_Exists', 'uploadparts_Part', 'uploadparts_Done'])
  const part = events.find((e) => e.ev === 'uploadparts_Part')!.form
  expect([part.offset, part.total]).toEqual(['0', String('new content'.length)])
  expect(w.emitted('done')).toBeTruthy()
  w.unmount()
})

// Several files, more than one there already: each new made a record by the
// component itself (createEvent), those there refused, not sent; the files
// and their messages kept (persist) for the component drawn again.
it('several files, some there: made by the component, kept for the next', async () => {
  const events = fakeServer({[hashOf('dup one')]: {id: '1', title: 'One'}, [hashOf('dup two')]: {id: '2', title: 'Two'}})
  const props = {url: '/admin/files', multiple: true, blockExisting: true, createEvent: 'files_CreateFromUpload', persist: 'test-key'}
  vi.useFakeTimers({toFake: ['setTimeout']})
  const w = mountTemplate(FileUpload, props)
  await (w.vm as any).add([fileOf('new a', 'a.txt'), fileOf('dup one', 'b.txt'), fileOf('new c', 'c.txt'), fileOf('dup two', 'd.txt')])
  await flushPromises()
  expect(events.filter((e) => e.ev === 'files_CreateFromUpload').length).toBe(2)
  expect(events.filter((e) => e.ev === 'uploadparts_Part').length).toBe(2)
  w.unmount()
  // drawn again: its files and messages there
  const again = mountTemplate(FileUpload, props)
  await flushPromises()
  const states = again.findAll('[data-file-upload-item]').map((e) => e.attributes('data-file-upload-item'))
  expect(states).toEqual(['done', 'exists', 'done', 'exists'])
  expect(again.findAll('[data-file-upload-existing]').map((e) => e.text())).toEqual(['One', 'Two'])
  // after 3s, the lines of the files sent gone; the refused stay
  vi.advanceTimersByTime(3000)
  await flushPromises()
  expect(again.findAll('[data-file-upload-item]').map((e) => e.attributes('data-file-upload-item'))).toEqual(['exists', 'exists'])
  vi.useRealTimers()
  again.unmount()
})
