<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref } from 'vue'
import AddressField, { type AddressValue } from './AddressField.vue'
import { loadMaps, mapsKeyLoaded, onMapsAuthFailure, suggest } from './googleMaps'

// <vx-maps-key-test>: tests a Google Maps key (the SEO settings'), the one
// given — the field's, even before it is saved —: a button that opens a dialog
// where, one after the other, the Maps JavaScript API is loaded with it (and,
// when Google refuses it, why: its error's code), the Places API asked for a
// place, the Geocoding API for an address; and an address field with its map,
// to try by hand.

interface Texts {
  button: string
  title: string
  noKey: string
  maps: string
  places: string
  geocoding: string
  waiting: string
  ok: string
  failed: string
  tryIt: string
  reload: string
  close: string
  // what a code of Google's means; {code} where it goes
  rejected: string
  hints: Record<string, string>
}

const props = withDefaults(
  defineProps<{
    apiKey?: string
    // what is searched: a place known to exist
    query?: string
    texts?: Partial<Texts>
  }>(),
  { apiKey: '', query: 'Boston, MA', texts: () => ({}) },
)

const defaults: Texts = {
  button: 'Test',
  title: 'Test of the Google Maps key',
  noKey: 'Type the key first.',
  maps: 'Maps JavaScript API',
  places: 'Places API (New) (search of places)',
  geocoding: 'Geocoding API (addresses)',
  waiting: 'Testing…',
  ok: 'Working',
  failed: 'Failed',
  tryIt: 'Search an address: it shows on the map.',
  reload: 'The page already loaded Google Maps with another key: save the settings and reload the page to test this one.',
  close: 'Close',
  rejected: 'Google refused the key: {code}.',
  hints: {
    RefererNotAllowedMapError: 'The key does not allow this site: add its address to the key’s website restrictions.',
    ApiNotActivatedMapError: 'The Maps JavaScript API is not enabled in the key’s project.',
    InvalidKeyMapError: 'The key is not valid.',
    ExpiredKeyMapError: 'The key expired.',
    BillingNotEnabledMapError: 'The key’s project has no billing account: enable billing in the Google Cloud console.',
    ZERO_RESULTS: 'Google found nothing for the search.',
    REQUEST_DENIED: 'The API is not enabled in the key’s project, or the key does not allow it.',
    OVER_QUERY_LIMIT: 'The key’s quota is over.',
  },
}
const t = computed<Texts>(() => ({ ...defaults, ...props.texts, hints: { ...defaults.hints, ...(props.texts.hints || {}) } }))

type State = 'idle' | 'waiting' | 'ok' | 'failed'
const dialog = ref(false)
const checks = reactive<Record<'maps' | 'places' | 'geocoding', { state: State; detail: string }>>({
  maps: { state: 'idle', detail: '' },
  places: { state: 'idle', detail: '' },
  geocoding: { state: 'idle', detail: '' },
})
const address = ref<AddressValue | null>(null)
const keyChanged = ref(false)
const loaded = ref(false)

const why = (code: string) => t.value.rejected.replace('{code}', code) + (t.value.hints[code] ? ' ' + t.value.hints[code] : '')

let stop: (() => void) | undefined
onBeforeUnmount(() => stop?.())

