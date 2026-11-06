import {describe, expect, it} from 'vitest'
import {changes, emptyEdit, parseEdit, resultName, resultType} from '../edit'

describe('the edit', () => {
  it('a fixed crop of a size is no new image (its CropOption)', () => {
    const e = emptyEdit('fixo')
    e.fixo = {x: 0, y: 0, width: 10, height: 10}
    expect(changes(e, true)).toBe(false)
    expect(changes(e, false)).toBe(true)
  })

  it('the colors and the annotations make a new image', () => {
    const e = emptyEdit()
    e.colors.contrast = 10
    expect(changes(e, true)).toBe(true)
  })

  it('the free crop is a PNG, named so', () => {
    const e = emptyEdit('livre')
    e.livre = {points: [{x: 0, y: 0}, {x: 9, y: 0}, {x: 0, y: 9}], rotation: 0}
    expect(resultType(e, 'image/jpeg')).toBe('image/png')
    expect(resultName('cheque.jpeg', 'image/png')).toBe('cheque.png')
    expect(resultType(emptyEdit(), 'image/jpeg')).toBe('image/jpeg')
  })

  it('an edit saved, completed with what it lacks', () => {
    const e = parseEdit({mode: 'perspectiva', colors: {contrast: 5}})
    expect(e.mode).toBe('perspectiva')
    expect(e.colors.contrast).toBe(5)
    expect(e.colors.gamma).toBe(1)
    expect(e.annotations).toEqual([])
  })
})

describe('a perspective on the whole image', () => {
  it('cuts nothing: no new image for it', async () => {
    const {normalized} = await import('../edit')
    const {defaultQuad} = await import('../perspective')
    const e = emptyEdit('perspectiva')
    e.perspectiva = defaultQuad(800, 600)
    expect(changes(normalized(e, 800, 600), true)).toBe(false)
    e.perspectiva[0] = {x: 10, y: 0}
    expect(changes(normalized(e, 800, 600), true)).toBe(true)
  })
})

describe('the fixed crop of no size', () => {
  it('cuts the image (a new image); the whole image cuts nothing', async () => {
    const {normalized} = await import('../edit')
    const e = emptyEdit('fixo')
    e.fixo = {x: 0, y: 0, width: 800, height: 600}
    expect(changes(normalized(e, 800, 600), false)).toBe(false)
    e.fixo = {x: 100, y: 0, width: 600, height: 600}
    expect(changes(normalized(e, 800, 600), false)).toBe(true)
    // of a size: its crop (CropOption), no new image
    expect(changes(normalized(e, 800, 600), true)).toBe(false)
  })
})

describe('the free rectangle', () => {
  it('its size is the result\'s; the whole image, straight, cuts nothing; no transparency', async () => {
    const {normalized, croppedSize, transparent} = await import('../edit')
    const e = emptyEdit('retangulo')
    e.retangulo = {cx: 400, cy: 300, width: 800, height: 600, angle: 0}
    expect(changes(normalized(e, 800, 600), false)).toBe(false)
    e.retangulo = {cx: 400, cy: 300, width: 500, height: 200, angle: 12}
    expect(changes(normalized(e, 800, 600), false)).toBe(true)
    expect(croppedSize(e, 800, 600, false)).toEqual({width: 500, height: 200})
    expect(transparent(e)).toBe(false)
  })
})
