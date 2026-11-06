import {expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
import CameraControls from '../ImageTools/webcam/CameraControls.vue'
import {constraintFor, controlsOf} from '../ImageTools/webcam/controls'
import {mountTemplate} from '@/lib/__tests__/testutils'

// a C922's capabilities (Chrome): no zoom without the PTZ permission
const caps = {
  brightness: {min: 0, max: 255, step: 1},
  contrast: {min: 0, max: 255, step: 1},
  saturation: {min: 0, max: 255, step: 1},
  sharpness: {min: 0, max: 255, step: 1},
  exposureMode: ['continuous', 'manual'],
  exposureTime: {min: 3, max: 2047, step: 1},
  whiteBalanceMode: ['continuous', 'manual'],
  colorTemperature: {min: 2000, max: 6500, step: 10},
  focusMode: ['continuous', 'manual'],
  focusDistance: {min: 0, max: 250, step: 5},
  width: {min: 1, max: 1920},
  frameRate: {min: 0, max: 60},
}

it('the controls a camera offers, by group; none, none', () => {
  const g = controlsOf(caps)
  expect(g.map((x) => x.key)).toEqual(['exposure', 'whiteBalance', 'focus', 'image'])
  expect(g[0].controls.map((c) => c.key)).toEqual(['exposureMode', 'exposureTime'])
  expect(controlsOf({width: {min: 1, max: 2}})).toEqual([])
  expect(controlsOf(null)).toEqual([])
})

it('a value of a manual control takes its group to manual', () => {
  const [exposure] = controlsOf(caps)
  const time = exposure.controls[1]
  expect(constraintFor(time, 300, {exposureMode: 'continuous'}, caps)).toEqual({exposureTime: 300, exposureMode: 'manual'})
  expect(constraintFor(time, 300, {exposureMode: 'manual'}, caps)).toEqual({exposureTime: 300})
  // the switch: automatic back
  expect(constraintFor(exposure.controls[0], true, {}, caps)).toEqual({exposureMode: 'continuous'})
})

it('the panel sets the camera (applyConstraints) and puts it back (Reset)', async () => {
  let settings: Record<string, any> = {brightness: 128, contrast: 128, saturation: 128, sharpness: 128, exposureMode: 'continuous', exposureTime: 156, whiteBalanceMode: 'continuous', colorTemperature: 4000, focusMode: 'continuous', focusDistance: 0}
  const applied: any[] = []
  const track = {
    getCapabilities: () => caps,
    getSettings: () => ({...settings}),
    applyConstraints: vi.fn(async (c: any) => {
      applied.push(c.advanced[0])
      settings = {...settings, ...c.advanced[0]}
    }),
  }
  const w = mountTemplate(CameraControls, {track})
  await flushPromises()
  expect(w.findAll('[data-camera-group]').map((g) => g.attributes('data-camera-group'))).toEqual(['exposure', 'whiteBalance', 'focus', 'image'])
  const vm = w.vm as any
  // through the component's own setter, as a slider does
  const time = vm.groups[0].controls[1]
  ;(w.vm as any).$?.setupState?.set?.(time, 500)
  await flushPromises()
  expect(applied.at(-1)).toEqual({exposureTime: 500, exposureMode: 'manual'})
  await w.find('[data-camera-reset]').trigger('click')
  await flushPromises()
  expect(applied.at(-1)).toMatchObject({exposureMode: 'continuous', exposureTime: 156, brightness: 128})
  w.unmount()
})

it('a camera with no controls says so', async () => {
  const track = {getCapabilities: () => ({width: {min: 1, max: 2}}), getSettings: () => ({}), applyConstraints: vi.fn()}
  const w = mountTemplate(CameraControls, {track})
  await flushPromises()
  expect(w.find('[data-camera-no-controls]').exists()).toBe(true)
  w.unmount()
})
