// The footnotes of a Markdown text — `text[^label]`, `[^label]: the note`
// (its next lines indented go on with it) —, for the footnotes builder
// (mdkit-footnote-builder): read (parse), and the text changed — a
// reference put, a note made, edited, renamed, removed; the notes put
// together at the end in order (organize) —; its labels changed by a map
// (relabel, labelsFor: texts put together, a book, keep their notes apart).
// The same reading as the server's (rvq x/markup/footnotes.go): what is
// code is none. Text only (no DOM).

// cleanLabel is s made a label: letters, digits, "_", "." and "-" (what an
// id may hold: the rich text keeps the labels in its ids).
export const cleanLabel = (s) => String(s || '').replace(/[^A-Za-z0-9_.-]+/g, '')

const defLineRe = /^[ ]{0,3}\[\^([^\]\s]+)\]:[ \t]?/
const refRe = /\[\^([^\]\s]+)\]/g
const fenceRe = /^[ ]{0,3}(`{3,}|~{3,})/

// codeRanges are the parts of the text that are code: its fenced blocks, its
// spans (a run of backticks to the next as long).
function codeRanges(s) {
  const out = []
  let pos = 0
  let fence = ''
  let start = 0
  for (const line of s.split(/(?<=\n)/)) {
    const trimmed = line.replace(/^ +/, '')
    if (!fence) {
      const m = fenceRe.exec(line)
      if (m) {
        fence = m[1]
        start = pos
      }
    } else if (trimmed.startsWith(fence) && trimmed.replace(new RegExp('^\\' + fence[0] + '+'), '').trim() === '') {
      out.push([start, pos + line.length])
      fence = ''
    }
    pos += line.length
  }
  if (fence) out.push([start, s.length])
  const inside = (i) => out.some(([a, b]) => i >= a && i < b)
  for (let i = 0; i < s.length;) {
    if (s[i] !== '`' || inside(i)) {
      i++
      continue
    }
    let n = 0
    while (s[i + n] === '`') n++
    let end = -1
    for (let j = i + n; j < s.length;) {
      if (s[j] !== '`') {
        j++
        continue
      }
      let k = 0
      while (s[j + k] === '`') k++
      if (k === n) {
        end = j + k
        break
      }
      j += k
    }
    if (end < 0) {
      i += n
      continue
    }
    out.push([i, end])
    i = end
  }
  return out
}

// parse reads the footnotes of text:
//   - notes: [{ label, number, refs, text, defined, from, to }] — referenced
//     ones by number, then the others; from/to, the lines of its definition
//   - spans: [{ start, end, label, def }] — each label in the text
export function parse(text) {
  const s = text || ''
  const code = codeRanges(s)
  const inCode = (i) => code.some(([a, b]) => i >= a && i < b)
  const spans = []
  const defs = {}
  const defAt = new Set()
  const lines = s.split(/(?<=\n)/)
  let pos = 0
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const m = defLineRe.exec(line)
    if (m && !inCode(pos)) {
      const lstart = pos + line.indexOf('[^') + 2
      spans.push({ start: lstart, end: lstart + m[1].length, label: m[1], def: true })
      defAt.add(lstart - 2)
      const from = pos
      const body = [line.slice(m[0].length).replace(/\n$/, '')]
      while (i + 1 < lines.length) {
        const next = lines[i + 1]
        const blank = next.trim() === ''
        if (blank && i + 2 < lines.length && /^( {4}|\t)/.test(lines[i + 2])) body.push('')
        else if (!blank && /^( {4}|\t)/.test(next)) body.push(next.replace(/^( {4}|\t)/, '').replace(/\n$/, ''))
        else break
        pos += lines[i].length
        i++
      }
      if (!defs[m[1]]) defs[m[1]] = { text: body.join('\n').trim(), from, to: pos + lines[i].length }
    }
    pos += lines[i].length
  }
  for (const m of s.matchAll(refRe)) {
    if (inCode(m.index) || defAt.has(m.index) || s[m.index + m[0].length] === ':') continue
    spans.push({ start: m.index + 2, end: m.index + 2 + m[1].length, label: m[1], def: false })
  }
  spans.sort((a, b) => a.start - b.start)
  const byLabel = new Map()
  const note = (label) => {
    if (!byLabel.has(label)) byLabel.set(label, { label, number: 0, refs: 0, text: '', defined: false, from: -1, to: -1 })
    return byLabel.get(label)
  }
  let number = 0
  for (const sp of spans) {
    const n = note(sp.label)
    if (sp.def) {
      Object.assign(n, { defined: true, text: defs[sp.label].text, from: defs[sp.label].from, to: defs[sp.label].to })
    } else {
      n.refs++
      if (!n.number) n.number = ++number
    }
  }
  const notes = [...byLabel.values()]
  notes.sort((a, b) => (a.number || Infinity) - (b.number || Infinity))
  return { notes, spans }
}

// problems are the notes with no reference, and the references with no note.
export function problems(text) {
  const { notes } = parse(text)
  return {
    unreferenced: notes.filter((n) => n.defined && !n.refs).map((n) => n.label),
    undefined: notes.filter((n) => !n.defined).map((n) => n.label),
  }
}

// nextLabel is the next number not a label of text.
export function nextLabel(text) {
  const used = new Set(parse(text).notes.map((n) => n.label))
  let n = 1
  while (used.has(String(n))) n++
  return String(n)
}

// randomLabel is a label of a new note none of used holds: "n" and four
// letters or digits ("n4k2x") — not renumbered when the notes are organized,
// and not another text's when texts are put together (a book) —; rand, the
// source of chance (Math.random).
export function randomLabel(used, rand = Math.random) {
  const chars = 'abcdefghijkmnpqrstuvwxyz23456789'
  for (;;) {
    let l = 'n'
    for (let i = 0; i < 4; i++) l += chars[Math.floor(rand() * chars.length)]
    if (!used.has(l)) return l
  }
}

