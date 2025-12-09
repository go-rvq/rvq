import { useI18n as vuseI18n, UseI18nOptions } from 'vue-i18n'
import pt from './pt'
import en from './en'

export function useI18n<Options extends UseI18nOptions = UseI18nOptions>(options?: Options) {
  const i18n = vuseI18n({
    messages: {
      en,
      en_US:en,
      pt,
      pt_BR:pt,
    }
  })
  return i18n
}
