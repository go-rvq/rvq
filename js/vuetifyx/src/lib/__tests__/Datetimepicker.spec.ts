import {expect, it} from 'vitest'

import Datetimepicker from '../Datetimepicker.vue'
import {mountTemplate} from '@/lib/__tests__/testutils'
import {nextTick} from 'vue'
import {flushPromises} from '@vue/test-utils'

it('Datetimepicker modelValue', async () => {
  const wrapper = mountTemplate(Datetimepicker, {
    modelValue: '2023-10-01 22:11'
  })
  await nextTick()
  await flushPromises()
  expect(wrapper.find('input').element.value).toContain('2023-10-01 22:11')
  expect(wrapper.find('input').element.value).not.toContain('2023-10-01 22:22')
})

it('Datetimepicker errorMessages', async () => {
  const wrapper = mountTemplate(Datetimepicker, {
    modelValue: '',
    errorMessages: ['Este campo não pode ser vazio']
  })
  await nextTick()
  await flushPromises()
  expect(wrapper.text()).toContain('Este campo não pode ser vazio')
})
