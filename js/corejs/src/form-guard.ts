// <go-plaid-form-guard :form="form">: a form in an overlay that holds a close
// the user asks for (ESC, a click outside, the close button: the closer's
// requestClose) while its CONTENT changed — what `form` holds against what it
// held when it opened. Renderless: the confirmation is the overlay's
// (closer.confirmClose).
//
// What a field held at first is the first value it had: the fields assign
// their values as they mount (v-assign), some later than others (a tab drawn
// when it is opened), so a key appearing is its initial value, not a change.
// The baseline is the closer's, for the opening: the form drawn again (a save
// refused, its errors) keeps it — the changes are still not saved.

import { defineComponent, inject, onBeforeUnmount, onMounted, watch, type PropType } from 'vue'

/** keyOf is a value as compared: JSON, a file by its name, size and date. */
export function keyOf(v: any): string {
  return JSON.stringify(v, (_k, x) => {
    if (typeof File !== 'undefined' && x instanceof File) {
      return { file: x.name, size: x.size, modified: x.lastModified }
    }
    return x
  })
}

/** fieldsOf are the fields of a form: its own keys, not the scope's ($…, __…). */
export function fieldsOf(form: any): string[] {
  if (!form || typeof form !== 'object') {
    return []
  }
  return Object.keys(form).filter((k) => !k.startsWith('$') && !k.startsWith('__'))
}

/** FormBaseline is what each field of a form held first. */
export class FormBaseline {
  values = new Map<string, string>()

  /** see records the fields not seen yet, their value now their first. */
  see(form: any) {
    for (const k of fieldsOf(form)) {
      if (!this.values.has(k)) {
        this.values.set(k, keyOf(form[k]))
      }
    }
  }

  /** changed says a field holds other than it held first. */
  changed(form: any): boolean {
    for (const k of fieldsOf(form)) {
      const first = this.values.get(k)
      if (first !== undefined && first !== keyOf(form[k])) {
        return true
      }
    }
    return false
  }
}

export default defineComponent({
  name: 'GoPlaidFormGuard',
  props: {
    form: { type: Object as PropType<any>, default: undefined }
  },
  setup(props) {
    const closer = inject<any>('closer', null)
    let stopWatch: (() => void) | undefined

    onMounted(() => {
      if (!closer || typeof closer.onCloseGuard !== 'function' || !props.form) {
        return
      }
      // the form of the opening: the latest drawn
      closer.__guardForm = props.form
      if (!closer.__guardBaseline) {
        closer.__guardBaseline = new FormBaseline()
        // a new opening: a new baseline
        closer.onOpen?.((c: any) => {
          c.__guardBaseline = new FormBaseline()
        })
        closer.onCloseGuard((c: any) => !!c.__guardBaseline?.changed(c.__guardForm))
      }
      // the fields of now, once the ones mounting with this have assigned
      // theirs; the ones that come later, as they come
      const baseline = () => closer.__guardBaseline as FormBaseline
      setTimeout(() => baseline().see(props.form), 0)
      stopWatch = watch(
        () => fieldsOf(props.form).join('\u0001'),
        () => baseline().see(props.form)
      )
    })

    onBeforeUnmount(() => stopWatch?.())

    return () => null
  }
})
