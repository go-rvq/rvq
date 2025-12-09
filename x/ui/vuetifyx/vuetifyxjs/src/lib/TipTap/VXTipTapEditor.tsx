import { genericComponent, type GenericProps, propsFactory, useRender } from 'vuetify/lib/util/index'

import { computed, markRaw, ModelRef, PropType, ref, Ref, toValue, useId, useModel, useTemplateRef, watch } from 'vue'
import { VAlert, VCard, VCardText, VDivider, VMessages, VToolbar, VToolbarTitle } from 'vuetify/components'

import { makeComponentProps } from 'vuetify/lib/composables/component'
import { type DensityProps, makeDensityProps } from 'vuetify/lib/composables/density'
import { makeThemeProps } from 'vuetify/lib/composables/theme'
import { useI18n } from '@/lib/locale'
import Card from '@/lib/Card.vue'

import { Heading as TiptapHeading, Level as TiptapHeadingLevel } from '@tiptap/extension-heading'
import { Blockquote as TiptapBlockquote } from '@tiptap/extension-blockquote'

import Portal from './../Portal/Portal.vue'
import { unrefElement, useResizeObserver } from '@vueuse/core'

// https://github.com/go-rvq/vuetify-pro-tiptap
import {
  BaseKit,
  BlockquoteMenuItem,
  BlockTypeSelector,
  Bold,
  BulletList,
  Clear,
  Code,
  CodeBlockMenuItem,
  Color,
  defaultBubbleList,
  FontSize,
  Fullscreen,
  Highlight,
  History,
  HorizontalRule,
  Image,
  Indent,
  Italic,
  Link,
  locale as ttLocale,
  OrderedList,
  Strike,
  SubAndSuperScript,
  Table,
  TaskList,
  TextAlign,
  Underline,
  Video,
  VuetifyTiptap
} from 'vuetify-pro-tiptap'

import './VXTipTapEditor.scss'

import preview from './extensions/preview'
import help from './extensions/help'
import LinkDialog from './components/LinkDialog.vue'
import { default as ImageTab } from './components/Image'

import 'vuetify-pro-tiptap/style.css'
import 'vuetify-pro-tiptap/styles/editor.css'
import 'vuetify-pro-tiptap/styles/markdown.css'
import './styles/markdown/github.scss'

import { Editor, JSONContent } from '@tiptap/core'
import { makeInputProps } from '../util'
import { InlineTemplateCode } from '@/lib/TipTap/extensions/script/extension'
import { TemplateCodeValue } from '@/lib/TipTap/extensions/script-value/extension'
import { ToolsButton } from '@/lib/TipTap/extensions/tools-button/extension'
import { CopyAndPasteContentItem } from '@/lib/TipTap/extensions/tools-button/tools/CopyAndPasteContent'
import { CodeBlock } from '@/lib/TipTap/extensions/codeblock/types'
import { ScriptBlockMenuItemBuilder } from '@/lib/TipTap/extensions/script-block/extension'
import { ScriptBlock } from '@/lib/TipTap/extensions/script-block/types'
import { VXBtn } from '@/lib/VXBtn'

export type ModelType = string | JSONContent

export type VXTipTapEditorSlotProps<T> = {
  sorting: ModelRef<boolean>
  label: Ref<string | undefined>
  items: Ref<T[]>
  itemTitle: (item: T) => any
  itemKey: (item: T) => any
}

export const makeProps = propsFactory({
  label: { type: String, default: '' },
  placeHolder: { type: String, default: '' },
  modelValue: { default: () => '' },
  readonly: { type: Boolean, default: false },
  template: { type: Boolean, default: false },
  output: {
    type: String as PropType<'html' | 'json' | 'text'>, default: 'json'
  },
  outlined: { type: Boolean, default: true },
  dense: { type: Boolean, default: false },
  hideToolbar: { type: Boolean, default: false },
  disableToolbar: { type: Boolean, default: false },
  maxWidth: { type: Number, default: 900 },
  maxHeight: { type: Number, default: 500 },
  readonlyClass: { type: String, default: '' },
  headings: { type: Array<TiptapHeadingLevel>, default: [1, 2, 3, 4, 5, 6] },
  viewerProps: { type: Object, default: () => ({}) },
  ...makeComponentProps(),
  ...makeDensityProps(),
  ...makeThemeProps(),
  ...makeInputProps()
}, 'VXTipTapEditor')

type ItemType<T> = T extends readonly (infer U)[] ? U : never

export type VXTipTapEditorSlots<T> = {
  help: Object
  preview: Object
  viewer: Object
  message: Object
  details: Object
  'select-image': Object
}

