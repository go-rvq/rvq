import {describe, expect, it} from 'vitest'
import {orientedSize, readOrientation} from '../exif'

// a JPEG with only an APP1 Exif of one entry, the orientation o
function jpeg(o: number, little = false): ArrayBuffer {
  const b = new ArrayBuffer(2 + 2 + 2 + 6 + 8 + 2 + 12 + 4)
  const v = new DataView(b)
  let p = 0
  v.setUint16(p, 0xffd8); p += 2
  v.setUint16(p, 0xffe1); p += 2
  v.setUint16(p, b.byteLength - 4); p += 2
  v.setUint32(p, 0x45786966); p += 4
  v.setUint16(p, 0); p += 2
  const tiff = p
  v.setUint16(p, little ? 0x4949 : 0x4d4d); p += 2
  v.setUint16(p, 42, little); p += 2
  v.setUint32(p, 8, little); p += 4
  v.setUint16(tiff + 8, 1, little)
  const e = tiff + 10
  v.setUint16(e, 0x0112, little)
  v.setUint16(e + 2, 3, little)
  v.setUint32(e + 4, 1, little)
  v.setUint16(e + 8, o, little)
  return b
}

describe('the camera\'s orientation', () => {
  it('reads it, big and little endian', () => {
    expect(readOrientation(jpeg(6))).toBe(6)
    expect(readOrientation(jpeg(8, true))).toBe(8)
  })

  it('1 when there is none, or not a JPEG', () => {
    expect(readOrientation(new Uint8Array([0x89, 0x50, 0x4e, 0x47]).buffer)).toBe(1)
  })

  it('5 to 8 swap the sides', () => {
    expect(orientedSize(4000, 3000, 6)).toEqual({width: 3000, height: 4000})
    expect(orientedSize(4000, 3000, 3)).toEqual({width: 4000, height: 3000})
  })
})
