<script lang="ts">
import {defineComponent, PropType} from 'vue'

// A selected place: enough to restore its marker and its postal code.
interface Place {
  placeId?: string
  name?: string
  lat: number
  lng: number
  zip?: string
}

// Loads the Google Maps JS API (with the Places library) once per page, keyed by
// api key. Resolves when window.google.maps is ready.
let mapsLoader: Promise<any> | null = null
function loadMaps(key: string): Promise<any> {
  const w = window as any
  if (w.google && w.google.maps) return Promise.resolve(w.google)
  if (mapsLoader) return mapsLoader
  mapsLoader = new Promise((resolve, reject) => {
    if (!key) {
      reject(new Error('no maps key'))
      return
    }
    const s = document.createElement('script')
    s.src = `https://maps.googleapis.com/maps/api/js?key=${encodeURIComponent(key)}&libraries=places`
    s.async = true
    s.onload = () => resolve((window as any).google)
    s.onerror = () => reject(new Error('maps load failed'))
    document.head.appendChild(s)
  })
  return mapsLoader
}

// SeoVarZipcodes edits a "zipcodes" SEO variable: a button opens a Google map to
// select places; their postal codes become the value (comma-separated). The
// selected places are kept in `data` so re-editing shows them back on the map.
export default defineComponent({
  name: 'VXSeoVarZipcodes',

  props: {
    value: {type: String, default: ''},
    data: {type: [Array, Object, String] as PropType<any>, default: () => []},
    mapsKey: {type: String, default: ''},
    action: {type: String, default: ''},
  },

  emits: {
    'update:value': (_v: string) => true,
    'update:data': (_d: any) => true,
  },

  data() {
    return {
      dialog: false,
      loading: false,
      error: '',
      places: [] as Place[],
      map: null as any,
      markers: [] as any[],
      search: '',
    }
  },

  created() {
    this.places = this.parsePlaces(this.data)
  },

  watch: {
    data(v) {
      this.places = this.parsePlaces(v)
    },
  },

  methods: {
    parsePlaces(v: any): Place[] {
      if (Array.isArray(v)) return v as Place[]
      if (typeof v === 'string' && v.trim()) {
        try {
          const a = JSON.parse(v)
          return Array.isArray(a) ? a : []
        } catch {
          return []
        }
      }
      return []
    },
    zipList(): string {
      const seen = new Set<string>()
      const zips: string[] = []
      for (const p of this.places) {
        if (p.zip && !seen.has(p.zip)) {
          seen.add(p.zip)
          zips.push(p.zip)
        }
      }
      return zips.join(', ')
    },
    commit() {
      this.$emit('update:data', this.places)
      this.$emit('update:value', this.zipList())
    },
    async open() {
      this.dialog = true
      this.error = ''
      this.loading = true
      try {
        await loadMaps(this.mapsKey)
        this.$nextTick(() => this.initMap())
      } catch (e: any) {
        this.error = this.mapsKey ? 'Não foi possível carregar o Google Maps.' : 'Configure a Google Maps API Key em SEOConfig.'
      } finally {
        this.loading = false
      }
    },
    initMap() {
      const g = (window as any).google
      const el = this.$refs.map as HTMLElement
      if (!g || !el) return
      const center = this.places.length ? {lat: this.places[0].lat, lng: this.places[0].lng} : {lat: 42.4792, lng: -71.1523}
      this.map = new g.maps.Map(el, {center, zoom: this.places.length ? 11 : 9})
      this.markers = []
      for (const p of this.places) this.addMarker(p)

      // Places Autocomplete search box.
      const input = this.$refs.search as HTMLInputElement
      if (input && g.maps.places) {
        const ac = new g.maps.places.Autocomplete(input, {fields: ['place_id', 'geometry', 'name', 'address_components']})
        ac.addListener('place_changed', () => {
          const place = ac.getPlace()
          if (!place || !place.geometry) return
          this.addPlace(place)
          this.map.panTo(place.geometry.location)
          this.map.setZoom(12)
          this.search = ''
        })
      }

      // Click the map to add a place (reverse-geocoded for its postal code).
      this.map.addListener('click', (ev: any) => {
        const geocoder = new g.maps.Geocoder()
        geocoder.geocode({location: ev.latLng}, (results: any[], status: string) => {
          if (status === 'OK' && results[0]) this.addPlace(results[0])
        })
      })
    },
    zipFromComponents(components: any[]): string {
      if (!components) return ''
      for (const c of components) {
        if ((c.types || []).includes('postal_code')) return c.long_name || c.short_name || ''
      }
      return ''
    },
    addPlace(place: any) {
      const loc = place.geometry.location
      const p: Place = {
        placeId: place.place_id,
        name: place.name || place.formatted_address,
        lat: typeof loc.lat === 'function' ? loc.lat() : loc.lat,
        lng: typeof loc.lng === 'function' ? loc.lng() : loc.lng,
        zip: this.zipFromComponents(place.address_components),
      }
      // Deduplicate by placeId or coordinates.
      if (this.places.some((x) => (p.placeId && x.placeId === p.placeId) || (x.lat === p.lat && x.lng === p.lng))) return
      this.places.push(p)
      this.addMarker(p)
    },
    addMarker(p: Place) {
      const g = (window as any).google
      if (!g || !this.map) return
      const marker = new g.maps.Marker({position: {lat: p.lat, lng: p.lng}, map: this.map, title: p.name})
      marker.addListener('click', () => this.removePlace(p, marker))
      this.markers.push(marker)
    },
    removePlace(p: Place, marker: any) {
      this.places = this.places.filter((x) => x !== p)
      marker.setMap(null)
      this.markers = this.markers.filter((m) => m !== marker)
    },
    confirm() {
      this.commit()
      this.dialog = false
    },
  },
})
</script>

