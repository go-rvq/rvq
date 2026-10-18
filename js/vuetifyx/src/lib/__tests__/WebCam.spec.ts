import {expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import WebCam from '../ImageTools/webcam/WebCam.vue'
import WebCamUI from '../ImageTools/webcam/WebCamUI.vue'
import {mountTemplate} from '@/lib/__tests__/testutils'

// a fake camera: getUserMedia's streams, their tracks stopped; the devices
function fakeMedia(opts: {deny?: boolean} = {}) {
  const calls: MediaStreamConstraints[] = []
  const tracks: {stop: ReturnType<typeof vi.fn>}[] = []
  const media = {
    getUserMedia: vi.fn(async (c: MediaStreamConstraints) => {
      calls.push(JSON.parse(JSON.stringify(c)))
      if (opts.deny) throw Object.assign(new Error('denied'), {name: 'NotAllowedError'})
      const track = {stop: vi.fn()}
      tracks.push(track)
      return {getTracks: () => [track]}
    }),
    enumerateDevices: vi.fn(async () => [
      {deviceId: 'cam-1', kind: 'videoinput', label: 'Integrated Camera'},
      {deviceId: 'cam-2', kind: 'videoinput', label: 'USB Camera'},
      {deviceId: 'mic', kind: 'audioinput', label: 'Mic'},
    ]),
  }
  Object.defineProperty(navigator, 'mediaDevices', {value: media, configurable: true})
  return {calls, tracks, media}
}


it('WebCam opens a camera (the usb one preferred), changes it, closes, the props untouched', async () => {
  const {calls, tracks} = fakeMedia()
  const constraints = {video: {width: {ideal: 640}, facingMode: 'environment'}}
  const w = mountTemplate(WebCam, {constraints})
  // the access tested, the devices listed, the device chosen (its watcher), opened
  for (let i = 0; i < 5; i++) await flushPromises()
  const vm = w.vm as any
  // the access tested, the devices listed, the preferred one opened
  expect(vm.cameras.map((c: any) => c.deviceId)).toEqual(['cam-1', 'cam-2'])
  const opened = calls.at(-1) as any
  expect(opened.video.deviceId).toEqual({exact: 'cam-2'})
  expect(opened.video.facingMode).toBeUndefined()
  expect(opened.video.width).toEqual({ideal: 640})
  // the props' constraints were not written
  expect(constraints.video).toEqual({width: {ideal: 640}, facingMode: 'environment'})
  // another camera: the open stream stopped, no error, the other opened
  vm.changeCamera('cam-1')
  await flushPromises()
  vm.changeCamera('cam-1')
  await flushPromises()
  expect(tracks.some((t) => t.stop.mock.calls.length > 0)).toBe(true)
  expect((calls.at(-1) as any).video.deviceId).toEqual({exact: 'cam-1'})
  // closed: stopped
  expect(() => vm.stop()).not.toThrow()
  w.unmount()
})

it('WebCamUI says the permission was denied', async () => {
  fakeMedia({deny: true})
  const w = mountTemplate(WebCamUI, {})
  await flushPromises()
  expect(w.find('[data-webcam-problem]').exists()).toBe(true)
  expect(w.find('[data-webcam-problem]').text()).toMatch(/denied|negada/)
  w.unmount()
})
