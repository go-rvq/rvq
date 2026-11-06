<script setup lang="ts">
// vx-image-editor: the image editor of the media box — Crop (fixed,
// perspective, free), Colors, Annotate —, with the zoom, the undo / redo of
// the whole dialog and the marks the same size on the screen at any zoom.
// It works on a URL (the original on the server) or a local File (an upload
// not yet sent): it never calls the server; it emits the result drawn.
import {computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, shallowRef, watch} from 'vue'
import {History, clone} from '@/core/history'
import {
  type View,
  type Point,
  actual,
  clampScale,
  fit,
  imageToScreen,
  nextStep,
  screenToImage,
  zoomAt,
  HIT_RADIUS,
  HIT_RADIUS_TOUCH,
} from '@/core/view'
import {
  type Rect,
  alignEdge,
  bounds as pointsBounds,
  center,
  clampRect,
  dist,
  fitAspect,
  insertPoint,
  nearestEdge,
  nearestPoint,
  normDeg,
  rectPoints,
  removePoint,
  resizeRect,
  rotatePoint,
  snapRotation,
  stepRotation,
  usefulArea,
  type TurnedRect,
  turnedCorners,
  toLocal,
  resizeTurned,
  snapAngle,
} from '@/core/geometry'
import {type Quad, defaultQuad} from '@/core/perspective'
import {detectCorners} from '@/core/detect'
import {type Colors, type Preset, analyze, fixLighting, neutral, preset as colorPreset} from '@/core/colors'
import {
  type Annotation,
  type Style,
  type Tool,
  TOOLS,
  bounds as annBounds,
  create as createAnn,
  defaultStyle,
  drawAll,
  duplicate as duplicateAnn,
  move as moveAnn,
  nextNumber,
  resize as resizeAnn,
  topmost,
} from '@/core/annotations'
import {
  type CropMode,
  type Edit,
  type Tab,
  canvas as newCanvas,
  changes,
  colored,
  croppedSize,
  cropped,
  annotated,
  oriented,
  normalized,
  parseEdit,
  render,
  resultName,
  resultType,
  toBlob,
} from '@/core/edit'
import {loadBitmap} from '@/core/exif'
import {type Messages, messagesFor, pageLocale} from './messages'
import {getCurrentInstance} from 'vue'

export interface SaveEvent {
  // the image made, when the edit makes one (else null: only cropOption)
  blob: Blob | null
  // the same, as a File named name (what a form posts)
  file: File | null
  name: string
  type: string
  edit: Edit
  // the fixed crop of the size, in the pixels of the result
  cropOption: Rect | null
}

const props = withDefaults(
  defineProps<{
    src?: string
    file?: Blob | null
    name?: string
    type?: string
    edit?: Partial<Edit> | null
    // the ratio of the size (width / height): the fixed crop; 0, none
    aspect?: number
    tab?: Tab
    mode?: CropMode
    messages?: Partial<Messages>
    // the texts' language (pt, en); the page's when not given
    locale?: string
    // the upload: the buttons "Save and upload", "Upload without editing";
    // queue: "1 of 3"; camera: "Take another"
    upload?: boolean
    queue?: {index: number; total: number} | null
    camera?: boolean
    // a photo larger than this (pixels) made smaller first; 0: never
    maxPixels?: number
    // the closing buttons (cancel, save…) out of the bar: the dialog's
    // header has them (askCancel, save)
    actionsOutside?: boolean
  }>(),
  {
    src: '',
    file: null,
    name: '',
    type: '',
    edit: null,
    aspect: 0,
    tab: 'recortar',
    mode: undefined,
    messages: () => ({}),
    locale: '',
    upload: false,
    queue: null,
    camera: false,
    maxPixels: 0,
    actionsOutside: false,
  },
)

const emit = defineEmits<{
  save: [SaveEvent]
  cancel: []
  skip: []
  skipAll: []
  retake: []
}>()

const app = getCurrentInstance()?.appContext.app
const m = computed<Messages>(() => ({...messagesFor(props.locale || pageLocale(app)), ...props.messages}))
// the fixed crop: always; of a size (its ratio, aspect), it is the size's
// crop (its CropOption: no new image); else a crop of the image, of the ratio
// chosen (free, 1:1, 4:3…)
const FIXO_IS_SIZE = props.aspect > 0
const RATIOS: {value: number; label: string}[] = [
  {value: 0, label: ''},
  {value: -1, label: ''},
  {value: 1, label: '1:1'},
  {value: 4 / 3, label: '4:3'},
  {value: 3 / 2, label: '3:2'},
  {value: 16 / 9, label: '16:9'},
]
// the ratio of the fixed crop now: the size's, else the one chosen (-1: the
// image's)
function fixoRatio(e: Edit = edit.value): number {
  if (props.aspect > 0) return props.aspect
  const r = e.fixoRatio ?? 0
  if (r === -1) return orientedSize.value.w / Math.max(1, orientedSize.value.h)
  return r
}
// a ratio chosen and turned (portrait): 3:4 is 4:3 turned
function sameRatio(a: number, b: number) {
  return a > 0 && b > 0 && (Math.abs(a - b) < 1e-3 || Math.abs(a - 1 / b) < 1e-3)
}
const PREVIEW_MAX = 1400

// --- the image ---

const loading = ref(true)
const loadError = ref('')
const saving = ref(false)
const source = shallowRef<ImageBitmap | HTMLImageElement | null>(null)
const sourceType = ref('image/jpeg')
// the image turned (rotate, flip): full and a preview (kPrev: its scale)
const orientedFull = shallowRef<HTMLCanvasElement | null>(null)
const orientedPrev = shallowRef<HTMLCanvasElement | null>(null)
const kPrev = ref(1)
// the image as cut, colored: the preview of the Colors and Annotate tabs
const cutPrev = shallowRef<HTMLCanvasElement | null>(null)
const colorPrev = shallowRef<HTMLCanvasElement | null>(null)

// --- the edit and its history ---

function initialMode(): CropMode {
  return props.mode || 'fixo'
}

const edit = ref<Edit>(parseEdit(props.edit, initialMode()))
const history = new History<Edit>(edit.value)
const counts = reactive({undo: 0, redo: 0})

function syncCounts() {
  counts.undo = history.undoCount
  counts.redo = history.redoCount
}

// commit: a step (a click, a value set)
function commit(next: Edit) {
  edit.value = next
  history.push(next)
  syncCounts()
}

// change: during a drag (begin … end: one step)
function change(next: Edit) {
  edit.value = next
  history.set(next)
}

function undo() {
  const s = history.undo()
  if (s) edit.value = s
  syncCounts()
}

function redo() {
  const s = history.redo()
  if (s) edit.value = s
  syncCounts()
}

const tab = ref<Tab>(props.tab)

// a phone's (narrower than 960 px): the bar of icons, the side panel a
// sheet from the bottom, opened by its button
const mq = typeof window !== 'undefined' && window.matchMedia ? window.matchMedia('(max-width: 959px)') : null
const mobile = ref(!!mq?.matches)
const showSide = ref(false)
const onMq = (e: MediaQueryListEvent) => (mobile.value = e.matches)
mq?.addEventListener?.('change', onMq)
onBeforeUnmount(() => mq?.removeEventListener?.('change', onMq))
const tool = ref<Tool>('select')
const style = reactive<Style>(defaultStyle())
const selectedPoint = ref(-1)
const selectedEdge = ref(-1)
const selectedAnn = ref<string | null>(null)
const compare = ref(false)
const whiteBalance = ref(true)

// --- the stage and the zoom ---

const stage = ref<HTMLDivElement | null>(null)
const cv = ref<HTMLCanvasElement | null>(null)
const lupe = ref<HTMLCanvasElement | null>(null)
const stageSize = reactive({w: 800, h: 500})
const view = ref<View>({scale: 1, x: 0, y: 0})
const spaceDown = ref(false)
const touch = ref(false)

// the size of the image the stage shows (its pixels: the marks' and the
// annotations' coordinates)
const orientedSize = computed(() => ({w: orientedFull.value?.width || 1, h: orientedFull.value?.height || 1}))
const cutSize = computed(() => {
  const s = croppedSize(edit.value, orientedSize.value.w, orientedSize.value.h, FIXO_IS_SIZE)
  return {w: s.width, h: s.height}
})
const imgSize = computed(() => (tab.value === 'recortar' ? orientedSize.value : cutSize.value))

function doFit() {
  view.value = fit(stageSize.w, stageSize.h, imgSize.value.w, imgSize.value.h)
}

function doActual() {
  view.value = actual(stageSize.w, stageSize.h, imgSize.value.w, imgSize.value.h)
}

function zoomStep(dir: 1 | -1) {
  view.value = zoomAt(view.value, nextStep(view.value.scale, dir), {x: stageSize.w / 2, y: stageSize.h / 2})
}

const zoomText = ref('100')
watch(
  () => view.value.scale,
  (s) => (zoomText.value = String(Math.round(s * 100))),
  {immediate: true},
)

function zoomTyped() {
  const v = parseFloat(zoomText.value)
  if (v > 0) view.value = zoomAt(view.value, clampScale(v / 100), {x: stageSize.w / 2, y: stageSize.h / 2})
  else zoomText.value = String(Math.round(view.value.scale * 100))
}

const S = (p: Point) => imageToScreen(view.value, p)
const hitR = () => (touch.value ? HIT_RADIUS_TOUCH : HIT_RADIUS) / view.value.scale

// --- loading ---

async function load() {
  loading.value = true
  loadError.value = ''
  try {
    const src = props.file || props.src
    if (!src) throw new Error('image-editor: no image')
    if (props.file) sourceType.value = props.file.type || props.type || 'image/jpeg'
    else sourceType.value = props.type || guessType(props.src)
    source.value = await loadBitmap(src, props.maxPixels)
    makeOriented()
    await nextTick()
    measure()
    doFit()
  } catch (e: any) {
    loadError.value = String(e?.message || e)
  } finally {
    loading.value = false
  }
}

