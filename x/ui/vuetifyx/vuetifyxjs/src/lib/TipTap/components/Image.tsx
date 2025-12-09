import { computed, defineComponent, ExtractPublicPropTypes, PropType, useModel } from 'vue'
import { ImageProperties } from 'vuetify-pro-tiptap'
import Input from '../../ImageTools/Input.vue'
import { VTextField } from 'vuetify/components'

const propsOptions = {
  modelValue: {
    type: Object as PropType<Record<string, string>>,
    required: true
  },

  upload: {
    type: Function as PropType<(file: File) => Promise<string>>,
    required: true
  },
  t : {
    type: Function as PropType<(path: string) => string>,
    required: true
  }
} as const

export type Props = ExtractPublicPropTypes<typeof propsOptions>

export default defineComponent({
  name: 'VXTipTapComponentImage',

  props: propsOptions,

  emits: {
    'update:modelValue': (value: Record<string, string>) => true
  },

  setup(props:Props, { emit }) {
    const t = props.t,
      form = useModel(props, 'modelValue'),
      isBinary = computed(() => form.value.src && form.value.src.substring(0, 5) == 'data:'),
      link = computed({
        get() {
          return isBinary.value ? 'binary' : form.value.src
        },
        set(val: string) {
          if (val !== 'binary') {
            form.value.src = val
          }
        }
      })

    return () => <>
      <div>
        <Input
          class="mb-3 mx-auto"
          link
          upload
          paste
          hint="Teste"
          initial-value={form.value.src}
          v-model={form.value.src}
          style="width:fit-content"
          field={{
            width: '310px',
            height: '200px'
          }} />
      </div>

      {isBinary.value ?
        <VTextField readonly v-model={link.value} clearable
                    label={t('editor.image.dialog.form.link')}
                    prepend-icon="mdi-link-variant" /> :
        <VTextField v-model={form.value.src} label={t('editor.image.dialog.form.link')} autofocus
                    prepend-icon="mdi-link-variant" />
      }

      <ImageProperties v-model={form.value} t={t} />
    </>
  }
})
