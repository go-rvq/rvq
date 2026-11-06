// The editor opened from anywhere in the admin: window.rvqImageEditor.open —
// the dialog mounted in the app (its Vuetify), a promise of what was done.
// The thumbnail of a media box ("Edit image") and the upload (before
// sending) open it so.

import {type App, createVNode, render} from 'vue'
import EditorDialog from './EditorDialog.vue'
import type {SaveEvent} from './ImageEditor.vue'
import type {Edit} from '@/core/edit'

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