function guessType(url: string): string {
  const ext = url.split('?')[0].split('.').pop()?.toLowerCase()
  return ext === 'png' ? 'image/png' : ext === 'webp' ? 'image/webp' : 'image/jpeg'
}

function makeOriented() {
  if (!source.value) return
  const full = oriented(source.value, edit.value.rotate, edit.value.flipH)
  orientedFull.value = full
  const k = Math.min(1, PREVIEW_MAX / Math.max(full.width, full.height))
  kPrev.value = k
  if (k < 1) {
    const p = newCanvas(full.width * k, full.height * k)
    const ctx = p.getContext('2d')!
    ctx.imageSmoothingQuality = 'high'
    ctx.drawImage(full, 0, 0, p.width, p.height)
    orientedPrev.value = p
  } else {
    orientedPrev.value = full
  }
  makeCut()
}

// the cut and its colors (a preview), again when the geometry or the colors
// change — at most once a frame
let cutKey = ''
function makeCut() {
  if (!orientedPrev.value) return
  const e = edit.value
  const key = JSON.stringify([e.rotate, e.flipH, e.mode, e.fixo, e.retangulo, e.perspectiva, e.livre])
  if (key !== cutKey || !cutPrev.value) {
    cutKey = key
    cutPrev.value = cropped(orientedPrev.value, e, FIXO_IS_SIZE, kPrev.value)
  }
  colorPrev.value = colored(cutPrev.value!, e)
  stats.value = null
}

let cutRaf = 0
function scheduleCut() {
  if (cutRaf) return
  cutRaf = requestAnimationFrame(() => {
    cutRaf = 0
    makeCut()
    draw()
  })
}

watch(
  () => [edit.value.rotate, edit.value.flipH],
  () => {
    makeOriented()
    draw()
  },
)
watch(
  () => [edit.value.mode, edit.value.fixo, edit.value.retangulo, edit.value.perspectiva, edit.value.livre, edit.value.colors],
  () => scheduleCut(),
  {deep: true},
)
watch(tab, async () => {
  selectedPoint.value = -1
  selectedEdge.value = -1
  if (tab.value !== 'anotar') selectedAnn.value = null
  makeCut()
  await nextTick()
  doFit()
})

// --- drawing the stage ---

let drawRaf = 0
function draw() {
  if (drawRaf) return
  drawRaf = requestAnimationFrame(() => {
    drawRaf = 0
    paint()
  })
}

function checker(ctx: CanvasRenderingContext2D) {
  const c = newCanvas(16, 16)
  const g = c.getContext('2d')!
  g.fillStyle = '#ccc'
  g.fillRect(0, 0, 16, 16)
  g.fillStyle = '#fff'
  g.fillRect(0, 0, 8, 8)
  g.fillRect(8, 8, 8, 8)
  return ctx.createPattern(c, 'repeat')!
}

// the image the tab shows, and its pixels per pixel of the image
function baseOf(): {img: HTMLCanvasElement | null; k: number} {
  if (tab.value === 'recortar') return {img: orientedPrev.value, k: kPrev.value}
  const img = compare.value ? cutPrev.value : colorPrev.value
  return {img, k: img ? img.width / Math.max(1, cutSize.value.w) : 1}
}

function paint() {
  const c = cv.value
  if (!c) return
  const dpr = window.devicePixelRatio || 1
  if (c.width !== Math.round(stageSize.w * dpr) || c.height !== Math.round(stageSize.h * dpr)) {
    c.width = Math.round(stageSize.w * dpr)
    c.height = Math.round(stageSize.h * dpr)
  }
  const ctx = c.getContext('2d')!
  ctx.setTransform(1, 0, 0, 1, 0, 0)
  ctx.fillStyle = '#3a3a3a'
  ctx.fillRect(0, 0, c.width, c.height)
  const {img, k} = baseOf()
  if (!img) return
  const v = view.value
  ctx.setTransform(dpr * v.scale, 0, 0, dpr * v.scale, dpr * v.x, dpr * v.y)
  const {w, h} = imgSize.value
  ctx.save()
  ctx.fillStyle = checker(ctx)
  ctx.scale(1 / v.scale, 1 / v.scale)
  ctx.fillRect(0, 0, w * v.scale, h * v.scale)
  ctx.restore()
  ctx.imageSmoothingEnabled = v.scale < 2
  ctx.drawImage(img, 0, 0, img.width, img.height, 0, 0, w, h)
  if (tab.value === 'anotar' && !compare.value) {
    const list = draft.value ? [...edit.value.annotations.filter((a) => a.id !== draft.value!.id), draft.value] : edit.value.annotations
    drawAll(ctx, list, img, k)
  }
}

watch([view, stageSize, () => edit.value.annotations, compare, colorPrev], () => draw(), {deep: true})

// --- the marks (screen coordinates: the same size at any zoom) ---

const HANDLES = ['nw', 'n', 'ne', 'e', 'se', 's', 'sw', 'w'] as const

function handlePoints(r: Rect): Record<string, Point> {
  const x0 = r.x, y0 = r.y, x1 = r.x + r.width, y1 = r.y + r.height
  const mx = (x0 + x1) / 2, my = (y0 + y1) / 2
  return {nw: {x: x0, y: y0}, n: {x: mx, y: y0}, ne: {x: x1, y: y0}, e: {x: x1, y: my}, se: {x: x1, y: y1}, s: {x: mx, y: y1}, sw: {x: x0, y: y1}, w: {x: x0, y: my}}
}

const pathOf = (pts: Point[]) => pts.map((p, i) => `${i ? 'L' : 'M'}${S(p).x.toFixed(1)},${S(p).y.toFixed(1)}`).join('') + 'Z'

// the shade outside the crop
const maskPath = computed(() => {
  if (tab.value !== 'recortar') return ''
  const e = edit.value
  let pts: Point[] | null = null
  if (e.mode === 'fixo' && e.fixo) pts = rectPoints(e.fixo)
  if (e.mode === 'perspectiva' && e.perspectiva) pts = e.perspectiva
  if (e.mode === 'livre' && e.livre && e.livre.points.length >= 3) pts = e.livre.points
  if (e.mode === 'retangulo' && e.retangulo) pts = turnedCorners(e.retangulo)
  if (!pts) return ''
  return `M0,0H${stageSize.w}V${stageSize.h}H0Z` + pathOf(pts)
})

const fixoMarks = computed(() => {
  const r = edit.value.fixo
  if (tab.value !== 'recortar' || edit.value.mode !== 'fixo' || !r) return null
  const hp = handlePoints(r)
  const thirds: string[] = []
  for (const t of [1 / 3, 2 / 3]) {
    const a = S({x: r.x + r.width * t, y: r.y}), b = S({x: r.x + r.width * t, y: r.y + r.height})
    thirds.push(`M${a.x},${a.y}L${b.x},${b.y}`)
    const c = S({x: r.x, y: r.y + r.height * t}), d = S({x: r.x + r.width, y: r.y + r.height * t})
    thirds.push(`M${c.x},${c.y}L${d.x},${d.y}`)
  }
  return {outline: pathOf(rectPoints(r)), handles: HANDLES.map((h) => ({h, ...S(hp[h])})), guides: thirds.join('')}
})

function lerp(a: Point, b: Point, t: number): Point {
  return {x: a.x + (b.x - a.x) * t, y: a.y + (b.y - a.y) * t}
}

const perspMarks = computed(() => {
  const q = edit.value.perspectiva
  if (tab.value !== 'recortar' || edit.value.mode !== 'perspectiva' || !q) return null
  const grid: string[] = []
  for (const t of [1 / 3, 2 / 3]) {
    const a = S(lerp(q[0], q[1], t)), b = S(lerp(q[3], q[2], t))
    grid.push(`M${a.x},${a.y}L${b.x},${b.y}`)
    const c = S(lerp(q[0], q[3], t)), d = S(lerp(q[1], q[2], t))
    grid.push(`M${c.x},${c.y}L${d.x},${d.y}`)
  }
  return {outline: pathOf(q), corners: q.map((p, i) => ({i, ...S(p)})), grid: grid.join('')}
})

const ROT_HANDLE = 36
// the free rectangle: its outline, its handles (corners, sides), the handle
// that turns it (above its top side), its thirds
const retMarks = computed(() => {
  const r = edit.value.retangulo
  if (tab.value !== 'recortar' || edit.value.mode !== 'retangulo' || !r) return null
  const c = {x: r.cx, y: r.cy}
  const local = handlePoints({x: -r.width / 2, y: -r.height / 2, width: r.width, height: r.height})
  const toImg = (q: Point) => rotatePoint({x: r.cx + q.x, y: r.cy + q.y}, r.angle, c)
  const handles = HANDLES.map((h) => ({h, ...S(toImg(local[h]))}))
  const top = S(toImg(local.n))
  // the turning handle: ROT_HANDLE px out of the top side, on the screen
  const ctr = S(c)
  const dx = top.x - ctr.x
  const dy = top.y - ctr.y
  const len = Math.hypot(dx, dy) || 1
  const rot = {cx: top.x, cy: top.y, x: top.x + (dx / len) * ROT_HANDLE, y: top.y + (dy / len) * ROT_HANDLE}
  const guides: string[] = []
  for (const t of [1 / 3, 2 / 3]) {
    const a = S(toImg({x: -r.width / 2 + r.width * t, y: -r.height / 2})), b = S(toImg({x: -r.width / 2 + r.width * t, y: r.height / 2}))
    guides.push(`M${a.x},${a.y}L${b.x},${b.y}`)
    const d = S(toImg({x: -r.width / 2, y: -r.height / 2 + r.height * t})), e = S(toImg({x: r.width / 2, y: -r.height / 2 + r.height * t}))
    guides.push(`M${d.x},${d.y}L${e.x},${e.y}`)
  }
  return {outline: pathOf(turnedCorners(r)), handles, rot, guides: guides.join('')}
})

