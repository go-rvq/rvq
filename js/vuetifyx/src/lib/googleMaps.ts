// The Google Maps JS API, loaded once per page — by the first key asked with
// —: <vx-seo-var-zipcodes>' map picker, <vx-address-field>'s,
// <vx-maps-key-test>'s. Resolves with window.google when it is ready.
//
// Places are asked of the Places API (New) — AutocompleteSuggestion,
// Place.fetchFields —: Google's legacy Autocomplete is not available to the
// projects created since March 2025.
let mapsLoader: Promise<any> | null = null
// the key it was loaded with: another can only be tried in a new page
let loadedKey = ''
// what Google said of the key: the code of its error ("RefererNotAllowedMapError")
let authError = ''
const authListeners: ((code: string) => void)[] = []

// Google calls gm_authFailure when it refuses the key; the reason it only
// writes to the console ("Google Maps JavaScript API error: XxxMapError",
// "Places API error: BillingNotEnabledMapError"): it is read from there
// while the API loads and works.
function watchAuth() {
  const w = window as any
  if (w.__vxMapsWatched) return
  w.__vxMapsWatched = true
  const orig = console.error
  console.error = function (...args: any[]) {
    const m = String(args[0] ?? '').match(/(?:Google Maps JavaScript|Places|Geocoding|Maps) API (?:error|warning): (\w+)/)
    if (m) {
      authError = m[1]
      authListeners.forEach((f) => f(authError))
    }
    return orig.apply(console, args as any)
  }
  const prev = w.gm_authFailure
  w.gm_authFailure = function () {
    if (!authError) authError = 'AuthFailure'
    authListeners.forEach((f) => f(authError))
    if (typeof prev === 'function') prev()
  }
}

// onMapsAuthFailure calls f with the code of Google's error when it refuses
// the key (now, if it did already); it returns how to stop.
export function onMapsAuthFailure(f: (code: string) => void): () => void {
  authListeners.push(f)
  if (authError) f(authError)
  return () => {
    const i = authListeners.indexOf(f)
    if (i >= 0) authListeners.splice(i, 1)
  }
}

// mapsKeyLoaded is the key the API was loaded with ("" when it was not).
export const mapsKeyLoaded = () => loadedKey

export function loadMaps(key: string): Promise<any> {
  const w = window as any
  if (w.google && w.google.maps && w.google.maps.importLibrary) return Promise.resolve(w.google)
  if (mapsLoader) return mapsLoader
  watchAuth()
  mapsLoader = new Promise((resolve, reject) => {
    if (!key) {
      reject(new Error('no maps key'))
      return
    }
    loadedKey = key
    const cb = '__vxMapsReady'
    w[cb] = () => resolve(w.google)
    const s = document.createElement('script')
    s.src = `https://maps.googleapis.com/maps/api/js?key=${encodeURIComponent(key)}&libraries=places&loading=async&callback=${cb}`
    s.async = true
    s.onerror = () => {
      // a failure is not kept: a later field may try again
      mapsLoader = null
      loadedKey = ''
      reject(new Error('maps load failed'))
    }
    document.head.appendChild(s)
  })
  return mapsLoader
}

// AddressValue is an address as schema.org's PostalAddress has it, with its
// coordinates and the place's id at Google (schemaform's AddressFields).
export interface AddressValue {
  formatted: string
  street?: string
  number?: string
  complement?: string
  neighborhood?: string
  city?: string
  region?: string
  postalCode?: string
  country?: string
  countryCode?: string
  lat?: number | null
  lng?: number | null
  placeId?: string
}

// Suggestion is a place Google suggests for what is typed.
export interface Suggestion {
  placeId: string
  text: string
  main: string
  secondary: string
  prediction: any
}

// newSession is a session of suggestions: the ones asked while one address
// is typed, and the place chosen, are one request to Google's billing.
export async function newSession(g: any): Promise<any> {
  const { AutocompleteSessionToken } = await g.maps.importLibrary('places')
  return new AutocompleteSessionToken()
}

// suggest are the places Google suggests for input.
export async function suggest(g: any, input: string, sessionToken?: any): Promise<Suggestion[]> {
  const { AutocompleteSuggestion } = await g.maps.importLibrary('places')
  const { suggestions } = await AutocompleteSuggestion.fetchAutocompleteSuggestions({ input, sessionToken })
  return (suggestions || [])
    .filter((s: any) => s.placePrediction)
    .map((s: any) => {
      const p = s.placePrediction
      return {
        placeId: p.placeId,
        text: p.text?.toString() ?? '',
        main: p.mainText?.toString() ?? p.text?.toString() ?? '',
        secondary: p.secondaryText?.toString() ?? '',
        prediction: p,
      }
    })
}

// addressOfComponents is the address of the parts Google gives — of the new
// Places ({longText, shortText, types}) or of a geocoding ({long_name,
// short_name, types}) —, its words and coordinates.
export function addressOfComponents(
  comps: any[],
  formatted: string,
  lat: number | null,
  lng: number | null,
  placeId: string,
): AddressValue {
  const comp = (type: string, short = false) => {
    const c = (comps || []).find((x: any) => (x.types || []).includes(type))
    if (!c) return ''
    return short ? c.shortText ?? c.short_name ?? '' : c.longText ?? c.long_name ?? ''
  }
  return {
    formatted,
    street: comp('route'),
    number: comp('street_number'),
    complement: comp('subpremise'),
    neighborhood: comp('sublocality_level_1') || comp('sublocality') || comp('neighborhood'),
    city: comp('locality') || comp('administrative_area_level_2') || comp('postal_town'),
    region: comp('administrative_area_level_1'),
    postalCode: comp('postal_code'),
    country: comp('country'),
    countryCode: comp('country', true),
    lat,
    lng,
    placeId,
  }
}

// placeAddress is the address of the place suggested: its words, its parts,
// its coordinates (Place.fetchFields).
export async function placeAddress(s: Suggestion, sessionToken?: any): Promise<AddressValue> {
  const place = s.prediction.toPlace()
  await place.fetchFields({ fields: ['id', 'formattedAddress', 'location', 'addressComponents', 'displayName'], sessionToken })
  return addressOfComponents(
    place.addressComponents,
    place.formattedAddress || place.displayName || s.text,
    place.location ? place.location.lat() : null,
    place.location ? place.location.lng() : null,
    place.id || s.placeId,
  )
}

// reverseAddress is the address at a point (a marker dragged there), else
// null.
export function reverseAddress(g: any, at: { lat: number; lng: number }): Promise<AddressValue | null> {
  return new Promise((resolve) => {
    new g.maps.Geocoder().geocode({ location: at }, (res: any[], status: string) => {
      if (status !== 'OK' || !res || !res[0]) return resolve(null)
      resolve(addressOfComponents(res[0].address_components, res[0].formatted_address, at.lat, at.lng, res[0].place_id || ''))
    })
  })
}