const printScript = `this.classList.remove('overflow-y-auto');this.style.maxHeight='inherit'`

export interface PI extends DensityProps {
  label?: String
  placeHolder?: String
  hint?: String
  modelValue?: ModelType
  readonly?: Boolean
  output?: string
  outlined?: Boolean
  dense?: Boolean
  hideToolbar?: Boolean
  disableToolbar?: Boolean
  errorMessages?: string | string[] | null
  maxWidth?: Number
  maxHeight?: Number
  headings?: TiptapHeadingLevel[]
  viewerProps?: Object
  'onUpdate:modelValue'?: (value: ModelType) => void
}

export interface Exposed {
  vuetifyTiptap: typeof VuetifyTiptap;
}

export const VXTipTapEditor = genericComponent<new <T extends readonly any[]>(
  props: PI,
  slots: VXTipTapEditorSlots<ItemType<T>>
) => GenericProps<typeof props, typeof slots>>()({
  name: 'VXTipTapEditor',

  props: makeProps(),

  emits: {
    'update:modelValue': (value: any) => true,
    'editor': (value: Editor) => true
  },

  setup(props, { expose, slots, emit }) {
    const { t, locale } = useI18n()

    ttLocale.setLang((locale.value as any) as string)

    const imageTabs: any[] = [
      { name: t('home'), component: markRaw(ImageTab as Object) }
    ]

    if (slots['select-image']) {
      imageTabs.push({
        name: t('select'),
        component: markRaw(slots['select-image'] as Object)
      })
    }

    const toSize = (str: any) => {
        const s = '' + str
        return /^\d+$/.test(s) ? s + 'px' : s
      },
      model = useModel(props, 'modelValue'),
      helpVisible = ref(false),
      helpExpanded = ref(false),
      extensions = [
        BaseKit.configure({
          placeholder: {
            placeholder: props.placeHolder
          },
          bubble: {
            defaultBubbleList: editor => {
              // You can customize the bubble menu here
              return defaultBubbleList(editor) // default customize bubble list
            }
          }
        }),
        BlockTypeSelector.configure({
          headingLevels: props.headings,
          items: ({ editor, t: tt }) => [
            BlockquoteMenuItem({ editor, t: tt }),
            CodeBlockMenuItem({ editor, t: tt }),
            ...(props.template ? [ScriptBlockMenuItemBuilder(editor, t)] : [])
          ]
        }),
        Bold,
        Italic,
        Underline,
        Strike,
        Code.configure({
          divider: true
        }),
        ...(props.template ? [
          InlineTemplateCode.configure({}),
          TemplateCodeValue.configure({ divider: true })
        ] : []),
        TiptapHeading,
        CodeBlock,
        TiptapBlockquote,
        ScriptBlock,
        TextAlign,
        FontSize,
        Color,
        Highlight.configure({
          divider: true
        }),
        SubAndSuperScript.configure({
          divider: true
        }),
        Clear.configure({
          divider: true
        }),
        BulletList,
        OrderedList,
        TaskList,
        Indent.configure({
          divider: true
        }),
        Link.configure({
          dialogComponent: () => LinkDialog
        }),
        Image.configure({
          imageTabs: imageTabs,
          width: 500,
          hiddenTabs: ['url', 'upload'],
          allowBase64: true,
          upload(file: File) {
            const url = URL.createObjectURL(file)
            return Promise.resolve(url)
          }
        }),
        Video,
        Table.configure({
          divider: true,
          cellMinWidth: 100
        }),
        HorizontalRule,
        // CodeBlock.configure({ divider: true }),
        History.configure({ divider: true }),
        preview.configure({ spacer: true, component: slots['preview'] }),
        Fullscreen.configure({
          useWindow: true
        }),
        ...(slots.help ? [help.configure({
          modelValue: helpVisible
        })] : []),
        ToolsButton.configure({
          items: [CopyAndPasteContentItem()]
        })
      ],
      messages = computed(() => {
        if (props.errorMessages?.length) {
          return props.errorMessages
        } else if (props.hint) {
          return props.hint
        } else {
          return props.messages
        }
      }),
      uid = useId(),
      id = computed(() => props.id || `input-${uid}`),
      messagesId = computed(() => `${id.value}-messages`),
      slotProps = computed(() => ({
        id,
        messagesId,
        density: props.density
      })),
      vuetifyTiptap = useTemplateRef<any>('vuetifyTiptap'),
      editor = computed(() => vuetifyTiptap.value.editor),
      editorHeight = ref(0)

    useResizeObserver(vuetifyTiptap, (entries) => {
      const entry = entries        [0]
      editorHeight.value = entry.contentRect.height
    })

    const targets = computed(() => {
      return unrefElement(toValue(vuetifyTiptap))
    })

    const stopWatch = watch(
      targets,
      (el: HTMLElement) => {
        if (el) {
          editorHeight.value = el.offsetHeight
        }
      },
      { immediate: true, flush: 'post' }
    )

    useRender(() => {
      const hasMessages = messages.value.length > 0
      const hasDetails = hasMessages
      const isError = props.error || props.errorMessages?.length
      const color = isError ? 'error' : undefined
      const messagesComp = () => (<>
          <VMessages
            active={hasMessages}
            messages={messages.value}
            v-slots={{ message: slots.message }}
          />

          {slots.details?.(slotProps.value)}
        </>
      )
      const detailComp = hasDetails && (
          (props.errorMessages?.length ? <><VAlert
            density="compact"
            icon="$error"
            type="error"
            variant="tonal"
            class="rounded-0"
          >
            {messagesComp()}
          </VAlert></> : <>
            <div
              id={messagesId.value}
              class="v-input__details py-2 px-3"
              role="alert"
              aria-live="polite"
            >
              {messagesComp()}
            </div>
          </>)
        ),
        editorComp = (
          <VuetifyTiptap
            class={props.label ? 'rounded-t-0' : (detailComp ? 'rounded-b-0' : undefined)}
            model-value={model.value}
            onUpdate:modelValue={(v) => emit('update:modelValue', v)}
            output={props.output}
            dense={props.dense}
            max-height={props.maxHeight}
            max-width={props.maxWidth}
            extensions={extensions}
            ref="vuetifyTiptap"
            style="width:100%"
          />),
        helpBtn = (title:string, icon:string, action:()=>void) => <VXBtn class="rounded me-1 ms-0" density="comfortable" size="small" icon={'mdi-'+icon} title={t(title)} variant="text" {...{onClick:action}} />

      const defaultComp = (ro: boolean) => {
        return <>
          <VCard color={isError ? 'error' : undefined} class={['border v-input']}
                 density={props.density} {...(ro ? props.viewerProps : {})}>
            {props.label && <>
              <VToolbar color={color} density={props.density} height="36">
                <VToolbarTitle style="margin-inline-start: 0" class="px-3">{props.label}</VToolbarTitle>
              </VToolbar>
              {detailComp && (
                <div class={['v-toolbar', { 'd-block': isError }]}>{detailComp}</div>
              )}
              <VDivider />
            </>}
            {ro ? <VCardText data-print-script={printScript}
                             class={['overflow-y-auto', props.readonlyClass, { 'bg-white': isError }]}
                             style={props.maxHeight ? { maxHeight: toSize(props.maxHeight) } : {}}>
                <Portal model-value={model.value} />
              </VCardText>
              : slots.help ? (
                <>
                  <div class="d-flex">
                    {helpExpanded.value ? undefined : (<>
                      {editorComp}
                      <VDivider vertical />
                    </>)}
                    {editorHeight.value && helpVisible.value ? (
                        <div class="d-flex flex-column" style={`height: ${editorHeight.value}px;width:${helpExpanded.value?'100%':'40%'}`}>
                            <VToolbar class="py-1" density="compact" height="24px" v-slots={{
                              append: () => (<>
                                {helpExpanded.value ?
                                  helpBtn('collapse' , 'arrow-collapse', () => helpExpanded.value = false):
                                  helpBtn('expand', 'arrow-expand', () => helpExpanded.value = true)
                                }
                                {helpBtn('close', 'close', () => {
                                  helpVisible.value = false
                                  helpExpanded.value = false
                                })}
                              </>),
                              default: () => (<VToolbarTitle class="text-subtitle-1">{t('help')}</VToolbarTitle>)
                            }}>
                            </VToolbar>
                          <VDivider />
                            <VCardText class="overflow-y-auto">
                              {slots.help({editor:editor.value})}
                            </VCardText>
                        </div>)
                      : undefined}
                  </div>
                </>
              ) : editorComp}
            {!props.label && <>
              <VDivider />
              {detailComp && (
                <div class="bg-grey-lighten-3">{detailComp}</div>
              )}
            </>}
          </VCard>
        </>
      }

      return <div class={[props.class, 'vx-tiptap-editor', { 'vx-tiptap-editor--error': isError }]} style={props.style}>
        {props.readonly ? (slots.viewer ? slots.viewer({
          maxHeight: props.maxHeight,
          value: model.value, ...props.viewerProps
        }) : defaultComp(true)) : defaultComp(false)}
      </div>
    })

    expose({
      vuetifyTiptap
    })

    return {}
  }
})

export type VXTipTapEditor = InstanceType<typeof VXTipTapEditor>

export default VXTipTapEditor
