import {expect, it} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import Lightbox from '../Lightbox/Lightbox.vue'
import {imagesOf, lightboxOf} from '../Lightbox/lightbox'
import {mountTemplate} from '@/lib/__tests__/testutils'

it('the images of a group, in the order of the page; alone without one', () => {
  const root = document.createElement('div')
  root.innerHTML = `
    <div data-lightbox="a" data-lightbox-src="/1.jpg" data-lightbox-caption="um"></div>
    <div data-lightbox="b" data-lightbox-src="/x.jpg"></div>
    <div data-lightbox="a" data-lightbox-src="/2.jpg"></div>
    <div data-lightbox="a" data-lightbox-src="/3.jpg"></div>
    <div data-lightbox-src="/solo.jpg"></div>`
  document.body.appendChild(root)
  const els = root.querySelectorAll('[data-lightbox-src]')
  const {images, index} = imagesOf(els[2])
  expect(images.map((i) => i.src)).toEqual(['/1.jpg', '/2.jpg', '/3.jpg'])
  expect(index).toBe(1)
  expect(images[0].caption).toBe('um')
  expect(imagesOf(els[4]).images.map((i) => i.src)).toEqual(['/solo.jpg'])
  root.remove()
})

it('goes through the images, round', async () => {
  // the dialog's overlay asks the viewport (none in jsdom)
  ;(globalThis as any).visualViewport ??= Object.assign(new EventTarget(), {width: 1024, height: 768, scale: 1, offsetLeft: 0, offsetTop: 0})
  const images = [{src: '/1.jpg'}, {src: '/2.jpg'}, {src: '/3.jpg'}]
  const w = mountTemplate(Lightbox, {images, index: 2})
  await flushPromises()
  const vm = w.vm as any
  expect(vm.i).toBe(2)
  vm.go(1)
  expect(vm.i).toBe(0)
  vm.go(-1)
  expect(vm.i).toBe(2)
  w.unmount()
})

it('a control over the image keeps its click', () => {
  const root = document.createElement('div')
  root.innerHTML = `
    <div data-lightbox-src="/1.jpg" id="img"><img id="pic"><button id="edit"><i id="icon"></i></button></div>
    <button data-lightbox-src="/2.jpg" id="btn"><span id="inner"></span></button>`
  const $ = (id: string) => root.querySelector('#' + id)
  expect(lightboxOf($('pic'))).toBe($('img'))
  expect(lightboxOf($('icon'))).toBeNull()
  expect(lightboxOf($('edit'))).toBeNull()
  // the element of the image a button itself: it opens
  expect(lightboxOf($('inner'))).toBe($('btn'))
  expect(lightboxOf(root)).toBeNull()
})

it('the size and the name of a file, for its chip', () => {
  const root = document.createElement('div')
  root.innerHTML = `
    <div data-lightbox="f" data-lightbox-src="/a.png" data-lightbox-size="2048" data-lightbox-name="a.png"></div>
    <div data-lightbox="f" data-lightbox-src="/b.png" data-lightbox-size="x"></div>`
  document.body.appendChild(root)
  const {images} = imagesOf(root.querySelector('[data-lightbox-src]')!)
  expect(images[0]).toMatchObject({size: 2048, name: 'a.png'})
  expect(images[1].size).toBeUndefined()
  root.remove()
})

it('the description below the image, hidden by its button; none, no button', async () => {
  ;(globalThis as any).visualViewport ??= Object.assign(new EventTarget(), {width: 1024, height: 768, scale: 1, offsetLeft: 0, offsetTop: 0})
  const w = mountTemplate(Lightbox, {images: [{src: '/1.jpg', description: 'uma descrição'}, {src: '/2.jpg'}]})
  await flushPromises()
  const q = (s: string) => document.querySelector(s)
  expect(q('[data-lightbox-description]')?.textContent).toBe('uma descrição')
  ;(q('[data-lightbox-toggle-description]') as HTMLElement).click()
  await flushPromises()
  expect(q('[data-lightbox-description]')).toBeNull()
  ;(w.vm as any).go(1)
  await flushPromises()
  expect(q('[data-lightbox-toggle-description]')).toBeNull()
  w.unmount()
})
