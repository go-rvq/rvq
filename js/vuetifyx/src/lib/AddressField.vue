<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue'
import {
  loadMaps,
  newSession,
  onMapsAuthFailure,
  placeAddress,
  reverseAddress,
  suggest,
  type AddressValue,
  type Suggestion,
} from './googleMaps'

// <vx-address-field>: an address, typed and found by Google — the places it
// suggests as it is typed (the Places API (New)) —, and the place on a map
// below it: a marker, dragged to the right spot (the address read back from
// there). The value (v-model) is an AddressValue: schema.org's PostalAddress,
// its coordinates, the place's id. Typed and not chosen among the
// suggestions, its words alone. Without a key (apiKey), only typed, no map.

export type { AddressValue } from './googleMaps'

const props = withDefaults(
  defineProps<{
    modelValue?: AddressValue | null
    // the Google Maps key (Maps JavaScript + Places (New))
    apiKey?: string
    label?: string
    hint?: string
    placeholder?: string
    required?: boolean
    readonly?: boolean
    disabled?: boolean
    // the map's height; "0" draws none
    mapHeight?: string
    // said when there is no key: the address is only typed
    noKeyText?: string
    // said when the maps could not be loaded, or Google refused the key
    loadErrorText?: string
    zoom?: number
  }>(),
  {
    modelValue: null,
    apiKey: '',
    label: '',
    hint: '',
    placeholder: '',
    mapHeight: '220px',
    noKeyText: '',
    loadErrorText: 'The map could not be loaded.',
    zoom: 16,
  },
)

const emit = defineEmits<{ 'update:modelValue': [v: AddressValue | null] }>()

const mapEl = ref<HTMLElement>()
const search = ref(props.modelValue?.formatted ?? '')
const items = ref<Suggestion[]>([])
const loading = ref(false)
const error = ref('')
const ready = ref(false)
let g: any
let session: any
const map = shallowRef<any>()
let marker: any
let timer: ReturnType<typeof setTimeout> | undefined
let seq = 0
let stopAuth: (() => void) | undefined

const located = computed(() => props.modelValue?.lat != null && props.modelValue?.lng != null)
const showMap = computed(() => !!props.apiKey && props.mapHeight !== '0' && ready.value && located.value)
// what the combobox holds: the address, by its words
const selected = computed(() => props.modelValue?.formatted || null)

function set(v: AddressValue | null) {
  search.value = v?.formatted ?? ''
  emit('update:modelValue', v)
}

// typed: Google's suggestions, once the typing stops
watch(search, (q) => {
  clearTimeout(timer)
  if (!g || !q || q.trim().length < 3 || q === props.modelValue?.formatted) {
    items.value = []
    return
  }
  timer = setTimeout(async () => {
    const n = ++seq
    loading.value = true
    try {
      session ??= await newSession(g)
      const found = await suggest(g, q, session)
      if (n === seq) items.value = found
    } catch (e: any) {
      if (n === seq) {
        items.value = []
        error.value = props.loadErrorText
      }
    } finally {
      if (n === seq) loading.value = false
    }
  }, 250)
})

// chosen: a suggestion — the place's address —, or the words typed
async function choose(v: Suggestion | string | null) {
  if (v == null || v === '') return set(null)
  if (typeof v === 'string') {
    if (v !== props.modelValue?.formatted) set({ formatted: v })
    return
  }
  loading.value = true
  try {
    const a = await placeAddress(v, session)
    session = undefined // a session ends with the place chosen
    set(a)
  } catch {
    set({ formatted: v.text })
  } finally {
    loading.value = false
    items.value = []
  }
}

function placeMarker() {
  if (!map.value || !g) return
  const v = props.modelValue
  if (!v || v.lat == null || v.lng == null) {
    marker?.setMap(null)
    return
  }
  const pos = { lat: v.lat, lng: v.lng }
  if (!marker) {
    marker = new g.maps.Marker({ map: map.value, position: pos, draggable: !props.readonly && !props.disabled })
    // dragged: the address read back from where it was left
    marker.addListener('dragend', async () => {
      const p = marker.getPosition()
      const at = { lat: p.lat(), lng: p.lng() }
      const a = await reverseAddress(g, at)
      set(a ?? { ...(props.modelValue || { formatted: search.value }), lat: at.lat, lng: at.lng, placeId: '' })
    })
  } else {
    marker.setMap(map.value)
    marker.setPosition(pos)
  }
  map.value.panTo(pos)
}

function buildMap() {
  if (!g || !mapEl.value || map.value || !located.value) return
  const v = props.modelValue!
  map.value = new g.maps.Map(mapEl.value, {
    center: { lat: v.lat, lng: v.lng },
    zoom: props.zoom,
    streetViewControl: false,
    mapTypeControl: false,
  })
  placeMarker()
}

onMounted(async () => {
  if (!props.apiKey) return
  stopAuth = onMapsAuthFailure((code) => (error.value = `${props.loadErrorText} (${code})`))
  try {
    g = await loadMaps(props.apiKey)
  } catch {
    error.value = props.loadErrorText
    return
  }
  ready.value = true
  await nextTick()
  buildMap()
})
onBeforeUnmount(() => {
  clearTimeout(timer)
  stopAuth?.()
  marker?.setMap(null)
})

// the value given anew (another record, a save): what is shown follows —
// watched by what it says, not by the object
watch(
  () => [props.modelValue?.formatted ?? '', props.modelValue?.lat, props.modelValue?.lng].join('|'),
  async () => {
    search.value = props.modelValue?.formatted ?? ''
    await nextTick()
    if (!map.value) buildMap()
    else placeMarker()
  },
)

defineExpose({ set })
</script>

<template>
  <div class="vx-address-field" :data-address-field="label">
    <v-combobox
      :model-value="selected"
      v-model:search="search"
      :items="items"
      item-title="text"
      item-value="placeId"
      :label="label"
      :hint="!apiKey && noKeyText ? noKeyText : hint"
      :persistent-hint="!!hint || (!apiKey && !!noKeyText)"
      :placeholder="placeholder"
      :readonly="readonly"
      :disabled="disabled"
      :clearable="!readonly && !disabled && !required"
      :loading="loading"
      :error-messages="error ? [error] : []"
      :hide-no-data="true"
      no-filter
      return-object
      variant="underlined"
      prepend-inner-icon="mdi-map-marker-outline"
      autocomplete="off"
      data-address-input
      @update:model-value="choose"
    >
      <template #item="{ props: itemProps, item }">
        <v-list-item
          v-bind="itemProps"
          :title="item.raw.main"
          :subtitle="item.raw.secondary"
          prepend-icon="mdi-map-marker-outline"
          data-address-suggestion
        />
      </template>
    </v-combobox>
    <div
      v-show="showMap"
      ref="mapEl"
      class="vx-address-map"
      :style="{ height: mapHeight }"
      data-address-map
    ></div>
  </div>
</template>

<style>
.vx-address-map {
  width: 100%;
  margin-top: 4px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
</style>
