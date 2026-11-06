// The editor opened from anywhere in the admin: window.rvqImageEditor.open —
// the dialog mounted in the app (its Vuetify), a promise of what was done.
// The thumbnail of a media box ("Edit image") and the upload (before
// sending) open it so.

import {type App, createVNode, render} from 'vue'
import EditorDialog from './EditorDialog.vue'
import type {SaveEvent} from './ImageEditor.vue'
import type {Edit} from '@/core/edit'
import {type Limit, fitSize, scaleRect} from '@/core/fit'

export interface OpenOptions {
  src?: string
  file?: Blob
  name?: string
  type?: string
  edit?: Partial<Edit> | null
  aspect?: number
  tab?: 'recortar' | 'cores' | 'anotar'
  mode?: 'fixo' | 'perspectiva' | 'livre'
  messages?: Record<string, string>
  locale?: string
  title?: string
  upload?: boolean
  queue?: {index: number; total: number}
  camera?: boolean
  maxPixels?: number
}

export interface OpenResult {
  action: 'save' | 'cancel' | 'skip' | 'skipAll' | 'retake'
  save?: SaveEvent
}

let app: App | null = null

export function setApp(a: App) {
  app = a
}

export function open(opts: OpenOptions): Promise<OpenResult> {
  return new Promise((resolve) => {
    const el = document.createElement('div')
    document.body.appendChild(el)
    const {title, ...editor} = opts
    const vnode = createVNode(EditorDialog, {
      editor,
      title,
      onDone: (r: OpenResult) => {
        // the dialog's leave transition, then away
        setTimeout(() => {
          render(null, el)
          el.remove()
        }, 300)
        resolve(r)
      },
    })
    if (app) vnode.appContext = app._context
    render(vnode, el)
  })
}

export interface UploadOptions {
  // the editor opens by itself (MediaBoxConfig.EditOnUpload); else only by
  // the "Edit" of a file of the queue
  aspect?: number
  tab?: OpenOptions['tab']
  mode?: OpenOptions['mode']
  messages?: Record<string, string>
  locale?: string
  camera?: boolean
  maxPixels?: number
}

export interface Uploaded {
  // what is sent: the result (or the file as chosen)
  file: File
  // the original and the edit — when the image was edited
  original?: File
  edit?: Edit
  // the fixed crop (of the field's largest size), in the result's pixels
  cropOption?: {x: number; y: number; width: number; height: number} | null
  // its thumbnail (fitUploads), made here: the server checks it, and makes
  // it only when it is not sent
  thumb?: File
}

export const isImage = (f: Blob): boolean => /^image\/(jpeg|png|webp|gif|bmp)$/.test(f.type)

// editUploads: the images of files, one after the other ("1 of 3"), in the
// editor before they are sent; the others as they are. null: cancelled (or
// "Take another" of a photo: retake).
export async function editUploads(files: File[], o: UploadOptions = {}): Promise<Uploaded[] | null | 'retake'> {
  const out: Uploaded[] = []
  const images = files.filter(isImage)
  let n = 0
  let skipAll = false
  for (const f of files) {
    if (!isImage(f) || skipAll) {
      out.push({file: f})
      continue
    }
    n++
    const r = await open({
      file: f,
      name: f.name,
      type: f.type,
      aspect: o.aspect,
      tab: o.tab,
      mode: o.mode,
      messages: o.messages,
      locale: o.locale,
      upload: true,
      queue: {index: n, total: images.length},
      camera: o.camera,
      maxPixels: o.maxPixels,
    })
    switch (r.action) {
      case 'cancel':
        return null
      case 'retake':
        return 'retake'
      case 'skipAll':
        skipAll = true
        out.push({file: f})
        break
      case 'skip':
        out.push({file: f})
        break
      case 'save': {
        const s = r.save!
        const file = s.file || f
        out.push({file, original: f, edit: s.edit, cropOption: s.cropOption})
        break
      }
    }
  }
  return out
}

// fitFile: the image f made smaller to fit max (its proportions kept, its
// EXIF orientation applied), of its type; f itself when it fits, is not an
// image that can be drawn (a GIF: its frames) or cannot be read. ratio: the
// scale it was made in.
export async function fitFile(f: File, max: Limit | null | undefined): Promise<{file: File; ratio: number}> {
  const same = {file: f, ratio: 1}
  if (!max || !isImage(f) || f.type === 'image/gif' || typeof createImageBitmap !== 'function') return same
  let bmp: ImageBitmap
  try {
    bmp = await createImageBitmap(f, {imageOrientation: 'from-image'} as ImageBitmapOptions)
  } catch {
    return same
  }
  try {
    const s = fitSize(bmp.width, bmp.height, max)
    if (s.ratio === 1) return same
    const canvas = document.createElement('canvas')
    canvas.width = s.width
    canvas.height = s.height
    const g = canvas.getContext('2d')
    if (!g) return same
    g.imageSmoothingEnabled = true
    g.imageSmoothingQuality = 'high'
    g.drawImage(bmp, 0, 0, s.width, s.height)
    // a BMP is sent as a PNG
    const type = f.type === 'image/bmp' ? 'image/png' : f.type
    const blob = await new Promise<Blob | null>((ok) => canvas.toBlob(ok, type, 0.92))
    if (!blob) return same
    const name = type === f.type ? f.name : f.name.replace(/\.[^.]*$/, '') + '.png'
    return {file: new window.File([blob], name, {type: blob.type, lastModified: f.lastModified}), ratio: s.ratio}
  } finally {
    bmp.close()
  }
}

// thumbOf: the thumbnail of the image f (made to fit max), sent with it —
// the server keeps it, not making it —; none when f cannot be drawn here or
// would change its type (a BMP: the server makes it).
async function thumbOf(f: File, max: Limit | null | undefined): Promise<File | undefined> {
  if (!max || !isImage(f) || f.type === 'image/gif' || f.type === 'image/bmp') return undefined
  const {file} = await fitFile(f, max)
  return file
}

// fitUploads: the images of list made smaller to fit max before they are
// sent — after the editor —, their fixed crops at the new scale, and their
// thumbnails (thumbMax) made; the originals kept for the edits (their pixels
// the edit's).
export async function fitUploads(list: Uploaded[], max: Limit | null | undefined, thumbMax?: Limit | null): Promise<Uploaded[]> {
  if (!max && !thumbMax) return list
  return Promise.all(
    list.map(async (u) => {
      const {file, ratio} = await fitFile(u.file, max)
      const out: Uploaded = ratio === 1 ? {...u} : {...u, file, cropOption: u.cropOption ? scaleRect(u.cropOption, ratio) : u.cropOption}
      const thumb = await thumbOf(out.file, thumbMax)
      if (thumb) out.thumb = thumb
      return out
    }),
  )
}

// fitSave: the image the editor saved (s.file) made smaller to fit max, its
// fixed crop at the new scale, and its thumbnail (thumbMax).
export async function fitSave(s: SaveEvent, max: Limit | null | undefined, thumbMax?: Limit | null): Promise<SaveEvent & {thumb?: File}> {
  if (!s.file) return s
  let out: SaveEvent & {thumb?: File} = s
  const {file, ratio} = await fitFile(s.file, max)
  if (ratio !== 1) {
    out = {...s, file, blob: file, type: file.type, name: file.name, cropOption: s.cropOption ? scaleRect(s.cropOption, ratio) : s.cropOption}
  }
  const thumb = await thumbOf(out.file!, thumbMax)
  return thumb ? {...out, thumb} : out
}
