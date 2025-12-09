import {
  genericComponent,
  type GenericProps,
  getPropertyFromItem,
  propsFactory,
  type SelectItemKey, useRender
} from 'vuetify/lib/util/index'

import { computed, markRaw, ModelRef, PropType, ref, Ref, toRaw, toRef, useModel } from 'vue'
import { VAvatar, VCard, VList, VListItem, VSpacer, VToolbar, VToolbarTitle } from 'vuetify/components'

import { makeComponentProps } from 'vuetify/lib/composables/component'
import { makeDensityProps } from 'vuetify/lib/composables/density'
import { makeThemeProps } from 'vuetify/lib/composables/theme'
import { useI18n } from './locale/index'
import { VXBtn } from './VXBtn'
import { VueDraggableNext } from 'vue-draggable-next'

export type VXArraySorterSlotProps<T> = {
  sorting: ModelRef<boolean>
  label: Ref<string|undefined>
  items: Ref<T[]>
  itemTitle: (item: T) => any
  itemKey: (item: T) => any
}

export type VXArraySorterItemSlotProps<T> = {
  global: VXArraySorterSlotProps<T>
  item: T
  itemIndex: number
}

export const makeSorterProps = propsFactory({
  label: String,
  modelValue: { type: Array<any>,   default: () => [] },
  sorting: { type: Boolean, default: false },
  readonly: { type: Boolean, default: false },
  itemKey: {    type: [String, Array, Function] as PropType<SelectItemKey>,    default: 'key'  },
  itemTitle: {    type: [String, Array, Function] as PropType<SelectItemKey>,    default: 'title'  },

  ...makeComponentProps(),
  ...makeDensityProps(),
  ...makeThemeProps()
}, 'VXArraySorter')

type ItemType<T> = T extends readonly (infer U)[] ? U : never

export type VXArraySorterSlots<T> = {
  default: VXArraySorterSlotProps<T>
  sort: VXArraySorterSlotProps<T>
  header: VXArraySorterSlotProps<T>
  item: VXArraySorterItemSlotProps<T>
}

export const VXArraySorter = genericComponent<new <T extends readonly any[]>(
  props: {
    modelValue?: T
    sorting?: Boolean
    readonly ?: Boolean
    label?: String
    itemTitle?: SelectItemKey<ItemType<T>>
    itemKey?: SelectItemKey<ItemType<T>>
    'onUpdate:modelValue'?: (value: T) => void
    'onUpdate:sorting'?: (value: Boolean) => void
  },
  slots: VXArraySorterSlots<ItemType<T>>
) => GenericProps<typeof props, typeof slots>>()({
  name: 'VXArraySorter',

  props: makeSorterProps(),

  emits: {
    'update:modelValue': (value: any[]) => true,
    'update:sorting': (value: boolean) => true,
    'sorted': (value: any[]) => true
  },

  setup(props, { slots, emit }) {
    const { t } = useI18n(),
      items = useModel(props, 'modelValue'),
      sorting = useModel(props, 'sorting'),
      isSorting = computed(() => sorting.value),
      toggle = () => {
        sorting.value = !sorting.value
        return sorting.value
      },
      label = toRef(props, 'label'),
      cItems = computed({
        get() {
          return items.value
        },
        set(v: any[]) {
          toRaw(v).forEach((e, i) => {
            items.value[i] = e
          })
        }
      }),
      itemTitle = (item: any) => getPropertyFromItem(item, props.itemTitle),
      itemKey = (item: any) => getPropertyFromItem(item, props.itemKey),
      slotProps:VXArraySorterSlotProps<any> = {
        label,
        sorting,
        itemTitle,
        itemKey,
        items: cItems
      },
      history: Ref<Array<any[]>> = ref([]),
      action = (e: string) => {
        if (e === 'edit') history.value.unshift(items.value)
        if (e === 'undo') {
          if (history.value.length > 0) {
            (history.value.shift() as any[]).forEach((e, i) => {
              items.value[i] = e
            })
          }
          return
        }
        if (e === 'save') {
          const result = history.value[0].map((e) => e)
          cItems.value = result
          history.value.splice(0, history.value.length)
          emit('update:modelValue', result)
        }
        toggle()
      },
      onDragEnd = (event: any) => {
        let raw = history.value[0].map((e) => e)
        const movedItem = raw.splice(event.oldIndex, 1)[0]
        raw.splice(event.newIndex, 0, movedItem)
        history.value.unshift(raw)
      },
      defaultItemComp = (itemProps:VXArraySorterItemSlotProps<any>) => {
        return <VListItem key={itemKey(itemProps.item)} title={itemTitle(itemProps.item)}>
          {{
            prepend: () => {
              return <>
                <VXBtn class="me-2" title={t('move')} density={props.density} icon="mdi-drag" variant="text" />
                <VAvatar density={props.density} color="secondary">{itemProps.itemIndex + 1}</VAvatar>
              </>
            }
          }}
        </VListItem>
      },
      itemComp = (item: any, itemIndex: number) => {
        const props: VXArraySorterItemSlotProps<any> = {
          item,
          itemIndex,
          global: slotProps
        }
        if (slots.item) {
          const newProps:any = { props, defaultItemComp: markRaw(defaultItemComp) }
          return slots.item?.(newProps)
        }
        return defaultItemComp(props)
      }

    useRender(() => {
      const tag:any = (() => {
        return (
          <VList tile density={props.density}>
            {items.value.map(itemComp)}
          </VList>
        )
      })

      return <VCard class={props.class} style={props.style} density={props.density}>
        <VToolbar density={props.density}>
          {isSorting.value ? <VXBtn title={t('back')} density={props.density} variant="text" icon="mdi-arrow-left"
            { ...{onClick: () => action('close')}} /> : undefined}
          {label.value && (<VToolbarTitle>{label.value}</VToolbarTitle>)}
          <VSpacer />
          {isSorting.value ? (<>
            <VXBtn title={t('apply')} density={props.density} class="me-3" variant="text" color="primary"
                   icon="mdi-check"
                   { ...{onClick: () => action('save')}} />
            {history.value.length > 1 ?
              <VXBtn title={t('undo')} density={props.density} variant="text" icon="mdi-undo"
                     { ...{onClick: () => action('undo')}} /> : undefined}

          </>) : props.readonly ? undefined : <VXBtn density={props.density} title={t('changeOrder')} variant="text" icon="mdi-sort"
                        { ...{onClick: () => action('edit')}} />}
        </VToolbar>
        {sorting.value ? (
          <VueDraggableNext {... {
            items,
            itemKey:props.itemKey,
            handle:".v-list-item__prepend > button",
            onEnd:onDragEnd,
            tag: tag as string,
            ghostClass:"ghost"
          }} />
        ) : (slots.default?.(slotProps) || (
          <VList density={props.density}>
            {items.value.map((item) => {
              return <VListItem title={itemTitle(item)} />
            })}
          </VList>
        ))}
      </VCard>
    })

    return {}
  }
})

export type VXArraySorter = InstanceType<typeof VXArraySorter>

export default VXArraySorter
