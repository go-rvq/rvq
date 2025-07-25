<script setup lang="ts">
import { FilterItem } from '@/lib/Filter/Model'
import { computed } from 'vue'
import * as constants from '@/lib/Filter/Constants'

const props = defineProps<{
  op: FilterItem
  compTranslations: any
}>()
const emit = defineEmits(['clear'])
const showValueComputed = computed(() => {
  let showValue = ''
  if (props.op.selected) {
    switch (props.op.itemType) {
      case 'DatetimeRangeItem':
      case 'DateRangeItem':
      case 'MonthRangeItem': {
        const mod = props.op.modifier || constants.ModifierBetween

        if (mod === constants.ModifierBetween) {
          if (props.op.valueFrom) {
            if (props.op.valueTo) {
              showValue = `${props.op.valueFrom} - ${props.op.valueTo}`
            } else {
              showValue = ` >= ${props.op.valueFrom}`
            }
          } else {
            if (props.op.valueTo) {
              showValue = ` < ${props.op.valueTo}`
            }
          }
        }
        break
      }
      case 'DateItem':
      case 'MonthItem': {
        showValue = props.op.valueIs
        break
      }
      case 'NumberItem': {
        const mod = props.op.modifier || 'equals'

        if (mod === 'equals') {
          const floatValue = parseFloat(props.op.valueIs)
          if (!isNaN(floatValue)) {
            showValue += floatValue
          }
        }

        if (mod === 'between') {
          const floatFrom = parseFloat(props.op.valueFrom || '')
          const floatTo = parseFloat(props.op.valueTo || '')
          const fromValid = !isNaN(floatFrom)
          const toValid = !isNaN(floatTo)
          if (fromValid) {
            if (toValid) {
              showValue = `${props.op.valueFrom} - ${props.op.valueTo}`
            } else {
              showValue = ` >= ${props.op.valueFrom}`
            }
          } else {
            if (toValid) {
              showValue = ` <= ${props.op.valueTo}`
            }
          }
        }

        if (mod === 'greaterThan') {
          const floatValue = parseFloat(props.op.valueIs)
          if (!isNaN(floatValue)) {
            showValue += ' > ' + props.op.valueFrom
          }
        }

        if (mod === 'lessThan') {
          const floatValue = parseFloat(props.op.valueIs)
          if (!isNaN(floatValue)) {
            showValue += ' < ' + props.op.valueTo
          }
        }
        break
      }
      case 'StringItem': {
        const mod = props.op.modifier || 'equals'
        if (mod === 'equals' && props.op.valueIs) {
          showValue = props.op.valueIs
        }

        if (mod === 'contains' && props.op.valueIs) {
          showValue = ' ~ ' + props.op.valueIs
        }
        break
      }
      case 'SelectItem': {
        const mod = props.op.modifier || 'equals'
        if (mod === 'equals' && props.op.valueIs) {
          showValue = props.op.options!.find((o) => o.value === props.op.valueIs)!.text
        }
        break
      }
      case 'MultipleSelectItem': {
        const mod = props.op.modifier || 'in'
        const textsAre = props.op
          .options!.filter((o) => props.op.valuesAre.includes(o.value))
          .map((o) => o.text)
        if (
          (mod === 'in' || mod === 'notIn') &&
          props.op.valuesAre &&
          props.op.valuesAre.length > 0
        ) {
          showValue = ` ${props.compTranslations[mod]} [ ${textsAre.join(', ')} ]`
        }
        break
      }
      case 'LinkageSelectItem': {
        const textsAre = props.op.valuesAre.map((o, i) => {
          const item = props.op.linkageSelectData?.items[i].find((x: any) => {
            return o === x.ID
          })
          return item.Name ? item.Name : item.ID
        })
        showValue = textsAre.join(',')
        break

        // const mod =  props.op.modifier || 'equals'
        // const textsAre =  props.op
        //   .options!.filter((o) =>  props.op.valuesAre.includes(o.value))
        //   .map((o) => o.text)
        // if (mod === 'equals' &&  props.op.valuesAre &&  props.op.valuesAre.length > 0) {
        //   showValue = textsAre.join(', ')
        // }
        // break
      }
      default:
        throw new Error(`itemType '${props.op.itemType}' not supported`)
    }
  }

  const showValueCopy = showValue
  showValue = ''
  let showLen = 0
  for (let i = 0; i < showValueCopy.length; i++) {
    showValue += showValueCopy.charAt(i)
    if (showValueCopy.charCodeAt(i) > 127) {
      showLen += 2
    } else {
      showLen++
    }
    if (showLen > 66) {
      showValue += '...'
      break
    }
  }
  return showValue
})

const clear = (e: any) => {
  emit('clear', e)
}
</script>

<template>
  <span class="cursor-pointer">
    <v-icon
      start
      @click="clear"
      :icon="op.selected ? 'mdi-close-circle' : 'mdi-plus-circle'"
    ></v-icon>
    {{ op.label }}
    <span v-if="op.selected">
      | <span class="text-primary">{{ showValueComputed }}</span>
    </span>
  </span>
</template>

<style scoped></style>
