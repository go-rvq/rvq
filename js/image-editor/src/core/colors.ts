// The colors: the controls of the Cores tab applied to pixels, the histogram,
// and "Corrigir iluminação" — the analysis of a photo taken in the dark or in
// too much light, that sets the controls (it does not paint: the controls
// stay to adjust, and it is one step of the history).

import type {Pixels} from './perspective'

export interface Colors {
  exposure: number // EV, -3..3
  brightness: number // -100..100
  contrast: number // -100..100
  highlights: number // -100..100 (negative recovers them)
  shadows: number // -100..100 (positive opens them)
  gamma: number // 0.2..5, 1: none (above 1 brightens the mid tones)
  saturation: number // -100..100
  temperature: number // -100 (cold) .. 100 (warm)
  tint: number // -100 (green) .. 100 (magenta)
  sharpness: number // 0..100
  grayscale: boolean
  invert: boolean
  // the levels: the black and the white of the photo, taken to 0 and 255
  black: number
  white: number
  // the white balance: a gain per channel (1: none)
  wb: [number, number, number]
}

export function neutral(): Colors {
  return {
    exposure: 0,
    brightness: 0,
    contrast: 0,
    highlights: 0,
    shadows: 0,
    gamma: 1,
    saturation: 0,
    temperature: 0,
    tint: 0,
    sharpness: 0,
    grayscale: false,
    invert: false,
    black: 0,
    white: 255,
    wb: [1, 1, 1],
  }
}

export function isNeutral(c: Colors): boolean {
  return JSON.stringify(c) === JSON.stringify(neutral())
}

// luminance (Rec. 601), 0..255
export function luma(r: number, g: number, b: number): number {
  return 0.299 * r + 0.587 * g + 0.114 * b
}

export interface Stats {
  // the histogram of the luminance, 256 bins
  histogram: number[]
  // per channel, for the panel
  red: number[]
  green: number[]
  blue: number[]
  mean: number
  p1: number
  p99: number
  // the mean of each channel (the white balance)
  avg: [number, number, number]
  // the share of the pixels at 255 (blown) and at 0 (crushed)
  clippedHigh: number
  clippedLow: number
}

function percentile(hist: number[], total: number, q: number): number {
  let acc = 0
  const target = total * q
  for (let i = 0; i < 256; i++) {
    acc += hist[i]
    if (acc >= target) return i
  }
  return 255
}

// analyze: the statistics of the pixels (the transparent ones left out).
export function analyze(px: Pixels): Stats {
  const histogram = new Array(256).fill(0)
  const red = new Array(256).fill(0)
  const green = new Array(256).fill(0)
  const blue = new Array(256).fill(0)
  let n = 0
  let sum = 0
  let sr = 0, sg = 0, sb = 0
  let hi = 0, lo = 0
  const d = px.data
  for (let i = 0; i < d.length; i += 4) {
    if (d[i + 3] === 0) continue
    const r = d[i], g = d[i + 1], b = d[i + 2]
    const l = Math.round(luma(r, g, b))
    histogram[l]++
    red[r]++
    green[g]++
    blue[b]++
    sum += l
    sr += r
    sg += g
    sb += b
    if (l >= 255) hi++
    if (l <= 0) lo++
    n++
  }
  if (n === 0) {
    return {histogram, red, green, blue, mean: 0, p1: 0, p99: 255, avg: [0, 0, 0], clippedHigh: 0, clippedLow: 0}
  }
  return {
    histogram,
    red,
    green,
    blue,
    mean: sum / n,
    p1: percentile(histogram, n, 0.01),
    p99: percentile(histogram, n, 0.99),
    avg: [sr / n, sg / n, sb / n],
    clippedHigh: hi / n,
    clippedLow: lo / n,
  }
}

const clamp = (v: number, a = 0, b = 255) => (v < a ? a : v > b ? b : v)

// fixLighting: the controls that correct a photo taken in the dark or in
// too much light — the levels stretched (its black and white to the ends),
// then the mid tones (gamma) brought to the middle; a dark one gets its
// shadows opened, a bright one its highlights recovered. whiteBalance:
// the gray world (the mean of each channel made the same), taking out the
// yellow of a lamp or the blue of a shade.
export function fixLighting(s: Stats, opts: {whiteBalance?: boolean} = {}): Colors {
  const c = neutral()
  // the levels, never closer than 64 (a flat photo is not made noise)
  let black = s.p1
  let white = s.p99
  if (white - black < 64) {
    const mid = (white + black) / 2
    black = clamp(mid - 32)
    white = clamp(mid + 32)
  }
  c.black = Math.round(black)
  c.white = Math.round(white)
  // where the mean lands after the levels, 0..1: only a photo still dark
  // (or still bright) has its mid tones moved — a sheet of paper, mostly
  // light, is no bright photo
  const m = clamp((s.mean - black) / Math.max(1, white - black), 0.01, 0.99)
  let gamma = 1
  if (m < 0.35) gamma = clamp(Math.log(m) / Math.log(0.45), 1, 2.5)
  else if (m > 0.8 && s.mean > 200) gamma = clamp(Math.log(m) / Math.log(0.65), 0.6, 1)
  c.gamma = Math.round(gamma * 100) / 100
  if (s.mean < 90) c.shadows = Math.round(clamp((90 - s.mean) / 2, 0, 40))
  if (s.clippedHigh > 0.02) c.highlights = -Math.round(clamp(s.clippedHigh * 300, 10, 50))
  if (opts.whiteBalance) c.wb = grayWorld(s.avg)
  return c
}

