// bun test (site/components): the footnotes of a Markdown text.
import { expect, test } from 'bun:test'
import * as F from './footnote-model.js'

const md = 'Um[^a] dois[^fonte] tres[^a] e `[^code]`.\n\n```\n[^x]: no\n```\n\n[^a]: Nota **a**.\n[^fonte]: Nota b\n    continua.\n[^sobra]: Sem referência.\n'

test('parse: the notes by number, their texts; code none', () => {
  const { notes } = F.parse(md)
  expect(notes.map((n) => [n.label, n.number, n.refs, n.text, n.defined])).toEqual([
    ['a', 1, 2, 'Nota **a**.', true],
    ['fonte', 2, 1, 'Nota b\ncontinua.', true],
    ['sobra', 0, 0, 'Sem referência.', true],
  ])
  expect(F.problems(md)).toEqual({ unreferenced: ['sobra'], undefined: [] })
  expect(F.problems('x[^nada]')).toEqual({ unreferenced: [], undefined: ['nada'] })
})

test('a new label: random, none of the text\'s', () => {
  // the chance in turn: 0, then the first char again (taken), then another
  const seq = [0, 0, 0, 0, 0, 0, 0, 0, 0.5, 0.5, 0.5, 0.5]
  let i = 0
  const rand = () => seq[i++ % seq.length]
  expect(F.newLabel('x[^naaaa]', rand)).toBe('nssss')
  expect(F.newLabel('')).toMatch(/^n[a-z2-9]{4}$/)
})

test('create, insert a reference, edit', () => {
  let r = F.create('Texto aqui.', 5, 'A nota\ncom duas linhas', '1')
  expect(r.label).toBe('1')
  expect(r.text).toBe('Texto[^1] aqui.\n\n[^1]: A nota\n    com duas linhas\n')
  expect(r.cursor).toBe(9)
  // a name taken, or none: a new one (random)
  r = F.create(r.text, 0, 'Outra', '1')
  expect(r.label).toMatch(/^n[a-z2-9]{4}$/)
  expect(F.create('t', 0, 'x').label).toMatch(/^n[a-z2-9]{4}$/)
  r = F.create(r.text, 0, 'Com nome', 'fonte ç!')
  expect(r.label).toBe('fonte')
  const i = F.insertRef('ab', 1, 'fonte')
  expect(i).toEqual({ text: 'a[^fonte]b', cursor: 9 })
  const e = F.edit('x[^1]\n\n[^1]: velha\n    linha\n\nfim', '1', 'nova')
  expect(e).toBe('x[^1]\n\n[^1]: nova\n\nfim')
})

test('rename, relabel by a map, remove', () => {
  expect(F.rename(md, 'a', 'b')).toContain('Um[^b] dois[^fonte] tres[^b]')
  expect(F.rename(md, 'a', 'fonte')).toBe(false)
  const re = F.relabel(md, { a: 'c1-a', fonte: 'c1 fonte' })
  expect(re).toContain('[^c1-a]: Nota')
  expect(re).toContain('[^c1fonte]: Nota b')
  expect(re).toContain('`[^code]`')
  expect(re).toContain('[^x]: no')
  const rm = F.remove(md, 'a')
  expect(rm).toContain('Um dois[^fonte] tres e')
  expect(rm).not.toContain('Nota **a**')
  expect(F.parse(rm).notes.map((n) => n.label)).toEqual(['fonte', 'sobra'])
})

test('organize: the notes at the end, in order; renumbered', () => {
  const t = '[^2]: Segunda.\n\nUm[^2] e[^1] e[^nome].\n\n[^1]: Primeira.\n\nMeio.\n\n[^nome]: Com nome.\n'
  expect(F.organize(t)).toBe('Um[^2] e[^1] e[^nome].\n\nMeio.\n\n[^2]: Segunda.\n[^1]: Primeira.\n[^nome]: Com nome.\n')
  expect(F.organize(t, true)).toBe('Um[^1] e[^2] e[^nome].\n\nMeio.\n\n[^1]: Segunda.\n[^2]: Primeira.\n[^nome]: Com nome.\n')
})

test('a book: each chapter its labels, together no conflict', () => {
  const ch1 = 'Um[^1].\n\n[^1]: Do primeiro.\n'
  const ch2 = 'Dois[^1].\n\n[^1]: Do segundo.\n'
  const book = F.relabel(ch1, F.labelsFor(ch1, 'cap1')) + '\n' + F.relabel(ch2, F.labelsFor(ch2, 'cap2'))
  expect(F.parse(book).notes.map((n) => [n.label, n.number, n.text])).toEqual([
    ['cap1-1', 1, 'Do primeiro.'], ['cap2-1', 2, 'Do segundo.'],
  ])
  expect(F.labelsFor(ch1, '')).toEqual({})
})

test('files put in a text: their labels none of the text\'s, nor each other\'s', () => {
  const text = 'Texto[^1] e[^cap-1].\n\n[^1]: Do texto.\n[^cap-1]: Já existe.\n'
  const f1 = 'Um[^1].\n\n[^1]: Do arquivo.\n'
  const m1 = F.uniqueLabels(f1, text, 'cap.md é')
  expect(m1).toEqual({ 1: 'cap.md-1' })
  const m2 = F.uniqueLabels(f1, text, 'cap')
  expect(m2).toEqual({ 1: 'cap-1-2' })
  const all = text + '\n' + F.relabel(f1, m2)
  const m3 = F.uniqueLabels(f1, all, 'cap')
  expect(m3).toEqual({ 1: 'cap-1-3' })
})
