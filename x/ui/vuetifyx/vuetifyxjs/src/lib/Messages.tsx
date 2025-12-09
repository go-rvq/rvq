import { computed, CSSProperties, defineComponent, ExtractPublicPropTypes, h, toRaw } from 'vue'
import { useTextColor } from 'vuetify/lib/composables/color'
import { VAlert, VIcon } from 'vuetify/components'
import { useRender } from 'vuetify/lib/util/index'
import { makeComponentProps } from 'vuetify/lib/composables/component'
import Dialog from './Dialog.vue'
import { VXBtn } from './VXBtn'
import { useI18n } from './locale'
import { componentByTemplate, ComponentDefinition } from './component-by-template'

export type Detail = string | ComponentDefinition

export interface Message {
  type?: string
  value: string
  detail?: Detail
}

export interface MessageColor {
  classes: string[]
  styles: CSSProperties
  icon: string
}

export const buildMessagesColor = (messages: Message[]) => {
  return messages.map(message => {
    const gtype = computed(() => message.type)
    const { textColorClasses: classes, textColorStyles: styles } = useTextColor(gtype)

    return {
      classes: classes.value,
      styles: styles.value
    } as MessageColor
  })
}

export type Slots = {
  message: Object
}

const propsOptions = {
  items: {
    type: Array,
    required: true
  },

  error: Boolean,
  ...makeComponentProps()
} as const

export type Props = ExtractPublicPropTypes<typeof propsOptions>

export default defineComponent({
  name: 'VXMessages',

  props: propsOptions,

  emits: {
    'update:modelValue': (value: string) => true
  },

  setup(props: Props, { slots: Slots }) {
    const { t } = useI18n(),
      messages = (props.items ?? []) as Message[],
      detailComp = (detail:Detail) => {
        if (typeof detail == 'string') {
          return (<div innerHTML={detail as string}></div>)
        }
        detail = toRaw(detail as any) as ComponentDefinition
        console.log("detail", detail)
        const Comp = componentByTemplate(detail)
        return (<><Comp /></>)
      }

    useRender(() => {
      if (!messages.length) {
        return <></>
      }

      const messagesColor = buildMessagesColor(messages),
        hasError = messages.filter(message => message.type == 'error').length > 0,
        messagesComp = () => (<div class="v-messages">
            {messages.map((message, i) =>
              <div class={['v-messages__message', ...messagesColor[i].classes]}
                   style={messagesColor[i].styles}
                   key={`${i}-${message.value}`}
              >
                <VIcon icon={'$' + message.type} />
                {message.value}
                { message.detail ? (<Dialog density='compact' closable expandable v-slots=
                  {{
                    activator: (arg:any) => (<VXBtn
                      title={t('showDetail')}
                      size="small"
                      density="compact"
                      variant='text'
                      class="ms-2"
                      icon="mdi-open-in-new"
                      {...arg.props} />),
                    body: () => {
                      let Comp;
                      if (typeof message.detail == "string" ) {
                        Comp = componentByTemplate({ template: message.detail as string })
                      } else {
                        Comp = componentByTemplate(message.detail as ComponentDefinition)
                      }
                      return (<div>
                        {h(Comp)}
                      </div>)
                    }
                  }} />) : undefined }
              </div>
            )}
          </div>
        )

      return <div
        class={['vx-messages', props.class]}
        style={props.style}>
        {hasError ?
          <VAlert
            density="compact"
            type="error"
            variant="tonal"
            class="rounded-0"
          >
            {messagesComp()}
          </VAlert> :
          <div
            class="v-input__details py-2 px-3"
            role="alert"
            aria-live="polite"
          >
            {messagesComp()}
          </div>
        }
      </div>
    })
    return {}
  }
})