// newLabel is the label of a new note of text (randomLabel).
export function newLabel(text, rand = Math.random) {
  return randomLabel(new Set(parse(text).notes.map((n) => n.label)), rand)
}

// definition is the definition of a note: its next lines indented.
const definition = (label, body) => {
  const lines = String(body || '').trim().split('\n')
  return '[^' + label + ']: ' + lines.map((l, i) => (i && l ? '    ' + l : l)).join('\n')
}

// insertRef is text with a reference to label at pos; cursor, after it.
export function insertRef(text, pos, label) {
  const ref = '[^' + label + ']'
  return { text: text.slice(0, pos) + ref + text.slice(pos), cursor: pos + ref.length }
}

// create is text with a new note — label (a new one, newLabel, when none
// or taken) and its body — referenced at pos, its definition at the end.
export function create(text, pos, body, label) {
  let l = cleanLabel(label)
  if (!l || parse(text).notes.some((n) => n.label === l)) l = newLabel(text)
  const r = insertRef(text, pos, l)
  const t = r.text.replace(/\s*$/, '') + '\n\n' + definition(l, body) + '\n'
  return { text: t, cursor: r.cursor, label: l }
}

// edit is text with the note label's body; a note it had not, at the end.
export function edit(text, label, body) {
  const n = parse(text).notes.find((x) => x.label === label)
  if (!n || !n.defined) return text.replace(/\s*$/, '') + '\n\n' + definition(label, body) + '\n'
  const end = text[n.to - 1] === '\n' ? '\n' : ''
  return text.slice(0, n.from) + definition(label, body) + end + text.slice(n.to)
}

// relabel is text with its labels changed as map says (old → new): in its
// references and its notes; a new one made a label.
export function relabel(text, map) {
  const { spans } = parse(text)
  let out = ''
  let last = 0
  for (const sp of spans) {
    const to = cleanLabel(map[sp.label])
    if (!(sp.label in map) || !to) continue
    out += text.slice(last, sp.start) + to
    last = sp.end
  }
  return out + text.slice(last)
}

// rename is text with the note from named to (references and note); false,
// to taken (or none).
export function rename(text, from, to) {
  const l = cleanLabel(to)
  if (!l || (l !== from && parse(text).notes.some((n) => n.label === l))) return false
  return relabel(text, { [from]: l })
}

// remove is text without the note label: its references and its note.
export function remove(text, label) {
  const { notes, spans } = parse(text)
  // what goes: its references ([^label]) and its note (its lines)
  const cut = spans.filter((sp) => sp.label === label && !sp.def).map((sp) => [sp.start - 2, sp.end + 1])
  const n = notes.find((x) => x.label === label)
  if (n && n.defined) cut.push([n.from, n.to])
  cut.sort((a, b) => a[0] - b[0])
  let out = ''
  let last = 0
  for (const [a, b] of cut) {
    out += text.slice(last, a)
    last = b
  }
  return (out + text.slice(last)).replace(/\n{3,}/g, '\n\n')
}

// organize is text with its notes together at its end, in the order of
// their first references (the ones with none, last); with renumber, the
// labels that are numbers made 1, 2, … in that order (names as they are).
export function organize(text, renumber = false) {
  const { notes } = parse(text)
  const defined = notes.filter((n) => n.defined)
  if (!defined.length) return text
  // the text without its notes
  let body = ''
  let last = 0
  for (const n of [...defined].sort((a, b) => a.from - b.from)) {
    body += text.slice(last, n.from)
    last = n.to
  }
  body = (body + text.slice(last)).replace(/\n{3,}/g, '\n\n').replace(/^\s*\n/, '').replace(/\s*$/, '')
  let out = body + '\n\n' + defined.map((n) => definition(n.label, n.text)).join('\n') + '\n'
  if (renumber) {
    const map = {}
    const taken = new Set(notes.filter((n) => !/^\d+$/.test(n.label)).map((n) => n.label))
    let k = 0
    for (const n of notes) {
      if (!/^\d+$/.test(n.label)) continue
      do k++
      while (taken.has(String(k)))
      map[n.label] = '\u0000' + k
    }
    // through a mark: no label taken while others are changed
    out = relabelRaw(out, map).replace(/\u0000(\d+)/g, '$1')
  }
  return out
}

// relabelRaw is relabel with labels as they are (the mark of organize).
function relabelRaw(text, map) {
  const { spans } = parse(text)
  let out = ''
  let last = 0
  for (const sp of spans) {
    if (!(sp.label in map)) continue
    out += text.slice(last, sp.start) + map[sp.label]
    last = sp.end
  }
  return out + text.slice(last)
}

// labelsFor is the map that makes the labels of text its own — prefix and
// the label ("cap2-1") — (relabel: a text put with others).
export function labelsFor(text, prefix) {
  const p = cleanLabel(prefix)
  const out = {}
  if (!p) return out
  for (const n of parse(text).notes) out[n.label] = p + '-' + n.label
  return out
}

// uniqueLabels is the map that makes the labels of text none of taken's (a
// text it is put in, other texts with it): prefix and the label
// ("capitulo1-1"), a number after it while it is taken ("capitulo1-1-2").
export function uniqueLabels(text, taken, prefix) {
  const used = new Set(parse(taken || '').notes.map((n) => n.label))
  const p = cleanLabel(prefix) || 'nota'
  const out = {}
  for (const n of parse(text).notes) {
    let l = p + '-' + n.label
    for (let k = 2; used.has(l); k++) l = p + '-' + n.label + '-' + k
    used.add(l)
    out[n.label] = l
  }
  return out
}
