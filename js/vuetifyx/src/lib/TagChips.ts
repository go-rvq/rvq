// <vx-tag-chips>: tags as chips. Each tag of items is a chip: one not chosen,
// outlined with a "+", chosen at a click; one chosen, filled in its color, a
// trash at its end that takes it off (only the trash: a click on the chip
// itself does nothing, so it is not taken off by mistake); its hint (who
// chose it, when, where) a tooltip. v-model is the values chosen; add and
// remove say each change (a parent that saves each one by itself listens to
// them). group shows only the items of that group (the tags of a type).
import { computed, defineComponent, h, type PropType, resolveComponent } from 'vue'
import { useVxMessages } from '@/lib/locale'

export interface TagChipItem {
  value: string
  title: string
  color?: string
  // the tooltip of the chip chosen (who chose it, when, where)
  hint?: string
  // the group the item is of (a type's tags): shown when group is it
  group?: string
}

export default defineComponent({
  name: 'VxTagChips',
  props: {
    items: { type: Array as PropType<TagChipItem[]>, default: () => [] },
    modelValue: { type: Array as PropType<string[]>, default: () => [] },
    group: { type: String, default: undefined },
    readonly: { type: Boolean, default: false },
    label: { type: String, default: '' },
    // nothing to choose from
    empty: { type: String, default: '' },
    // the item being saved (a parent saving each change): its chip waits
    busy: { type: String, default: '' }
  },
  emits: ['update:modelValue', 'add', 'remove'],
  setup(props, { emit }) {
    const vx = useVxMessages()
    const shown = computed(() =>
      props.items.filter((it) => props.group === undefined || props.group === null || it.group === props.group)
    )
    const chosen = computed(() => new Set((props.modelValue || []).map(String)))
    const add = (v: string) => {
      if (props.readonly || chosen.value.has(v)) return
      emit('update:modelValue', [...(props.modelValue || []), v])
      emit('add', v)
    }
    const remove = (v: string) => {
      if (props.readonly || !chosen.value.has(v)) return
      emit(
        'update:modelValue',
        (props.modelValue || []).filter((x) => String(x) !== v)
      )
      emit('remove', v)
    }
    return () => {
      const VChip = resolveComponent('v-chip')
      const VTooltip = resolveComponent('v-tooltip')
      const chips = shown.value.map((it) => {
        const v = String(it.value)
        const on = chosen.value.has(v)
        const busy = props.busy === v
        if (!on) {
          return h(
            VChip,
            {
              key: v,
              variant: 'outlined',
              size: 'small',
              prependIcon: busy ? 'mdi-loading mdi-spin' : 'mdi-plus',
              disabled: props.readonly || busy,
              class: 'vx-tag-chips__chip ma-1',
              'data-tag': v,
              'data-tag-on': 'false',
              onClick: () => add(v)
            },
            () => it.title
          )
        }
        const chip = (extra: Record<string, any> = {}) =>
          h(
            VChip,
            {
              key: v,
              variant: 'flat',
              size: 'small',
              color: it.color || 'primary',
              closable: !props.readonly,
              closeIcon: busy ? 'mdi-loading mdi-spin' : 'mdi-trash-can-outline',
              'close-label': vx.value.tagRemove,
              class: 'vx-tag-chips__chip vx-tag-chips__chip--on ma-1',
              'data-tag': v,
              'data-tag-on': 'true',
              'onClick:close': (e: Event) => {
                e?.stopPropagation?.()
                remove(v)
              },
              ...extra
            },
            () => it.title
          )
        if (!it.hint) return chip()
        return h(
          VTooltip,
          { key: v, text: it.hint, location: 'bottom' },
          { activator: ({ props: p }: any) => chip(p) }
        )
      })
      return h('div', { class: 'vx-tag-chips' }, [
        props.label ? h('div', { class: 'text-caption text-medium-emphasis mb-1' }, props.label) : null,
        chips.length
          ? h('div', { class: 'd-flex flex-wrap align-center ma-n1' }, chips)
          : h('div', { class: 'text-body-2 text-medium-emphasis' }, props.empty || vx.value.tagsNone)
      ])
    }
  }
})