<template>
  <div class="vx-seo-var-zipcodes">
    <div class="d-flex align-center ga-2">
      <v-text-field :model-value="value" density="compact" variant="underlined" hide-details readonly
                    placeholder="Selecione lugares no mapa" />
      <v-btn prepend-icon="mdi-map-marker" variant="outlined" size="small" @click="open">Mapa</v-btn>
    </div>

    <v-dialog v-model="dialog" max-width="900">
      <v-card>
        <v-toolbar density="compact" color="primary">
          <v-toolbar-title>Selecionar CEPs no mapa</v-toolbar-title>
          <v-spacer />
          <v-btn icon="mdi-close" @click="dialog = false" />
        </v-toolbar>
        <v-card-text>
          <v-alert v-if="error" type="warning" density="compact" variant="tonal" class="mb-2">{{ error }}</v-alert>
          <v-progress-linear v-if="loading" indeterminate class="mb-2" />
          <input ref="search" v-model="search" placeholder="Buscar lugar…" class="vx-zip-search" />
          <div ref="map" class="vx-zip-map"></div>
          <div class="mt-3">
            <div class="text-caption text-medium-emphasis mb-1">Lugares selecionados (clique no marcador para remover):</div>
            <v-chip v-for="(p, i) in places" :key="i" class="ma-1" closable
                    @click:close="places.splice(i, 1)">
              {{ p.name }}<template v-if="p.zip"> · {{ p.zip }}</template>
            </v-chip>
            <div v-if="!places.length" class="text-medium-emphasis">Nenhum lugar.</div>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">Cancelar</v-btn>
          <v-btn color="primary" variant="flat" @click="confirm">Aplicar</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.vx-zip-search {
  width: 100%;
  padding: 8px;
  margin-bottom: 8px;
  border: 1px solid rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}
.vx-zip-map {
  width: 100%;
  height: 420px;
  border-radius: 4px;
  background: #eee;
}
</style>