async function run() {
  dialog.value = true
  for (const c of Object.values(checks)) Object.assign(c, { state: 'idle', detail: '' })
  keyChanged.value = false
  loaded.value = false
  const key = (props.apiKey || '').trim()
  if (!key) {
    Object.assign(checks.maps, { state: 'failed', detail: t.value.noKey })
    return
  }
  if (mapsKeyLoaded() && mapsKeyLoaded() !== key) {
    keyChanged.value = true
    return
  }
  checks.maps.state = 'waiting'
  stop?.()
  stop = onMapsAuthFailure((code) => {
    Object.assign(checks.maps, { state: 'failed', detail: why(code) })
  })
  let g: any
  try {
    g = await loadMaps(key)
  } catch (e: any) {
    Object.assign(checks.maps, { state: 'failed', detail: String(e?.message || e) })
    return
  }
  // Google refuses a key a moment after the script loads
  await new Promise((r) => setTimeout(r, 1200))
  if ((checks.maps.state as State) === 'failed') return
  Object.assign(checks.maps, { state: 'ok', detail: '' })
  loaded.value = true

  checks.places.state = 'waiting'
  // Places API (New): a suggestion for the query; refused, the code Google
  // wrote (BillingNotEnabledMapError…), else what it threw
  suggest(g, props.query)
    .then((found) => {
      if (found.length) Object.assign(checks.places, { state: 'ok', detail: found[0].text })
      else Object.assign(checks.places, { state: 'failed', detail: why('ZERO_RESULTS') })
    })
    .catch((e: any) => {
      setTimeout(() => {
        const code = String(e?.message || e).match(/\b([A-Z][A-Za-z]+(?:MapError|_DENIED|_LIMIT))\b/)?.[1]
        Object.assign(checks.places, { state: 'failed', detail: code ? why(code) : String(e?.message || e) })
      }, 300)
    })
  checks.geocoding.state = 'waiting'
  new g.maps.Geocoder().geocode({ address: props.query }, (res: any[], status: string) => {
    if (status === 'OK' && res?.length) {
      const loc = res[0].geometry.location
      Object.assign(checks.geocoding, { state: 'ok', detail: `${res[0].formatted_address} (${loc.lat().toFixed(4)}, ${loc.lng().toFixed(4)})` })
    } else Object.assign(checks.geocoding, { state: 'failed', detail: why(status) })
  })
}

const icon = (s: State) =>
  ({ idle: 'mdi-circle-outline', waiting: 'mdi-progress-clock', ok: 'mdi-check-circle', failed: 'mdi-alert-circle' })[s]
const color = (s: State) => ({ idle: 'grey', waiting: 'info', ok: 'success', failed: 'error' })[s]
</script>

<template>
  <v-btn variant="tonal" size="small" prepend-icon="mdi-map-check-outline" data-maps-key-test @click="run">
    {{ t.button }}
  </v-btn>
  <v-dialog v-model="dialog" max-width="720">
    <v-card :title="t.title" data-maps-key-test-dialog>
      <v-card-text>
        <v-alert v-if="keyChanged" type="warning" variant="tonal" class="mb-4">{{ t.reload }}</v-alert>
        <v-list density="compact" class="mb-2">
          <v-list-item
            v-for="name in (['maps', 'places', 'geocoding'] as const)"
            lines="three"
            :key="name"
            :data-maps-check="name"
            :data-maps-state="checks[name].state"
            :title="t[name]"
            :subtitle="checks[name].state === 'waiting' ? t.waiting : checks[name].detail"
          >
            <template #prepend>
              <v-icon :icon="icon(checks[name].state)" :color="color(checks[name].state)" />
            </template>

<style>
/* the reason of a failure whole, not cut */
[data-maps-key-test-dialog] .v-list-item-subtitle {
  -webkit-line-clamp: unset;
  line-clamp: unset;
  white-space: normal;
  overflow-wrap: anywhere;
}
</style>
            <template #append>
              <v-chip
                v-if="checks[name].state === 'ok' || checks[name].state === 'failed'"
                size="x-small"
                variant="tonal"
                :color="color(checks[name].state)"
                >{{ checks[name].state === 'ok' ? t.ok : t.failed }}</v-chip
              >
            </template>

<style>
/* the reason of a failure whole, not cut */
[data-maps-key-test-dialog] .v-list-item-subtitle {
  -webkit-line-clamp: unset;
  line-clamp: unset;
  white-space: normal;
  overflow-wrap: anywhere;
}
</style>
          </v-list-item>
        </v-list>
        <template v-if="loaded">
          <p class="text-body-2 text-medium-emphasis mb-1">{{ t.tryIt }}</p>
          <AddressField v-model="address" :api-key="apiKey" map-height="260px" />
        </template>

<style>
/* the reason of a failure whole, not cut */
[data-maps-key-test-dialog] .v-list-item-subtitle {
  -webkit-line-clamp: unset;
  line-clamp: unset;
  white-space: normal;
  overflow-wrap: anywhere;
}
</style>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="dialog = false">{{ t.close }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style>
/* the reason of a failure whole, not cut */
[data-maps-key-test-dialog] .v-list-item-subtitle {
  -webkit-line-clamp: unset;
  line-clamp: unset;
  white-space: normal;
  overflow-wrap: anywhere;
}
</style>
