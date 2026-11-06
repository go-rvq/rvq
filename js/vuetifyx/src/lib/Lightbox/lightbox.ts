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
    images.push({src, caption: e.getAttribute('data-lightbox-caption') || undefined})
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

function onClick(ev: MouseEvent) {
  if (ev.defaultPrevented || ev.button !== 0 || ev.ctrlKey || ev.metaKey || ev.shiftKey) return
  const el = (ev.target as Element | null)?.closest?.('[data-lightbox-src]')
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
