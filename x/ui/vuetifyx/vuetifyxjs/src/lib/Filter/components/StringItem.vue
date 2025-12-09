<script setup lang="ts">
import { onMounted, PropType, ref } from 'vue'

interface ItemValue {
  v: string
}

export interface Model {
  modifier?: string
  valueIs?: string
}

export interface Translations {
  equals: string
  contains: string
}

export interface Config {
  rows?: number
  hint?: string
  width?: any
  maxWidth?: any
  many?: boolean
  modifierDisabled?:boolean
}

const props = defineProps(
  {
    modelValue: {
      type: Object as PropType<Model>, default: () => {
      }
    },
    translations: {
      type: Object as PropType<Translations>, default: () => {
      }
    },
    config: {
      type: Object as PropType<Config>, default: () => {
      }
    }
  }
)

props.modelValue.modifier = props.modelValue.modifier || 'contains'

const emit = defineEmits(['update:modelValue', 'postSetup'])
const t = props.translations
const items = ref([
  { text: t.equals, value: 'equals' },
  { text: t.contains, value: 'contains' }
])
const values = ref<ItemValue[]>([])

const config = props.config || {}
config.rows ||= 1

const addValue = () => {
  values.value.push({ v: '' })
}

const delValue = (index: number) => {
  values.value.splice(index, 1)
}

const done = (cb:(empty:boolean)=>void) => {
  if (config.many) {
    const value = values.value.filter(e => e.v.length).map(e => e.v).join(',')
    props.modelValue.valueIs = value
    emit('update:modelValue', props.modelValue)
    cb(value.length === 0)
  } else {
    cb(!props.modelValue.valueIs || props.modelValue.valueIs?.length === 0)
  }
}

defineExpose(done)

onMounted(()=>{
  emit('postSetup', {beforeDoneHandler: done})
})

if (config.many) {
  if (props.modelValue.valueIs) {
    props.modelValue.valueIs.split(',').filter(v => v.length).forEach((v) => {
      values.value.push({v:v})
    })
  } else {
    addValue()
  }
}

</script>

<template>
  <div>
    <div v-if="!config.modifierDisabled">
      <v-select
        class="d-inline-block"
        style="width: 200px"
        v-model="props.modelValue.modifier"
        :items="items"
        item-title="text"
        item-value="value"
        variant="underlined"
        hide-details
      ></v-select>
    </div>
    <div class="d-flex">
      <v-icon v-if="!config.modifierDisabled"
        class="py-8 mr-4 float-md-start"
        icon="mdi-subdirectory-arrow-right"
        size="large"
      ></v-icon>
      <template v-if="config.many">
        <div>
          <v-list density="compact">
            <v-list-item density="compact" v-for="(item, i) in values">
              <v-text-field density="compact" :width="config.width" hide-details v-model="item.v"></v-text-field>
              <template #prepend>
                <v-avatar density="compact" color="secondary">{{ i + 1 }}</v-avatar>
              </template>
              <template #append>
                <v-btn
                  density="compact"
                  variant="text"
                  icon="mdi-close"
                  @click="delValue(i)"
                />
              </template>
            </v-list-item>
          </v-list>
          <div class="d-flex">
            <v-spacer />
            <v-btn
              density="compact"
              variant="text"
              icon="mdi-plus"
              @click="addValue"
            />
            <v-spacer />
          </div>
        </div>
      </template>
      <v-textarea v-else-if="config.rows && config.rows > 1"
                  :rows="config.rows as number"
                  class="d-inline-block"
                  :width="config.width"
                  :max-width="config.maxWidth"
                  variant="underlined"
                  type="text"
                  v-model="props.modelValue.valueIs"
                  persistent-hint
                  :hint="config.hint"
                  :hide-details="!config.hint"
      />
      <v-text-field v-else
                    class="d-inline-block"
                    style="width: 120px"
                    variant="underlined"
                    type="text"
                    v-model="props.modelValue.valueIs"
                    persistent-hint
                    :hint="config.hint"
                    :hide-details="!config.hint"
      />
    </div>
  </div>
</template>

<style scoped></style>
