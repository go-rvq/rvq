<template>
  <div>
    <v-dialog :width="dialogWidth">
      <template v-slot:activator="{ isActive: isActive, props: activatorProps }">
        <v-text-field
          v-bind="activatorProps"
          :disabled="disabled"
          :loading="loading"
          :label="label"
          v-model="formattedDatetime"
          :hide-details="hideDetails"
          variant="underlined"
          @click="init()"
          readonly
        >
          <template v-slot:prepend>
            <v-icon
              icon="mdi-calendar-edit"
              :color="isActive ? 'primary' : ''"
              size="x-large"
            ></v-icon>
          </template>
          <template v-slot:loader>
            <v-progress-linear
              color="primary"
              indeterminate
              absolute
              height="2"
            ></v-progress-linear>
          </template>
        </v-text-field>
      </template>

      <template v-slot:default="{ isActive }">
        <v-card>
          <v-card-text class="px-0 py-0">
            <v-container>
              <v-row>
                <v-col>
                  <v-select
                    v-if="years.length"
                    :items="years"
                    :label="yearText"
                    v-model="year"
                  ></v-select>
                  <v-number-input
                    :min="min?.year || 0"
                    :max="max?.year || 2025"
                    type="number"
                    clearable
                    v-else
                    :label="yearText"
                    v-model="year"
                    controlVariant="split"
                  ></v-number-input>
                  <v-select
                    :v-if="year > 0"
                    :label="monthText"
                    :items="monthsAvailable"
                    :item-title="(item: number) => monthNames[item]"
                    :item-value="(item: number) => item"
                    v-model="month"
                  ></v-select>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="grey lighten-1" variant="text" @click.native="clearHandler(isActive)"
              >{{ clearText }}
            </v-btn>
            <v-btn color="green darken-1" variant="text" @click="okHandler(isActive)"
              >{{ okText }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </template>
    </v-dialog>
  </div>
</template>

<script lang="ts" setup>
import { format, parse } from 'date-fns'

import { computed, nextTick, onMounted, Ref, ref } from 'vue'

const props = defineProps({
    modelValue: {
      type: String,
      default: null
    },
    disabled: {
      type: Boolean
    },
    loading: {
      type: Boolean
    },
    label: {
      type: String,
      default: ''
    },
    dialogWidth: {
      type: Number,
      default: 380
    },
    clearText: {
      type: String,
      default: 'CLEAR'
    },
    okText: {
      type: String,
      default: 'OK'
    },
    yearText: {
      type: String,
      default: 'Year'
    },
    monthText: {
      type: String,
      default: 'Month'
    },
    years: {
      type: Array<number>,
      default: []
    },
    months: {
      type: Array<number>,
      default: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    },
    monthNames: {
      type: Array<String>,
      default: [
        '',
        'January',
        'February',
        'March',
        'April',
        'May',
        'June',
        'July',
        'August',
        'September',
        'October',
        'November',
        'December'
      ]
    },
    format: {
      type: String,
      default: 'yyyy-MM'
    },
    textFieldProps: {
      type: Object
    },
    datePickerProps: {
      type: Object
    },
    hideDetails: {
      type: Boolean
    },
    min: {
      type: String
    },
    max: {
      type: String
    }
  }),
  min = { year: 0, month: 0 },
  max = { year: 0, month: 0 },
  display = ref(false),
  monthsAvailable = ref<Array<number>>([]),
  monthsOfMin = () => {
    if (!min?.year) return props.months
    if (year.value < min.year) return []
    let items = props.months
    if (min.year == year.value) items = items.filter((m: number) => m >= min.month)
    return items
  },
  updateMonthsAvailable = (year: number) => {
    console.log('update months of ', year)
    if (!year) {
      monthsAvailable.value = []
      return
    }

    let items = monthsOfMin()
    if (items.length && max?.year && year >= max.year) {
      items = items.filter((m: number) => m <= max.month)
    }
    monthsAvailable.value = items
  },
  date = ref(),
  setDate = (val: any) => {
    date.value = val
  },
  year = computed<number>({
    // getter
    get() {
      return date.value ? (date.value as Date).getFullYear() : 0
    },
    // setter
    set(newValue) {
      console.log('year', newValue)
      const d: Date = date.value ? (date.value as Date) : new Date()
      let m = d.getUTCMonth()

      if (newValue && min?.year == newValue) {
        m = min.month
      }

      setDate(
        newValue
          ? new Date(
              Date.UTC(
                newValue,
                m,
                d.getUTCDate(),
                d.getUTCHours(),
                d.getUTCMinutes(),
                d.getUTCSeconds(),
                d.getUTCMilliseconds()
              )
            )
          : null
      )

      updateMonthsAvailable(newValue)
    }
  }),
  month = computed<number>({
    // getter
    get() {
      console.log(date.value)
      return date.value ? (date.value as Date).getMonth() : 0
    },
    // setter
    set(newValue) {
      if (!date.value) return
      const d = date.value ? date.value : new Date()
      setDate(
        new Date(
          Date.UTC(
            d.getUTCFullYear(),
            newValue,
            d.getUTCDate(),
            d.getUTCHours(),
            d.getUTCMinutes(),
            d.getUTCSeconds(),
            d.getUTCMilliseconds()
          )
        )
      )
    }
  }),
  dateTimeFormat = computed(() => {
    return props.format
  }),
  formattedDatetime = computed(() => {
    return date.value ? format(<Date>date.value, dateTimeFormat.value) : ''
  }),
  initValue = () => {
    if (!props.modelValue || props.modelValue.length < 7) {
      setDate(null)
      return
    }

    // see https://stackoverflow.com/a/9436948
    let d = parse(props.modelValue, dateTimeFormat.value, new Date())
    if (min?.year) {
      d = new Date(
        Date.UTC(
          min.year,
          min.month,
          d.getUTCDate(),
          d.getUTCHours(),
          d.getUTCMinutes(),
          d.getUTCSeconds(),
          d.getUTCMilliseconds()
        )
      )
    }
    setDate(d)
  },
  init = () => {
    if (props.min) {
      const d = parse(props.min, dateTimeFormat.value, new Date())
      min.year = d.getUTCFullYear()
      min.month = d.getUTCMonth()
    }

    if (props.max) {
      const d = parse(props.max, dateTimeFormat.value, new Date())
      max.year = d.getUTCFullYear()
      max.month = d.getUTCMonth()
    }

    initValue()
    updateMonthsAvailable(year.value)

    console.log('init')
    console.log(props.min, props.max)
    console.log(min, max)
  },
  emit = defineEmits(['update:modelValue']),
  okHandler = (isActive: Ref) => {
    isActive.value = false
    if (!date.value) {
      setDate(new Date())
    }
    emit('update:modelValue', formattedDatetime.value)
  },
  clearHandler = (isActive: Ref) => {
    isActive.value = false
    setDate(null)
    emit('update:modelValue', null)
  },
  resetPicker = () => {
    display.value = false
  }

onMounted(() => {
  nextTick(() => {
    init()
  })
})
</script>
