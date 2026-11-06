import {expect, it} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import Lightbox from '../Lightbox/Lightbox.vue'
import {imagesOf} from '../Lightbox/lightbox'
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
