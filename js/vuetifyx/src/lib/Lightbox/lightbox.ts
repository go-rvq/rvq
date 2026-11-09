// The lightbox: a click on an element with data-lightbox-src opens that
// image over the page — the elements of the same data-lightbox (a group: the
// images of an item, say) in their order on the page, gone through by the
// arrows, a swipe or the thumbnails. data-lightbox-caption: what is said
// under it. Any page, any component: the click is caught on the document.
//
//   <div data-lightbox="item-1" data-lightbox-src="/a.jpg" data-lightbox-caption="a.jpg">…</div>

import {type App, createVNode, render} from 'vue'
import Lightbox from './Lightbox.vue'

export interface LightboxImage {
  src: string
  caption?: string
  // its size in bytes (data-lightbox-size) and the name it is downloaded
  // as (data-lightbox-name; else the caption): its chip
  size?: number
  name?: string
  // description: what is said below it (data-lightbox-description), hidden
  // by the button of the bar; none, nothing
  description?: string
}

let app: App | null = null
let mounted: HTMLElement | null = null

// imagesOf: the images of the group of el (el alone, when it has none) and
// where el is among them.
export function imagesOf(el: Element, root: ParentNode = document): {images: LightboxImage[]; index: number} {
  const group = el.getAttribute('data-lightbox')
  const els = group
    ? Array.from(root.querySelectorAll('[data-lightbox-src]')).filter((e) => e.getAttribute('data-lightbox') === group)
    : [el]
  const images: LightboxImage[] = []
  let index = 0
  for (const e of els) {
    const src = e.getAttribute('data-lightbox-src')
    if (!src) continue
    // the same image twice in a group (an element repeated): once
    if (images.some((i) => i.src === src) && e !== el) continue
    if (e === el) index = images.length
    const size = Number(e.getAttribute('data-lightbox-size'))
    images.push({
      src,
      caption: e.getAttribute('data-lightbox-caption') || undefined,
      size: size > 0 ? size : undefined,
      name: e.getAttribute('data-lightbox-name') || undefined,
      description: e.getAttribute('data-lightbox-description') || undefined,
    })
  }
  return {images, index}
}

// open shows images from index.
export function open(images: LightboxImage[], index = 0) {
  close()
  if (!images.length) return
  const el = document.createElement('div')
  document.body.appendChild(el)
  mounted = el
  const vnode = createVNode(Lightbox, {images, index, onClose: close})
  if (app) vnode.appContext = app._context
  render(vnode, el)
}

export function close() {
  if (!mounted) return
  render(null, mounted)
  mounted.remove()
  mounted = null
}

// CONTROLS are what, inside an image of the lightbox, keeps its own click.
const CONTROLS = 'button, a[href], input, textarea, select, label, [role="button"], [data-lightbox-ignore]'

// lightboxOf: the image of the lightbox a click on target opens; null, none —
// out of one, or on a control over it (its edit, its delete: the control's).
export function lightboxOf(target: Element | null): Element | null {
  const el = target?.closest?.('[data-lightbox-src]')
  if (!el) return null
  const control = target!.closest(CONTROLS)
  if (control && control !== el && el.contains(control)) return null
  return el
}

function onClick(ev: MouseEvent) {
  if (ev.defaultPrevented || ev.button !== 0 || ev.ctrlKey || ev.metaKey || ev.shiftKey) return
  const el = lightboxOf(ev.target as Element | null)
  if (!el) return
  ev.preventDefault()
  ev.stopPropagation()
  const {images, index} = imagesOf(el)
  open(images, index)
}

// install catches the clicks of the page, once.
export function install(a: App) {
  app = a
  const w = window as any
  if (w.__vxLightbox) return
  w.__vxLightbox = {open, close}
  document.addEventListener('click', onClick, true)
}
