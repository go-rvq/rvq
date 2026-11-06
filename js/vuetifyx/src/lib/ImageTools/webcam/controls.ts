// The controls of the camera itself (its track: MediaStreamTrack): what it
// offers (getCapabilities) — exposure, white balance, focus, the image
// (brightness, contrast…), zoom, the torch — and the constraint that sets one
// (applyConstraints): a value that needs the manual mode of its group
// (exposureTime: exposureMode manual) takes the group to manual with it.

export type ModeKey = 'exposureMode' | 'whiteBalanceMode' | 'focusMode'

export interface RangeControl {
  kind: 'range'
  key: string
  min: number
  max: number
  step: number
  // the mode it needs manual (none: it works in any)
  mode?: ModeKey
}

export interface ModeControl {
  kind: 'mode'
  key: ModeKey
  // auto: the automatic mode it has (continuous, else single-shot)
  auto: string
  manual: boolean
}

export interface BoolControl {
  kind: 'bool'
  key: string
}

export type Control = RangeControl | ModeControl | BoolControl

export interface Group {
  key: 'exposure' | 'whiteBalance' | 'focus' | 'image' | 'position' | 'light'
  controls: Control[]
}

// The groups, in order, and the controls of each.
const GROUPS: {key: Group['key']; mode?: ModeKey; ranges: {key: string; mode?: boolean}[]; bool?: string}[] = [
  {key: 'exposure', mode: 'exposureMode', ranges: [{key: 'exposureCompensation'}, {key: 'exposureTime', mode: true}, {key: 'iso', mode: true}]},
  {key: 'whiteBalance', mode: 'whiteBalanceMode', ranges: [{key: 'colorTemperature', mode: true}]},
  {key: 'focus', mode: 'focusMode', ranges: [{key: 'focusDistance', mode: true}]},
  {key: 'image', ranges: [{key: 'brightness'}, {key: 'contrast'}, {key: 'saturation'}, {key: 'sharpness'}]},
  {key: 'position', ranges: [{key: 'zoom'}, {key: 'pan'}, {key: 'tilt'}]},
  {key: 'light', ranges: [], bool: 'torch'},
]

function isRange(c: any): c is {min: number; max: number; step?: number} {
  return c && typeof c === 'object' && typeof c.min === 'number' && typeof c.max === 'number' && c.max > c.min
}

// controlsOf: the groups of controls a camera offers (its capabilities);
// none offered, none.
export function controlsOf(caps: Record<string, any> | null | undefined): Group[] {
  if (!caps) return []
  const out: Group[] = []
  for (const g of GROUPS) {
    const controls: Control[] = []
    let modeKey: ModeKey | undefined
    if (g.mode && Array.isArray(caps[g.mode])) {
      const modes: string[] = caps[g.mode]
      const auto = modes.includes('continuous') ? 'continuous' : modes.includes('single-shot') ? 'single-shot' : ''
      if (auto && modes.includes('manual')) {
        controls.push({kind: 'mode', key: g.mode, auto, manual: true})
        modeKey = g.mode
      }
    }
    for (const r of g.ranges) {
      const c = caps[r.key]
      if (!isRange(c)) continue
      controls.push({
        kind: 'range',
        key: r.key,
        min: c.min,
        max: c.max,
        step: c.step && c.step > 0 ? c.step : (c.max - c.min) / 100,
        mode: r.mode ? modeKey ?? (g.mode as ModeKey) : undefined,
      })
    }
    if (g.bool && Array.isArray(caps[g.bool]) && caps[g.bool].includes(true)) controls.push({kind: 'bool', key: g.bool})
    if (controls.length) out.push({key: g.key, controls})
  }
  return out
}

// constraintFor: what sets the control c to v — with the manual mode of its
// group when it needs it and the camera is not in it (settings).
export function constraintFor(c: Control, v: any, settings: Record<string, any>, caps: Record<string, any>): MediaTrackConstraintSet {
  const set: Record<string, any> = {}
  if (c.kind === 'mode') {
    set[c.key] = v ? c.auto : 'manual'
    return set
  }
  set[c.key] = c.kind === 'range' ? Number(v) : !!v
  if (c.kind === 'range' && c.mode && settings[c.mode] !== 'manual' && Array.isArray(caps[c.mode]) && caps[c.mode].includes('manual')) {
    set[c.mode] = 'manual'
  }
  return set
}

// initialOf: the values of the controls now (the camera's settings) — what
// "Reset" puts back.
export function initialOf(groups: Group[], settings: Record<string, any>): MediaTrackConstraintSet {
  const set: Record<string, any> = {}
  for (const g of groups) {
    for (const c of g.controls) {
      if (settings[c.key] !== undefined) set[c.key] = settings[c.key]
    }
  }
  return set
}

// isAuto: the group's mode is automatic (its range controls follow the
// camera).
export function isAuto(c: Control, settings: Record<string, any>): boolean {
  return c.kind === 'range' && !!c.mode && settings[c.mode] !== undefined && settings[c.mode] !== 'manual'
}
