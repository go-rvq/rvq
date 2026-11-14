<template>
  <div class="vx-camera-controls" data-camera-controls>
    <div class="d-flex align-center mb-2">
      <span class="text-subtitle-2">{{ words.cameraControls }}</span>
      <v-spacer />
      <v-btn
        size="x-small"
        variant="text"
        prepend-icon="mdi-restore"
        data-camera-reset
        @click="reset"
        >{{ words.cameraReset }}</v-btn
      >
    </div>
    <!-- the image mirrored (the image editor's buttons): of every camera -->
    <div class="mb-3" data-camera-group="flip">
      <div class="text-caption text-medium-emphasis mb-1">{{ words.cameraFlip }}</div>
      <div class="d-flex ga-1">
        <v-btn
          icon="mdi-flip-horizontal"
          size="small"
          :variant="mirroredH ? 'tonal' : 'text'"
          :title="words.cameraFlipHorizontal"
          :aria-label="words.cameraFlipHorizontal"
          :color="mirroredH ? 'primary' : undefined"
          data-camera-flip-h
          @click="mirror('h')"
        />
        <v-btn
          icon="mdi-flip-vertical"
          size="small"
          :variant="mirroredV ? 'tonal' : 'text'"
          :title="words.cameraFlipVertical"
          :aria-label="words.cameraFlipVertical"
          :color="mirroredV ? 'primary' : undefined"
          data-camera-flip-v
          @click="mirror('v')"
        />
      </div>
    </div>
    <div v-if="!groups.length" class="text-caption text-medium-emphasis" data-camera-no-controls>
      {{ words.cameraNoControls }}
    </div>
    <div v-for="g in groups" :key="g.key" class="mb-3" :data-camera-group="g.key">
      <div class="text-caption text-medium-emphasis mb-1">{{ groupLabel(g.key) }}</div>
      <template v-for="c in g.controls" :key="c.key">
        <v-switch
          v-if="c.kind === 'mode'"
          :model-value="settings[c.key] !== 'manual'"
          :label="words.cameraAuto"
          density="compact"
          hide-details
          color="primary"
          :data-camera-control="c.key"
          @update:model-value="set(c, $event)"
        />
        <v-switch
          v-else-if="c.kind === 'bool'"
          :model-value="!!settings[c.key]"
          :label="label(c.key)"
          density="compact"
          hide-details
          color="primary"
          :data-camera-control="c.key"
          @update:model-value="set(c, $event)"
        />
        <div v-else class="vx-camera-range" :class="{auto: isAuto(c, settings)}">
          <div class="d-flex text-caption">
            <span>{{ label(c.key) }}</span>
            <v-spacer />
            <span>{{ shown(c) }}</span>
          </div>
          <v-slider
            :model-value="settings[c.key] ?? c.min"
            :min="c.min"
            :max="c.max"
            :step="c.step"
            density="compact"
            hide-details
            color="primary"
            :data-camera-control="c.key"
            @update:model-value="set(c, $event)"
          />
        </div>
      </template>
    </div>
    <div v-if="failed" class="text-caption text-error">{{ failed }}</div>
  </div>
</template>

<script lang="ts" setup>
// The camera's own controls (the panel beside the image of vx-webcam): what
// the track offers, set on it — the camera changes, not the photo.
import {computed, ref, shallowRef, watch} from 'vue'
import {useVxMessages} from '@/lib/locale'
import {type Control, type Group, constraintFor, controlsOf, initialOf, isAuto} from './controls'
import {flipped, type Orientation} from '@rvq/image-editor/orient'

const props = defineProps<{track: MediaStreamTrack | null}>()

// the orientation of the image (the image editor's): mirrored as it is seen
const orientation = defineModel<Orientation>('orientation', {default: () => ({rotate: 0, flipH: false})})
// mirrored horizontally: once, not turned; vertically: turned 180°, once —
// both, turned 180° alone
const mirroredH = computed(() => (orientation.value.flipH && orientation.value.rotate === 0) || (!orientation.value.flipH && orientation.value.rotate === 180))
const mirroredV = computed(() => orientation.value.rotate === 180)
function mirror(axis: 'h' | 'v') {
  orientation.value = flipped(orientation.value.rotate, orientation.value.flipH, axis)
}

const words = useVxMessages()
const groups = shallowRef<Group[]>([])
const caps = shallowRef<Record<string, any>>({})
const settings = ref<Record<string, any>>({})
const failed = ref('')
let initial: MediaTrackConstraintSet = {}

function read() {
  const t = props.track
  if (!t || typeof (t as any).getCapabilities !== 'function') {
    groups.value = []
    return
  }
  caps.value = (t as any).getCapabilities() || {}
  settings.value = {...(t.getSettings() as any)}
  groups.value = controlsOf(caps.value)
}

watch(
  () => props.track,
  () => {
    read()
    initial = initialOf(groups.value, settings.value)
  },
  {immediate: true},
)

// a slider moves many times: one apply at a time, the last value kept
let busy = false
let next: MediaTrackConstraintSet | null = null
async function apply(set: MediaTrackConstraintSet) {
  const t = props.track
  if (!t) return
  if (busy) {
    next = {...(next || {}), ...set}
    return
  }
  busy = true
  try {
    await t.applyConstraints({advanced: [set]} as MediaTrackConstraints)
    failed.value = ''
  } catch (e: any) {
    failed.value = String(e?.message || e)
  } finally {
    busy = false
    settings.value = {...(t.getSettings() as any)}
    if (next) {
      const n = next
      next = null
      apply(n)
    }
  }
}

function set(c: Control, v: any) {
  const s = constraintFor(c, v, settings.value, caps.value)
  // shown at once, the camera's answer after
  settings.value = {...settings.value, ...s}
  apply(s)
}

function reset() {
  orientation.value = {rotate: 0, flipH: false}
  if (groups.value.length) apply(initial)
}

const LABELS: Record<string, keyof typeof words.value> = {
  exposureCompensation: 'cameraExposureCompensation',
  exposureTime: 'cameraExposureTime',
  iso: 'cameraIso',
  colorTemperature: 'cameraColorTemperature',
  focusDistance: 'cameraFocusDistance',
  brightness: 'cameraBrightness',
  contrast: 'cameraContrast',
  saturation: 'cameraSaturation',
  sharpness: 'cameraSharpness',
  zoom: 'cameraZoom',
  pan: 'cameraPan',
  tilt: 'cameraTilt',
  torch: 'cameraTorch',
}
const GROUP_LABELS: Record<Group['key'], keyof typeof words.value> = {
  exposure: 'cameraExposure',
  whiteBalance: 'cameraWhiteBalance',
  focus: 'cameraFocus',
  image: 'cameraImage',
  position: 'cameraPosition',
  light: 'cameraLight',
}
const label = (k: string) => (LABELS[k] ? (words.value[LABELS[k]] as string) : k)
const groupLabel = (k: Group['key']) => words.value[GROUP_LABELS[k]] as string

function shown(c: Control): string {
  const v = settings.value[c.key]
  if (typeof v !== 'number') return ''
  if (c.key === 'colorTemperature') return `${Math.round(v)} K`
  if (c.key === 'zoom') return `${v.toFixed(1)}×`
  return c.kind === 'range' && c.step < 1 ? v.toFixed(2) : String(Math.round(v))
}

defineExpose({groups: computed(() => groups.value), read})
</script>

<style scoped>
.vx-camera-controls {
  padding: 12px;
}
.vx-camera-range.auto {
  opacity: 0.6;
}
</style>
