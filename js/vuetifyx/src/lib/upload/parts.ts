// The sending of a file of any size (admin/uploadparts): its sha256
// calculated in the browser, slice by slice; asked to the server (a content
// it has needs no sending); else sent in parts — 8MB each, 3 at once, a part
// failed tried again —, its file and hash given back (a token).

import {Sha256} from './sha256'

export const PART_SIZE = 8 << 20
const PARALLEL = 3
const TRIES = 3

export interface Token {
  id: string
  name: string
  size: number
  hash: string
  // mtime: when the file was last modified (ms, its lastModified)
  mtime?: number
}

// bytesOf are the bytes of b (FileReader where Blob.arrayBuffer is not).
function bytesOf(b: Blob): Promise<ArrayBuffer> {
  if (typeof (b as any).arrayBuffer === 'function') return b.arrayBuffer()
  return new Promise((resolve, reject) => {
    const r = new FileReader()
    r.onload = () => resolve(r.result as ArrayBuffer)
    r.onerror = () => reject(r.error)
    r.readAsArrayBuffer(b)
  })
}

// hashFile is the sha256 of file, read in slices (onProgress: of 0 to 1).
export async function hashFile(file: Blob, onProgress?: (p: number) => void): Promise<string> {
  const h = new Sha256()
  for (let at = 0; at < file.size; at += PART_SIZE) {
    const piece = file.slice(at, Math.min(at + PART_SIZE, file.size))
    h.update(new Uint8Array(await bytesOf(piece)))
    onProgress?.(Math.min(1, (at + PART_SIZE) / file.size))
  }
  onProgress?.(1)
  return h.hex()
}

// call posts the event ev of the model at url (its listing) with form: its
// data.
export async function call(url: string, ev: string, form: Record<string, string | Blob>): Promise<any> {
  const body = new FormData()
  for (const [k, v] of Object.entries(form)) body.append(k, v as any)
  const u = url + (url.includes('?') ? '&' : '?') + '__execute_event__=' + encodeURIComponent(ev)
  const res = await fetch(u, {method: 'POST', body, credentials: 'same-origin', headers: {'X-Plaid-Request': '1'}})
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  const out = await res.json()
  if (out && out.runScript && /color: "error"/.test(out.runScript)) {
    const m = /message: "((?:[^"\\]|\\.)*)"/.exec(out.runScript)
    throw new Error(m ? JSON.parse('"' + m[1] + '"') : 'error')
  }
  return out ? out.data : null
}

// exists is the server's data of the content of hash, null of none.
export function exists(url: string, hash: string): Promise<any> {
  return call(url, 'uploadparts_Exists', {hash})
}

function newID(): string {
  const r = typeof crypto !== 'undefined' && crypto.getRandomValues ? crypto.getRandomValues(new Uint8Array(12)) : null
  return 'up-' + (r ? Array.from(r, (b) => b.toString(16).padStart(2, '0')).join('') : Math.random().toString(16).slice(2) + Date.now())
}

// send sends file in parts (onProgress: of 0 to 1): its token.
export async function send(url: string, file: File, hash: string, onProgress?: (p: number) => void): Promise<Token> {
  const id = newID()
  const total = file.size
  const offsets: number[] = []
  for (let at = 0; at < total; at += PART_SIZE) offsets.push(at)
  let sent = 0
  const one = async (at: number) => {
    const part = file.slice(at, Math.min(at + PART_SIZE, total))
    for (let i = 1; ; i++) {
      try {
        await call(url, 'uploadparts_Part', {id, offset: String(at), total: String(total), part})
        break
      } catch (e) {
        if (i >= TRIES) throw e
        await new Promise((r) => setTimeout(r, 500 * i))
      }
    }
    sent += part.size
    onProgress?.(total ? sent / total : 1)
  }
  try {
    const queue = [...offsets]
    await Promise.all(Array.from({length: Math.min(PARALLEL, queue.length)}, async () => {
      while (queue.length) await one(queue.shift()!)
    }))
    return await call(url, 'uploadparts_Done', {id, total: String(total), name: file.name, hash, mtime: String(file.lastModified || 0)})
  } catch (e) {
    call(url, 'uploadparts_Abort', {id}).catch(() => {})
    throw e
  }
}
