// tiptapFootnoteAdapter is the footnotes builder's (footnote-builder) way to
// the text of a tiptap editor that has the nodes of footnote-tiptap — the
// rich text editor's, the admin's (vx-tiptap-editor, its plugin) —: a state
// of its own (the editor's at first), the cursor followed in it, and the
// reference clicked (refPos: its place) to point at another note; done, the
// editor's text made it in a change of its history — the notes in the order
// of their references, the labels that are numbers 1, 2, … —.
import { numbers } from './footnote-tiptap.js'
import { cleanLabel, randomLabel } from './footnote-model.js'

export function tiptapFootnoteAdapter(editor, refPos = null) {
  let st = editor.state
  let pos = st.selection.from
  // the reference clicked (its place, followed as the text changes)
  let ref = refPos
  const refType = st.schema.nodes.footnoteRef
  const notesType = st.schema.nodes.footnotes
  const apply = (tr) => {
    pos = tr.mapping.map(pos)
    if (ref != null) {
      const r = tr.mapping.mapResult(ref)
      ref = r.deleted ? null : r.pos
    }
    st = st.apply(tr)
  }
  // keep applies tr, a change of labels only (positions as they are):
  // the reference clicked kept where it is
  const keep = (tr) => {
    const at = ref
    apply(tr)
    if (at != null && st.doc.nodeAt(at) && st.doc.nodeAt(at).type === refType) ref = at
  }
  const refNode = () => {
    const n = ref != null && ref < st.doc.content.size ? st.doc.nodeAt(ref) : null
    return n && n.type === refType ? n : null
  }
  const notesAt = () => {
    let found = null
    st.doc.descendants((n, p) => {
      if (found) return false
      if (n.type === notesType) found = { node: n, pos: p }
    })
    return found
  }
  const stored = () => { const f = notesAt(); return f ? [...(f.node.attrs.notes || [])] : [] }
  const setNotes = (notes) => {
    const f = notesAt()
    const tr = st.tr
    if (f) tr.setNodeMarkup(f.pos, null, { notes })
    else tr.insert(st.doc.content.size, notesType.create({ notes }))
    apply(tr)
  }
  const refs = () => {
    const out = []
    st.doc.descendants((n, p) => { if (n.type === refType) out.push({ label: n.attrs.label, pos: p }) })
    return out
  }
  const relabel = (map) => {
    const tr = st.tr
    for (const r of refs()) if (r.label in map) tr.setNodeMarkup(r.pos, null, { label: map[r.label] })
    keep(tr)
    setNotes(stored().map((n) => (n.label in map ? { ...n, label: map[n.label] } : n)))
  }
  const list = () => {
    const nums = numbers(st.doc)
    const count = new Map()
    for (const r of refs()) count.set(r.label, (count.get(r.label) || 0) + 1)
    const notes = stored()
    const labels = [...new Set([...nums.keys(), ...notes.map((n) => n.label)])]
    return labels.map((label) => {
      const n = notes.find((x) => x.label === label)
      return { label, number: nums.get(label) || 0, refs: count.get(label) || 0, text: n ? n.html : '', defined: !!n }
    }).sort((a, b) => (a.number || Infinity) - (b.number || Infinity))
  }
  // the label of a new note (random: footnote-model randomLabel)
  const next = () => randomLabel(new Set(list().map((n) => n.label)))
  return {
    kind: 'rich',
    notes: list,
    next,
    create: (html, label) => {
      let l = cleanLabel(label)
      if (!l || list().some((n) => n.label === l)) l = next()
      setNotes([...stored(), { label: l, html }])
      return l
    },
    insert: (label) => apply(st.tr.insert(pos, refType.create({ label }))),
    // the note of the reference clicked; point makes it another's
    pointed: () => { const n = refNode(); return n ? n.attrs.label : null },
    point: (label) => { if (refNode()) keep(st.tr.setNodeMarkup(ref, null, { label })) },
    edit: (label, html) => {
      const notes = stored()
      const i = notes.findIndex((n) => n.label === label)
      if (i < 0) notes.push({ label, html })
      else notes[i] = { ...notes[i], html }
      setNotes(notes)
    },
    rename: (from, to) => {
      const l = cleanLabel(to)
      if (!l || (l !== from && list().some((n) => n.label === l))) return false
      relabel({ [from]: l })
      return true
    },
    remove: (label) => {
      const tr = st.tr
      for (const r of refs().reverse()) if (r.label === label) tr.delete(r.pos, r.pos + 1)
      apply(tr)
      setNotes(stored().filter((n) => n.label !== label))
    },
    done: () => {
      // organized: the notes in the order of their references, the
      // labels that are numbers renumbered (through marks: none taken)
      const order = list()
      const named = new Set(order.filter((n) => !/^\d+$/.test(n.label)).map((n) => n.label))
      const map = {}
      let k = 0
      for (const n of order) {
        if (!/^\d+$/.test(n.label)) continue
        do k++
        while (named.has(String(k)))
        map[n.label] = '\u0000' + k
      }
      relabel(map)
      relabel(Object.fromEntries(Object.values(map).map((v) => [v, v.slice(1)])))
      const byNumber = numbers(st.doc)
      setNotes(stored().map((n, i) => ({ n, i }))
        .sort((a, b) => ((byNumber.get(a.n.label) || Infinity) - (byNumber.get(b.n.label) || Infinity)) || a.i - b.i)
        .map((x) => x.n))
      // the notes at the end of the text
      const f = notesAt()
      if (f && f.pos + f.node.nodeSize !== st.doc.content.size) {
        const tr = st.tr
        tr.delete(f.pos, f.pos + f.node.nodeSize)
        tr.insert(tr.doc.content.size, f.node)
        apply(tr)
      }
      // the editor's text made this one: a change of its history
      const tr = editor.state.tr.replaceWith(0, editor.state.doc.content.size, st.doc.content)
      editor.view.dispatch(tr)
      editor.commands.focus(Math.min(pos, editor.state.doc.content.size))
    },
  }
}
