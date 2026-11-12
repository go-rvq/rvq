import { describe, expect, it } from 'vitest'
import { nextTick } from 'vue'

import PermissionInput from '../PermissionInput.vue'
import { mountTemplate } from '@/lib/__tests__/testutils'

// the viewport of the overlays (VDialog's): none in jsdom
if (!(globalThis as any).visualViewport) {
  ;(globalThis as any).visualViewport = {
    width: 1024,
    height: 768,
    offsetLeft: 0,
    offsetTop: 0,
    scale: 1,
    addEventListener() {},
    removeEventListener() {}
  }
}

const flush = async () => {
  for (let i = 0; i < 5; i++) await nextTick()
  await new Promise((r) => setTimeout(r, 0))
}

const tree = (value: string) => [
  {
    title: 'Admin',
    kind: 'scope',
    value: 'admin:',
    edit: 'admin:',
    canAdd: true,
    expand: true,
    children: value
      ? [
          {
            title: 'List',
            kind: 'verb',
            value: 'admin:posts:@list',
            edit: 'admin:posts:@list',
            explicit: true,
            prefix: 'grant'
          }
        ]
      : []
  }
]

// a loader of the server's: what it was asked, a value of the edits made
function server(calls: any[]) {
  return (op: string, p: Record<string, any>) => {
    calls.push({ op, ...p })
    switch (op) {
      case 'children':
        return {
          records: true,
          choices: [
            { segment: '@list', resource: 'admin:posts:@list', kind: 'verb', title: 'List' },
            { segment: '@create', resource: 'admin:posts:@create', kind: 'verb', title: 'Create' }
          ]
        }
      case 'records':
        return [{ id: '7', label: 'Primeiro <post>' }]
      case 'merge':
        if (p.text.includes('{')) return { error: 'at 7: a group not closed: a } is missing' }
        if (p.text.startsWith('unset'))
          return {
            value: p.value,
            items: tree(p.value),
            warning: 'unset admin:x:* took nothing out'
          }
        return { value: 'admin:posts:{@list, @get}', items: tree('x') }
      case 'tree':
        return {
          value: p.value,
          items: tree(p.value),
          formatted: p.value ? 'admin:{\n  posts:@list\n}' : ''
        }
      default: {
        const value = op === 'remove' ? '' : 'admin:posts:{@list, @create}'
        return { value, items: tree(value) }
      }
    }
  }
}

