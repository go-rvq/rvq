// bun test (site/components): the table of the table builder.
import { expect, test } from 'bun:test'
import * as T from './table-model.js'

const filled = () => {
  const t = T.create(3, 3)
  t.cells = [['A', 'B', 'C'], ['1', '2', '3'], ['4', '5', '6']]
  return t
}

test('GFM: a header, no merge; its alignment, pipes escaped', () => {
  const t = filled()
  t.aligns = ['', 'center', 'right']
  t.cells[1][0] = 'a|b'
  expect(T.toMarkdown(t)).toBe('| A | B | C |\n| --- | :-: | --: |\n| a\\|b | 2 | 3 |\n| 4 | 5 | 6 |')
  const back = T.fromMarkdown(T.toMarkdown(t).split('\n'))
  expect(back.cells).toEqual(t.cells)
  expect(back.aligns).toEqual(t.aligns)
})

test('merged cells: HTML, read back the same', () => {
  const t = filled()
  expect(T.merge(t, 0, 0, 0, 1)).toBe(true)
  expect(t.cells[0][0]).toBe('A B')
  expect(T.merge(t, 1, 2, 2, 2)).toBe(true)
  const html = T.toMarkdown(t)
  expect(html).toContain('<th colspan="2">A B</th>')
  expect(html).toContain('<td rowspan="2">3 6</td>')
  expect(html).not.toContain('<td>6</td>')
  const back = T.fromHTML(html)
  expect(back.merges).toEqual(t.merges)
  expect(back.cells).toEqual(t.cells)
  expect(back.header).toBe(true)
})

test('no header: HTML; each cell its Markdown, read back the same', () => {
  const t = filled()
  t.header = false
  t.cells[0][0] = '**x** & `a<b`'
  t.cells[0][1] = ''
  const html = T.toMarkdown(t)
  expect(html).not.toContain('<thead>')
  expect(html).toContain('<td>\n\n**x** & `a<b`\n\n</td>')
  expect(html).toContain('<td>4</td>')
  expect(html).toContain('<td></td>')
  expect(T.fromHTML(html).cells).toEqual(t.cells)
})

test('GFM: the lines of a cell by <br>, read back', () => {
  const t = filled()
  t.cells[1][1] = '**one**\ntwo'
  expect(T.toMarkdown(t)).toContain('| **one**<br>two |')
  expect(T.fromMarkdown(T.toMarkdown(t).split('\n')).cells[1][1]).toBe('**one**\ntwo')
})

test('the HTML of another: its tags made Markdown', () => {
  const t = T.fromHTML('<table><tr><td><strong>B</strong> <em>i</em> <code>c</code> <a href="https://e.com">l</a><br>x &amp; y</td></tr></table>')
  expect(t.cells[0][0]).toBe('**B** _i_ `c` [l](https://e.com)\nx & y')
})

test('a selection grows to the merges it touches', () => {
  const t = filled()
  T.merge(t, 0, 1, 1, 2)
  expect(T.span(t, 0, 0, 0, 1)).toEqual({ r1: 0, c1: 0, r2: 1, c2: 2 })
  // a merge within a merge: one
  T.merge(t, 0, 0, 2, 2)
  expect(t.merges).toEqual([{ r: 0, c: 0, rs: 3, cs: 3 }])
  expect(T.unmerge(t, 1, 1)).toBe(true)
  expect(t.merges).toEqual([])
})