const livreMarks = computed(() => {
  const l = edit.value.livre
  if (tab.value !== 'recortar' || edit.value.mode !== 'livre' || !l) return null
  const pts = l.points
  const edges = pts.length >= 2 ? pts.map((p, i) => ({i, a: S(p), b: S(pts[(i + 1) % pts.length])})) : []
  let rot = null as null | {cx: number; cy: number; x: number; y: number; area: string; angle: number}
  if (pts.length >= 3) {
    const b = pointsBounds(pts)
    const top = S({x: b.x + b.width / 2, y: b.y})
    // the useful area after the turn, shown turned back on the image
    const area = usefulArea(pts, l.rotation)
    const c = center(pts)
    const back = rectPoints(area).map((p) => rotatePoint(p, -l.rotation, c))
    rot = {cx: top.x, cy: top.y, x: top.x, y: top.y - ROT_HANDLE, area: pathOf(back), angle: l.rotation}
  }
  return {
    outline: pts.length >= 2 ? pathOf(pts) : '',
    points: pts.map((p, i) => ({i, ...S(p)})),
    edges,
    rot,
  }
})

const annMarks = computed(() => {
  if (tab.value !== 'anotar' || !selectedAnn.value) return null
  const a = edit.value.annotations.find((x) => x.id === selectedAnn.value)
  if (!a) return null
  const b = annBounds(a)
  const hp = handlePoints(b)
  return {outline: pathOf(rectPoints(b)), handles: HANDLES.map((h) => ({h, ...S(hp[h])}))}
})

// --- the modes ---

function setMode(mode: CropMode) {
  const e = clone(edit.value)
  e.mode = mode
  const {w, h} = orientedSize.value
  if (mode === 'fixo' && !e.fixo) e.fixo = fitAspect(w, h, fixoRatio(e))
  if (mode === 'perspectiva' && !e.perspectiva) e.perspectiva = defaultQuad(w, h)
  if (mode === 'livre' && !e.livre) e.livre = {points: [], rotation: 0}
  if (mode === 'retangulo' && !e.retangulo) e.retangulo = defaultTurned(w, h)
  selectedPoint.value = -1
  selectedEdge.value = -1
  commit(e)
}

function resetMode() {
  const e = clone(edit.value)
  const {w, h} = orientedSize.value
  if (e.mode === 'fixo') e.fixo = fitAspect(w, h, fixoRatio(e))
  if (e.mode === 'perspectiva') e.perspectiva = defaultQuad(w, h)
  if (e.mode === 'livre') e.livre = {points: [], rotation: 0}
  if (e.mode === 'retangulo') e.retangulo = defaultTurned(w, h)
  selectedPoint.value = -1
  selectedEdge.value = -1
  commit(e)
}

// setRatio: the fixed crop of the ratio r (0: free; -1: the image's) — the
// largest of it around the crop now, in the image
function setRatio(r: number) {
  const e = clone(edit.value)
  e.fixoRatio = r
  const ratio = fixoRatio(e)
  const {w, h} = orientedSize.value
  if (ratio > 0) {
    const cur = e.fixo || {x: 0, y: 0, width: w, height: h}
    const cx = cur.x + cur.width / 2
    const cy = cur.y + cur.height / 2
    const f = fitAspect(w, h, ratio)
    e.fixo = clampRect({x: cx - f.width / 2, y: cy - f.height / 2, width: f.width, height: f.height}, w, h)
  } else if (!e.fixo) {
    e.fixo = {x: 0, y: 0, width: w, height: h}
  }
  commit(e)
}

// swapRatio: the ratio chosen turned (4:3 ↔ 3:4)
function swapRatio() {
  const r = fixoRatio()
  if (!(r > 0) || props.aspect > 0) return
  setRatio(1 / r)
}

const ratioNow = computed(() => (props.aspect > 0 ? props.aspect : edit.value.fixoRatio ?? 0))

// defaultTurned: the free rectangle at the start — the image, a little in
function defaultTurned(w: number, h: number): TurnedRect {
  return {cx: w / 2, cy: h / 2, width: w * 0.8, height: h * 0.8, angle: 0}
}

function removeSelectedPoint() {
  const l = edit.value.livre
  if (!l || selectedPoint.value < 0) return
  const pts = removePoint(l.points, selectedPoint.value)
  if (!pts) return
  const e = clone(edit.value)
  e.livre!.points = pts
  selectedPoint.value = -1
  commit(e)
}

function doAlignEdge() {
  const l = edit.value.livre
  if (!l || selectedEdge.value < 0) return
  const e = clone(edit.value)
  e.livre!.rotation = Math.round(alignEdge(l.points, selectedEdge.value, l.rotation) * 100) / 100
  commit(e)
}

function setAngle(v: string | number) {
  const a = parseFloat(String(v))
  if (!isFinite(a)) return
  const e = clone(edit.value)
  if (e.mode === 'retangulo' && e.retangulo) e.retangulo.angle = normDeg(a)
  else if (e.livre) e.livre.rotation = normDeg(a)
  else return
  commit(e)
}

function detectEdges() {
  const src = orientedPrev.value
  if (!src) return
  const d = src.getContext('2d')!.getImageData(0, 0, src.width, src.height)
  const q = detectCorners({width: d.width, height: d.height, data: d.data})
  if (!q) return
  const e = clone(edit.value)
  e.perspectiva = q.map((p) => ({x: p.x / kPrev.value, y: p.y / kPrev.value})) as Quad
  commit(e)
}

// a turn of 90° or a mirror: the shapes made on the image before start over
function turn(deg: 90 | -90) {
  const e = clone(edit.value)
  e.rotate = (((e.rotate + deg) % 360) + 360) % 360 as Edit['rotate']
  e.fixo = null
  e.perspectiva = null
  e.livre = e.mode === 'livre' ? {points: [], rotation: 0} : null
  commit(e)
  nextTick(() => {
    setMode(edit.value.mode)
    doFit()
  })
}

function mirror() {
  const e = clone(edit.value)
  e.flipH = !e.flipH
  e.fixo = null
  e.perspectiva = null
  e.livre = e.mode === 'livre' ? {points: [], rotation: 0} : null
  commit(e)
  nextTick(() => setMode(edit.value.mode))
}

function backToOriginal() {
  const e = parseEdit(null, initialMode())
  commit(e)
  nextTick(() => {
    setMode(e.mode)
    doFit()
  })
}

// --- the colors ---

const stats = shallowRef<ReturnType<typeof analyze> | null>(null)
function statsOf(c: HTMLCanvasElement | null) {
  if (!c) return null
  const k = Math.min(1, 256 / Math.max(c.width, c.height))
  const s = newCanvas(c.width * k, c.height * k)
  s.getContext('2d')!.drawImage(c, 0, 0, s.width, s.height)
  const d = s.getContext('2d')!.getImageData(0, 0, s.width, s.height)
  return analyze({width: d.width, height: d.height, data: d.data})
}

const shownStats = computed(() => {
  if (tab.value !== 'cores') return null
  if (!stats.value) stats.value = statsOf(compare.value ? cutPrev.value : colorPrev.value)
  return stats.value
})
watch(compare, () => (stats.value = null))

const histPaths = computed(() => {
  const s = shownStats.value
  if (!s) return null
  const W = 256, H = 80
  const path = (h: number[]) => {
    const max = Math.max(1, ...h.slice(1, 255))
    return `M0,${H}` + h.map((v, i) => `L${i},${(H - Math.min(H, (v / max) * H)).toFixed(1)}`).join('') + `L${W},${H}Z`
  }
  return {luma: path(s.histogram), r: path(s.red), g: path(s.green), b: path(s.blue)}
})

function setColors(c: Partial<Colors>, drag = false) {
  const e = clone(edit.value)
  e.colors = {...e.colors, ...c}
  if (drag) {
    history.begin()
    change(e)
  } else commit(e)
}

function endColorDrag() {
  history.end()
  syncCounts()
}

function doFixLighting() {
  const s = statsOf(cutPrev.value)
  if (!s) return
  const e = clone(edit.value)
  e.colors = fixLighting(s, {whiteBalance: whiteBalance.value})
  commit(e)
}

function applyPreset(p: Preset) {
  const s = statsOf(cutPrev.value)
  if (!s) return
  const e = clone(edit.value)
  e.colors = colorPreset(p, s)
  if (p === 'auto' && !whiteBalance.value) e.colors.wb = [1, 1, 1]
  commit(e)
}

const SLIDERS: {key: keyof Colors; label: keyof Messages; min: number; max: number; step: number}[] = [
  {key: 'exposure', label: 'exposure', min: -3, max: 3, step: 0.05},
  {key: 'brightness', label: 'brightness', min: -100, max: 100, step: 1},
  {key: 'contrast', label: 'contrast', min: -100, max: 100, step: 1},
  {key: 'highlights', label: 'highlights', min: -100, max: 100, step: 1},
  {key: 'shadows', label: 'shadows', min: -100, max: 100, step: 1},
  {key: 'gamma', label: 'gamma', min: 0.2, max: 5, step: 0.01},
  {key: 'saturation', label: 'saturation', min: -100, max: 100, step: 1},
  {key: 'temperature', label: 'temperature', min: -100, max: 100, step: 1},
  {key: 'tint', label: 'tint', min: -100, max: 100, step: 1},
  {key: 'sharpness', label: 'sharpness', min: 0, max: 100, step: 1},
]

// --- the annotations ---

