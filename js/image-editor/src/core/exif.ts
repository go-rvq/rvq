// The camera's orientation (EXIF 0x0112) of a JPEG, so a phone's photo does
// not open lying down: the image is read with it applied
// (createImageBitmap's imageOrientation "from-image") — readOrientation says
// what it was.

export function readOrientation(buf: ArrayBuffer): number {
  const v = new DataView(buf)
  if (v.byteLength < 4 || v.getUint16(0) !== 0xffd8) return 1
  let off = 2
  while (off + 4 <= v.byteLength) {
    const marker = v.getUint16(off)
    const len = v.getUint16(off + 2)
    if (marker === 0xffe1 && off + 10 <= v.byteLength && v.getUint32(off + 4) === 0x45786966) {
      const tiff = off + 10
      const little = v.getUint16(tiff) === 0x4949
      const u16 = (p: number) => v.getUint16(p, little)
      const u32 = (p: number) => v.getUint32(p, little)
      const ifd = tiff + u32(tiff + 4)
      if (ifd + 2 > v.byteLength) return 1
      const n = u16(ifd)
      for (let i = 0; i < n; i++) {
        const e = ifd + 2 + i * 12
        if (e + 12 > v.byteLength) break
        if (u16(e) === 0x0112) {
          const o = u16(e + 8)
          return o >= 1 && o <= 8 ? o : 1
        }
      }
      return 1
    }
    if ((marker & 0xff00) !== 0xff00 || marker === 0xffda) break
    off += 2 + len
  }
  return 1
}

// orientedSize: the size of an image of w × h shown in orientation o (5 to
// 8 turn it a quarter: the sides swap).
export function orientedSize(w: number, h: number, o: number): {width: number; height: number} {
  return o >= 5 && o <= 8 ? {width: h, height: w} : {width: w, height: h}
}

// loadBitmap: the image of a file or a URL, its orientation applied; never
// larger than maxPixels (a photo of 48 MP made smaller, to edit and send
// faster) when it is given.
export async function loadBitmap(src: Blob | string, maxPixels = 0): Promise<ImageBitmap | HTMLImageElement> {
  let blob: Blob
  if (typeof src === 'string') {
    const res = await fetch(src, {credentials: 'same-origin'})
    if (!res.ok) throw new Error(`image-editor: ${src}: ${res.status}`)
    blob = await res.blob()
  } else {
    blob = src
  }
  let bmp: ImageBitmap | HTMLImageElement
  if (typeof createImageBitmap === 'function') {
    bmp = await createImageBitmap(blob, {imageOrientation: 'from-image'} as ImageBitmapOptions)
  } else {
    bmp = await new Promise<HTMLImageElement>((resolve, reject) => {
      const img = new Image()
      img.onload = () => resolve(img)
      img.onerror = reject
      img.src = URL.createObjectURL(blob)
    })
  }
  const w = 'naturalWidth' in bmp ? bmp.naturalWidth : bmp.width
  const h = 'naturalHeight' in bmp ? bmp.naturalHeight : bmp.height
  if (maxPixels > 0 && w * h > maxPixels && typeof createImageBitmap === 'function') {
    const k = Math.sqrt(maxPixels / (w * h))
    return createImageBitmap(bmp, {resizeWidth: Math.round(w * k), resizeHeight: Math.round(h * k), resizeQuality: 'high'})
  }
  return bmp
}
