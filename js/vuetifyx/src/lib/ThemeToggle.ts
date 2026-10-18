// <vx-theme-toggle>: a button that turns the theme of the page light or dark,
// the choice remembered in the browser (storage-key) and taken again when the
// page opens.
import { computed, defineComponent, h, onMounted, resolveComponent } from 'vue'
import {useVxMessages} from '@/lib/locale'
import { useTheme } from 'vuetify'

const read = (key: string): string | null => {
  try {
    return localStorage.getItem(key)
  } catch {
    return null
  }
}
const write = (key: string, v: string) => {
  try {
    localStorage.setItem(key, v)
  } catch {
    /* private window */
  }
}

export default defineComponent({
  name: 'VxThemeToggle',
  props: {
    storageKey: { type: String, default: 'vx.theme' },
    // the titles of the button: to light, to dark
    lightTitle: { type: String, default: undefined },
    darkTitle: { type: String, default: undefined }
  },
  setup(props) {
    const theme = useTheme()
    // the titles: the props given, else the page's vue-i18n's (vx)
    const vx = useVxMessages()
    const dark = computed(() => theme.global.current.value.dark)
    const set = (name: string) => {
      if (theme.themes.value[name]) theme.global.name.value = name
    }
    onMounted(() => {
      const saved = read(props.storageKey)
      if (saved) set(saved)
    })
    return () =>
      h(resolveComponent('v-btn'), {
        icon: dark.value ? 'mdi-weather-sunny' : 'mdi-weather-night',
        variant: 'text',
        title: dark.value ? props.lightTitle ?? vx.value.themeLight : props.darkTitle ?? vx.value.themeDark,
        'data-theme-toggle': dark.value ? 'dark' : 'light',
        onClick: () => {
          const name = dark.value ? 'light' : 'dark'
          set(name)
          write(props.storageKey, name)
        }
      })
  }
})