const TOOL_ICONS: Record<Tool, string> = {
  select: 'mdi-cursor-default-outline',
  pen: 'mdi-pencil',
  highlighter: 'mdi-marker',
  line: 'mdi-vector-line',
  arrow: 'mdi-arrow-top-right',
  rect: 'mdi-rectangle-outline',
  ellipse: 'mdi-ellipse-outline',
  text: 'mdi-format-text',
  number: 'mdi-numeric-1-circle-outline',
  pixelate: 'mdi-checkerboard',
  blur: 'mdi-blur',
}
const TOOL_LABELS: Record<Tool, keyof Messages> = {
  select: 'toolSelect',
  pen: 'toolPen',
  highlighter: 'toolHighlighter',
  line: 'toolLine',
  arrow: 'toolArrow',
  rect: 'toolRect',
  ellipse: 'toolEllipse',
  text: 'toolText',
  number: 'toolNumber',
  pixelate: 'toolPixelate',
  blur: 'toolBlur',
}

const draft = shallowRef<Annotation | null>(null)
const selected = computed(() => edit.value.annotations.find((a) => a.id === selectedAnn.value) || null)

// the properties shown: the selected annotation's, else the style of the
// next ones
const props_ = computed<Style>(() => (selected.value ? (selected.value as Style) : style))

function setProp<K extends keyof Style | 'text'>(k: K, v: any) {
  if (k !== 'text') (style as any)[k] = v
  if (!selected.value) return
  const e = clone(edit.value)
  const a = e.annotations.find((x) => x.id === selectedAnn.value)!
  ;(a as any)[k] = v
  commit(e)
}

function deleteAnn(id = selectedAnn.value) {
  if (!id) return
  const e = clone(edit.value)
  e.annotations = e.annotations.filter((a) => a.id !== id)
  if (selectedAnn.value === id) selectedAnn.value = null
  commit(e)
}

function duplicateSelected() {
  if (!selected.value) return
  const e = clone(edit.value)
  const d = duplicateAnn(selected.value, e.annotations)
  e.annotations.push(d)
  selectedAnn.value = d.id
  commit(e)
}

function moveLayer(id: string, dir: 1 | -1) {
  const e = clone(edit.value)
  const i = e.annotations.findIndex((a) => a.id === id)
  const j = i + dir
  if (i < 0 || j < 0 || j >= e.annotations.length) return
  ;[e.annotations[i], e.annotations[j]] = [e.annotations[j], e.annotations[i]]
  commit(e)
}

const textInput = ref<any>(null)

// --- the pointer ---

type Drag =
  | {kind: 'pan'; sx: number; sy: number; vx: number; vy: number}
  | {kind: 'fixo-move'; start: Point; rect: Rect}
  | {kind: 'fixo-resize'; handle: string}
  | {kind: 'persp'; i: number}
  | {kind: 'livre-point'; i: number}
  | {kind: 'livre-rotate'; c: Point; a0: number; r0: number}
  | {kind: 'ret-move'; start: Point; r: TurnedRect}
  | {kind: 'ret-resize'; handle: string}
  | {kind: 'ret-rotate'; c: Point; a0: number; r0: number}
  | {kind: 'ann-create'}
  | {kind: 'ann-move'; start: Point; a: Annotation}
  | {kind: 'ann-resize'; handle: string; a: Annotation; box: Rect}

let drag: Drag | null = null
const pointers = new Map<number, Point>()
let pinch: {d: number; s: number; mid: Point; v: View} | null = null
const lupeAt = ref<{x: number; y: number; img: Point} | null>(null)

function local(ev: PointerEvent | MouseEvent | WheelEvent): Point {
  const r = stage.value!.getBoundingClientRect()
  return {x: ev.clientX - r.left, y: ev.clientY - r.top}
}

function hitHandle(r: Rect, p: Point): string | null {
  const hp = handlePoints(r)
  const rad = hitR()
  for (const h of HANDLES) if (dist(hp[h], p) <= rad) return h
  return null
}

function inside(r: Rect, p: Point) {
  return p.x >= r.x && p.x <= r.x + r.width && p.y >= r.y && p.y <= r.y + r.height
}

function onPointerDown(ev: PointerEvent) {
  if (loading.value) return
  touch.value = ev.pointerType === 'touch'
  // the moves of this pointer here even out of the stage (a pointer no more
  // active — gone already — has none to capture)
  try {
    ;(ev.currentTarget as Element).setPointerCapture?.(ev.pointerId)
  } catch {
    /* no capture: the moves over the stage still come */
  }
  const sp = local(ev)
  pointers.set(ev.pointerId, sp)
  if (pointers.size === 2) {
    // a pinch: the drag in course undone as a step
    if (drag && drag.kind !== 'pan') history.end()
    const [a, b] = [...pointers.values()]
    pinch = {d: dist(a, b), s: view.value.scale, mid: {x: (a.x + b.x) / 2, y: (a.y + b.y) / 2}, v: {...view.value}}
    drag = null
    return
  }
  if (ev.button === 1 || spaceDown.value || tab.value === 'cores') {
    drag = {kind: 'pan', sx: sp.x, sy: sp.y, vx: view.value.x, vy: view.value.y}
    ev.preventDefault()
    return
  }
  if (ev.button !== 0) return
  const p = screenToImage(view.value, sp)
  const e = edit.value
  if (tab.value === 'recortar') {
    if (e.mode === 'fixo' && e.fixo) {
      const h = hitHandle(e.fixo, p)
      if (h) {
        history.begin()
        drag = {kind: 'fixo-resize', handle: h}
      } else if (inside(e.fixo, p)) {
        history.begin()
        drag = {kind: 'fixo-move', start: p, rect: {...e.fixo}}
      }
    } else if (e.mode === 'retangulo' && e.retangulo) {
      const r = e.retangulo
      const marks = retMarks.value!
      const rad = touch.value ? HIT_RADIUS_TOUCH : HIT_RADIUS
      // the nearest handle (the screen's: the same size at any zoom)
      const hd = marks.handles.reduce((a, b) => (dist(sp, a) <= dist(sp, b) ? a : b))
      if (dist(sp, {x: marks.rot.x, y: marks.rot.y}) <= rad) {
        const cs = S({x: r.cx, y: r.cy})
        history.begin()
        drag = {kind: 'ret-rotate', c: cs, a0: Math.atan2(sp.y - cs.y, sp.x - cs.x), r0: r.angle}
      } else if (dist(sp, hd) <= rad) {
        history.begin()
        drag = {kind: 'ret-resize', handle: hd.h}
      } else {
        const l = toLocal(r, p)
        if (Math.abs(l.x) <= r.width / 2 && Math.abs(l.y) <= r.height / 2) {
          history.begin()
          drag = {kind: 'ret-move', start: p, r: {...r}}
        }
      }
    } else if (e.mode === 'perspectiva' && e.perspectiva) {
      const i = nearestPoint(e.perspectiva, p, hitR())
      if (i >= 0) {
        history.begin()
        drag = {kind: 'persp', i}
        selectedPoint.value = i
      }
    } else if (e.mode === 'livre') {
      const l = e.livre || {points: [], rotation: 0}
      const marks = livreMarks.value
      if (marks?.rot && dist(sp, {x: marks.rot.x, y: marks.rot.y}) <= (touch.value ? HIT_RADIUS_TOUCH : HIT_RADIUS)) {
        const c = center(l.points)
        const cs = S(c)
        history.begin()
        drag = {kind: 'livre-rotate', c: cs, a0: Math.atan2(sp.y - cs.y, sp.x - cs.x), r0: l.rotation}
        return
      }
      const i = nearestPoint(l.points, p, hitR())
      if (i >= 0) {
        selectedPoint.value = i
        selectedEdge.value = -1
        history.begin()
        drag = {kind: 'livre-point', i}
        return
      }
      if (l.points.length >= 3) {
        const ei = nearestEdge(l.points, p, hitR())
        if (ei >= 0) {
          selectedEdge.value = ei
          selectedPoint.value = -1
          return
        }
      }
      // a point: on the nearest edge; dragged at once (one step with it)
      history.begin()
      const ins = insertPoint(l.points, p)
      const ne = clone(e)
      ne.livre = {points: ins.points, rotation: l.rotation}
      change(ne)
      selectedPoint.value = ins.index
      selectedEdge.value = -1
      drag = {kind: 'livre-point', i: ins.index}
    }
    if (drag) updateLupe(sp, p)
    return
  }
  if (tab.value === 'anotar') {
    const tol = hitR() / 2
    if (tool.value === 'select') {
      if (selected.value) {
        const h = hitHandle(annBounds(selected.value), p)
        if (h) {
          history.begin()
          drag = {kind: 'ann-resize', handle: h, a: clone(selected.value), box: annBounds(selected.value)}
          return
        }
      }
      const a = topmost(e.annotations, p, tol)
      selectedAnn.value = a ? a.id : null
      if (a) {
        history.begin()
        drag = {kind: 'ann-move', start: p, a: clone(a)}
      }
      return
    }
    // a new annotation: drawn while dragging, one step at the end
    const a = createAnn(tool.value as Annotation['type'], p, {...style}, nextNumber(e.annotations))
    draft.value = a
    drag = {kind: 'ann-create'}
    if (a.type === 'text' || a.type === 'number') finishDraft()
  }
}

function updateLupe(sp: Point, p: Point) {
  lupeAt.value = {x: sp.x, y: sp.y, img: p}
  drawLupe()
}

function drawLupe() {
  const l = lupeAt.value
  const c = lupe.value
  const src = orientedPrev.value
  if (!l || !c || !src) return
  const ctx = c.getContext('2d')!
  const size = c.width
  const zoom = Math.max(4, view.value.scale * 3)
  const k = kPrev.value
  const half = size / 2 / zoom
  ctx.clearRect(0, 0, size, size)
  ctx.imageSmoothingEnabled = false
  ctx.drawImage(src, (l.img.x - half) * k, (l.img.y - half) * k, 2 * half * k, 2 * half * k, 0, 0, size, size)
  ctx.strokeStyle = '#fff'
  ctx.lineWidth = 1
  ctx.beginPath()
  ctx.moveTo(size / 2, 0)
  ctx.lineTo(size / 2, size)
  ctx.moveTo(0, size / 2)
  ctx.lineTo(size, size / 2)
  ctx.stroke()
}

