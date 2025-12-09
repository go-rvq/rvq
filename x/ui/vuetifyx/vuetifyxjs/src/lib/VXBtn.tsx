import { genericComponent, propsFactory, useRender } from 'vuetify/lib/util/index'
import { VBtn, VTooltip } from 'vuetify/components'
import { makeVBtnProps, type VBtnSlots } from 'vuetify/lib/components/VBtn/VBtn'

export const makeVXBtnProps = propsFactory({
  title: String,
  ...makeVBtnProps()
}, 'VXBtn')

export const VXBtn = genericComponent<VBtnSlots>()({
  name: 'VXBtn',

  props: makeVXBtnProps(),

  setup(props, { slots, attrs }) {
    useRender(() => {
      const mustBtnProps = { ...props, ...attrs }

      return (
        <VTooltip text={props.title} location="top">
          {{
            activator: ({ props: btnProps }) => (
              <VBtn {...{ ...mustBtnProps, ...btnProps }} v-slots={slots} />
            )
          }}
        </VTooltip>
      )
    })
    return {}
    }
})

export type VXBtn = InstanceType<typeof VXBtn>
