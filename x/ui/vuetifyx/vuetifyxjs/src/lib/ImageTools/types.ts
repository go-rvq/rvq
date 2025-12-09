import ImageEditor from 'tui-image-editor'
import { ref, Ref } from 'vue'
import * as url from 'node:url'

export class URLValue {
  value: string

  constructor(value: string) {
    this.value = value
  }

  toRaw(): RawValue {
    const rImageType = /data:(image\/.+);base64,/,
      data = this.value
    let mimeString = ''
    let raw, uInt8Array, i

    raw = data.replace(rImageType, (header, imageType) => {
      mimeString = imageType
      return ''
    })

    return new RawValue(atob(raw), mimeString)
  }
}

export class RawValue {
  raw: string
  mime: string

  constructor(raw: string, mime: string) {
    this.raw = raw
    this.mime = mime
  }

  toBinary(): Promise<BinaryValue> {
    return new Promise<BinaryValue>((resolve) => {
      const rawLength = this.raw.length
      let uInt8Array, i
      uInt8Array = new Uint8Array(rawLength) // eslint-disable-line

      for (i = 0; i < rawLength; i += 1) uInt8Array[i] = this.raw.charCodeAt(i)

      resolve(new BinaryValue(uInt8Array, this.mime))
    })
  }
}

export class BinaryValue {
  data: Uint8Array
  mime: string

  constructor(data: Uint8Array, mime: string) {
    this.data = data
    this.mime = mime
  }

  toBlob(): Promise<BlobValue> {
    return new Promise<BlobValue>((resolve) => {
      resolve(new BlobValue(new Blob([this.data], { type: this.mime })))
    })
  }
}

export class BlobValue {
  value: Blob

  constructor(value: Blob) {
    this.value = value
  }

  toFile(): Promise<File> {
    return new Promise<File>((resolve) => {
      const typ = this.value.type.split('/')[1]

      let ext = ''
      if (typ === 'png' || typ == 'gif') {
        ext = typ
      } else if (typ === 'jpeg') {
        ext = 'jpg'
      }

      resolve(
        new File([this.value], 'image.' + ext, {
          type: this.value.type
        })
      )
    })
  }
}

type VisibilityHandler = (v: boolean) => void

export class Controller {
  element: HTMLDivElement
  editor: ImageEditor
  hidder: Hidder | null
  visibilityHandlers: VisibilityHandler[]
  _visible: boolean

  constructor(element: HTMLDivElement, editor: tuiImageEditor.ImageEditor, hidder: Hidder | null) {
    this.element = element
    this.editor = editor
    this.hidder = hidder
    this.visibilityHandlers = []
    this._visible = true
  }

  toURL(): Promise<URLValue> {
    return new Promise<URLValue>((resolve) => {
      resolve(new URLValue(this.editor.toDataURL()))
    })
  }

  getMethod(instance: any, methodName: string): any {
    const { first, rest } = parseDotMethodName(methodName)
    const isInstance = instance.constructor.name !== 'Object'
    const type = typeof instance[first]
    let obj

    if (isInstance && type === 'function') {
      obj = instance[first].bind(instance)
    } else {
      obj = instance[first]
    }

    if (rest.length > 0) {
      return this.getMethod(obj, rest)
    }

    return obj
  }

  invoke(methodName: string, ...args: any) {
    let ei = this.editor as any,
      result = null
    if (ei[methodName]) {
      result = ei[methodName](...args)
    } else if (methodName.indexOf('.') > -1) {
      const func = this.getMethod(ei, methodName)

      if (typeof func === 'function') {
        result = func(...args)
      }
    }

    return result
  }

  load(v: any): Promise<null> {
    if (v.constructor == String) {
      const s = v as string
      return new Promise<null>((resolve, reject) => {
        loadImage(s)
          .catch(reject)
          .then((el) => {
            if (!el) {
              reject("Failed to load image.");
              return;
            }
            const img = el as HTMLImageElement
            const canvas = document.createElement("canvas");
            const ctx = canvas.getContext("2d");
            canvas.width = img.width;
            canvas.height = img.height;
            ctx?.drawImage(img, 0, 0);

            // Convert canvas content to a Blob
            canvas.toBlob((blob) => {
              if (blob) {
                const f = new File([blob], baseName(s), { type: blob.type })
                ;(this.editor.ui as any)._actions.main.load(f)
              } else {
                reject("Failed to create Blob from image.");
              }
            }, "image/png")
        })
      })
    } else if (v.constructor == File) {
      return new Promise<null>((resolve) => {
        ;(this.editor.ui as any)._actions.main.load(v as File)
        resolve(null)
      })
    }
    return new Promise<null>((resolve, reject) => {
      reject(`unsupported value of type ${typeof v}`)
    })
  }