function onPointerMove(ev: PointerEvent) {
  const sp = local(ev)
  if (pointers.has(ev.pointerId)) pointers.set(ev.pointerId, sp)
  if (pinch && pointers.size === 2) {
    const [a, b] = [...pointers.values()]
    const mid = {x: (a.x + b.x) / 2, y: (a.y + b.y) / 2}
    const z = zoomAt(pinch.v, (pinch.s * dist(a, b)) / Math.max(1, pinch.d), pinch.mid)
    view.value = {...z, x: z.x + mid.x - pinch.mid.x, y: z.y + mid.y - pinch.mid.y}
    return
  }
  if (!drag) return
  const p = screenToImage(view.value, sp)
  const e = clone(edit.value)
  const {w, h} = orientedSize.value
  switch (drag.kind) {
    case 'pan':
      view.value = {...view.value, x: drag.vx + sp.x - drag.sx, y: drag.vy + sp.y - drag.sy}
      return
    case 'fixo-move':
      e.fixo = clampRect({...drag.rect, x: drag.rect.x + p.x - drag.start.x, y: drag.rect.y + p.y - drag.start.y}, w, h)
      change(e)
      break
    case 'fixo-resize': {
      const q = {x: Math.max(0, Math.min(w, p.x)), y: Math.max(0, Math.min(h, p.y))}
      let r = resizeRect(e.fixo!, drag.handle, q, fixoRatio(e))
      if (r.x < 0 || r.y < 0 || r.x + r.width > w || r.y + r.height > h) r = e.fixo!
      e.fixo = r
      change(e)
      updateLupe(sp, q)
      break
    }
    case 'persp':
      e.perspectiva![drag.i] = {x: Math.max(0, Math.min(w, p.x)), y: Math.max(0, Math.min(h, p.y))}
      change(e)
      updateLupe(sp, p)
      break
    case 'livre-point':
      e.livre!.points[drag.i] = {x: Math.max(0, Math.min(w, p.x)), y: Math.max(0, Math.min(h, p.y))}
      change(e)
      updateLupe(sp, p)
      break
    case 'ret-move':
      e.retangulo = {...drag.r, cx: drag.r.cx + p.x - drag.start.x, cy: drag.r.cy + p.y - drag.start.y}
      change(e)
      break
    case 'ret-resize':
      e.retangulo = resizeTurned(e.retangulo!, drag.handle, p)
      change(e)
      updateLupe(sp, p)
      break
    case 'ret-rotate': {
      const a = Math.atan2(sp.y - drag.c.y, sp.x - drag.c.x)
      let deg = normDeg(drag.r0 + ((a - drag.a0) * 180) / Math.PI)
      deg = ev.shiftKey ? stepRotation(deg) : snapAngle(deg)
      e.retangulo!.angle = Math.round(deg * 10) / 10
      change(e)
      break
    }
    case 'livre-rotate': {
      const a = Math.atan2(sp.y - drag.c.y, sp.x - drag.c.x)
      let deg = normDeg(drag.r0 + ((a - drag.a0) * 180) / Math.PI)
      deg = ev.shiftKey ? stepRotation(deg) : snapRotation(e.livre!.points, deg)
      e.livre!.rotation = Math.round(deg * 10) / 10
      change(e)
      break
    }
    case 'ann-create': {
      const d = clone(draft.value!)
      if (d.type === 'pen' || d.type === 'highlighter') d.points.push(p)
      else d.points[d.points.length - 1] = p
      draft.value = d
      draw()
      break
    }
    case 'ann-move': {
      const moved = moveAnn(drag.a, p.x - drag.start.x, p.y - drag.start.y)
      e.annotations = e.annotations.map((a) => (a.id === moved.id ? moved : a))
      change(e)
      break
    }
    case 'ann-resize': {
      const box = resizeRect(drag.box, drag.handle, p, 0, 2)
      const r = resizeAnn(drag.a, drag.box, box)
      e.annotations = e.annotations.map((a) => (a.id === r.id ? r : a))
      change(e)
      break
    }
  }
}

function finishDraft() {
  const d = draft.value
  draft.value = null
  if (!d) return
  // a click with no drag makes nothing (but a text, a number)
  const b = annBounds(d)
  if (!(d.type === 'text' || d.type === 'number') && b.width < 2 && b.height < 2 && d.points.length < 3) {
    draw()
    return
  }
  const e = clone(edit.value)
  e.annotations.push(d)
  commit(e)
  selectedAnn.value = d.id
  if (d.type === 'text') nextTick(() => textInput.value?.focus?.())
}

function onPointerUp(ev: PointerEvent) {
  pointers.delete(ev.pointerId)
  if (pinch) {
    if (pointers.size < 2) pinch = null
    return
  }
  lupeAt.value = null
  if (!drag) return
  const d = drag
  drag = null
  if (d.kind === 'pan') return
  if (d.kind === 'ann-create') {
    finishDraft()
    return
  }
  history.end()
  syncCounts()
}

function onDblClick(ev: MouseEvent) {
  if (tab.value !== 'recortar' || edit.value.mode !== 'livre' || !edit.value.livre) return
  const p = screenToImage(view.value, local(ev))
  // the click before it inserted a point there (one step, undone): the
  // point under the cursor is the one to remove
  const i = nearestPoint(edit.value.livre.points, p, hitR())
  if (i < 0) return
  selectedPoint.value = i
  removeSelectedPoint()
}

function onWheel(ev: WheelEvent) {
  ev.preventDefault()
  const sp = local(ev)
  if (ev.ctrlKey || ev.metaKey) {
    const f = Math.exp(-ev.deltaY * 0.0025)
    view.value = zoomAt(view.value, view.value.scale * f, sp)
  } else {
    view.value = {...view.value, x: view.value.x - ev.deltaX, y: view.value.y - ev.deltaY}
  }
}

// --- the keyboard ---

function typing(ev: KeyboardEvent) {
  const t = ev.target as HTMLElement
  return t && (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable)
}

function nudge(dx: number, dy: number) {
  const e = clone(edit.value)
  if (tab.value === 'anotar' && selected.value) {
    e.annotations = e.annotations.map((a) => (a.id === selectedAnn.value ? moveAnn(a, dx, dy) : a))
  } else if (tab.value === 'recortar' && selectedPoint.value >= 0) {
    const arr = e.mode === 'livre' ? e.livre?.points : e.mode === 'perspectiva' ? e.perspectiva : null
    if (!arr) return
    const p = arr[selectedPoint.value]
    arr[selectedPoint.value] = {x: p.x + dx, y: p.y + dy}
  } else if (tab.value === 'recortar' && e.mode === 'retangulo' && e.retangulo) {
    e.retangulo = {...e.retangulo, cx: e.retangulo.cx + dx, cy: e.retangulo.cy + dy}
  } else if (tab.value === 'recortar' && e.mode === 'fixo' && e.fixo) {
    e.fixo = clampRect({...e.fixo, x: e.fixo.x + dx, y: e.fixo.y + dy}, orientedSize.value.w, orientedSize.value.h)
  } else return
  commit(e)
}

function onKey(ev: KeyboardEvent) {
  if (typing(ev)) return
  const k = ev.key.toLowerCase()
  if ((ev.ctrlKey || ev.metaKey) && k === 'z' && !ev.shiftKey) {
    ev.preventDefault()
    undo()
  } else if ((ev.ctrlKey || ev.metaKey) && (k === 'y' || (k === 'z' && ev.shiftKey))) {
    ev.preventDefault()
    redo()
  } else if ((ev.ctrlKey || ev.metaKey) && k === 'd' && tab.value === 'anotar') {
    ev.preventDefault()
    duplicateSelected()
  } else if (k === 'delete' || k === 'backspace') {
    if (tab.value === 'anotar' && selectedAnn.value) {
      ev.preventDefault()
      deleteAnn()
    } else if (tab.value === 'recortar' && edit.value.mode === 'livre' && selectedPoint.value >= 0) {
      ev.preventDefault()
      removeSelectedPoint()
    }
  } else if (k === ' ') {
    spaceDown.value = true
    ev.preventDefault()
  } else if (k.startsWith('arrow')) {
    const step = ev.shiftKey ? 10 : 1
    const d = {arrowleft: [-step, 0], arrowright: [step, 0], arrowup: [0, -step], arrowdown: [0, step]}[k]
    if (d) {
      ev.preventDefault()
      nudge(d[0], d[1])
    }
  } else if (k === 'escape') {
    ev.preventDefault()
    askCancel()
  }
}

function onKeyUp(ev: KeyboardEvent) {
  if (ev.key === ' ') spaceDown.value = false
}

// --- the preview (the side panel of Crop) ---

const preview = ref<HTMLCanvasElement | null>(null)
const resultSize = computed(() => {
  const e = edit.value
  if (e.mode === 'fixo' && e.fixo) return {w: Math.round(e.fixo.width), h: Math.round(e.fixo.height)}
  return cutSize.value
})
const outType = computed(() => resultType(edit.value, sourceType.value))

function paintPreview() {
  const c = preview.value
  if (!c || !orientedPrev.value || tab.value !== 'recortar') return
  const e = edit.value
  let img: HTMLCanvasElement
  if (e.mode === 'fixo' && e.fixo) {
    const r = e.fixo
    const k = kPrev.value
    img = newCanvas(r.width * k, r.height * k)
    img.getContext('2d')!.drawImage(orientedPrev.value, r.x * k, r.y * k, r.width * k, r.height * k, 0, 0, img.width, img.height)
    img = annotated(colored(img, e), {...e, annotations: []})
  } else {
    img = compare.value ? orientedPrev.value : annotated(colorPrev.value || orientedPrev.value, e, (colorPrev.value?.width || 1) / Math.max(1, cutSize.value.w))
  }
  const max = 268
  const k = Math.min(max / img.width, 200 / img.height, 1)
  c.width = Math.max(1, Math.round(img.width * k))
  c.height = Math.max(1, Math.round(img.height * k))
  const ctx = c.getContext('2d')!
  ctx.clearRect(0, 0, c.width, c.height)
  ctx.drawImage(img, 0, 0, c.width, c.height)
}

