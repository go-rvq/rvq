import { Message } from '@/lib/Messages'

export const makeMessagesAttribute = () => ({
  messages: {
    renderHTML: () => ({}),
    parseHTML(element:HTMLElement) {
      const v = element.getAttribute('data-messages')
      let r: Message[] = []
      if (v) {
        r = JSON.parse(v) as Message[]
      }
      return r
    }
  }
})

export const makeInnerTextAsValueAttribute = () => ({
  value: {
    default: '',
    renderHTML: () => {
      return {}
    },
    parseHTML: (element:HTMLElement) => {
      const v = element.innerText
      element.innerHTML = ''
      return v
    }
  }
})

export const makeAttributes = () => ({
  ...makeMessagesAttribute(),
  ...makeInnerTextAsValueAttribute()
})