  onVisibilityChange(h: VisibilityHandler) {
    this.visibilityHandlers.push(h)
  }

  setVisibility(v: boolean) {
    if (v === this._visible) return
    this._visible = v

    if (this.hidder) {
      this.hidder.set(v)
    }

    this.visibilityHandlers.forEach((func) => {
      func(v)
    })
  }

  setContainer(el: HTMLElement | null) {
    if (el) {
      this.hidder = new Hidder(el as HTMLElement)
    } else {
      this.hidder = null
    }
  }
}

declare const window: any

enum HiddenMode {
  ALONE,
  BEFORE,
  AFTER
}

export class Hidder {
  container: HTMLElement
  refElement: HTMLElement
  mode: HiddenMode
  tempContainer: HTMLElement | null

  constructor(container: HTMLElement) {
    this.container = container
    this.tempContainer = null

    if (container.nextSibling) {
      this.refElement = container.nextSibling as HTMLElement
      this.mode = HiddenMode.BEFORE
    } else if (container.previousSibling) {
      this.refElement = container.previousSibling as HTMLElement
      this.mode = HiddenMode.AFTER
    } else {
      this.refElement = container.parentNode as HTMLElement
      this.mode = HiddenMode.ALONE
    }
  }

  hide() {
    if (this.tempContainer) return
    const temp = window.document.createElement('div')
    temp.style.display = 'none'
    temp.appendChild(this.container)
    window.document.body.appendChild(temp)
    this.tempContainer = temp
  }

  show() {
    if (!this.tempContainer) return

    if (this.mode == HiddenMode.ALONE) {
      this.refElement.appendChild(this.container)
    } else if (this.mode == HiddenMode.AFTER) {
      if (this.refElement.nextSibling) {
        this.refElement.parentNode?.insertBefore(this.container, this.refElement.nextSibling)
      } else {
        this.refElement.parentNode?.appendChild(this.container)
      }
    } else if (this.mode == HiddenMode.BEFORE) {
      this.refElement.parentNode?.insertBefore(this.container, this.refElement)
    }

    this.tempContainer = null
  }

  set(v: boolean) {
    if (v) this.show()
    else this.hide()
  }
}

export class Timer {
  private timerId: number | null = null
  private readonly callback: () => void
  private readonly duration: number

  constructor(duration: number, callback: () => void) {
    this.duration = duration
    this.callback = callback
  }

  start(): void {
    if (this.timerId === null) {
      this.timerId = window.setTimeout(() => {
        this.callback()
        this.timerId = null
      }, this.duration)
    }
  }

  restart(): void {
    this.stop()
    this.start()
  }

  stop(): void {
    if (this.timerId !== null) {
      clearTimeout(this.timerId)
      this.timerId = null
    }
  }
}

const parseDotMethodName = (methodName: string) => {
    const firstDotIdx = methodName.indexOf('.')
    let firstMethodName = methodName
    let restMethodName = ''

    if (firstDotIdx > -1) {
      firstMethodName = methodName.substring(0, firstDotIdx)
      restMethodName = methodName.substring(firstDotIdx + 1, methodName.length)
    }

    return {
      first: firstMethodName,
      rest: restMethodName
    }
  },
  baseName = (v: string): string => {
    const pos = v.lastIndexOf('/')
    return pos >= 0 ? v.substring(pos + 1) : v
  }

class TempImageStackEntry {
  public url: string
  public  load: (img: HTMLImageElement, next:()=>void)=>void
  public  error: (reason: any) => void

  constructor(url: string, load: (img: HTMLImageElement, next: () => void) => void, error: (reason: any) => void) {
    this.url = url
    this.load = load
    this.error = error
  }
}

function loadImage(url: string) {
  return new Promise<HTMLImageElement>((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = "anonymous"; // Important for images from different origins
    img.src = url; // Replace with your image path or URL

    img.onload = function() {
      resolve(img)
    };

    img.onerror = function() {
      console.error("Error loading image.");
    };
  })
}