watch([colorPrev, () => edit.value.fixo, tab, compare], () => nextTick(paintPreview), {deep: true})

// --- save, cancel ---

const confirmClose = ref(false)

function askCancel() {
  if (history.undoCount > 0) confirmClose.value = true
  else emit('cancel')
}

async function save() {
  if (!source.value || saving.value) return
  saving.value = true
  try {
    await nextTick()
    await new Promise((r) => setTimeout(r, 30))
    // a perspective on the image's corners cuts nothing
    const e = normalized(clone(edit.value), orientedSize.value.w, orientedSize.value.h)
    const cropOption =
      e.mode === 'fixo' && e.fixo && FIXO_IS_SIZE
        ? {x: Math.round(e.fixo.x), y: Math.round(e.fixo.y), width: Math.round(e.fixo.width), height: Math.round(e.fixo.height)}
        : null
    let blob: Blob | null = null
    const type = outType.value
    if (changes(e, FIXO_IS_SIZE)) {
      const out = render(source.value, e, FIXO_IS_SIZE)
      blob = await toBlob(out, type)
    }
    const name = resultName(props.name || (props.src.split('?')[0].split('/').pop() ?? 'image'), type)
    emit('save', {blob, file: blob ? new File([blob], name, {type}) : null, name, type, edit: e, cropOption})
  } finally {
    saving.value = false
  }
}

// --- the layout ---

let ro: ResizeObserver | null = null
function measure() {
  const el = stage.value
  if (!el) return
  stageSize.w = el.clientWidth
  stageSize.h = el.clientHeight
  draw()
}

onMounted(() => {
  window.addEventListener('keydown', onKey)
  window.addEventListener('keyup', onKeyUp)
  if (typeof ResizeObserver !== 'undefined' && stage.value) {
    ro = new ResizeObserver(() => measure())
    ro.observe(stage.value)
  }
  measure()
  load().then(() => {
    if (!edit.value.fixo && edit.value.mode === 'fixo') setInitial()
    else if (edit.value.mode === 'perspectiva' && !edit.value.perspectiva) setInitial()
    else if (edit.value.mode === 'livre' && !edit.value.livre) setInitial()
    else if (edit.value.mode === 'retangulo' && !edit.value.retangulo) setInitial()
    draw()
    paintPreview()
  })
})

// the shape of the mode at the start: no step (the history begins there)
function setInitial() {
  const e = clone(edit.value)
  const {w, h} = orientedSize.value
  if (e.mode === 'fixo') e.fixo = fitAspect(w, h, fixoRatio(e))
  if (e.mode === 'perspectiva') e.perspectiva = defaultQuad(w, h)
  if (e.mode === 'livre') e.livre = {points: [], rotation: 0}
  if (e.mode === 'retangulo') e.retangulo = defaultTurned(w, h)
  edit.value = e
  history.reset(e)
  syncCounts()
}

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKey)
  window.removeEventListener('keyup', onKeyUp)
  ro?.disconnect()
})

const cursor = computed(() => {
  if (spaceDown.value || tab.value === 'cores') return 'grab'
  if (tab.value === 'anotar' && tool.value !== 'select') return 'crosshair'
  if (tab.value === 'recortar' && edit.value.mode === 'livre') return 'crosshair'
  return 'default'
})

const MODES: {mode: CropMode; label: keyof Messages; icon: string}[] = [
  {mode: 'fixo', label: 'modeFixed', icon: 'mdi-crop'},
  {mode: 'retangulo', label: 'modeRect', icon: 'mdi-crop-rotate'},
  {mode: 'perspectiva', label: 'modePerspective', icon: 'mdi-perspective-less'},
  {mode: 'livre', label: 'modeFree', icon: 'mdi-vector-polygon'},
]

defineExpose({undo, redo, save, askCancel, edit, history, saving})
</script>

