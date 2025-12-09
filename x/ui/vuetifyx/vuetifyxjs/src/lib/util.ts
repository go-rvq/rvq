import { propsFactory } from 'vuetify/lib/util/propsFactory'
import { PropType } from 'vue'

export const makeInputProps = propsFactory({
  id: String,
  hideDetails: [Boolean, String] as PropType<boolean | 'auto'>,
  hint: String,
  messages: {
    type: [Array, String] as PropType<string | readonly string[]>,
    default: () => ([]),
  },
  error: Boolean,
  errorMessages: {
    type: [Array, String] as PropType<string | readonly string[] | null>,
    default: () => ([]),
  },
}, 'vx-util-input')


export const makeInputWithPersistentHintProps = propsFactory({
  persistentHint: Boolean,
  ...makeInputProps()
}, 'vx-util-input')
