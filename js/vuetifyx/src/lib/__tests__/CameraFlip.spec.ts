import {afterEach, expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import CameraControls from '../ImageTools/webcam/CameraControls.vue'
import WebCam from '../ImageTools/webcam/WebCam.vue'
import WebCamUI from '../ImageTools/webcam/WebCamUI.vue'
import {mountTemplate} from '@/lib/__tests__/testutils'

// The image of the camera mirrored, with the image editor's orientation
// ({rotate, flipH}): the buttons of the panel, the preview, the photo, and
// what is remembered by camera.

afterEach(() => {
  vi.restoreAllMocks()
})

// a storage of the page, in memory (the test's environment has none)
function memoryStorage() {
  const m = new Map<string, string>()
  const storage = {
    getItem: (k: string) => (m.has(k) ? m.get(k)! : null),
    setItem: (k: string, v: string) => void m.set(k, String(v)),
    removeItem: (k: string) => void m.delete(k),
    clear: () => m.clear(),
  }
  Object.defineProperty(window, 'localStorage', {value: storage, configurable: true})
  return m
}

it('the panel mirrors as the image is seen, its buttons of every camera, Reset undoing it', async () => {
  let orientation = {rotate: 0, flipH: false}
  const w = mountTemplate(CameraControls, {
    track: null,
    orientation,
    'onUpdate:orientation': (o: any) => {
      orientation = o
      w.setProps({orientation: o})
    },
  })
  // no control of the camera: the mirror still there
  expect(w.find('[data-camera-no-controls]').exists()).toBe(true)
  await w.find('[data-camera-flip-h]').trigger('click')
  expect(orientation).toEqual({rotate: 0, flipH: true})
  await w.find('[data-camera-flip-v]').trigger('click')
  // both: turned 180°, not mirrored
  expect(orientation).toEqual({rotate: 180, flipH: false})
  await w.find('[data-camera-flip-h]').trigger('click')
  // vertically alone
  expect(orientation).toEqual({rotate: 180, flipH: true})
  await w.find('[data-camera-reset]').trigger('click')
  expect(orientation).toEqual({rotate: 0, flipH: false})
  w.unmount()
})

// a fake camera, of a video of 4 × 2
function fakeMedia() {
  const media = {
    getUserMedia: vi.fn(async () => ({getTracks: () => [{stop: vi.fn()}]})),
    enumerateDevices: vi.fn(async () => [{deviceId: 'cam-1', kind: 'videoinput', label: 'USB Camera'}]),
  }
  Object.defineProperty(navigator, 'mediaDevices', {value: media, configurable: true})
}

// canvases that note what is drawn on them
function recordCanvases() {
  const calls: {canvas: HTMLCanvasElement; op: string; args: any[]}[] = []
  vi.spyOn(HTMLCanvasElement.prototype, 'getContext').mockImplementation(function (this: HTMLCanvasElement) {
    const canvas = this
    const note = (op: string) => (...args: any[]) => calls.push({canvas, op, args})
    return {translate: note('translate'), rotate: note('rotate'), scale: note('scale'), drawImage: note('drawImage')} as any
  })
  vi.spyOn(HTMLCanvasElement.prototype, 'toDataURL').mockReturnValue('data:')
  vi.spyOn(HTMLCanvasElement.prototype, 'toBlob').mockImplementation((cb: any) => cb(new Blob()))
  return calls
}

it('the preview and the photo mirrored alike', async () => {
  fakeMedia()
  const w = mountTemplate(WebCam, {orientation: {rotate: 0, flipH: true}, autoStart: false})
  await flushPromises()
  const video = w.find('video').element as HTMLVideoElement
  expect(video.style.transform).toBe('rotate(0deg) scaleX(-1)')
  Object.defineProperty(video, 'videoWidth', {value: 4})
  Object.defineProperty(video, 'videoHeight', {value: 2})

  const calls = recordCanvases()
  await (w.vm as any).takePhoto()
  // the frame oriented (mirrored), then drawn on the photo
  const mirrored = calls.find((c) => c.op === 'scale')
  expect(mirrored?.args).toEqual([-1, 1])
  const drawn = calls.filter((c) => c.op === 'drawImage')
  expect(drawn[0].args[0]).toBe(video)
  expect(drawn.at(-1)!.args[0]).toBe(mirrored!.canvas)
  expect(w.emitted('photoTaken')).toBeTruthy()

  // not mirrored: the video itself, no transform
  await w.setProps({orientation: {rotate: 0, flipH: false}})
  expect(video.style.transform).toBe('')
  calls.length = 0
  await (w.vm as any).takePhoto()
  expect(calls.some((c) => c.op === 'scale')).toBe(false)
  expect(calls.filter((c) => c.op === 'drawImage')[0].args[0]).toBe(video)
  w.unmount()
})

it('the mirror is remembered by camera', async () => {
  const stored = memoryStorage()
  window.localStorage.setItem('_vwl_orientation:cam-1', JSON.stringify({rotate: 180, flipH: true}))
  fakeMedia()
  const w = mountTemplate(WebCamUI, {})
  for (let i = 0; i < 6; i++) await flushPromises()
  expect(w.findComponent(WebCam).props('orientation')).toEqual({rotate: 180, flipH: true})
  w.unmount()
  // kept as it was
  expect(stored.get('_vwl_orientation:cam-1')).toBe(JSON.stringify({rotate: 180, flipH: true}))
})
