import { describe, expect, it } from 'vitest'
import { masked } from '../mask'

// A mask writes a value as it is typed: # a digit, A a letter, * either,
// anything else itself; what does not fit its place is left out.
describe('masked', () => {
  it('writes a phone', () => {
    expect(masked('(###) ###-####', '7815550100')).toBe('(781) 555-0100')
    expect(masked('(###) ###-####', '(781) 555-0100')).toBe('(781) 555-0100')
    expect(masked('(###) ###-####', '781')).toBe('(781')
  })
  it('leaves out what does not fit', () => {
    expect(masked('#####-###', '12a345b678')).toBe('12345-678')
    expect(masked('#####-###', '1234567890')).toBe('12345-678')
  })
  it('takes letters', () => {
    expect(masked('AAA-#A##', 'abc1d23')).toBe('abc-1d23')
    expect(masked('A#A #A#', 'k1a0b1')).toBe('k1a 0b1')
  })
  it('an empty value is empty', () => {
    expect(masked('###', '')).toBe('')
  })
})