// grayWorld: the gains that make the mean of the channels the same gray.
export function grayWorld(avg: [number, number, number]): [number, number, number] {
  const [r, g, b] = avg
  const gray = (r + g + b) / 3
  if (!(r > 0 && g > 0 && b > 0)) return [1, 1, 1]
  const k = (v: number) => Math.round(clamp(gray / v, 0.5, 2) * 1000) / 1000
  return [k(r), k(g), k(b)]
}

// document: grays, the contrast up and the background (the paper) cleared —
// a check, a sheet. From the photo's statistics.
export function documentPreset(s: Stats): Colors {
  const c = neutral()
  c.grayscale = true
  c.contrast = 35
  c.black = Math.round(clamp(s.p1 + 10))
  // the paper (the most of the light tones) to white
  c.white = Math.round(clamp(Math.max(c.black + 40, s.p99 - 25)))
  c.sharpness = 30
  return c
}

// lut: the table of the tone controls (levels, exposure, gamma, shadows,
// highlights, brightness, contrast) for one value 0..255.
export function toneLUT(c: Colors): Uint8ClampedArray {
  const lut = new Uint8ClampedArray(256)
  const ev = Math.pow(2, c.exposure)
  const bri = c.brightness / 100
  const con = c.contrast / 100
  const k = con >= 0 ? 1 + con * 2 : 1 + con
  for (let i = 0; i < 256; i++) {
    let v = (i - c.black) / Math.max(1, c.white - c.black)
    v = clamp(v, 0, 1)
    v = clamp(v * ev, 0, 1)
    if (c.gamma !== 1) v = Math.pow(v, 1 / c.gamma)
    // shadows: lift the dark (a curve strongest at the bottom)
    if (c.shadows) v = clamp(v + (c.shadows / 100) * 0.5 * Math.pow(1 - v, 3) * v * 4, 0, 1)
    // highlights: pull (or push) the light
    if (c.highlights) v = clamp(v + (c.highlights / 100) * 0.5 * Math.pow(v, 3) * (1 - v) * 4 + (c.highlights < 0 ? (c.highlights / 100) * 0.05 * v * v : 0), 0, 1)
    v += bri * 0.5
    v = (v - 0.5) * k + 0.5
    lut[i] = Math.round(clamp(v, 0, 1) * 255)
  }
  return lut
}

// applyColors: the controls on the pixels (a copy).
export function applyColors(px: Pixels, c: Colors): Pixels {
  const out = new Uint8ClampedArray(px.data)
  if (isNeutral(c)) return {width: px.width, height: px.height, data: out}
  const lut = toneLUT(c)
  const sat = 1 + c.saturation / 100
  // temperature: red and blue apart; tint: green against magenta
  const t = c.temperature / 100
  const tn = c.tint / 100
  const gr = c.wb[0] * (1 + 0.15 * t + 0.05 * tn)
  const gg = c.wb[1] * (1 - 0.1 * tn)
  const gb = c.wb[2] * (1 - 0.15 * t + 0.05 * tn)
  for (let i = 0; i < out.length; i += 4) {
    let r = clamp(out[i] * gr)
    let g = clamp(out[i + 1] * gg)
    let b = clamp(out[i + 2] * gb)
    r = lut[Math.round(r)]
    g = lut[Math.round(g)]
    b = lut[Math.round(b)]
    if (c.grayscale || sat !== 1) {
      const l = luma(r, g, b)
      const s = c.grayscale ? 0 : sat
      r = clamp(l + (r - l) * s)
      g = clamp(l + (g - l) * s)
      b = clamp(l + (b - l) * s)
    }
    if (c.invert) {
      r = 255 - r
      g = 255 - g
      b = 255 - b
    }
    out[i] = r
    out[i + 1] = g
    out[i + 2] = b
  }
  const res = {width: px.width, height: px.height, data: out}
  return c.sharpness > 0 ? sharpen(res, c.sharpness / 100) : res
}

// sharpen: an unsharp mask of radius 1 (the 3×3 box), amount 0..1.
export function sharpen(px: Pixels, amount: number): Pixels {
  const {width: w, height: h, data: s} = px
  const out = new Uint8ClampedArray(s)
  const k = amount * 1.5
  for (let y = 1; y < h - 1; y++) {
    for (let x = 1; x < w - 1; x++) {
      const o = (y * w + x) * 4
      for (let c = 0; c < 3; c++) {
        let sum = 0
        for (let dy = -1; dy <= 1; dy++) for (let dx = -1; dx <= 1; dx++) sum += s[((y + dy) * w + x + dx) * 4 + c]
        const blur = sum / 9
        out[o + c] = clamp(s[o + c] + (s[o + c] - blur) * k)
      }
    }
  }
  return {width: w, height: h, data: out}
}

export type Preset = 'auto' | 'document' | 'original'

export function preset(p: Preset, s: Stats): Colors {
  switch (p) {
    case 'auto':
      return fixLighting(s, {whiteBalance: true})
    case 'document':
      return documentPreset(s)
    default:
      return neutral()
  }
}