<template>
  <div class="vx-image-editor" :class="{'vx-ie--mobile': mobile}" data-image-editor>
    <!-- the bar: undo/redo, zoom, the mode's actions, cancel/save -->
    <div class="vx-ie-bar">
      <!-- the tools: wrapped on a wide screen, scrolled on a phone -->
      <div class="vx-ie-tools-row">
      <v-btn :disabled="!counts.undo" icon variant="text" size="small" :title="m.undo + ' (Ctrl+Z)'" @click="undo" data-ie-undo>
        <v-badge v-if="counts.undo" :content="counts.undo" color="grey" floating><v-icon>mdi-undo</v-icon></v-badge>
        <v-icon v-else>mdi-undo</v-icon>
      </v-btn>
      <v-btn :disabled="!counts.redo" icon variant="text" size="small" :title="m.redo + ' (Ctrl+Y)'" @click="redo" data-ie-redo>
        <v-badge v-if="counts.redo" :content="counts.redo" color="grey" floating><v-icon>mdi-redo</v-icon></v-badge>
        <v-icon v-else>mdi-redo</v-icon>
      </v-btn>
      <v-divider vertical class="mx-1" />
      <v-btn icon="mdi-magnify-minus-outline" variant="text" size="small" :title="m.zoomOut" @click="zoomStep(-1)" />
      <input v-model="zoomText" class="vx-ie-zoom" @change="zoomTyped" @keydown.enter="zoomTyped" data-ie-zoom />
      <span class="text-caption mr-1">%</span>
      <v-btn icon="mdi-magnify-plus-outline" variant="text" size="small" :title="m.zoomIn" @click="zoomStep(1)" />
      <v-btn variant="text" size="small" prepend-icon="mdi-fit-to-screen-outline" :title="m.zoomFit" @click="doFit"><span class="vx-ie-label">{{ m.zoomFit }}</span></v-btn>
      <v-btn variant="text" size="small" prepend-icon="mdi-magnify-scan" :title="m.zoomActual" @click="doActual"><span class="vx-ie-label">{{ m.zoomActual }}</span></v-btn>
      <v-divider vertical class="mx-1" />
      <template v-if="tab === 'recortar'">
        <v-btn icon="mdi-rotate-left" variant="text" size="small" :title="m.rotateLeft" @click="turn(-90)" />
        <v-btn icon="mdi-rotate-right" variant="text" size="small" :title="m.rotateRight" @click="turn(90)" />
        <v-btn icon="mdi-flip-horizontal" variant="text" size="small" :title="m.flip" @click="mirror" />
        <template v-if="edit.mode === 'livre'">
          <v-btn variant="text" size="small" prepend-icon="mdi-vector-point-minus" :disabled="selectedPoint < 0 || (edit.livre?.points.length || 0) <= 3" @click="removeSelectedPoint" data-ie-remove-point :title="m.removePoint"><span class="vx-ie-label">{{ m.removePoint }}</span></v-btn>
          <v-btn variant="text" size="small" prepend-icon="mdi-format-align-bottom" :disabled="selectedEdge < 0" @click="doAlignEdge" data-ie-align-edge :title="m.alignEdge"><span class="vx-ie-label">{{ m.alignEdge }}</span></v-btn>
          <span class="d-inline-flex align-center text-no-wrap">
            <span class="text-caption ml-2 mr-1">{{ m.angle }}</span>
            <input class="vx-ie-zoom" :value="+(edit.livre?.rotation ?? 0).toFixed(1)" @change="setAngle(($event.target as HTMLInputElement).value)" data-ie-angle />
            <span class="text-caption">°</span>
          </span>
        </template>
        <template v-if="edit.mode === 'retangulo'">
          <span class="d-inline-flex align-center text-no-wrap">
            <span class="text-caption ml-2 mr-1">{{ m.angle }}</span>
            <input class="vx-ie-zoom" :value="+(edit.retangulo?.angle ?? 0).toFixed(1)" @change="setAngle(($event.target as HTMLInputElement).value)" data-ie-angle />
            <span class="text-caption">°</span>
          </span>
        </template>
        <v-btn v-if="edit.mode === 'perspectiva'" variant="text" size="small" prepend-icon="mdi-auto-fix" @click="detectEdges" data-ie-detect :title="m.detectEdges"><span class="vx-ie-label">{{ m.detectEdges }}</span></v-btn>
        <v-btn variant="text" size="small" prepend-icon="mdi-restore" @click="resetMode" :title="m.reset"><span class="vx-ie-label">{{ m.reset }}</span></v-btn>
      </template>
      <template v-if="tab === 'anotar' && selected">
        <v-btn variant="text" size="small" prepend-icon="mdi-content-copy" @click="duplicateSelected" :title="m.duplicate"><span class="vx-ie-label">{{ m.duplicate }}</span></v-btn>
        <v-btn variant="text" size="small" prepend-icon="mdi-delete-outline" color="error" @click="deleteAnn()" :title="m.delete"><span class="vx-ie-label">{{ m.delete }}</span></v-btn>
      </template>
      </div>
      <!-- the closing buttons: together, at the right, out of the scroll -->
      <div v-if="!actionsOutside" class="vx-ie-end">
      <span v-if="queue && queue.total > 1" class="text-caption mr-2" data-ie-queue>{{ queue.index }} {{ m.of }} {{ queue.total }}</span>
      <v-btn variant="text" size="small" prepend-icon="mdi-close" :title="m.cancel" @click="askCancel" data-ie-cancel><span class="vx-ie-label">{{ m.cancel }}</span></v-btn>
      <template v-if="upload">
        <v-btn v-if="camera" variant="text" size="small" prepend-icon="mdi-camera-retake-outline" :title="m.retake" @click="emit('retake')"><span class="vx-ie-label">{{ m.retake }}</span></v-btn>
        <v-btn v-if="queue && queue.total > 1" variant="text" size="small" prepend-icon="mdi-upload-multiple" :title="m.uploadAllWithoutEditing" @click="emit('skipAll')"><span class="vx-ie-label">{{ m.uploadAllWithoutEditing }}</span></v-btn>
        <v-btn variant="tonal" size="small" class="mr-1" prepend-icon="mdi-upload-outline" :title="m.uploadWithoutEditing" @click="emit('skip')" data-ie-skip><span class="vx-ie-label">{{ m.uploadWithoutEditing }}</span></v-btn>
        <v-btn variant="flat" color="primary" size="small" prepend-icon="mdi-check" :title="m.saveAndUpload" :loading="saving" @click="save" data-ie-save><span class="vx-ie-label">{{ m.saveAndUpload }}</span></v-btn>
      </template>
      <v-btn v-else variant="flat" color="primary" size="small" prepend-icon="mdi-check" :title="m.save" :loading="saving" @click="save" data-ie-save><span class="vx-ie-label">{{ m.save }}</span></v-btn>
      </div>
    </div>

    <div class="d-flex align-center">
      <v-tabs v-model="tab" density="compact" color="primary" class="vx-ie-tabs">
        <v-tab value="recortar" prepend-icon="mdi-crop" data-ie-tab="recortar">{{ m.tabCrop }}</v-tab>
        <v-tab value="cores" prepend-icon="mdi-palette-outline" data-ie-tab="cores">{{ m.tabColors }}</v-tab>
        <v-tab value="anotar" prepend-icon="mdi-draw" data-ie-tab="anotar">{{ m.tabAnnotate }}</v-tab>
      </v-tabs>
      <v-spacer />
      <v-btn
        v-if="mobile"
        :icon="showSide ? 'mdi-chevron-down' : 'mdi-tune-variant'"
        variant="text"
        size="small"
        :color="showSide ? 'primary' : undefined"
        data-ie-side-toggle
        @click="showSide = !showSide"
      />
    </div>

    <div class="vx-ie-body">
      <div
        ref="stage"
        class="vx-ie-stage"
        :style="{cursor}"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerUp"
        @pointercancel="onPointerUp"
        @dblclick="onDblClick"
        @wheel="onWheel"
        @contextmenu.prevent
        data-ie-stage
      >
        <canvas ref="cv" class="vx-ie-canvas" />
        <svg class="vx-ie-marks" :width="stageSize.w" :height="stageSize.h">
          <path v-if="maskPath" :d="maskPath" fill="rgba(0,0,0,.5)" fill-rule="evenodd" />
          <!-- fixed -->
          <g v-if="fixoMarks">
            <path :d="fixoMarks.guides" class="vx-ie-guide" />
            <path :d="fixoMarks.outline" class="vx-ie-line" />
            <rect v-for="hd in fixoMarks.handles" :key="hd.h" :x="hd.x - 5" :y="hd.y - 5" width="10" height="10" class="vx-ie-handle" :data-ie-handle="hd.h" />
          </g>
          <!-- perspective -->
          <g v-if="perspMarks">
            <path :d="perspMarks.grid" class="vx-ie-guide" />
            <path :d="perspMarks.outline" class="vx-ie-line" />
            <circle v-for="c in perspMarks.corners" :key="c.i" :cx="c.x" :cy="c.y" r="7" class="vx-ie-handle" :class="{selected: selectedPoint === c.i}" :data-ie-corner="c.i" />
          </g>
          <!-- the free rectangle -->
          <g v-if="retMarks">
            <path :d="retMarks.guides" class="vx-ie-guide" />
            <path :d="retMarks.outline" class="vx-ie-line" />
            <line :x1="retMarks.rot.cx" :y1="retMarks.rot.cy" :x2="retMarks.rot.x" :y2="retMarks.rot.y" class="vx-ie-line" />
            <circle :cx="retMarks.rot.x" :cy="retMarks.rot.y" r="7" class="vx-ie-handle rot" data-ie-rotate />
            <rect v-for="hd in retMarks.handles" :key="hd.h" :x="hd.x - 5" :y="hd.y - 5" width="10" height="10" class="vx-ie-handle" :data-ie-handle="hd.h" />
          </g>
          <!-- free -->
          <g v-if="livreMarks">
            <path v-if="livreMarks.rot && livreMarks.rot.angle" :d="livreMarks.rot.area" class="vx-ie-guide" />
            <path v-if="livreMarks.outline" :d="livreMarks.outline" class="vx-ie-line" />
            <line v-if="selectedEdge >= 0 && livreMarks.edges[selectedEdge]" :x1="livreMarks.edges[selectedEdge].a.x" :y1="livreMarks.edges[selectedEdge].a.y" :x2="livreMarks.edges[selectedEdge].b.x" :y2="livreMarks.edges[selectedEdge].b.y" class="vx-ie-edge" data-ie-edge-selected />
            <template v-if="livreMarks.rot">
              <line :x1="livreMarks.rot.cx" :y1="livreMarks.rot.cy" :x2="livreMarks.rot.x" :y2="livreMarks.rot.y" class="vx-ie-line" />
              <circle :cx="livreMarks.rot.x" :cy="livreMarks.rot.y" r="7" class="vx-ie-handle rot" data-ie-rotate />
            </template>
            <circle v-for="c in livreMarks.points" :key="c.i" :cx="c.x" :cy="c.y" r="6" class="vx-ie-handle" :class="{selected: selectedPoint === c.i}" :data-ie-point="c.i" />
          </g>
          <!-- the selected annotation -->
          <g v-if="annMarks">
            <path :d="annMarks.outline" class="vx-ie-guide sel" />
            <rect v-for="hd in annMarks.handles" :key="hd.h" :x="hd.x - 4" :y="hd.y - 4" width="8" height="8" class="vx-ie-handle" />
          </g>
        </svg>
        <canvas v-show="lupeAt" ref="lupe" width="120" height="120" class="vx-ie-lupe" :style="lupeAt ? {left: Math.min(stageSize.w - 130, Math.max(10, lupeAt.x + 24)) + 'px', top: Math.max(10, lupeAt.y - 144) + 'px'} : {}" />
        <div v-if="loading || saving" class="vx-ie-busy">
          <v-progress-circular indeterminate color="primary" />
          <span class="ml-3">{{ loading ? m.loading : m.saving }}</span>
        </div>
        <div v-if="loadError" class="vx-ie-busy text-error">{{ loadError }}</div>
      </div>

      <!-- the side panel -->
      <div v-show="!mobile || showSide" class="vx-ie-side" data-ie-side>
        <template v-if="tab === 'recortar'">
          <div class="vx-ie-chips">
            <v-chip
              v-for="md in MODES"
              :key="md.mode"
              :prepend-icon="md.icon"
              :color="edit.mode === md.mode ? 'primary' : undefined"
              :variant="edit.mode === md.mode ? 'flat' : 'outlined'"
              size="small"
              class="mr-1 mb-1"
              :data-ie-mode="md.mode"
              @click="setMode(md.mode)"
              >{{ m[md.label] }}</v-chip
            >
          </div>
          <template v-if="edit.mode === 'fixo'">
            <div class="text-caption text-medium-emphasis mt-2">{{ m.ratio }}</div>
            <div v-if="aspect > 0" class="text-caption" data-ie-ratio-size>{{ m.ratioOfSize }}</div>
            <div v-else class="vx-ie-chips d-flex flex-wrap align-center">
              <v-chip
                v-for="r in RATIOS"
                :key="r.value"
                size="small"
                class="mr-1 mb-1"
                :color="(r.value === ratioNow || sameRatio(r.value, ratioNow)) ? 'primary' : undefined"
                :variant="(r.value === ratioNow || sameRatio(r.value, ratioNow)) ? 'flat' : 'outlined'"
                :data-ie-ratio="r.label || (r.value === 0 ? 'free' : 'original')"
                @click="setRatio(r.value)"
                >{{ r.label || (r.value === 0 ? m.ratioFree : m.ratioOriginal) }}</v-chip
              >
              <v-btn
                icon="mdi-phone-rotate-landscape"
                size="x-small"
                variant="text"
                :disabled="!(ratioNow > 0)"
                :title="m.swapRatio"
                data-ie-ratio-swap
                @click="swapRatio"
              />
            </div>
          </template>
          <div class="text-caption text-medium-emphasis mt-2">{{ m.preview }}</div>
          <div class="vx-ie-preview" @pointerdown="compare = true" @pointerup="compare = false" @pointerleave="compare = false">
            <canvas ref="preview" data-ie-preview />
          </div>
          <div class="text-caption">
            {{ m.size }}: <b data-ie-result-size>{{ resultSize.w }} × {{ resultSize.h }}</b> px · {{ m.format }}: <b>{{ outType.replace('image/', '').toUpperCase() }}</b>
          </div>
          <v-btn variant="text" size="small" class="mt-2 px-0" prepend-icon="mdi-image-refresh-outline" @click="backToOriginal">{{ m.backToOriginal }}</v-btn>
        </template>

        <template v-if="tab === 'cores'">
          <v-btn block color="primary" variant="tonal" prepend-icon="mdi-white-balance-sunny" @click="doFixLighting" data-ie-fix-lighting>{{ m.fixLighting }}</v-btn>
          <v-switch v-model="whiteBalance" :label="m.whiteBalance" density="compact" hide-details color="primary" />
          <div class="vx-ie-chips mb-2">
            <v-chip size="small" class="mr-1" @click="applyPreset('auto')" data-ie-preset="auto">{{ m.presetAuto }}</v-chip>
            <v-chip size="small" class="mr-1" @click="applyPreset('document')" data-ie-preset="document">{{ m.presetDocument }}</v-chip>
            <v-chip size="small" @click="applyPreset('original')" data-ie-preset="original">{{ m.presetOriginal }}</v-chip>
          </div>
          <div class="d-flex align-center">
            <span class="text-caption text-medium-emphasis">{{ m.histogram }}</span>
            <v-spacer />
            <v-btn size="x-small" variant="tonal" @pointerdown="compare = true" @pointerup="compare = false" @pointerleave="compare = false">{{ compare ? m.before : m.after }}</v-btn>
          </div>
          <svg v-if="histPaths" viewBox="0 0 256 80" preserveAspectRatio="none" class="vx-ie-hist" data-ie-histogram>
            <path :d="histPaths.r" fill="rgba(229,57,53,.35)" />
            <path :d="histPaths.g" fill="rgba(67,160,71,.35)" />
            <path :d="histPaths.b" fill="rgba(30,136,229,.35)" />
            <path :d="histPaths.luma" fill="none" stroke="currentColor" stroke-width="1" />
          </svg>
          <div v-for="sl in SLIDERS" :key="sl.key" class="vx-ie-slider">
            <div class="d-flex text-caption">
              <span>{{ m[sl.label] }}</span>
              <v-spacer />
              <span>{{ (edit.colors[sl.key] as number).toFixed(sl.step < 1 ? 2 : 0) }}</span>
            </div>
            <v-slider
              :model-value="edit.colors[sl.key] as number"
              :min="sl.min"
              :max="sl.max"
              :step="sl.step"
              density="compact"
              hide-details
              color="primary"
              :data-ie-slider="sl.key"
              @update:model-value="setColors({[sl.key]: $event} as any, true)"
              @end="endColorDrag"
              @dblclick="setColors({[sl.key]: (neutral() as any)[sl.key]} as any)"
            />
          </div>
          <v-switch :model-value="edit.colors.grayscale" :label="m.grayscale" density="compact" hide-details color="primary" @update:model-value="setColors({grayscale: !!$event})" />
          <v-switch :model-value="edit.colors.invert" :label="m.invert" density="compact" hide-details color="primary" @update:model-value="setColors({invert: !!$event})" />
        </template>

        <template v-if="tab === 'anotar'">
          <div class="vx-ie-tools">
            <v-btn
              v-for="t in TOOLS"
              :key="t"
              :icon="TOOL_ICONS[t]"
              :title="m[TOOL_LABELS[t]]"
              size="small"
              :variant="tool === t ? 'flat' : 'text'"
              :color="tool === t ? 'primary' : undefined"
              :data-ie-tool="t"
              @click="tool = t"
            />
          </div>
          <div class="text-caption text-medium-emphasis mt-2">{{ m.properties }}</div>
          <div class="vx-ie-props">
            <label>{{ m.stroke }} <input type="color" :value="props_.stroke" @change="setProp('stroke', ($event.target as HTMLInputElement).value)" /></label>
            <label>
              {{ m.fill }}
              <input type="color" :value="props_.fill || '#ffffff'" :disabled="!props_.fill" @change="setProp('fill', ($event.target as HTMLInputElement).value)" />
              <input type="checkbox" :checked="!!props_.fill" @change="setProp('fill', ($event.target as HTMLInputElement).checked ? '#ffffff' : '')" />
            </label>
            <label>{{ m.width }} <input type="range" min="1" max="40" :value="props_.width" @change="setProp('width', +($event.target as HTMLInputElement).value)" /></label>
            <label>{{ m.opacity }} <input type="range" min="0.1" max="1" step="0.05" :value="props_.opacity" @change="setProp('opacity', +($event.target as HTMLInputElement).value)" /></label>
            <label><input type="checkbox" :checked="props_.shadow" @change="setProp('shadow', ($event.target as HTMLInputElement).checked)" /> {{ m.shadow }}</label>
            <label>
              {{ m.font }}
              <select :value="props_.font" @change="setProp('font', ($event.target as HTMLSelectElement).value)">
                <option value="sans-serif">Sans</option>
                <option value="serif">Serif</option>
                <option value="monospace">Mono</option>
              </select>
            </label>
            <label>{{ m.fontSize }} <input type="number" min="6" max="400" :value="props_.fontSize" style="width: 64px" @change="setProp('fontSize', +($event.target as HTMLInputElement).value)" /></label>
          </div>
          <v-textarea
            v-if="selected && selected.type === 'text'"
            ref="textInput"
            :model-value="selected.text"
            :label="m.text"
            rows="1"
            auto-grow
            density="compact"
            variant="outlined"
            hide-details
            class="mt-2"
            data-ie-text
            @update:model-value="setProp('text', $event)"
          />
          <div class="text-caption text-medium-emphasis mt-3">{{ m.layers }}</div>
          <div class="vx-ie-layers">
            <div
              v-for="a in [...edit.annotations].reverse()"
              :key="a.id"
              class="vx-ie-layer"
              :class="{selected: a.id === selectedAnn}"
              data-ie-layer
              @click="selectedAnn = a.id; tool = 'select'"
            >
              <v-icon size="small" class="mr-1" :color="a.stroke">{{ TOOL_ICONS[a.type] }}</v-icon>
              <span class="text-caption">{{ m[TOOL_LABELS[a.type]] }}{{ a.type === 'number' ? ' ' + a.n : a.type === 'text' ? ': ' + a.text.slice(0, 12) : '' }}</span>
              <v-spacer />
              <v-btn icon="mdi-arrow-up" size="x-small" variant="text" @click.stop="moveLayer(a.id, 1)" />
              <v-btn icon="mdi-arrow-down" size="x-small" variant="text" @click.stop="moveLayer(a.id, -1)" />
              <v-btn icon="mdi-close" size="x-small" variant="text" @click.stop="deleteAnn(a.id)" />
            </div>
          </div>
        </template>
      </div>
    </div>

    <v-dialog v-model="confirmClose" max-width="380">
      <v-card :title="m.unsavedTitle" :text="m.unsavedText">
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="confirmClose = false">{{ m.keepEditing }}</v-btn>
          <v-btn color="error" variant="flat" data-ie-discard @click="confirmClose = false; emit('cancel')">{{ m.discard }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style>
.vx-image-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 420px;
  user-select: none;
}
.vx-ie-bar {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px 8px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-ie-tools-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 2px;
  flex: 1;
  min-width: 0;
}
.vx-ie--mobile .vx-ie-tools-row {
  flex-wrap: nowrap;
  overflow-x: auto;
  scrollbar-width: none;
}
.vx-ie--mobile .vx-ie-label,
.vx-ie--mobile .vx-ie-end .text-caption {
  display: none;
}
.vx-ie--mobile .vx-ie-bar {
  padding: 2px 4px;
}
.vx-ie--mobile .vx-ie-bar .v-btn {
  min-width: 0;
  padding: 0 6px;
}
.vx-ie--mobile .vx-ie-bar .v-btn :deep(.v-btn__prepend),
.vx-ie--mobile .vx-ie-bar .v-btn .v-btn__prepend {
  margin-inline: 0;
}
/* the closing buttons: apart from the tools that scroll */
.vx-ie--mobile .vx-ie-end {
  padding-left: 4px;
  border-left: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-ie--mobile .vx-ie-body {
  position: relative;
}
.vx-ie--mobile .vx-ie-side {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  width: auto;
  max-height: 55%;
  z-index: 2;
  border-left: none;
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 -4px 16px rgba(0, 0, 0, 0.3);
}
.vx-ie-end {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-left: auto;
  flex: none;
  white-space: nowrap;
}
.vx-ie-zoom {
  width: 48px;
  text-align: right;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 1px 4px;
  font-size: 0.8rem;
  color: inherit;
}
.vx-ie-body {
  display: flex;
  flex: 1;
  min-height: 0;
}
.vx-ie-stage {
  position: relative;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  touch-action: none;
  background: #3a3a3a;
}
.vx-ie-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
.vx-ie-marks {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: visible;
}
.vx-ie-line {
  fill: none;
  stroke: #fff;
  stroke-width: 1.5;
  vector-effect: non-scaling-stroke;
  filter: drop-shadow(0 0 1px rgba(0, 0, 0, 0.8));
}
.vx-ie-guide {
  fill: none;
  stroke: rgba(255, 255, 255, 0.7);
  stroke-width: 1;
  stroke-dasharray: 4 4;
  vector-effect: non-scaling-stroke;
}
.vx-ie-guide.sel {
  stroke: rgb(var(--v-theme-primary));
  stroke-width: 1.5;
}
.vx-ie-edge {
  stroke: rgb(var(--v-theme-primary));
  stroke-width: 4;
  vector-effect: non-scaling-stroke;
}
.vx-ie-handle {
  fill: #fff;
  stroke: rgb(var(--v-theme-primary));
  stroke-width: 2;
  vector-effect: non-scaling-stroke;
}
.vx-ie-handle.selected,
.vx-ie-handle.rot {
  fill: rgb(var(--v-theme-primary));
  stroke: #fff;
}
.vx-ie-lupe {
  position: absolute;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
  pointer-events: none;
  background: #222;
}
.vx-ie-busy {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
}
.vx-ie-side {
  width: 300px;
  flex: none;
  overflow-y: auto;
  padding: 10px 12px;
  border-left: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-ie-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80px;
  margin: 6px 0;
  background: repeating-conic-gradient(#ddd 0% 25%, #fff 0% 50%) 50% / 12px 12px;
  border-radius: 4px;
  cursor: pointer;
}
.vx-ie-hist {
  width: 100%;
  height: 80px;
  display: block;
  margin: 4px 0 8px;
  background: rgba(var(--v-theme-on-surface), 0.04);
}
.vx-ie-slider {
  margin-bottom: 2px;
}
.vx-ie-tools {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
}
.vx-ie-props {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.8rem;
}
.vx-ie-props label {
  display: flex;
  align-items: center;
  gap: 6px;
}
.vx-ie-props select,
.vx-ie-props input[type='number'] {
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  padding: 1px 4px;
  color: inherit;
}
.vx-ie-layers {
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  max-height: 220px;
  overflow-y: auto;
}
.vx-ie-layer {
  display: flex;
  align-items: center;
  padding: 2px 6px;
  cursor: pointer;
}
.vx-ie-layer.selected {
  background: rgba(var(--v-theme-primary), 0.12);
}
</style>