describe('vx-permission-input', () => {
  it('asks the tree of its value, of the view and the groups', async () => {
    const calls: any[] = []
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server(calls),
      stringView: true
    })
    await flush()
    expect(calls[0]).toMatchObject({
      op: 'tree',
      value: 'admin:posts:@list',
      view: 'verb',
      groups: ''
    })
    const vm = w.vm as any
    expect(vm.items[0].children[0].title).toBe('List')
    vm.st.groups = true
    vm.st.view = 'kind'
    await flush()
    expect(calls.at(-1)).toMatchObject({ op: 'tree', view: 'kind', groups: '1', only: '' })
    // the icons of the bar: only what is defined
    await w.find('[data-permission-only]').trigger('click')
    await flush()
    expect(calls.at(-1)).toMatchObject({ op: 'tree', groups: '1', only: '1' })
    // the panel hidden: its icon alone, shown again by it
    await w.find('[data-permission-string-hide]').trigger('click')
    await flush()
    expect(w.find('[data-permission-string]').exists()).toBe(false)
    await w.find('[data-permission-string-show]').trigger('click')
    await flush()
    expect(w.find('[data-permission-string]').exists()).toBe(true)
    // the value to read: the server's, indented
    expect(w.find('[data-permission-line]').element.textContent).toBe('admin:{\n  posts:@list\n}')
  })

  it('says the prefixes in the language of the page', async () => {
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server([]),
      messages: { grant: 'CONCEDE', deny: 'NEGA' }
    })
    await flush()
    const vm = w.vm as any
    expect(vm.prefixLabel('grant')).toBe('CONCEDE')
    expect(vm.prefixLabel('deny')).toBe('NEGA')
    expect(vm.prefixLabel('allow')).toBe('ALLOW')
    expect(vm.prefixLabel('')).toBe('—')
  })

  it('marks the rows of something defined in them or under them', async () => {
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server([])
    })
    await flush()
    // the scope: a line under it (expand); its button of the primary color
    const btns = w.findAll('[data-permission-expand]')
    expect(btns.length).toBeGreaterThan(0)
    expect(btns[0].attributes('data-permission-defined')).toBe('true')
    expect(btns[0].classes()).toContain('text-primary')
  })

  it('shows the permission as text only when asked', async () => {
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server([])
    })
    await flush()
    expect(w.find('[data-permission-string]').exists()).toBe(false)
  })

  it('edits by the server: a prefix, out, parts added', async () => {
    const calls: any[] = []
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server(calls)
    })
    await flush()
    const vm = w.vm as any
    const list = vm.items[0].children[0]
    await vm.setPrefix(list, 'deny')
    expect(calls.at(-1)).toMatchObject({
      op: 'set',
      resources: '["admin:posts:@list"]',
      prefix: 'deny'
    })
    await vm.setPrefix(list, '')
    expect(calls.at(-1)).toMatchObject({ op: 'set', prefix: 'none' })
    await vm.remove(vm.items[0])
    expect(calls.at(-1)).toMatchObject({
      op: 'remove',
      resources: '["admin:","admin:posts:@list"]'
    })
    expect(w.emitted('update:modelValue')!.at(-1)).toEqual([''])
  })

  it('adds what the dialog offers, records with their labels', async () => {
    const calls: any[] = []
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server(calls)
    })
    await flush()
    const vm = w.vm as any
    await vm.openAdd({ title: 'Posts', edit: 'admin:posts:', children: vm.items[0].children })
    await flush()
    expect(calls.find((c) => c.op === 'children')).toMatchObject({
      path: 'admin:posts:',
      groups: ''
    })
    // there already: checked, not chosen again
    expect(vm.dlg.present.has('admin:posts:@list')).toBe(true)
    vm.toggle('admin:posts:@list')
    vm.toggle('admin:posts:@create')
    vm.toggleRecord({ id: '7', label: 'Primeiro <post>' })
    await vm.add()
    expect(calls.at(-1)).toMatchObject({
      op: 'add',
      prefix: 'grant',
      resources: JSON.stringify(['admin:posts:@create', 'admin:posts:<7<Primeiro \\<post\\>>>:'])
    })
    expect(w.emitted('update:modelValue')!.at(-1)).toEqual(['admin:posts:{@list, @create}'])
  })

  it('adds lines from a text, or says what is wrong with it', async () => {
    const calls: any[] = []
    const w = mountTemplate(PermissionInput, {
      modelValue: 'admin:posts:@list',
      loader: server(calls)
    })
    await flush()
    const vm = w.vm as any
    expect(w.find('[data-permission-text]').exists()).toBe(true)
    vm.txt.text = 'admin:{posts'
    await vm.applyText()
    expect(calls.at(-1)).toMatchObject({
      op: 'merge',
      text: 'admin:{posts',
      value: 'admin:posts:@list'
    })
    expect(vm.txt.error).toContain('a group not closed')
    expect(vm.txt.text).toBe('admin:{posts')
    vm.txt.text = 'admin:posts:@get'
    await vm.applyText()
    expect(w.emitted('update:modelValue')!.at(-1)).toEqual(['admin:posts:{@list, @get}'])
    expect(vm.txt.text).toBe('')
    expect(vm.txt.error).toBe('')
    // an unset of nothing: said
    vm.txt.text = 'unset admin:x:*'
    await vm.applyText()
    await flush()
    expect(w.find('[data-permission-text-warning]').text()).toContain('took nothing out')
    // off by its prop
    const off = mountTemplate(PermissionInput, {
      modelValue: '',
      loader: server([]),
      textInput: false
    })
    await flush()
    expect(off.find('[data-permission-text]').exists()).toBe(false)
  })

  it('of a key: deny only', async () => {
    const calls: any[] = []
    const w = mountTemplate(PermissionInput, {
      modelValue: '',
      loader: server(calls),
      prefixes: ['deny']
    })
    await flush()
    const vm = w.vm as any
    await vm.openAdd(null)
    await flush()
    vm.toggle('admin:posts:@list')
    await vm.add()
    expect(calls.at(-1)).toMatchObject({ op: 'add', prefix: 'deny' })
  })
})