test('rows and columns in and out; the merges follow', () => {
  const t = filled()
  T.merge(t, 0, 0, 1, 0)
  T.insertRow(t, 1) // within the merge: it grows
  expect(t.merges[0]).toEqual({ r: 0, c: 0, rs: 3, cs: 1 })
  T.insertRow(t, 0) // above: it moves
  expect(t.merges[0]).toEqual({ r: 1, c: 0, rs: 3, cs: 1 })
  expect(t.rows).toBe(5)
  T.removeRow(t, 1) // its first row: its text goes down
  expect(t.merges[0]).toEqual({ r: 1, c: 0, rs: 2, cs: 1 })
  expect(t.cells[1][0]).toBe('A 1')
  T.removeRow(t, 1)
  expect(t.merges).toEqual([]) // one cell: no merge
  T.insertCol(t, 1)
  expect(t.cols).toBe(4)
  expect(t.aligns.length).toBe(4)
  T.removeCol(t, 1)
  expect(t.cells[1]).toEqual(['A 1', '2', '3'])
  const one = T.create(1, 1)
  T.removeRow(one, 0)
  T.removeCol(one, 0)
  expect([one.rows, one.cols]).toEqual([1, 1])
})

test('the table under the cursor', () => {
  const lines = ['text', '| A | B |', '| --- | --- |', '| 1 | 2 |', '', '<table>', '<tr><td colspan="2">x</td></tr>', '</table>', 'after']
  const g = T.tableAt(lines, 3)
  expect([g.from, g.to, g.table.cells[1][1]]).toEqual([1, 3, '2'])
  const h = T.tableAt(lines, 6)
  expect([h.from, h.to, h.table.merges.length]).toEqual([5, 7, 1])
  expect(T.tableAt(lines, 0)).toBe(null)
  expect(T.tableAt(lines, 8)).toBe(null)
})

test('the widths of the columns: HTML (col width), read back; with a column in and out', () => {
  const t = filled()
  t.widths = ['120', '', '']
  let html = T.toMarkdown(t)
  expect(html.startsWith('<table>\n<colgroup><col width="120"><col><col></colgroup>')).toBe(true)
  expect(T.fromHTML(html).widths).toEqual(['120', '', ''])
  t.widths = ['50%', '25%', '25%']
  html = T.toMarkdown(t)
  expect(html.startsWith('<table width="100%">')).toBe(true)
  expect(T.fromHTML(html).widths).toEqual(['50%', '25%', '25%'])
  T.insertCol(t, 1)
  expect(t.widths).toEqual(['50%', '', '25%', '25%'])
  T.removeCol(t, 0)
  expect(t.widths).toEqual(['', '25%', '25%'])
  // no width: GFM again
  t.widths = ['', '', '']
  expect(T.toMarkdown(t).startsWith('|')).toBe(true)
  // a table of before: its widths
  expect(T.normalize({ rows: 1, cols: 2, cells: [['a', 'b']], merges: [], header: true, aligns: [''] }).widths).toEqual(['', ''])
})

test('the rich text editor\'s table: from its HTML (colwidth), to it', () => {
  const html = '<table style="min-width: 100px"><colgroup><col style="width: 150px"><col style="min-width: 25px"></colgroup><tbody>' +
    '<tr><th colspan="1" rowspan="1" colwidth="150"><p>A <strong>b</strong></p></th><th colspan="1" rowspan="1"><p>B</p></th></tr>' +
    '<tr><td colspan="2" rowspan="1"><p>x</p></td></tr></tbody></table>'
  const t = T.fromHTML(html)
  expect(t.cells).toEqual([['A **b**', 'B'], ['x', '']])
  expect(t.header).toBe(true)
  expect(t.merges).toEqual([{ r: 1, c: 0, rs: 1, cs: 2 }])
  expect(t.widths).toEqual(['150', ''])
  t.aligns[1] = 'center'
  t.widths = ['50%', '']
  const out = T.toRichHTML(t, (s) => s.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>'), 600)
  expect(out).toBe('<table><tbody><tr><th colspan="1" rowspan="1" colwidth="300"><p>A <strong>b</strong></p></th>' +
    '<th colspan="1" rowspan="1"><p style="text-align: center">B</p></th></tr>' +
    '<tr><td colspan="2" rowspan="1" colwidth="300,0"><p>x</p></td></tr></tbody></table>')
})
